package controller

import (
	"net/http"
	"time"

	"RobinPenn974/OpenAI-mocker/models"

	"github.com/gin-gonic/gin"
)

// ModelLoadRequest 加载模型的请求结构
type ModelLoadRequest struct {
	ModelID   string `json:"model_id" binding:"required"`
	ModelType string `json:"model_type" binding:"required"`
	OwnedBy   string `json:"owned_by,omitempty"`
}

// HandleLoadModel 处理加载模型的请求
func HandleLoadModel(c *gin.Context) {
	var req ModelLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid request: " + err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 创建模型信息
	modelInfo := models.ModelInfo{
		ID:        req.ModelID,
		Object:    "model",
		Created:   time.Now().Unix(),
		OwnedBy:   req.OwnedBy,
		ModelType: req.ModelType,
	}

	// 如果没有提供OwnedBy，设置默认值
	if modelInfo.OwnedBy == "" {
		modelInfo.OwnedBy = "openai-mocker"
	}

	// 注册模型
	models.RegisterModel(modelInfo)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Model loaded successfully",
		"model":   modelInfo,
	})
}

// HandleUnloadModel 处理卸载模型的请求
func HandleUnloadModel(c *gin.Context) {
	var req struct {
		ModelID string `json:"model_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid request: " + err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 卸载模型
	if err := models.UnloadModel(req.ModelID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Model unloaded successfully",
	})
}

// HandleUnloadAllModels 处理卸载所有模型的请求
func HandleUnloadAllModels(c *gin.Context) {
	// 卸载所有模型
	models.UnloadAllModels()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All models unloaded successfully",
	})
}

// HandlePreloadModels 处理预加载默认模型的请求
func HandlePreloadModels(c *gin.Context) {
	// 初始化默认模型
	models.InitDefaultModels()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Default models preloaded successfully",
	})
}
