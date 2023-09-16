package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getUserInfo(c *gin.Context) {
	auth_token, err := c.Cookie("todo_auth_token")
	if err != nil {
		c.JSON(http.StatusNotAcceptable, bson.M{
			"message": "Error in cookie",
			"err":     err,
		})
		return
	}

	// var token
	token, err := jwt.Parse(auth_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(token_key), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{
			"error":   err,
			"message": "Internal Server error (Parse)",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		id := claims["_id"].(string)
		user, err := getUserDocumentById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{
				"error":   err,
				"message": "Internal Server error (Claim)",
			})
			return
		}
		c.JSON(http.StatusAccepted, bson.M{
			"msg":    "Successfull",
			"claims": claims,
			"user":   user,
		})
	} else {
		c.JSON(http.StatusNotAcceptable, bson.M{
			"msg": "Invalid token",
		})
	}
}

func getUserById(c *gin.Context) {

}

func deleteUser(c *gin.Context) {

}

func editUser(c *gin.Context) {

}

func getUserDocumentById(id string) (UserDocs, error) {

	var user UserDocumentModel
	basonId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return UserDocs{}, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: basonId}}

	err = UserColl.FindOne(Ctx, filter).Decode(&user)
	if err != nil {
		return UserDocs{}, err
	}

	return UserDocs{
		ID:       user.ID,
		Name:     user.Name,
		UserName: user.UserName,
	}, nil
}
