package bridge

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/icpz/sftp-bridge/common"
	"github.com/icpz/sftp-bridge/config"
	"github.com/icpz/sftp-bridge/sshtun"
)

type Instance struct {
	listener net.Listener
	config   config.Instance
	wait     *sync.WaitGroup
}

func (ins *Instance) Close() {
	ins.listener.Close()
	ins.wait.Wait()
	log.Printf("[bridge] instance closed\n")
}

func (ins *Instance) handleConn(conn net.Conn) {
	defer conn.Close()

	cfg := &ins.config
	stconn, err := sshtun.Dial(cfg.SshHost, net.JoinHostPort(cfg.TargetIP, strconv.Itoa(cfg.TargetPort)))
	if err != nil {
		log.Printf("[bridge] failed to connect to %s: %v\n", cfg.SshHost, err)
		return
	}
	defer stconn.Close()

	common.Relay(conn, stconn)
}

func (ins *Instance) mainLoop() {
	ins.wait.Add(1)
	for {
		conn, err := ins.listener.Accept()
		if err != nil {
			log.Printf("[bridge] failed to accept: %v\n", err)
			break
		}
		log.Printf("[bridge] new conn from %s -> %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())

		go ins.handleConn(conn)
	}
	ins.wait.Done()
}

func StartInstance(cfg *config.Instance) (*Instance, error) {
	lis, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(cfg.ListenPort)))
	if err != nil {
		log.Printf("[bridge] failed to listen on port %d: %v\n", cfg.ListenPort, err)
		return nil, err
	}
	log.Printf("[bridge] instance listening on port %d, host %s, target %s:%d\n", cfg.ListenPort, cfg.SshHost, cfg.TargetIP, cfg.TargetPort)

	ins := &Instance{
		listener: lis,
		config:   *cfg,
		wait:     &sync.WaitGroup{},
	}

	go ins.mainLoop()

	return ins, nil
}
