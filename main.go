package main

import (
	"log"

	"github.com/vilamslep/psql.maintenance/lib/config"
)

func main() {
	c, err := config.LoadSetting("setting.yaml")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(c)
}