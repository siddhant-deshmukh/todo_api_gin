package todos

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getTodoList(c *gin.Context) {
	filter := bson.D{}
	var todos []TodoDocumentModel

	cursor, err := TodoColl.Find(Ctx, filter)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	if err = cursor.All(Ctx, &todos); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, todos)
}

func getTodoById(c *gin.Context) {

	var todo TodoDocumentModel

	id := c.Param("id")
	basonId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invlaid todo object id"})
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: basonId}}

	err = TodoColl.FindOne(Ctx, filter).Decode(&todo)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, todo)
}

func postTodo(c *gin.Context) {
	var newTodo TodoModel

	err := c.BindJSON(&newTodo)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	result, err := TodoColl.InsertOne(Ctx, newTodo)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}

func editTodo(c *gin.Context) {
	var newTodo TodoDocumentModel
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

	newTodo.ID = basonId
	result, err := TodoColl.ReplaceOne(Ctx, filter, newTodo)
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

	result, err := TodoColl.DeleteOne(Ctx, filter)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
