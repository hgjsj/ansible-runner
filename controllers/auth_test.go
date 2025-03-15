package controllers

import (
	"ansible-runner/configs"
	"ansible-runner/utils"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func init() {
	configs.LoadConfig("../")
	utils.InitLogger("../logs/test.log")
}

func TestAuth(t *testing.T) {
	user := "sec\\fsun4"
	pwd := "lh0kp)incT@tvtdc"
	user = user[strings.Index(user, "\\") + 1 :]
	// LDAP configuration
	ldapServer := viper.GetString("ldap.server")
	ldapPort := viper.GetInt("ldap.port")
	baseDN := viper.GetString("ldap.baseDN")
	err := Auth(ldapServer, ldapPort, user, pwd, baseDN)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	t.Log("Auth successful")
}
