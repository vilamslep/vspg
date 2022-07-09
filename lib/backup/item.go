package backup

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vilamslep/psql.maintenance/lib/config"
	"github.com/vilamslep/psql.maintenance/lib/fs"
	"github.com/vilamslep/psql.maintenance/logger"
	"github.com/vilamslep/psql.maintenance/postgres/pgdump"
	"github.com/vilamslep/psql.maintenance/postgres/psql"
	"github.com/vilamslep/psql.maintenance/render"
)

type RestoreConfig struct {
	Data struct {
		Tables []string
	}
}

type Item struct {
	psql.Database
	Status       int
	StartTime    time.Time
	FinishTime   time.Time
	DatabaseSize int64
	BackupSize   int64
	BackupPath   string
	Details      string
	config       config.Config
}

func (i *Item) ExecuteBackup(tempDir string, targetDir string) (err error) {
	if i.OID == 0 {
		i.Status = render.StatusWarning
		i.Details = "oid is empty. Database isn't found in server"
		return nil
	}

	i.StartTime = time.Now()
	i.setDatabaseSize()

	err = i.backup(tempDir, targetDir)

	if err != nil {
		i.Details += err.Error() + ";"
		i.Status = render.StatusError
	} else {
		i.Status = render.StatusSuccess
	}
	i.FinishTime = time.Now()

	return err
}

func (i *Item) backup(tempDir string, targetDir string) (err error) {

	logger.Debug("checking space in template directory")

	if ok, err := i.checkSpace(tempDir); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("%s doesn't have enough space. Db %s. Size %d", tempDir, i.Name, i.DatabaseSize)
	} else {
		logger.Debug("success")
	}

	logger.Debug("checking space in target directory")
	if ok, err := i.checkSpace(targetDir); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("%s doesn't have enough space. Db %s. Size %d", targetDir, i.Name, i.DatabaseSize)
	} else {
		logger.Debug("success")
	}
	chdir := make([]string, 0, 2)
	chdir = append(chdir, "logical")
	nc := PGConnectionConfig
	nc.Database = i.Database
	excludeTabls, err := psql.ExcludedTables(nc)
	if err != nil {
		return err
	}

	if len(excludeTabls) > 0 {
		logger.Debugf("excluded tables are %s", strings.Join(excludeTabls, ","))
		chdir = append(chdir, "binary")
	}

	locations, err := fs.CreateDirectories(tempDir, i.Name, chdir)
	if err != nil {
		return err
	}
	logger.Debug("start dumping")
	if err := i.dump(locations["logical"], excludeTabls); err != nil {
		return err
	}
	logger.Debug("success")

	if len(excludeTabls) > 0 {
		logger.Debug("uploading binary data")
		if biniriesFiles, err := i.unloadBinaryTable(locations["binary"], excludeTabls); err == nil {
			err := i.writeRestoreFile(locations["main"], biniriesFiles)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	pathBackDb := locations["main"]
	archive := fmt.Sprintf("%s.zip", pathBackDb)
	logger.Debug("compressing")
	if err := fs.Compress(pathBackDb, archive); err != nil {
		return err
	}
	logger.Debug("success")

	logger.Debug("coping backup to target directory")
	dstName := filepath.Base(archive)
	i.BackupPath = fmt.Sprintf("%s\\%s", targetDir, dstName)
	if err := fs.CopyFile(archive, i.BackupPath); err == nil {
		i.BackupSize, err = fs.GetSize(i.BackupPath)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	logger.Debug("success")
	logger.Debug("removing temp files")
	if err := fs.Remove(pathBackDb); err != nil {
		return err
	}
	if err := fs.Remove(archive); err != nil {
		return err
	}
	logger.Debug("success")

	return nil
}

func (i *Item) checkSpace(path string) (bool, error) {
	dbStora := fmt.Sprintf("%s\\base\\%d", DatabaseLocation, i.OID)
	return fs.IsEnoughSpace(dbStora, path, i.BackupSize)
}

func (i *Item) setDatabaseSize() error {
	dbStora := fmt.Sprintf("%s\\base\\%d", DatabaseLocation, i.OID)
	if c, err := fs.GetSize(dbStora); err == nil {
		i.DatabaseSize = c
	} else {
		return err
	}
	return nil
}

func (i *Item) dump(lpath string, excludeTabls []string) error {
	fout := fmt.Sprintf("%s.log", i.Name)
	out, err := os.Create(fout)
	if err != nil {
		return err
	}

	stdout, stderr, err := pgdump.Dump(i.Name, lpath, excludeTabls)
	if err != nil {
		return fmt.Errorf(stderr.String(), err)
	}

	wrt := bufio.NewWriter(out)
	if stdout.Len() > 0 {
		wrt.Write(stdout.Bytes())
	}

	if stderr.Len() > 0 {
		wrt.Write(stderr.Bytes())
	}
	wrt.Flush()
	isErrors, err := i.findErrorInDumpLog(fout)

	if err != nil {
		return err
	}

	if isErrors {
		out.Close()
		return fmt.Errorf("dumping ended with errors. check dumping log %s", fout)
	}
	out.Close()
	if err := fs.Remove(fout); err != nil {
		i.Details = err.Error() + ";"
	}

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
