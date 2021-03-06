package app

import (
	"github.com/bubulearn/bubucore"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
)

// i18nFileDft is a default i18n file path
const i18nFileDft = "./i18n.yml"

// Config is a basic Bubulearn service config
type Config struct {
	Port     string
	LogLevel log.Level

	CORSEnable    bool
	CORSAllowAll  bool
	CORSAllowCred bool
	CORSAllowWS   bool
	CORSAllowExt  bool
	CORSMethods   []string
	CORSHeaders   []string
	CORSOrigins   []string

	NotificationsHost  string
	NotificationsToken string

	UsersServiceHost     string
	UsersServiceToken    string
	UsersServiceUseRedis bool
	UsersServiceTTL      int

	StaticServiceHost string
	StaticServiceSign string

	RedisHost     string
	RedisDb       int
	RedisPassword string

	MongoHost     string
	MongoUser     string
	MongoPassword string
	MongoDatabase string

	JWTPassword []byte

	I18nFile string
}

// SetFromViper applies values from the viper config to the Config instance
func (c *Config) SetFromViper(conf *viper.Viper) {
	logLvl, err := log.ParseLevel(conf.GetString("log_level"))
	if err != nil {
		logLvl = bubucore.Opt.LogLevelDft
	}

	c.Port = conf.GetString("bubu_service_port")
	c.LogLevel = logLvl

	// CORS
	{
		c.CORSEnable = conf.GetBool("cors_enable")
		c.CORSAllowAll = conf.GetBool("cors_allow_all")
		c.CORSAllowCred = conf.GetBool("cors_allow_cred")
		c.CORSAllowWS = conf.GetBool("cors_allow_ws")
		c.CORSAllowExt = conf.GetBool("cors_allow_ext")

		values := strings.TrimSpace(conf.GetString("cors_methods"))
		if values == "" {
			c.CORSMethods = []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete}
		} else {
			c.CORSMethods = strings.Split(values, ",")
		}

		values = strings.TrimSpace(conf.GetString("cors_headers"))
		c.CORSHeaders = strings.Split(values, ",")

		values = strings.TrimSpace(conf.GetString("cors_origins"))
		c.CORSOrigins = strings.Split(values, ",")
	}

	c.NotificationsHost = conf.GetString("bubu_notifications_host")
	c.NotificationsToken = conf.GetString("bubu_notifications_token")

	c.UsersServiceHost = conf.GetString("bubu_users_host")
	c.UsersServiceToken = conf.GetString("bubu_users_token")
	c.UsersServiceUseRedis = conf.GetBool("bubu_users_use_redis")
	c.UsersServiceTTL = conf.GetInt("bubu_users_ttl")

	c.StaticServiceHost = conf.GetString("bubu_staticservice_host")
	c.StaticServiceSign = conf.GetString("bubu_staticservice_sign")

	c.RedisHost = conf.GetString("redis_host")
	c.RedisDb = conf.GetInt("redis_db")
	c.RedisPassword = conf.GetString("redis_password")

	c.MongoHost = conf.GetString("mongo_host")
	c.MongoUser = conf.GetString("mongo_username")
	c.MongoPassword = conf.GetString("mongo_password")
	c.MongoDatabase = conf.GetString("mongo_db")

	c.JWTPassword = []byte(conf.GetString("bubu_jwt_password"))

	c.I18nFile = conf.GetString("i18n_file")
	if c.I18nFile == "" {
		if _, err := os.Stat(i18nFileDft); err == nil {
			c.I18nFile = i18nFileDft
		}
	}

	c.ApplyToGlobals()
}

// ApplyToGlobals applies values from the Config instance to global instances
func (c *Config) ApplyToGlobals() {
	log.SetLevel(c.LogLevel)
	bubucore.Opt.JWTPassword = c.JWTPassword
}
