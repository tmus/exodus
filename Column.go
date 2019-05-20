package exodus

import "fmt"

// Columnable elements can be converted to an SQL representation.
// Using this interface rather than a concrete struct allows for
// users to pass in their own implementations - an `IDColumn()`
// method, for example.
type Columnable interface {
	toSQL() string
}

// Column defines a single column on a database table.
type Column struct {
	Name       string
	Values     ColumnType
	increments bool
	primaryKey bool
	nullable   bool
	length     int
}

func (c Column) toSQL() string {
	// TODO: Make this better. Really, all the "meta" info
	// should be put into a slice and iterated through and
	// appended to the "core" column data - the name and type.
	sql := fmt.Sprintf("%s %s", c.Name, c.Values)
	// TODO: Tidy this up.
	if (c.Values == String || c.Values == Char) && c.length != 0 {
		sql = sql + fmt.Sprintf("(%d)", c.length)
	}
	if c.nullable == false {
		sql = sql + " not null"
	}

	if c.primaryKey == true {
		sql = sql + " primary key"
	}

	if c.increments == true {
		sql = sql + " autoincrement"
	}

	return sql
}

// Is allows the user to define a column type in a fluent syntax.
func (c Column) Is(t ColumnType) Column {
	c.Values = t
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

// StringColumn creates a column with a type of String.
func StringColumn(name string, len int) Column {
	return Column{
		Name:   name,
		Values: String,
		length: len,
	}
}

func CharColumn(name string, len int) Column {
	return Column{
		Name:   name,
		Values: Char,
		length: len,
	}
}
