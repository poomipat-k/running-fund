-- users
INSERT INTO users (first_name, last_name, email, user_role)
VALUES ('reviewer1', 'x', 'r1@test.com', 'reviewer');

INSERT INTO users (first_name, last_name, email, user_role)
VALUES ('reviewer2', 'x', 'r2@test.com', 'reviewer');

INSERT INTO users (first_name, last_name, email, user_role)
VALUES ('user1', 'x', 'u1@test.com', 'user');

-- project
-- 2023-11-12 16:47:25.152158+00
INSERT INTO project (project_code, created_at)
VALUES ('OCT66_15', '2023-10-15 19:47:25.152158+00');

INSERT INTO project (project_code, created_at)
VALUES ('OCT66_16', '2023-10-16 20:47:25.152158+00');

INSERT INTO project (project_code, created_at)
VALUES ('OCT66_20', '2023-10-20 22:47:25.152158+00');


-- project_history
INSERT INTO project_history (project_code, project_name, created_at, download_link, admin_comment, project_id)
VALUES ('OCT66_15', 'Project_A', '2023-10-15 19:47:25.152158+00', 'https://google.com', 'Admin comment1', (SELECT id FROM project where project_code = 'OCT66_15'));

UPDATE project SET project_history_id = (
    SELECT id FROM project_history WHERE project_code = 'OCT66_15'
    AND project_version = (SELECT project_version FROM project WHERE project_code = 'OCT66_15')) 
    WHERE project_code = 'OCT66_15' ;


INSERT INTO project_history (project_code, project_name, created_at, download_link, admin_comment, project_id)
VALUES ('OCT66_16', 'Project_B', '2023-10-16 20:47:25.152158+00', 'https://youtube.com', 'Admin comment2', (SELECT id FROM project where project_code = 'OCT66_16'));

UPDATE project SET project_history_id = (
    SELECT id FROM project_history WHERE project_code = 'OCT66_16'
    AND project_version = (SELECT project_version FROM project WHERE project_code = 'OCT66_16')) 
    WHERE project_code = 'OCT66_16' ;



INSERT INTO project_history (project_code, project_name, created_at, download_link, admin_comment, project_id)
VALUES ('OCT66_20', 'Project_A', '2023-10-20 22:47:25.152158+00', NULL, NULL, (SELECT id FROM project where project_code = 'OCT66_20'));

UPDATE project SET project_history_id = (
    SELECT id FROM project_history WHERE project_code = 'OCT66_20'
    AND project_version = (SELECT project_version FROM project WHERE project_code = 'OCT66_20')) 
    WHERE project_code = 'OCT66_20' ;
-- END project_history

