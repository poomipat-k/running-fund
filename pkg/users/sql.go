package users

const getUserByIdSQL = "SELECT id, first_name, last_name, email, user_role, activated FROM users WHERE id = $1"

const getUserByEmailSQL = "SELECT id, email, password, first_name, last_name, user_role, activated, activate_before, created_at FROM users WHERE email = $1"

const addUserSQL = "INSERT INTO users (email, password, first_name, last_name, user_role, activated, activate_code) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"

const DeleteUserByIdSQL = "DELETE FROM users WHERE id = $1 RETURNING id;"

const activateEmailSQL = "UPDATE users SET activated = true, activate_code = NULL WHERE activate_code = $1 AND activated = false AND activate_before >= now();"
