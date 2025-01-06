package model

import (
	"time"
)

type User struct {
	Id        string     `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid()" bson:"_id,omitempty"`
	Name      string     `json:"name" bson:"name"`
	Email     string     `json:"email" bson:"email"`
	Subjects  []string   `json:"subjects" gorm:"type:text[]" bson:"subjects"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;default:current_timestamp" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;default:current_timestamp" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at" bson:"deleted_at,omitempty"`
}
