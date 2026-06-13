package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func ValidateRegisterFields(nameFull, phone, idNumber, dateOfBirth, email, password string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("validation error: %v", r)
		}
	}()

	if strings.TrimSpace(nameFull) == "" {
		return errors.New("name is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return errors.New("invalid email address")
	}

	return nil
}