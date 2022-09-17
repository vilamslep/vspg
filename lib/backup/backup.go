package backup

import (
	"fmt"
	"time"

	"github.com/vilamslep/vspg/lib/config"
	"github.com/vilamslep/vspg/logger"
	"github.com/vilamslep/vspg/notice"
	"github.com/vilamslep/vspg/notice/email"
	"github.com/vilamslep/vspg/postgres/psql"
	"github.com/vilamslep/vspg/render"
)

var (
	DatabaseLocation   string
	LogsErrors         []string
	PGConnectionConfig psql.ConnectionConfig
)

type BackupProcess struct {
	config config.Config
	date   time.Time
	tasks  []Task
	sender notice.Sender
	status int
}

func (b *BackupProcess) Run() {
	logger.Info("start backuping")

	for _, t := range b.tasks {
		logger.Infof("handling of %s", config.GetKindPrewiew(t.Kind))
		if err := t.Run(b.config); err != nil {
			logger.Errorf("handling task is failed. %v", err)
		}
	}

	if err := b.sendNotification(); err != nil {
		logger.Error("Notification is failed. %v", err)
	}
}

func (b *BackupProcess) sendNotification() error {
	if len(b.tasks) == 0 {
		return nil
	}

	if content, err := b.renderReport(); err == nil {
		letter := email.Letter{
			Subject:  fmt.Sprintf("%s [%s]", b.config.Email.Subject, render.GetStatusPreview(b.status)),
			From:     b.config.Email.User,
			FromName: b.config.Email.SenderName,
			To:       b.config.Email.Recivers,
			Body:     string(content),
		}
		return b.sender.Send(letter)
	} else {
		return err
	}
}

func (b *BackupProcess) renderReport() ([]byte, error) {
	report := render.BackupReport{}
	b.countSetStatus(&report)
	b.copyBuildInStructToReport(&report)

	return render.RenderReport(report, b.config.App.Templates)
}

func (b *BackupProcess) countSetStatus(report *render.BackupReport) {
	for _, t := range b.tasks {
		cerr, cwarn, csuc := t.CountStatuses()
		report.ErrorCount += cerr
		report.WarningCount += cwarn
		report.SuccessCount += csuc
	}

	if report.ErrorCount > 0 {
		report.Status = render.StatusError
	} else if report.WarningCount > 0 {
		report.Status = render.StatusWarning
	} else {
		report.Status = render.StatusSuccess
	}
	b.status = report.Status
}

func (b *BackupProcess) copyBuildInStructToReport(report *render.BackupReport) {
	report.Date = b.date.Format("Monday, 02 January 2006")
	gb := 1024 * 1024 * 1024
	mb := 1024 * 1024
	for _, t := range b.tasks {
		nt := render.Task{}
		nt.Name = t.Name

		for _, i := range t.Items {
			ni := render.Item{}
			if i.Type == POSTGRES {
				ni.Name = i.Name
			} else {
				ni.Name = i.File
			}
			ni.OID = i.OID
			ni.StartTime = i.StartTime.Format("03:04:05")
			ni.FinishTime = i.FinishTime.Format("03:04:05")
			ni.Status = i.Status
			ni.BackupPath = i.BackupPath

			if i.BackupSize > int64(gb) {
				ni.BackupSize = fmt.Sprintf("%.2fGB", float64(i.BackupSize)/float64(gb))
				ni.DatabaseSize = fmt.Sprintf("%.2fGB", float64(i.DatabaseSize)/float64(gb))
			} else {
				ni.BackupSize = fmt.Sprintf("%.2fMB", float64(i.BackupSize)/float64(mb))
				ni.DatabaseSize = fmt.Sprintf("%.2fMB", float64(i.DatabaseSize)/float64(mb))
			}
			ni.Details = i.Details

			nt.Items = append(nt.Items, ni)
		}
		report.Tasks = append(report.Tasks, nt)
	}
}

func NewBackupProcess(conf config.Config) (*BackupProcess, error) {

	b := BackupProcess{config: conf}

	DatabaseLocation = conf.GetDataLocation()
	LogsErrors = getPGDumpErrorEvents()

	PGConnectionConfig = psql.ConnectionConfig{
		User:     conf.GetUser(),
		Password: conf.GetPassword(),
		Database: psql.Database{Name: "postgres"},
		SSlMode:  false,
	}

	b.date = time.Now()
	tasks, err := CreateTaskBySchedules(conf.Schedule)
	if err != nil {
		return nil, err
	}
	b.tasks = tasks

	b.sender = conf.GetSender()

	return &b, nil
}

func getPGDumpErrorEvents() []string {
	return []string{
		"pg_dump: ошибка:",
		"pg_dump: error:",
	}
}
