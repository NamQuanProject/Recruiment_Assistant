package apigateway

import (
	"net/http"

	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/handlers"
	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/initializers"
	"github.com/gin-gonic/gin"
)

func Init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func RunServer() {
	r := gin.Default()

	// Load HTML templates
	// r.LoadHTMLGlob("./templates/*")

	// // Init()

	// // Authentication routes
	// r.POST("/signup", user.Signup)
	// r.POST("/login", user.Login)
	// r.POST("/logout", user.Logout)

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

	r.Run(":8081")
}
