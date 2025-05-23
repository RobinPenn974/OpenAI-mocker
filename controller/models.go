package controller

import (
	"net/http"

	"RobinPenn974/OpenAI-mocker/api"
	"RobinPenn974/OpenAI-mocker/models"

	"github.com/gin-gonic/gin"
)

// HandleListModels 处理获取模型列表请求
func HandleListModels(c *gin.Context) {
	modelsList := models.ListModels()

	// 转换为API响应格式
	response := api.ModelListResponse{
		Object: "list",
		Data:   make([]api.ModelData, 0, len(modelsList)),
	}

	for _, model := range modelsList {
		response.Data = append(response.Data, api.ModelData{
			ID:      model.ID,
			Object:  model.Object,
			Created: model.Created,
			OwnedBy: model.OwnedBy,
		})
	}

	c.JSON(http.StatusOK, response)
}
