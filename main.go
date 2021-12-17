package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sevlyar/go-daemon"
)

var (
	daemonFlag bool
	logPath    string
	vFlag      bool
)

func init() {
	flag.BoolVar(&daemonFlag, "d", false, "start as Daemon")
	flag.BoolVar(&vFlag, "v", false, "forward-server version")
	flag.StringVar(&logPath, "l", "/tmp/forward-server.log", "log file path")
}

func main() {
	flag.Parse()
	if vFlag {
		PrintVersion()
		return
	}
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

func PrintVersion() {
	fmt.Printf("Version:             %s\n", Version)
	fmt.Printf("Git Commit Hash:     %s\n", GitHash)
	fmt.Printf("UTC Build Time :     %s\n", BuildStamp)
	fmt.Printf("Go Version:          %s\n", GoVersion)
}
