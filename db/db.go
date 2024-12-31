package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	postgresDB *sql.DB
	mongoDB    *mongo.Client
)

// initPostgresDB initializes PostgreSQL connection
func InitPostgresDB() error {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	postgresDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = postgresDB.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Successfully connected to PostgreSQL")
	return nil
}

// initMongoDB initializes MongoDB connection
func InitMongoDB() error {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Failed to connect to MongoDB:", err)
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Println("Failed to ping MongoDB:", err)
		return err
	}

	mongoDB = client
	log.Println("Connected to MongoDB")
	return nil
}

// GetMongoDB returns the global MongoDB client
func GetMongoDB() (*mongo.Client, error) {
	if mongoDB == nil {
		return nil, fmt.Errorf("MongoDB client is not initialized")
	}
	return mongoDB, nil
}

// GetPostgresDB returns the PostgreSQL connection
func GetPostgresDB() *sql.DB {
	return postgresDB
}
