package exodus

import (
	"github.com/tmus/exodus/column"
)

// Create generates an SQL command to create a table using the
// schema provided.
func Create(table string, columns []column.Definition) MigrationPayload {
	return MigrationPayload{
		Operation: CREATE_TABLE,
		Payload:   columns,
		Table:     table,
	}
}
