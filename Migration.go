package exodus

type MigrationCommand string

// Migration does x.
type Migration interface {
	Up() MigrationCommand
	Down() MigrationCommand
}

type operation int

const (
	UNKNOWN_OPERATION operation = iota
	CREATE_TABLE
	DROP_TABLE
	RENAME_TABLE
)

type MigrationPayload struct {
	Operation operation
	Payload   interface{}
	Table     string
}

/*
needs more than this, as the different dialects need to be able to inspect the column
and do stuff with it.

So something like:

operation:

create table - will need to contain a payload of columns to create - column.Definition
drop table - easy - pass in string and turn into `drop table string;`
rename table - easy - pass in strings and turn into `rename table string1 to string2`


alter table
alter table add column
alter table drop column
*/
