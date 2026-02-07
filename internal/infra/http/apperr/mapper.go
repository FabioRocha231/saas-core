package apperr

import (
	"net/http"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/infra/http/response"
)

func ToHTTP(err error) (int, response.Envelope) {
	if err == nil {
		return http.StatusOK, response.Ok(nil)
	}

	code := errx.CodeOf(err)
	msg := errx.MsgOf(err)

	switch code {
	case errx.CodeInvalid:
		return http.StatusBadRequest, response.Fail(string(code), msg)
	case errx.CodeNotFound:
		return http.StatusNotFound, response.Fail(string(code), msg)
	case errx.CodeConflict:
		return http.StatusConflict, response.Fail(string(code), msg)
	case errx.CodeUnauthorized:
		return http.StatusUnauthorized, response.Fail(string(code), msg)
	case errx.CodeForbidden:
		return http.StatusForbidden, response.Fail(string(code), msg)
	default:
		return http.StatusInternalServerError, response.Fail(string(errx.CodeInternal), "internal error")
	}
}
