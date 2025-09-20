package main

import (
	"net/http"

	"github.com/trevtemba/richrecommend/models"
	"github.com/trevtemba/richrecommend/orchestrator"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/recommend", func(c *gin.Context) {
		var req models.RecommendationRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		params := models.OrchestratorParams{
			SystemPrompt:               req.SystemPrompt,
			UserPrompt:                 req.UserPrompt,
			Categories:                 req.Categories,
			RecommendationsPerCategory: req.RecommendationsPerCategory,
			RecommendationSchema:       req.RecommendationSchema,
		}
		results, err := orchestrator.RunPipelineWithParams(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	})

	r.Run(":8080") // listens on localhost:8080
}
