package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	SshHost    string   `json:"ssh-host,omitempty"`
	ListenPort int      `json:"listen-port,omitempty"`
}

func ReadConfig(cfgFile string) (*Config, error) {
	cfg := &Config{
		ListenPort: 8123,
	}

	buf, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Printf("[config] failed to read %s: %v\n", cfgFile, err)
		return nil, err
	}

	err = json.Unmarshal(buf, cfg)
	if err != nil {
		log.Printf("[config] failed to parse config file %s: %v\n", cfgFile, err)
		return nil, err
	}
	return cfg, nil
}
