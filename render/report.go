package render

type Item struct {
	Name         string
	OID          int
	Status       int
	StartTime    string
	FinishTime   string
	DatabaseSize string
	BackupSize   string
	BackupPath   string
	Details      string
}

type Task struct {
	Name  string
	Items []Item
}

type BackupReport struct {
	Status       int
	ErrorCount   int
	SuccessCount int
	WarningCount int
	Tasks        []Task
	Date         string
}
