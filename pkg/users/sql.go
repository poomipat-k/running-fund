package users

const getUserByIdSQL = "SELECT id, first_name, last_name, email, user_role, activated FROM users WHERE id = $1"

const getUserByEmailSQL = "SELECT id, email, password, first_name, last_name, user_role, activated, activate_before, created_at FROM users WHERE email = $1"

const addUserSQL = "INSERT INTO users (email, password, first_name, last_name, user_role, activated) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"

const DeleteUserByIdSQL = "DELETE FROM users WHERE id = $1 RETURNING id;"
