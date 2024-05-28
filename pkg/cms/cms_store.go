package cms

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const CONTENT_LANDING_PAGE_CACHE_KEY = "content_landing_page"
const CONTENT_CMS_DATA_CACHE_KEY = "content_cms"
const CONTENT_FAQ_CACHE_KEY = "content_faq"
const CONTENT_HOW_TO_CREATE_CACHE_KEY = "content_how_to_create"
const CONTENT_FOOTER_CACHE_KEY = "content_footer"

var ReCacheOnUpdateCmsData = []string{
	CONTENT_LANDING_PAGE_CACHE_KEY,
	CONTENT_CMS_DATA_CACHE_KEY,
	CONTENT_FAQ_CACHE_KEY,
	CONTENT_HOW_TO_CREATE_CACHE_KEY,
	CONTENT_FOOTER_CACHE_KEY,
}

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
	// Fetch banners from the db
	banners, err := s.getLandingPageBanners(configId)
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

// Website cms data
func (s *store) GetWebsiteConfigData() (AdminUpdateWebsiteConfigRequest, string, error) {
	// check cache first
	raw, found := s.c.Get(CONTENT_CMS_DATA_CACHE_KEY)
	if found {
		cachedData, ok := raw.(AdminUpdateWebsiteConfigRequest)
		if ok {
			return cachedData, "cache", nil
		}
	}
	// Fetch data from db
	// dashboard
	period, err := s.GetReviewPeriod()
	if err != nil && err != sql.ErrNoRows {
		return AdminUpdateWebsiteConfigRequest{}, "reviewPeriod", err
	}
	fromDate := *period.FromDate
	toDate := *period.ToDate
	loc, err := utils.GetTimeLocation()
	if err != nil {
		return AdminUpdateWebsiteConfigRequest{}, "timeZone loc", err
	}
	locFromDate := fromDate.In(loc)
	locToDate := toDate.Add(time.Duration(-1 * time.Minute)).In(loc)

	// Check website config table is not empty
	var websiteConfigId int
	row := s.db.QueryRow(getLatestWebsiteConfigIdSQL)
	err = row.Scan(&websiteConfigId)

	if err == sql.ErrNoRows {
		data := AdminUpdateWebsiteConfigRequest{
			Dashboard: DashboardConfig{
				FromYear:  locFromDate.Year(),
				FromMonth: int(locFromDate.Month()),
				FromDay:   locFromDate.Day(),
				ToYear:    locToDate.Year(),
				ToMonth:   int(locToDate.Month()),
				ToDay:     locToDate.Day(),
			},
		}
		// Set cache
		s.c.Set(CONTENT_CMS_DATA_CACHE_KEY, data, cache.NoExpiration)
		return data, "", nil
	}
	if err != nil {
		return AdminUpdateWebsiteConfigRequest{}, "", err
	}

	// landing
	landingPage, err := s.GetLandingPageContent()
	if err != nil {
		return AdminUpdateWebsiteConfigRequest{}, "landingPage", err
	}

	// faq
	faqList, err := s.GetFAQ()
	if err != nil {
		return AdminUpdateWebsiteConfigRequest{}, "faq", err
	}
	// howToCreate
	howToCreateList, err := s.GetHowToCreate()
	if err != nil {
		return AdminUpdateWebsiteConfigRequest{}, "howToCreate", err
	}

	// footer
	footer, err := s.GetFooter()
	if err != nil {
		return AdminUpdateWebsiteConfigRequest{}, "faq", err
	}
	operatingHour := getOperatingHourFromStr(footer.Contact.OperatingHour)

	cmsData := AdminUpdateWebsiteConfigRequest{
		Landing: landingPage,
		Dashboard: DashboardConfig{
			FromYear:  locFromDate.Year(),
			FromMonth: int(locFromDate.Month()),
			FromDay:   locFromDate.Day(),
			ToYear:    locToDate.Year(),
			ToMonth:   int(locToDate.Month()),
			ToDay:     locToDate.Day(),
		},
		Faq:         faqList,
		HowToCreate: howToCreateList,
		Footer: FooterConfig{
			Logo: footer.Logo,
			Contact: FooterContact{
				Email:       footer.Contact.Email,
				PhoneNumber: footer.Contact.PhoneNumber,
				FromHour:    operatingHour[0],
				FromMinute:  operatingHour[1],
				ToHour:      operatingHour[2],
				ToMinute:    operatingHour[3],
			},
		},
	}
	// Set cache
	s.c.Set(CONTENT_CMS_DATA_CACHE_KEY, cmsData, cache.NoExpiration)

	return cmsData, "", nil
}

func getOperatingHourFromStr(s string) []string {
	trim := strings.Trim(s, " น.")
	split := strings.Split(trim, " - ")
	if len(split) != 2 {
		return []string{"", "", "", ""}
	}
	result := []string{}
	for _, str := range split {
		sp := strings.Split(str, ".")
		result = append(result, sp...)
	}
	if len(result) != 4 {
		return []string{"", "", "", ""}
	}
	return result
}

func (s *store) GetFooter() (FooterResponse, error) {
	// check cache first
	raw, found := s.c.Get(CONTENT_FOOTER_CACHE_KEY)
	if found {
		cachedData, ok := raw.(FooterResponse)
		if ok {
			return cachedData, nil
		}
	}

	var footerData FooterResponse
	row := s.db.QueryRow(getLatestWebsiteConfigWithFooterSQL)
	err := row.Scan(&footerData.Id, &footerData.Contact.Email, &footerData.Contact.PhoneNumber, &footerData.Contact.OperatingHour)
	if err != nil {
		return FooterResponse{}, err
	}

	// logos
	rows, err := s.db.Query(getFooterLogoSQL, footerData.Id, "footer_logo")
	if err != nil {
		return FooterResponse{}, err
	}
	defer rows.Close()

	var logoList []Image
	for rows.Next() {
		var row Image
		err := rows.Scan(&row.Id, &row.FullPath, &row.ObjectKey, &row.LinkTo)
		if err != nil {
			return FooterResponse{}, err
		}
		logoList = append(logoList, row)
	}
	err = rows.Err()
	if err != nil {
		return FooterResponse{}, err
	}

	footerData.Logo = logoList
	// Set cache
	s.c.Set(CONTENT_FOOTER_CACHE_KEY, footerData, cache.NoExpiration)
	return footerData, nil
}

func (s *store) GetFAQ() ([]FAQ, error) {
	// check cache first
	raw, found := s.c.Get(CONTENT_FAQ_CACHE_KEY)
	if found {
		cachedData, ok := raw.([]FAQ)
		if ok {
			return cachedData, nil
		}
	}

	// fetch from db
	var websiteConfigId int
	row := s.db.QueryRow(getLatestWebsiteConfigIdSQL)
	err := row.Scan(&websiteConfigId)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(getFaqSQL, websiteConfigId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faqList []FAQ
	for rows.Next() {
		var row FAQ
		err := rows.Scan(&row.Id, &row.Question, &row.Answer)
		if err != nil {
			return nil, err
		}
		faqList = append(faqList, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Set cache
	s.c.Set(CONTENT_FAQ_CACHE_KEY, faqList, cache.NoExpiration)
	return faqList, nil
}

func (s *store) GetHowToCreate() ([]HowToCreate, error) {
	// check cache first
	raw, found := s.c.Get(CONTENT_HOW_TO_CREATE_CACHE_KEY)
	if found {
		cachedData, ok := raw.([]HowToCreate)
		if ok {
			return cachedData, nil
		}
	}

	// fetch from db
	var websiteConfigId int
	row := s.db.QueryRow(getLatestWebsiteConfigIdSQL)
	err := row.Scan(&websiteConfigId)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(getHowToCreateSQL, websiteConfigId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var howToCreateList []HowToCreate
	for rows.Next() {
		var row HowToCreate
		err := rows.Scan(&row.Id, &row.Header, &row.Content)
		if err != nil {
			return nil, err
		}
		howToCreateList = append(howToCreateList, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Set cache
	s.c.Set(CONTENT_HOW_TO_CREATE_CACHE_KEY, howToCreateList, cache.NoExpiration)
	return howToCreateList, nil
}

func (s *store) getLandingPageBanners(websiteConfigId int) ([]Image, error) {
	rows, err := s.db.Query(getLandingPageBannerSQL, websiteConfigId, "banner")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []Image
	for rows.Next() {
		var row Image
		err := rows.Scan(&row.Id, &row.FullPath, &row.ObjectKey, &row.LinkTo)
		if err != nil {
			return nil, err
		}
		banners = append(banners, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return banners, nil
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
	// Landing banner
	if len(payload.Landing.Banner) > 0 {
		_, err := s.addLandingPageBanners(ctx, tx, payload.Landing.Banner, webConfigId)
		if err != nil {
			return err
		}
	}
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

	// FAQ
	if len(payload.Faq) > 0 {
		_, err := s.addFaqList(ctx, tx, payload.Faq, webConfigId)
		if err != nil {
			return err
		}
	}
	// HowToCreate
	if len(payload.HowToCreate) > 0 {
		_, err := s.addHowToCreate(ctx, tx, payload.HowToCreate, webConfigId)
		if err != nil {
			return err
		}
	}

	// Footer.Logo
	if len(payload.Footer.Logo) > 0 {
		_, err := s.addFooterLogos(ctx, tx, payload.Footer.Logo, webConfigId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	// Remove cache to have it refreshed later on the first visit
	for _, cacheKey := range ReCacheOnUpdateCmsData {
		s.c.Delete(cacheKey)
	}
	return nil
}

func (s *store) addWebsiteConfig(ctx context.Context, tx *sql.Tx, payload AdminUpdateWebsiteConfigRequest) (int, error) {
	var id int
	operatingHour := fmt.Sprintf("%s.%s - %s.%s น.",
		payload.Footer.Contact.FromHour,
		payload.Footer.Contact.FromMinute,
		payload.Footer.Contact.ToHour,
		payload.Footer.Contact.ToMinute,
	)
	err := tx.QueryRowContext(
		ctx,
		addWebsiteConfigSQL,
		payload.Landing.Content,
		payload.Footer.Contact.Email,
		payload.Footer.Contact.PhoneNumber,
		operatingHour,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *store) addLandingPageBanners(ctx context.Context, tx *sql.Tx, banners []Image, websiteConfigId int) (int64, error) {
	const initialSQL = `
	INSERT INTO website_image (code, full_path, object_key, link_to, order_number, website_config_id)
	VALUES 
	`

	valuesStrPlaceholder := []string{}
	values := []any{}
	for i, banner := range banners {
		valuesStrPlaceholder = append(valuesStrPlaceholder, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", 6*i+1, 6*i+2, 6*i+3, 6*i+4, 6*i+5, 6*i+6))
		values = append(values, "banner", banner.FullPath, banner.ObjectKey, banner.LinkTo, i+1, websiteConfigId)
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

func (s *store) addFooterLogos(ctx context.Context, tx *sql.Tx, logos []Image, websiteConfigId int) (int64, error) {
	const initialSQL = `
	INSERT INTO website_image (code, full_path, object_key, link_to, order_number, website_config_id)
	VALUES 
	`

	valuesStrPlaceholder := []string{}
	values := []any{}
	for i, logo := range logos {
		valuesStrPlaceholder = append(valuesStrPlaceholder, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", 6*i+1, 6*i+2, 6*i+3, 6*i+4, 6*i+5, 6*i+6))
		values = append(values, "footer_logo", logo.FullPath, logo.ObjectKey, logo.LinkTo, i+1, websiteConfigId)
	}
	customSQL := initialSQL + strings.Join(valuesStrPlaceholder, ",") + ";"
	stmt, err := tx.Prepare(customSQL)
	if err != nil {
		slog.Error("error prepare add logos sql", "error", err)
		return 0, err
	}
	result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		slog.Error("execContext on add logos sql", "error", err)
		return 0, err
	}
	return result.RowsAffected()
}

func (s *store) addFaqList(ctx context.Context, tx *sql.Tx, faqList []FAQ, websiteConfigId int) (int64, error) {
	const initialSQL = `
	INSERT INTO faq (question, answer, order_number, website_config_id)
	VALUES 
	`

	valuesStrPlaceholder := []string{}
	values := []any{}
	for i, faq := range faqList {
		valuesStrPlaceholder = append(valuesStrPlaceholder, fmt.Sprintf("($%d, $%d, $%d, $%d)", 4*i+1, 4*i+2, 4*i+3, 4*i+4))
		values = append(values, faq.Question, faq.Answer, i+1, websiteConfigId)
	}
	customSQL := initialSQL + strings.Join(valuesStrPlaceholder, ",") + ";"
	stmt, err := tx.Prepare(customSQL)
	if err != nil {
		slog.Error("error prepare add faqs sql", "error", err)
		return 0, err
	}
	result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		slog.Error("execContext on add faqs sql", "error", err)
		return 0, err
	}
	return result.RowsAffected()
}

func (s *store) addHowToCreate(ctx context.Context, tx *sql.Tx, howToCreateList []HowToCreate, websiteConfigId int) (int64, error) {
	const initialSQL = `
	INSERT INTO how_to_create (header, content, order_number, website_config_id)
	VALUES 
	`

	valuesStrPlaceholder := []string{}
	values := []any{}
	for i, item := range howToCreateList {
		valuesStrPlaceholder = append(valuesStrPlaceholder, fmt.Sprintf("($%d, $%d, $%d, $%d)", 4*i+1, 4*i+2, 4*i+3, 4*i+4))
		values = append(values, item.Header, item.Content, i+1, websiteConfigId)
	}
	customSQL := initialSQL + strings.Join(valuesStrPlaceholder, ",") + ";"
	stmt, err := tx.Prepare(customSQL)
	if err != nil {
		slog.Error("error prepare add howToCreate sql", "error", err)
		return 0, err
	}
	result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		slog.Error("execContext on add howToCreate sql", "error", err)
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
