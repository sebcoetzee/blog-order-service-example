package main

import (
	"github.com/SebastianCoetzee/blog-order-service-example/application"
	"github.com/SebastianCoetzee/blog-order-service-example/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	app.GET("/users/:id/orders", handlers.FindOrdersForUser)
	app.Run()

	defer application.CloseDB()
}
