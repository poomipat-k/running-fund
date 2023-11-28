package projects

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const dbTimeout = time.Second * 5

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

	now := time.Now()
	// Insert improvement if reviewSummary = to_be_revised
	improvementId := 0
	if payload.Review.ReviewSummary == "to_be_revised" {
		err = tx.QueryRowContext(
			ctx,
			insertImprovementSQL,
			payload.Review.Improvement.Benefit,
			payload.Review.Improvement.ExperienceAndReliability,
			payload.Review.Improvement.FundAndOutput,
			payload.Review.Improvement.ProjectQuality,
			payload.Review.Improvement.ProjectStandard,
			payload.Review.Improvement.VisionAndImage,
		).Scan(&improvementId)
		if err != nil {
			return fail(err)
		}
	}
	// insert review
	reviewId := 0
	sqlImprovementId := sql.NullInt64{}
	if improvementId != 0 {
		sqlImprovementId = sql.NullInt64{
			Valid: true,
			Int64: int64(improvementId),
		}
	}
	err = tx.QueryRowContext(
		ctx,
		insertReviewSQL,
		userId,
		payload.ProjectHistoryId,
		payload.Ip.IsInterestedPerson,
		payload.Ip.InterestedPersonType,
		now,
		payload.Review.ReviewSummary,
		sqlImprovementId,
		payload.Comment,
	).Scan(&reviewId)
	if err != nil {
		return fail(err)
	}

	// insert review_details
	valuesString := []string{}
	values := []any{}

	for i := 0; i < len(criteriaList); i++ {
		valuesString = append(valuesString, fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3))
		scoreName := fmt.Sprintf("q_%d_%d", criteriaList[i].CriteriaVersion, criteriaList[i].OrderNumber)
		score, exist := payload.Review.Scores[scoreName]
		if !exist {
			return fail(fmt.Errorf("score %s is not exist", scoreName))
		}
		values = append(values, reviewId, criteriaList[i].CriteriaId, score)
	}
	customSQL := insertReviewDetailsSQL + strings.Join(valuesString, ",") + ";"

	stmt, err := tx.Prepare(customSQL)
	if err != nil {
		return fail(err)
	}
	_, err = stmt.ExecContext(ctx, values...)
	if err != nil {
		return fail(err)
	}

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}
	return reviewId, nil
}
