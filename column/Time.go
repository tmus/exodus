package column

// Time returns a time column.
func Time(name string) Column {
	return Column{
		Name:   name,
		Values: "time",
	}
}
