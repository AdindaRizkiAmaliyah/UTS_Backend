package main

import (
	"log"
	"os"
	"clean-archi/config"
	"clean-archi/database"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()
	defer database.DB.Close()

	app := config.NewApp(database.DB)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
