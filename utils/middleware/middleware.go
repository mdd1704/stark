package middleware

import (
	"net/http"
	"os"
	"stark/database"
	"stark/respond"
	"stark/services/client"
	"stark/utils"

	"github.com/gin-gonic/gin"
)

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func InternalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		internalID := c.Request.Header.Get("X-Internal-ID")
		if internalID != os.Getenv("INTERNAL_ID") {
			c.Abort()
			respond.Error(c, "", http.StatusBadRequest, "BadRequest", "Internal ID is invalid")
			return
		}

		c.Next()
	}
}

func ClientMiddleware(clientService *client.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		bearerKey, err := utils.GetBearerKey(bearerToken)
		if err != nil {
			c.Abort()
			respond.Error(c, "", http.StatusUnauthorized, "Unauthorized", err.Error())
			return
		}

		filter := client.Filter{
			BearerKeys: []string{bearerKey},
		}

		client, err := clientService.FindAllByFilter(filter)
		if err != nil || len(client) == 0 {
			c.Abort()
			respond.Error(c, "", http.StatusUnauthorized, "Unauthorized", "Client not found")
			return
		}

		c.Set("client_id", client[0].ID.String())
		c.Next()
	}
}

func AuthMiddleware(redisDB *database.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		bearerKey, err := utils.GetBearerKey(bearerToken)
		if err != nil {
			c.Abort()
			respond.Error(c, "", http.StatusUnauthorized, "Unauthorized", err.Error())
			return
		}

		metadata, err := utils.ExtractAccessTokenMetadata(bearerKey)
		if err != nil {
			c.Abort()
			respond.Error(c, "", http.StatusUnauthorized, "Unauthorized", err.Error())
			return
		}

		userID, err := utils.FetchAccessAuth(metadata, redisDB)
		if err != nil {
			c.Abort()
			respond.Error(c, "", http.StatusUnauthorized, "Unauthorized", err.Error())
			return
		}

		c.Set("access_uuid", metadata.AccessUuid)
		c.Set("user_id", userID)
		c.Next()
	}
}
