package surlane

import (
	"net"
	"fmt"
	"strconv"
)

type TcpServer struct {
	Name string
	Port int
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
		select {
		case <-ctx.Done():
			available = false
			listener.Close()
		}
	}()
	for {
		if !available {
			return
		}
		conn, err := listener.Accept()
		if err !=nil {
			fmt.Println("test run", server.Name, "handle conn error:", err)
			continue
		}
		go server.Handler(conn)
	}
}

func RunClient(ctx *LocalContext, config ClientConfig) {
	TcpServer{
		"surlane-client",
		config.Port,
		func(conn net.Conn) {
			cpn++
			CreateClientPipe(NewContext(ctx, "client pipe " + strconv.Itoa(cpn)), config, conn)
		},
	}.Run(ctx)
}

func RunServer(ctx *LocalContext, config ServerConfig, dialServer func(*LocalContext, string)(net.Conn, error)) {
	TcpServer{
		"surlane-server",
		config.Port,
		func(conn net.Conn) {
			spn++
			CreateServerPipe(NewContext(ctx, "server pipe " + strconv.Itoa(spn)), config, conn, dialServer)
		},
	}.Run(ctx)
}
