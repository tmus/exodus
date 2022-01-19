package column

// String creates a column with a type of String.
func String(name string, len int) Definition {
	meta := make(map[string]interface{}, 1)
	meta["length"] = len
	return Definition{
		Name:     name,
		Kind:     "string",
		Metadata: meta,
	}
}
