package controller

import (
	manager "fitness-api/managers"
	"fitness-api/request"
	"fitness-api/response"
	"log"
	"net/http"
	"strconv"

	// "fitness-api/config"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	manager *manager.UserManager
}

func NewUserController(mn *manager.UserManager) *UserController {
	// manager:=manager.NewUserManager()
	return &UserController{manager: mn}
}
func (uc *UserController) CreateUser(c echo.Context) error {

	var req request.CreateUserRequest

	// Bind request body to CreateUserRequest
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

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Successfully created the user, return HTTP 201 Created
	return c.JSON(http.StatusCreated, createdUser)
}

func (uc *UserController) UpdateUser(c echo.Context) error {
	// Parse user ID from the request URL
	idParam := c.Param("id")
	log.Println("ID from URL:", idParam)
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		log.Println("Invalid ID:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	log.Println("Converted ID:", id)

	var req request.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		log.Println("Request body error:", err) // Log error details
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Flag == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Flag is required and should be 'true' for MongoDB or 'false' for PostgreSQL"})
	}

	// Call manager to perform the update
	updatedUser, err := uc.manager.UpdateUser(id, req, req.Flag)
	if err != nil {
		log.Println("Update user error:", err) // Log error details
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return the updated user
	return c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles user deletion
func (uc *UserController) DeleteUser(c echo.Context) error {
	// Get the flag value from query parameters
	flag := c.QueryParam("flag")

	// Validate the flag value
	if flag != "true" && flag != "false" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Get the user ID from the URL parameters
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	// Call the manager's DeleteUser method
	err = uc.manager.DeleteUser(id, flag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// GetAllUsers retrieves all users with pagination and filtering options
func (uc *UserController) GetAllUsers(c echo.Context) error {

	pageSize := c.QueryParam("per_page")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt <= 0 {
		pageSizeInt = 10 // Default page size
	}

	pageNo := c.QueryParam("page_no")
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil || pageNoInt <= 0 {
		pageNoInt = 1 // Default page number
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

// GetUserByID retrieves a user by ID
func (uc *UserController) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	flag := c.QueryParam("flag")
	if flag != "true" && flag != "false" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request "})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	user, err := uc.manager.GetUserByID(id, flag)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response.UserResponse{
		ID:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Subjects: user.Subjects,
	})
}
