package main

import (
	"fitness-api/cmd/db"
	"fitness-api/cmd/router"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize the database connection
	// Replace with your PostgreSQL connection string
	db.InitDB()

	// Create a new Echo instance
	e := echo.New()

	// Register routes
	router.RegisterRoutes(e)

	// Start the server
	e.Logger.Fatal(e.Start(":8081"))
}
