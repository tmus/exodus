package column

import "fmt"

type Constraint struct {
	Name    string
	Columns []string
}

// UniqueSet ...
func UniqueSet(name string, columns ...string) *Constraint {
	return &Constraint{
		Name:    name,
		Columns: columns,
	}
}

func (us *Constraint) ToSQL() string {
	sql := fmt.Sprintf(
		"CONSTRAINT %s UNIQUE ( %s )",
		us.Name,
		us.getSet(),
	)

	return sql
}

func (us *Constraint) getSet() (inline string) {
	for i, column := range us.Columns {
		inline = inline + column
		if i < len(us.Columns)-1 {
			inline = inline + ","
		}
	}

	return
}
