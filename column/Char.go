package column

// Char returns a char column.
func Char(name string, len int) *Column {
	col := &Column{
		Name:     name,
		datatype: "char",
		meta:     make(map[string]interface{}),
	}

	col.Length(len)
	return col
}
