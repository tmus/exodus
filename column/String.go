package column

// String creates a column with a type of String.
func String(name string, len int) Definition {
	meta := getBaseMeta()
	meta["length"] = len
	return Definition{
		Name:     name,
		Kind:     "string",
		Metadata: meta,
	}
}
