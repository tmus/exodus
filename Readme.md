# Exodus

Exodus is a database migration runner for Go. It uses a database-agnostic schema
DSL to enable you to write database migrations that can be ran against any
underlying database.

> Currently, only MySQL is supported. Support for other database drivers is
> planned. If you want to help out, feel free to send a PR!

## Installation

Exodus comes in two parts:

1. It is a library that enables you to write platform agnostic database
   migrations.
2. It is a command line program that enables you to run and manage these
   migrations.

Firstly, you should install the library in your project with
`go get github.com/tmus/exodus@latest`.

Then, install the command line program with
`go install github.com/tmus/exodus/cmd/exodus@latest`.

## Usage

### CLI

Once installed, you can run `exodus help` to get more information on how to use
the CLI.

#### `init`

The `init` command bootstraps a `migrations/` directory in the current working
directory. This will contain a `run.go` file and is where created migrations
will live.

> The `run.go` file is generated and should not be edited manually. However,
> given the fact that this library is a work–in–progress, there may be times
> that you do need to manually make changes, for example to populate the
> database credentials and datasource.

#### `create`

The `create` command will add a new migration file to the directory created by
`init`. You must pass in a single argument, which will be the name of the
migration. The datetime is prepended to this name so that migrations are stored
in order and are guaranteed to be unique.

For example `exodus create create_users_table` would yield a new file called
`20220120211339_create_users_table.go`.

Once the file is created, you may write migration scripts inside it.

#### `run`

The `run` command runs the migrations against the database. Any previously ran
migrations are not ran again.

There are two directions that the migrations can run: `up` and `down`, which can
be specified as part of the command. If a direction is not provided, `up` is
assumed.

`up` will run all migrations that have not yet been ran against the database, by
running the `Up` function inside the migration file.

`down` will revert the last "batch" of migrations that have been ran, by running
the `Down` function inside the migration file.

### Library

To create a new migration, run `exodus create <migration_name>`. This will
create a corresponding migration inside the migrations directory.

Two functions will exist against this migration: `Up` and `Down`, which both
take a single argument of `*exodus.MigrationPayload`. You can use this payload
to register commands to run against the database.

Let's make a `users` table, by calling the `schema.Create` function with a table
name and a slice of `column.Definition`:

```go
func (migration20220121120254create_users_table) Up(schema *exodus.MigrationPayload) {
	schema.Create("users", []column.Definition{})
}
```

This would create an empty table. Let's add some columns. For a `users` table,
let's add an incrementing integer for the primary key, as well as an email (that
needs to be unique!) and password.

```go
func (migration20220121120254create_users_table) Up(schema *exodus.MigrationPayload) {
	schema.Create("users", []column.Definition{
		column.Increments("id"),
		column.String("email", 100).Unique(),
		column.String("password", 255),
	})
}
```

Notice the `Increments` convenience function and the chained `Unique` call.
That'll do for now, but you can check the documentation for a full list of
supported column types and modifiers.

We should write an equivalent `Down` call for this operation. If the `Up` call
creates the table, the `Down` call should drop it:

```go
func (migration20220121120254create_users_table) Down(schema *exodus.MigrationPayload) {
	schema.Drop("users")
}
```

We're in a good position to run our migrations. `cd` into the `migrations`
directory and run `exodus run up`. You'll get visible feedback that the
migrations have ran, or notices of any errors and the reason for failure.

If all went well, you should be able to check your database and see that the
table exists, as well as a record of the migration in the `migrations` table.

To revert this migration, simply run `exodus run down`, which will revert the
last "batch" of migrations.
