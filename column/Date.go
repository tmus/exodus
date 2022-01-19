package column

// Date creates a column with the type `date`.
func Date(name string) Definition {
	return Definition{
		Name: name,
		Kind: "date",
	}
}
