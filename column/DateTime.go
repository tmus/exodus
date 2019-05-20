package column

// DateTime returns a DateTime column.
func DateTime(name string) Column {
	return Column{
		Name:   name,
		Values: "datetime",
	}
}
