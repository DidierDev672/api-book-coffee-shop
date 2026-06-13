package models

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type RegisterRequest struct {
	NameFull    string `json:"name_full" validate:"required,min=3,max=100"`
	Phone       string `json:"phone" validate:"required,min=8,max=20"`
	IDNumber    string `json:"id_number" validate:"required,min=5,max=20"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Email       string `json:"email" validate:"required,email,max=100"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
	Token       string `json:"token" validate:"omitempty"`
}

func (r *RegisterRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if err := validate.Struct(r); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Field()] = getValidationMessage(e)
		}
	}
	return errs
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Minimum length is " + err.Param()
	case "max":
		return "Maximum length is " + err.Param()
	case "datetime":
		return "Invalid date format (use YYYY-MM-DD)"
	default:
		return "Invalid value"
	}
}
