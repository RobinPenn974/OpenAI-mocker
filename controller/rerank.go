package controller

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"RobinPenn974/OpenAI-mocker/api"
	"RobinPenn974/OpenAI-mocker/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HandleRerank 处理文档重排序请求
func HandleRerank(c *gin.Context) {
	var req api.RerankRequest
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
		modelID = "mock-rerank-v1" // 默认模型
	}

	model, err := models.GetModel(modelID)
	if err != nil || model.ModelType != models.ModelTypeRerank {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code,omitempty"`
			}{
				Message: fmt.Sprintf("Model '%s' not found or not a rerank model", modelID),
				Type:    "model_not_found_error",
			},
		})
		return
	}

	// 生成模拟的重排序结果
	response := generateMockRerank(req)
	c.JSON(http.StatusOK, response)
}

// generateMockRerank 生成模拟的重排序结果
func generateMockRerank(req api.RerankRequest) api.RerankResponse {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 限制topN
	topN := req.TopN
	if topN <= 0 || topN > len(req.Documents) {
		topN = len(req.Documents)
	}

	// 准备结果
	results := make([]api.RerankResult, len(req.Documents))

	// 计算相关性分数 - 基于简单的字符匹配和随机因子
	query := strings.ToLower(req.Query)
	queryWords := strings.Fields(query)

	for i, doc := range req.Documents {
		// 计算基础分数 - 字符匹配越多分数越高
		docLower := strings.ToLower(doc)

		// 计算查询词在文档中出现的次数
		matchScore := 0.0
		for _, word := range queryWords {
			if strings.Contains(docLower, word) {
				matchScore += 0.2 + rand.Float64()*0.1
			}
		}

		// 添加随机因子，保证分数有差异
		randomFactor := 0.2 + rand.Float64()*0.3

		// 计算最终分数，确保范围在0到1之间
		score := math.Min(0.2+matchScore+randomFactor, 0.99)

		results[i] = api.RerankResult{
			Index:          i,
			Document:       doc,
			Score:          score,
			RelevanceScore: score,
		}
	}

	// 按分数降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// 只保留topN个结果
	if topN < len(results) {
		results = results[:topN]
	}

	// 构建响应
	return api.RerankResponse{
		ID:      "rerank-" + uuid.NewString()[:8],
		Object:  "rerank-list",
		Results: results,
		Model:   req.Model,
		Created: time.Now().Unix(),
	}
}
