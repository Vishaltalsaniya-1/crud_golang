package service

import (
	"context"
	"database/sql"
	"fitness-api/config"
	"fitness-api/db"
	"fitness-api/model"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateUser inserts a new user into both PostgreSQL and MongoDB
func CreateUser(user model.User) (model.User, error) {
	flagConfig, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading flag config: %v", err)
	}
	if flagConfig.FlagValue == "TRUE" {
		// MongoDB logic
		log.Println(" MongoDB flag is true, processing user creation in MongoDB")

		mongoClient, err := db.GetMongoDB()
		if err != nil {
			log.Printf("MongoDB initialization error: %v\n", err)
			return model.User{}, fmt.Errorf(" MongoDB is not initialized: %v", err)
		}

		mongoCollection := mongoClient.Database("fitness").Collection("users")
		var id = uuid.New().String()
		// Generate a new unique _id
		mongoUser := bson.M{
			"_id":      id,
			"name":     user.Name,
			"email":    user.Email,
			"subjects": user.Subjects,
		}

		log.Printf(" Inserting user into MongoDB: %+v\n", mongoUser)

		// Insert into MongoDB
		_, err = mongoCollection.InsertOne(context.Background(), mongoUser)
		if err != nil {
			log.Printf("Failed to create user in MongoDB: %v", err)
			return model.User{}, fmt.Errorf("database error")
		}

		log.Println(" User successfully created in MongoDB")
		user.Id = id
		return user, nil
	}

	log.Println(" PostgreSQL flag is false, processing user creation in PostgreSQL")
	postgresDB := db.GetPostgresDB()

	// Check if the user already exists by email or name
	var existingUser model.User
	err = postgresDB.QueryRow("SELECT id FROM users WHERE email = $1 OR name = $2 LIMIT 1", user.Email, user.Name).Scan(&existingUser.Id)
	if err != nil {
		//log.Printf(" Error checking user existence in PostgreSQL: %v\n", err)
		return model.User{}, fmt.Errorf(" failed to check user existence: %v", err)
	}

	// Insert the new user into PostgreSQL
	sqlStatement := `
        INSERT INTO users (name, email, subjects)
        VALUES ($1, $2, $3)
        RETURNING id, name, email, subjects`

	var createdUser model.User
	var subjects []string

	err = postgresDB.QueryRow(
		sqlStatement,
		user.Name,
		user.Email,
		pq.Array(user.Subjects),
	).Scan(
		&createdUser.Id,
		&createdUser.Name,
		&createdUser.Email,
		pq.Array(&subjects),
	)
	if err != nil {
		log.Printf("Error inserting user into PostgreSQL: %v", err)

		return model.User{}, fmt.Errorf("database error")
	}
	// if err != nil {
	// 	log.Printf(" Error inserting user into PostgreSQL: %v\n", err)
	// 	return model.User{}, fmt.Errorf(" failed to create user in PostgreSQL: %v", err)
	// }
	createdUser.Subjects = subjects

	log.Println(" User successfully created in PostgreSQL")
	return createdUser, nil
}

func UpdateUser(user model.User, id string) (model.User, error) {
	flagConfig, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading flag config: %v", err)
	}

	if flagConfig.FlagValue == "TRUE" {
		// MongoDB logic
		log.Println(" MongoDB flag is true, processing user creation in MongoDB")

		mongoClient, err := db.GetMongoDB()
		if err != nil {
			log.Printf(" MongoDB initialization error: %v\n", err)
			return model.User{}, fmt.Errorf("MongoDB is not initialized: ")
		}

		mongoCollection := mongoClient.Database("fitness").Collection("users")
		filter := bson.D{{Key: "_id", Value: id}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: user.Name},
				{Key: "email", Value: user.Email},
				{Key: "subjects", Value: user.Subjects},
			}},
		}

		result, err := mongoCollection.UpdateOne(
			context.Background(),
			filter,
			update,
			&options.UpdateOptions{Upsert: &[]bool{false}[0]},
		)
		if err != nil {
			log.Printf(" MongoDB update error: %v\n", err)
			return model.User{}, fmt.Errorf(" failed to update user in MongoDB: %v", err)
		}
		if result.MatchedCount == 0 {
			return model.User{}, fmt.Errorf("no user found with the given ID: %s", id)
		}
		log.Printf(" Update result: %+v\n", result)

		var updatedUser model.User
		err = mongoCollection.FindOne(context.Background(), filter).Decode(&updatedUser)
		if err != nil {
			log.Printf(" Failed to fetch updated user: %v\n", err)
			return model.User{}, fmt.Errorf(" failed to fetch updated user from MongoDB ")
		}

		log.Println(" User successfully updated in MongoDB")
		return updatedUser, nil
	}

	pgDB := db.GetPostgresDB()

	sqlStatement := `
            UPDATE users 
            SET name = $1, email = $2, subjects = $3
            WHERE id = $4
            RETURNING id, name, email, subjects`
	log.Printf("Updating user: ID: %s, Name: %s, Email: %s, Subjects: %v", id, user.Name, user.Email, user.Subjects)

	var updatedUser model.User
	var subjects pq.StringArray

	err = pgDB.QueryRow(sqlStatement, user.Name, user.Email, pq.Array(user.Subjects), id).Scan(
		&updatedUser.Id, &updatedUser.Name, &updatedUser.Email, &subjects)

	if err != nil {
		return model.User{}, err
	}
	// if err != nil {
	// 	return model.User{}, fmt.Errorf(" failed to update user in PostgreSQL: %v", err)
	// }

	updatedUser.Subjects = subjects
	return updatedUser, nil

}

// DeleteUser deletes a user from PostgreSQL or MongoDB based on the flag
func DeleteUser(id string) error {
	flagConfig, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading flag config: %v", err)
	}
	if flagConfig.FlagValue == "TRUE" {
		log.Println(" MongoDB flag is true, processing user creation in MongoDB")

		mongoClient, err := db.GetMongoDB()

		if err != nil {
			log.Printf(" MongoDB initialization error: %v\n", err)
			return fmt.Errorf(" MongoDB is not initialized: %v", err)
		}

		mongoCollection := mongoClient.Database("fitness").Collection("users")
		filter := bson.D{{Key: "_id", Value: id}}
		result, err := mongoCollection.DeleteOne(context.Background(), filter)
		if err != nil {
			log.Printf(" MongoDB deletion error: %v\n", err)
			return fmt.Errorf(" failed to delete user from MongoDB")
		}

		if result.DeletedCount == 0 {
			return fmt.Errorf(" no user found with id %s", id)
		}

		log.Println(" User successfully deleted from MongoDB")
		return nil
	}

	db := db.GetPostgresDB()
	sqlStatement := `DELETE FROM users WHERE id = $1`

	// Execute the DELETE statement
	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf(" failed to delete user: %v", err)
	}

	// Check if any rows were affected (i.e., if the user was found and deleted)
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf(" no user found with id %s", id)
	}

	log.Println(" User successfully deleted from PostgreSQL")
	return nil
}

func GetAllUsers(pageSize int, pageNo int, subject string, order string, orderby string) ([]model.User, int, int, error) {
	flagConfig, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading flag config: %v", err)
	}
	if flagConfig.FlagValue == "TRUE" {
		log.Println(" MongoDB flag is true, fetching users from in MongoDB")
		mongoClient, err := db.GetMongoDB()
		if err != nil {
			log.Printf(" MongoDB initialization error: %v\n", err)
			return nil, 0, 0, fmt.Errorf(" MongoDB is not initialized: %v", err)
		}
		mongoCollection := mongoClient.Database("fitness").Collection("users")

		filter := bson.M{}
		if subject != "" {
			filter["subjects"] = subject
		}

		totalDocuments, err := mongoCollection.CountDocuments(context.Background(), filter)
		if err != nil {
			log.Printf(" MongoDB count error: %v\n", err)
			return nil, 0, 0, fmt.Errorf("failed to count users: %v", err)
		}

		if pageSize == -1 {
			pageSize = int(totalDocuments)
		}

		skip := (pageNo - 1) * pageSize

		sortOrder := 1
		if order == "DESC" {
			sortOrder = -1
		}
		sortField := orderby
		if sortField == "" {
			sortField = "id"
		}

		var opts *options.FindOptions
		if pageSize == -1 {
			opts = options.Find().
				SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
				SetSkip(int64(skip))
		} else {
			opts = options.Find().
				SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
				SetSkip(int64(skip)).
				SetLimit(int64(pageSize))
		}

		cursor, err := mongoCollection.Find(context.Background(), filter, opts)
		if err != nil {
			log.Printf(" MongoDB find error: %v\n", err)
			return nil, 0, 0, fmt.Errorf("failed to fetch users: %v", err)
		}
		defer cursor.Close(context.Background())

		// Decode results
		var users []model.User
		for cursor.Next(context.Background()) {
			var user model.User
			if err := cursor.Decode(&user); err != nil {
				log.Printf(" MongoDB decode error: %v\n", err)
				return nil, 0, 0, fmt.Errorf("failed to decode user data: %v", err)
			}
			users = append(users, user)
		}

		if err := cursor.Err(); err != nil {
			log.Printf(" MongoDB cursor iteration error: %v\n", err)
			return nil, 0, 0, fmt.Errorf("error iterating MongoDB cursor: %v", err)
		}

		lastPage := (int(totalDocuments) + pageSize - 1) / pageSize

		return users, lastPage, int(totalDocuments), nil
	}

	db := db.GetPostgresDB()

	// Validate `orderby` and `order` fields to prevent SQL injection
	validColumns := map[string]bool{"id": true, "name": true, "email": true}
	if !validColumns[orderby] {
		orderby = "id" // Default column
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC" // Default sorting order
	}

	var sqlStatement string
	var rows *sql.Rows

	if pageSize == -1 {
		sqlStatement = fmt.Sprintf(`
			SELECT id, name, email, subjects
			FROM users
			WHERE $1 = ANY(subjects) OR $1 = ''
			ORDER BY %s %s`, orderby, order)

		rows, err = db.Query(sqlStatement, subject)
	} else {
		offset := (pageNo - 1) * pageSize
		sqlStatement = fmt.Sprintf(`
			SELECT id, name, email, subjects
			FROM users
			WHERE $1 = ANY(subjects) OR $1 = ''
			ORDER BY %s %s
			LIMIT $2 OFFSET $3`, orderby, order)

		rows, err = db.Query(sqlStatement, subject, pageSize, offset)
	}

	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var subjects []string

		err := rows.Scan(&user.Id, &user.Name, &user.Email, pq.Array(&subjects))
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to scan user: %v", err)
		}
		user.Subjects = subjects
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch users: %v", err)
	}

	var totalDocuments int
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE $1 = ANY(subjects) OR $1 = ''`
	err = db.QueryRow(countQuery, subject).Scan(&totalDocuments)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count total users: ")
	}

	lastPage := 1
	if pageSize != -1 {
		lastPage = (totalDocuments + pageSize - 1) / pageSize
	}

	return users, lastPage, totalDocuments, nil
}

func GetUserByID(id string) (model.User, error) {
	flagConfig, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading flag config: %v", err)
	}
	if flagConfig.FlagValue == "TRUE" {
		log.Println(" MongoDB flag is true, fetching users from in MongoDB")
		mongoClient, err := db.GetMongoDB()

		if err != nil {
			// log.Printf(" MongoDB initialization error: %v\n", err)
			return model.User{}, fmt.Errorf(" MongoDB is not initialized: ")
		}
		mongoCollection := mongoClient.Database("fitness").Collection("users")
		var user model.User
		err = mongoCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&user)
		if err != nil {
			log.Printf(" MongoDB error: %v", err)
			return model.User{}, fmt.Errorf(" user not found in MongoDB: %v", err)
		}
		return user, nil
	}
	// PostgreSQL logic
	db := db.GetPostgresDB()

	sqlStatement := `SELECT id, name, email, subjects FROM users WHERE id = $1`

	var user model.User
	var subjects []string
	err = db.QueryRow(sqlStatement, id).Scan(
		&user.Id, &user.Name, &user.Email, pq.Array(&subjects))

	if err != nil {
		return model.User{},fmt.Errorf("user not found in PostgreSQL%v\n:",err)
	}
	user.Subjects = subjects
	return user, nil
}
