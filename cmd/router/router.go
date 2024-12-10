package router

import (
	"fitness-api/cmd/controller"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes sets up all application routes
func RegisterRoutes(e *echo.Echo) {
	// User routes
	e.POST("/users", controller.CreateUser)
	e.GET("/users", controller.GetAllUsers)
	e.PUT("/users/:id", controller.UpdateUser)
	e.GET("/users/:id", controller.GetUserByID)
	e.DELETE("/users/:id", controller.DeleteUser)
}
