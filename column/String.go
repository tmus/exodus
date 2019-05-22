package column

// String creates a column with a type of String.
func String(name string, len int) *Column {
	col := &Column{
		Name:     name,
		datatype: "varchar",
	}

	col.Length(len)
	return col
}
