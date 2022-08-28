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

var (
	tAllDBs      string = "all_databases.sql"
	tSearchDbs   string = "search_database.sql"
	tLargeTables string = "large_tables.sql"
)

func main() {
	c, err := config.LoadSetting("setting.yaml")
	if err != nil {
		logger.Fatalf("loading config is failed. %v", err)
	}

	initModules(c)

	b, err := backup.NewBackupProcess(c)
	if err != nil {
		logger.Fatalf("creating backup process is failed. %v", err)
	}

	b.Run()
}

func initModules(conf config.Config) {
	var err error

	psql.PsqlExe = conf.Psql
	pgdump.PGDumpExe = conf.Dump
	fs.CompressExe = conf.Compress

	fs.LoadEnvfile(conf.Envfile)

	if conf.Queries != "" {
		if psql.AllDatabasesTxt, err = exportQueryFromFile(fmt.Sprintf("%s\\%s", conf.Queries, tAllDBs)); err != nil {
			logger.Fatal(err)
		}
		if psql.SearchDatabases, err = exportQueryFromFile(fmt.Sprintf("%s\\%s", conf.Queries, tSearchDbs)); err != nil {
			logger.Fatal(err)
		}
		if psql.LargeTablesTxt, err = exportQueryFromFile(fmt.Sprintf("%s\\%s", conf.Queries, tLargeTables)); err != nil {
			logger.Fatal(err)
		}
	}
}

func exportQueryFromFile(path string) (string, error) {
	if t, err := ioutil.ReadFile(path); err == nil {
		return string(t), err
	} else {
		return "", fmt.Errorf("can't read file %s, %v", path, err)
	}
}