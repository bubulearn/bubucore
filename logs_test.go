package bubucore_test

import (
	"github.com/bubulearn/bubucore"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitLogs(t *testing.T) {
	bubucore.Opt.LogLevelDft = log.WarnLevel
	bubucore.InitLogs()
	assert.Equal(t, log.WarnLevel, log.GetLevel())
	log.SetLevel(log.DebugLevel)
}
