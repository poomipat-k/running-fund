package users

const getReviewersSQL = "SELECT id, first_name, last_name, email, user_role, created_at FROM users WHERE user_role = $1"

const getReviewerByIdSQL = "SELECT id, first_name, last_name, email FROM users WHERE id = $1 and  user_role = 'reviewer'"

const getUserByEmailSQL = "SELECT id, email, first_name, last_name, user_role, activated, activate_before, created_at FROM users WHERE email = $1"

const addUserSQL = "INSERT INTO users (email, password, first_name, last_name, user_role, activated) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"

const DeleteUserByIdSQL = "DELETE FROM users WHERE id = $1 RETURNING id;"
