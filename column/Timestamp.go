package column

// Timestamp returns a timestamp column.
func Timestamp(name string) *Column {
	return &Column{
		Name:     name,
		datatype: "timestamp",
	}
}
