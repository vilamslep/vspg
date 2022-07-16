package config

import (
	"os"

	"github.com/vilamslep/vspg/notice"
	"github.com/vilamslep/vspg/notice/email"
	"gopkg.in/yaml.v2"
)

type App struct {
	Folders struct {
		Templates string `yaml:"templates"`
		Queries string `yaml:"queries"`
	} `yaml:"folders"`
}

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

type Config struct {
	App      `yaml:"app"`
	Postgres `yaml:"postgres"`
	Folder   `yaml:"target_folder"`
	Utils    `yaml:"utils"`
	Email    `yaml:"email"`
	Schedule `yaml:"schedules"`
}

type Folder struct {
	Path string `yaml:"path"`
	User string `yaml:"user"`
	Password string	`yaml:"password"`
}

type Postgres struct {
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	DataLocation string	`yaml:"data_location"`
}

type Utils struct {
	Dump string `yaml:"dump"`
	Psql string `yaml:"psql"`
	Compress string `yaml:"compress"`
}

func (c Config) GetSender() notice.Sender {
	return email.NewSmptClient(
		c.Email.User,
		c.Email.Password,
		c.Email.SmtpHost,
		c.Email.SmtpPort)
}

func LoadSetting(path string) (config Config, err error) {

	config = Config{}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	err = d.Decode(&config)

	return config, err
}
