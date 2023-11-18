package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/poomipat-k/running-fund/pkg/review"
)

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
	// To check if the user exists in the db
	_, err := getAuthUserId(r)
	if err != nil {
		log.Println("Panic 1")
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	var payload review.AddReviewRequest
	err = decoder.Decode(&payload)
	if err != nil {
		log.Println("Panic 2")
		panic(err)
	}
	log.Println(payload)

	cv := os.Getenv("CRITERIA_VERSION")
	log.Println("cv:", cv)
	w.Write([]byte("Add review"))
}
