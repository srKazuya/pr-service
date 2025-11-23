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
	const op = "handlers.PostPullRequestCreate"

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
		responseErr(w, http.StatusBadRequest, transport.ErrInvalidRequest.Error())
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
	pullRequestOK(w, svcPr)
}

// Пометить PR как MERGED (идемпотентная операция)
// (POST /pullRequest/merge)
func (h *API) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostPullRequestMerge"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)

	var req dto.PostPullRequestMergeJSONBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Log.Error("bad request",
			slog.String("type", err.Error()),
			sl.Err(err),
		)
		responseErr(w, http.StatusBadRequest, transport.ErrInvalidRequest.Error())
		return
	}

	if req.PullRequestId == "" {
		h.Log.Warn("pull_request_id is empty")
		responseErr(w, http.StatusBadRequest, "pull_request_id is required")
		return
	}
	h.Log.Info("request body decoded", slog.Any("req", req))

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		h.Log.Error("invalid request", sl.Err(err))
		_ = transport.WriteJSON(w, http.StatusBadRequest, validateResp.ValidationError(validateErr))
		return
	}

	svcPr, err := h.Svc.PullRequestMerge(r.Context(), req.PullRequestId)
	if errors.Is(err, postgres.ErrNotFound) {
		h.Log.Error("bad request",
			slog.String("type", err.Error()),
			sl.Err(err),
		)
		responseErr(w, http.StatusInternalServerError, postgres.ErrNotFound.Error())
		return
	}
	if err != nil {
		h.Log.Error("bad request",
			slog.String("type", err.Error()),
			sl.Err(err),
		)
		responseErr(w, http.StatusInternalServerError, "failed to PullRequestMerge")
		return
	}

	h.Log.Info("pullRequPullRequestMerged", slog.Any("title", svcPr.PullRequestName))
	pullRequestOK(w, svcPr)
}

// Переназначить конкретного ревьювера на другого из его команды
// (POST /pullRequest/reassign)
func (h *API) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostPullRequestMerge"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)

	var req openapi.PostPullRequestReassignJSONBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Log.Error("bad request",
			slog.String("type", err.Error()),
			sl.Err(err),
		)
		responseErr(w, http.StatusBadRequest, transport.ErrInvalidRequest.Error())
		return
	}
	h.Log.Info("request body decoded", slog.Any("req", req))

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		h.Log.Error("invalid request", sl.Err(err))
		_ = transport.WriteJSON(w, http.StatusBadRequest, validateResp.ValidationError(validateErr))
		return
	}

	prReassign := dto.PostPullRequestReassignToModel(req)
	updatedPR, err := h.Svc.PullRequestReassign(r.Context(), prReassign)

	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrNotFound):
			responseErr(w, http.StatusNotFound, postgres.ErrNotFound.Error())
		case errors.Is(err, postgres.ErrReviewerNotInPR):
			responseErr(w, http.StatusBadRequest, postgres.ErrReviewerNotInPR.Error())
		case errors.Is(err, postgres.ErrNotAssigned):
			responseErr(w, http.StatusBadRequest, postgres.ErrNotAssigned.Error())
		case errors.Is(err, postgres.ErrAlreadyMerged):
			responseErr(w, http.StatusBadRequest, postgres.ErrAlreadyMerged.Error())
		default:
			h.Log.Error("reassign failed", sl.Err(err))
			responseErr(w, http.StatusInternalServerError, "failed to reassign reviewer")
		}
		return
	}

	h.Log.Info("reviewer reassigned",
		slog.String("pr_id", updatedPR.PullRequestId),
	)

	pullRequestOK(w, updatedPR)
}

// Создать команду с участниками (создаёт/обновляет пользователей)
// (POST /team/add)
func (h *API) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostTeamAdd"
	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestID(r)),
	)

	var req openapi.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Error("failed to decode request", sl.Err(err))
		responseErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := validator.New().Struct(req); err != nil {
		h.Log.Warn("validation failed", sl.Err(err))
		_ = transport.WriteJSON(w, http.StatusBadRequest, validateResp.ValidationError(err.(validator.ValidationErrors)))
		return
	}

	teamDomain := dto.TeamMapToModel(req)

	createdTeam, err := h.Svc.TeamAdd(r.Context(), teamDomain)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrTeamExists):
			h.Log.Warn("team already exists", slog.String("team", teamDomain.TeamName))
			responseErr(w, http.StatusConflict, "команда уже существует")
			return

		case errors.Is(err, postgres.ErrNoCandidate):
			h.Log.Warn("attempt to create empty team", slog.String("team", teamDomain.TeamName))
			responseErr(w, http.StatusBadRequest, "в команде должен быть хотя бы один участник")
			return

		case errors.Is(err, postgres.ErrNotFound):
			h.Log.Warn("user not found when creating team", slog.String("team", teamDomain.TeamName))
			responseErr(w, http.StatusNotFound, "один или несколько пользователей не найдены")
			return

		default:
			h.Log.Error("failed to create team", sl.Err(err), slog.String("team", teamDomain.TeamName))
			responseErr(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
			return
		}
	}

	h.Log.Info("team created successfully", slog.String("team", createdTeam.TeamName))
	teamRequestOK(w, createdTeam)
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

func responseErr(w http.ResponseWriter, c int, m string) {
	resp := openapi.ErrorResponse{
		Error: struct {
			Code    openapi.ErrorResponseErrorCode "json:\"code\""
			Message string                         "json:\"message\""
		}{
			Code:    openapi.ErrorResponseErrorCode(http.StatusText(c)),
			Message: m,
		},
	}
	transport.WriteJSON(w, c, resp)
}

func pullRequestOK(w http.ResponseWriter, pr pr.PullRequest) {
	r := openapi.PullRequest{
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		AuthorId:          pr.AuthorId,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            openapi.PullRequestStatus(pr.Status),
	}
	transport.WriteJSON(w, http.StatusOK, r)
}

func teamRequestOK(w http.ResponseWriter, pr pr.Team) {
	members := make([]openapi.TeamMember, 0, len(pr.Members))
	for _, m := range pr.Members {
		members = append(members, openapi.TeamMember{
			IsActive: m.IsActive,
			UserId:   m.UserId,
			Username: m.Username,
		})
	}
	r := openapi.Team{
		Members:  members,
		TeamName: pr.TeamName,
	}
	transport.WriteJSON(w, http.StatusOK, r)
}
