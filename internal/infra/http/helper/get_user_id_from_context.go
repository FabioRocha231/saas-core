package helper

import (
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "user_id"

func GetUserIDFromContext(c *gin.Context) (string, error) {
	v, ok := c.Get(CtxUserIDKey)
	if !ok {
		return "", errx.New(errx.CodeUnauthorized, "missing user")
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", errx.New(errx.CodeUnauthorized, "invalid user")
	}
	return s, nil
}
