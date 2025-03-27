package api

import (
	"ansible-runner/internal/utils"
	"github.com/spf13/viper"
	"testing"
)

func init() {
	utils.LoadConfig("../../configs/config")
	utils.InitLogger()
}

func TestAuth(t *testing.T) {
	user := "sec\\fsun4"
	pwd := "lh0kp)incT@tvtdc"
	// LDAP configuration
	ldapServer := viper.GetString("ldap.server")
	ldapPort := viper.GetInt("ldap.port")
	baseDN := viper.GetString("ldap.baseDN")
	requiredGroup := viper.GetString("ldap.requiredGroup")
	err := Auth(ldapServer, ldapPort, user, pwd, baseDN, requiredGroup)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	t.Log("Auth successful")
}
