package main

import (
	"fmt"
	"io/ioutil"

	"github.com/vilamslep/vspg/lib/backup"
	"github.com/vilamslep/vspg/lib/config"
	"github.com/vilamslep/vspg/lib/fs"
	"github.com/vilamslep/vspg/logger"
	"github.com/vilamslep/vspg/postgres/pgdump"
	"github.com/vilamslep/vspg/postgres/psql"
)

func main() {
	c, err := config.LoadSetting("setting.yaml")
	if err != nil {
		logger.Fatalf("loading config is failed. %v", err)
	}

	psql.PsqlExe = c.Utils.Psql
	pgdump.PGDumpExe = c.Utils.Dump
	fs.CompressExe = c.Utils.Compress

	var (
		tAllDBs      string = "all_databases.sql"
		tSearchDbs   string = "search_database.sql"
		tLargeTables string = "large_tables.sql"
	)
	if c.App.Folders.Queries != "" {
		if t, err := ioutil.ReadFile(fmt.Sprintf("%s\\%s", c.App.Folders.Queries, tAllDBs)); err == nil {
			psql.AllDatabasesTxt = string(t)
		} else {
			logger.Fatalf("can't read file %s, %v", tAllDBs, err)
		}
	
		if t, err := ioutil.ReadFile(fmt.Sprintf("%s\\%s", c.App.Folders.Queries, tSearchDbs)); err == nil {
			psql.SearchDatabases = string(t)
		} else {
			logger.Fatalf("can't read file %s, %v", tSearchDbs, err)
		}
	
		if t, err := ioutil.ReadFile(fmt.Sprintf("%s\\%s", c.App.Folders.Queries, tLargeTables)); err == nil {
			psql.LargeTablesTxt = string(t)
		} else {
			logger.Fatalf("can't read file %s, %v", tAllDBs, err)
		}
	}
	
	b, err := backup.NewBackupProcess(c)
	if err != nil {
		logger.Fatalf("creating backup process is failed. %v", err)
	}

	b.Run()
}
