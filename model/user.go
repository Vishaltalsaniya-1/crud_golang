package model

type User struct {
	Id    string `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid()" bson:"_id,omitempty"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	//	Flag     string    `json:"flag"`
	Subjects []string `json:"subjects" gorm:"type:text[]" bson:"subjects"`
}
