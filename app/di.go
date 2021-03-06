package app

import (
	"context"
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/di"
	"github.com/bubulearn/bubucore/ginsrv"
	"github.com/bubulearn/bubucore/i18n"
	"github.com/bubulearn/bubucore/mongodb"
	"github.com/bubulearn/bubucore/notifications"
	"github.com/bubulearn/bubucore/staticservice"
	"github.com/bubulearn/bubucore/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

// Dependencies names
const (
	// DIConfigViper contains app config read to viper.Viper instance
	DIConfigViper = "bubu_config_viper"

	// DIConfig contains initialized Config instance
	DIConfig = "bubu_config"

	// DII18n contains initialized i18n.TextsSource instance
	// In case of renaming, see ginsrv/context.go:42
	DII18n = "bubu_i18n"

	// DIRouter contains gin router (gin.Engine) instance
	DIRouter = "bubu_router"

	// DINotifications contains notifications.Client instance
	DINotifications = "bubu_notifications"

	// DIUsersService contains users.Client instance
	DIUsersService = "bubu_users_service"

	// DIStaticService contains users.Client instance
	DIStaticService = "bubu_static_service"

	// DIMongo contains mongodb.MongoDB connection instance, or nil if no mongo host provided in config
	DIMongo = "bubu_mongo"

	// DIRedis contains redis.Client instance, or nil if no redis host provided in config
	DIRedis = "bubu_redis"
)

// GetDefaultDIBuilder returns default DI builder
func GetDefaultDIBuilder() (*di.Builder, error) {
	builder := &di.Builder{}

	err := builder.Add(DIDefConfigViper(), DIDefConfig(), DIDefI18n(), DIDefRouter())
	if err != nil {
		return nil, err
	}

	err = builder.Add(DIDefMongo(), DIDefRedis())
	if err != nil {
		return nil, err
	}

	err = builder.Add(DIDefNotifications(), DIDefUsersService(), DIDefStaticService())
	if err != nil {
		return nil, err
	}

	return builder, nil
}

// BuildDefaultContainer builds Container with default di.Builder
func BuildDefaultContainer() *di.Container {
	builder, err := GetDefaultDIBuilder()
	if err != nil {
		log.Fatal(logTag, "failed to create default DI container: ", err)
	}

	ctn, err := builder.Build()
	if err != nil {
		log.Fatal(logTag, "failed to build default DI container: ", err)
	}

	return ctn
}

// DIDefConfigViper returns app config read from config file to the viper.Viper instance
func DIDefConfigViper() di.Def {
	return di.Def{
		Name: DIConfigViper,
		Build: func(ctn *di.Container) (interface{}, error) {
			return bubucore.ReadConfig()
		},
	}
}

// DIDefConfig returns default Config dependency definition
func DIDefConfig() di.Def {
	return di.Def{
		Name: DIConfig,
		Build: func(ctn *di.Container) (interface{}, error) {
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
		Build: func(ctn *di.Container) (interface{}, error) {
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
		Build: func(ctn *di.Container) (interface{}, error) {
			conf := ctn.Get(DIConfig).(*Config)
			router := ginsrv.GetDefaultRouter()

			// CORS
			if conf.CORSEnable {
				cc := cors.Config{
					AllowWildcard: true,

					AllowAllOrigins:        conf.CORSAllowAll,
					AllowCredentials:       conf.CORSAllowCred,
					AllowWebSockets:        conf.CORSAllowWS,
					AllowBrowserExtensions: conf.CORSAllowExt,

					AllowHeaders: conf.CORSHeaders,
					AllowMethods: conf.CORSMethods,
				}

				if !cc.AllowAllOrigins {
					cc.AllowOrigins = conf.CORSOrigins
				}

				router.Use(cors.New(cc))
			}

			return router, nil
		},
	}
}

// DIDefNotifications returns default notifications.Client dependency definition
func DIDefNotifications() di.Def {
	return di.Def{
		Name: DINotifications,
		Build: func(ctn *di.Container) (interface{}, error) {
			conf := ctn.Get(DIConfig).(*Config)
			client := notifications.NewClient(conf.NotificationsHost, conf.NotificationsToken)
			if conf.NotificationsHost != "" {
				err := client.Ping()
				if err != nil {
					return nil, err
				}
			}
			return client, nil
		},
		Close: func(obj interface{}) error {
			return obj.(*notifications.Client).Close()
		},
	}
}

// DIDefUsersService returns default users.Client dependency definition
func DIDefUsersService() di.Def {
	return di.Def{
		Name: DIUsersService,
		Build: func(ctn *di.Container) (interface{}, error) {
			conf := ctn.Get(DIConfig).(*Config)
			client := users.NewClient(conf.UsersServiceHost, conf.UsersServiceToken)

			if conf.UsersServiceUseRedis {
				redisClient := DIGetRedis(ctn)
				client.SetRedis(redisClient, conf.UsersServiceTTL)
			}

			return client, nil
		},
		Close: func(obj interface{}) error {
			return obj.(*users.Client).Close()
		},
	}
}

// DIDefStaticService returns default staticservice.Client dependency definition
func DIDefStaticService() di.Def {
	return di.Def{
		Name: DIStaticService,
		Build: func(ctn *di.Container) (interface{}, error) {
			conf := ctn.Get(DIConfig).(*Config)
			client := staticservice.NewClient(conf.StaticServiceHost, conf.StaticServiceSign)
			return client, nil
		},
		Close: func(obj interface{}) error {
			return obj.(*staticservice.Client).Close()
		},
	}
}

// DIDefMongo returns default mongodb.MongoDB dependency definition.
// Returns nil if no Config.MongoHost defined in config.
func DIDefMongo() di.Def {
	return di.Def{
		Name: DIMongo,
		Build: func(ctn *di.Container) (interface{}, error) {
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
		Build: func(ctn *di.Container) (interface{}, error) {
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

// DIGetConfigViper returns config viper.Viper from the DI container
func DIGetConfigViper(ctn *di.Container) *viper.Viper {
	return ctn.Get(DIConfigViper).(*viper.Viper)
}

// DIGetConfig returns Config from the DI container
func DIGetConfig(ctn *di.Container) *Config {
	return ctn.Get(DIConfig).(*Config)
}

// DIGetI18n returns i18n.TextsSource from the DI container
func DIGetI18n(ctn *di.Container) *i18n.TextsSource {
	return ctn.Get(DII18n).(*i18n.TextsSource)
}

// DIGetRouter returns gin.Engine router from the DI container
func DIGetRouter(ctn *di.Container) *gin.Engine {
	return ctn.Get(DIRouter).(*gin.Engine)
}

// DIGetNotifications returns notifications.Client from the DI container
func DIGetNotifications(ctn *di.Container) *notifications.Client {
	return ctn.Get(DINotifications).(*notifications.Client)
}

// DIGetUsersService returns users.Client from the DI container
func DIGetUsersService(ctn *di.Container) *users.Client {
	return ctn.Get(DIUsersService).(*users.Client)
}

// DIGetStaticService returns staticservice.Client from the DI container
func DIGetStaticService(ctn *di.Container) *staticservice.Client {
	return ctn.Get(DIStaticService).(*staticservice.Client)
}

// DIGetMongoDB returns mongodb.MongoDB from the DI container
func DIGetMongoDB(ctn *di.Container) *mongodb.MongoDB {
	m := ctn.Get(DIMongo).(*mongodb.MongoDB)
	if m == nil {
		log.Fatal(logTag, "attempt to access nil MongoDB instance")
	}
	return m
}

// DIGetRedis returns redis.Client from the DI container
func DIGetRedis(ctn *di.Container) *redis.Client {
	r := ctn.Get(DIRedis).(*redis.Client)
	if r == nil {
		log.Fatal(logTag, "attempt to access nil redis client instance")
	}
	return r
}
