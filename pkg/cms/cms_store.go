package cms

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const CONTENT_LANDING_PAGE_CACHE_KEY = "content_landing_page"

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

func (s *store) GetLandingPageContent() (LandingConfig, error) {
	// check cache first
	raw, found := s.c.Get(CONTENT_LANDING_PAGE_CACHE_KEY)
	if found {
		cachedData, ok := raw.(LandingConfig)
		if ok {
			return cachedData, nil
		}
	}

	// Fetch data from the db
	var configId int
	var landingContent string
	row := s.db.QueryRow(getLandingPageContentSQL)
	err := row.Scan(&configId, &landingContent)
	if err != nil {
		return LandingConfig{}, err
	}
	// Fetch data from the db
	rows, err := s.db.Query(getLandingPageBannerSQL, configId)
	if err != nil {
		return LandingConfig{}, err
	}
	defer rows.Close()

	var banners []Banner
	for rows.Next() {
		var row Banner
		err := rows.Scan(&row.FullPath, &row.ObjectKey, &row.LinkTo)
		if err != nil {
			return LandingConfig{}, err
		}
		banners = append(banners, row)
	}

	err = rows.Err()
	if err != nil {
		return LandingConfig{}, err
	}
	responseBody := LandingConfig{
		WebsiteConfigId: configId,
		Content:         landingContent,
		Banner:          banners,
	}
	// Set cache
	s.c.Set(CONTENT_LANDING_PAGE_CACHE_KEY, responseBody, cache.NoExpiration)

	return responseBody, nil
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
	// create a website_config record
	webConfigId, err := s.addWebsiteConfig(ctx, tx, payload)
	if err != nil {
		return err
	}
	log.Println("==webConfigId", webConfigId)
	// Landing banner
	bannerAddedCount, err := s.addLandingPageBanners(ctx, tx, payload.Landing.Banner, webConfigId)
	if err != nil {
		return err
	}
	log.Println("===bannerAddedCount", bannerAddedCount)

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
	// Remove cache to have it refreshed later on the first visit
	s.c.Delete(CONTENT_LANDING_PAGE_CACHE_KEY)
	return nil
}

func (s *store) addWebsiteConfig(ctx context.Context, tx *sql.Tx, payload AdminUpdateWebsiteConfigRequest) (int, error) {
	var id int
	err := tx.QueryRowContext(ctx, addWebsiteConfigSQL, payload.Landing.Content).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *store) addLandingPageBanners(ctx context.Context, tx *sql.Tx, banners []Banner, websiteConfigId int) (int64, error) {
	const initialSQL = `
	INSERT INTO banner (full_path, object_key, link_to, website_config_id)
	VALUES 
	`

	valuesStrPlaceholder := []string{}
	values := []any{}
	for i, banner := range banners {
		valuesStrPlaceholder = append(valuesStrPlaceholder, fmt.Sprintf("($%d, $%d, $%d, $%d)", 4*i+1, 4*i+2, 4*i+3, 4*i+4))
		values = append(values, banner.FullPath, banner.ObjectKey, banner.LinkTo, websiteConfigId)
	}
	customSQL := initialSQL + strings.Join(valuesStrPlaceholder, ",") + ";"
	stmt, err := tx.Prepare(customSQL)
	if err != nil {
		slog.Error("error prepare add banners sql", "error", err)
		return 0, err
	}
	result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		slog.Error("execContext on add banners sql", "error", err)
		return 0, err
	}
	return result.RowsAffected()
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
