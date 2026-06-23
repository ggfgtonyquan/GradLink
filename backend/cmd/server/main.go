package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gradlink/internal/db"
	apphttp "gradlink/internal/http"
)

func main() {
	database, err := db.Open(databasePath())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}
	if err := db.SeedWithOptions(database, db.SeedOptions{
		AdminPassword: os.Getenv("GRADLINK_ADMIN_PASSWORD"),
	}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	apphttp.RegisterRoutes(r, &apphttp.Deps{DB: database})

	if err := r.Run(listenAddr()); err != nil {
		log.Fatal(err)
	}
}

func databasePath() string {
	if value := os.Getenv("GRADLINK_DB_PATH"); value != "" {
		return value
	}
	return "data/gradlink.db"
}

func listenAddr() string {
	if value := os.Getenv("GRADLINK_ADDR"); value != "" {
		return value
	}
	return ":8080"
}
