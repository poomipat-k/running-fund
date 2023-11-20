package projects

import (
	"context"
	"fmt"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

func (s *store) AddReview(payload AddReviewRequest, userId int, criteriaList []ProjectReviewCriteriaMinimal) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Create a helper function for preparing failure results.
	fail := func(err error) (int, error) {
		return 0, fmt.Errorf("addReview: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	// Insert improvement if valid
	// insert review
	// insert review_details
	id := 0
	err = tx.QueryRowContext(ctx, addReviewSQL, userId, 2, false, payload.Review.ReviewSummary).Scan(&id)

	if err != nil {
		return fail(err)
	}
	log.Println(criteriaList)
	log.Println("Payload===")
	log.Println(payload)
	log.Println("id:", id)

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return id, nil
}
