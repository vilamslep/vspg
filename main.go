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


///compress
// import subprocess
// from configuration import config

// def compress_dir(src:str, dst)->bool:
//     utill = config.compress_tool()
    
//     exit_code = subprocess.Popen( [ utill, 'a', '-tzip', '-mx5', dst, src],stdout=subprocess.PIPE).wait()

//     return int(exit_code) == 0


