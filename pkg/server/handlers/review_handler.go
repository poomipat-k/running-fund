package server

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/poomipat-k/running-fund/pkg/projects"
	"github.com/poomipat-k/running-fund/pkg/users"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

func (h *ProjectHandler) AddReview(w http.ResponseWriter, r *http.Request) {
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

	var payload projects.AddReviewRequest
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

func (h *ProjectHandler) getCriteriaList() ([]projects.ProjectReviewCriteriaMinimal, error) {
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
