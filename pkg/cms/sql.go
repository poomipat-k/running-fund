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
INSERT INTO website_config (landing_page, footer_email, footer_phone_number, footer_operate_hour)
VALUES ($1, $2, $3, $4) RETURNING id;
`

const getLandingPageContentSQL = `
SELECT
website_config.id as id,
landing_page as landing_page
FROM website_config
ORDER BY website_config.id DESC LIMIT 1
;
`
const getLatestWebsiteConfigIdSQL = `
SELECT
website_config.id as id
FROM website_config
ORDER BY website_config.id DESC LIMIT 1
;
`

const getLandingPageBannerSQL = `
SELECT  website_image.id, website_image.full_path, website_image.object_key, website_image.link_to
FROM website_image WHERE website_image.website_config_id = $1 AND website_image.code = $2
ORDER BY website_image.order_number;
`

const getFaqSQL = `
SELECT  faq.id, faq.question, faq.answer
FROM faq WHERE faq.website_config_id = $1
ORDER BY faq.order_number;
`

const getHowToCreateSQL = `
SELECT  how_to_create.id, how_to_create.header, how_to_create.content
FROM how_to_create WHERE how_to_create.website_config_id = $1
ORDER BY how_to_create.order_number;
`

const getFooterLogoSQL = `
SELECT  
website_image.id, 
website_image.full_path, 
website_image.object_key, 
website_image.link_to
FROM website_image WHERE website_image.website_config_id = $1 AND website_image.code = $2
ORDER BY website_image.order_number;
`

const getLatestWebsiteConfigWithFooterSQL = `
SELECT
website_config.id as id,
website_config.footer_email as footer_email,
website_config.footer_phone_number as footer_phone_number,
website_config.footer_operate_hour as footer_operate_hour
FROM website_config
ORDER BY website_config.id DESC LIMIT 1
;
`

const updateOperationConfigSQL = `
UPDATE operation_config SET allow_new_project = $1 WHERE id = (select id FROM operation_config ORDER BY id DESC limit 1) RETURNING id;
`
