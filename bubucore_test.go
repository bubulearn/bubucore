package bubucore_test

import (
	"github.com/bubulearn/bubucore"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetServiceInfo(t *testing.T) {
	info := bubucore.GetServiceInfo()
	assert.Equal(t, bubucore.Opt.APIVersion, info.Version)
	assert.Equal(t, bubucore.Opt.APIBasePath, info.BasePath)
}

func TestOptions_GetHostname(t *testing.T) {
	opt := &bubucore.Options{
		Hostname: "test_hostname",
	}
	assert.Equal(t, "test_hostname", opt.GetHostname())

	osHost, _ := os.Hostname()

	opt = &bubucore.Options{}
	assert.Equal(t, osHost, opt.GetHostname())
}

func TestReadConfig(t *testing.T) {
	bubucore.Opt.ConfigFilePath = ".test.env"
	bubucore.Opt.ConfigFileType = "env"

	conf, err := bubucore.ReadConfig()
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "test1", conf.GetString("test_env_var"))
	assert.Equal(t, "test2", conf.GetString("OTHER_TEST_ENV_VAR"))

	bubucore.Opt.ConfigFilePath = "__unkown_file__.env"
	_, err = bubucore.ReadConfig()
	assert.Error(t, err)
}
