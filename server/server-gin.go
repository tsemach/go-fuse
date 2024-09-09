package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

func SendError(status int, code int, err error, c *gin.Context) {
	log.Fatalln("Fatal Error", err)

	c.JSON(
		status,
		gin.H{"error-code": code, "error": err})
}

func CreateGin() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	serverFuseFS(context.TODO(), r)

	return r
}
