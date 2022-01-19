package column

// BaseColumn defines common properties that apply to all columns.
type Definition struct {
	Name     string
	Kind     string
	Metadata map[string]interface{}
}

func (d Definition) Nullable() Definition {
	d.Metadata["nullable"] = true
	return d
}

func (d Definition) NotNullable() Definition {
	d.Metadata["nullable"] = false
	return d
}

func getBaseMeta() map[string]interface{} {
	meta := make(map[string]interface{})
	meta["nullable"] = false
	meta["increments"] = false
	meta["unsigned"] = false
	meta["primary_key"] = false

	return meta
}
