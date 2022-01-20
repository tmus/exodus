package exodus

import (
	"fmt"
)

type Migrator struct {
	driver Driver
}

func NewMigrator(driver Driver) (*Migrator, error) {
	m := &Migrator{
		driver: driver,
	}

	return m, nil
}

// TODO replace dir with opts, which includes direction and fresh commands, maybe more.
func (m *Migrator) Run(dir string, migrations ...Migration) error {
	if err := m.driver.Init(); err != nil {
		return fmt.Errorf("cannot initialise migration driver: %w", err)
	}

	if err := m.driver.ProcessBatch(migrations); err != nil {
		return fmt.Errorf("unable to run migrations: %w", err)
	}

	return nil
}
