package column

// Date returns a date column.
func Date(name string) *Column {
	return &Column{
		Name:     name,
		datatype: "date",
	}
}
