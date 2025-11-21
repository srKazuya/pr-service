package handlers

import (
	"log/slog"
	"net/http"
	"pr-service/internal/domain/pr"
	"pr-service/internal/infrastructure/http/middleware"
	openapi "pr-service/internal/infrastructure/http/openapi"
)

type API struct {
	Log *slog.Logger
	Svc pr.Service
}

// Создать PR и автоматически назначить до 2 ревьюверов из команды автора
// (POST /pullRequest/create)
func (h *API) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	const op = "hanlers.PostPullRequestCreate"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)
}

// Пометить PR как MERGED (идемпотентная операция)
// (POST /pullRequest/merge)
func (h *API) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
}

// Переназначить конкретного ревьювера на другого из его команды
// (POST /pullRequest/reassign)
func (h *API) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
}

// Создать команду с участниками (создаёт/обновляет пользователей)
// (POST /team/add)
func (h *API) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
}

// Получить команду с участниками
// (GET /team/get
func (h *API) GetTeamGet(w http.ResponseWriter, r *http.Request, params openapi.GetTeamGetParams) {
}

// Получить PR'ы, где пользователь назначен ревьювером
// (GET /users/getReview)
func (h *API) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params openapi.GetUsersGetReviewParams) {
}

// Установить флаг активности пользователя
// (POST /users/setIsActive)
func (h *API) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
}

// Установить флаг активности пользователя
// (POST /users/setIsActive)
