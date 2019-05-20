package exodus

import "fmt"

// Create generates a Migration that contains SQL to create
// the schema against the table name provided.
func Create(table string, schema Schema) Migration {
	sql := parseColumns(loadColumnSQL(schema))

	return Migration{
		SQL: fmt.Sprintf("CREATE TABLE %s ( %s );", table, sql),
	}
}

// loadColumnSQL iterates through the Columns defined in the
// Schema and calls the toSQL command on them. The resulting
// SQL for each column is stored in a slice of strings and
// returned.
func loadColumnSQL(schema Schema) (commands []string) {
	for _, col := range schema {
		commands = append(commands, col.toSQL())
	}

	return
}

// parseColumns implodes the slice of SQL commands into a
// single string containing all the column definitions.
func parseColumns(sql []string) (formatted string) {
	for i, command := range sql {
		formatted = formatted + command
		// If it is not the last column, add ", " after the
		// formatted string.
		if i != len(sql)-1 {
			formatted = formatted + ", "
		}
	}

	return
}
