package exodus

import (
	"fmt"
)

// Migration is a fully-formed SQL command that can be ran against
// a database connection.
type Migration string

type BaseMigration struct{}

func (BaseMigration) Up() Migration {
	fmt.Println("WARNING: Migration empty")
	return ""
}

func (BaseMigration) Down() Migration {
	return ""
}

// MigrationInterface ...
type MigrationInterface interface {
	Up() Migration
	Down() Migration
}
