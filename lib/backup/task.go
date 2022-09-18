package backup

import (
	"path/filepath"
	"time"

	"github.com/vilamslep/vspg/cloud/yandex"
	"github.com/vilamslep/vspg/lib/config"
	"github.com/vilamslep/vspg/lib/fs"
	"github.com/vilamslep/vspg/logger"
	"github.com/vilamslep/vspg/postgres/psql"
	"github.com/vilamslep/vspg/render"
)

type Task struct {
	Name      string
	Kind      int
	Items     []*Item
	KeepCount int
	BucketName string
	BucketRoot string
}

func (t *Task) Run(config config.Config) (err error) {
	var tmpath, rpath string
	if tmpath, err = fs.TempDir(); err != nil {
		return err
	}

	if rpath, err = fs.GetRootDir(config.Folder.Path, t.Name, t.Kind); err != nil {
		return err
	}
	logger.Infof("temp directory is %s, root directory is %s", tmpath, rpath)

	for _, item := range t.Items {
		logger.Infof("start handlind '%s'", item.Database.Name)
		if err := item.ExecuteBackup(tmpath, rpath); err == nil {
			logger.Infof("finish handling '%s'", item.Database.Name)
		} else {
			logger.Errorf("handling database '%s' is failed. %v", item.Database.Name, err)
		}
	}

	logger.Info("removing old copies")

	if err := fs.ClearOldBackup(filepath.Dir(rpath), t.KeepCount); err != nil {
		return err
	}

	logger.Info("removing old copies in cloud")

	s3client, err := yandex.NewClient(t.BucketRoot)
	if err != nil {
		if err == yandex.ErrLoadingConfiguration {
			logger.Error(err)
		}
	} else {
		if err := s3client.KeepNecessaryQuantity(t.BucketName, t.KeepCount); err != nil {
			return err
		}
	}
	
	return nil
}

func (t *Task) CountStatuses() (cerr int, cwarn int, csuc int) {
	for _, i := range t.Items {
		switch i.Status {
		case render.StatusError:
			cerr += 1
		case render.StatusSuccess:
			csuc += 1
		default:
			cwarn += 1
		}
	}
	return
}

func NewTask(schItem config.ScheduleItem) (*Task, error) {
	t := Task{
		Name:      config.GetKindPrewiew(schItem.Kind),
		Kind:      schItem.Kind,
		KeepCount: schItem.KeepCount,
		BucketName: schItem.BucketName,
		BucketRoot: schItem.BucketRoot,
	}

	if len(schItem.Databases) > 0 || len(schItem.Files) > 0 {
		if len(schItem.Databases) > 0 {
			err := t.appendDatabases(schItem.Databases)
			if err != nil {
				return &t, err
			}
		}

		if len(schItem.Files) > 0 {
			t.appendFiles(schItem.Files)
		}

		return &t, nil
	} else {
		return &t, nil
	}
}

func (t *Task) appendDatabases(dbs []string) error {

	if databasesInServer, err := psql.Databases(PGConnectionConfig, dbs); err == nil {
		databasesInServer = addNotFoundDatabases(dbs, databasesInServer)

		for _, db := range databasesInServer {
			item := NewItem(POSTGRES, db, "")
			item.BucketName = t.BucketName
			item.BucketRoot = t.BucketRoot + "/" + time.Now().Format("02-01-2006")
			t.Items = append(t.Items, &item)
		}
	} else {
		return err
	}

	return nil
}

func (t *Task) appendFiles(files []string) {
	for _, f := range files {
		item := NewItem(FILE, psql.Database{}, f)
		t.Items = append(t.Items, &item)
	}
}

func addNotFoundDatabases(dbs []string, dbsInServer []psql.Database) []psql.Database {
	for _, i := range dbs {
		found := false
		for _, j := range dbsInServer {
			found = i == j.Name
			if found {
				break
			}
		}
		if !found {
			dbsInServer = append(dbsInServer, psql.Database{Name: i})
		}
	}
	return dbsInServer
}

func CreateTaskBySchedules(schedules config.Schedule) ([]Task, error) {
	tasks := make([]Task, 0, 3)
	
	sch := []config.ScheduleItem{
		schedules.Daily,
		schedules.Weekly,
		schedules.Monthly,
	}

	for _, it := range sch {
		if task, exist, err := createTask(it); err == nil && exist {
			tasks = append(tasks, task)
		} else if err != nil {
			return nil, err
		}
	}
	
	return tasks, nil
}

func createTask(sch config.ScheduleItem) (t Task, ok bool, err error) {
	if sch.Empty() {
		return
	}
	sch.Today = time.Now()

	if sch.NeedToRun() {
		if t, err := NewTask(sch); err == nil {
			return *t, true, nil
		} else {
			return Task{}, false, err
		}
	} else {
		return
	}
}
