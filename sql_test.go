package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-test/deep"
	"github.com/martinohmann/godog-helpers/datatable"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jakubknejzlik/godog-sql/database"
)

type sqlFeature struct {
	db    *sql.DB
	rows  *[][]string
	count *int
	table *string
}

func (s *sqlFeature) iSelect(query string) (err error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return
	}

	rs, count, err := database.RowsToTable(rows)
	if err != nil {
		return
	}
	s.rows = &rs
	s.count = &count

	return
}

func (s *sqlFeature) iRun(query string) (err error) {
	_, err = s.db.Query(query)

	s.rows = nil

	return
}

func (s *sqlFeature) theResponseShouldHaveRow(count int) (err error) {
	if count != *s.count {
		err = fmt.Errorf("Wrong number of rows returned, expected %d, returned %d", count, *s.count)
	}
	return
}

func (s *sqlFeature) theResponseShouldBe(expected *gherkin.DataTable) (err error) {
	options := &datatable.Options{
		RequiredFields: []string{},
		OptionalFields: []string{},
	}

	result, err := datatable.FromGherkinWithOptions(options, expected)
	if err != nil {
		return err
	}

	expectedTable := [][]string{result.Fields()}
	expectedTable = append(expectedTable, result.RowValues()...)

	if diff := deep.Equal(s.rows, &expectedTable); diff != nil {
		err = errors.New(strings.Join(diff, "\n"))
	}
	return
}

func (s *sqlFeature) iHaveTable(tablename string) (err error) {
	s.table = &tablename
	return
}

func (s *sqlFeature) theNumberOfRowsShouldBeGreaterThan(count int) (err error) {
	res, err := s.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM %s", *s.table))
	if err != nil {
		return
	}

	var actualCount int
	for res.Next() {
		err = res.Scan(&actualCount)
		if err != nil {
			return
		}
	}

	if actualCount <= count {
		err = fmt.Errorf("number of rows should be greater then %d, but is %d", count, actualCount)
	}
	return
}

func FeatureContext(s *godog.Suite) {
	feature := &sqlFeature{}
	s.Step(`^I select "([^"]*)"$`, feature.iSelect)
	s.Step(`^I run "([^"]*)"$`, feature.iRun)
	s.Step(`^the response should be:$`, feature.theResponseShouldBe)
	s.Step(`^the response should have (\d+) rows?$`, feature.theResponseShouldHaveRow)
	s.Step(`^I have table "([^"]*)"$`, feature.iHaveTable)
	s.Step(`^the number of rows should be greater than (\d+)$`, feature.theNumberOfRowsShouldBeGreaterThan)

	s.BeforeScenario(func(interface{}) {
		feature.db = database.NewDBWithString(os.Getenv("DATABASE_URL"))
		if err := feature.db.Ping(); err != nil {
			panic(err)
		}
	})

	s.AfterScenario(func(interface{}, error) {
		feature.db.Close()
	})
}
