package request

// CreateUserRequest defines the structure of the incoming request to create a user.
type CreateUserRequest struct {
	Name     string   `json:"fname" bson:"name" validate:"required,min=3"`
	Email    string   `json:"email" bson:"email" validate:"required,email"`
	Flag     string   `json:"flag"`
	Subjects []string `json:"subjects" bson:"subjects"`
}

// UpdateUserRequest defines the structure for updating user data.
type UpdateUserRequest struct {
	Name     *string  `json:"name" `    // Optional field
	Email    *string  `json:"email"`    // Optional field
	Subjects []string `json:"subjects"` // Optional field
	Flag     string   `json:"flag"`     // "true" for MongoDB, "false" for PostgreSQL
}
