package column

// Binary creates a binary column.
func Binary(name string) *Column {
	return &Column{
		Name:     name,
		datatype: "binary",
		meta:     make(map[string]interface{}),
	}
}
