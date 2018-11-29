package surlane

import (
	"net"
	"io"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
)

func LaneAck(ctx *LocalContext, conn net.Conn, rawAddr []byte, iv []byte) (err1 error) {
	if _, err := conn.Write(append(iv, rawAddr...)); err != nil {
		err1 = errors.Wrap(err, "")
		return
	}
	return nil
}

func LaneAuth(ctx *LocalContext, config Config, conn net.Conn) (iv []byte, address string, err1 error) {
	ivLen := GetIvLen(config)
	buffer := BufferPool.Borrow()
	defer BufferPool.GetBack(buffer)
	if _, err := io.ReadFull(conn, buffer[:ivLen+1]); err != nil {
		err1 = errors.Wrap(err, "")
		return
	}
	iv = buffer[:ivLen]
	var addrLen, addrType int
	addrType = int(buffer[ivLen])
	switch addrType {
	case AddrTypeIP4:
		addrLen = 4
	case AddrTypeIP6:
		addrLen = 16
	default:
		addrLen = addrType
	}
	if _, err := io.ReadFull(conn, buffer[ivLen+1:addrLen+ivLen+3]); err != nil {
		err1 = errors.Wrap(err, "")
		return
	}
	switch addrType {
	case AddrTypeIP4, AddrTypeIP6:
		address = net.IP(buffer[ivLen+1:addrLen+ivLen+1]).String()
	default:
		address = string(buffer[ivLen+1:addrLen+ivLen+1])
	}
	address += fmt.Sprintf(":%d", binary.BigEndian.Uint16(buffer[ivLen+addrLen+1:ivLen+addrLen+3]))
	return
}
