package manager

import (
	"fitness-api/model"
	"fitness-api/request"
	"fitness-api/service"
	"fmt"
	"log"
	"time"
	//"github.com/jinzhu/now"
)

type UserManager struct {
}

func NewUserManager() *UserManager {

	return &UserManager{}
}

// func derefString(ptr *string) string {
// 	if ptr == nil {
// 		return ""
// 	}
// 	return *ptr
// }

func (um *UserManager) CreateUser(req request.UserRequest) (model.User, error) {

	if req.CreatedAt == nil {
		now := time.Now()
		req.CreatedAt = &now
	}

	user := model.User{
		Name:      req.Name,
		Email:     req.Email,
		Subjects:  req.Subjects,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.CreatedAt,
		DeletedAt: nil,
	}
	createdUser, err := service.CreateUser(user)
	if err != nil {
		log.Println("Failed to create user:", err)
		return model.User{}, fmt.Errorf("error unable to create user please try again: %w", err)
	}
	return createdUser, nil
}

func (um *UserManager) UpdateUser(id string, req request.UserRequest) (model.User, error) {

	user := model.User{
		Name:      req.Name,
		Email:     req.Email,
		Subjects:  req.Subjects,
		CreatedAt: req.CreatedAt,
		UpdatedAt: &time.Time{},
		DeletedAt: nil,
	}

	updatedUser, err := service.UpdateUser(user, id)
	if err != nil {

		return model.User{}, err
	}

	return updatedUser, nil
}

func (um *UserManager) DeleteUser(id string) error {
	err := service.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}

func (um *UserManager) GetAllUsers(pageSize int, pageNo int, subject string, order string, orderby string) ([]model.User, int, int, error) {

	users, lastPage, totalDocuments, err := service.GetAllUsers(pageSize, pageNo, subject, order, orderby)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch users: ")
	}
	return users, lastPage, totalDocuments, nil
}

func (um *UserManager) GetUserByID(id string) (model.User, error) {
	user, err := service.GetUserByID(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
