package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trevtemba/richrecommend/internal/agents/orchestrator"
	"github.com/trevtemba/richrecommend/internal/models"
)

// func StartBase(c *gin.Context) {
// 	var req models.RequestBase
// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
// 		return
// 	}

// 	params := models.OrchestratorParams{
// 		SystemPrompt:               req.SystemPrompt,
// 		UserPrompt:                 req.UserPrompt,
// 		Categories:                 req.Categories,
// 		RecommendationsPerCategory: req.RecommendationsPerCategory,
// 		ContextSchema:              req.ContextSchema,
// 	}
// 	results, err := orchestrator.RunBasePipelineWithParams(params)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, results)
// }

func StartAdvanced(c *gin.Context) {
	var req models.RequestAdvanced

	key := c.GetHeader("X-Provider-Key")

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	val, exists := c.Get("request_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading request ID"})
		return
	}

	requestId, ok := val.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "request ID is not a string"})
		return
	}

	params := models.OrchestratorParams(req)
	results, err := orchestrator.RunAdvPipelineWithParams(params, key, requestId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
