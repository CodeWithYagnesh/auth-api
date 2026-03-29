package middlewares

import (
	"fmt"
	"gin_jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(c *gin.Context) {
	// Retrieve the cookie or Authorization header from the request
	tokenStr, err := c.Cookie("Auth")
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:]
			err = nil
		}
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No auth token"})
		c.Abort()
		return
	}

	// Extract the JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		hmacSampleSecret := os.Getenv("JWT_SECRET")
		if hmacSampleSecret == "" {
			hmacSampleSecret = "e1bed9f5-81d7-4810-9f9b-307d2761c4d4"
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse JWT token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT Claims failed"})
		c.Abort()
		return
	}

	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token expired!"})
		c.Abort()
		return
	}

	// Extract the user from the token
	userIDStr, ok := claims["userID"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID in token"})
		c.Abort()
		return
	}

	objID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ObjectID format"})
		c.Abort()
		return
	}

	user := models.UserFromId(objID)

	if user.ID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not find the user!"})
		c.Abort()
		return
	}

	// Set the current user in the context
	c.Set("user", user)

	// Go to the next in chain
	c.Next()
}
