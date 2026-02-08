package middleware

import (
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/infra/http/handlers"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService ports.JwtInterface
}

func NewAuthMiddleware(jwtService ports.JwtInterface) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

const (
	CtxUserIDKey = "user_id"
	CtxRoleKey   = "user_role"
)

func (m *AuthMiddleware) Middleware(c *gin.Context) {
	h := c.GetHeader("Authorization")
	if h == "" {
		handlers.RespondErr(c, errx.New(errx.CodeUnauthorized, "unauthorized"))
		c.Abort()
		return
	}

	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		handlers.RespondErr(c, errx.New(errx.CodeUnauthorized, "invalid auth header"))
		c.Abort()
		return
	}

	claims, err := m.jwtService.Parse(parts[1])
	if err != nil {
		handlers.RespondErr(c, errx.New(errx.CodeUnauthorized, "invalid token"))
		c.Abort()
		return
	}

	c.Set(CtxUserIDKey, claims.UserID)
	c.Set(CtxRoleKey, claims.Role)
	c.Next()
}
