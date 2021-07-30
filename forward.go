package main

import (
	"fmt"
	"io"
	"net"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/ssh"
)

type Forward struct {
	ID      string `json:"id"`
	client  *ssh.Client
	logger  echo.Logger
	dstAddr string

	ln   net.Listener
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (s *Forward) Start() error {
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return err
	}
	tcpAddr := ln.Addr().(*net.TCPAddr)
	s.Host = tcpAddr.IP.String()
	s.Port = tcpAddr.Port
	s.ln = ln
	go s.run()
	return nil

}
func (s *Forward) Stop() {
	if s.ln != nil {
		if err := s.ln.Close(); err != nil {
			s.logger.Error(err)
		}
	}
	if err := s.client.Close(); err != nil {
		s.logger.Error(err)
	}
}

func (s *Forward) String() string {
	return fmt.Sprintf("%s", s.dstAddr)
}

func (s *Forward) run() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			s.logger.Errorf("listen accept failed: %v", err)
			return
		}
		go s.forward(conn)
	}
}

func (s *Forward) forward(conn net.Conn) {
	defer conn.Close()

	proxyCon, err := s.client.Dial("tcp", s.dstAddr)
	if err != nil {
		s.logger.Errorf("ssh.Dial failed: %s\n", err)
		return
	}
	go func() {
		defer proxyCon.Close()
		if _, err = io.Copy(proxyCon, conn); err != nil {
			s.logger.Errorf("io.Copy local-> proxy err: %s\n", err)
		}
	}()
	if _, err = io.Copy(conn, proxyCon); err != nil {
		s.logger.Errorf("io.Copy proxy -> local err: %s\n", err)
	}
}

