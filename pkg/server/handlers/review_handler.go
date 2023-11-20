package server

import (
	"log/slog"
	"net/http"

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
	// Validate Payload
	err = validateAddPayload(payload, h.store)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, payload)
}
