package exodus

import (
	"fmt"
	"log"
	"reflect"
)

type Migrator struct {
	driver driver
}

func NewMigrator(driver driver) (*Migrator, error) {
	m := &Migrator{
		driver: driver,
	}

	return m, nil
}

// Fresh drops all tables in the database.
func (m *Migrator) Fresh() error {
	return m.driver.Fresh()
}

// nextBatchNumber retreives the highest batch number from the
// migrations table and increments it by one.
func (m *Migrator) nextBatchNumber() int {
	return m.lastBatchNumber() + 1
}

// lastBatchNumber retrieves the number of the last batch ran
// on the migrations table.
func (m *Migrator) lastBatchNumber() int {
	r := m.driver.GetDB().QueryRow("SELECT MAX(batch) FROM migrations")
	var num int
	r.Scan(&num)
	return num
}

func (m *Migrator) Run(migrations ...Migration) error {
	m.Fresh() // temp

	if err := m.verifyMigrationsTable(); err != nil {
		return fmt.Errorf("cannot verify state of `migrations` table: %w", err)
	}

	batch := m.nextBatchNumber()
	for _, migration := range migrations {
		if err := m.driver.Process(migration.Up()); err != nil {
			return fmt.Errorf("unable to execute SQL: %w", err)
		}

		m.addBatchToMigrationsTable(migration, batch)
	}

	return nil
}

func (m *Migrator) addBatchToMigrationsTable(migration Migration, batch int) {
	stmt, err := m.driver.GetDB().Prepare("INSERT INTO migrations (migration, batch) VALUES ( ?, ? )")
	if err != nil {
		log.Fatalln("Cannot create `migrations` batch statement. ")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(reflect.TypeOf(migration).String(), batch); err != nil {
		log.Fatalln(err)
	}
}

func (m *Migrator) verifyMigrationsTable() error {
	if ok, _ := m.driver.TableExists("migrations"); ok {
		return nil
	}

	if err := m.driver.CreateMigrationsTable(); err != nil {
		return err
	}

	return nil
}
