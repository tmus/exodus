package column

// Boolean returns a bool column.
func Boolean(name string) *booleanColumn {
	base := BaseColumn{
		name: name,
	}

	return &booleanColumn{
		BaseColumn: base,
	}
}

type booleanColumn struct {
	BaseColumn
}

func (c booleanColumn) ToSQL() string {
	return c.name + " BOOLEAN"
}
