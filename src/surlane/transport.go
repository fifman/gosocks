package surlane

import (
	"net"
	"io"
	"time"
	"github.com/pkg/errors"
)

type Pipe struct {
	downConn net.Conn
	upConn net.Conn
	*LocalContext
}

func (pipe *Pipe) Run() {
	signChannel := make(chan interface{}, 2)
	go transfer(NewContext(pipe.LocalContext, "down channel"), signChannel, pipe.upConn, pipe.downConn)
	go transfer(NewContext(pipe.LocalContext, "up channel"), signChannel, pipe.downConn, pipe.upConn)
	for i:=0; i<2; i++ {
		select {
		case <-pipe.Done():
			return
		case <- signChannel:
		}
	}
}

func transfer(ctx *LocalContext, signChannel chan interface{}, src, dst net.Conn) {
	for once(ctx, src, dst) {}
	signChannel <- nil
}

func once(ctx *LocalContext, src, dst net.Conn) bool {
	buffer := BufferPool.Borrow()
	defer BufferPool.GetBack(buffer)
	n, err := src.Read(buffer)
	ctx.Debug("transfer bytes:", n, err)
	if n > 0 {
		if _, err := dst.Write(buffer[:n]); err != nil {
			ctx.logError(err, "once write wrong!")
			return false
		}
	}
	if err != nil {
		if err != io.EOF {
			ctx.logError(err, "once read wrong!")
			ctx.Cancel()
		}
		return false
	}
	return true
}

func CreateClientPipe (ctx *LocalContext, config ClientConfig, conn net.Conn) error {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(config.Timeout))
	rawAddr, err := Socks5Auth(ctx, conn);
	if err != nil {
		return errors.Wrap(err, "client pipe 1")
	}
	upRawConn, err := (&net.Dialer{}).DialContext(ctx, "tcp", config.Server)
	if err != nil {
		return errors.Wrap(err, "client pipe 2")
	}
	defer upRawConn.Close()
	iv := GetIV(config)
	if err = LaneAck(ctx, upRawConn, rawAddr, iv); err != nil {
		return errors.Wrap(err, "client pipe 3")
	}
	upConn, err := NewClientSecureConn(upRawConn, config, iv)
	if err != nil {
		return errors.Wrap(err, "client pipe 4")
	}
	pipe := Pipe{ conn, upConn, ctx}
	pipe.Run()
	return nil
}

func CreateServerPipe(ctx *LocalContext, config ServerConfig, conn net.Conn, dialServer func(*LocalContext, string)(net.Conn, error)) error {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(config.Timeout))
	iv, address, err := LaneAuth(ctx, config, conn)
	if  err != nil {
		return errors.Wrap(err, "server pipe 1")
	}
	downConn, err := NewServerSecureConn(conn, config, iv)
	if err != nil {
		return errors.Wrap(err, "server pipe 2")
	}
	defer downConn.Close()
	upConn, err := dialServer(ctx, address)
	if err != nil {
		return errors.Wrap(err, "server pipe 3")
	}
	pipe := Pipe{ downConn, upConn, ctx}
	pipe.Run()
	return nil
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

func NewClientSecureConn(conn net.Conn, config ClientConfig, iv []byte) (secureConn *SecureConn, err error) {
	cipher, err := NewSurCipher4Client(config, iv)
	if err != nil {
		return
	}
	secureConn = &SecureConn{
		cipher,
		conn,
	}
	return
}

func NewServerSecureConn(conn net.Conn, config ServerConfig, iv []byte) (secureConn *SecureConn, err error) {
	cipher, err := NewSurCipher4Server(config, iv)
	if err != nil {
		return
	}
	secureConn = &SecureConn{
		cipher,
		conn,
	}
	return
}