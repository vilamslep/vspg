package backup

import "github.com/vilamslep/psql.maintenance/lib/config"

type Task struct{
	Name string
	Kind int
	Items []Item
	KeepCount int
}

func (t *Task) Run() error {
	return nil
}

func (t *Task) CountStatuses() (cerr int, cwarn int, csuc int) {
	for _, i := range t.Items {
		switch i.Status {
		case "error":
			cerr +=1
		case "success":
			csuc +=1
		default:
			cwarn +=1
		}
	}
	return	
}

func (t *Task) addNotFoundDatabases() {}



func NewTask() {}

func CreateTaskBySchedules(conf config.Config) (rs []Task, err error) {
	
	return
}