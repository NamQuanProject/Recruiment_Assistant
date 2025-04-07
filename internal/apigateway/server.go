package apigateway

import (
	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/initializers"
	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/middleware"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/user"
	"github.com/gin-gonic/gin"
)

func Init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func RunServer() {
	r := gin.Default()
	Init()

	r.POST("/signup", user.Signup)
	r.POST("/login", user.Login)
	r.GET("/validate", middleware.RequireAuth, user.Validate)
	r.POST("/logout", user.Logout) 
	r.GET("/test", middleware.RequireAuth, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	r.Run()
}