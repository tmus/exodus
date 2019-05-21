package column

// Boolean returns a bool column.
func Boolean(name string) *Column {
	return &Column{
		Name:     name,
		datatype: "boolean",
	}
}
