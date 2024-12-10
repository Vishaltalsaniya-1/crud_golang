package request

// CreateUserRequest defines the structure of the incoming request to create a user.
type CreateUserRequest struct {
	Name     string   `json:"name" validate:"required,min=3"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
	Subjects []string `json:"subjects"`
}

// UpdateUserRequest defines the structure for updating user data.
type UpdateUserRequest struct {
	Name     *string  `json:"name" validate:"omitempty,min=2,max=100"` // Optional, with validation rules
	Email    *string  `json:"email" validate:"omitempty,email"`        // Optional, valid email format
	Password *string  `json:"password" validate:"omitempty,min=6"`     // Optional, with a minimum length
	Subjects []string `json:"subjects"`
}

// GetUserRequest defines the structure for filtering users with pagination
// type GetUserRequest struct {
// 	Name     string   `query:"name"`     // Filter by name (optional)
// 	Email    string   `query:"email"`    // Filter by email (optional)
// 	Subjects []string `query:"subjects"` // Filter by subjects (optional)
// }
