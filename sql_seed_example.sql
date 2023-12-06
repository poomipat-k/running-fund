-- users
INSERT INTO users (
        first_name,
        last_name,
        email,
        password,
        user_role,
        activated
    )
VALUES (
        'John',
        'Doe',
        'reviewer1@test.com',
        '$2a$10$TakkT.8E/YedwJ1iUAz7/OeZM6WTaXQnYVEb3UNS8yu2Hq9lE2vO._ayc1TQyw',
        'reviewer',
        true
    );
INSERT INTO users (
        first_name,
        last_name,
        email,
        password,
        user_role,
        activated
    )
VALUES (
        'Jane',
        'Doe',
        'reviewer2@test.com',
        '$2a$10$pmoaZfX/NBdBl9fhuammOeuS9YpbaUwqccaNIzFlfLLPYBiNhjfyu_AIStMEmR',
        'reviewer',
        true
    );
INSERT INTO users (
        first_name,
        last_name,
        email,
        password,
        user_role,
        activated
    )
VALUES (
        'applicant1',
        'test',
        'applicant1@test.com',
        '$2a$10$tm0vILIvIUUaKz0dlf3Fl.4vg0bs9WLwYW2zyog2pudjMT7yloWN2_jjQ3JlHC',
        'applicant',
        true
    );
-- project
-- 2023-11-12 16:47:25.152158+00
INSERT INTO project (project_code, created_at)
VALUES ('OCT66_16', '2023-10-15 19:47:25.152158+00');
INSERT INTO project (project_code, created_at)
VALUES ('OCT66_17', '2023-10-16 20:47:25.152158+00');
INSERT INTO project (project_code, created_at)
VALUES ('OCT66_21', '2023-10-20 22:47:25.152158+00');
INSERT INTO project (project_code, created_at)
VALUES ('NOV66_10', '2023-11-10 10:47:25.152158+00');
INSERT INTO project (project_code, created_at)
VALUES ('NOV66_18', '2023-11-18 11:47:25.152158+00');
-- project_history
INSERT INTO project_history (
        project_code,
        project_name,
        created_at,
        download_link,
        admin_comment,
        project_id
    )
VALUES (
        'OCT66_16',
        'Project_A',
        '2023-10-15 19:47:25.152158+00',
        'https://google.com',
        'Admin comment1',
        (
            SELECT id
            FROM project
            where project_code = 'OCT66_16'
        )
    );
UPDATE project
SET project_history_id = (
        SELECT id
        FROM project_history
        WHERE project_code = 'OCT66_16'
            AND project_version = (
                SELECT project_version
                FROM project
                WHERE project_code = 'OCT66_16'
            )
    )
WHERE project_code = 'OCT66_16';
INSERT INTO project_history (
        project_code,
        project_name,
        created_at,
        download_link,
        admin_comment,
        project_id
    )
VALUES (
        'OCT66_17',
        'Project_B',
        '2023-10-16 20:47:25.152158+00',
        'https://youtube.com',
        'Admin comment2',
        (
            SELECT id
            FROM project
            where project_code = 'OCT66_17'
        )
    );
UPDATE project
SET project_history_id = (
        SELECT id
        FROM project_history
        WHERE project_code = 'OCT66_17'
            AND project_version = (
                SELECT project_version
                FROM project
                WHERE project_code = 'OCT66_17'
            )
    )
WHERE project_code = 'OCT66_17';
INSERT INTO project_history (
        project_code,
        project_name,
        created_at,
        download_link,
        admin_comment,
        project_id
    )
VALUES (
        'OCT66_21',
        'Project_C',
        '2023-10-20 22:47:25.152158+00',
        NULL,
        NULL,
        (
            SELECT id
            FROM project
            where project_code = 'OCT66_21'
        )
    );
UPDATE project
SET project_history_id = (
        SELECT id
        FROM project_history
        WHERE project_code = 'OCT66_21'
            AND project_version = (
                SELECT project_version
                FROM project
                WHERE project_code = 'OCT66_21'
            )
    )
WHERE project_code = 'OCT66_21';
INSERT INTO project_history (
        project_code,
        project_name,
        created_at,
        download_link,
        admin_comment,
        project_id
    )
VALUES (
        'NOV66_10',
        'Project_D',
        '2023-11-10 10:47:25.152158+00',
        'https://google.com',
        NULL,
        (
            SELECT id
            FROM project
            where project_code = 'NOV66_10'
        )
    );
UPDATE project
SET project_history_id = (
        SELECT id
        FROM project_history
        WHERE project_code = 'NOV66_10'
            AND project_version = (
                SELECT project_version
                FROM project
                WHERE project_code = 'NOV66_10'
            )
    )
WHERE project_code = 'NOV66_10';
INSERT INTO project_history (
        project_code,
        project_name,
        created_at,
        download_link,
        admin_comment,
        project_id
    )
VALUES (
        'NOV66_18',
        'Project_E',
        '2023-11-18 11:47:25.152158+00',
        'https://google.com',
        NULL,
        (
            SELECT id
            FROM project
            where project_code = 'NOV66_18'
        )
    );
UPDATE project
SET project_history_id = (
        SELECT id
        FROM project_history
        WHERE project_code = 'NOV66_18'
            AND project_version = (
                SELECT project_version
                FROM project
                WHERE project_code = 'NOV66_18'
            )
    )
WHERE project_code = 'NOV66_18';
-- END project_history
-- improvement
INSERT INTO improvement (
        benefit,
        experience_and_reliability,
        fund_and_output,
        project_quality,
        project_standard,
        vision_and_image
    )
VALUES (TRUE, TRUE, FALSE, FALSE, TRUE, TRUE);
-- review
INSERT INTO review (
        user_id,
        project_history_id,
        is_interested_person,
        interested_person_type,
        summary,
        improvement_id,
        comment
    )
VALUES (
        (
            SELECT id
            FROM users
            WHERE email = 'reviewer1@test.com'
        ),
        (
            SELECT id
            FROM project_history
            WHERE project_code = 'OCT66_16'
                AND project_version = 1
        ),
        TRUE,
        'connection',
        'to_be_revised',
        (
            SELECT id
            FROM improvement
            WHERE benefit = TRUE
                AND experience_and_reliability = TRUE
                AND fund_and_output = FALSE
                AND project_quality = FALSE
                AND project_standard = TRUE
                AND vision_and_image = TRUE
        ),
        'Test DB comment'
    );
-- review_details
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 1
        ),
        4
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 2
        ),
        3
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 3
        ),
        4
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 4
        ),
        5
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 5
        ),
        4
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 6
        ),
        3
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 7
        ),
        2
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 8
        ),
        1
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 9
        ),
        2
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 10
        ),
        3
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 11
        ),
        4
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 12
        ),
        5
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 13
        ),
        4
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 14
        ),
        3
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 15
        ),
        2
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 16
        ),
        1
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 17
        ),
        2
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 18
        ),
        3
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 19
        ),
        4
    );
INSERT INTO review_details (review_id, review_criteria_id, score)
VALUES (
        (
            SELECT id
            FROM review
            WHERE user_id =(
                    SELECT id
                    FROM users
                    WHERE email = 'reviewer1@test.com'
                )
                AND project_history_id =(
                    SELECT id
                    FROM project_history
                    WHERE project_code = 'OCT66_16'
                        AND project_version = 1
                )
        ),
        (
            SELECT id
            FROM review_criteria
            WHERE criteria_version = 1
                AND order_number = 20
        ),
        5
    );
-- End Review details