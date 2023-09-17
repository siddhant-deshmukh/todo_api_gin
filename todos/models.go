package todos

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TodoColl *mongo.Collection
var Ctx context.Context

type TodoDocumentModel struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Title       string             `json:"title" binding:"required"`
	State       bool               `json:"state" default:"false"`
	Description string             `json:"description"`
	Author      primitive.ObjectID `bson:"author_id" json:"author_id" binding:"required"`
	CreatedAt   primitive.DateTime `bson:"created_at" binding:"required"`
}

type TodoModel struct {
	Title       string             `json:"title" binding:"required"`
	State       bool               `json:"state" default:"false"`
	Description string             `json:"description"`
	Author      primitive.ObjectID `bson:"author_id" binding:"required"`
	CreatedAt   primitive.DateTime `bson:"created_at" binding:"required"`
}

type TodoForm struct {
	Title       string `json:"title" binding:"required"`
	State       bool   `json:"state" default:"false"`
	Description string `json:"description"`
}

func SetTodoCollection(pointer1 *mongo.Collection, pointer2 context.Context) {
	TodoColl = pointer1
	Ctx = pointer2
}

func ManageTodoCollection(client *mongo.Client) {
	Ctx = context.TODO()
	database := client.Database("test")

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"title":    "todo object validation",
			"required": []string{"title", "description", "author_id", "created_at", "state"},
			"properties": bson.M{
				"title": bson.M{
					"bsonType":    "string",
					"minLength":   1,
					"maxLength":   50,
					"description": "string with length between 1 to 50",
				},
				"description": bson.M{
					"bsonType":    "string",
					"minLength":   0,
					"maxLength":   200,
					"description": "string with length between 0 to 200",
				},
				"state": bson.M{
					"bsonType": "bool",
				},
				"author_id": bson.M{
					"bsonType": "objectId",
				},
				"created_at": bson.M{
					"bsonType": "date",
				},
			},
		},
	}
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("created_at"),
		},
		{
			Keys:    bson.D{{Key: "author_id", Value: 1}},
			Options: options.Index().SetName("author_id"),
		},
	}

	opts := options.CreateCollection().SetValidator(validator)
	err := database.CreateCollection(Ctx, "todo", opts)

	todoColl := client.Database("test").Collection("todo")
	TodoColl = todoColl

	if err != nil {
		fmt.Println("Collection already exist!")
	} else {
		_, err := todoColl.Indexes().CreateMany(Ctx, indexModels)
		if err != nil {
			log.Fatal(err)
		}
	}
}
