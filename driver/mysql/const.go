package mysql

const (
	dropTableSchema   = `DROP TABLE %s;`
	renameTableSchema = `ALTER TABLE %s RENAME TO %s;`
	createTableSchema = `CREATE TABLE %s (
	%s
);`
	migrationSchema = `CREATE TABLE migrations (
	id integer not null primary key auto_increment,
	migration varchar(255) not null,
	batch integer not null
);`
	dropSchema = `SELECT table_name FROM information_schema.tables WHERE table_schema = SCHEMA()`
)
