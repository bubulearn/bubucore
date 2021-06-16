package bubucore

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// Opt shares package options
var Opt = &Options{
	ConfigFilePath: "./.env",
	ConfigFileType: "env",

	LogsPath:    "data/logs",
	LogFileGin:  "gin.log",
	LogFileApp:  "app.log",
	LogLevelDft: log.InfoLevel,
}

// Options represents package options
type Options struct {
	ConfigFilePath string
	ConfigFileType string

	LogsPath    string
	LogFileGin  string
	LogFileApp  string
	LogLevelDft log.Level

	APIVersion  string
	APIBasePath string

	ServiceName        string
	ServiceDescription string
	ServiceRepo        string

	Hostname string
}

// GetHostname returns hostname from options or OS
func (o *Options) GetHostname() string {
	if o.Hostname == "" {
		o.Hostname, _ = os.Hostname()
	}
	return o.Hostname
}

// ServiceInfo represents info about the service
type ServiceInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	BasePath    string `json:"base_path"`
	Description string `json:"description"`
	Repo        string `json:"repo"`
}

// GetServiceInfo returns map with service info
func GetServiceInfo() *ServiceInfo {
	return &ServiceInfo{
		Name:        Opt.ServiceName,
		Version:     Opt.APIVersion,
		BasePath:    Opt.APIBasePath,
		Description: Opt.ServiceDescription,
		Repo:        Opt.ServiceRepo,
	}
}

// ReadConfig loads configuration from the config file to the viper instance
func ReadConfig() (*viper.Viper, error) {
	conf := viper.New()

	conf.SetConfigFile(Opt.ConfigFilePath)
	conf.SetConfigType(Opt.ConfigFileType)

	err := conf.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return conf, nil
}
