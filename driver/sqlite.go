package driver

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tmus/exodus"
	"github.com/tmus/exodus/column"
)

const migrationSchema = `CREATE TABLE migrations (
	id integer not null primary key autoincrement,
	migration varchar not null,
	batch integer not null
);`

const createTableSchema = `CREATE TABLE %s (%s);`

const renameTableSchema = `ALTER TABLE %s RENAME TO %s;`

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

func (d sqliteDriver) Process(payload exodus.MigrationPayload) error {
	switch payload.Operation {
	case exodus.CREATE_TABLE:
		return d.CreateTable(payload)
	case exodus.RENAME_TABLE:
		return d.RenameTable(payload)
	default:
		return errors.New("operation not supported")
	}
}

func (d sqliteDriver) RenameTable(payload exodus.MigrationPayload) error {
	from := payload.Table
	to := payload.Payload.(string)
	sql := fmt.Sprintf(renameTableSchema, from, to)

	if _, err := d.db.Exec(sql); err != nil {
		return fmt.Errorf("error renaming `%s` table to `%s`: %w", from, to, err)
	}

	return nil
}

func (d sqliteDriver) CreateTable(payload exodus.MigrationPayload) error {
	var cols []string

	columns, ok := payload.Payload.([]column.Definition)
	if !ok {
		return errors.New("incorrect payload creating a table")
	}

	for _, col := range columns {
		cols = append(cols, col.ToSQL())
	}

	colSql := strings.Join(cols, ",\n	")

	sql := fmt.Sprintf(createTableSchema, payload.Table, colSql)

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
