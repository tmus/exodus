package mysql

import (
	"fmt"

	"github.com/tmus/exodus/column"
)

type columnFormatter func(column column.Definition) (string, error)

func (d mysqlDriver) formatters() map[string]columnFormatter {
	return map[string]columnFormatter{
		"string":  makeStringColumn,
		"boolean": makeBooleanColumn,
		"date":    makeDateColumn,
		"time":    makeTimeColumn,
		"int":     makeIntColumn,
	}
}

// makeColumn takes a column.Definition and turns it into a string representation
// of that column, in a format that is understood by the driver.
func (d mysqlDriver) makeColumn(c column.Definition) (string, error) {
	if fn, ok := d.formatters()[c.Kind]; ok {
		return fn(c)
	}

	return "", fmt.Errorf("unable to make column `%s`. unknown kind `%s`", c.Name, c.Kind)
}

func makeStringColumn(c column.Definition) (string, error) {
	modifiers := []string{"nullable", "default", "unique"}
	sql := makeModifierSQL(modifiers, c)

	len, ok := c.Metadata["length"].(int)
	if !ok {
		return "", fmt.Errorf("invalid length (%v) for column %s", c.Metadata["length"], c.Name)
	}

	return fmt.Sprintf("%s VARCHAR(%d)%s", c.Name, len, sql), nil
}

func makeBooleanColumn(c column.Definition) (string, error) {
	modifiers := []string{"nullable", "default"}
	sql := makeModifierSQL(modifiers, c)

	return fmt.Sprintf("%s BOOLEAN%s", c.Name, sql), nil
}

func makeDateColumn(c column.Definition) (string, error) {
	modifiers := []string{"nullable", "default"}
	sql := makeModifierSQL(modifiers, c)

	return fmt.Sprintf("%s DATE%s", c.Name, sql), nil
}

func makeTimeColumn(c column.Definition) (string, error) {
	modifiers := []string{"nullable", "default"}
	sql := makeModifierSQL(modifiers, c)

	return fmt.Sprintf("%s TIME%s", c.Name, sql), nil
}

func makeIntColumn(c column.Definition) (string, error) {
	modifiers := []string{"nullable", "default", "primary_key", "increments"}
	sql := makeModifierSQL(modifiers, c)

	return fmt.Sprintf("%s INT%s", c.Name, sql), nil
}
