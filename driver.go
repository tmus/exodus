package exodus

import "database/sql"

type driver interface {
	Process(payload MigrationPayload) error
	CreateMigrationsTable() error
	TableExists(name string) (bool, error)
	GetDB() *sql.DB
	Fresh() error
}
