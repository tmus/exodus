package exodus

// ColumnType provides enum access to common database column
// types.
type ColumnType string

// Possible column types are defined below.
const (
	Int       ColumnType = "integer" // TODO: Allow different sizes of int
	String    ColumnType = "varchar" // TODO: Make length customisable
	Text      ColumnType = "text"
	Binary    ColumnType = "binary"
	Boolean   ColumnType = "boolean"
	Char      ColumnType = "char" // TODO: make length customisable
	Date      ColumnType = "date"
	Time      ColumnType = "time"
	DateTime  ColumnType = "datetime"
	Timestamp ColumnType = "timestamp"
)
