package main

import (
	"github.com/tmus/exodus"
	"github.com/tmus/exodus/column"
)

type migration20220121120254create_users_table struct{}

func (migration20220121120254create_users_table) Up(schema *exodus.MigrationPayload) {
	schema.Create("users", []column.Definition{
		column.Increments("id"),
		column.String("email", 100).Unique(),
		column.String("password", 255),
	})

	schema.Drop("users")

	schema.Raw("create table `what` (`id` varchar(255));")
}

func (migration20220121120254create_users_table) Down(schema *exodus.MigrationPayload) {
	schema.Drop("users")
}
