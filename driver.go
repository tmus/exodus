package exodus

type Driver interface {
	// Init allows the migration driver to do any one time set up when ran.
	Init() error
	// Run takes a slice of exodus.Migration and runs it against the
	// underlying data store defined in the driver.
	Run(options Options, payload []Migration) error
	// Close closes the underlying database instance.
	Close() error
}
