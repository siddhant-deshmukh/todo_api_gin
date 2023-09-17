package todos

import (
	"net/http"
	"time"

	"example.com/gin_01/users"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getTodoList(c *gin.Context) {
	user := c.MustGet("user").(users.UserDocs)

	filter := bson.D{{Key: "author_id", Value: user.ID}}
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

	if todos != nil {
		c.IndentedJSON(http.StatusOK, todos)
	} else {
		c.IndentedJSON(http.StatusOK, []string{})
	}
}

func getTodoById(c *gin.Context) {
	user := c.MustGet("user").(users.UserDocs)
	var todo TodoDocumentModel

	id := c.Param("id")
	basonId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invlaid todo object id"})
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: basonId}, {Key: "author_id", Value: user.ID}}

	singleResult := TodoColl.FindOne(Ctx, filter)
	if singleResult.Err() == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}

	err = singleResult.Decode(&todo)
	if err != nil {

		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, todo)
}

func postTodo(c *gin.Context) {
	user := c.MustGet("user").(users.UserDocs)

	var newTodoForm TodoForm

	err := c.BindJSON(&newTodoForm)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	newTodo := TodoModel{
		Title:       newTodoForm.Title,
		Description: newTodoForm.Description,
		State:       false,
		Author:      user.ID,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}
	result, err := TodoColl.InsertOne(Ctx, newTodo)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}

func editTodo(c *gin.Context) {
	user := c.MustGet("user").(users.UserDocs)

	var newTodo TodoForm
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

	filter := bson.D{primitive.E{Key: "_id", Value: basonId}, {Key: "author_id", Value: user.ID}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "title",
					Value: newTodo.Title,
				},
				{
					Key:   "description",
					Value: newTodo.Description,
				},
				{
					Key:   "state",
					Value: newTodo.State,
				},
			},
		},
	}
	opts := options.Update().SetUpsert(false)

	result, err := TodoColl.UpdateOne(Ctx, filter, update, opts)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	if result.MatchedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, bson.M{
			"message": "incorrect id",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

func deleteTodo(c *gin.Context) {
	user := c.MustGet("user").(users.UserDocs)

	id := c.Param("id")
	basonId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invlaid todo object id"})
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: basonId}, {Key: "author_id", Value: user.ID}}

	result, err := TodoColl.DeleteOne(Ctx, filter)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	if result.DeletedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, bson.M{
			"message": "incorrect id",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
