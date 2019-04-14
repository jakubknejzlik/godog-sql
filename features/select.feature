Feature: It should be possible to run select query

    Scenario: Basic select
        When I select "SELECT 1+1 as value"
        Then the response should have 1 row
        Then the response should be:
            | value |
            | 2     |