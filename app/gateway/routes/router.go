package routes

import (
	"go-micro-toDoList/app/gateway/internal/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	// maybe some other common middleware here
	// session setting here
	apiGroup := router.Group("/api/")

	// User service
	userApiGroup := apiGroup.Group("/user/")
	userApiGroup.POST("login", http.UserLogin)
	userApiGroup.POST("register", http.UserRegister)

	// Task service

	return router
}
