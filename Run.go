package exodus

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
)

// Run uses the passed migration to update the passed database.
func Run(migration Migration, database *sql.DB) error {

	fmt.Println("DEBUG: Running SQL:", migration.String())

	if err := prepMigrations(database); err != nil {
		log.Fatalln(err)
	}

	_, err := database.Exec(migration.String())
	if err != nil {
		return err
	}

	return nil
}

// prepMigrations ensures that the migrations are ready to
// be ran.
func prepMigrations(database *sql.DB) error {
	if !TableExists("migrations", database) {
		if err := createMigrationsTable(database); err != nil {
			return err
		}
	}

	return nil
}

// Fresh drops all tables in the database.
func Fresh(database *sql.DB) {
	if err := dropAllTables(database); err != nil {
		log.Fatalln(err)
	}
}

// createMigrationsTable runs SQL to create a table to hold
// all of the migrations and the order that they were executed.
func createMigrationsTable(database *sql.DB) error {
	migrationSchema := fmt.Sprintf(
		"CREATE TABLE migrations ( %s, %s, %s )",
		"id integer not null primary key autoincrement",
		"migration varchar not null",
		"batch integer not null",
	)

	if _, err := database.Exec(migrationSchema); err != nil {
		return fmt.Errorf("error creating migrations table: %s", err)
	}

	return nil
}

// TableExists determines if a table exists on the database.
// TODO: Probably a better way of doing this.
func TableExists(table string, database *sql.DB) bool {
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT 1", table)
	if _, err := database.Exec(sql); err != nil {
		return false
	}

	return true
}

func getDriverName(driver driver.Driver) string {
	sqlDriverNamesByType := map[reflect.Type]string{}

	for _, driverName := range sql.Drivers() {
		// Tested empty string DSN with MySQL, PostgreSQL, and SQLite3 drivers.
		db, _ := sql.Open(driverName, "")

		if db != nil {
			driverType := reflect.TypeOf(db.Driver())
			sqlDriverNamesByType[driverType] = driverName
		}
	}

	driverType := reflect.TypeOf(driver)
	if driverName, found := sqlDriverNamesByType[driverType]; found {
		return driverName
	}

	return ""
}

func getDropSQLForDriver(d driver.Driver) (string, error) {
	driver := getDriverName(d)

	// TODO: Add more driver support.
	if driver == "sqlite3" {
		return "SELECT name FROM sqlite_master WHERE type='table'", nil
	}

	return "", fmt.Errorf("`%s` driver is not yet supported", driver)
}

// dropAllTables grabs the tables from the database and drops
// them in turn, stopping if there is an error.
// TODO: Wrap this in a transaction, so it is cancelled if any
// of the drops fail?
func dropAllTables(database *sql.DB) error {
	// Get the SQL command to drop all tables for the current
	// SQL driver provided in the database connection.
	dropSQL, err := getDropSQLForDriver(database.Driver())
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
