package column

// Text returns a column of type text.
func Text(name string) *Column {
	return &Column{
		Name:     name,
		datatype: "text",
		meta:     make(map[string]interface{}),
	}
}
