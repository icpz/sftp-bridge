package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/icpz/sftp-bridge/common"
	"github.com/icpz/sftp-bridge/config"
	"github.com/icpz/sftp-bridge/sshtun"
)

var (
	cfgFile = flag.String("config", common.DefaultConfigFile, "path to config file")
	cfg *config.Config = nil
)

func handle(conn net.Conn) {
	defer conn.Close()

	stconn, err := sshtun.Dial(cfg.SshHost)
	if err != nil {
		log.Printf("[main] failed to connect to %s: %v\n", cfg.SshHost, err)
		return
	}
	defer stconn.Close()

	common.Relay(conn, stconn)
}

func mainLoop(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("[main] failed to accept: %v\n", err)
			break
		}
		log.Printf("[main] new conn from %s\n", conn.RemoteAddr().String())

		go handle(conn)
	}
}

func init() {
	flag.Parse()

	var err error
	cfg, err = config.ReadConfig(*cfgFile)
	if err != nil {
		log.Fatalln("[main] failed to load config: %v\n", err)
	}

	if cfg.SshHost == "" {
		log.Fatalln("[main] please specify ssh host")
	}
}

func main() {
	common.Init()
	log.Printf("[main] using TmpDir %s\n", common.TmpDir)
	defer common.DeInit()

	lis, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(cfg.ListenPort)))
	if err != nil {
		log.Printf("[main] failed to listen on port %d: %v\n", cfg.ListenPort, err)
		return
	}
	defer lis.Close()
	log.Printf("[main] listening on port %d\n", cfg.ListenPort)

	go mainLoop(lis)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
