package column

// Int returns a column of type Int.
// TODO: Add big int, tiny int, etc.
func Int(name string) Column {
	return Column{
		Name:   name,
		Values: "integer",
	}
}
