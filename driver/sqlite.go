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

const createTableSchema = `CREATE TABLE %s (
	%s
);`

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

func (d sqliteDriver) Process(migration exodus.Migration) error {
	builder := &exodus.MigrationPayload{}
	migration.Up(builder)

	for _, p := range builder.Operations() {
		switch p.Operation() {
		case exodus.CREATE_TABLE:
			if err := d.CreateTable(p); err != nil {
				return err
			}
		case exodus.RENAME_TABLE:
			if err := d.RenameTable(p); err != nil {
				return err
			}
		default:
			return errors.New("operation not supported")
		}
	}

	return nil
}

func (d sqliteDriver) RenameTable(payload *exodus.MigrationOperation) error {
	from := payload.Table()
	to := payload.Payload().(string)
	sql := fmt.Sprintf(renameTableSchema, from, to)
	fmt.Println(sql)

	if _, err := d.db.Exec(sql); err != nil {
		return fmt.Errorf("error renaming `%s` table to `%s`: %w", from, to, err)
	}

	return nil
}

func (d sqliteDriver) CreateTable(payload *exodus.MigrationOperation) error {
	var cols []string

	columns, ok := payload.Payload().([]column.Definition)
	if !ok {
		return errors.New("incorrect payload creating a table")
	}

	for _, col := range columns {
		def, err := d.makeColumn(col)
		if err != nil {
			return fmt.Errorf("unable to create table: %w", err)
		}
		cols = append(cols, def)
	}

	colSql := strings.Join(cols, ",\n	")

	sql := fmt.Sprintf(createTableSchema, payload.Table(), colSql)
	fmt.Println(sql)

	if _, err := d.db.Exec(sql); err != nil {
		return fmt.Errorf("error creating `%s` table: %w", payload.Table(), err)
	}

	return nil
}

// makeColumn takes a column.Definition and turns it into a string representation
// of that column, in a format that is understood by the driver.
func (d sqliteDriver) makeColumn(c column.Definition) (string, error) {
	switch c.Kind {
	case "string":
		return d.makeStringColumn(c)
	case "boolean":
		return d.makeBooleanColumn(c)
	default:
		return "", fmt.Errorf("unable to make column `%s`. unknown kind `%s`", c.Name, c.Kind)
	}
}

func (d sqliteDriver) makeStringColumn(c column.Definition) (string, error) {
	len, ok := c.Metadata["length"].(int)
	if !ok {
		return "", fmt.Errorf("invalid length (%v) for column %s", c.Metadata["length"], c.Name)
	}

	return fmt.Sprintf("%s VARCHAR(%d)", c.Name, len), nil
}

func (d sqliteDriver) makeBooleanColumn(c column.Definition) (string, error) {
	return fmt.Sprintf("%s BOOLEAN", c.Name), nil
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
