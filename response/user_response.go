package response

type UserResponse struct {
	ID       string   `json:"id" bson:"id"`
	Name     string   `json:"name" bson:"name"`
	Email    string   `json:"email" bson:"email"`
	Subjects []string `json:"subjects" bson:"subjects"`

	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
	DeletedAt string `json:"deleted_at" bson:"deleted_at"`
}
