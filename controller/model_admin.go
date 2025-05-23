package controller

import (
	"net/http"
	"time"

	"RobinPenn974/OpenAI-mocker/models"
	"RobinPenn974/OpenAI-mocker/templates"

	"github.com/gin-gonic/gin"
)

// ModelLoadRequest 加载模型的请求结构
type ModelLoadRequest struct {
	ModelID      string `json:"model_id" binding:"required"`
	ModelType    string `json:"model_type" binding:"required"`
	OwnedBy      string `json:"owned_by,omitempty"`
	ResponseType string `json:"response_type,omitempty"` // chat, completion, reasoning

	// 可选的响应模板
	Template *TemplateConfig `json:"template,omitempty"`
}

// TemplateConfig 模型响应模板配置
type TemplateConfig struct {
	Prefix           string `json:"prefix,omitempty"`
	Greeting         string `json:"greeting,omitempty"`
	Question         string `json:"question,omitempty"`
	HelpRequest      string `json:"help_request,omitempty"`
	Default          string `json:"default,omitempty"`
	SupportReasoning bool   `json:"support_reasoning,omitempty"`
	ReasoningPrefix  string `json:"reasoning_prefix,omitempty"`
	CompletionPrefix string `json:"completion_prefix,omitempty"`
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

	// 如果提供了模板，注册模板
	if req.Template != nil {
		template := templates.ResponseTemplate{
			ModelID:          req.ModelID,
			Prefix:           req.Template.Prefix,
			Greeting:         req.Template.Greeting,
			Question:         req.Template.Question,
			HelpRequest:      req.Template.HelpRequest,
			Default:          req.Template.Default,
			SupportReasoning: req.Template.SupportReasoning,
			ReasoningPrefix:  req.Template.ReasoningPrefix,
			CompletionPrefix: req.Template.CompletionPrefix,
		}

		// 设置默认值
		if template.Prefix == "" {
			template.Prefix = "[" + req.ModelID + "] "
		}
		if template.Greeting == "" {
			template.Greeting = "Hello! I'm a mock AI model based on " + req.ModelID + ". How can I assist you today?"
		}
		if template.Question == "" {
			template.Question = "That's an interesting question. As a " + req.ModelID + " mock, I'll provide a simulated answer."
		}
		if template.HelpRequest == "" {
			template.HelpRequest = "I'm here to help! As a " + req.ModelID + " mock, I can provide simulated assistance."
		}
		if template.Default == "" {
			template.Default = "I understand. As a " + req.ModelID + " mock, I'm providing this simulated response."
		}
		if template.ReasoningPrefix == "" && template.SupportReasoning {
			template.ReasoningPrefix = "REASONING: "
		}

		// 注册模板
		if err := templates.RegisterTemplate(template); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message": "Failed to register template: " + err.Error(),
					"type":    "internal_server_error",
				},
			})
			return
		}
	}

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

	// 删除模板
	templates.DeleteTemplate(req.ModelID)

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
