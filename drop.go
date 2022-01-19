package exodus

// Drop generates an SQL command to drop the given table.
func (p *MigrationPayload) Drop(table string) *MigrationPayload {
	op := &MigrationOperation{
		operation: DROP_TABLE,
		table:     table,
	}

	p.ops = append(p.ops, op)
	return p
}
