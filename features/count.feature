Feature: Checking number of rows

    Scenario: Check count on table
        Given I run "CREATE TABLE IF NOT EXISTS test (id int, blah varchar(255), foo varchar(255))"
        And I run "TRUNCATE test"
        And I run "INSERT INTO test (id,blah,foo) VALUES (1,'aa','bb')"
        When I have table "test"
        Then the number of rows should be greater than 0