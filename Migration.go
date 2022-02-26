package exodus

type MigrationCommand string

// Migration does x.
type Migration interface {
	Up(migrate *MigrationPayload)
	Down(migrate *MigrationPayload)
}

type operation int

const (
	UNKNOWN_OPERATION operation = iota
	CREATE_TABLE
	DROP_TABLE
	RENAME_TABLE
	RAW_SQL
)

type MigrationOperation struct {
	operation operation
	payload   interface{}
	table     string
}

func (o MigrationOperation) Operation() operation {
	return o.operation
}

func (o MigrationOperation) Payload() interface{} {
	return o.payload
}

func (o MigrationOperation) Table() string {
	return o.table
}

type MigrationPayload struct {
	ops []*MigrationOperation
}

func (p *MigrationPayload) Operations() []*MigrationOperation {
	return p.ops
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
