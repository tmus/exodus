package exodus

// Migration is a complete, raw SQL command that can be ran
// against a database.
type Migration struct {
	SQL   string
	Batch int
}

func (m Migration) String() string {
	return string(m.SQL)
}

type MigrationInterface interface {
	Up() Migration
	Down() Migration
}
