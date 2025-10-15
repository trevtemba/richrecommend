package main

import (
	"github.com/trevtemba/richrecommend/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.SetupRoutes(r)

	r.Run(":8080")
}
