package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// bindError turns binding/validation errors into a readable 400 response.
func bindError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make(map[string]string, len(ve))
		for _, fe := range ve {
			details[lowerFirst(fe.Field())] = describe(fe)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": details})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
}

func describe(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("must be at most %s characters", fe.Param())
		}
		return "must be at most " + fe.Param()
	case "gte":
		return "must be >= " + fe.Param()
	case "gt":
		return "must be > " + fe.Param()
	case "lte":
		return "must be <= " + fe.Param()
	case "oneof":
		return "must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	case "alphanum":
		return "must contain only letters and numbers"
	}
	return "is invalid (" + fe.Tag() + ")"
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
