package response

import (
	"github.com/go-playground/validator/v10" // Ensure this import exists
)

var validate *validator.Validate

// Initialize validate in init function
func init() {
	validate = validator.New() // This initializes the validator
}

// UserResponse defines the structure of the outgoing response for a user.
type UserResponse struct {
	ID       string   `json:"id" bson:"id"`
	Name     string   `json:"name" bson:"name" validate:"required"`
	Email    string   `json:"email" bson:"email" validate:"required,email"`
	Subjects []string `json:"subjects" bson:"subjects"`
}

// MessageResponse is a generic response structure for status messages.
type MessageResponse struct {
	Message string `json:"message" bson:"message" validate:"required"`
}

// Validate method for UserResponse
func (u *UserResponse) Validate() error {
	return validate.Struct(u) // This validates the struct
}

// Validate method for MessageResponse
func (m *MessageResponse) Validate() error {
	return validate.Struct(m) // This validates the struct
}
