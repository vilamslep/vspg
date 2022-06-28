package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App `yaml:"app"`
	Postgres `yaml:"postgres"`
	Folder `yaml:"target_folder"`
	Utils `yaml:"utils"`
	Email `yaml:"email"`
	Schedule `yaml:"schedules"` 
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