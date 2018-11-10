package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFullRun(t *testing.T) {
	count := trySuccessRun("123456", "www.baidu.com", 443, 300)
	assert.Zero(t, count)
}
