package surlane

import (
	"net"
	"fmt"
)

type TcpServer struct {
	Name string
	Port uint16
	Handler func(conn net.Conn)
}

func (server TcpServer) Run(ctx *LocalContext) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Port))
	if err != nil {
		fmt.Println("test run", server.Name, "listener error:", err)
		return
	}
	available := true
	go func() {
		for ;available; {
			conn, err := listener.Accept()
			if err !=nil {
				fmt.Println("test run", server.Name, "handle conn error:", err)
				continue
			}
			go server.Handler(conn)
		}
	}()
	select {
	case <-ctx.Done():
		available = false
		listener.Close()
	}
}

func RunClient(ctx *LocalContext, config ClientConfig) {
	TcpServer{
		"surlane-client",
		config.Port,
		func(conn net.Conn) {
			ctx := NewContext(ctx, "client conn handler")
			if err := CreateClientPipe(ctx, config, conn); err != nil {
				ctx.Errorf("handle client error: %+v\n\n", err)
			}
		},
	}.Run(ctx)
}

func RunServer(ctx *LocalContext, config ServerConfig, dialServer func(*LocalContext, string)(net.Conn, error)) {
	TcpServer{
		"surlane-server",
		config.Port,
		func(conn net.Conn) {
			ctx := NewContext(ctx, "server conn handler")
			if err := CreateServerPipe(ctx, config, conn, dialServer); err != nil {
				ctx.Errorf("handle server error: %+v\n\n", err)
			}
		},
	}.Run(ctx)
}
