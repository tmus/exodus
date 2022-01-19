package exodus

import (
	"github.com/tmus/exodus/column"
)

// Create generates an SQL command to create a table using the
// schema provided.
func (p *MigrationPayload) Create(table string, columns []column.Definition) *MigrationPayload {
	op := &MigrationOperation{
		operation: CREATE_TABLE,
		payload:   columns,
		table:     table,
	}

	p.ops = append(p.ops, op)
	return p
}
