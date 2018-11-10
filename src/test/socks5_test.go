package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/fifman/gosocks/src/surlane"
	"fmt"
)

var (
	inputAddress = "www.test.com"
	inputContent = "aaa"
	inputPort = 9999
	inputMethodNo = byte(1)
	inputMethods = []byte{surlane.Socks5Method}
)

func TestSocksPass(t *testing.T) {
	runSuccessTest(t, inputContent, inputAddress, inputPort, surlane.Socks5AtypHost, []byte{1,5})
	runSuccessTest(t, inputContent, "127.3.4.5", inputPort, surlane.Socks5AtypIP4, []byte{1,5})
}

func TestSocksFail(t *testing.T) {
	runTestWithErrorVersion(t)
	runTestWithErrorCommand(t)
	surlane.DebugMark = 1
	runTestWithErrorMethod(t)
	runTestWithErrorMethodNo(t)
}

func runSuccessTest(t *testing.T, content, addr string, port int, atyp byte, methods []byte) {
	address, content1, port1, err := runTest(content, addr, port, surlane.Socks5Version, byte(len(methods)+1), surlane.Socks5Command, surlane.Socks5RSV, atyp, append([]byte{surlane.Socks5Method}, methods...))
	assert.Nil(t, err)
	assert.Equal(t, content, content1)
	assert.Equal(t, addr, address)
	assert.Equal(t, uint16(port), port1)
}

func runTestWithErrorVersion(t *testing.T) {
	_, _, _, err := runTest("test content", inputAddress, inputPort, surlane.Socks5Version-1, inputMethodNo, surlane.Socks5Command, surlane.Socks5RSV, surlane.AddrTypeIP4, inputMethods)
	fmt.Println(err)
	assert.NotNil(t, err)
}

func runTestWithErrorMethod(t *testing.T) {
	_, _, _, err := runTest("test content", inputAddress, inputPort, surlane.Socks5Version, inputMethodNo+1, surlane.Socks5Command, surlane.Socks5RSV, surlane.AddrTypeIP4, []byte{1})
	fmt.Println(err)
	assert.NotNil(t, err)
}

func runTestWithErrorMethodNo(t *testing.T) {
	_, _, _, err := runTest("test content", inputAddress, inputPort, surlane.Socks5Version, inputMethodNo+1, surlane.Socks5Command+1, surlane.Socks5RSV, surlane.AddrTypeIP4, inputMethods)
	fmt.Println(err)
	assert.NotNil(t, err)
}

func runTestWithErrorCommand(t *testing.T) {
	_, _, _, err := runTest("test content", inputAddress, inputPort, surlane.Socks5Version, inputMethodNo, surlane.Socks5Command+1, surlane.Socks5RSV, surlane.AddrTypeIP4, inputMethods)
	fmt.Println(err)
	assert.NotNil(t, err)
}