package main

import (
	"net"
	"github.com/fifman/gosocks/src/surlane"
	"fmt"
	"time"
)

func main() {
	surlane.TcpServer{
		Name: "surlane-test-server",
		Port: 1191,
		Handler: func(conn net.Conn) {
			conn.SetDeadline(time.Now().Add(time.Second * 300))
			defer conn.Close()
			fmt.Println("new handled client conn:", conn.RemoteAddr())
			buffer := make([]byte, 10)
			var err error
			for err == nil {
				_, err = conn.Read(buffer)
			}
			fmt.Println("read err:", err)
			_, err = conn.Write([]byte{1,2,3,4})
			fmt.Println("write err:", err)
		},
	}.Run(&surlane.RootContext)
}
