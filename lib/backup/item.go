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

	"github.com/vilamslep/vspg/cloud/yandex"
	"github.com/vilamslep/vspg/lib/fs"
	"github.com/vilamslep/vspg/logger"
	"github.com/vilamslep/vspg/postgres/pgdump"
	"github.com/vilamslep/vspg/postgres/psql"
	"github.com/vilamslep/vspg/render"
)

var (
	POSTGRES = 1
	FILE     = 2
)

type RestoreConfig struct {
	Data struct {
		Tables []string
	}
}

type Item struct {
	psql.Database
	File         string
	Status       int
	StartTime    time.Time
	FinishTime   time.Time
	DatabaseSize int64
	BackupSize   int64
	BackupPath   string
	Details      string
	BucketName string
	BucketRoot string
	Type         int
}

func (i *Item) ExecuteBackup(tempDir string, targetDir string) (err error) {

	switch i.Type {
	case POSTGRES:
		err = i.pgBackup(tempDir, targetDir)
	case FILE:
		err = i.fileBackup(targetDir)
	default:
		err = fmt.Errorf("unexpected type of item, type is %d", i.Type)
	}

	if err != nil {
		i.Details += err.Error() + ";"
		i.Status = render.StatusError
	} else {
		i.Status = render.StatusSuccess
	}
	i.FinishTime = time.Now()

	return err
}

func (i *Item) fileBackup(targetDir string) (err error) {

	i.StartTime = time.Now()
	i.setDatabaseSize()

	logger.Debug("checking space in target directory")
	{
		if ok, err := i.checkSpace(targetDir); err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("%s doesn't have enough space. Db %s. Size %d", targetDir, i.Name, i.DatabaseSize)
		}
	}
	logger.Debug("coping backup to target directory")
	{
		dstName := filepath.Base(i.File)
		i.BackupPath = fmt.Sprintf("%s\\%s", targetDir, dstName)
		if err := fs.CreateIfNotExists(i.BackupPath); err != nil {
			return err
		}

		if err := fs.Copy(i.File, i.BackupPath); err == nil {
			i.BackupSize, err = fs.GetSize(i.BackupPath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (i *Item) pgBackup(tempDir string, targetDir string) (err error) {

	if i.OID == 0 {
		i.Status = render.StatusWarning
		i.Details = "oid is empty. Database isn't found in server"
		return nil
	}

	i.StartTime = time.Now()
	i.setDatabaseSize()

	logger.Debug("checking free space")
	{
		if err := i.checkingFreeSpace(tempDir, targetDir); err != nil {
			return err
		}
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
	logger.Debug("dumping")
	{
		if err := i.dump(locations["logical"], excludeTabls); err != nil {
			return err
		}
	}

	if len(excludeTabls) > 0 {
		logger.Debug("upload %s as binary data", strings.Join(excludeTabls, ","))
		{
			biniriesFiles, err := i.unloadBinaryTable(locations["binary"], excludeTabls)
			if err != nil {
				return err
			}

			if err := i.writeRestoreFile(locations["main"], biniriesFiles); err != nil {
				return err
			}
		}
	}

	pathBackDb := locations["main"]
	archive := fmt.Sprintf("%s.zip", pathBackDb)
	logger.Debug("compressing")
	{
		if err := fs.Compress(pathBackDb, archive); err != nil {
			return err
		}
	}

	logger.Debug("coping backup to target directory")
	{
		dstName := filepath.Base(archive)
		i.BackupPath = fmt.Sprintf("%s\\%s", targetDir, dstName)
		err := fs.Copy(archive, i.BackupPath)
		if err != nil {
			return err
		}

		if i.BackupSize, err = fs.GetSize(i.BackupPath); err != nil {
			return err
		}
	}

	if i.BucketName != "" && i.BucketRoot != ""  {
		logger.Debug("coping backup to cloud storage")
		{
			if s3client, err := yandex.NewClient(i.BucketRoot); err == nil {
				if err := s3client.Add(archive, i.BucketName); err != nil {
					return err
				}
			} else if err != yandex.ErrLoadingConfiguration {
				return err
			}
		}
	}
	

	logger.Debug("removing temp files")
	{
		if err := fs.Remove(pathBackDb, archive); err != nil {
			return err
		}
	}
	return nil
}

func (i *Item) checkingFreeSpace(path ...string) error {
	for _, target := range path {
		if ok, err := i.checkSpace(target); err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("%s doesn't have enough space. Db %s. Size %d", target, i.Name, i.DatabaseSize)
		}
	}
	return nil
}

func (i *Item) checkSpace(path string) (bool, error) {
	var dbStora string
	if i.Type == POSTGRES {
		dbStora = fmt.Sprintf("%s\\base\\%d", DatabaseLocation, i.OID)
	} else {
		dbStora = i.File
	}

	return fs.IsEnoughSpace(dbStora, path, i.BackupSize)
}

func (i *Item) setDatabaseSize() error {
	if i.Type == POSTGRES {
		dbStora := fmt.Sprintf("%s\\base\\%d", DatabaseLocation, i.OID)
		if c, err := fs.GetSize(dbStora); err == nil {
			i.DatabaseSize = c
		} else {
			return err
		}
	} else {
		if c, err := fs.GetSize(i.File); err == nil {
			i.DatabaseSize = c
		} else {
			return err
		}
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

func NewItem(kind int, db psql.Database, file string) Item {
	return Item{
		Database: db,
		File:     file,
		Type:     kind,
	}
}
