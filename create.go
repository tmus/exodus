package exodus

import (
	"fmt"
	"strings"
)

// Create generates an SQL command to create a table using the
// schema provided.
func Create(table string, schema Schema) Migration {
	var cols []string

	for _, col := range schema {
		cols = append(cols, col.ToSQL())
	}

	sql := strings.Join(cols, ",")

	return Migration(fmt.Sprintf("create table %s (%s);", table, sql))
}
