package main

import (
	"net"
	"github.com/fifman/gosocks/src/surlane"
	"fmt"
	"io"
	"time"
)

func main() {
	surlane.TcpServer{
		Name: "surlane-test-server",
		Port: 1191,
		Handler: func(conn net.Conn) {
			//conn.SetDeadline(time.Now().Add(time.Second * 3000))
			defer conn.Close()
			time.Sleep(time.Second * 5)
			//fmt.Println("new handled client conn:", conn.RemoteAddr())
			buffer := make([]byte, 4)
			n, err := io.ReadFull(conn, buffer)
			if err != nil {
				fmt.Println("read wrong:", err)
			} else if n < 4 {
				fmt.Println("less bytes read!")
				return
			} else {
				fmt.Println("bytes:", buffer[:n])
			}
			/*
			_, err = conn.Write([]byte{1,2,3,4})
			if err != nil {
				fmt.Println("write err:", err)
			}
			_, err = conn.Read(buffer)
			if err != nil {
				fmt.Println("end read err:", err)
				return
			}
			*/
		},
	}.Run(&surlane.RootContext)
}
