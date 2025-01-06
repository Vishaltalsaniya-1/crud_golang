package controller

import (
	manager "fitness-api/managers"
	"fitness-api/request"
	"fitness-api/response"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	manager *manager.UserManager
}

func NewUserController(mn *manager.UserManager) *UserController {
	return &UserController{manager: mn}
}

func (uc *UserController) CreateUser(c echo.Context) error {

	var req request.UserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	log.Println("REQ-------->", req)
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	createdUser, err := uc.manager.CreateUser(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, createdUser)
}

func (uc *UserController) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	var req request.UserRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	updatedUser, err := uc.manager.UpdateUser(id, req)
	if err != nil {
		log.Printf("Error updating user with ID %s: %v", id, err)

		if strings.Contains(err.Error(), "no user found with the given ID") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("No user found with the given ID: %s. Please check the ID and try again.", id),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "An unexpected error occurred while updating the user. Please try again later.",
		})
	}
	return c.JSON(http.StatusOK, updatedUser)
}

func (uc *UserController) DeleteUser(c echo.Context) error {

	id := c.Param("id")

	err := uc.manager.DeleteUser(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user found with id %s", id) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (uc *UserController) GetAllUsers(c echo.Context) error {

	pageSize := c.QueryParam("per_page")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if pageSize == "-1" {
		pageSizeInt = -1
	} else if err != nil || pageSizeInt <= 0 {
		pageSizeInt = 10
	}

	pageNo := c.QueryParam("page_no")
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil || pageNoInt <= 0 {
		pageNoInt = 1
	}

	order := c.QueryParam("order")
	orderby := c.QueryParam("orderby")
	subject := c.QueryParam("subject")

	users, lastPage, totalDocuments, err := uc.manager.GetAllUsers(pageSizeInt, pageNoInt, subject, order, orderby)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"page_no":         pageNoInt,
		"per_page":        pageSizeInt,
		"last_page":       lastPage,
		"total_documents": totalDocuments,
		"users":           users,
	})
}
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
func (uc *UserController) GetUserByID(c echo.Context) error {
	id := c.Param("id")

	user, err := uc.manager.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	//return c.JSON(http.StatusInternalServerError, map[string]string{"error":"Internal server error"})

	return c.JSON(http.StatusOK, response.UserResponse{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Subjects:  user.Subjects,
		CreatedAt: formatTime(user.CreatedAt),
		UpdatedAt: formatTime(user.UpdatedAt),
		DeletedAt: formatTime(user.DeletedAt),
	})

}
