package email

import(
	"testing"
)

func TestClientConnection(t *testing.T) {
	user := "USER"
	password := "HOST"
	if err := SendEmail(user, password); err != nil {
		t.Fatal(err)
	}
	// smtp_host := "smtp.yandex.ru"
	// smtp_port := 587



}