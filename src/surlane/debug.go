package surlane

import (
	"net"
	"fmt"
	"strconv"
)

var (
	cpn = 0
	spn = 0
)

func logPipe(ctx LocalContext, local, remote net.Conn, msg string) {
	fmt.Println(fmt.Sprintf("%s-%d:%d", ctx.Name,
		getPort(local.RemoteAddr().String()),
		getPort(remote.LocalAddr().String()),
	), "###", msg)
}

func getPort(addr string) int {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return -1
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = -1
	}
	return port
}
