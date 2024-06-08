package tests

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/stretchr/testify/assert"
)

func TestCheckIp(t *testing.T) {
	blackList := []string{"127.0.0.1", "128.1.*.*", "129.*.0.1"}
	assert := assert.New(t)
	assert.False(utils.CheckIpAddress(blackList, "127.0.0.1"), "false")
	assert.False(utils.CheckIpAddress(blackList, "128.1.0.2"), "false")
	assert.False(utils.CheckIpAddress(blackList, "128.1.256.2"), "false")
	assert.False(utils.CheckIpAddress(blackList, "129.128.0.1"), "false")
	assert.True(utils.CheckIpAddress(blackList, "158.5.8.1"), "true")
}
