package column

func Int(name string) Definition {
	meta := getBaseMeta()

	return Definition{
		Name:     name,
		Kind:     "int",
		Metadata: meta,
	}
}
