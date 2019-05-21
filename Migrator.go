package exodus

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

// supportedDrivers lists the drivers that currently work with
// the migration framework.
var supportedDrivers = []string{
	"sqlite3",
}

// Migrator is responsible for receiving the incoming migrations
// and running their SQL.
type Migrator struct {
	DB    *sql.DB
	Batch int
}

// NewMigrator creates a new instance of a migrator.
func NewMigrator(db *sql.DB) (*Migrator, error) {
	m := Migrator{
		DB: db,
	}

	if !m.driverIsSupported(m.getDriverName()) {
		return nil, fmt.Errorf("the %s driver is currently unsupported", m.getDriverName())
	}

	return &m, nil
}

// TableExists determines if a table exists on the database.
// TODO: Probably a better way of doing this.
func (m *Migrator) TableExists(table string, database *sql.DB) bool {
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT 1", table)
	if _, err := database.Exec(sql); err != nil {
		return false
	}

	return true
}

// Fresh drops all tables in the database.
func (m *Migrator) Fresh(database *sql.DB) {
	if err := m.dropAllTables(database); err != nil {
		log.Fatalln(err)
	}
}

// dropAllTables grabs the tables from the database and drops
// them in turn, stopping if there is an error.
// TODO: Wrap this in a transaction, so it is cancelled if any
// of the drops fail?
func (m *Migrator) dropAllTables(database *sql.DB) error {
	// Get the SQL command to drop all tables for the current
	// SQL driver provided in the database connection.
	dropSQL, err := m.getDropSQLForDriver(m.getDriverName())
	if err != nil {
		// If support for the driver does not exist, log a
		// fatal error.
		log.Fatalln("Unable to drop tables:", err)
	}

	rows, err := database.Query(dropSQL)
	if err != nil {
		return err
	}
	defer rows.Close()

	// tables is the list of tables returned from the database.
	var tables []string

	// for each row returned, add the name of it to the
	// tables slice.
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
		if _, err := database.Exec("DROP TABLE IF EXISTS " + table); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) getDropSQLForDriver(d string) (string, error) {
	// TODO: Add more driver support.
	// Postgres? Then that'll do.
	if d == "sqlite3" {
		return "SELECT name FROM sqlite_master WHERE type='table'", nil
	}

	if d == "mysql" {
		return "SHOW FULL TABLES WHERE table_type = 'BASE TABLE'", nil
	}

	return "", fmt.Errorf("`%s` driver is not yet supported", d)
}

// nextBatchNumber retreives the highest batch number from the
// migrations table and increments it by one.
func (m *Migrator) nextBatchNumber() int {
	return m.lastBatchNumber() + 1
}

// lastBatchNumber retrieves the number of the last batch ran
// on the migrations table.
func (m *Migrator) lastBatchNumber() int {
	r := m.DB.QueryRow("SELECT MAX(batch) FROM migrations")
	var num int
	r.Scan(&num)
	return num
}

// TODO: Wrap this in a transaction and reverse it
func (m *Migrator) Run(migrations ...MigrationInterface) error {
	m.verifyMigrationsTable()

	batch := m.nextBatchNumber()

	for _, migration := range migrations {
		if _, err := m.DB.Exec(migration.Up().String()); err != nil {
			return err
		}

		m.addBatchToMigrationsTable(migration, batch)
	}

	return nil
}

func (m *Migrator) addBatchToMigrationsTable(migration MigrationInterface, batch int) {
	stmt, err := m.DB.Prepare("INSERT INTO migrations (migration, batch) VALUES ( ?, ? )")
	if err != nil {
		log.Fatalln("Cannot create `migrations` batch statement. ")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(reflect.TypeOf(migration).String(), batch); err != nil {
		log.Fatalln(err)
	}
}

// prepMigrations ensures that the migrations are ready to
// be ran.
func (m *Migrator) verifyMigrationsTable() {
	if !m.TableExists("migrations", m.DB) {
		if err := m.createMigrationsTable(); err != nil {
			log.Fatalln("Could not create `migrations` table: ", err)
		}
	}
}

func (m *Migrator) driverIsSupported(driver string) bool {
	for _, d := range supportedDrivers {
		if d == driver {
			return true
		}
	}

	return false
}

// getDriverName returns the name of the SQL driver currently
// associated with the Migrator.
func (m *Migrator) getDriverName() string {
	sqlDriverNamesByType := map[reflect.Type]string{}

	for _, driverName := range sql.Drivers() {
		// Tested empty string DSN with MySQL, PostgreSQL, and SQLite3 drivers.
		db, _ := sql.Open(driverName, "")

		if db != nil {
			driverType := reflect.TypeOf(db.Driver())
			sqlDriverNamesByType[driverType] = driverName
		}
	}

	driverType := reflect.TypeOf(m.DB.Driver())
	if driverName, found := sqlDriverNamesByType[driverType]; found {
		return driverName
	}

	return ""
}

// createMigrationsTable makes a table to hold migrations and
// the order that they were executed.
func (m *Migrator) createMigrationsTable() error {
	migrationSchema := fmt.Sprintf(
		"CREATE TABLE migrations ( %s, %s, %s )",
		"id integer not null primary key autoincrement",
		"migration varchar not null",
		"batch integer not null",
	)

	if _, err := m.DB.Exec(migrationSchema); err != nil {
		return fmt.Errorf("error creating migrations table: %s", err)
	}

	return nil
}
