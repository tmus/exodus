package exodus

import "fmt"

// Drop generates an SQL command to drop the given table.
func Drop(table string) MigrationCommand {
	return MigrationCommand(fmt.Sprintf("drop table %s", table))
}
