package driver

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tmus/exodus"
)

const migrationSchema = `CREATE TABLE migrations (
	id integer not null primary key autoincrement,
	migration varchar not null,
	batch integer not null
);`

const createTableSchema = `CREATE TABLE %s (id int);`

type sqliteDriver struct {
	db *sql.DB
}

func NewSqliteDriver(datasource string) (*sqliteDriver, error) {
	db, err := sql.Open("sqlite3", datasource)
	if err != nil {
		return nil, fmt.Errorf("unable to open source for sqlite3 driver: %w", err)
	}

	return &sqliteDriver{
		db: db,
	}, nil
}

func (d sqliteDriver) Close() error {
	return d.db.Close()
}

func (sqliteDriver) Ref() string {
	return "sqlite3"
}

func (d sqliteDriver) CreateTable(payload exodus.MigrationPayload) error {
	sql := fmt.Sprintf(createTableSchema, payload.Table)

	if _, err := d.db.Exec(sql); err != nil {
		return fmt.Errorf("error creating `%s` table: %w", payload.Table, err)
	}

	return nil
}

func (d sqliteDriver) CreateMigrationsTable() error {
	if _, err := d.db.Exec(migrationSchema); err != nil {
		return fmt.Errorf("error creating `migrations` table: %w", err)
	}

	return nil
}

func (d sqliteDriver) GetDB() *sql.DB {
	return d.db
}

func (d sqliteDriver) Fresh() error {
	dropSql := "SELECT name FROM sqlite_master WHERE type='table'"

	rows, err := d.db.Query(dropSql)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		if name == "sqlite_sequence" {
			continue
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, table := range tables {
		if _, err := d.db.Exec("DROP TABLE IF EXISTS " + table); err != nil {
			return err
		}
	}

	return nil
}

func (d sqliteDriver) TableExists(table string) (bool, error) {
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT 1", table)
	if _, err := d.db.Exec(sql); err != nil {
		if strings.Contains(err.Error(), "no such table") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
