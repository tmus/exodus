package exodus

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
)

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

// Make Migrator struct with a pointer to a sql DB. This can then be used instead of passing the db in every time:

// Migator.Run(migration)
// Migrator.NextBatchNumber() //

// Migration is a struct, should contain a sql.DB so we don't need to pass
