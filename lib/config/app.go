package config

type App struct {
	Folders struct {
		Templates string `yaml:"templates"`
		Queries string `yaml:"queries"`
	} `yaml:"folders"`
}

