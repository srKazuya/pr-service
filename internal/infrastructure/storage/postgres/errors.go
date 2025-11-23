package postgres

import (
	"errors"
	"pr-service/internal/infrastructure/http/openapi"
)

var (
	ErrOpenDB          = errors.New("failed to open database")
	ErrMigration       = errors.New("failed to run migrations")
	ErrGormOpen        = errors.New("failed to gorm open")
	ErrReviewerNotInPR = errors.New("reviewer not in pr")
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
	ErrAlreadyMerged = codedError{
		code:    openapi.NOTASSIGNED,
		message: "pull request уже смержен",
	}
	ErrNotAssigned = codedError{
		code:    openapi.PRMERGED,
		message: "у пользователя нет команды",
	}
	ErrTeamExists = codedError{
		code:    openapi.TEAMEXISTS,
		message: "Команда существует",
	}
)
