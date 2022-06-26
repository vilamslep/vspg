package email

import (
	"fmt"
	"strings"
)

type Letter struct {
	Subject string
	From    string
	To      []string
	Body    string
}

func (l Letter) BuildHtmlLetter(fromName string) string {
	return fmt.Sprintf(
		(header() + from() + to() + subject() + body()), 
		fromName, l.From, strings.Join(l.To, ";"), l.Subject, l.Body,
	)
}

func header() string {
	return "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
}

func from() string {
	return "From: %s<%s>\r\n"
}

func to() string {
	return "To: %s\r\n"
}

func subject() string {
	return "Subject: %s\r\n"
}

func body() string {
	return "\r\n%s\n"
}


