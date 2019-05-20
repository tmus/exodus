package column

// Char returns a char column.
func Char(name string, len int) Column {
	return Column{
		Name:   name,
		Values: "char",
	}.Length(len)
}
