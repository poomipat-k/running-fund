package cms

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

type store struct {
	db *sql.DB
	c  *cache.Cache
}

func NewStore(db *sql.DB, c *cache.Cache) *store {
	return &store{
		db: db,
		c:  c,
	}
}

func (s *store) GetReviewPeriod() (ReviewPeriod, error) {
	var period ReviewPeriod
	row := s.db.QueryRow(getReviewPeriodSQL)
	err := row.Scan(&period.Id, &period.FromDate, &period.ToDate)
	switch err {
	case sql.ErrNoRows:
		slog.Error("GetReviewPeriod(): no row were returned!")
		return ReviewPeriod{}, err
	case nil:
		return period, nil
	default:
		slog.Error(err.Error())
		return ReviewPeriod{}, fmt.Errorf("GetReviewPeriod() unknown error")
	}
}

func (s *store) GetAdminWebsiteDashboardDateConfigPreview(fromDate, toDate time.Time, limit, offset int) ([]AdminDateConfigPreviewRow, error) {
	rows, err := s.db.Query(getAdminWebsiteDashboardDateConfigPreviewSQL, fromDate, toDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []AdminDateConfigPreviewRow
	for rows.Next() {
		var row AdminDateConfigPreviewRow
		err := rows.Scan(
			&row.Count,
			&row.ProjectCode,
			&row.ProjectCreatedAt,
			&row.ProjectName,
			&row.ProjectStatus,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *store) AdminUpdateWebsiteConfig(payload AdminUpdateWebsiteConfigRequest) error {
	// start transaction
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// // create a website_config record
	// webConfigId, err := s.addWebsiteConfig(payload)
	// if err != nil {
	// 	return err
	// }
	// log.Println("==webConfigId", webConfigId)

	// Landing

	// Dashboard
	needUpdateDashboardConfig, err := s.shouldUpdateDashboardConfig(payload.Dashboard)
	if err != nil {
		return err
	}
	if needUpdateDashboardConfig {
		err := s.updateDashboardConfig(ctx, tx, payload.Dashboard)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *store) shouldUpdateDashboardConfig(payload DashboardConfig) (bool, error) {
	period, err := s.GetReviewPeriod()
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	// no record in current review_period table
	if period.FromDate == nil || period.ToDate == nil {
		return true, nil
	}
	// validate that dashboard config values have been changed
	curFromDate := period.FromDate
	curToDate := period.ToDate

	loc, err := utils.GetTimeLocation()
	if err != nil {
		return false, err
	}
	newFromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	newToDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay+1, 0, 0, 0, 0, loc)
	return (!curFromDate.Equal(newFromDate) || !curToDate.Equal(newToDate)), nil
}

func (s *store) updateDashboardConfig(
	ctx context.Context,
	tx *sql.Tx,
	payload DashboardConfig,
) error {
	loc, err := utils.GetTimeLocation()
	if err != nil {
		return err
	}
	newFromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	newToDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay+1, 0, 0, 0, 0, loc)
	var id int
	err = tx.QueryRowContext(
		ctx,
		adminUpdateReviewerPeriodSQL,
		newFromDate,
		newToDate,
	).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
