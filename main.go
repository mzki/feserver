package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"

	"github.com/mzki/feserver/server"
)

const pkgPath = "github.com/mzki/feserver"

const defaultConfFile = "config.toml"

var confPath string

func init() {
	flag.StringVar(&confPath, "config", "", "path for the server config file")
	flag.Parse()
}

func main() {
	if confPath == "" {
		// get the directory of the feserver repository under GOPATH.
		// referenced from https://golang.org/x/tools/cmd/present/local.go.
		p, err := build.Default.Import(pkgPath, "", build.FindOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't find default config path: %v\n", err)
			fmt.Fprintf(os.Stderr, confPathMessage, pkgPath)
			os.Exit(1)
		}
		confPath = filepath.Join(p.Dir, defaultConfFile)
	}

	// load config file.
	conf, err := server.LoadConfigFile(confPath)
	if err != nil {
		log.Fatalf("FATAL: Invalid config file: %v", err)
	}

	// launch server process.
	if err := server.ListenAndServe(conf); err != nil {
		log.Fatalf("FATAL: %v", err)
	}
}

const confPathMessage = `
By default, feserver locates the server config file by 
looking for a %q package in your Go workspaces (GOPATH).
You may use the -config flag to specify an alternate location.
`
