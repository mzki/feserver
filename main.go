package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
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
	conf, err := loadConfig(confPath)
	if err != nil {
		log.Println(err)
		log.Println("\nCan not find any config path.\nuse builtin config insteadly")
		conf = &server.DefaultConfig
	}

	// launch server process.
	if err := server.ListenAndServe(conf); err != nil {
		log.Fatalf("FATAL: %v", err)
	}
}

func loadConfig(confPath string) (*server.Config, error) {
	if confPath == "" {
		// get the directory of the feserver repository under GOPATH.
		// referenced from https://golang.org/x/tools/cmd/present/local.go.
		p, err := build.Default.Import(pkgPath, "", build.FindOnly)
		if err != nil {
			return nil, fmt.Errorf(
				"Couldn't find default config path: %v\n"+confPathMessage,
				err,
				pkgPath,
			)
		}
		confPath = filepath.Join(p.Dir, defaultConfFile)
	}
	return server.LoadConfigFile(confPath)
}

const confPathMessage = `
By default, feserver locates the server config file by 
looking for a %q package in your Go workspaces (GOPATH).
You may use the -config flag to specify an alternate location.
`
