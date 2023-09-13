package main

import (
	"context"
	"log"
	"os"

	"example.com/gin_01/todos"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	router := gin.Default()

	todoRoutesGroup := router.Group("/")
	todos.RegisterTodoRoutes(todoRoutesGroup)

	router.Run("localhost:8080")
}
