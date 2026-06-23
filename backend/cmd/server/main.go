package main

import (
	"log"

	"github.com/gin-gonic/gin"
	apphttp "gradlink/internal/http"
)

func main() {
	r := gin.Default()
	apphttp.RegisterRoutes(r, nil)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
