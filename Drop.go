package exodus

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
)

// Drop generates an SQL command to drop a table.
func Drop(table string) Migration {
	return Migration{
		SQL: fmt.Sprintf("DROP TABLE %s", table),
	}
}

// Fresh drops all tables in the database.
func Fresh(database *sql.DB) {
	if err := dropAllTables(database); err != nil {
		log.Fatalln(err)
	}
}

func getDropSQLForDriver(d driver.Driver) (string, error) {
	driver := getDriverName(d)

	// TODO: Add more driver support.
	// Postgres? Then that'll do.
	if driver == "sqlite3" {
		return "SELECT name FROM sqlite_master WHERE type='table'", nil
	}

	if driver == "mysql" {
		return "SHOW FULL TABLES WHERE table_type = 'BASE TABLE'", nil
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
