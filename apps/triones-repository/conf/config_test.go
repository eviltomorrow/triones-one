package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {
	assert := assert.New(t)

	config, err := ReadConfigFile("./etc/config.toml")
	assert.Nil(err)
	t.Log(config.String())
}
