package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"RobinPenn974/OpenAI-mocker/api"
	"RobinPenn974/OpenAI-mocker/models"
	"RobinPenn974/OpenAI-mocker/responses"

	"github.com/gin-gonic/gin"
)

// HandleCompletions 处理文本完成请求
func HandleCompletions(c *gin.Context) {
	var req api.CompletionRequest
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
		modelID = "mock-davinci-002" // 默认模型
	}

	_, err := models.GetModel(modelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code,omitempty"`
			}{
				Message: fmt.Sprintf("Model '%s' not found", modelID),
				Type:    "model_not_found_error",
			},
		})
		return
	}

	// 根据Stream参数决定响应方式
	if req.Stream {
		handleStreamingCompletion(c, req)
	} else {
		// 生成模拟回复
		response := generateCompletion(req)
		c.JSON(http.StatusOK, response)
	}
}

// handleStreamingCompletion 处理流式返回
func handleStreamingCompletion(c *gin.Context, req api.CompletionRequest) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 生成唯一ID
	responseID := responses.GenerateID("cmpl")
	now := responses.GetCurrentTimestamp()

	// 获取响应生成器
	generator := responses.ModelFactory(req.Model)

	// 生成响应内容
	responseContent := generator.GenerateResponse(req.Prompt, req.Model)

	// 流式发送内容
	words := strings.Split(responseContent.Content, " ")
	chunkSize := 2 // 每次发送2个词

	// 发送所有块，除了最后一块
	for i := 0; i < len(words)-chunkSize; i += chunkSize {
		chunk := strings.Join(words[i:i+chunkSize], " ")
		if i == 0 {
			chunk = strings.TrimSpace(chunk)
		}

		var finishReason *string

		chunkResponse := api.CompletionChunkResponse{
			ID:      responseID,
			Object:  "text_completion",
			Created: now,
			Model:   req.Model,
			Choices: []api.CompletionChunkChoice{
				{
					Index:        0,
					Text:         chunk,
					FinishReason: finishReason,
				},
			},
		}

		// 发送数据
		c.SSEvent("", chunkResponse)
		c.Writer.Flush()
		time.Sleep(100 * time.Millisecond) // 模拟流式生成延迟
	}

	// 发送最后一块
	if len(words) > 0 {
		lastChunk := strings.Join(words[max(0, len(words)-chunkSize):], " ")
		finishReason := responseContent.FinishReason

		chunkResponse := api.CompletionChunkResponse{
			ID:      responseID,
			Object:  "text_completion",
			Created: now,
			Model:   req.Model,
			Choices: []api.CompletionChunkChoice{
				{
					Index:        0,
					Text:         lastChunk,
					FinishReason: &finishReason,
				},
			},
		}

		c.SSEvent("", chunkResponse)
		c.Writer.Flush()
	}

	// 发送结束事件
	c.SSEvent("", "")
}

// generateCompletion 生成模拟的文本完成回复
func generateCompletion(req api.CompletionRequest) api.CompletionResponse {
	// 获取响应生成器
	generator := responses.ModelFactory(req.Model)

	// 生成响应内容
	responseContent := generator.GenerateResponse(req.Prompt, req.Model)

	// 创建模拟的Token使用量
	promptTokens := len(strings.Split(req.Prompt, " "))
	completionTokens := len(strings.Split(responseContent.Content, " "))
	totalTokens := promptTokens + completionTokens

	// 构建响应
	now := responses.GetCurrentTimestamp()
	response := api.CompletionResponse{
		ID:      responses.GenerateID("cmpl"),
		Object:  "text_completion",
		Created: now,
		Model:   req.Model,
		Choices: []api.CompletionChoice{
			{
				Text:         responseContent.Content,
				Index:        0,
				FinishReason: responseContent.FinishReason,
			},
		},
		Usage: api.ChatCompletionUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
		},
	}

	return response
}
