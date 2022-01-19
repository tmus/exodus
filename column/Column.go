package column

// BaseColumn defines common properties that apply to all columns.
type Definition struct {
	Name     string
	Kind     string
	Metadata map[string]interface{}
}
