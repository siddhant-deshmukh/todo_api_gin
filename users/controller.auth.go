package users

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var token_key = getTokenKey()

func getTokenKey() string {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token_key := os.Getenv("TOKEN_KEY")
	if token_key == "" {
		log.Fatal("No connection string")
	}

	return token_key
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

	if id, isTypeId := result.InsertedID.(primitive.ObjectID); isTypeId {
		err = saveTokenString(c, id.Hex())
	} else {
		c.JSON(http.StatusInternalServerError, bson.M{
			"message": "Unexpcted outcome",
			"error":   err,
		})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{
			"message": "Some error occured while creating token",
			"error":   err,
		})
	}

	c.IndentedJSON(http.StatusCreated, bson.M{
		"result": result,
	})
}

// Login user
func userLogin(c *gin.Context) {

	var checkUser UserCredentials // Credentials submitted by user meaning username/email and password
	var user UserDocumentModel    // user details

	err := c.BindJSON(&checkUser) // checking if Credentials are in correct format
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Enter correct user meow"})
		return
	}
	filter := bson.D{primitive.E{Key: "email", Value: checkUser.Identifier}}

	err = UserColl.FindOne(Ctx, filter).Decode(&user) // match credentials with all emails
	if err != nil {
		filter = bson.D{primitive.E{Key: "username", Value: checkUser.Identifier}}
		err = UserColl.FindOne(Ctx, filter).Decode(&user) // match credentials with all UserName field
		if err != nil {                                   // the given credentail is neither in email nor in UserName
			c.JSON(http.StatusNotFound, gin.H{"message": "Enter correct user email or username"})
			return
		}
	}

	// check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(checkUser.Password))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Enter correct credentials"})
		return
	}

	err = saveTokenString(c, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{
			"message": "Some error occured while creating token",
			"error":   err,
		})
	}

	c.JSON(http.StatusAccepted, bson.M{
		"user": user,
	})
}

func saveTokenString(c *gin.Context, ID string) error {

	fmt.Println(ID)

	signing_key := []byte(token_key)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"_id": ID,
	})

	tokenString, err := token.SignedString(signing_key)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, bson.M{
			"message": "While creating token",
			"err":     err,
		})
		return err
	}

	c.SetCookie("todo_auth_token", tokenString, 364000, "/", "http://www.localhost:8080.com", false, true)
	return nil
}
