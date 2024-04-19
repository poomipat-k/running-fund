package projects

const getReviewPeriodSQL = "SELECT id, from_date, to_date FROM review_period ORDER BY id DESC LIMIT 1;"

const getReviewerDashboardSQL = `
SELECT 
project.user_id as user_id,
project.id as project_id, project.project_code, 
project.created_at, project_history.project_name, 
review.id as review_id, review.created_at as reviewed_at
FROM project
INNER JOIN project_history
ON project.project_history_id = project_history.id
LEFT JOIN review
ON project.project_history_id = review.project_history_id AND review.user_id = $1
WHERE project.created_at >= $2
AND project.created_at < $3
ORDER BY project_name;
`
const getReviewerProjectDetailsSQL = `
SELECT 
project.user_id as user_id,
project.id as project_id, 
project_history.id as project_history_id, 
project.project_code, 
project.created_at as project_created_at, 
project_history.project_name,

contact.prefix as project_head_prefix,
contact.first_name as project_head_first_name,
contact.last_name as project_head_last_name,
project_history.from_date as from_date,
project_history.to_date as to_date,
address.address as address_details,
province.name as province_name,
district.name as district_name,
subdistrict.name as subdistrict_name,
distance.type as distance_type,
distance.is_dynamic as distance_is_dynamic,
project_history.expected_participants as expected_participants,
project_history.collaborated as collaborated,

review.id as review_id, 
review.created_at as reviewed_at, 
review.is_interested_person, 
review.interested_person_type,
review.summary as review_summary, 
review.comment as reviewer_comment, 
improvement.benefit, 
improvement.experience_and_reliability,
improvement.fund_and_output, improvement.project_quality, improvement.project_standard,
improvement.vision_and_image
FROM project
INNER JOIN project_history ON project.project_history_id = project_history.id
INNER JOIN contact ON project_history.project_head_contact_id = contact.id
INNER JOIN address ON project_history.address_id = address.id
INNER JOIN postcode ON address.postcode_id = postcode.id
INNER JOIN subdistrict ON postcode.subdistrict_id = subdistrict.id
INNER JOIN district ON subdistrict.district_id = district.id
INNER JOIN province ON district.province_id = province.id
INNER JOIN distance ON project_history.id = distance.project_history_id

LEFT JOIN review ON project.project_history_id = review.project_history_id AND review.user_id = $1
LEFT JOIN improvement ON review.improvement_id = improvement.id
WHERE project.project_code = $2;
`

const getProjectCriteriaSQL = `
SELECT criteria_version ,order_number, group_number, in_group_number, display_text
FROM review_criteria WHERE criteria_version = $1 ORDER BY order_number ASC;
`

const getReviewDetailsByReviewIdSQL = `
SELECT review_details.id as review_details_id, review_criteria.criteria_version,
review_criteria.order_number as criteria_order_number, review_details.score
FROM review_details INNER JOIN review_criteria 
ON review_criteria.id = review_details.review_criteria_id 
WHERE review_details.review_id = $1
ORDER BY criteria_order_number ASC;
`

const countProjectCreatedToday = `
SELECT count(*) FROM project 
WHERE created_at >= date_trunc('day', now()) + interval '- 7 hour' 
AND created_at < date_trunc('day', now()) + interval '1 day - 7 hour - 1 microsecond';
`

const getApplicantCriteriaSQL = `
SELECT id, criteria_version, order_number, display
FROM applicant_criteria WHERE criteria_version = $1 AND code = 'project_self_score' 
ORDER BY order_number ASC;
`
const getApplicantCriteriaPdfSQL = `
SELECT id, criteria_version, order_number, pdf_display
FROM applicant_criteria WHERE criteria_version = $1 AND code = 'project_self_score' 
ORDER BY order_number ASC;
`

const addAddressSQL = `
INSERT INTO address (address, postcode_id) VALUES ($1, $2) RETURNING id;
`

const addContactMainSQL = `
INSERT INTO contact 
(prefix, first_name, last_name, organization_position, event_position) 
VALUES ($1, $2, $3, $4, $5) RETURNING id;
`

const addContactFullSQL = `
INSERT INTO contact 
(prefix, first_name, last_name, organization_position, event_position, address_id, email, line_id, phone_number) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;
`

const addContactOnlyRequiredParamSQL = `
INSERT INTO contact 
(prefix, first_name, last_name) 
VALUES ($1, $2, $3) RETURNING id;
`

const addProjectSQL = `
INSERT INTO project
(project_code, created_at, project_history_id, user_id)
VALUES ($1, $2, $3, $4) RETURNING id;
`

const addProjectHistorySQL = `
INSERT INTO project_history 
(
	project_code,
	project_version,
	created_at,
	updated_at,
	status,
	collaborated,
	project_name,
	from_date,
	to_date,
	address_id,
	start_point,
	finish_point,
	cat_road_race,
	cat_trail_running,
	cat_has_other,
	cat_other_type,
	vip,
	vip_fee,
	expected_participants,
	has_organizer,
	organizer_name,
	project_head_contact_id,
	project_manager_contact_id,
	project_coordinator_contact_id,
	project_race_director_contact_id,
	organization_type,
	organization_name,
	background,
	objective,
	mkt_has_facebook,
	mkt_facebook,
	mkt_has_website,
	mkt_website,
	mkt_use_online_page,
	mkt_online_page,
	mkt_use_other_online_marketing,
	mkt_other_online_marketing,
	mkt_pr,
	mkt_local_official,
	mkt_booth,
	mkt_billboard,
	mkt_tv,
	mkt_use_other_offline_marketing,
	mkt_other_offline_marketing,
	st_runner_info,
	st_health_decider,
	st_ambulance,
	st_first_aid,
	st_aed,
	st_aed_count,
	st_volunteer_doctor,
	st_insurance,
	st_other,
	st_addition,
	measure_athletics_association,
	measure_calibrated_bicycle,
	measure_self_measurement,
	measure_self_tool,
	traffic_ask_permission,
	traffic_has_supporter,
	traffic_road_closure,
	traffic_signs,
	traffic_lighting,
	judge_type,
	judge_other_type,
	support_provincial_admin,
	support_safety,
	support_health,
	support_volunteer,
	support_community,
	support_other,
	support_addition,
	feedback,
	exp_this_first_time,
	exp_this_ordinal_number,
	exp_this_latest_date,
	exp_this_completed1_year,
	exp_this_completed1_participant,
	exp_this_completed2_year,
	exp_this_completed2_participant,
	exp_this_completed3_year,
	exp_this_completed3_participant,
	exp_other_done_before,
	exp_other_completed1_year,
	exp_other_completed1_name,
	exp_other_completed1_participant,
	exp_other_completed2_year,
	exp_other_completed2_name,
	exp_other_completed2_participant,
	exp_other_completed3_year,
	exp_other_completed3_name,
	exp_other_completed3_participant,
	fund_total,
	fund_support_organization,
	fund_req_fund,
	fund_req_fund_amount,
	fund_req_bib,
	fund_req_bib_amount,
	fund_req_pr,
	fund_req_seminar,
	fund_req_seminar_topic,
	fund_req_other,
	fund_req_other_type,
	files_prefix
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, 
	$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, 
	$39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, 
	$57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, 
	$75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89, $90, $91, $92, 
	$93, $94, $95, $96, $97, $98, $99, $100, $101, $102, $103, $104
) RETURNING id;
`

const addManyDistanceSQL = `
INSERT INTO distance (type, fee, is_dynamic, project_history_id) VALUES 
`

const addManyApplicantScoreSQL = `
INSERT INTO applicant_score (project_history_id, applicant_criteria_id, score) VALUES 
`
const getAllProjectDashboardByApplicantIdSQL = `
SELECT
project.id as id,
project.project_code as project_code,
project.created_at as project_created_at,
project_history.project_name as project_name,
project_history.status as project_status,
project_history.updated_at as latest_status_date,
project_history.admin_comment as admin_comment
FROM project
INNER JOIN project_history ON project.project_history_id = project_history.id
WHERE project.user_id = $1
ORDER BY project.created_at DESC
;
`

const getApplicantProjectDetailsSQL = `
SELECT 
project.project_code as project_code,
project.user_id as user_id,
project_history.project_name  as project_name,
project_history.status as project_status,
project_history.admin_score as admin_score,
project_history.fund_approved_amount as fund_approved_amount,
project_history.admin_comment as admin_comment,
review.id as review_id,
review.user_id as reviewer_id,
review.created_at as reviewed_at
FROM project
INNER JOIN project_history ON project.project_code = project_history.project_code
LEFT JOIN review ON review.project_history_id = project_history.id
WHERE project.project_code = $1 AND project.user_id = $2
ORDER BY reviewed_at ASC
;
`

const getApplicantProjectDetailsByAdminSQL = `
SELECT 
project.project_code as project_code,
project.user_id as user_id,
project_history.project_name  as project_name,
project_history.status as project_status,
project_history.admin_score as admin_score,
project_history.fund_approved_amount as fund_approved_amount,
project_history.admin_comment as admin_comment,
review.id as review_id,
review.user_id as reviewer_id,
review.created_at as reviewed_at,
SUM(review_details.score)
FROM project
INNER JOIN project_history ON project.project_code = project_history.project_code
LEFT JOIN review ON review.project_history_id = project_history.id
LEFT JOIN review_details ON review.id = review_details.review_id
WHERE project.project_code = $1 
GROUP BY project.project_code, project.user_id, project_history.project_name, 
project_history.status, project_history.admin_score, project_history.fund_approved_amount,
project_history.admin_comment, review.id
ORDER BY reviewed_at ASC
;
`

const hasRightToAddAdditionalFilesSQL = `
SELECT project.id as project_id, project_history.status as project_status
FROM project INNER JOIN project_history ON project.project_history_id = project_history.id 
WHERE project.user_id = $1 AND project.project_code = $2;
`

const getAddressDetailsSQL = `
SELECT postcode.code as postcode, subdistrict."name" as subdistrict_name,
district."name" as district_name, province."name" as province_name
FROM postcode
INNER JOIN subdistrict ON postcode.subdistrict_id = subdistrict.id
INNER JOIN district ON subdistrict.district_id = district.id
INNER JOIN province ON district.province_id = province.id
WHERE postcode.id = $1;
`

const getProjectForAdminUpdateByProjectCodeSQL = `
SELECT
project.user_id as user_id,
project_history.id as project_history_id,
project_history.status as project_status,
project_history.admin_score as admin_score,
project_history.fund_approved_amount as fund_approved_amount,
project_history.admin_comment as admin_comment,
project_history.admin_approved_at as admin_approved_at,
project_history.updated_at as updated_at
FROM project
INNER JOIN project_history ON project.project_history_id = project_history.id
WHERE project.project_code = $1;
`

const updateProjectByAdminSQL = `
UPDATE project_history
SET
status = $2,
admin_score = $3,
fund_approved_amount = $4,
admin_comment = $5,
admin_approved_at = $6,
updated_at = $7
WHERE project_history.id = $1 RETURNING id;
`

const getAdminSummarySQL = `
SELECT 
project_history.status,
COUNT(*) as count,
SUM(project_history.fund_approved_amount)
FROM project INNER JOIN project_history ON project.project_history_id = project_history.id
WHERE project.created_at >= $1 AND project.created_at < $2
GROUP BY project_history.status;
`
