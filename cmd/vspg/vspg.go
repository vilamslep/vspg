package main

import (
	"errors"

	"github.com/spf13/pflag"
	"github.com/vilamslep/vspg/internal/backup"
	"github.com/vilamslep/vspg/internal/config"
	"github.com/vilamslep/vspg/pkg/fs"
	"github.com/vilamslep/vspg/pkg/logger"
	"github.com/vilamslep/vspg/pkg/postgres/pgdump"
	"github.com/vilamslep/vspg/pkg/postgres/psql"
)

//cli args
var (
	envfile string
	settingFile string
	showHelp bool
)

//errors
var (
	configErr error = errors.New("not defined setting file")
	envErr error = errors.New("not defined enviroment file")
)

func main() {
	
	setAndParseFlags()
	
	if showHelp {
		pflag.Usage()
		return
	}

	if err := checkArgs(); err != nil {
		logger.Fatal(err)
	}

	c, err := config.LoadSetting(settingFile)
	if err != nil {
		logger.Fatalf("loading config is failed. %v", err)
	}

	if err := initModules(c); err != nil {
		logger.Fatalf("module initing is falled; %v", err)
	}

	if b, err := backup.NewBackupProcess(c); err == nil {
		b.Run()
	} else {
		logger.Fatalf("creating backup process is failed. %v", err)
	}
}

func initModules(conf config.Config) error {
	psql.PsqlExe = conf.Psql
	pgdump.PGDumpExe = conf.Dump
	fs.CompressExe = conf.Compress
	fs.WIN_OS_PROGDATA = conf.TempPath

	setQueriesText()

	if err := fs.LoadEnvfile(envfile); err != nil {
		return err
	}

	return nil
}

func setQueriesText() {
	psql.AllDatabasesTxt = `
		SELECT datname, oid 
		FROM pg_database 
		WHERE NOT datname IN ('postgres', 'template1', 'template0')`
	
	psql.SearchDatabases = `
		SELECT datname, oid 
		FROM pg_database WHERE datname IN ($1)`

	psql.LargeTablesTxt = `
		SELECT table_name as name 
		FROM (SELECT table_name,pg_total_relation_size(table_name) AS total_size
				FROM (SELECT (table_schema || '.' || table_name) AS table_name FROM information_schema.tables) AS all_tables 
				ORDER BY total_size DESC) AS pretty_sizes 
		WHERE total_size > 4294967296;`
}

func setAndParseFlags() {
	pflag.BoolVarP(&showHelp, "help", "",
		false,
		"Print help message")
	pflag.StringVarP(&settingFile, "setting", "s",
		"",
		"File common setting")
	pflag.StringVarP(&envfile, "env", "e",
		"",
		"File with enviroment variables")

	pflag.Parse()
}

func checkArgs() error {
	if settingFile == "" {
		return configErr
	} 

	if envfile == "" {
		return envErr
	}
	return nil
}