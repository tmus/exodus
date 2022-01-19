package exodus

func (p *MigrationPayload) Rename(from string, to string) *MigrationPayload {
	op := &MigrationOperation{
		operation: RENAME_TABLE,
		table:     from,
		payload:   to,
	}

	p.ops = append(p.ops, op)
	return p
}
