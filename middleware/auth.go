package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware : to verify all authorized operations
func AuthMiddleware(c *gin.Context) {
	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)

	authorizationToken := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))

	if idToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"title": "Authentication Required", "msg": "You need to login in order to use this service"})
		c.Abort()
		return
	}

	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"title": "Token Invalid", "msg": "Access Token is not valid"})
		c.Abort()
		return
	}

	c.Set("UUID", token.UID)
	c.Next()
}
