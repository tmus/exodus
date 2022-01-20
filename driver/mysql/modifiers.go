package mysql

import (
	"fmt"

	"github.com/tmus/exodus/column"
)

type columnModifier func(column column.Definition) string

var modifiers = map[string]columnModifier{
	"nullable":    nullableModifier,
	"default":     defaultModifier,
	"increments":  autoincrementModifier,
	"primary_key": primaryKeyModifier,
}

func makeModifierSQL(desired []string, column column.Definition) string {
	sql := ""
	for _, d := range desired {
		if fn, ok := modifiers[d]; ok {
			sql = sql + fn(column)
		}
	}

	return sql
}

func nullableModifier(column column.Definition) string {
	value, ok := column.Metadata["nullable"].(bool)
	if ok && value {
		return ""
	}

	return " NOT NULL"
}

func defaultModifier(column column.Definition) string {
	value, ok := column.Metadata["default"].(string)
	if !ok || value == "" {
		return ""
	}

	return fmt.Sprintf(` DEFAULT "%s"`, value)
}

func autoincrementModifier(column column.Definition) string {
	value, ok := column.Metadata["increments"].(bool)
	if !ok || !value {
		return ""
	}

	return " AUTO_INCREMENT"
}

func primaryKeyModifier(column column.Definition) string {
	value, ok := column.Metadata["primary_key"].(bool)
	if ok && value {
		return " PRIMARY KEY"
	}

	return ""
}

func uniqueModifier(column column.Definition) string {
	value, ok := column.Metadata["unique"].(bool)
	if ok && value {
		return " UNIQUE"
	}

	return ""
}
