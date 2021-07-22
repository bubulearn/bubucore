package app

import (
	"context"
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/ginsrv"
	"github.com/bubulearn/bubucore/i18n"
	"github.com/bubulearn/bubucore/mongodb"
	"github.com/bubulearn/bubucore/notifications"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di"
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
