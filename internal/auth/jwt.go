package auth

import (
	"log"
	"net/http"

	"github.com/ali-shokoohi/micro-gopia/scripts"
	"github.com/gin-gonic/gin"
)

type Auth interface {
	JwtAuthMiddleware() gin.HandlerFunc
}

// auth represents the service.
type auth struct {
}

// NewAuth returns a new instance of Auth.
func NewAuth() Auth {
	return &auth{}
}

func (a *auth) JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := scripts.ValidateToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Authentication required"})
			log.Printf("An error while verifying token in middleware: %s", err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
