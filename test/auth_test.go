package test

import (
	"fmt"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

func TestAuth(t *testing.T) {
	auth := "sec\fsun4"
	pawd := "lh0kp)incT@tvtdc"
	// LDAP configuration
	ldapServer := viper.GetString("ldap.server")
	ldapPort := viper.GetInt("ldap.port")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		t.Fatalf("LDAP connection failed: %v", err)
	}
	defer l.Close()

	// Bind with the user's credentials
	err = l.Bind(auth, pawd)
	if err != nil {
		t.Fatalf("LDAP bind failed: %v", err)
	}
	t.Log("LDAP bind successful")
}
