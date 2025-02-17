package main

import (
	"ansible-runner/configs"
	"ansible-runner/routers"
)

func main() {

	configs.LoadConfig()
	routers.Server()
}
