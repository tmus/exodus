package exodus

// ColumnType provides enum access to common database column
// types.
type ColumnType string

// Possible column types are defined below.
const (
	Int       ColumnType = "integer" // TODO: Allow different sizes of int
	String    ColumnType = "varchar"
	Text      ColumnType = "text"
	Binary    ColumnType = "binary"
	Boolean   ColumnType = "boolean"
	Char      ColumnType = "char"
	Date      ColumnType = "date"
	Time      ColumnType = "time"
	DateTime  ColumnType = "datetime"
	Timestamp ColumnType = "timestamp"
)
