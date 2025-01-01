package manager

import (
	"fitness-api/model"
	"fitness-api/request"
	"fitness-api/service"
	"fmt"
	"log"
)

type UserManager struct {
}

func NewUserManager() *UserManager {
	// Initialize and return a new instance of UserManager
	return &UserManager{}
}
func derefString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (um *UserManager) CreateUser(req request.CreateUserRequest) (model.User, error) {

	// Validate the incoming request data
	if req.Name == "" || req.Email == "" {
		return model.User{}, fmt.Errorf("manager: name and email are required fields")
	}

	// Map the request fields to the User model
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Subjects: req.Subjects,
	}

	// Call the service layer to create the user in the database
	createdUser, err := service.CreateUser(user)
	if err != nil {
		// Handle specific errors if needed
		log.Printf("manager: failed to create user: %v\n", err)
		return model.User{}, fmt.Errorf("manager: failed to create user: %w", err)
	}

	return createdUser, nil
}

func (um *UserManager) UpdateUser(id string, req request.UpdateUserRequest) (model.User, error) {

	// log.Printf("Updating user: ID: %d, Name: %v, Email: %v, Subjects: %v", id, req.Name, req.Email, req.Subjects)

	if req.Name == nil || req.Email == nil || *req.Name == "" || *req.Email == "" {
		return model.User{}, fmt.Errorf("manager: name and email cannot be empty")
	}

	user := model.User{
		Name:     derefString(req.Name),
		Email:    derefString(req.Email),
		Subjects: req.Subjects,
	}

	// Pass user, ID, and flag to the service layer
	updatedUser, err := service.UpdateUser(user, id)
	if err != nil {
		return model.User{}, fmt.Errorf("manager: failed to update user: %v", err)
	}

	return updatedUser, nil
}

// DeleteUser handles the deletion of a user
func (um *UserManager) DeleteUser(id string) error {
	err := service.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("manager: failed to delete user: %v", err)
	}

	return nil
}

// GetAllUsers retrieves all users with pagination and filtering options
func (um *UserManager) GetAllUsers(pageSize int, pageNo int, subject string, order string, orderby string) ([]model.User, int, int, error) {

	// Call service with pagination arguments
	users, lastPage, totalDocuments, err := service.GetAllUsers(pageSize, pageNo, subject, order, orderby)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("manager: failed to fetch users: %v", err)
	}
	return users, lastPage, totalDocuments, nil
}

// GetUserByID retrieves a user by ID

func (um *UserManager) GetUserByID(id string) (model.User, error) {
	user, err := service.GetUserByID(id) // Updated method call
	if err != nil {
		return model.User{}, fmt.Errorf("manager: failed to fetch user: %v", err)
	}
	return user, nil
}