package server

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mzki/feserver/src"
)

// default Server configuration file.
const ConfigFile = "config.toml"

// Configuration for server behavior.
// it must construct by LoadConfig() or LoadConfigFile().
type Config struct {
	// server port
	Port string

	// server URL
	URL string

	// Source location for get questions.
	Sources []Source
}

func (conf *Config) validates() error {
	// check Port
	if p := conf.Port; p == "" {
		return fmt.Errorf("Config: empty Port, must be exist.")
	}

	// check sources
	for _, s := range conf.Sources {
		if err := s.Source.ValidatesSelf(); err != nil {
			return fmt.Errorf("Config: invalid source: %v", err)
		}
		// check WaitSecond
		if ws := s.WaitSecond; ws < 0 {
			return fmt.Errorf("Config: incorrect WaitSecond %d, must be positive.", ws)
		}
	}
	return nil
}

const (
	DefaultURL  = "localhost"
	DefaultPort = "8080"

	DefaultWaitSecond = 2
)

var DefaultConfig = Config{
	URL:  DefaultURL,
	Port: DefaultPort,
	Sources: []Source{
		DefaultSource,
		FESource,
		APSource,
	},
}

// DefaultSource has F.E. examination source.
var DefaultSource = Source{
	SubAddr:    "",
	Source:     src.FE,
	WaitSecond: DefaultWaitSecond,
}

//	FESource has F.E. examination source.
var FESource = Source{
	SubAddr:    "/fe",
	Source:     src.FE,
	WaitSecond: DefaultWaitSecond,
}

//	APSource has A.P. examination source.
var APSource = Source{
	SubAddr:    "/ap",
	Source:     src.AP,
	WaitSecond: DefaultWaitSecond,
}

// Source is the source definition for getting questions.
type Source struct {
	src.Source

	// sub address.
	SubAddr string
	// Wait time for the requesting, in second.
	WaitSecond int
}

// it loads the configuration from file.
// it returns loaded config and load error.
func LoadConfigFile(file string) (*Config, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return LoadConfig(fp)
}

// it loads the configuration from io.Reader.
// it returns loaded config and load error.
func LoadConfig(r io.Reader) (*Config, error) {
	conf := &Config{}
	if err := decode(r, conf); err != nil {
		return nil, fmt.Errorf("LoadConfig: %v", err)
	}
	return conf, conf.validates()
}

// decode from reader and store it to data.
func decode(r io.Reader, data interface{}) error {
	meta, err := toml.DecodeReader(r, data)
	if undecoded := meta.Undecoded(); undecoded != nil && len(undecoded) > 0 {
		log.Println("Config.Decode:", "undecoded keys exist,", undecoded)
	}
	return err
}
