package main

const (
	ErrGateWay = "ErrGateway"
	ErrParams  = "ErrParams"
	ErrListen  = "ErrListen"
)

type Response struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Err  string `json:"error"`
	Msg  string `json:"message"`
}

const (
	MsgOk = "ok"
)

func NewSuccessResponse(f *Forward) Response {
	return Response{
		ID:   f.ID,
		Host: f.Host,
		Port: f.Port,
		Msg:  MsgOk,
	}
}

func NewErrResponse(err, msg string) Response {
	return Response{
		Err: err,
		Msg: msg,
	}
}
