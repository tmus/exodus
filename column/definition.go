package column

type Definition interface {
	ToSQL() string
}
