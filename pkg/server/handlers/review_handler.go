package server

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/poomipat-k/running-fund/pkg/projects"
)

func (h *ProjectHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	// To check if the user exists in the db
	userId, err := getAuthUserId(r)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err, http.StatusForbidden)
		return
	}

	_, err = h.uStore.GetReviewerById(userId)
	if err != nil {
		slog.Error("Don't have reviewer permission")
		errorJSON(w, err, http.StatusForbidden)
		return
	}

	var payload projects.AddReviewRequest
	err = readJSON(w, r, &payload)

	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err)
		return
	}

	criteriaList, err := h.getCriteriaList()
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate Payload
	err = validateAddPayload(payload, criteriaList)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	id, err := h.store.AddReview(payload, userId, criteriaList)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err)
		return
	}
	writeJSON(w, http.StatusOK, id)
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
