package column

// Binary creates a binary column.
func Binary(name string) Column {
	return Column{
		Name:   name,
		Values: "binary",
	}
}
