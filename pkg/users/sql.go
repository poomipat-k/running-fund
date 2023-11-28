package users

const getReviewersSQL = "SELECT id, first_name, last_name, email, user_role, created_at FROM users WHERE user_role = $1"

const getReviewerByIdSQL = "SELECT id, first_name, last_name, email FROM users WHERE id = $1 and  user_role = 'reviewer'"
