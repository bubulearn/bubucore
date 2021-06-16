package bubucore_test

import (
	"github.com/bubulearn/bubucore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetServiceInfo(t *testing.T) {
	info := bubucore.GetServiceInfo()
	assert.Equal(t, bubucore.Opt.APIVersion, info.Version)
	assert.Equal(t, bubucore.Opt.APIBasePath, info.BasePath)
}
