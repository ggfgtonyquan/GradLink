package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gradlink/internal/db"
	apphttp "gradlink/internal/http"
)

func main() {
	database, err := db.Open("data/gradlink.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}
	if err := db.Seed(database); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	apphttp.RegisterRoutes(r, &apphttp.Deps{DB: database})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
