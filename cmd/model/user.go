package model

// User represents a user in the system
type User struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Subjects []string `json:"subjects"` // Add subjects as a slice of strings
}
