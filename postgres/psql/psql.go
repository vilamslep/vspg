package psql

import (
	"bytes"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
	vos "github.com/vilamslep/vspg/os"
)

type Database struct {
	Name string
	OID  int
}

type ConnectionConfig struct {
	User     string
	Password string
	Database
	SSlMode bool
}

func (c ConnectionConfig) String() string {
	var mode string
	if c.SSlMode {
		mode = "enable"
	} else {
		mode = "disable"
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", c.User, c.Password, c.Name, mode)
}

var (
	AllDatabasesTxt string
	LargeTablesTxt  string
	SearchDatabases string
	PsqlExe         string
)

//TODO I have to just define query args witout replacing substring in the query text
func Databases(pgConf ConnectionConfig, dbsFilter []string) ([]Database, error) {

	db, err := createConnection(pgConf)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var txt string
	if len(dbsFilter) > 0 {
		nf := make([]string, 0, len(dbsFilter))
		for i := range dbsFilter {
			nf = append(nf, fmt.Sprintf("'%s'", dbsFilter[i]))
		}
		txt = strings.ReplaceAll(SearchDatabases, "$1", strings.Join(nf, ","))
	} else {
		txt = AllDatabasesTxt
	}

	rows, err := db.Query(txt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dbs := []Database{}

	for rows.Next() {
		db := Database{}
		if err := rows.Scan(&db.Name, &db.OID); err == nil {
			dbs = append(dbs, db)
		} else {
			return nil, err
		}
	}
	return dbs, nil
}

func ExcludedTables(pgConf ConnectionConfig) ([]string, error) {
	db, err := createConnection(pgConf)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(LargeTablesTxt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables := make([]string, 0, 0)

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err == nil {
			tables = append(tables, table)
		} else {
			return nil, err
		}
	}
	return tables, nil
}

func createConnection(pgConf ConnectionConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", pgConf.String())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CopyBinary(db string, src string, dst string) (err error) {
	var stderr bytes.Buffer

	command := fmt.Sprintf("COPY %s TO '%s' WITH BINARY;", src, dst)
	cmd := exec.Command(PsqlExe, "--dbname", db, "--command", command)
	
	cmd.Stderr = &stderr

	if err := vos.ExecCommand(cmd); err != nil {
		return errors.Wrapf(err, "binary copying is failed. Command %s. \n stderr: %s", command, stderr.String())
	}
	return err
}
