package di

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/elusiv0/medods_test/internal/app"
	"github.com/elusiv0/medods_test/internal/config"
	tokenRepository "github.com/elusiv0/medods_test/internal/repo/token"
	userRepository "github.com/elusiv0/medods_test/internal/repo/user"
	httpRouter "github.com/elusiv0/medods_test/internal/router/http"
	authService "github.com/elusiv0/medods_test/internal/service/auth"
	tokenManager "github.com/elusiv0/medods_test/internal/util/token"
	"github.com/elusiv0/medods_test/pkg/httpserver"
	"github.com/elusiv0/medods_test/pkg/logger"
	mongo "github.com/elusiv0/medods_test/pkg/mongo"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
)

const (
	Config          = "config"
	Logger          = "logger"
	App             = "app"
	Router          = "router"
	Httpserver      = "httpserver"
	TokenManager    = "tokenManager"
	Mongo           = "mongo"
	TokenRepository = "tokenRepository"
	UserRepository  = "userRepository"
	AuthService     = "authService"
)

func InitContainer() (di.Container, error) {
	builder, err := initBuilder()
	if err != nil {
		return nil, err
	}

	container := builder.Build()

	return container, nil
}

func initBuilder() (*di.Builder, error) {
	builder, err := di.NewBuilder()

	if err != nil {
		return nil, fmt.Errorf("error with building di container: %w", err)
	}

	RegisterDeps(builder)

	return builder, nil
}

func RegisterDeps(b *di.Builder) {
	//building config
	b.Add(di.Def{
		Name: Config,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.GetConfig()
		},
	})

	//building logger
	b.Add(di.Def{
		Name: Logger,
		Build: func(ctn di.Container) (interface{}, error) {
			return logger.New(
				ctn.Get("config").(*config.Config).App.Environment,
			), nil
		},
	})

	//building token manager
	b.Add(di.Def{
		Name: TokenManager,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)

			return tokenManager.New(
				cfg.Jwt.LifeTime,
				cfg.Jwt.Secret,
			), nil
		},
	})

	//building mongo
	b.Add(di.Def{
		Name: Mongo,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*config.Config)
			fmt.Println(cfg.Mongo.Host)
			logger := ctn.Get("logger").(*slog.Logger)
			mongoConn := mongo.NewMongoConn(
				cfg.Mongo.Host,
				cfg.Mongo.Port,
				cfg.Mongo.DbName,
				mongo.WithConnectionAttempts(cfg.Mongo.ConnectionAttempts),
				mongo.WithTimeout(cfg.Mongo.ConnectionTimeout),
				mongo.WithCredentials(cfg.Mongo.User, cfg.Mongo.Password),
			)

			return mongo.New(
				context.Background(),
				mongoConn,
				logger,
			)
		},
	})

	//building repositories
	b.Add(di.Def{
		Name: TokenRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			mongoClient := ctn.Get("mongo").(*mongo.MongoClient)
			logger := ctn.Get("logger").(*slog.Logger)

			return tokenRepository.New(
				mongoClient,
				logger,
			), nil
		},
	})
	b.Add(di.Def{
		Name: UserRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			mongoClient := ctn.Get("mongo").(*mongo.MongoClient)
			logger := ctn.Get("logger").(*slog.Logger)

			return userRepository.New(
				mongoClient,
				logger,
			), nil
		},
	})

	//building services
	b.Add(di.Def{
		Name: AuthService,
		Build: func(ctn di.Container) (interface{}, error) {
			userRepo := ctn.Get("userRepository").(*userRepository.UserRepo)
			tokenRepo := ctn.Get("tokenRepository").(*tokenRepository.TokenRepo)
			logger := ctn.Get("logger").(*slog.Logger)
			tokenManager := ctn.Get("tokenManager").(*tokenManager.TokenManager)

			return authService.New(
				userRepo,
				tokenRepo,
				logger,
				tokenManager,
			), nil
		},
	})

	//building router
	b.Add(di.Def{
		Name: Router,
		Build: func(ctn di.Container) (interface{}, error) {
			logger := ctn.Get("logger").(*slog.Logger)
			tokenManager := ctn.Get("tokenManager").(*tokenManager.TokenManager)
			authService := ctn.Get("authService").(*authService.AuthService)

			return httpRouter.InitRoutes(
				logger,
				tokenManager,
				authService,
			), nil
		},
	})

	//building http server
	b.Add(di.Def{
		Name: Httpserver,
		Build: func(ctn di.Container) (interface{}, error) {
			router := ctn.Get("router").(*gin.Engine)
			cfg := ctn.Get("config").(*config.Config)

			server := httpserver.New(
				router,
				httpserver.Port(cfg.Http.Port),
				httpserver.ReadTimeout(cfg.Http.ReadTimeout),
				httpserver.WriteTimeout(cfg.Http.WriteTimeout),
				httpserver.ShutdownTimeout(cfg.Http.ShutdownTimeout),
			)

			return server, nil
		},
	})

	//building app
	b.Add(di.Def{
		Name: App,
		Build: func(ctn di.Container) (interface{}, error) {
			server := ctn.Get("httpserver").(*httpserver.HttpServer)
			logger := ctn.Get("logger").(*slog.Logger)

			return app.New(
				server,
				logger,
			), nil
		},
	})
}
