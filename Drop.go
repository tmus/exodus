package exodus

import "fmt"

// Drop generates an SQL command to drop a table.
func Drop(table string) Migration {
	return Migration{
		SQL: fmt.Sprintf("DROP TABLE %s", table),
	}
}
