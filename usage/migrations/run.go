// Code generated by Exodus [https://github.com/tmus/exodus]. DO NOT EDIT.

package main

import (
	"os"

	"github.com/tmus/exodus"
	"github.com/tmus/exodus/driver/mysql"
)

var migrations []exodus.Migration = []exodus.Migration{	migration20220121120254create_users_table{},
} // END OF MIGRATIONS

func main() {
	var  driver exodus.Driver
	driver, err := mysql.NewDriver("root:root@/db")
	if err != nil {
		panic(err)
	}
	defer driver.Close()

	migrator, err := exodus.NewMigrator(driver)
	if err != nil {
		panic(err)
	}

	if err := migrator.Run(os.Args, migrations...); err != nil {
		panic(err)
	}
}
