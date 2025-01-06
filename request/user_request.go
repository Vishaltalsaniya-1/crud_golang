package request

import "time"

type UserRequest struct {
	Name      string     `json:"name" bson:"name" `
	Email     string     `json:"email" bson:"email" validate:"required,email"`
	Subjects  []string   `json:"subjects" bson:"subjects"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;default:current_timestamp" bson:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;default:current_timestamp" bson:"updated_at,omitempty"`
}
