package surlane

import (
	"time"
	"net"
	"flag"
	"fmt"
	"strconv"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Password string
	Method   int
	Port     int
	Timeout  time.Duration
	Server   string
}

func (config *Config ) ApplyTimeout(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(config.Timeout))
}

var (
	BufferPool = NewPool(4096, 4000)
)

func loadConfigFile(file string, config *Config) {
	if _, err := toml.DecodeFile(file, config); err != nil {
		panic(fmt.Sprint("config file cannot be parsed: ", err))
	}
}

func GetConfig() Config {
	config := Config{}
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
	flag.StringVar(&config.Server, "v", "", "the address of the surlane server, will be used by go with dial()")
	var configFile string
	flag.StringVar(&configFile, "c", "", "config file location")
	flag.Parse()
	if configFile != "" {
		loadConfigFile(configFile, &config)
	}
	fmt.Println("server address is: " + config.Server)
	fmt.Println("listening with port: " + strconv.Itoa(config.Port))
	fmt.Println("choose method: " + strconv.Itoa(config.Method))
	return config
}
