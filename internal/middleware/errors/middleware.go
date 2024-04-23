package errors

import (
	"errors"
	"net/http"

	api "github.com/elusiv0/medods_test/internal/model/api"
	token "github.com/elusiv0/medods_test/internal/model/token"
	user "github.com/elusiv0/medods_test/internal/model/user"

	"github.com/gin-gonic/gin"
)

type response struct {
	Error string `json:"error_message"`
}

func InitErrors() map[error]int {
	errs := make(map[error]int)

	errs[api.ErrNoUUID] = http.StatusBadRequest
	errs[api.ErrNoAccessTokenFound] = http.StatusUnauthorized
	errs[api.ErrInvalidAccessToken] = http.StatusUnauthorized
	errs[api.ErrAccessTokenExpired] = http.StatusUnauthorized
	errs[api.ErrBadRefreshRequest] = http.StatusUnauthorized
	errs[api.ErrTokenMismatch] = http.StatusUnauthorized

	errs[token.ErrRefreshTokenNotRegistered] = http.StatusUnauthorized

	errs[user.ErrUserNotFound] = http.StatusUnauthorized

	return errs
}

func ErrorsMiddleware(errs map[error]int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		firstError := err

		for err != nil {
			firstError = err
			err = errors.Unwrap(err)
		}

		if code, ok := errs[firstError]; ok {
			c.JSON(code, &response{
				Error: firstError.Error(),
			})
		} else {
			c.Status(http.StatusInternalServerError)
		}

		c.Errors = c.Errors[:0]
	}
}
