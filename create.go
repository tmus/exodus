package exodus

import (
	"fmt"
	"strings"

	"github.com/tmus/exodus/column"
)

// Create generates an SQL command to create a table using the
// schema provided.
func Create(table string, columns []column.Definition) Migration {
	var cols []string

	for _, col := range columns {
		cols = append(cols, col.ToSQL())
	}

	sql := strings.Join(cols, ",\n	")

	return Migration(fmt.Sprintf("create table %s (\n	%s\n);", table, sql))
}
