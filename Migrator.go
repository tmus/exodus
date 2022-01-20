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

func (m *Migrator) Run(args []string, migrations ...Migration) error {
	opts := gatherOptions(args)

	if err := m.driver.Init(); err != nil {
		return fmt.Errorf("cannot initialise migration driver: %w", err)
	}

	if err := m.driver.Run(opts, migrations); err != nil {
		return fmt.Errorf("unable to run migrations: %w", err)
	}

	return nil
}

func gatherOptions(args []string) Options {
	dir := Unknown
	if len(args) > 0 {
		dir = directionFromString(args[1])
	}

	return Options{
		direction: dir,
	}
}
