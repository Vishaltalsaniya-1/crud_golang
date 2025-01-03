package main

import (
	// "fitness-api/config"
	"fitness-api/config"
	controller "fitness-api/controller"
	"fitness-api/db"
	manager "fitness-api/managers"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {

	// Initialize the database
	flagConfig, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading flag config: %v", err)
	}

	if flagConfig.FlagValue == "TRUE" {
		
		fmt.Println("MongoDB URL:", flagConfig.FlagValue)
		if err := db.InitMongoDB(); err != nil {
			log.Fatalf("Failed to initialize MongoDB: %v", err)
		}
		fmt.Println("MongoDB Initialized")
	} else {

		if err := db.InitPostgresDB(); err != nil {
			log.Fatalf("Failed to initialize PostgreSQL: %v", err)
		}
		fmt.Println("PostgreSQL Initialized") 
	}

	// instanc of manager
	userManager := manager.NewUserManager()
	userController := controller.NewUserController(userManager)
	// Create a new Echo instance
	e := echo.New()

	// Define routes and their handlers
	e.POST("/users", userController.CreateUser)
	e.GET("/users", userController.GetAllUsers)
	e.PUT("/users/:id", userController.UpdateUser)
	e.GET("/users/:id", userController.GetUserByID)
	e.DELETE("/users/:id", userController.DeleteUser)

	// Start the server on port 8081
	e.Logger.Fatal(e.Start(":8081"))
}
