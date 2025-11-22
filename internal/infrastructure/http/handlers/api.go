package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"pr-service/internal/domain/pr"
	dto "pr-service/internal/infrastructure/http/handlers/dto"
	"pr-service/internal/infrastructure/http/middleware"
	openapi "pr-service/internal/infrastructure/http/openapi"
	"pr-service/internal/infrastructure/http/transport"
	"pr-service/internal/infrastructure/storage/postgres"
	"pr-service/pkg/sl_logger/sl"
	validateResp "pr-service/pkg/validator"

	"github.com/go-playground/validator"
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

	var req dto.PostPullRequestCreateJSONBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Log.Error("bad request",
			slog.String("type", err.Error()),
			sl.Err(err),
		)
		responseErr(w, http.StatusBadRequest, "invalid json body")
		return
	}

	h.Log.Info("request body decoded", slog.Any("req", req))

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		h.Log.Error("invalid request", sl.Err(err))
		_ = transport.WriteJSON(w, http.StatusBadRequest, validateResp.ValidationError(validateErr))
		return
	}

	pr := dto.PostPullRequestMapToModel(req)

	svcPr, err := h.Svc.PullRequestCreate(r.Context(), pr)
	if errors.Is(err, postgres.ErrNoCandidate) {
		h.Log.Error("bad request",
			slog.String("type", err.Error()),
			sl.Err(err),
		)
		responseErr(w, http.StatusInternalServerError, postgres.ErrNoCandidate.Error())
		return
	}
	if errors.Is(err, postgres.ErrPrExists) {
		h.Log.Error("bad request",
			slog.String("type", err.Error()),
			sl.Err(err),
		)
		responseErr(w, http.StatusInternalServerError, postgres.ErrPrExists.Error())
		return
	}
	if err != nil {
		h.Log.Error("bad request",
			slog.String("type", err.Error()),
			sl.Err(err),
		)
		responseErr(w, http.StatusInternalServerError, "failed to add pullRequest")
		return
	}

	h.Log.Info("pullRequest CREATE", slog.Any("title", svcPr.PullRequestName))
	pullRequestCreateOK(w, svcPr)
}

func responseErr(w http.ResponseWriter, c int, m string) {
	resp := openapi.ErrorResponse{
		Error: struct {
			Code    openapi.ErrorResponseErrorCode "json:\"code\""
			Message string                         "json:\"message\""
		}{
			Code:    openapi.ErrorResponseErrorCode(c),
			Message: m,
		},
	}
	transport.WriteJSON(w, c, resp)
}

func pullRequestCreateOK(w http.ResponseWriter, pr pr.PullRequest) {
	r := openapi.PullRequestShort{
		AuthorId:        pr.AuthorId,
		PullRequestId:   pr.PullRequestId,
		PullRequestName: pr.PullRequestName,
		Status:          openapi.PullRequestShortStatus(pr.Status),
	}
	transport.WriteJSON(w, http.StatusOK, r)
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
