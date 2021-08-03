package main

import (
	"flag"
	"log"

	"github.com/sevlyar/go-daemon"
)

var (
	daemonFlag bool
	logPath    string
)

func init() {
	flag.BoolVar(&daemonFlag, "d", false, "start as Daemon")
	flag.StringVar(&logPath, "l", "/tmp/forward-server.log", "log file path")
}

func main() {
	flag.Parse()
	if daemonFlag {
		ctx := &daemon.Context{
			PidFileName: "/tmp/forward-server.pid",
			PidFilePerm: 0644,
			LogFileName: logPath,
			Umask:       027,
			WorkDir:     "./",
		}
		d, err := ctx.Reborn()
		if err != nil {
			log.Fatalf("run failed: %v", err)
		}
		if d != nil {
			return
		}
		defer ctx.Release()
		runServer()
	} else {
		runServer()
	}

}
