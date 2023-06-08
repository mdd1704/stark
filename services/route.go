package services

import (
	"context"

	"github.com/gin-gonic/gin"

	"stark/database"
	"stark/services/auth"
	"stark/services/client"
	"stark/services/profile"
	"stark/services/user"
	"stark/services/user_detail"
	"stark/services/user_location"
	"stark/utils/log"
	"stark/utils/middleware"
)

func InitRoute(
	ctx context.Context,
	router *gin.Engine,
	redisDB *database.Redis,
	clientService *client.Service,
	clientHandler *client.Handler,
	userHandler *user.Handler,
	userDetailHandler *user_detail.Handler,
	userLocationHandler *user_location.Handler,
	authHandler *auth.Handler,
	profileHandler *profile.Handler,
) {
	// Internal group
	internal := router.Group("/internal")
	internal.Use(middleware.InternalMiddleware())

	// Client service
	internal.POST("/client", clientHandler.HandleCreate)
	internal.GET("/client/:id", clientHandler.HandleDetail)
	internal.PUT("/client/:id", clientHandler.HandleUpdate)
	internal.POST("/client/filter", clientHandler.HandleAllByFilter)
	internal.GET("/client", clientHandler.HandlePage)

	// Client group
	client := router.Group("/client")
	client.Use(middleware.ClientMiddleware(clientService))

	// User service
	client.POST("/user", userHandler.HandleCreate)
	client.GET("/user/:id", userHandler.HandleDetail)
	client.PUT("/user/:id", userHandler.HandleUpdate)
	client.POST("/user/filter", userHandler.HandleAllByFilter)
	client.GET("/user", userHandler.HandlePage)

	// User detail service
	client.POST("/user-detail", userDetailHandler.HandleCreate)
	client.GET("/user-detail/:id", userDetailHandler.HandleDetail)
	client.PUT("/user-detail/:id", userDetailHandler.HandleUpdate)
	client.POST("/user-detail/filter", userDetailHandler.HandleAllByFilter)
	client.GET("/user-detail", userDetailHandler.HandlePage)

	// User location service
	client.POST("/user-location", userLocationHandler.HandleCreate)
	client.GET("/user-location/:id", userLocationHandler.HandleDetail)
	client.PUT("/user-location/:id", userLocationHandler.HandleUpdate)
	client.POST("/user-location/filter", userLocationHandler.HandleAllByFilter)
	client.GET("/user-location", userLocationHandler.HandlePage)

	// API group
	api := router.Group("/api")

	// Auth service
	api.POST("/login", authHandler.HandleLogin)
	api.POST("/register", authHandler.HandleRegister)
	api.POST("/refresh-token", authHandler.HandleRegister)
	api.Use(middleware.AuthMiddleware(redisDB))
	api.GET("/logout", authHandler.HandleLogout)

	// Profile service
	api.POST("/update-profile", profileHandler.HandleUpdateProfile)
	api.POST("/change-password", profileHandler.HandleChangePassword)

	router.GET("/ping", func(c *gin.Context) {
		log.WithContext(ctx).Info("when you ping, then you get pong!")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
