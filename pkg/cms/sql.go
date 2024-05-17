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
