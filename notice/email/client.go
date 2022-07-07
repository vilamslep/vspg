package email

import (
	"fmt"
	"net/smtp"

	"github.com/vilamslep/psql.maintenance/notice"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

type SmptClient struct {
	user     string
	password string
	server   string
	port     int
}

func (sc SmptClient) Send(msgr notice.Messager) error {

	to := []string{
		sc.user,
	}

	addr := fmt.Sprintf("%s:%d", sc.server, sc.port)

	if msg, err := msgr.BuildMessage(); err == nil {
		auth := smtp.PlainAuth("", sc.user, sc.password, sc.server)
		return smtp.SendMail(addr, auth, sc.user, to, msg)
	} else {
		return err
	}
}

func NewSmptClient(user string, password string, server string, port int) (client *SmptClient) {
	return &SmptClient{
		user: user,
		password: password,
		server: server,
		port: port,
	}
}
