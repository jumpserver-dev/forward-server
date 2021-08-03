package main

import (
	"flag"
	"github.com/sevlyar/go-daemon"
	"log"
)

var (
	daemonFlag bool
)

func init() {
	flag.BoolVar(&daemonFlag, "d", false, "start as Daemon")
}

func main() {
	flag.Parse()
	if daemonFlag {
		ctx := &daemon.Context{
			PidFileName: "/tmp/forward-server.pid",
			PidFilePerm: 0644,
			LogFileName: "forward-server.log",
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
