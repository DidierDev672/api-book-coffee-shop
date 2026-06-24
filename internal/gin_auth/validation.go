package gin_auth

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func parseValidationError(err error) map[string]string {
	errs := make(map[string]string)

	if _, ok := err.(validator.ValidationErrors); !ok {
		if strings.Contains(err.Error(), "required") {
			field := extractField(err.Error())
			errs[field] = "this field is required"
		} else {
			errs["_error"] = err.Error()
		}
		return errs
	}

	for _, e := range err.(validator.ValidationErrors) {
		switch e.Tag() {
		case "required":
			errs[e.Field()] = "this field is required"
		case "email":
			errs[e.Field()] = "invalid email format"
		case "min":
			errs[e.Field()] = "minimum length is " + e.Param()
		case "max":
			errs[e.Field()] = "maximum length is " + e.Param()
		default:
			errs[e.Field()] = "invalid value"
		}
	}
	return errs
}

func extractField(msg string) string {
	idx := strings.Index(msg, "'")
	if idx == -1 {
		return "field"
	}
	end := strings.Index(msg[idx+1:], "'")
	if end == -1 {
		return "field"
	}
	return strings.ToLower(msg[idx+1 : idx+1+end])
}
