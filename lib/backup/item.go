package backup

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/vilamslep/psql.maintenance/lib/config"
	"github.com/vilamslep/psql.maintenance/lib/fs"
	"github.com/vilamslep/psql.maintenance/postgres/pgdump"
	"github.com/vilamslep/psql.maintenance/postgres/psql"
)

type RestoreConfig struct {
	Data struct {
		Tables []string
	}
}

type Item struct {
	psql.Database
	Status       string
	StartTime    time.Time
	FinishTime   time.Time
	DatabaseSize int64
	BackupSize   int64
	BackupPath   string
	Details      string
	config       config.Config
}

func (i *Item) ExecuteBackup(tempDir string, targetDir string) (err error) {
	//         if len(self.database.oid) == 0:
	//             self.status = 'warning'
	//             self.details = 'oid is empty. Database isn\'t found in server'
	//             return
	//         self.start_time = datetime.now()
	//         self.__set_db_size()
	//         self.details = ''
	//         try:
	//             self.__backup(tempdir, targetdir)
	//         except Exception:
	//             exc_info = sys.exc_info()

	//             self.details = exc_info[1]
	//             self.status = 'error'
	//             self.end_time = datetime.now()

	//             raise Exception(f'backup of {self.database.name} is failed').with_traceback(exc_info[2])

	//         self.status = 'success'
	//         self.end_time = datetime.now()

	return nil
}

func (i *Item) backup(tempDir string, targetDir string) (err error) {
	//         logger.info('checking space in template directory')

	//         if not self.__check_space_for_backup(tempdir):
	//             raise Exception(f'there isn\'t enought space in disk. checking path is {tempdir}')

	//         logger.info('success')
	//         logger.info('checking space in target directory')
	//         if not self.__check_space_for_backup(targetdir):
	//             raise Exception(f'there isn\'t enought space in disk. checking path is {targetdir}')

	//         logger.info('success')

	//         chdir = ['logical']

	//         exclude_tabls = excluded_tables(self.database.name)

	//         if len(exclude_tabls) > 0:
	//             logger.info('excluded tables are ' + ','.join(exclude_tabls))
	//             chdir.append('binary')

	//         try:
	//             locations = fs.generate_directories(tempdir, self.database.name, chdir)
	//         except:
	//             trace = sys.exc_info()[2]
	//             raise Exception('generating of directories end with error').with_traceback(trace)

	//         logger.info('dumping...')
	//         self.__dump(locations['logical'], exclude_tabls)
	//         logger.info('success')

	//         if len(exclude_tabls) > 0:
	//             logger.info('uploading binary data')
	//             binaries_files = self.__unload_as_binary(locations['binary'], exclude_tabls)
	//             self.__write_restore_file(locations['main'], binaries_files)
	//             logger.info('success')

	//         path_back_db = locations['main']
	//         archive = f'{path_back_db}.zip'

	//         logger.info('compressing...')
	//         if not compress_dir(path_back_db, archive):
	//             raise Exception('compressing is not success')
	//         logger.info('success')

	//         logger.info('coping backup to target directory')
	//         try:
	//             dst_name = basename(archive)
	//             self.backup_path = f'{targetdir}\\{dst_name}'
	//             fs.copy_file(archive, self.backup_path)
	//             self.size_backup = fs.get_size(self.backup_path)
	//         except:
	//             trace = sys.exc_info()[2]
	//             raise Exception('can\'t copy the backup to target directory').with_traceback(trace)
	//         logger.info('success')

	//         logger.info('removing temp files')
	//         try:
	//             fs.remove(path_back_db)
	//             fs.remove(archive)
	//         except:
	//             trace = sys.exc_info()[2]
	//             raise Exception(f'removing is not success').with_traceback(trace)
	//         logger.info('success')

	return nil
}

func (i *Item) checkSpace(path string) (bool, error) {
	dbStora := fmt.Sprintf("%s\\base\\%s", DatabaseLocation, i.OID)
	return fs.IsEnoughSpace(dbStora, path)
}

func (i *Item) setDatabaseSize() error {
	dbStora := fmt.Sprintf("%s\\base\\%s", DatabaseLocation, i.OID)
	if c, err := fs.GetSize(dbStora); err == nil {
		i.DatabaseSize = c
	} else {
		return err
	}
	return nil
}

func (i *Item) dump(lpath string, excludeTabls []string) error {
	fout := fmt.Sprintf("log\\%s.log", i.Name)
	out, err := os.Create(fout)
	if err != nil {
		return err
	}
	out.Close()

	err = pgdump.Dump(i.Name, lpath, fout, excludeTabls)
	if err != nil {
		return err
	}

	isErrors, err := i.findErrorInDumpLog(fout)

	if err != nil {
		return err
	}

	if isErrors {
		return fmt.Errorf("dumping ended with errors. check dumping log %s", fout)
	}
	fs.Remove(fout)

	return nil
}

func (i *Item) findErrorInDumpLog(logFile string) (bool, error) {
	f, err := os.Open(logFile)
	if err != nil {
		return false, err
	}
	defer f.Close()
	rd := bufio.NewScanner(f)
	for rd.Scan() {
		s := rd.Text()
		for _, er := range LogsErrors {
			if strings.Contains(s, er) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (i *Item) unloadBinaryTable(binPath string, tabls []string) ([]string, error) {
	binfiles := make([]string, 0, len(tabls))

	for _, t := range tabls {
		tblPath := fmt.Sprintf("%s\\%s", binPath, t)
		if err := psql.CopyBinary(i.Name, t, tblPath); err == nil {
			binfiles = append(binfiles, tblPath)
		} else {
			return nil, err
		}
	}
	return binfiles, nil
}

func (i *Item) writeRestoreFile(mainPath string, binfiles []string) error {
	data := RestoreConfig{
		Data: struct{ Tables []string }{
			Tables: binfiles,
		},
	}
	f := fmt.Sprintf("%s\\map.json", mainPath)
	if c, err := json.Marshal(data); err == nil {
		if err := ioutil.WriteFile(f, c, 0777); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func NewItem(db psql.Database) Item {
	return Item{Database: db}
}
