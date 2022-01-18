package exodus

import "fmt"

func Rename(from string, to string) Migration {
	return Migration(fmt.Sprintf("rename table %s to %s", from, to))
}
