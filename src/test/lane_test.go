package test

import (
	"testing"
	"github.com/fifman/gosocks/src/surlane"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestLane(t *testing.T) {
	err := runLaneTest("www.baidu.com", 1099, surlane.Config{
		"123456",
		surlane.CES_128_CFB,
		1999,
		time.Second * 5,
	})
	ERROR.Println(err)
	assert.Nil(t, err)
}
