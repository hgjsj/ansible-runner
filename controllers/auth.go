package controllers

import (
	"ansible-runner/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

func auth(server string, port int, user string, pwd string, baseDN string)  error {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		utils.Log.Error(fmt.Sprintf("LDAP connection failed: %v", err))
		return err
	}
	defer l.Close()

	// Bind with the user's credentials
	err = l.Bind(user, pwd)

	if err != nil {
		utils.Log.Error(fmt.Sprintf("LDAP bind failed for user %s: %v", user, err))
		return err
	}

	searchRequest := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", user), []string{"memberOf"}, nil)

	sr, err := l.Search(searchRequest)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("LDAP search failed for user %s: %v", user, err))
		return err
	}
	for _, entry := range sr.Entries {
		for _, memberOf := range entry.GetAttributeValues("memberOf") {
			utils.Log.Info(fmt.Sprintf("User %s is member of %s", user, memberOf))
		}
	}
	return nil
}



// AuthMiddleware is a simple authentication middleware example.
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	utils.Log.Info("Username in header", username)
	utils.Log.Info("Password in header", password)

	// LDAP configuration
	ldapServer := viper.GetString("ldap.server")
	ldapPort := viper.GetInt("ldap.port")
	baseDN := viper.GetString("ldap.baseDN")
	utils.Log.Info("ldapserver in config", ldapServer)
	utils.Log.Info("port in config", ldapPort)
	err = auth(ldapServer, ldapPort, username, password, baseDN)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} 

	// Continue to the next handler if authenticated
	c.Next()

}
