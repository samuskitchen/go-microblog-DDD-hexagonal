package persistence

const(

	// selectAllUser is a query that selects all rows in the user table
	selectAllUser = "SELECT id, first_name, last_name, username, email, picture, created_at, updated_at FROM users;"

	// selectUserById is a query that selects a row from the users table based off of the given id.
	selectUserById = "SELECT id, first_name, last_name, username, email, picture, created_at, updated_at FROM users WHERE id = $1;"

	// selectUSerByUsername is a query that selects a row from the users table based off of the given username
	selectUSerByUsername = "SELECT id, first_name, last_name, username, email, picture, password, created_at, updated_at FROM users WHERE username = $1;"

	// insertUser is a query that inserts a new row in the user table using the values
	// given in order for first_name, last_name, username, email, picture, password, created_at and updated_at.
	insertUser = "INSERT INTO users (first_name, last_name, username, email, picture, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"

	// updateUser is a query that updates a row in the users table based off of id.
	// The values able to be updated are first_name, last_name, email, picture and updated_at.
	updateUser = "UPDATE users SET first_name=$1, last_name=$2, email=$3, picture=$4, updated_at=$5 WHERE id=$6;"

	// deleteUser is a query that deletes a row in the users table given a id.
	deleteUser = "DELETE FROM users WHERE id=$1;"
)
