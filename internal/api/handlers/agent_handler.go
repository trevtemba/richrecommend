package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trevtemba/richrecommend/internal/agents/orchestrator"
	"github.com/trevtemba/richrecommend/internal/models"
)

func Start(c *gin.Context) {
	var req models.Request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	params := models.OrchestratorParams{
		SystemPrompt:               req.SystemPrompt,
		UserPrompt:                 req.UserPrompt,
		Categories:                 req.Categories,
		RecommendationsPerCategory: req.RecommendationsPerCategory,
		ResponseFormat:             req.ResponseFormat,
	}
	results, err := orchestrator.RunPipelineWithParams(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
