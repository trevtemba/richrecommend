package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trevtemba/richrecommend/internal/api/handlers"
)

func SetupRoutes(router *gin.Engine) {

	// userGroup := router.Group("/users")
	// {
	// 	userGroup.GET("/", handlers.GetUsers)
	// 	userGroup.POST("/", handlers.CreateUser)
	// 	userGroup.GET("/:id", handlers.GetUserByID)
	// 	userGroup.PATCH("/:id", handlers.UpdateUser)
	// 	userGroup.DELETE("/:id", handlers.DeleteUser)
	// 	userGroup.GET("/:id/recommendation", handlers.GetRecommendation)
	// 	userGroup.GET("/:id/fetch", handlers.FetchUserData)
	// }

	// recommendationBaseGroup := router.Group("/recommendation/base")
	// {
	// 	recommendationBaseGroup.POST("/:id", handlers.StartBase)
	// }

	recommendationAdvGroup := router.Group("/recommendation/advanced")
	{
		recommendationAdvGroup.POST("/:id", handlers.StartAdvanced)
	}
}
