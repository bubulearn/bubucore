package app

import (
	"github.com/bubulearn/bubucore"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_SetFromViper(t *testing.T) {
	bubucore.Opt.ConfigFilePath = "../.test.env"
	v, err := bubucore.ReadConfig()
	if !assert.NoError(t, err) {
		return
	}

	conf := &Config{}
	conf.SetFromViper(v)

	assert.Equal(t, "80", conf.Port)

	assert.Equal(t, "http://localhost/notifications", conf.NotificationsHost)
	assert.Equal(t, "notifications-token", conf.NotificationsToken)

	assert.Equal(t, "localhost:6379", conf.RedisHost)
	assert.Equal(t, "localhost:27017", conf.MongoHost)

	conf.ApplyToGlobals()

	assert.Equal(t, []byte("12345"), bubucore.Opt.JWTPassword)
	assert.Equal(t, log.InfoLevel, log.GetLevel())
}
