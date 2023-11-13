package projects

const getReviewPeriodSQL = "SELECT id, from_date, to_date FROM review_period ORDER BY id DESC LIMIT 1;"

const getReviewerDashboardSQL = `
SELECT project.id as project_id, project.project_code, 
project.created_at, project_history.project_name, 
review.id as review_id, review.created_at as reviewed_at,
project_history.download_link
FROM project
INNER JOIN project_history
ON project.project_history_id = project_history.id
LEFT JOIN review
ON project.project_history_id = review.project_history_id AND user_id = $1
WHERE project.created_at >= $2
AND project.created_at < $3
ORDER BY project_name;
`
const getReviewerProejctDetailsSQL = `
SELECT project.id as project_id, project.project_code, project_history.project_name, 
review.id as review_id, review.created_at as reviewed_at
FROM project
INNER JOIN project_history ON project.project_history_id = project_history.id
LEFT JOIN review ON project.project_history_id = review.project_history_id AND user_id = $1
WHERE project.project_code = $2
LIMIT 1;
`
const getProjectCriteriaSQL = `
SELECT criteria_version ,order_number, group_number, in_group_number, display_text
FROM review_criteria WHERE criteria_version = $1;
`
