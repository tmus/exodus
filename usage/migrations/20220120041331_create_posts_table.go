package main

import (
	"github.com/tmus/exodus"
	"github.com/tmus/exodus/column"
)

type migration20220120041331create_posts_table struct{}

func (migration20220120041331create_posts_table) Up(schema *exodus.MigrationPayload) {
	schema.Create("posts", []column.Definition{
		column.Increments("id"),
		column.String("title", 255).Unique(),
	})
}

func (migration20220120041331create_posts_table) Down(schema *exodus.MigrationPayload) {
	schema.Drop("posts")
}
