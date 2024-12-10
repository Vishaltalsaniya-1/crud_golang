package manager

import (
	"fitness-api/cmd/model"
	"fitness-api/cmd/request"
	"fitness-api/cmd/service"
	"fmt"
	"log"
)

// Helper function to safely dereference string pointers for UpdateUserRequest
func derefString(ptr *string) string {
	if ptr == nil {
		return "" // Default value for nil pointers
	}
	return *ptr
}

// CreateUser delegates the user creation to the service layer
func CreateUser(req request.CreateUserRequest) (model.User, error) {
	// Directly assign fields from CreateUserRequest (no derefString needed)
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // You can hash the password here
		Subjects: req.Subjects,
	}

	createdUser, err := service.CreateUser(user)
	if err != nil {
		return model.User{}, fmt.Errorf("manager: failed to create user: %v", err)
	}
	return createdUser, nil
}

// UpdateUser delegates the user update to the service layer
func UpdateUser(id int, req request.UpdateUserRequest) (model.User, error) {
	// Use derefString to safely handle nil pointers in UpdateUserRequest
	user := model.User{
		Name:     derefString(req.Name),
		Email:    derefString(req.Email),
		Password: derefString(req.Password), // You can hash the password here
		Subjects: req.Subjects,
	}

	updatedUser, err := service.UpdateUser(user, id)
	if err != nil {
		return model.User{}, fmt.Errorf("manager: failed to update user: %v", err)
	}
	return updatedUser, nil
}

// Other unchanged functions
func DeleteUser(id int) error {
	err := service.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("manager: failed to delete user: %v", err)
	}
	return nil
}

func GetAllUsers(pageSize int, pageNo int, subject string, order string, orderby string) ([]model.User,int,int, error) {
	// Log for debugging
	log.Println("PageSize ----->", pageSize)
	log.Println("PageNo-------->", pageNo)
	log.Println("Subject-------->", subject)

	// Call service with pagination arguments0
	users,lastPage, totalDocuments, err := service.GetAllUsers(pageSize, pageNo, subject,order, orderby)
	if err != nil {
		return nil,0,0, fmt.Errorf("manager: failed to fetch users: %v", err)
	}
	return users,lastPage, totalDocuments,  nil
}

func GetUserByID(id int) (model.User, error) {
	user, err := service.GetUserByID(id)
	if err != nil {
		return model.User{}, fmt.Errorf("manager: failed to fetch user: %v", err)
	}
	return user, nil
}

// func GetUsers(req request.GetUserRequest) ([]model.User, error) {
// 	dbConn := db.GetDB()
// 	query := "SELECT * FROM users WHERE 1=1"
// 	args := []interface{}{}
// 	argIndex := 1

// 	// Apply filters for name, email, subjects, etc.
// 	if req.Name != "" {
// 		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
// 		args = append(args, "%"+req.Name+"%")
// 		argIndex++
// 	}
// 	if req.Email != "" {
// 		query += fmt.Sprintf(" AND email ILIKE $%d", argIndex)
// 		args = append(args, "%"+req.Email+"%")
// 		argIndex++
// 	}
// 	if len(req.Subjects) > 0 {
// 		query += fmt.Sprintf(" AND subjects && $%d", argIndex) // Use PostgreSQL array overlap operator
// 		args = append(args, req.Subjects)
// 		argIndex++
// 	}

// 	// Pagination (Limit and Offset)
// 	offset := (req.Page - 1) * req.Limit
// 	query += fmt.Sprintf(" LIMIT %d OFFSET %d", req.Limit, offset)

// 	// Execute the query
// 	rows, err := dbConn.Query(query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var users []model.User
// 	for rows.Next() {
// 		var user model.User
// 		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Subjects, &user.CreatedAt, &user.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		users = append(users, user)
// 	}

// 	return users, nil
// }
