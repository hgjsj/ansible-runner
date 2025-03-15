package controllers

import (
	"ansible-runner/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

func Auth(server string, port int, user string, pwd string, baseDN string, group string)  error {
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

	// user is start with sec\\
	if strings.HasPrefix(strings.ToLower(user), "sec\\") {
		user = user[strings.Index(user, "\\") + 1 :]
	} else if (strings.Contains(user, "@")) {
		user = user[:strings.Index(user, "@")]
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
			if strings.HasPrefix(strings.ToLower(memberOf), strings.ToLower(group)) {
				return nil
			}
		}
	}
	return fmt.Errorf("user %s is not member of group %s", user, group)
}



// AuthMiddleware is a simple authentication middleware example.
func AuthMiddleware(c *gin.Context) {
	authH := c.GetHeader("Authorization")
	if authH == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}
	// Basic authentication format: "Basic <base64 encoded string>"
	// Extract base64 encoded credentials
	authType, credentials, found := utils.ParseBasicAuth(authH)

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

	utils.Log.Info(fmt.Sprintf("Username in header %s", username))
	utils.Log.Info(fmt.Sprintf("Password in header %s", password))

	// LDAP configuration
	ldapServer := viper.GetString("ldap.server")
	ldapPort := viper.GetInt("ldap.port")
	baseDN := viper.GetString("ldap.baseDN")
	requiredGroup := viper.GetString("ldap.requiredGroup")


	err = Auth(ldapServer, ldapPort, username, password, baseDN, requiredGroup)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} 

	// Continue to the next handler if authenticated
	c.Next()

}
