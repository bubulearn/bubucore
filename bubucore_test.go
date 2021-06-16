package bubucore_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/bubulearn/bubucore"
	"testing"
)

func TestGetServiceInfo(t *testing.T) {
	info := bubucore.GetServiceInfo()
	assert.Equal(t, bubucore.Opt.APIVersion, info.Version)
	assert.Equal(t, bubucore.Opt.APIBasePath, info.BasePath)
}
