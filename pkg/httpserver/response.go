package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeBadRequest            = "BAD_REQUEST"
	CodeInvalidCredentials    = "INVALID_CREDENTIALS"
	CodeUnauthenticated       = "UNAUTHENTICATED"
	CodeForbidden             = "FORBIDDEN"
	CodeValidation            = "VALIDATION_ERROR"
	CodeInternal              = "INTERNAL_ERROR"
	CodeEmailDuplicate        = "EMAIL_DUPLICATE"
	CodeNotFound              = "NOT_FOUND"
	CodeMethodNotAllowed      = "METHOD_NOT_ALLOWED"
	CodeInvalidParameter      = "INVALID_PARAMETER"
	CodeDivisionCodeDuplicate = "DIVISION_CODE_DUPLICATE"
	CodeRoleCodeDuplicate     = "ROLE_CODE_DUPLICATE"
	CodeUserNotFound          = "USER_NOT_FOUND"
	CodeCannotDeactivateSelf  = "CANNOT_DEACTIVATE_SELF"
	CodeUserAlreadyActive     = "USER_ALREADY_ACTIVE"
	CodeUserAlreadyInactive   = "USER_ALREADY_INACTIVE"
	CodeUserReferenceInactive = "USER_REFERENCE_INACTIVE"
)

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"success": true, "data": data})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   gin.H{"code": code, "message": message},
	})
}

func FailFields(c *gin.Context, message string, fields map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"success": false,
		"error": gin.H{
			"code":    CodeValidation,
			"message": message,
			"fields":  fields,
		},
	})
}
