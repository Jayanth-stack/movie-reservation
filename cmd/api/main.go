package main

import (
	"log"
	"movie-reservation/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	db.Connect()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the Movie Reservation API"})
	})

	log.Println("Starting server on port 8080")
	r.Run(":8080")

}
