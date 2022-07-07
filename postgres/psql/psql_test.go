package psql

import (
	"os"
	"testing"
)

func TestDatabases(t *testing.T) {

	AllDatabasesTxt = `
	SELECT datname, oid 
	FROM pg_database 
	WHERE NOT datname IN ('postgres', 'template1', 'template0')`
	SearchDatabases = `
	SELECT datname, oid 
	FROM pg_database 
	WHERE datname IN ($1)`

	conf := ConnectionConfig{
		User:     "postgres",
		Password: "142543",
		Database: Database{Name: "postgres"},
		SSlMode:  false,
	}
	dbs := make([]string, 0, 1)
	dbs = append(dbs, "'postgres'", "'template1'")
	Databases(conf, dbs)
}
func TestExcludedTables(t *testing.T) {
	LargeTablesTxt = `
	SELECT table_name as name
	FROM (
		SELECT table_name,pg_total_relation_size(table_name) AS total_size
		FROM (
			SELECT (table_schema || '.' || table_name) AS table_name 
			FROM information_schema.tables) AS all_tables
			 ORDER BY total_size DESC) AS pretty_sizes 
	WHERE total_size > 4294967296;`
	conf := ConnectionConfig{
		User:     "postgres",
		Password: "142543",
		Database: Database{Name: "kfk"},
		SSlMode:  false,
	}
	if tabls, err := ExcludedTables(conf); err != nil {
		t.Fatal(err)
	} else if len(tabls) == 0 {
		t.Fatalf("Should be more than %d large tables", len(tabls))
	}
}
func TestCopyBinary(t *testing.T) {
	PsqlExe = "C:\\PostgreSQL\\bin\\psql.exe"
	os.Setenv("PGPASSWORD", "142543")

	if err := CopyBinary("kfk", "public._inforg12487", "C:\\ut\\bin\\public._inforg12487"); err != nil {
		t.Fatal(err)
	}
}
