package surlane

import "time"

type Config struct {
	Password string
	Method   int16
	Port     uint16
	Timeout  time.Duration
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
