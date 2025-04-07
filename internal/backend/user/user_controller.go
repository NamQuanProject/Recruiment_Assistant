package user

import (
	"net/http"
	"os"
	"time"

	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway/initializers"
	"github.com/KietAPCS/test_recruitment_assistant/internal/database/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get the email and password from the request body
	var body struct {
		Email string
		Password string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}


	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create a new user in the database
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	// Get the email and password from the request body
	var body struct {
		Email string
		Password string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Check if the user exists in the database
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare the provided password with the hashed password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate a JWT token for the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Return the JWT token in the response
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET"))) // Use a secure secret key in production

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// For now, just return a success message
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
		"status":  "ok",
	})
}

func Logout(c *gin.Context) {
    c.SetSameSite(http.SameSiteNoneMode)
    c.SetCookie("Authorization", "", -1, "", "", false, true)
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Logout successful",
    })
}