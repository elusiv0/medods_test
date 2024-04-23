package router

import (
	"log/slog"
	"net/http"

	authMiddleware "github.com/elusiv0/medods_test/internal/middleware/auth"
	errorsMiddleware "github.com/elusiv0/medods_test/internal/middleware/errors"
	authRouter "github.com/elusiv0/medods_test/internal/router/http/v1/auth"
	authService "github.com/elusiv0/medods_test/internal/service/auth"
	tokenManager "github.com/elusiv0/medods_test/internal/util/token"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func InitRoutes(
	log *slog.Logger,
	tokenM *tokenManager.TokenManager,
	authS *authService.AuthService,
) *gin.Engine {
	router := gin.New()

	router.Use(sloggin.New(log))
	router.Use(errorsMiddleware.ErrorsMiddleware(errorsMiddleware.InitErrors()))

	router.GET("ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	auth := router.Group("api/auth")
	{
		authRouter.New(
			authS,
			log,
			auth,
		)
	}
	v1 := router.Group("api/v1", authMiddleware.Auth(tokenM, log))
	{
		v1.GET("/test", func(c *gin.Context) {
			log.Info("Inside protected end-point")
		})
	}

	return router
}
