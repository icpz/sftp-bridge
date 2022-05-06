package sshtun

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/icpz/sftp-bridge/common"
)

type STConn struct {
	net.Conn
	ssh *exec.Cmd
}

func (s *STConn) Close() error {
	s.Conn.Close()
	s.ssh.Process.Signal(os.Interrupt)
	s.ssh.Wait()
	log.Printf("[sshtun.Close] STConn closed")
	return nil
}

func waitForUnix(usock string, timeout time.Duration) error {
	start := time.Now()

	for {
		if _, err := os.Stat(usock); !os.IsNotExist(err) {
			break
		}

		log.Printf("[sshtun] unix sock %s does not exist yet...\n", usock)
		time.Sleep(500 * time.Millisecond)

		if time.Now().Sub(start) > timeout {
			return fmt.Errorf("ssh forward timeout")
		}
	}
	return nil
}

func Dial(sshHost, forwardTarget string) (*STConn, error) {
	conn := &STConn{}

	uaddr := filepath.Join(common.TmpDir, fmt.Sprintf("%s-%x.sock", sshHost, (uintptr(unsafe.Pointer(conn))/8)&0xffff))
	log.Printf("[sshtun.Dial] will use internal unix sock %s\n", uaddr)

	os.Remove(uaddr)
	cmd := exec.Command("ssh", "-N", "-L", fmt.Sprintf("%s:%s", uaddr, forwardTarget), sshHost)
	err := cmd.Start()
	if err != nil {
		log.Printf("[sshtun.Dial] failed to start cmd: %v\n", err)
		return nil, err
	}

	err = waitForUnix(uaddr, 10*time.Second)
	if err != nil {
		cmd.Process.Signal(os.Interrupt)
		cmd.Wait()
		return nil, err
	}

	uc, err := net.Dial("unix", uaddr)
	if err != nil {
		log.Printf("[sshtun.Dial] failed to dial %s: %v\n", sshHost, err)
		cmd.Process.Signal(os.Interrupt)
		cmd.Wait()
		return nil, err
	}
	conn.Conn = uc
	conn.ssh = cmd
	return conn, nil
}
