package cms

const getReviewPeriodSQL = "SELECT id, from_date, to_date FROM review_period ORDER BY id DESC LIMIT 1;"

const getAdminWebsiteDashboardDateConfigPreviewSQL = `
SELECT
(
	SELECT COUNT(*) FROM project 
	INNER JOIN project_history ON project.project_history_id = project_history.id
	WHERE project.created_at >= $1 AND project.created_at < $2
) as count,
project.project_code as project_code,
project.created_at as created_at,
project_history.project_name as project_name,
project_history.status as project_status
FROM project 
INNER JOIN project_history ON project.project_history_id = project_history.id
WHERE project.created_at >= $1 AND project.created_at < $2
ORDER BY project.created_at ASC
LIMIT $3 OFFSET $4
;`

const adminUpdateReviewerPeriodSQL = `
INSERT INTO review_period (from_date, to_date)
VALUES ($1, $2) RETURNING id;
`

const addWebsiteConfigSQL = `
INSERT INTO website_config (landing_page)
VALUES ($1) RETURNING id;
`

const getLandingPageContentSQL = `
SELECT
website_config.id as id,
landing_page as landing_page
FROM website_config
ORDER BY website_config.id DESC LIMIT 1
;
`
const getLandingPageBannerSQL = `
SELECT  banner.full_path, banner.object_key, banner.link_to
FROM banner WHERE banner.website_config_id = $1;
`
