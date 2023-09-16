package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUserInfo(c *gin.Context) {

}

func registerUser(c *gin.Context) {
	var newUser UserModel

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Enter correct user meow"})
		return
	}

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
