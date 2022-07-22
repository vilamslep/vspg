package backup

import (
	"path/filepath"
	"time"

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
}

func (t *Task) Run(config config.Config) (err error) {
	var tmpath, rpath string
	if tmpath, err = fs.TempDir(); err != nil {
		return err
	}

	if rpath, err = fs.GetRootDir(config.Folder.Path, t.Name, t.Kind); err != nil {
		return err
	}
	logger.Infof("temp directory is %s", tmpath)
	logger.Infof("root directory is %s", rpath)

	for _, item := range t.Items {
		logger.Infof("start handlind '%s'", item.Database.Name)
		if err := item.ExecuteBackup(tmpath, rpath); err == nil {
			logger.Info("finish handling '%s'", item.Database.Name)
		} else {
			logger.Errorf("handling database '%s' is failed. %v", item.Database.Name, err)
		}
	}

	logger.Info("removing old copies")

	if err := fs.ClearOldBackup(filepath.Dir(rpath), t.KeepCount); err != nil {
		return err
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

func NewTask(name string, kind int, dbs []string, files []string, keepCount int) (*Task, error) {
	t := Task{
		Name:      name,
		Kind:      kind,
		KeepCount: keepCount,
	}

	if len(dbs) > 0 || len(files) > 0 {
		if len(dbs) > 0 {
			err := t.appendPgDbs(dbs)
			if err != nil {
				return &t, err
			}
		}

		if len(files) > 0 {
			t.appendFiles(files)
		}

		return &t, nil
	} else {
		return &t, nil
	}
}

func (t *Task) appendPgDbs(dbs []string) error {
	if dbsInServer, err := psql.Databases(PGConnectionConfig, dbs); err == nil {
		dbsInServer = addNotFoundDatabases(dbs, dbsInServer)

		for _, db := range dbsInServer {
			item := NewItem(POSTGRES, db, "")
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
	if daily, exist, err := createTask(schedules.Daily); err == nil {
		if exist {
			tasks = append(tasks, daily)
		}
	} else {
		return nil, err
	}

	if weekly, exist, err := createTask(schedules.Weekly); err == nil {
		if exist {
			tasks = append(tasks, weekly)
		}

	} else {
		return nil, err
	}

	if monthly, exist, err := createTask(schedules.Monthly); err == nil {
		if exist {
			tasks = append(tasks, monthly)
		}
	} else {
		return nil, err
	}
	return tasks, nil
}

func createTask(sch config.ScheduleItem) (t Task, ok bool, err error) {
	if sch.Empty() {
		return
	}
	sch.Today = time.Now()

	if sch.NeedToRun() {

		name := config.GetKindPrewiew(sch.Kind)
		if t, err := NewTask(name, sch.Kind, sch.Dbs, sch.Files, sch.KeepCount); err == nil {
			return *t, true, nil
		} else {
			return Task{}, false, err
		}

	} else {
		return
	}
}
