package auth_route

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"ppl-m1-be/pkg/httpserver"
)

const msgValidation = "Data yang dikirim tidak valid"

func writeBindError(c *gin.Context, err error) {
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		httpserver.FailFields(c, msgValidation, toFieldErrors(verr))
		return
	}
	httpserver.Fail(c, http.StatusBadRequest, httpserver.CodeBadRequest, "Format permintaan tidak valid")
}

func toFieldErrors(verr validator.ValidationErrors) map[string]string {
	fields := make(map[string]string, len(verr))
	for _, fe := range verr {
		fields[strings.ToLower(fe.Field())] = fieldMessage(fe)
	}
	return fields
}

func fieldMessage(fe validator.FieldError) string {
	if fe.Tag() == "email" {
		return "Format email tidak valid"
	}
	return "Isian ini wajib diisi"
}