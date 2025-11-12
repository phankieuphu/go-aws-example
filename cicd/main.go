package main

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/scoring", func(c *gin.Context) {
		score := rand.Float64()
		approved_limit := rand.Intn(10000000)
		c.JSON(http.StatusOK, gin.H{
			"score":          score,
			"approved_limit": approved_limit,
		})

	})
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
