package controller

import (
	"net/http"

	"RobinPenn974/OpenAI-mocker/templates"

	"github.com/gin-gonic/gin"
)

// HandleListTemplates 处理列出所有模板的请求
func HandleListTemplates(c *gin.Context) {
	templateList := templates.ListTemplates()
	c.JSON(http.StatusOK, gin.H{
		"templates": templateList,
		"count":     len(templateList),
	})
}

// HandleGetTemplate 处理获取指定模型模板的请求
func HandleGetTemplate(c *gin.Context) {
	modelID := c.Param("model_id")
	if modelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Model ID is required",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	template := templates.GetTemplate(modelID)
	c.JSON(http.StatusOK, template)
}

// HandleUpdateTemplate 处理更新模板的请求
func HandleUpdateTemplate(c *gin.Context) {
	var template templates.ResponseTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid request: " + err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	modelID := c.Param("model_id")
	if modelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Model ID is required",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 确保ModelID匹配
	template.ModelID = modelID

	// 注册模板
	if err := templates.RegisterTemplate(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to update template: " + err.Error(),
				"type":    "internal_server_error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Template updated successfully",
		"template": template,
	})
}

// HandleDeleteTemplate 处理删除模板的请求
func HandleDeleteTemplate(c *gin.Context) {
	modelID := c.Param("model_id")
	if modelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Model ID is required",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 删除模板
	if err := templates.DeleteTemplate(modelID); err != nil {
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
		"message": "Template deleted successfully",
	})
}
