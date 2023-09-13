package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `json:"title" binding:"required"`
	State       bool               `json:"state" default:"false"`
	Description string             `json:"description"`
}

var todoColl *mongo.Collection

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connection_string := os.Getenv("MONGODB_CONNECTION_STRING")
	if connection_string == "" {
		log.Fatal("No connection string")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connection_string))
	if err != nil {
		log.Fatal("Unable to connect the mongodb database")
	}

	todoColl = client.Database("test").Collection("todo")

	router := gin.Default()

	router.GET("/", getTodoList)
	// router.GET("/:id", getTodoById)
	router.POST("/", postTodo)
	router.PUT("/:id", editTodo)
	router.DELETE("/:id", deleteTodo)

	router.Run("localhost:8080")
}

func getTodoList(c *gin.Context) {
	filter := bson.D{}
	var todos []TodoModel

	cursor, err := todoColl.Find(context.TODO(), filter)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	if err = cursor.All(context.TODO(), &todos); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, todos)
}

func postTodo(c *gin.Context) {
	var newTodo TodoModel

	err := c.BindJSON(&newTodo)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	// albums = append(albums, newTodo)
	result, err := todoColl.InsertOne(context.TODO(), newTodo)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, result)
}

func editTodo(c *gin.Context) {
	var newTodo TodoModel
	id := c.Param("id")

	err := c.BindJSON(&newTodo)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	basonId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invlaid todo object id"})
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: basonId}}
	// replacement :=
	newTodo.ID = basonId
	// albums = append(albums, newTodo)
	result, err := todoColl.ReplaceOne(context.TODO(), filter, newTodo)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	basonId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invlaid todo object id"})
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: basonId}}

	result, err := todoColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
