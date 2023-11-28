package review

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/poomipat-k/running-fund/pkg/users"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

type reviewStore interface {
	AddReview(payload AddReviewRequest, userId int, criteriaList []ProjectReviewCriteriaMinimal) (int, error)
	GetProjectCriteriaMinimalDetails(cv int) ([]ProjectReviewCriteriaMinimal, error)
}

type ReviewHandler struct {
	store  reviewStore
	uStore users.UserStore
}

func NewProjectHandler(s reviewStore, uStore users.UserStore) *ReviewHandler {
	return &ReviewHandler{
		store:  s,
		uStore: uStore,
	}
}

func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	// To check if the user exists in the db
	userId, err := users.GetAuthUserId(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	_, err = h.uStore.GetReviewerById(userId)
	if err != nil {
		slog.Error("Don't have reviewer permission")
		utils.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	var payload AddReviewRequest
	err = utils.ReadJSON(w, r, &payload)

	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err)
		return
	}

	criteriaList, err := h.getCriteriaList()
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate Payload
	err = validateAddPayload(payload, criteriaList)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	id, err := h.store.AddReview(payload, userId, criteriaList)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, id)
}

func (h *ReviewHandler) getCriteriaList() ([]ProjectReviewCriteriaMinimal, error) {
	cv := os.Getenv("CRITERIA_VERSION")
	v, err := strconv.Atoi(cv)
	if err != nil {
		return nil, errors.New("failed to convert CRITERIA_VERSION to int")
	}
	criteriaList, err := h.store.GetProjectCriteriaMinimalDetails(v)
	if err != nil {
		return nil, err
	}
	return criteriaList, nil
}
