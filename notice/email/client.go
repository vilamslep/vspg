package email

import (
	"net/smtp"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func SendEmail(user string, password string) error {
	to := []string{
		"almuertservice@yandex.ru",
	}
	subject := "Test subject"
	body := `<p>An old <b>falcon</b> in the sky.</p>`

	letter := Letter{
		Subject: subject,
		From:    user,
		To:      to,
		Body:    body,
	}

	addr := "smtp.yandex.ru:587"
	host := "smtp.yandex.ru"

	msg := letter.BuildHtmlLetter("vitaliy novak")
	auth := smtp.PlainAuth("", user, password, host)
	return smtp.SendMail(addr, auth, user, to, []byte(msg))

}

// func BuildMessage(mail Mail) string {
// 	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
// 	msg += fmt.Sprintf("From: Vitaliy Novak<%s>\r\n", mail.Sender)
// 	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
// 	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
// 	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

// 	return msg
// }
