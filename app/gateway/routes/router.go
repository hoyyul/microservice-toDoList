package routes

import (
	"micro-toDoList/app/gateway/internal/http"
	"micro-toDoList/global"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(global.Config.Server.Env)
	router := gin.Default()
	// maybe some other common middleware here
	// session setting here
	apiGroup := router.Group("/api/")

	// User service
	userApiGroup := apiGroup.Group("/user/")
	userApiGroup.POST("login", http.UserLogin)
	userApiGroup.POST("register", http.UserRegister)
	userApiGroup.DELETE("logout", http.UserLogout)

	// Task service
	taskApiGroup := apiGroup.Group("/task/")
	taskApiGroup.POST("/", http.TaskCreate)
	taskApiGroup.DELETE("/", http.TaskDelete)
	taskApiGroup.PUT("/", http.TaskUpdate)
	taskApiGroup.GET("/", http.TaskShow)

	return router
}
