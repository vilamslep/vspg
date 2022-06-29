package backup

import (
	"github.com/vilamslep/psql.maintenance/lib/config"
	"github.com/vilamslep/psql.maintenance/lib/fs"
	"github.com/vilamslep/psql.maintenance/logger"
	"github.com/vilamslep/psql.maintenance/postgres/psql"
	"github.com/vilamslep/psql.maintenance/render"
)

type Task struct{
	Name string
	Kind int
	Items []Item
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
	logger.Info()//'temp directory is {tmpath}'
	logger.Info()//'root path is {rpath}'

	for _, item := range t.Items {
		logger.Info()//f'start handling \'{item.database.name}\''
		if err := item.ExecuteBackup(tmpath, rpath); err == nil {
			logger.Info()//f'finish handling \'{item.database.name}\''
		} else {
			logger.Error()//error
		}
	}

    logger.Info()//'removing old copies'
	//kind_path = os.path.dirname(rpath)
	//fs.clear_old_backup(kind_path, self.keep_count)

	return nil
}

func (t *Task) CountStatuses() (cerr int, cwarn int, csuc int) {
	for _, i := range t.Items {
		switch i.Status {
		case render.StatusError:
			cerr +=1
		case render.StatusSuccess:
			csuc +=1
		default:
			cwarn +=1
		}
	}
	return	
}

func NewTask(name string, kind int, dbs []string, keepCount int) (*Task, error){
	t := Task{
		Name: name,
		Kind:  kind,
		KeepCount: keepCount,
	}

	if dbsInServer, err := psql.Databases(PGConnectionConfig, dbs); err == nil {
		addNotFoundDatabases(dbs, dbsInServer)
		
		for _, db := range dbsInServer {
			item := NewItem(db)
			t.Items = append(t.Items, item)
		}
		
		return &t, err
	
	} else {
		return nil, err
	}
}

func addNotFoundDatabases(dbs[]string, dbsInServer []psql.Database) {
	for _, i := range dbs {
		found := false
		for _, j := range dbsInServer {
			found =  i == j.Name 
			if found { 
				break
			}
		}
		if !found {
			dbsInServer = append(dbsInServer, psql.Database{Name: i})
		}
	}	
}

func CreateTaskBySchedules(schedules config.Schedule) ([]Task, error) {

	tasks := make([]Task,0,3)
	if daily, exist, err := createTask(schedules.Daily); err == nil && exist{
		tasks = append(tasks, daily)
	} else {
		return nil, err
	}
	
	if weekly, exist, err := createTask(schedules.Weekly); err == nil && exist {
		tasks = append(tasks, weekly)
	} else {
		return nil, err
	}

	if monthly, exist, err := createTask(schedules.Monthly); err == nil && exist{
		tasks = append(tasks, monthly)
	} else {
		return nil, err
	}
	return tasks, nil 
}

func createTask(sch config.ScheduleItem) (t Task, ok bool, err error) {
	if sch.Empty() {
		return
	}

    if sch.NeedToRun(){
		
		name := sch.GetKindPrewiew()
		if t, err := NewTask(name, sch.Kind, sch.Dbs, sch.KeepCount); err == nil {
			return *t, true, nil
		} else {
			return Task{}, false, err
		}

	} else {
		return
	}
}