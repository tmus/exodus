package exodus

func (p *MigrationPayload) Raw(sql string) *MigrationPayload {
	op := &MigrationOperation{
		operation: RAW_SQL,
		payload:   sql,
	}

	p.ops = append(p.ops, op)
	return p
}
