package test

import (
	"net"
	"fmt"
	"github.com/fifman/gosocks/src/surlane"
	"io"
	"time"
	"context"
	"github.com/pkg/errors"
)

func RunServer(textLen, port int, clientRun func()error) (address, content string, port1 uint16, err error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	result := make(chan string, 2)
	result2 := make(chan uint16, 1)
	if err != nil {
		return
	}
	go func(){
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		address1, content1, port2 := handle(textLen, conn)
		duration, _ := time.ParseDuration("3s")
		conn.SetDeadline(time.Now().Add(duration))
		result <- address1
		result <- content1
		result2 <- port2
	}()
	if err = clientRun(); err != nil {
		return
	}
	address = <- result
	content = <- result
	port1 = <- result2
	return
}

func handle(textLen int, conn net.Conn) (address, content string, port uint16) {
	defer conn.Close()
	rawAddr, err := surlane.Socks5Auth(&surlane.LocalContext{context.TODO(), 1000, "handle socks5 conn", nil, nil,}, conn)
	if err != nil {
		fmt.Printf("%+v\n\n", errors.Wrap(err, "test handle error!"))
		return
	}
	address, port = surlane.ParseRawAddr(rawAddr)
	buffer := make([]byte, textLen)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		fmt.Printf("%+v\n\n", errors.Wrap(err, "test read error!"))
	}
	content = string(buffer)
	return
}

func runTest(content, addr string, port int, version, nmethod, command, rsv, atyp byte, methods []byte) (string, string, uint16, error) {
	return RunServer(len(content), port, func() error {
		client := NewSocks5Client(port)
		duration, _ := time.ParseDuration("3s")
		client.SetDeadline(time.Now().Add(duration))
		defer client.Close()
		return client.Run(content, addr, uint16(port), version, nmethod, command, rsv, atyp, methods)
	})
}

