package driver

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tmus/exodus"
	"github.com/tmus/exodus/column"
)

type mysqlDriver struct {
	db *sql.DB
}

func NewMysqlDriver(datasource string) (*mysqlDriver, error) {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, fmt.Errorf("unable to open source for mysql driver: %w", err)
	}

	return &mysqlDriver{
		db: db,
	}, nil
}

func (d mysqlDriver) Close() error {
	return d.db.Close()
}

func (d mysqlDriver) Process(migration exodus.Migration) error {
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

func (d mysqlDriver) RenameTable(payload *exodus.MigrationOperation) error {
	from := payload.Table()
	to := payload.Payload().(string)
	sql := fmt.Sprintf(renameTableSchema, from, to)
	fmt.Println(sql)

	if _, err := d.db.Exec(sql); err != nil {
		return fmt.Errorf("error renaming `%s` table to `%s`: %w", from, to, err)
	}

	return nil
}

func (d mysqlDriver) CreateTable(payload *exodus.MigrationOperation) error {
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
func (d mysqlDriver) makeColumn(c column.Definition) (string, error) {
	switch c.Kind {
	case "string":
		return d.makeStringColumn(c)
	case "boolean":
		return d.makeBooleanColumn(c)
	case "date":
		return d.makeDateColumn(c)
	case "int":
		return d.makeIntColumn(c)
	default:
		return "", fmt.Errorf("unable to make column `%s`. unknown kind `%s`", c.Name, c.Kind)
	}
}

func isNullable(c column.Definition) (bool, error) {
	nullable, ok := c.Metadata["nullable"].(bool)
	if !ok {
		return false, errors.New("nullable value is not bool")
	}

	return nullable, nil
}

func (d mysqlDriver) nullableModifier(value bool) string {
	if value {
		return ""
	}

	return " NOT NULL"
}

func (d mysqlDriver) makeIntColumn(c column.Definition) (string, error) {
	autoincrement, ok := c.Metadata["increments"].(bool)
	if !ok {
		return "", fmt.Errorf("cannot create column %s: increments value is not bool", c.Name)
	}

	primaryKey, ok := c.Metadata["primary_key"].(bool)
	if !ok {
		return "", fmt.Errorf("cannot create column %s: primary_key value is not bool", c.Name)
	}

	nullable, ok := c.Metadata["nullable"].(bool)
	if !ok {
		return "", fmt.Errorf("cannot create column %s: nullable value is not bool", c.Name)
	}

	return fmt.Sprintf("%s INT%s%s%s", c.Name, d.autoincrementModifier(autoincrement), d.primaryKeyModifier(primaryKey), d.nullableModifier(nullable)), nil
}

func (d mysqlDriver) autoincrementModifier(value bool) string {
	if value {
		return " AUTO_INCREMENT"
	}

	return ""
}

func (d mysqlDriver) primaryKeyModifier(value bool) string {
	if value {
		return " PRIMARY KEY"
	}

	return ""
}

func (d mysqlDriver) makeStringColumn(c column.Definition) (string, error) {
	len, ok := c.Metadata["length"].(int)
	if !ok {
		return "", fmt.Errorf("invalid length (%v) for column %s", c.Metadata["length"], c.Name)
	}

	nullable, err := isNullable(c)
	if err != nil {
		return "", fmt.Errorf("cannot add column `%s`: %w", c.Name, err)
	}

	return fmt.Sprintf("%s VARCHAR(%d)%s", c.Name, len, d.nullableModifier(nullable)), nil
}

func (d mysqlDriver) makeBooleanColumn(c column.Definition) (string, error) {
	nullable, err := isNullable(c)
	if err != nil {
		return "", fmt.Errorf("cannot add column `%s`: %w", c.Name, err)
	}

	return fmt.Sprintf("%s BOOLEAN%s", c.Name, d.nullableModifier(nullable)), nil
}

func (d mysqlDriver) makeDateColumn(c column.Definition) (string, error) {
	nullable, err := isNullable(c)
	if err != nil {
		return "", fmt.Errorf("cannot add column `%s`: %w", c.Name, err)
	}

	return fmt.Sprintf("%s DATE%s", c.Name, d.nullableModifier(nullable)), nil
}

func (d mysqlDriver) migrationSchema() string {
	return `CREATE TABLE migrations (
	id integer not null primary key auto_increment,
	migration varchar(255) not null,
	batch integer not null
);`
}

func (d mysqlDriver) CreateMigrationsTable() error {
	if _, err := d.db.Exec(d.migrationSchema()); err != nil {
		return fmt.Errorf("error creating `migrations` table: %w", err)
	}

	return nil
}

func (d mysqlDriver) GetDB() *sql.DB {
	return d.db
}

func (d mysqlDriver) Fresh() error {
	dropSql := `SELECT table_name FROM information_schema.tables WHERE table_schema = SCHEMA()`

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

func (d mysqlDriver) TableExists(table string) (bool, error) {
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT 1", table)
	if _, err := d.db.Exec(sql); err != nil {
		if strings.Contains(err.Error(), "no such table") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
