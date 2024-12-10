package controller

import (
	manager "fitness-api/cmd/managers"
	"fitness-api/cmd/request"
	"fitness-api/cmd/response"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CreateUser handles user creation
func CreateUser(c echo.Context) error {
	flag := c.QueryParam("flag")

	if flag != "true" && flag != "false"{
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request "})
	}
	var req request.CreateUserRequest

	// Bind request body to CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Validate the request data
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// If validation fails, return the error message
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Proceed with creating user if validation passes
	createdUser, err := manager.CreateUser(req, flag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, createdUser)
}

// UpdateUser handles user updates
func UpdateUser(c echo.Context) error {
	
	flag := c.QueryParam("flag")
	if flag != "true" && flag !="false"{
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request "})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	var req request.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Validate the request data
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// If validation fails, return the error message
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	updatedUser, err := manager.UpdateUser(id, req, flag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles user deletion
func DeleteUser(c echo.Context) error {
	flag := c.QueryParam("flag")
	
	if flag != "true" && flag !="false"{
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request "})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	err = manager.DeleteUser(id, flag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// GetAllUsers retrieves all users with pagination and filtering options
func GetAllUsers(c echo.Context) error {
	flag := c.QueryParam("flag")

	if flag != "true" && flag !="false"{
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request "})
	}
	// Parse `page_size` query parameter
	pageSize := c.QueryParam("per_page")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt <= 0 {
		pageSizeInt = 10 // Default page size
	}

	// Parse `page_no` query parameter
	pageNo := c.QueryParam("page_no")
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil || pageNoInt <= 0 {
		pageNoInt = 1 // Default page number
	}

	order := c.QueryParam("order")
	orderby := c.QueryParam("orderby")

	subject := c.QueryParam("subject")

	// Call manager's `GetAllUsers` function with pagination arguments
	users, lastPage, totalDocuments, err := manager.GetAllUsers(pageSizeInt, pageNoInt, subject, order, orderby, flag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return the paginated response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"page_no":         pageNoInt,
		"per_page":        pageSizeInt,
		"last_page":       lastPage,
		"total_documents": totalDocuments,
		"users":           users,
	})
}

// GetUserByID retrieves a user by ID
func GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	flag := c.QueryParam("flag")
	if flag != "true" && flag !="false"{
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request "})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	user, err := manager.GetUserByID(id, flag)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response.UserResponse{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Subjects:  user.Subjects,
	})
}
