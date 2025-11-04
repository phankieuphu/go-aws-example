package handlers

import (
	"github.com/phankieuphu/go-aws-example/internal/domain/services"

	"github.com/gin-gonic/gin"
)


func SetupRouters(handler *services.UserService) *gin.Engine{
	r := gin.Default()

	r.GET("/get-user", handler.GetUser)
	
	return r
}
