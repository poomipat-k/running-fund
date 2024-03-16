package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const dbTimeout = time.Second * 5

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

func (s *store) AddReview(payload AddReviewRequest, userId int, criteriaList []ProjectReviewCriteriaMinimal) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

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
	valuesStrStatement := []string{}
	values := []any{}

	for i := 0; i < len(criteriaList); i++ {
		valuesStrStatement = append(valuesStrStatement, fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3))
		scoreName := fmt.Sprintf("q_%d_%d", criteriaList[i].CriteriaVersion, criteriaList[i].OrderNumber)
		score, exist := payload.Review.Scores[scoreName]
		if !exist {
			return fail(fmt.Errorf("score %s is not exist", scoreName))
		}
		values = append(values, reviewId, criteriaList[i].CriteriaId, score)
	}
	customSQL := insertReviewDetailsSQL + strings.Join(valuesStrStatement, ",") + ";"

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

func (s *store) GetProjectCriteriaMinimalDetails(cv int) ([]ProjectReviewCriteriaMinimal, error) {
	if cv == 0 {
		cv = 1
	}
	rows, err := s.db.Query(getProjectCriteriaMinimalSQL, cv)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ProjectReviewCriteriaMinimal
	for rows.Next() {
		var row ProjectReviewCriteriaMinimal

		err := rows.Scan(&row.CriteriaId, &row.CriteriaVersion, &row.OrderNumber)
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	// get any error occur during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("criteria version not found")
	}
	return data, nil
}

func fail(err error) (int, error) {
	return 0, fmt.Errorf("addReview: %w", err)
}
