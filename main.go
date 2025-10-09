package main

import (
	"ansible-runner/internal/cmd"
)

// @title Ansible Runner API
// @version 1.0
// @description This is an API server for managing and executing Ansible playbooks.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email fu-qi.sun@dxc.com
// @license.name Apache 2.0	
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth

func main() {
	cmd.Execute()
}
