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

func InitPostgresDB() error {
	
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
	err = ExistsTable(postgresDB)
	if err != nil {
		return fmt.Errorf("failed to ensure table exists: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL")
	return nil
}


func InitMongoDB() error {
	
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

func GetMongoDB() (*mongo.Client, error) {
	if mongoDB == nil {
		return nil, fmt.Errorf("MongoDB client is not initialized")
	}
	return mongoDB, nil
}


func GetPostgresDB() *sql.DB {
	return postgresDB
}

func ExistsTable(db *sql.DB) error {

	createTableSQL := `
	 CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100) UNIQUE NOT NULL,
        subjects TEXT[],
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP
    );
    `
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error user table%v\n", err)
		return fmt.Errorf("failed to create table:%v", err)

	}
	return nil
}
