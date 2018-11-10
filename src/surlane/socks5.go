package surlane

import (
	"net"
	"io"
	"github.com/pkg/errors"
)

const (
	Socks5Version = 5
	Socks5Method = 0
	Socks5AtypIP4 = 1
	Socks5AtypIP6 = 4
	Socks5AtypHost = 3
	Socks5Command = 1
	Socks5RSV = 0
	AddrTypeIP4 = 0
	AddrTypeIP6 = 1
)


func Socks5Auth(ctx *LocalContext, conn net.Conn) ([]byte, error) {
	buffer := BufferPool.Borrow()
	defer BufferPool.GetBack(buffer)
	//buffer := make([]byte, 4000)
	ctx.Trace("server mth vali")
	if err := methodValidation(ctx, conn, buffer); err != nil {
		ctx.Trace("vali end")
		confirmError(conn)
		return nil, errors.WithStack(err)
	}
	ctx.Trace("vali end")
	if err := confirm(conn); err != nil {
		return nil, errors.WithStack(err)
	}
	rawAddr, err := parseRequest(ctx, conn, buffer)
	if err != nil {
		switch err {
		case ProtocolError:
			reply(conn, 0x02)
		default:
			reply(conn, 0x01)
		}
		return nil, errors.WithStack(err)
	}
	if err := response(conn); err != nil {
		return nil, errors.WithStack(err)
	}
	return rawAddr, nil
}

var (
	VersionError = errors.New("Not a socks5 version")
	MethodError = errors.New("Invalid auth method")
	ProtocolError = errors.New("Wrong Protocol")
)

func confirmError(conn net.Conn) {
	conn.Write([]byte{Socks5Version, 0xFF})
}

func methodValidation(ctx *LocalContext, conn net.Conn, buffer []byte) error {
	if _, err := io.ReadFull(conn, buffer[:2]); err != nil {
		return err
	}
	if buffer[0] != Socks5Version {
		return VersionError
	}
	nmethod := buffer[1]
	if _, err := io.ReadFull(conn, buffer[2:2+nmethod]); err != nil {
		return err
	}
	ctx.Debug(buffer)
	for methodByte := range buffer[2:2+nmethod] {
		if methodByte == Socks5Method {
			return nil
		}
	}
	return MethodError
}

func confirm(conn net.Conn) (err error) {
	_, err = conn.Write([]byte{Socks5Version,Socks5Method})
	return
}

func parseRequest(ctx *LocalContext, conn net.Conn, buffer []byte) (rawAddr []byte, err error) {
	if _, err = io.ReadFull(conn, buffer[:4]); err != nil {
		return
	}
	if buffer[0] != Socks5Version || buffer[1] != Socks5Command || buffer[2] != Socks5RSV {
		return nil, ProtocolError
	}
	var addrlen, addrType byte
	switch buffer[3] {
	case Socks5AtypIP4:
		addrlen = 4
		addrType = AddrTypeIP4
	case Socks5AtypHost:
		if _, err = io.ReadFull(conn, buffer[4:5]); err != nil {
			return
		}
		addrlen = buffer[4]
		addrType = addrlen
	case Socks5AtypIP6:
		addrlen = 16
		addrType = AddrTypeIP6
	default:
		return nil, ProtocolError
	}
	buffer[0] = addrType
	if _, err = io.ReadFull(conn, buffer[1:addrlen+3]); err != nil {
		return
	}
	rawAddr = make([]byte, addrlen+3)
	copy(rawAddr, buffer[:addrlen+3])
	ctx.Debug(rawAddr)
	return
}

func response(conn net.Conn) (err error) {
	return reply(conn, 0x00)
}

func reply(conn net.Conn, reply byte) (err error) {
	_, err = conn.Write([]byte{0x05, reply, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	return
}
