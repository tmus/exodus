package column

import (
	"fmt"
	"strings"
)

// Columnable elements can be converted to an SQL representation.
// Using this interface rather than a concrete struct allows for
// users to pass in their own implementations - an `IDColumn()`
// method, for example.
type Columnable interface {
	ToSQL() string
}

// Column defines a single column on a database table.
type Column struct {
	Name     string
	datatype string
	meta     map[string]interface{}
}

// ToSQL converts the column struct to an SQL command.
func (c *Column) ToSQL() string {
	var metadata []string
	for md, val := range c.meta {
		if md == "length" {
			// skip length, that's defined on the datatype
			continue
		}
		if val == true {
			metadata = append(metadata, md)
		} else {
			metadata = append(metadata, fmt.Sprint(md, " ", val))
		}
	}

	sql := fmt.Sprintf(
		"%s %s",
		c.Name,
		c.getDatatype(),
	)

	// TODO: Make this a function.
	optional := strings.Join(metadata, " ")
	if optional != "" {
		sql = sql + " " + optional
	}

	return sql
}

// getDatatype returns the datatype for the column. If the datatype
// has a length property, it is foratted correctly.
func (c *Column) getDatatype() string {
	if val, exists := c.meta["length"]; exists {
		return fmt.Sprintf("%s(%d)", c.datatype, val)
	}

	return c.datatype
}

// Unique makes a columns value unique in the table.
func (c *Column) Unique() *Column {
	c.meta["unique"] = true
	return c
}

// Default sets the default value for a Column.
// TODO: value should be an interface to reflect.
func (c *Column) Default(value string) *Column {
	c.meta["default"] = value
	return c
}

// Increments determines if the column auto-increments or not.
// The default value is false.
func (c *Column) Increments() *Column {
	c.meta["autoincrement"] = true
	return c
}

// PrimaryKey determines if a column is the table's primary key.
// The default value for a column is false.
func (c *Column) PrimaryKey() *Column {
	c.meta["primary key"] = true
	return c
}

// NotNullable determines if a column's value can be null.
func (c *Column) NotNullable() *Column {
	c.meta["not null"] = true
	return c
}

// Nullable determines if a column's value can be null.
func (c *Column) Nullable() *Column {
	c.meta["not null"] = false
	return c
}

// Length adds a length constraint to applicable columns.
// TODO: Should this throw an error on columns that can't
// have a length modifier? Like TEXT?
func (c *Column) Length(len int) *Column {
	c.meta["length"] = len
	return c
}
