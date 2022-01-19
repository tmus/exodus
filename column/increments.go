package column

// Increments creates an auto-incrementing unsigned INT equivalent column as a primary key.
func Increments(name string) Definition {
	meta := getBaseMeta()
	meta["increments"] = true
	meta["unsigned"] = true
	meta["primary_key"] = true

	return Definition{
		Name:     name,
		Kind:     "int",
		Metadata: meta,
	}
}
