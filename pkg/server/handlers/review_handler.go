package server

import "net/http"

type reviewerStore interface {
	AddReview()
}

type ReviewHandler struct {
	store reviewerStore
}

func NewReviewHandler(s reviewerStore) *ReviewHandler {
	return &ReviewHandler{
		store: s,
	}
}

func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Add review"))
}
