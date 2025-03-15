package controllers

import (
	"ansible-runner/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

// AuthMiddleware is a simple authentication middleware example.
// You should replace this with your actual authentication logic.
func AuthMiddleware(c *gin.Context) {
	// Check for an authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}
	// Basic authentication format: "Basic <base64 encoded string>"
	// Extract base64 encoded credentials
	authType, credentials, found := utils.ParseBasicAuth(authHeader)

	if !found || authType != "Basic" {

		utils.Log.Error("Invalid Authorization header format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	// decode
	username, password, err := utils.DecodeBasicAuth(credentials)
	if err != nil {
		utils.Log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// LDAP configuration
	ldapServer := viper.GetString("ldap.server")
	ldapPort := viper.GetInt("ldap.port")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		utils.Log.Error(fmt.Sprintf("LDAP connection failed: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "LDAP connection failed"})
		return
	}
	defer l.Close()

	// Bind with the user's credentials
	err = l.Bind(username, password)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("LDAP bind failed for user %s: %v", username, err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	c.Set("username", username)

	// Continue to the next handler if authenticated
	c.Next()

}
