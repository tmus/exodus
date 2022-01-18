package column

import "fmt"

// String creates a column with a type of String.
func String(name string, len int) *stringColumn {
	base := &BaseColumn{
		name: name,
	}

	return &stringColumn{
		length:     len,
		BaseColumn: base,
	}
}

type stringColumn struct {
	*BaseColumn
	length int
}

func (c stringColumn) ToSQL() string {
	return fmt.Sprintf("%s VARCHAR(%d)", c.name, c.length)
}
