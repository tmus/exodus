package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const migrationStub = `package main

import "github.com/tmus/exodus"

type REPLACE_ME struct {
	exodus.BaseMigration
}

// func (REPLACE_ME) Up() exodus.Migration {
//	return ""
// }

// func (REPLACE_ME) Down() exodus.Migration {
//	return ""
// }`

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new migration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		date := getCreateDate()
		file := args[0]
		path := fmt.Sprintf("./migrations/%s_%s.go", date, file)
		err := os.WriteFile(path, generateStub(file), 0755)
		if err != nil {
			log.Fatalf("unable to create file: %s\n", err)
		}

		addMigrationToList("migration" + date + file)
	},
}

// getCreateDate returns a string to prefix a newly created migration, in the format
// YYYYMMDDHHMMSS
func getCreateDate() string {
	return time.Now().Format("20060102030405")
}

func generateStub(name string) []byte {
	return []byte(strings.ReplaceAll(migrationStub, "REPLACE_ME", "migration"+getCreateDate()+name))
}

func addMigrationToList(name string) {
	curr, err := os.ReadFile("./migrations/run.go")
	if err != nil {
		panic(err)
	}

	ok, err := regexp.MatchString(`var migrations \[\]exodus\.MigrationInterface = \[\]exodus\.MigrationInterface{} \/\/ END OF MIGRATIONS`, string(curr))
	if err != nil {
		panic(err)
	}

	var newContents string
	if ok {
		fmt.Println("matched!")
		replacement := fmt.Sprintf(`var migrations []exodus.MigrationInterface = []exodus.MigrationInterface{
	%s{},
} // END OF MIGRATIONS`, name)
		newContents = strings.Replace(string(curr), "var migrations []exodus.MigrationInterface = []exodus.MigrationInterface{} // END OF MIGRATIONS", replacement, 1)
	} else {
		fmt.Println("did not match")
		replacement := fmt.Sprintf(`	%s{},
} // END OF MIGRATIONS`, name)
		newContents = strings.Replace(string(curr), "} // END OF MIGRATIONS", replacement, 1)
	}

	os.WriteFile("./migrations/run.go", []byte(newContents), 0755)
}
