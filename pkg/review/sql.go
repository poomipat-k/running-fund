package review

const insertImprovementSQL = `
INSERT INTO improvement (benefit, experience_and_reliability, fund_and_output, project_quality, project_standard, vision_and_image)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
`

const insertReviewSQL = `
INSERT INTO review (user_id, project_history_id, is_interested_person, interested_person_type, created_at, summary, improvement_id, comment) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;
`

const insertReviewDetailsSQL = `
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES 
`

const getProjectCriteriaMinimalSQL = `
SELECT id, criteria_version ,order_number FROM review_criteria WHERE criteria_version = $1 ORDER BY order_number ASC;
`

const updateProjectStatusToReviewed = `
UPDATE project_history
SET status = 'Reviewed', updated_at = $2
WHERE project_history.id = $1 AND (SELECT COUNT(*) as review_count
FROM review INNER JOIN project_history ON project_history.id = review.project_history_id
WHERE project_history.id = $1
) = $3 RETURNING id;
`
