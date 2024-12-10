package service

import (
	"fitness-api/cmd/db"
	"fitness-api/cmd/model"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// CreateUser inserts a new user into the database
// CreateUser function (with error handling)
func CreateUser(user model.User) (model.User, error) {
	db := db.GetDB()

	// Check if user already exists by email or name (filtering logic)
	var existingUser model.User
	err := db.QueryRow("SELECT id FROM users WHERE email = $1 OR name = $2 LIMIT 1", user.Email, user.Name).Scan(&existingUser.Id)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return model.User{}, fmt.Errorf("service: failed to check user existence: %v", err)
	}

	if existingUser.Id != 0 {
		return model.User{}, fmt.Errorf("service: user with this email or name already exists")
	}

	// Insert the user into the database
	sqlStatement := `
        INSERT INTO users (name, email, password, created_at, updated_at, subjects)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, name, email, created_at, updated_at, subjects`

	var createdUser model.User
	var subjects []string // Temporary slice to scan subjects

	err = db.QueryRow(
		sqlStatement,
		user.Name,
		user.Email,
		user.Password,
		time.Now(),
		time.Now(),
		pq.Array(user.Subjects), // Use pq.Array for the slice
	).Scan(
		&createdUser.Id,
		&createdUser.Name,
		&createdUser.Email,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
		pq.Array(&subjects), // Use pq.Array to decode array
	)
	createdUser.Subjects = subjects
	if err != nil {
		return model.User{}, fmt.Errorf("service: failed to create user: %v", err)
	}
	return createdUser, nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(user model.User, id int) (model.User, error) {
	db := db.GetDB()

	// Correct the SQL statement
	sqlStatement := `
        UPDATE users 
        SET name = $1, email = $2, password = $3, updated_at = CURRENT_TIMESTAMP, subjects = $4
        WHERE id = $5
        RETURNING id, name, email, updated_at, subjects`

	var updatedUser model.User
	var subjects pq.StringArray // This allows handling of NULL values correctly

	// Pass the subjects array correctly using pq.Array
	err := db.QueryRow(sqlStatement, user.Name, user.Email, user.Password, pq.Array(user.Subjects), id).Scan(
		&updatedUser.Id, &updatedUser.Name, &updatedUser.Email, &updatedUser.UpdatedAt, &subjects)

	if err != nil {
		return model.User{}, fmt.Errorf("service: failed to update user: %v", err)
	}

	if subjects == nil {
		updatedUser.Subjects = []string{}
	} else {
		updatedUser.Subjects = subjects
	}
	return updatedUser, nil
}

// DeleteUser deletes a user from the database by ID
func DeleteUser(id int) error {
	db := db.GetDB()

	sqlStatement := `DELETE FROM users WHERE id = $1`

	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("service: failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("service: no user found with id %d", id)
	}
	return nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(pageSize int, pageNo int, subject string, order string, orderby string) ([]model.User, int, int, error) {
	db := db.GetDB()

	// Calculate offset for pagination
	offset := (pageNo - 1) * pageSize

	validColumns := map[string]bool{"id": true, "name": true, "email": true, "created_at": true, "updated_at": true}
	if !validColumns[orderby] {
		orderby = "created_at" // Default column
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC" // Default sorting order
	}
	// Query to get the latest user first, followed by paginated users
	sqlStatement := fmt.Sprintf(`
	SELECT id, name, email, created_at, updated_at, subjects
	FROM users
	WHERE $1 = ANY(subjects) OR $1 = ''
	ORDER BY %s %s
	LIMIT $2 OFFSET $3`, orderby, order)

	rows, err := db.Query(sqlStatement, subject, pageSize, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("service: failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var subjects []string

		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, pq.Array(&subjects))
		if err != nil {
			return nil, 0, 0, fmt.Errorf("service: failed to scan user: %v", err)
		}

		user.Subjects = subjects
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("service: failed to fetch users: %v", err)
	}

	// If no users were fetched, return an empty list
	if len(users) == 0 {
		return nil, 0, 0, fmt.Errorf("service: no users found")
	}

	// Count the total number of users that match the filter
	var totalDocuments int
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE $1 = ANY(subjects) OR $1 = ''`
	err = db.QueryRow(countQuery, subject).Scan(&totalDocuments)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("service: failed to count total users: %v", err)
	}

	// Calculate last page based on total documents and pageSize
	lastPage := (totalDocuments + pageSize - 1) / pageSize

	return users, lastPage, totalDocuments, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id int) (model.User, error) {
	db := db.GetDB()

	sqlStatement := `SELECT id, name, email, created_at, updated_at, subjects FROM users WHERE id = $1`

	var user model.User
	var subjects []string
	err := db.QueryRow(sqlStatement, id).Scan(
		&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, pq.Array(&subjects))

	if err != nil {
		return model.User{}, fmt.Errorf("service: user not found: %v", err)
	}
	user.Subjects = subjects
	return user, nil
}
