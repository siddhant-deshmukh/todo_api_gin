package todos

import "github.com/gin-gonic/gin"

func RegisterTodoRoutes(todoRoutesGroup *gin.RouterGroup) {

	todoRoutesGroup.GET("/", getTodoList)
	todoRoutesGroup.GET("/:id", getTodoById)
	todoRoutesGroup.POST("/", postTodo)
	todoRoutesGroup.PUT("/:id", editTodo)
	todoRoutesGroup.DELETE("/:id", deleteTodo)

}
