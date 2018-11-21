package surlane

import (
	"time"
	"net"
)

type Config struct {
	Password string
	Method   int16
	Port     uint16
	Timeout  time.Duration
}

func (config *Config ) ApplyTimeout(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(config.Timeout))
}

type ServerConfig struct {
	Config
}

type ClientConfig struct {
	Config
	Server string
}

var (
	BufferPool = NewPool(4096, 4000)
)
