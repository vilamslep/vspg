package psql

import (
	"testing"
)


func TestDatabases(t *testing.T) {

	AllDatabasesTxt = `SELECT datname, oid 
	FROM pg_database 
	WHERE NOT datname IN ('postgres', 'template1', 'template0')`
	SearchDatabases = `SELECT datname, oid 
	FROM pg_database 
	WHERE datname IN ($1)`
	
	conf := ConnectionConfig{
		User:     "postgres",
		Password: "secret",
		Database: Database{Name: "postgres"},
		SSlMode:  false,
	}
	dbs := make([]string, 0, 1)
	dbs = append(dbs, "'postgres'", "'template1'")
	Databases(conf, dbs)
}
