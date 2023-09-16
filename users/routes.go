package users

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(userRoutesGroup *gin.RouterGroup) {

	userRoutesGroup.GET("/", getUserInfo)
	userRoutesGroup.GET("/:id", getUserById)
	userRoutesGroup.POST("/", registerUser)
	userRoutesGroup.PUT("/:id", editUser)
	userRoutesGroup.DELETE("/:id", deleteUser)

}
