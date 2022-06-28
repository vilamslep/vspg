package config

type Email struct {
	User       string   `yaml:"user"`
	SenderName string   `yaml:"fromName"`
	Password   string   `yaml:"password"`
	SmtpHost   string   `yaml:"smtp_host"`
	SmtpPort   int      `yaml:"smtp_port"`
	Recivers   []string `yaml:"recivers"`
	Letter     `yaml:"letter"`
}

type Letter struct {
	Subject string `yaml:"subject"`
}
