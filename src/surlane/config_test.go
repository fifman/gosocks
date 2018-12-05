package surlane

import (
	"testing"
	"fmt"
)

func TestLoadConfig(t *testing.T) {
	//config := CreateDefaultConfig()
	config := Config{}
	loadConfigFile("./config_sample.toml", &config)
	//assert.Equal(t, ":9000", config.Kcptun.Listen)
	fmt.Println(config)
}
