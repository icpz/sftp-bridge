package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Instance struct {
	SshHost    string   `json:"ssh-host,omitempty"`
	ListenPort int      `json:"listen-port,omitempty"`
	TargetIP   string   `json:"target-ip,omitempty"`
	TargetPort int      `json:"target-port,omitempty"`
}

type Config struct {
	Instances  []*Instance `json:"instances"`
}

func ReadConfig(cfgFile string) (*Config, error) {
	cfg := &Config{}

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
	for _, ins := range cfg.Instances {
		if ins.TargetIP == "" {
			ins.TargetIP = "127.0.0.1"
		}
		if ins.TargetPort == 0 {
			ins.TargetPort = 22
		}
	}
	return cfg, nil
}
