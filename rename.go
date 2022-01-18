package exodus

import "fmt"

func Rename(from string, to string) MigrationCommand {
	return MigrationCommand(fmt.Sprintf("rename table %s to %s", from, to))
}
