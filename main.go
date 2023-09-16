package main

import (
	"context"
	"log"
	"os"

	"example.com/gin_01/todos"
	"example.com/gin_01/users"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var client *mongo.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connection_string := os.Getenv("MONGODB_CONNECTION_STRING")
	if connection_string == "" {
		log.Fatal("No connection string")
	}

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection_string))
	if err != nil {
		log.Fatal("Unable to connect the mongodb database")
	}

	todoColl := client.Database("test").Collection("todo")
	todos.SetTodoCollection(todoColl, ctx)

	// userColl := client.Database("test").Collection("users")
	users.ManageUserCollection(client)
	// users.SetUserCollection(userColl, ctx)

	router := gin.Default()

	todoRoutesGroup := router.Group("/todo")
	todos.RegisterTodoRoutes(todoRoutesGroup)

	userRoutesGroup := router.Group("/user")
	userRoutesGroup.Use(users.AuthUserMiddleware())
	users.RegisterUserRoutes(userRoutesGroup)

	userAuthRoutesGroup := router.Group("/")
	users.RegisterUserAuthRoutes(userAuthRoutesGroup)

	router.Run("localhost:8080")
}
