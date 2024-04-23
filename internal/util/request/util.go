package request

import (
	"fmt"

	"github.com/elusiv0/medods_test/internal/model/api"
	"github.com/gin-gonic/gin"
)

const (
	uuid = "uuid"
)

func GetIdFromQueryPath(c *gin.Context) (string, error) {
	uuid := c.Query(uuid)

	if uuid == "" {
		return "", fmt.Errorf("request utils - GetIdFromQueryPath: %w", api.ErrNoUUID)
	}

	return uuid, nil
}
