package todos

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var TodoColl *mongo.Collection
var Ctx context.Context

type TodoDocumentModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `json:"title" binding:"required"`
	State       bool               `json:"State" default:"false"`
	Description string             `json:"description"`
}

type TodoModel struct {
	Title       string `json:"title" binding:"required"`
	State       bool   `json:"State" default:"false"`
	Description string `json:"description"`
}

func SetTodoCollection(pointer1 *mongo.Collection, pointer2 context.Context) {
	TodoColl = pointer1
	Ctx = pointer2
}
