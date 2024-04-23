package auth

import (
	"log/slog"
	"strings"

	"github.com/elusiv0/medods_test/internal/model/api"
	tokenManager "github.com/elusiv0/medods_test/internal/util/token"
	"github.com/gin-gonic/gin"
)

func Auth(
	tokenManager *tokenManager.TokenManager,
	logger *slog.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenH := c.GetHeader("Authorization")
		if tokenH == "" {
			logger.Error("AuthMiddleware: no authorization token found in request headers")
			c.Error(api.ErrNoAccessTokenFound)
			c.Abort()
			return
		}

		token := strings.Split(tokenH, " ")
		if len(token) != 2 || token[0] != "Bearer" || len(token[1]) == 0 {
			logger.Error("AuthMiddleware: invalid token")
			c.Error(api.ErrInvalidAccessToken)
			c.Abort()
			return
		}

		tokenInfo, err := tokenManager.ValidateJWT(token[1])
		if err != nil {
			logger.Error("AuthMiddleware: ", err.Error())
			c.Error(err)
			return
		}
		c.Set("tokenInfo", tokenInfo)

		c.Next()
	}
}
