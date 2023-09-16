package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func getUserInfo(c *gin.Context) {

}

func registerUser(c *gin.Context) {
	var newUser UserModel

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Enter user information in correct format"})
		return
	}

	var bytes []byte
	bytes, err = bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "While hashing password"})
		return
	}
	newUser.Password = string(bytes)

	result, err := UserColl.InsertOne(Ctx, newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}

func getUserById(c *gin.Context) {

}

func deleteUser(c *gin.Context) {

}

func editUser(c *gin.Context) {

}

func userLogin(c *gin.Context) {
	var checkUser UserCredentials
	var user UserDocumentModel

	err := c.BindJSON(&checkUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Enter correct user meow"})
		return
	}
	filter := bson.D{primitive.E{Key: "email", Value: checkUser.Identifier}}

	err = UserColl.FindOne(Ctx, filter).Decode(&user)
	if err != nil {
		filter = bson.D{primitive.E{Key: "username", Value: checkUser.Identifier}}
		err = UserColl.FindOne(Ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Enter correct user email or username"})
			return
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(checkUser.Password))
	if err == nil {
		c.JSON(http.StatusAccepted, gin.H{"message": "Sucessfull"})
	} else {
		c.JSON(http.StatusConflict, gin.H{"message": "Enter correct credentials"})
	}
}
