package main

import (
	"fmt"
	"module_git/content"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	conf_all, err := content.Load_config()
	server_conf := conf_all.Server
	if err != nil {
		fmt.Println("Error loading config:", err)

		return
	}

	router := gin.Default()

	// Exchanging routre to sending emails
	router.POST("/v1/api/emails", content.Post_contents)

	// Statrting server on port 8080
	if err := router.Run(server_conf.Port); err != nil {
		panic("Failed to start server: " + err.Error())
	}

}
