package config

import (
	"os"

	"github.com/vilamslep/vspg/pkg/notice"
	"github.com/vilamslep/vspg/pkg/notice/email"
	"gopkg.in/yaml.v2"
)

type Config struct {
	App      `yaml:"app"`
	Postgres 
	Folder   `yaml:"target_folder"`
	Utils    `yaml:"utils"`
	Email    `yaml:"email"`
	Schedule `yaml:"schedules"`
}

type App struct {
	SettingsFolders `yaml:"folders"`
}

type Postgres struct{}

func (p Postgres) GetUser() string {
	return os.Getenv("PGUSER")
}

func (p Postgres) GetPassword() string {
	return os.Getenv("PGPASSWORD")
}

func (p Postgres) GetServer() string {
	return os.Getenv("PGHOST")
}

func (p Postgres) GetPort() string {
	return os.Getenv("PGPORT")
}

func (p Postgres) GetDataLocation() string {
	return os.Getenv("PGDATA")
}

type Folder struct {
	Path     string `yaml:"path"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Utils struct {
	Dump     string `yaml:"dump"`
	Psql     string `yaml:"psql"`
	Compress string `yaml:"compress"`
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

type SettingsFolders struct {
	Templates string `yaml:"templates"`
	TempPath string `yaml:"tempate_place"`
}

type Letter struct {
	Subject string `yaml:"subject"`
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
