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
