package test

import (
	"net"
	"io"
	"github.com/fifman/surlane/src/surlane"
	"encoding/binary"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
)

type Socks5Client struct {
	net.Conn
}

func (client *Socks5Client) validate(version, nmethod byte, methods []byte) (err error) {
	header := append([]byte{version, nmethod}, methods...)
	if _, err = client.Write(header); err != nil {
		return
	}
	buffer := make([]byte, 2)
	if _, err = io.ReadFull(client, buffer); err != nil {
		return
	}
	if buffer[0] != surlane.Socks5Version {
		return errors.WithStack(surlane.VersionError)
	}
	if buffer[1] != surlane.Socks5Method {
		return errors.WithStack(surlane.MethodError)
	}
	return
}

func (client *Socks5Client) request(version, command, rsv, atyp byte, addr string, port uint16) (err error) {
	buffer := make([]byte, 1000)
	copy(buffer[:4], []byte{version, command, rsv, atyp})
	var addrLen byte
	switch atyp {
	case surlane.Socks5AtypIP4:
		addrLen = 4
		copy(buffer[4:8], net.ParseIP(addr).To4())
	case surlane.Socks5AtypIP6:
		addrLen = 16
		copy(buffer[4:20], net.ParseIP(addr).To16())
	default:
		addrBytes := []byte(addr)
		addrLen = byte(len(addrBytes))
		buffer[4] = addrLen
		copy(buffer[5:5+addrLen], addrBytes)
		addrLen = addrLen + 1
	}
	binary.BigEndian.PutUint16(buffer[4+addrLen:6+addrLen], port)
	if _, err = client.Write(buffer[:addrLen+6]); err != nil {
		return
	}
	if _, err = io.ReadFull(client, buffer[:10]); err != nil {
		return
	}
	if !bytes.Equal(buffer[:10], []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43}) {
		return surlane.ProtocolError
	}
	return
}

func (client *Socks5Client) send(content string) (err error) {
	if len(content) > 255 {
		return errors.New("content length is longer than 255")
	}
	_, err = client.Write(append([]byte{byte(len(content))}, []byte(content)...))
	return
}

func (client *Socks5Client) Run(content, addr string, port uint16, version, nmethod, command, rsv, atyp byte, methods []byte) (err error) {
	if err = client.validate(version, nmethod, methods); err != nil {
		return
	}
	if err = client.request(version, command, rsv, atyp, addr, port); err != nil {
		return
	}
	err = client.send(content)
	return
}

func NewSocks5Client(port int) *Socks5Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Printf("%+v\n\n", errors.Wrap(err, "socks5 client dial error"))
		return nil
	}
	return &Socks5Client{conn}
}
