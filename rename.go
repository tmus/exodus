package exodus

func Rename(from string, to string) MigrationPayload {
	return MigrationPayload{
		Operation: RENAME_TABLE,
		Table:     from,
		Payload:   to,
	}
}
