package main

import (
	"module_git/content"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	port := "0.0.0.0:8080"

	router := gin.Default()

	// Exchanging routre to sending emails
	router.POST("/v1/api/emails", content.Post_contents)
	router.GET("/v1/api/emails", content.Get_contents)

	// Statrting server on port 8080
	if err := router.Run(port); err != nil {
		panic("Failed to start server: " + err.Error())
	}

}
