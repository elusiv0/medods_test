package auth

import (
	"log/slog"
	"net/http"

	"github.com/elusiv0/medods_test/internal/model/api"
	tokenDto "github.com/elusiv0/medods_test/internal/model/token"
	authService "github.com/elusiv0/medods_test/internal/service/auth"
	reqUtils "github.com/elusiv0/medods_test/internal/util/request"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authService *authService.AuthService
	logger      *slog.Logger
}

func New(
	authService *authService.AuthService,
	log *slog.Logger,
	group *gin.RouterGroup,
) {
	authRouter := &AuthRouter{
		logger:      log,
		authService: authService,
	}

	group.POST("/sign-in", authRouter.signIn)
	group.POST("/refresh", authRouter.refresh)
}

func (authRouter *AuthRouter) signIn(c *gin.Context) {
	uuid, err := reqUtils.GetIdFromQueryPath(c)

	if err != nil {
		authRouter.logger.Error("AuthRouter - signIn: ", err.Error())
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	tokenResponse, err := authRouter.authService.SignIn(ctx, uuid)
	if err != nil {
		authRouter.logger.Error("AuthRouter - signIn: ", err.Error())
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func (authRouter *AuthRouter) refresh(c *gin.Context) {
	refreshReponse := tokenDto.RefreshRequest{}

	err := c.ShouldBindJSON(&refreshReponse)
	if err != nil {
		authRouter.logger.Error("AuthRouter - refresh - ", err.Error())
		c.Error(api.ErrBadRefreshRequest)
		return
	}

	ctx := c.Request.Context()
	accessToken := refreshReponse.AccessToken
	refreshToken := refreshReponse.RefreshToken
	tokenResponse, err := authRouter.authService.Refresh(ctx, refreshToken, accessToken)
	if err != nil {
		authRouter.logger.Error("AuthRouter - refresh - ", err.Error())
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}
