package surlane

import (
	"time"
	"net"
	"flag"
	"fmt"
	"strconv"
)

type Config struct {
	Password string
	Method   int
	Port     int
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

func setConfigFlags(config *Config) {
	flag.IntVar(&config.Method, "m", 0,
		"choose encryption method (default 0):\n" +
		"  0: CES_128_CFB\n" +
		"  1: CES_256_CFB\n")
	flag.IntVar(&config.Port, "p", 1180,
		"Set listening port")
	flag.StringVar(&config.Password, "s", "123456", "password for encryption, " +
		"please change it to a non-empty string!")
	flag.DurationVar(&config.Timeout, "d", time.Second * 600, "timeout, parsed by go with time.ParseDuration(). " +
		"Check go docs for details")
}

func GetServerConfig() ServerConfig {
	config := ServerConfig{ Config{} }
	setConfigFlags(&config.Config)
	flag.Parse()
	fmt.Println("listening with port: " + strconv.Itoa(config.Port))
	fmt.Println("choose method: " + strconv.Itoa(config.Method))
	return config
}

func GetClientConfig() ClientConfig {
	config := ClientConfig{Config{}, "" }
	setConfigFlags(&config.Config)
	flag.StringVar(&config.Server, "v", "", "the address of the surlane server, will be used by go with dial()")
	flag.Parse()
	if config.Server == "" {
		panic("server address (-v) flag required")
	}
	fmt.Println("server address is: " + config.Server)
	fmt.Println("listening with port: " + strconv.Itoa(config.Port))
	fmt.Println("choose method: " + strconv.Itoa(config.Method))
	return config
}
