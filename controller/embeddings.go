package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"RobinPenn974/OpenAI-mocker/api"
	"RobinPenn974/OpenAI-mocker/models"

	"github.com/gin-gonic/gin"
)

// HandleEmbeddings 处理Embeddings请求
func HandleEmbeddings(c *gin.Context) {
	var req api.EmbeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code,omitempty"`
			}{
				Message: "Invalid request: " + err.Error(),
				Type:    "invalid_request_error",
			},
		})
		return
	}

	// 检查模型是否存在
	modelID := req.Model
	if modelID == "" {
		modelID = "mock-embedding-ada-002" // 默认模型
	}

	model, err := models.GetModel(modelID)
	if err != nil || model.ModelType != models.ModelTypeEmbedding {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code,omitempty"`
			}{
				Message: fmt.Sprintf("Model '%s' not found or not an embedding model", modelID),
				Type:    "model_not_found_error",
			},
		})
		return
	}

	// 生成模拟嵌入向量
	response := generateMockEmbeddings(req)
	c.JSON(http.StatusOK, response)
}

// generateMockEmbeddings 生成模拟的嵌入向量
func generateMockEmbeddings(req api.EmbeddingRequest) api.EmbeddingResponse {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 为每个输入生成模拟嵌入向量
	data := make([]api.EmbeddingData, 0, len(req.Input))
	totalTokens := 0

	for i, input := range req.Input {
		// 生成1536维的随机向量并标准化
		embedding := make([]float64, 1536)
		var sum float64
		for j := range embedding {
			embedding[j] = rand.Float64() - 0.5
			sum += embedding[j] * embedding[j]
		}

		// 标准化向量
		magnitude := 0.0
		if sum > 0 {
			magnitude = 1.0 / (sum * 0.5)
		}
		for j := range embedding {
			embedding[j] *= magnitude
		}

		// 添加到结果集
		data = append(data, api.EmbeddingData{
			Object:    "embedding",
			Embedding: embedding,
			Index:     i,
		})

		// 估算token数量
		tokens := len([]rune(input)) / 4
		if tokens < 1 {
			tokens = 1
		}
		totalTokens += tokens
	}

	// 构建响应
	return api.EmbeddingResponse{
		Object: "list",
		Data:   data,
		Model:  req.Model,
		Usage: struct {
			PromptTokens int `json:"prompt_tokens"`
			TotalTokens  int `json:"total_tokens"`
		}{
			PromptTokens: totalTokens,
			TotalTokens:  totalTokens,
		},
	}
}
