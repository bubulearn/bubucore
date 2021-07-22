package app

import (
	"github.com/bubulearn/bubucore/i18n"
	"github.com/bubulearn/bubucore/mongodb"
	"github.com/bubulearn/bubucore/notifications"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewApp creates new App instance.
// If ctn is nil, DIBuilderDft() will be called.
func NewApp(ctn *di.Container) *App {
	var c di.Container
	if ctn == nil {
		builder, err := DIBuilderDft()
		if err != nil {
			log.Fatal("failed to init default DI builder: ", err)
		}
		c = builder.Build()
	} else {
		c = *ctn
	}
	return &App{
		ctn: c,
	}
}

// InitRouterFn is a function to prepare router before run
type InitRouterFn = func(router *gin.Engine) error

// PrepareContainerFn is a function to prepare DI container before run
type PrepareContainerFn = func(ctn di.Container) error

// App is a Bubulearn service app
type App struct {
	ctn di.Container

	initRouterFn InitRouterFn
	prepareCtnFn PrepareContainerFn
}

// SetInitRouterFn sets init router hook
func (a *App) SetInitRouterFn(fn InitRouterFn) {
	a.initRouterFn = fn
}

// SetPrepareContainerFn sets prepare DI container hook
func (a *App) SetPrepareContainerFn(fn PrepareContainerFn) {
	a.prepareCtnFn = fn
}

// GetContainer returns an App's di.Container instance
func (a *App) GetContainer() di.Container {
	return a.ctn
}

// GetConfigViper returns config viper.Viper from the DI container
func (a *App) GetConfigViper() *viper.Viper {
	return a.ctn.Get(DIConfigViper).(*viper.Viper)
}

// GetConfig returns Config from the DI container
func (a *App) GetConfig() *Config {
	return a.ctn.Get(DIConfig).(*Config)
}

// GetI18n returns i18n.TextsSource from the DI container
func (a *App) GetI18n() *i18n.TextsSource {
	return a.ctn.Get(DII18n).(*i18n.TextsSource)
}

// GetRouter returns gin.Engine router from the DI container
func (a *App) GetRouter() *gin.Engine {
	return a.ctn.Get(DIRouter).(*gin.Engine)
}

// GetNotifications returns notifications.Client from the DI container
func (a *App) GetNotifications() *notifications.Client {
	return a.ctn.Get(DINotifications).(*notifications.Client)
}

// GetMongoDB returns mongodb.MongoDB from the DI container
func (a *App) GetMongoDB() *mongodb.MongoDB {
	m := a.ctn.Get(DIMongo).(*mongodb.MongoDB)
	if m == nil {
		log.Fatal("[bubucore.App] attempt to access nil MongoDB instance")
	}
	return m
}

// GetRedis returns redis.Client from the DI container
func (a *App) GetRedis() *redis.Client {
	r := a.ctn.Get(DIRedis).(*redis.Client)
	if r == nil {
		log.Fatal("[bubucore.App] attempt to access nil redis client instance")
	}
	return r
}

// Init initializes App without starting the server
func (a *App) Init() {
	if a.prepareCtnFn != nil {
		err := a.prepareCtnFn(a.ctn)
		if err != nil {
			log.Fatal("[bubucore.App] failed to prepare DI container: ", err)
		}
	}

	router := a.GetRouter()

	if a.initRouterFn != nil {
		err := a.initRouterFn(router)
		if err != nil {
			log.Fatal("[bubucore.App] failed to init router: ", err)
		}
	}
}

// Run starts the App's server
func (a *App) Run() {
	a.Init()
	defer a.Close()

	router := a.GetRouter()
	conf := a.GetConfig()

	log.Info("[bubucore.App] starting gin server")

	err := router.Run(":" + conf.Port)
	if err != nil {
		log.Fatal(err)
	}
}

// Close finalizes the App
func (a *App) Close() {
	err := a.ctn.Delete()
	if err != nil {
		log.Fatal(err)
	}
}
