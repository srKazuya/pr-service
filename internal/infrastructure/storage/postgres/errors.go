package postgres

import (
	"errors"
	"pr-service/internal/infrastructure/http/openapi"
)

var (
	ErrOpenDB          = errors.New("failed to open database")
	ErrMigration       = errors.New("failed to run migrations")
	ErrGormOpen        = errors.New("failed to gorm open")
)

type codedError struct {
	code    openapi.ErrorResponseErrorCode
	message string
}

func (e codedError) Error() string {
	return e.message
}

func (e codedError) Code() openapi.ErrorResponseErrorCode {
	return e.code
}

var (
	ErrNoCandidate = codedError{
		code:    openapi.NOCANDIDATE,
		message: "нет доступных ревьюеров в команде",
	}
	ErrNotFound = codedError{
		code:    openapi.NOTFOUND,
		message: "не найдено",
	}
	ErrPrExists = codedError{
		code:    openapi.PREXISTS,
		message: "pull request уже существует",
	}
	ErrNotAssigned = codedError{
		code:    openapi.NOTASSIGNED,
		message: "у пользователя нет команды",
	}
)
