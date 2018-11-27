package surlane

import (
	"net"
	"github.com/pkg/errors"
	"io"
)

type Pipe struct {
	downConn net.Conn
	upConn net.Conn
	config Config
	*LocalContext
}

func (pipe *Pipe) Run() {
	signChannel := make(chan interface{}, 2)
	go transfer(pipe.LocalContext, signChannel, pipe.upConn, pipe.downConn, pipe.config)
	go transfer(pipe.LocalContext, signChannel, pipe.downConn, pipe.upConn, pipe.config)
	for i:=0; i<2; i++ {
		select {
		case <-pipe.Done():
			return
		case <- signChannel:
		}
	}
}

func transfer(ctx *LocalContext, signChannel chan interface{}, src, dst net.Conn, config Config) {
	for {
		config.ApplyTimeout(src)
		//buffer := BufferPool.Borrow()
		//defer BufferPool.GetBack(buffer)
		buffer := make([]byte, 4096)
		n, err := src.Read(buffer)
		ctx.Debug("transfer bytes:", n, err)
		if n > 0 {
			if _, err2 := dst.Write(buffer[:n]); err2 != nil {
				ctx.LogError(err2, "once write wrong!")
				ctx.Cancel()
				return
			}
		}
		if err == nil {
            continue
		}
		if err != io.EOF && !CheckConnReset(err) {
			ctx.LogError(err, "once read wrong!")
			ctx.Cancel()
		} else {
			signChannel <- nil
		}
		return
	}
}

func CreateClientPipe (ctx *LocalContext, config ClientConfig, conn net.Conn) {
	defer conn.Close()
	config.ApplyTimeout(conn)
	rawAddr, err := Socks5Auth(ctx, conn)
	if err != nil {
		ctx.LogError(err, "client pipe 1")
		return
	}
	upRawConn, err := (&net.Dialer{}).DialContext(NewContext(ctx, "client dial"), "tcp", config.Server)
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
	pipe := Pipe{ conn, upConn, config.Config, ctx}
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
	pipe := Pipe{ downConn, upConn, config.Config, ctx}
	pipe.Run()
}

func DialWeb(ctx *LocalContext, address string)(net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "tcp", address)
}

type SecureConn struct {
	*SurCipher
	net.Conn
}

func (secureConn *SecureConn) Read(buffer []byte) (n int, err error) {
	n, err = secureConn.Conn.Read(buffer)
	if n > 0 {
		copy(secureConn.SurCipher.decrypt(buffer[:n]), buffer[:n])
	}
	return
}

func (secureConn *SecureConn) Write(buffer []byte) (n int, err error) {
	n, err = secureConn.Conn.Write(buffer)
	if n > 0 {
		copy(secureConn.SurCipher.encrypt(buffer[:n]), buffer[:n])
	}
	return
}

func NewClientSecureConn(conn net.Conn, config ClientConfig, iv []byte) (*SecureConn, error) {
	cipher, err := NewSurCipher4Client(config, iv)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &SecureConn{
		cipher,
		conn,
	}, nil
}

func NewServerSecureConn(conn net.Conn, config ServerConfig, iv []byte) (*SecureConn, error) {
	cipher, err := NewSurCipher4Server(config, iv)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &SecureConn{
		cipher,
		conn,
	}, nil
}