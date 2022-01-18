package exodus

// Drop generates an SQL command to drop the given table.
func Drop(table string) MigrationPayload {
	return MigrationPayload{
		Operation: DROP_TABLE,
		Table:     table,
	}
}
