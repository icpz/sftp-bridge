package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/icpz/sftp-bridge/bridge"
	"github.com/icpz/sftp-bridge/common"
	"github.com/icpz/sftp-bridge/config"
)

var (
	cfgFile     = flag.String("config", common.DefaultConfigFile, "path to config file")
	showVersion = flag.Bool("version", false, "show version")

	cfg *config.Config = nil
)

func init() {
	flag.Parse()

	if *showVersion {
		log.Printf("sftp-bridge %s\n", common.Version)
		os.Exit(0)
	}

	var err error
	cfg, err = config.ReadConfig(*cfgFile)
	if err != nil {
		log.Fatalln("[main] failed to load config: %v\n", err)
	}

	if len(cfg.Instances) == 0 {
		log.Fatalln("[main] please specify at least one instance")
	}
	for idx, ins := range cfg.Instances {
		if ins.SshHost == "" {
			log.Fatalf("[main] please specify ssh host for instance %d\n", idx)
		}
	}
}

func main() {
	common.Init()
	log.Printf("[main] using TmpDir %s\n", common.TmpDir)
	defer common.DeInit()

	for _, insCfg := range cfg.Instances {
		ins, err := bridge.StartInstance(insCfg)
		if err != nil {
			log.Fatalf("[main] failed to start instance: %v\n", err)
		}
		defer ins.Close()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
