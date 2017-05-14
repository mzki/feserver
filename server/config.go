package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
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

	// Wait time for the requesting, in second.
	WaitSecond int
	waitTime   time.Duration
}

func (conf *Config) validates() error {
	// check Port
	if p := conf.Port; p == "" {
		return fmt.Errorf("Config: empty Port, must be exist.")
	}
	// check WaitSecond
	if ws := conf.WaitSecond; ws < 0 {
		return fmt.Errorf("Config: incorrect WaitSecond %d, must be positive.", ws)
	} else {
		conf.waitTime = time.Duration(ws) * time.Second
	}
	return nil
}

const (
	DefaultURL  = "localhost"
	DefaultPort = "8080"

	DefaultWaitSecond = 2
)

var DefaultConfig = Config{
	URL:        DefaultURL,
	Port:       DefaultPort,
	WaitSecond: DefaultWaitSecond,
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
