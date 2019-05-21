package column

// Char returns a char column.
func Char(name string, len int) *Column {
	col := &Column{
		Name:     name,
		datatype: "char",
	}

	col.Length(len)
	return col
}
