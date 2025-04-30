package apigateway

import (
	"net/http"
	"os"

	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/handlers"
	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/initializers"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/user"
	"github.com/gin-contrib/cors" //chau added this
	"github.com/gin-gonic/gin"
)

func Init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func RunServer() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow requests from your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})) // Chau added this

	r.Static("/storage", "./storage")
	r.Static("/internal", "./internal")

	// Authentication routes
	r.POST("/signup", user.Signup)
	r.POST("/login", user.Login)
	r.POST("/logout", user.Logout)

	// Serve HTML form upload
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	// Create a route for testing the upload without auth (temporary)
	r.GET("/upload-test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	// Job description routes
	r.POST("/submitJD", handlers.SubmitJDHandler)
	r.POST("/submitCVs", handlers.SubmitCVsHandler)
	r.POST("/getHlCV", handlers.GetHlCVHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback khi chạy local
	}

	r.Run(":" + port)
}
