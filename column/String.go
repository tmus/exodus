package column

// String creates a column with a type of String.
func String(name string, len int) Column {
	return Column{
		Name:   name,
		Values: "varchar",
	}.Length(len)
}
