package main

import (
	"github.com/tmus/exodus"
)

type migration20220115115948create_users_table struct{}

func (migration20220115115948create_users_table) Up() exodus.MigrationPayload {
	return exodus.Rename("users", "boosers")
}

func (migration20220115115948create_users_table) Down() exodus.MigrationPayload {
	return exodus.Drop("users")
}
