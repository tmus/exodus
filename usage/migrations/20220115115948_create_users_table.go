package main

import (
	"github.com/tmus/exodus"
	"github.com/tmus/exodus/column"
)

type migration20220115115948create_users_table struct {
	exodus.BaseMigration
}

func (migration20220115115948create_users_table) Up() exodus.Migration {
	return exodus.Create("users", []column.Definition{
		column.String("username", 255),
		column.Boolean("can_email"),
	})
}

func (migration20220115115948create_users_table) Down() exodus.Migration {
	return ""
}
