package handlers

import (
	"github.com/FabioRocha231/saas-core/internal/infra/http/apperr"
	"github.com/FabioRocha231/saas-core/internal/infra/http/response"
	"github.com/gin-gonic/gin"
)

func RespondOK(c *gin.Context, status int, data any) {
	c.JSON(status, response.Ok(data))
}

func RespondErr(c *gin.Context, err error) {
	status, body := apperr.ToHTTP(err)
	c.JSON(status, body)
}
