package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wchr-aun/wasted-backend/config"
	"github.com/wchr-aun/wasted-backend/endpoints"
	"github.com/wchr-aun/wasted-backend/middleware"
)

func main() {
	const port string = ":3000"

	router := gin.Default()
	firebaseAuth := config.SetupFirebase()
	firestoreCon := config.ConnectFirestore()

	router.Use(middleware.CORSMiddleware())
	router.Use(func(c *gin.Context) {
		c.Set("firebaseAuth", firebaseAuth)
	})
	router.Use(func(c *gin.Context) {
		c.Set("firestoreCon", firestoreCon)
	})
	router.Use(middleware.AuthMiddleware)

	router.GET("/api/auth", endpoints.Authentication)

	router.Run(port)
}
