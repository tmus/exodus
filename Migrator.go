package exodus

import (
	"database/sql"
	"fmt"
	"log"
)

// Migrator is responsible for receiving the incoming migrations
// and running their SQL.
type Migrator struct {
	DB *sql.DB
}

// nextBatchNumber retreives the highest batch number from the
// migrations table and increments it by one.
func (m Migrator) nextBatchNumber() int {
	return m.lastBatchNumber() + 1
}

// lastBatchNumber retrieves the number of the last batch ran
// on the migrations table.
func (m Migrator) lastBatchNumber() int {
	r := m.DB.QueryRow("SELECT MAX(batch) FROM migrations")
	var num int
	r.Scan(&num)
	return num
}

// Run uses the passed migration to update the passed database.
func (m Migrator) Run(migration MigrationInterface) error {
	if err := m.prepMigrations(); err != nil {
		log.Fatalln(err)
	}

	_, err := m.DB.Exec(migration.Up().String())
	if err != nil {
		return err
	}

	// TODO: Make this actually work.
	stmt, err := m.DB.Prepare("INSERT INTO migrations (migration, batch) VALUES ( ?, ? )")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(migration.Up().SQL, m.nextBatchNumber()); err != nil {
		return err
	}

	return nil
}

// prepMigrations ensures that the migrations are ready to
// be ran.
func (m Migrator) prepMigrations() error {
	if !TableExists("migrations", m.DB) {
		if err := m.createMigrationsTable(); err != nil {
			return err
		}
	}

	return nil
}

// createMigrationsTable runs SQL to create a table to hold
// all of the migrations and the order that they were executed.
func (m Migrator) createMigrationsTable() error {
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
