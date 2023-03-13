package main

import (
	"github.com/sirupsen/logrus"
	"todo/database"
	"todo/server"
)

func main() {
	// Create a server instance
	srv := server.SetupRoutes()
	if err := database.ConnectAndMigrate(
		"localhost",
		"5432",
		"postgres",
		"local",
		"local",
		database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Info("migration successful!!")

	if err := srv.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to run server with error: %+v", err)
	}
	logrus.Print("Server started at :8080")
}
