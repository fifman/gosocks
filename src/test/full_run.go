package test

import (
	"github.com/fifman/surlane/src/surlane"
	"net"
	"io"
	"time"
	"fmt"
	"strconv"
	"math/rand"
	"github.com/pkg/errors"
	"sync"
)

func runRemote(ctx *surlane.LocalContext) {
	surlane.TcpServer{
		"remote",
		1999,
		func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 4000)
			_ , err := io.ReadFull(conn, buf[0:1])
			if err != nil {
				fmt.Printf("receive error: %+v\n\n", errors.Wrap(err, ""))
				return
			}
			m := int(buf[0])
			surlane.RootContext.Level(surlane.LevelInfo)
			surlane.RootContext.Info("get the size of content: ", m)
			n, err := io.ReadFull(conn, buf[0:m])
			var resp []byte
			if err == nil {
				if n > 1 {
					for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
						buf[i], buf[j] = buf[j], buf[i]
					}
				}
				if n > 0 {
					resp = buf[:n]
				} else {
					resp = []byte("xxx")
				}
				_, err = conn.Write(resp)
				surlane.RootContext.Info("resp is:", string(resp), err)
				surlane.RootContext.Level(surlane.LevelError)
			} else {
				fmt.Printf("handle conn {%T} error: %+v\n", err, errors.Wrap(err, ""))
				return
			}
		},
	}.Run(ctx)
}

func runLocal(content, addr string, port uint16, version, nMethod, command, rsv, atyp byte, methods []byte) error {
	ctx := surlane.NewContext(&surlane.RootContext, "run local")
	ctx.Level(surlane.LevelDebug)
	client := NewSocks5Client(1977)
	client.SetDeadline(time.Now().Add(time.Second * 10))
	defer client.Close()
	err := client.Run(content, addr, port, version, nMethod,command, rsv, atyp, methods)
	if err != nil {
		return errors.Wrap(err, "")
	}
	buf := make([]byte, 4000)
	n, err := io.ReadFull(client, buf[:len(content)])
	ctx.Debug("client read result:", n)
	if err != nil {
		fmt.Printf("DEBUG: reply {%d} bytes: {%s}\n\n", n, string(buf[:n]))
		return errors.Wrap(err, "")
	}
	result := buf[:n]
	if n < 2 {
		return errors.New(fmt.Sprintf("result {%s} is different from reversed expected {%s}", result, content))
	}
	for i,j := 0, n-1; i < j; i,j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	if string(result) != content {
		return errors.New(fmt.Sprintf("reversed result {%s} is different from expected {%s}", string(result), content))
	}
	return nil
}

func runFullPipe(ctx *surlane.LocalContext, expectedAddress, pwd string, method int) {
	go runRemote(ctx)
	go surlane.RunServer(ctx, surlane.ServerConfig{
		surlane.Config{
			pwd,
			method,
			1988,
			time.Second * 60,
		},
	}, func(ctx *surlane.LocalContext, address string) (net.Conn, error) {
		if address != expectedAddress {
			return nil, errors.New(fmt.Sprintf("address {%s} is different from expected {%s}", address, expectedAddress))
		}
		return (&net.Dialer{}).DialContext(ctx, "tcp", fmt.Sprintf("127.0.0.1:%d", 1999))
	})
	go surlane.RunClient(ctx, surlane.ClientConfig{
		surlane.Config{
			pwd,
			method,
			1977,
			time.Second * 60,
		},
		fmt.Sprintf("127.0.0.1:%d", 1988),
	})
}

func trySuccessRun(pwd, address string, port uint16, num int) int {
	ctx := surlane.NewContext(&surlane.RootContext, "test context")
	var waiter sync.WaitGroup
	runFullPipe(ctx, address+":"+strconv.Itoa(int(port)), pwd, surlane.Ces128Cfb)
	errChan := make(chan error, num)
	time.Sleep(time.Second * 3)
	for i:=0; i<num; i++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			content := strconv.Itoa(rand.Intn(100000000) + 400)
			surlane.RootContext.Level(surlane.LevelInfo)
			surlane.RootContext.Info("content is: " + content)
			surlane.RootContext.Level(surlane.LevelError)
			errChan <- runLocal(content, address, port, surlane.Socks5Version,
				1, surlane.Socks5Command, surlane.Socks5RSV, surlane.Socks5AtypHost, []byte{surlane.Socks5Method})
		}()
	}
	waiter.Wait()
	close(errChan)
	count := 0
	for err := range errChan {
		if err != nil {
			count ++
			fmt.Printf("try run result {%T} error: %+v\n\n", errors.Cause(err), err)
		}
	}
	ctx.Cancel()
	return count
}
