package main

import (
	"module_git/content"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)


func main() {

	router := gin.Default()

	// Exchanging routre to sending emails
	router.POST("/v1/api/emails", content.PostContents)

	// Statrting server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}

}
