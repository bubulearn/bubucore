package app

import (
	"context"
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/ginsrv"
	"github.com/bubulearn/bubucore/i18n"
	"github.com/bubulearn/bubucore/mongodb"
	"github.com/bubulearn/bubucore/notifications"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

// Dependencies names
const (
	// DIConfigViper contains app config read to viper.Viper instance
	DIConfigViper = "config_viper"

	// DIConfig contains initialized Config instance
	DIConfig = "config"

	// DII18n contains initialized i18n.TextsSource instance
	DII18n = "i18n"

	// DIRouter contains gin router (gin.Engine) instance
	DIRouter = "router"

	// DINotifications contains notifications.Client instance
	DINotifications = "notifications"

	// DIMongo contains mongodb.MongoDB connection instance, or nil if no mongo host provided in config
	DIMongo = "mongo"

	// DIRedis contains redis.Client instance, or nil if no redis host provided in config
	DIRedis = "redis"
)

// BuildContainer builds Container with di.Builder
func BuildContainer(builder *di.Builder) *Container {
	c := builder.Build()
	return &Container{
		Container: c,
	}
}

// BuildDefaultContainer builds Container with default di.Builder
func BuildDefaultContainer() *Container {
	builder, err := DIBuilderDft()
	if err != nil {
		log.Fatal(logTag, "failed to build default DI container: ", err)
	}
	return BuildContainer(builder)
}

// Container is a DI container
type Container struct {
	di.Container
}

// GetConfigViper returns config viper.Viper from the DI container
func (c *Container) GetConfigViper() *viper.Viper {
	return c.Get(DIConfigViper).(*viper.Viper)
}

// GetConfig returns Config from the DI container
func (c *Container) GetConfig() *Config {
	return c.Get(DIConfig).(*Config)
}

// GetI18n returns i18n.TextsSource from the DI container
func (c *Container) GetI18n() *i18n.TextsSource {
	return c.Get(DII18n).(*i18n.TextsSource)
}

// GetRouter returns gin.Engine router from the DI container
func (c *Container) GetRouter() *gin.Engine {
	return c.Get(DIRouter).(*gin.Engine)
}

// GetNotifications returns notifications.Client from the DI container
func (c *Container) GetNotifications() *notifications.Client {
	return c.Get(DINotifications).(*notifications.Client)
}

// GetMongoDB returns mongodb.MongoDB from the DI container
func (c *Container) GetMongoDB() *mongodb.MongoDB {
	m := c.Get(DIMongo).(*mongodb.MongoDB)
	if m == nil {
		log.Fatal(logTag, "attempt to access nil MongoDB instance")
	}
	return m
}

// GetRedis returns redis.Client from the DI container
func (c *Container) GetRedis() *redis.Client {
	r := c.Get(DIRedis).(*redis.Client)
	if r == nil {
		log.Fatal(logTag, "attempt to access nil redis client instance")
	}
	return r
}

// DIBuilderDft returns default DI builder
func DIBuilderDft() (*di.Builder, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	err = builder.Add(DIDefConfigViper(), DIDefConfig(), DIDefI18n(), DIDefRouter())
	if err != nil {
		return nil, err
	}

	err = builder.Add(DIDefNotifications())
	if err != nil {
		return nil, err
	}

	err = builder.Add(DIDefMongo(), DIDefRedis())
	if err != nil {
		return nil, err
	}

	return builder, nil
}

// DIDefConfigViper returns app config read from config file to the viper.Viper instance
func DIDefConfigViper() di.Def {
	return di.Def{
		Name: DIConfigViper,
		Build: func(ctn di.Container) (interface{}, error) {
			return bubucore.ReadConfig()
		},
	}
}

// DIDefConfig returns default Config dependency definition
func DIDefConfig() di.Def {
	return di.Def{
		Name: DIConfig,
		Build: func(ctn di.Container) (interface{}, error) {
			vpr := ctn.Get(DIConfigViper).(*viper.Viper)

			conf := &Config{}
			conf.SetFromViper(vpr)

			return conf, nil
		},
	}
}

// DIDefI18n returns default i18n.TextsSource dependency definition
func DIDefI18n() di.Def {
	return di.Def{
		Name: DII18n,
		Build: func(ctn di.Container) (interface{}, error) {
			var err error
			conf := ctn.Get(DIConfig).(*Config)
			if conf.I18nFile != "" {
				i18n.Source, err = i18n.NewSourceFromFile(conf.I18nFile)
			}
			return i18n.Source, err
		},
	}
}

// DIDefRouter returns default gin.Engine dependency definition
func DIDefRouter() di.Def {
	return di.Def{
		Name: DIRouter,
		Build: func(ctn di.Container) (interface{}, error) {
			return ginsrv.GetDefaultRouter(), nil
		},
	}
}

// DIDefNotifications returns default notifications.Client dependency definition
func DIDefNotifications() di.Def {
	return di.Def{
		Name: DINotifications,
		Build: func(ctn di.Container) (interface{}, error) {
			conf := ctn.Get(DIConfig).(*Config)
			return notifications.NewClient(conf.NotificationsHost, conf.NotificationsToken), nil
		},
		Close: func(obj interface{}) error {
			return obj.(*notifications.Client).Close()
		},
	}
}

// DIDefMongo returns default mongodb.MongoDB dependency definition.
// Returns nil if no Config.MongoHost defined in config.
func DIDefMongo() di.Def {
	return di.Def{
		Name: DIMongo,
		Build: func(ctn di.Container) (interface{}, error) {
			conf := ctn.Get(DIConfig).(*Config)
			if conf.MongoHost == "" {
				return nil, nil
			}
			opt := &mongodb.Options{
				Hosts:    []string{conf.MongoHost},
				Database: conf.MongoDatabase,
				User:     conf.MongoUser,
				Password: conf.MongoPassword,
			}
			return mongodb.NewMongoDB(opt)
		},
		Close: func(obj interface{}) error {
			m, ok := obj.(*mongodb.MongoDB)
			if ok && m != nil {
				return obj.(*mongodb.MongoDB).Close()
			}
			return nil
		},
	}
}

// DIDefRedis returns default redis.Client dependency definition.
// Returns nil if no Config.RedisHost defined in config.
func DIDefRedis() di.Def {
	return di.Def{
		Name: DIRedis,
		Build: func(ctn di.Container) (interface{}, error) {
			conf := ctn.Get(DIConfig).(*Config)
			if conf.RedisHost == "" {
				return nil, nil
			}

			client := redis.NewClient(&redis.Options{
				Addr:     conf.RedisHost,
				Password: conf.RedisPassword,
				DB:       conf.RedisDb,
			})

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := client.Ping(ctx).Err()
			if err != nil {
				return nil, err
			}

			return client, nil
		},
		Close: func(obj interface{}) error {
			client, ok := obj.(*redis.Client)
			if ok && client != nil {
				return client.Close()
			}
			return nil
		},
	}
}
