package users

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserColl *mongo.Collection
var Ctx context.Context

type UserDocumentModel struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name" binding:"required"`
	UserName string             `json:"username" binding:"required"`
	Email    string             `json:"email" binding:"required"`
	Password string             `json:"password" binding:"required"`
}

type UserDocs struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name" binding:"required"`
	UserName string             `json:"username" binding:"required"`
}

type UserModel struct {
	Name     string `json:"name" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserCredentials struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// type UserTokenClaims struct {
// 	*jwt.Cla
// 	TokenType string
// 	CustomerInfo
// }

func SetUserCollection(pointer1 *mongo.Collection, pointer2 context.Context) {
	UserColl = pointer1
	Ctx = pointer2
}

func ManageUserCollection(client *mongo.Client) {
	Ctx = context.TODO()
	database := client.Database("test")
	// fmt.Println(database.Collection("users").Name())

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"title":    "user object validation",
			"required": []string{"name", "email", "username", "password"},
			"properties": bson.M{
				"name": bson.M{
					"bsonType":    "string",
					"minLength":   3,
					"maxLength":   20,
					"description": "string with length between 3 to 20",
				},
				"email": bson.M{
					"bsonType": "string",
					"pattern":  "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
				},
				"username": bson.M{
					"bsonType":    "string",
					"minLength":   1,
					"maxLength":   10,
					"description": "string with length between 1 to 10",
				},
				"password": bson.M{
					"bsonType":    "string",
					"description": "string with length between 5 to 20",
				},
			},
		},
	}
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetName("email").SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetName("username").SetUnique(true),
		},
	}

	opts := options.CreateCollection().SetValidator(validator)
	err := database.CreateCollection(Ctx, "users", opts)

	userColl := client.Database("test").Collection("users")
	UserColl = userColl

	if err != nil {
		fmt.Println("Collection already exist!")
	} else {
		_, err := userColl.Indexes().CreateMany(Ctx, indexModels)
		if err != nil {
			log.Fatal(err)
		}
	}
}
