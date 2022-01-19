package main

import (
	"github.com/tmus/exodus"
	"github.com/tmus/exodus/column"
)

type migration20220115115948create_users_table struct{}

func (migration20220115115948create_users_table) Up(schema *exodus.MigrationPayload) {
	schema.Create("example", []column.Definition{
		column.Increments("id"),
		column.String("username", 255).Nullable().Default("greenhorn"),
		column.String("email", 100).Unique(),
		column.Boolean("is_verified"),
		column.Date("created_at"),
		column.Time("eats_breakfast_at"),
		column.Int("visits").Default("5"),
	})
}

func (migration20220115115948create_users_table) Down(schema *exodus.MigrationPayload) {

}
