package main

import (
	"local/feserver/server"
	"log"
)

func main() {
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
