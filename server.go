package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ForwardManager struct {
	lns map[string]*Forward

	sync.Mutex
}

func (s *ForwardManager) AddForward(key string, forward *Forward) {
	s.Lock()
	defer s.Unlock()
	s.lns[key] = forward

}

func (s *ForwardManager) RemoveForward(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.lns, key)

}

func (s *ForwardManager) GetForward(key string) *Forward {
	s.Lock()
	defer s.Unlock()
	return s.lns[key]
}

func runServer(){
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.HideBanner = true
	forwardManager := &ForwardManager{lns: make(map[string]*Forward)}
	e.GET("/forward/:id", func(c echo.Context) error {
		id := c.Param("id")
		if ValidUUIDString(id) {
			if forward := forwardManager.GetForward(id); forward != nil {
				return c.JSON(http.StatusOK, NewSuccessResponse(forward))
			}
		}
		return c.JSON(http.StatusNotFound, NewErrResponse(ErrParams, "not found id"))
	})

	e.POST("/forward", func(c echo.Context) (err error) {
		req := new(RequestForward)
		if err = c.Bind(req); err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusBadRequest, NewErrResponse(ErrParams, err.Error()))
		}
		if req.RemoteAddr == "" {
			return c.JSON(http.StatusBadRequest, NewErrResponse(ErrParams, "no forward dst addr"))
		}
		sshOptions := make([]SSHClientOption, 0, 5)

		sshOptions = append(sshOptions, SSHClientHost(req.Host))
		sshOptions = append(sshOptions, SSHClientPort(req.Port))
		sshOptions = append(sshOptions, SSHClientUsername(req.Username))
		sshOptions = append(sshOptions, SSHClientPassword(req.Password))
		if req.Timeout > defaultTimeout {
			sshOptions = append(sshOptions, SSHClientTimeout(req.Timeout))
		}
		if req.PrivateKey != "" {
			if req.Passphrase == "" {
				sshOptions = append(sshOptions, SSHClientPassphrase(req.Password))
			} else {
				sshOptions = append(sshOptions, SSHClientPassphrase(req.Passphrase))
			}
			sshOptions = append(sshOptions, SSHClientPrivateKey(req.PrivateKey))
		}
		sshClient, err := NewSSHClient(sshOptions...)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusBadRequest, NewErrResponse(ErrGateWay, err.Error()))
		}
		id := UUID()
		forward := &Forward{
			ID:      id,
			client:  sshClient,
			logger:  e.Logger,
			dstAddr: req.RemoteAddr,
		}
		err = forward.Start()
		if err != nil {
			return c.JSON(http.StatusBadRequest, NewErrResponse(ErrListen, err.Error()))
		}
		forwardManager.AddForward(id, forward)
		e.Logger.Infof("Add forward(%s) id %s", forward, id)
		return c.JSON(http.StatusOK, NewSuccessResponse(forward))
	})

	e.DELETE("/forward/:id", func(c echo.Context) error {
		id := c.Param("id")
		if ValidUUIDString(id) {
			if forward := forwardManager.GetForward(id); forward != nil {
				forward.Stop()
				forwardManager.RemoveForward(id)
				e.Logger.Infof("Remove forward(%s) id %s", forward, id)
				return c.JSON(http.StatusOK, NewSuccessResponse(forward))
			}
		}
		return c.JSON(http.StatusNotFound, NewErrResponse(ErrParams, "not found id"))
	})

	e.Logger.Fatal(e.Start(":8088"))
}

type RequestForward struct {
	Host       string `json:"host" xml:"host" form:"host" query:"host"`
	Port       int    `json:"port" xml:"port" form:"port" query:"port"`
	Username   string `json:"username" xml:"username" form:"username" query:"username"`
	Password   string `json:"password" xml:"password" form:"password" query:"password"`
	PrivateKey string `json:"private_key" xml:"private_key" form:"private_key" query:"private_key"`
	Passphrase string `json:"passphrase" xml:"passphrase" form:"passphrase" query:"passphrase"`

	RemoteAddr string `json:"remote_addr" xml:"remote_addr" form:"remote_addr" query:"remote_addr"`

	Timeout int `json:"timeout" xml:"timeout" form:"timeout" query:"timeout"` // 超时时间
}
