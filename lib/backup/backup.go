package backup

import (
	"fmt"
	"os"
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
		logger.Infof("handling of", config.GetKindPrewiew(t.Kind))
		if err := t.Run(b.config); err != nil {
			logger.Errorf("handling task is failed. %v", err)
		}
	}

	if err := b.sendNotification(); err != nil {
		logger.Error("Notification is failed. %v", err)
	}
}

func (b *BackupProcess) sendNotification() error {
	if content, err := b.renderReport(); err == nil {
		letter := email.Letter{
			Subject:  fmt.Sprintf("%s [%s]", b.config.Email.Subject, render.GetStatusPreview(b.status)) ,
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

	return render.RenderReport(report, b.config.App.Folders.Templates)
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
			if i.BackupSize > (1024 * 1024 * 1024) {
				ni.BackupSize = fmt.Sprintf("%.2dGB", (i.BackupSize / 1024 / 1024 / 1024))
				ni.DatabaseSize = fmt.Sprintf("%.2dGB", (i.DatabaseSize / 1024 / 1024 / 1024))
			} else {
				ni.BackupSize = fmt.Sprintf("%.2dMB", (i.BackupSize / 1024 / 1024))
				ni.DatabaseSize = fmt.Sprintf("%.2dMB", (i.DatabaseSize / 1024 / 1024))
			}
			ni.Details = i.Details

			nt.Items = append(nt.Items, ni)
		}
		report.Tasks = append(report.Tasks, nt)
	}
}

func NewBackupProcess(conf config.Config) (*BackupProcess, error) {

	b := BackupProcess{
		config: conf,
	}

	DatabaseLocation = conf.Postgres.DataLocation
	LogsErrors = make([]string, 0, 2)
	LogsErrors = append(LogsErrors, "pg_dump: ошибка:")
	LogsErrors = append(LogsErrors, "pg_dump: error:")

	PGConnectionConfig = psql.ConnectionConfig{
		User:     conf.Postgres.User,
		Password: conf.Postgres.Password,
		Database: psql.Database{Name: "postgres"},
		SSlMode:  false,
	}

	os.Setenv("PGUSER", PGConnectionConfig.Name)
	os.Setenv("PGPASSWORD", PGConnectionConfig.Password)

	b.date = time.Now()
	tasks, err := CreateTaskBySchedules(conf.Schedule)
	if err != nil {
		return nil, err
	}
	b.tasks = tasks

	b.sender = conf.GetSender()

	return &b, nil
}
