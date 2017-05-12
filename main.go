package main

import (
	"log"

	"github.com/mzki/feserver/server"
)

func main() {
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
