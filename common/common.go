package common

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

var (
	DefaultConfigFile = "./config.json"

	TmpDir string
)

func Init() {
	dir, err := ioutil.TempDir("", "sftp-bridge-*")
	if err != nil {
		log.Fatalf("[common] failed to create tmp dir: %v\n", err)
	}
	TmpDir = dir
}

func DeInit() {
	log.Printf("[common] cleanup %s\n", TmpDir)
	os.RemoveAll(TmpDir)
}

func Relay(left, right net.Conn) {
	ch := make(chan error)

	go func() {
		_, err := io.Copy(left, right)
		left.SetReadDeadline(time.Now())
		ch <- err
	}()

	io.Copy(right, left)
	right.SetReadDeadline(time.Now())
	<-ch
}
