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
	"github.com/icpz/sftp-bridge/sshtun"
)

var (
	shost = flag.String("shost", "", "ssh tunnel host")
	lport = flag.Int("lport", 8123, "local bridge port")
)

func handle(conn net.Conn) {
	defer conn.Close()

	stconn, err := sshtun.Dial(*shost)
	if err != nil {
		log.Printf("[main] failed to connect to %s: %v\n", *shost, err)
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
		log.Printf("[main] new conn from %s\n", conn.LocalAddr().String())

		go handle(conn)
	}
}

func init() {
	flag.Parse()

	if *shost == "" {
		log.Fatalln("[main] please specify shost")
	}
}

func main() {
	common.Init()
	log.Printf("[main] using TmpDir %s\n", common.TmpDir)
	defer common.DeInit()

	lis, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(*lport)))
	if err != nil {
		log.Printf("[main] failed to listen on port %d: %v\n", *lport, err)
		return
	}
	defer lis.Close()

	go mainLoop(lis)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
