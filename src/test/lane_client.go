package test

import (
"net"
"fmt"
	"github.com/fifman/gosocks/src/surlane"
)

type LaneClient struct {
	net.Conn
}

func (client *LaneClient) validate(ctx *surlane.LocalContext, config surlane.Config, address string, port uint16) (iv []byte, err error) {
	rawAddr := surlane.GenRawAddr(address, port)
	iv = surlane.GetIV(surlane.ClientConfig{config, "aaa"})
	fmt.Println("gen iv is: ", iv)
	fmt.Println("gen rawAddr is: ", rawAddr)
	err = surlane.LaneAck(ctx, client.Conn, rawAddr, iv)
	return
}

func (client *LaneClient) send(content string) (err error) {
	_, err = client.Write([]byte(content))
	return
}

func (client *LaneClient) Run(parent *surlane.LocalContext, clientChan chan interface{}, signChan chan interface{}, config surlane.Config, addr string, port uint16) {
	ctx := surlane.NewContext(parent, "laneclient")
	deadline, ok := ctx.Deadline()
	if ok {
		client.SetDeadline(deadline)
	}
	defer client.Close()
	iv, err := client.validate(ctx, config, addr, port)
	clientChan <- err
	clientChan <- iv
	signChan <- struct {}{}
}

func NewLaneClient(config surlane.Config) *LaneClient {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", config.Port))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return &LaneClient{conn}
}
