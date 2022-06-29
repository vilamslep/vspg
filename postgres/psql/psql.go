package psql

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	Name string
	OID  string
}

func Databases(dbsFilter []string) ([]Database, error) {
	connStr := "user=postgres password=123456 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT datname, oid FROM pg_database WHERE NOT datname IN ('postgres', 'template1', 'template0')")
	if err != nil {
		panic(err)
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

// def excluded_tables(db: str) -> list[str]:
//     tables = list()
//     with closing(config.get_connection(psycopg2.connect,db)) as conn:
//         with conn.cursor() as cursor:
//             cursor.execute( txt_stat_tables() )
//             tables = [item[0] for item in cursor.fetchall() ]

//     return tables

func CopyBinary(db string, src string, dst string) (err error) {
	//     tool = config.psql()

	//     args = [tool, '--dbname', db]

	//         path_save = __prepare_path(dst)

	//         cmd = ["COPY", src, "TO", path_save, "WITH BINARY;"]
	//         return __execute(args, " ".join(cmd)) == 0

	//     def __prepare_path(path:str)->str:
	//         npath = path.replace("\\", "\\\\")
	//         return f'\'{npath}\''

	//     def __execute(args:list, command:str)->int:
	//         args.append('--command')
	//         args.append(command)

	//         exit_code = subprocess.Popen(args, stdout=subprocess.PIPE).wait()
	//         return int(exit_code)

	return
}

// def copy_binary(db:str, src:str, dst:str)->dict:

// def txt_custom_databases(dbs:list)->str:
//     if len(dbs) > 0:
//         txtdb = ','.join(list(map(lambda x: f'\'{x}\'',dbs)))
//         filter = f'WHERE datname IN( {txtdb} )'
//     else:
//         filter = 'WHERE NOT datname IN (\'postgres\', \'template1\', \'template0\')'

//     return 'SELECT datname, oid FROM pg_database {};'.format(filter)

// def txt_stat_tables() -> str:
//     '''select tables which are bigger that 4 GB'''
//     return '''
//     SELECT table_name as name
//     FROM (
// 	    SELECT table_name,pg_total_relation_size(table_name) AS total_size
// 	    FROM (
// 		    SELECT (table_schema || '.' || table_name) AS table_name
//             FROM information_schema.tables) AS all_tables
//  	ORDER BY total_size DESC) AS pretty_sizes WHERE total_size > 4294967296;
//     '''
