package column

import "fmt"

// Columnable elements can be converted to an SQL representation.
// Using this interface rather than a concrete struct allows for
// users to pass in their own implementations - an `IDColumn()`
// method, for example.
type Columnable interface {
	ToSQL() string
}

// Column defines a single column on a database table.
type Column struct {
	Name         string
	datatype     string
	increments   bool
	primaryKey   bool
	nullable     bool
	unique       bool
	length       int
	defaultValue string
}

// ToSQL converts the column struct to an SQL command.
func (c Column) ToSQL() string {
	// TODO: Make this better. Really, all the "meta" info
	// should be put into a slice and iterated through and
	// appended to the "core" column data - the name and type.
	sql := fmt.Sprintf("%s %s", c.Name, c.datatype)
	// TODO: Tidy this up.
	if (c.datatype == "string" || c.datatype == "char") && c.length != 0 {
		sql = sql + fmt.Sprintf("(%d)", c.length)
	}
	if c.nullable == false {
		sql = sql + " not null"
	}

	if c.unique == true {
		sql = sql + " unique"
	}

	if c.primaryKey == true {
		sql = sql + " primary key"
	}

	if c.increments == true {
		sql = sql + " autoincrement"
	}

	if c.defaultValue != "" {
		sql = sql + " default " + c.defaultValue
	}

	return sql
}

// Unique makes a columns value unique in the table.
func (c Column) Unique() Column {
	c.unique = true
	return c
}

// Default sets the default value for a Column.
// TODO: value should be an interface to reflect.
func (c Column) Default(value string) Column {
	c.defaultValue = value
	return c
}

// Increments determines if the column auto-increments or not.
// The default value is false.
func (c Column) Increments() Column {
	c.increments = true
	return c
}

// PrimaryKey determines if a column is the table's primary key.
// The default value for a column is false.
func (c Column) PrimaryKey() Column {
	c.primaryKey = true
	return c
}

// NotNullable determines if a column's value can be null.
func (c Column) NotNullable() Column {
	c.nullable = false
	return c
}

// Nullable determines if a column's value can be null.
func (c Column) Nullable() Column {
	c.nullable = true
	return c
}

// Length adds a length constraint to applicable columns.
// TODO: Should this throw an error on columns that can't
// have a length modifier? Like TEXT?
func (c *Column) Length(len int) {
	c.length = len
}