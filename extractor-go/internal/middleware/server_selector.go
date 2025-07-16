package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

// ServerType represents the different server types
const (
	ServerKRO   = "kro"
	ServerLATAM = "latam"
)

// ServerSelectorMiddleware is a middleware that checks the x-server header and sets the server type in context
func ServerSelectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("x-server")
		switch header {
		case ServerKRO:
			c.Set("x-server", server.GetKROMain())
		case ServerLATAM:
			c.Set("x-server", server.GetLATAM())
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid x-server header value. Must be either 'kro' or 'latam'",
			})
			return
		}
		c.Next()
	}
}
