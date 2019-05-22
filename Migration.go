package exodus

import (
	"fmt"
	"strings"
)

// Migration is a fully-formed SQL command that can be ran against
// a database connection.
type Migration string

// MigrationInterface ...
type MigrationInterface interface {
	Up() Migration
	Down() Migration
}

// Create generates an SQL command to create a table using the
// schema provided.
func Create(table string, schema Schema) Migration {
	sql := strings.Join(loadColumnSQL(schema), ", ")

	return Migration(fmt.Sprintf("CREATE TABLE %s ( %s );", table, sql))
}

// Drop generates an SQL command to drop the given table.
func Drop(table string) Migration {
	return Migration(fmt.Sprintf("DROP TABLE %s", table))
}

// loadColumnSQL iterates through the Columns defined in the
// Schema and calls the toSQL command on them. The resulting
// SQL for each column is stored in a slice of strings and
// returned.
func loadColumnSQL(schema Schema) (commands []string) {
	for _, col := range schema {
		commands = append(commands, col.ToSQL())
	}

	return
}
