package app

import (
	"gin_jwt/controllers"
	"gin_jwt/middlewares"
	"gin_jwt/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(middlewares.CORSMiddleware())

	models.ConnectDatabase()
	models.DBMigrate()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.DELETE("/logout", controllers.Logout)

	r.GET("/me", middlewares.AuthMiddleware, controllers.GetMe)
	r.Any("/verify", middlewares.AuthMiddleware, controllers.GetMe)

	return r
}
