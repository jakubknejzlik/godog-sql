package database

import (
	"fmt"
	"net/url"
	"strings"

	"database/sql"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewDBWithString ...
func NewDBWithString(urlString string) *sql.DB {
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	urlString = strings.Replace(urlString, u.Scheme+"://", "", 1)

	db, err := sql.Open(u.Scheme, urlString)
	if err != nil {
		panic(err)
	}
	// db.LogMode(true)
	return db
}

func RowsToTable(rows *sql.Rows) (res [][]string, count int, err error) {
	cols, err := rows.Columns()
	if err != nil {
		return
	}

	res = append(res, cols)
	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}

		res = append(res, result)
		count++
	}

	return
}
