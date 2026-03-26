package controllers

import (
	"gin_jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type formData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *gin.Context) {
	var data formData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the user exists already
	if !models.CheckUserAvailability(data.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email missing or already taken"})
		return
	}

	// Create the user
	user := models.UserCreate(data.Email, data.Password)
	if user == nil || user.ID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user creation failed"})
		return
	}

	// Create JWT token
	tokenString, err := createAndSignJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JWT creation failed"})
		return
	}

	// 2. Send the token in a cookie
	setCookie(c, tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": "Signup successful",
		"token":   tokenString,
	})
}

func Login(c *gin.Context) {
	var data formData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Match password
	user := models.UserMatchPassword(data.Email, data.Password)
	if user.ID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		return
	}

	// Create JWT token
	tokenString, err := createAndSignJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JWT creation failed"})
		return
	}

	// 2. Send the token in a cookie
	setCookie(c, tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func Logout(c *gin.Context) {
	// Add the JWT token to the block list.
	// or change expiry time of the cookie.

	c.SetCookie("Auth", "deleted", 0, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func GetMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func createAndSignJWT(user *models.User) (string, error) {
	// 1. Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID.Hex(),
		"ttl":    time.Now().Add(time.Hour * 24 * 100).Unix(),
	})

	hmacSampleSecret := os.Getenv("JWT_SECRET")
	if hmacSampleSecret == "" {
		hmacSampleSecret = "e1bed9f5-81d7-4810-9f9b-307d2761c4d4"
	}

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(hmacSampleSecret))
}

func setCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*100, "", "", false, true)
}
