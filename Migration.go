package exodus

type MigrationCommand string

// Migration does x.
type Migration interface {
	Up() MigrationCommand
	Down() MigrationCommand
}
