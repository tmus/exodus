package column

// DateTime returns a DateTime column.
func DateTime(name string) *Column {
	return &Column{
		Name:     name,
		datatype: "datetime",
		meta:     make(map[string]interface{}),
	}
}
