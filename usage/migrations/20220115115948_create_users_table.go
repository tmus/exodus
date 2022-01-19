package main

import (
	"github.com/tmus/exodus"
	"github.com/tmus/exodus/column"
)

type migration20220115115948create_users_table struct{}

func (migration20220115115948create_users_table) Up(schema *exodus.MigrationPayload) {
	schema.Create("example", []column.Definition{
		column.String("username", 255),
		column.Boolean("is_verified"),
		column.Date("created_at"),
	})

	schema.Rename("example", "users")
}

func (migration20220115115948create_users_table) Down(schema *exodus.MigrationPayload) {

}
