package column

// Boolean returns a bool column.
func Boolean(name string) Definition {
	return Definition{
		Name:     name,
		Kind:     "boolean",
		Metadata: getBaseMeta(),
	}
}
