package middleware

import (
	"encoding/json"
	"net/http"

	"book-coffee-shop/internal/utils"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidatePayload(target any) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewDecoder(r.Body).Decode(target); err != nil {
				utils.WriteError(w, "invalid request body", http.StatusBadRequest)
				return
			}

			if err := validate.Struct(target); err != nil {
				if errs, ok := err.(validator.ValidationErrors); ok {
					details := make(map[string]string)
					for _, e := range errs {
						details[e.Field()] = validationMessage(e)
					}
					utils.WriteValidationError(w, details, http.StatusBadRequest)
					return
				}
				utils.WriteError(w, "validation failed", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Minimum length is " + e.Param()
	case "max":
		return "Maximum length is " + e.Param()
	case "datetime":
		return "Invalid date format (use YYYY-MM-DD)"
	case "omitempty":
		return ""
	default:
		return "Invalid value"
	}
}
