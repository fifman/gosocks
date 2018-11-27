package test

import (
	"net"
	"fmt"
	"github.com/fifman/gosocks/src/surlane"
	"bytes"
	"github.com/pkg/errors"
	"strconv"
)

func RunLaneServer(parent *surlane.LocalContext, result chan interface{}, signChan chan interface{}, listener net.Listener, config surlane.Config) {
	ctx := surlane.NewContext(parent, "laneServer")
	defer listener.Close()
	conn, err := listener.Accept()
	deadline, ok := ctx.Deadline()
	if ok {
		conn.SetDeadline(deadline)
	}
	if err != nil {
		result <- err
		signChan <- struct {}{}
		return
	}
	address, iv, err := handleLane(ctx, config, conn)
	result <- err
	result <- address
	result <- iv
	signChan <- struct {}{}
	return
}

func handleLane(ctx *surlane.LocalContext, config surlane.Config, conn net.Conn) (address string, iv []byte, err error) {
	defer conn.Close()
	iv, address, err = surlane.LaneAuth(ctx, surlane.ServerConfig{config}, conn)
	fmt.Println("server iv is: ", iv)
	fmt.Println("server iv address: ", address)
	return
}

func runLaneTest(address string, port uint16, config surlane.Config) (error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.Port))
	if err != nil {
		return err
	}
	client := NewLaneClient(config)
	serverChan := make(chan interface{}, 3)
	clientChan := make(chan interface{}, 2)
	signChan := make(chan interface{}, 2)
	ctx := surlane.NewContextWithDeadline(&surlane.RootContext, "LaneTest", config.Timeout)
	go RunLaneServer(ctx, serverChan, signChan, listener, config)
	go client.Run(ctx, clientChan, signChan, config, address, port)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case _, _ = <- signChan:
	}
	if err1 := <-serverChan; err1 != nil {
		return err1.(error)
	}
	address_result := (<-serverChan).(string)
	iv1 := (<-serverChan).([]byte)
	if err1 := <-clientChan; err1 != nil {
		return err1.(error)
	}
	iv := (<-clientChan).([]byte)
	if bytes.Compare(iv, iv1) != 0 {
		fmt.Println("iv is: ",  iv)
		fmt.Println("iv1 is: ",  iv1)
		return errors.New("iv comparison wrong!")
	}
	if address_result != address + ":" + strconv.Itoa(int(port)) {
		fmt.Println("address: ", address + ":" + strconv.Itoa(int(port)))
		fmt.Println("result_address: ", address_result)
		return errors.New("address parsing wrong!")
	}
	return nil
}
