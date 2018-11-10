package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/binary"
	"fmt"
	"net"
)

func TestErrorCheck(t *testing.T) {
	a, b := sup()
	assert.Equal(t, 4, a)
	assert.Equal(t, nil, b)
}

func sub() (int, error) {
	return 3, nil
}

func sup() (a int, err error) {
	x, err := sub()
	a = x + 1
	return
}

func TestDefer(t *testing.T) {
	assert.Equal(t, 13, simDefer())
}

func simDefer() (a int) {
	condition := false
	x := 3
	defer func(param int) {
		if condition {
			a = 10 + param
		} else {
			a = 20 + param
		}
	}(x)
	condition = true
	x = 6
	return a
}

func TestBigEndian(t *testing.T) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, 443)
	println(fmt.Sprintf(":%d", binary.BigEndian.Uint16(b)))
}

func TestAdhoc(t *testing.T) {
	buffer := make([]byte, 10)
	println(buffer)
}

func TestBuffer2String(t *testing.T) {
	buffer := make([]byte, 10)
	ip := net.ParseIP("1.2.4.5").To4()
	copy(buffer[2:6], ip)
}