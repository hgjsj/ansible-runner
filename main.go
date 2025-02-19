package main

import (
	"ansible-runner/configs"
	"ansible-runner/routers"
	"ansible-runner/utils"
	"ansible-runner/db"
	"os"
	"os/signal"
	"syscall"
)

func cleanup() {
	routers.ShutDown()
	db.MongoClient.Close()
	utils.CleanUp()
}

func main() {
	utils.InitLogger()
	configs.LoadConfig("./")
	routers.Server()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	utils.Log.Info("Cleanup and shutting down server...")
	cleanup()
}
