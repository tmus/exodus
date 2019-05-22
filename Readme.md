# Exodus

Database Migrations in Go.

> Currently, only the SQLite3 driver is supported. Obviously this is not ideal. Support
> for at least MySQL and Postgresql will come in the future.

> **Notice:** This is very beta, and is very subject to change. It may be that eventually
> the package will be rereleased with breaking changes and improvements down the line.
> Please don't rely on this for anything critical.

## Installation

Use Go Modules.

`TODO: Add Installation instructions`

## Usage

There's not exactly a Laravel / Rails / Zend / <Framework> way of running the migrations,
yet, in that there is no command line utility to run or generate the migrations. Much
of this will be streamlined in future releases.

1. Create a new struct type. The type should be the name of the migration:

```go
type CreateUsersTable struct{}
```

2. Define two methods on the created struct: `Up()` and `Down()`. These should both
return an `exodus.Migration`. This satisfies the `exodus.MigrationInterface`.

The `Up()` function should run the *creative* side of the migration, e.g., creating
a new table. The `Down()` function should run the *destructive* side of the migration,
e.g., dropping the table.

```go
func (m CreateUsersTable) Up() exodus.Migration {
	return exodus.Create("users", exodus.Schema{
		column.Int("id").Increments().PrimaryKey(),
		column.String("email", 100).NotNullable().Unique(),
		column.String("name", 60).NotNullable(),
		column.Timestamp("activated_at"),
		column.Date("birthday"),

		column.UniqueSet("unique_name_birthday", "name", "birthday"),
	})
}

// Down reverts the changes on the database.
func (m CreateUsersTable) Down() exodus.Migration {
	return exodus.Drop("users")
}
```

3. As you can see above, there exists a Create method and a Drop method. More methods
(change, add, remove column) will be added at some point.

The `exodus.Create` method accepts a table name as a string, and an `exodus.Schema`, which
is a slice of items that implement the [`exodus.Columnable`](column/Column.go) interface.
It's easy to add columns to this schema, as you can see in the above `Up()` migration.

The supported column types are:

- `column.Binary`: creates a `binary` column.
- `column.Boolean`: creates a `boolean` column.
- `column.Char`: creates a `char` column. Must be passed a length as the second parameter.
- `column.Date`: creates a `date` column.
- `column.DateTime`: creates a `datetime` column.
- `column.Int`: creates an `int` column. Currently only `int` is supported.
- `column.String`: creates a `varchar` column. Must be passed a length as the second parameter.
- `column.Text`: creates a `text` column.
- `column.Timestamp`: creates a `timestamp` column.

These columns can have modifiers chained to them, as you can see in the `Up()` migration
above. Their effects should be obvious:

- `Unique()`
- `Default(value string)`
- `Increments()`
- `PrimaryKey()`
- `NotNullable()`
- `Nullable()`
- `Length()`

4. When your migrations have been created, create an `exodus.Migrator`, and pass it an `*sql.DB`.
The function will return an error if the DB driver passed in is not supported.

```go
db, _ := sql.Open("sqlite3", "./database.db")
defer db.Close()

migrator, err := exodus.NewMigrator(db)
if err != nil {
    log.Fatalln(err)
}
```

5. Finally, use the migrator to run the Migrations. You can pass as many migrations
as you like into the Run function:

```go
migrator.Run(migrations ...MigrationInterface)
```

The tables should now exist in your database.