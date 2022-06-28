package backup

import "time"

type Database struct{
	Name string
	OID string	
}

type Item struct{
	Database
	Status string
	StartTime time.Time
	FinishTime time.Time
	DatabaseSize float64
	BackupSize float64
	BackupPath string
	Details string
}


func (i *Item) ExecuteBackup() {}

func (i *Item) backup() {}

func (i *Item) checkSpace() {}

func (i *Item) setDatabaseSize() {}

func (i *Item) dump() {}

func (i *Item) findErrorInDumpLog() {}

func (i *Item) unloadBinaryTable() {}

func (i *Item) writeRestoreFile() {}
