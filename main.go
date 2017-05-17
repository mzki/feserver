package main

import (
	"log"
	// "os"

	"github.com/mzki/feserver/server"
)

func main() {
	var conf *server.Config = nil
	// if _, err := os.Stat(server.ConfigFile); err == nil {
	// 	// file exists, load it
	// 	conf, err = server.LoadConfigFile(server.ConfigFile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	if err := server.ListenAndServe(conf); err != nil {
		log.Fatal(err)
	}
}
