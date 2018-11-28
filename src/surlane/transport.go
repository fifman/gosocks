package surlane

import (
	"net"
	"io"
)

type Pipe struct {
	downConn net.Conn
	upConn net.Conn
	config Config
	channel chan interface{}
	*LocalContext
}

func (pipe *Pipe) Run() {
	pipe.channel = make(chan interface{}, 2)
	go pipe.transfer(pipe.upConn, pipe.downConn)
	go pipe.transfer(pipe.downConn, pipe.upConn)
	for i:=0; i<2; i++ {
		select {
		case <-pipe.Done():
			return
		case <- pipe.channel:
		}
	}
}

func (pipe *Pipe) transfer(src, dst net.Conn) {
	for ;pipe.once(src, dst); {}
}

func (pipe *Pipe) once(src, dst net.Conn) bool {
	pipe.config.ApplyTimeout(src)
	buffer := BufferPool.Borrow()
	defer BufferPool.GetBack(buffer)
	n, err := src.Read(buffer)
	if n > 0 {
		if _, err2 := dst.Write(buffer[:n]); err2 != nil {
			pipe.LogError(err2, "once write wrong!")
			pipe.Cancel()
			return false
		}
	}
	if err == nil {
		return true
	}
	if err != io.EOF && !CheckConnReset(err) {
		pipe.LogError(err, "once read wrong!")
		pipe.Cancel()
	} else {
		pipe.channel <- nil
	}
	return false
}

func CreateClientPipe (ctx *LocalContext, config ClientConfig, conn net.Conn) {
	defer conn.Close()
	config.ApplyTimeout(conn)
	rawAddr, err := Socks5Auth(ctx, conn)
	if err != nil {
		ctx.LogError(err, "client pipe 1")
		return
	}
	upRawConn, err := (&net.Dialer{}).DialContext(NewContextWithDeadline(ctx, "client dial", config.Timeout), "tcp", config.Server)
	if err != nil {
		ctx.LogError(err, "client pipe 2")
		return
	}
	config.ApplyTimeout(upRawConn)
	defer upRawConn.Close()
	iv := GetIV(config)
	config.ApplyTimeout(conn)
	if err = LaneAck(ctx, upRawConn, rawAddr, iv); err != nil {
		ctx.LogError(err, "client pipe 3")
		return
	}
	upConn, err := NewClientSecureConn(upRawConn, config, iv)
	if err != nil {
		ctx.LogError(err, "client pipe 4")
		return
	}
	pipe := Pipe{ conn, upConn, config.Config, make(chan interface{}, 2), ctx}
	pipe.Run()
}

func CreateServerPipe(ctx *LocalContext, config ServerConfig, conn net.Conn, dialServer func(*LocalContext, string)(net.Conn, error)) {
	defer conn.Close()
	config.ApplyTimeout(conn)
	iv, address, err := LaneAuth(ctx, config, conn)
	if  err != nil {
		ctx.LogError(err, "server pipe 1")
		return
	}
	downConn, err := NewServerSecureConn(conn, config, iv)
	if err != nil {
		ctx.LogError(err, "server pipe 2")
		return
	}
	defer downConn.Close()
	config.ApplyTimeout(downConn)
	upConn, err := dialServer(NewContextWithDeadline(ctx, "server dial", config.Timeout), address)
	if err != nil {
		ctx.LogError(err, "server pipe 3")
		return
	}
	pipe := Pipe{ downConn, upConn, config.Config, make(chan interface{}, 2), ctx}
	pipe.Run()
}

func DialWeb(ctx *LocalContext, address string)(net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "tcp", address)
}

