package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tmus/exodus"
	"github.com/tmus/exodus/column"
)

type mysqlDriver struct {
	db *sql.DB
}

func NewDriver(datasource string) (*mysqlDriver, error) {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, fmt.Errorf("unable to open source for mysql driver: %w", err)
	}

	return &mysqlDriver{
		db: db,
	}, nil
}

func (d mysqlDriver) Init() error {
	// if err := d.fresh(); err != nil {
	// 	return err
	// }

	if ok, _ := d.tableExists("migrations"); !ok {
		if err := d.createMigrationsTable(); err != nil {
			return err
		}
	}

	return nil
}

func (d mysqlDriver) Close() error {
	return d.db.Close()
}

func (d mysqlDriver) ProcessBatch(migrations []exodus.Migration) error {
	previous, err := d.getRan()
	if err != nil {
		return fmt.Errorf("unable to process migrations: %w", err)
	}
	migrations = filterPendingMigrations(migrations, previous)

	if len(migrations) == 0 {
		// TODO: Add logger to the driver
		fmt.Println("Nothing to migrate.")
		return nil
	}

	for _, migration := range migrations {
		if err := d.process(migration); err != nil {
			return fmt.Errorf("unable to execute SQL: %w", err)
		}
	}

	d.addBatchToMigrationsTable(migrations, d.nextBatchNumber())

	return nil
}

func (d mysqlDriver) getRan() ([]string, error) {
	rows, err := d.db.Query("SELECT migration FROM migrations")
	if err != nil {
		return []string{}, fmt.Errorf("unable to get previous migrations from database: %w", err)
	}

	var ran []string
	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return []string{}, fmt.Errorf("unable to get previous migrations from database: %w", err)
		}
		ran = append(ran, migration)
	}
	if err := rows.Err(); err != nil {
		return []string{}, fmt.Errorf("unable to get previous migrations from database: %w", err)
	}

	return ran, nil
}

func (d mysqlDriver) process(migration exodus.Migration) error {
	builder := &exodus.MigrationPayload{}
	migration.Up(builder)

	for _, p := range builder.Operations() {
		switch p.Operation() {
		case exodus.CREATE_TABLE:
			if err := d.createTable(p); err != nil {
				return err
			}
		case exodus.RENAME_TABLE:
			if err := d.renameTable(p); err != nil {
				return err
			}
		default:
			return errors.New("operation not supported")
		}
	}

	return nil
}

func (d mysqlDriver) addBatchToMigrationsTable(migrations []exodus.Migration, batch int) error {
	stmt, err := d.db.Prepare("INSERT INTO migrations (migration, batch) VALUES ( ?, ? )")
	if err != nil {
		log.Fatalln("Cannot create `migrations` batch statement. ")
	}
	defer stmt.Close()

	for _, migration := range migrations {
		if _, err = stmt.Exec(reflect.TypeOf(migration).String(), batch); err != nil {
			return err
		}
	}

	return nil
}

// nextBatchNumber retreives the highest batch number from the
// migrations table and increments it by one.
func (d mysqlDriver) nextBatchNumber() int {
	return d.lastBatchNumber() + 1
}

// lastBatchNumber retrieves the number of the last batch ran
// on the migrations table.
func (d mysqlDriver) lastBatchNumber() int {
	r := d.db.QueryRow("SELECT MAX(batch) FROM migrations")
	var num int
	r.Scan(&num)
	return num
}

func (d mysqlDriver) renameTable(payload *exodus.MigrationOperation) error {
	from := payload.Table()
	to := payload.Payload().(string)
	sql := fmt.Sprintf(renameTableSchema, from, to)
	fmt.Println(sql)

	if _, err := d.db.Exec(sql); err != nil {
		return fmt.Errorf("error renaming `%s` table to `%s`: %w", from, to, err)
	}

	return nil
}

func (d mysqlDriver) createTable(payload *exodus.MigrationOperation) error {
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

func (d mysqlDriver) createMigrationsTable() error {
	if _, err := d.db.Exec(migrationSchema); err != nil {
		return fmt.Errorf("error creating `migrations` table: %w", err)
	}

	return nil
}

func (d mysqlDriver) fresh() error {
	rows, err := d.db.Query(dropSchema)
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

func (d mysqlDriver) tableExists(table string) (bool, error) {
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT 1", table)
	if _, err := d.db.Exec(sql); err != nil {
		if strings.Contains(err.Error(), "no such table") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func filterPendingMigrations(migrations []exodus.Migration, existing []string) []exodus.Migration {
	var response []exodus.Migration

	for _, migration := range migrations {
		if !exists(migration, existing) {
			response = append(response, migration)
		}
	}

	return response
}

func exists(migration exodus.Migration, existing []string) bool {
	for _, ex := range existing {
		if reflect.TypeOf(migration).String() == ex {
			return true
		}
	}

	return false
}
