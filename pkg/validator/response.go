package validators

import (
	"fmt"
	"github.com/go-playground/validator"
)

type ValidationResponse struct {
	Status string            `json:"status"`
	Errors map[string]string `json:"errors"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func Error(msg string) ValidationResponse {
	return ValidationResponse{
		Status: StatusError,
		Errors: map[string]string{"error": msg},
	}
}

func OK() ValidationResponse {
	return ValidationResponse{
		Status: StatusOK,
	}
}

func ValidationError(errs validator.ValidationErrors) ValidationResponse {
	errorsMap := make(map[string]string)

	for _, err := range errs {
		fieldName := err.Field()
		tag := err.Tag()
		param := err.Param()

		var message string

		switch tag {
		case "required":
			message = "Это поле обязательно"
		case "alphanum":
			message = "Допустимы только латинские буквы и цифры"
		case "min":
			message = fmt.Sprintf("Минимум %s символов", param)
		case "oneof":
			message = fmt.Sprintf("Ввидите валидное значение: %s", param)
		}

		errorsMap[fieldName] = message
	}

	return ValidationResponse{
		Status: StatusError,
		Errors: errorsMap,
	}
}
