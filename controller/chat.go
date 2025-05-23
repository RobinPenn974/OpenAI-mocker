package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"RobinPenn974/OpenAI-mocker/api"
	"RobinPenn974/OpenAI-mocker/models"
	"RobinPenn974/OpenAI-mocker/responses"

	"github.com/gin-gonic/gin"
)

// HandleChatCompletions 处理Chat Completions请求
func HandleChatCompletions(c *gin.Context) {
	var req api.ChatCompletionRequest
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
		modelID = "mock-gpt-3.5-turbo" // 默认模型
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
		handleStreamingChatCompletion(c, req)
	} else {
		// 生成模型响应
		response := generateChatResponse(req)
		c.JSON(http.StatusOK, response)
	}
}

// handleStreamingChatCompletion 处理流式聊天完成请求
func handleStreamingChatCompletion(c *gin.Context, req api.ChatCompletionRequest) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 生成唯一ID
	responseID := responses.GenerateID("chatcmpl")
	now := responses.GetCurrentTimestamp()

	// 获取最后一条消息内容以便生成相关回复
	var lastContent string
	if len(req.Messages) > 0 {
		lastContent = req.Messages[len(req.Messages)-1].Content
	}

	// 获取响应生成器
	generator := responses.ModelFactory(req.Model)

	// 生成响应内容
	responseContent := generator.GenerateResponse(lastContent, req.Model)

	// 首先发送role
	roleChunk := api.ChatCompletionChunkResponse{
		ID:      responseID,
		Object:  "chat.completion.chunk",
		Created: now,
		Model:   req.Model,
		Choices: []api.ChatCompletionChunkChoice{
			{
				Index: 0,
				Delta: api.ChatCompletionChunkDelta{
					Role: stringPtr("assistant"),
				},
			},
		},
	}

	c.SSEvent("", roleChunk)
	c.Writer.Flush()
	time.Sleep(50 * time.Millisecond)

	// 流式处理推理内容（如果有）
	if responseContent.ReasoningContent != nil {
		// 启用了推理功能，将推理内容作为 reasoning_content 字段流式返回
		reasoningWords := strings.Split(*responseContent.ReasoningContent, " ")
		reasoningChunkSize := 3 // 每次发送3个词

		// 发送推理内容的所有块
		for i := 0; i < len(reasoningWords); i += reasoningChunkSize {
			end := min(i+reasoningChunkSize, len(reasoningWords))
			chunk := strings.Join(reasoningWords[i:end], " ")
			if i == 0 {
				chunk = strings.TrimSpace(chunk)
			}

			chunkResponse := api.ChatCompletionChunkResponse{
				ID:      responseID,
				Object:  "chat.completion.chunk",
				Created: now,
				Model:   req.Model,
				Choices: []api.ChatCompletionChunkChoice{
					{
						Index: 0,
						Delta: api.ChatCompletionChunkDelta{
							ReasoningContent: &chunk,
						},
					},
				},
			}

			c.SSEvent("", chunkResponse)
			c.Writer.Flush()
			time.Sleep(50 * time.Millisecond)
		}
	}

	// 解析并流式发送主要内容
	content := responseContent.Content

	// 检查是否包含思考内容（<think>标签）
	if strings.Contains(content, "<think>") && strings.Contains(content, "</think>") {
		// 提取思考内容和实际回答
		parts := strings.SplitN(content, "</think>", 2)
		if len(parts) == 2 {
			// 提取思考部分（移除<think>标签）
			thinkingPart := strings.TrimPrefix(parts[0], "<think>")
			// 提取回答部分（去除前导空白）
			answerPart := strings.TrimSpace(parts[1])

			// 流式发送思考部分
			thinkingWords := strings.Split(thinkingPart, " ")
			chunkSize := 3 // 思考部分每次发送3个词

			// 先发送思考开始标记
			startThinkTag := "<think>"
			chunkResponse := api.ChatCompletionChunkResponse{
				ID:      responseID,
				Object:  "chat.completion.chunk",
				Created: now,
				Model:   req.Model,
				Choices: []api.ChatCompletionChunkChoice{
					{
						Index: 0,
						Delta: api.ChatCompletionChunkDelta{
							Content: &startThinkTag,
						},
					},
				},
			}
			c.SSEvent("", chunkResponse)
			c.Writer.Flush()
			time.Sleep(50 * time.Millisecond)

			// 发送思考内容
			for i := 0; i < len(thinkingWords); i += chunkSize {
				end := min(i+chunkSize, len(thinkingWords))
				chunk := strings.Join(thinkingWords[i:end], " ")
				if i == 0 {
					chunk = strings.TrimSpace(chunk)
				}

				chunkContent := chunk
				chunkResponse := api.ChatCompletionChunkResponse{
					ID:      responseID,
					Object:  "chat.completion.chunk",
					Created: now,
					Model:   req.Model,
					Choices: []api.ChatCompletionChunkChoice{
						{
							Index: 0,
							Delta: api.ChatCompletionChunkDelta{
								Content: &chunkContent,
							},
						},
					},
				}

				c.SSEvent("", chunkResponse)
				c.Writer.Flush()
				time.Sleep(50 * time.Millisecond)
			}

			// 发送思考结束标记
			endThinkTag := "</think>\n\n"
			chunkResponse = api.ChatCompletionChunkResponse{
				ID:      responseID,
				Object:  "chat.completion.chunk",
				Created: now,
				Model:   req.Model,
				Choices: []api.ChatCompletionChunkChoice{
					{
						Index: 0,
						Delta: api.ChatCompletionChunkDelta{
							Content: &endThinkTag,
						},
					},
				},
			}
			c.SSEvent("", chunkResponse)
			c.Writer.Flush()
			time.Sleep(50 * time.Millisecond)

			// 接着流式发送回答部分
			content = answerPart
		}
	}

	// 流式发送内容（思考部分处理完后的回答部分，或原始内容）
	// 检测是否包含中文字符，或者内容较长但空格很少（通常是中文内容的特征）
	containsChinese := false
	for _, r := range content {
		if unicode.Is(unicode.Han, r) {
			containsChinese = true
			break
		}
	}

	words := strings.Split(content, " ")

	// 判断是否应该使用按字符分割的方式（适用于中文内容）
	if containsChinese || (len(words) <= 5 && len(content) > 15) {
		// 按字符分割，处理中文内容
		runes := []rune(content)
		runeChunkSize := 5 // 每次发送5个字符

		// 发送除最后一块外的所有内容
		for i := 0; i < len(runes); i += runeChunkSize {
			end := min(i+runeChunkSize, len(runes))
			if end == len(runes) && i+runeChunkSize >= len(runes) {
				// 最后一块将在后面处理
				break
			}

			chunk := string(runes[i:end])
			chunkContent := chunk
			chunkResponse := api.ChatCompletionChunkResponse{
				ID:      responseID,
				Object:  "chat.completion.chunk",
				Created: now,
				Model:   req.Model,
				Choices: []api.ChatCompletionChunkChoice{
					{
						Index: 0,
						Delta: api.ChatCompletionChunkDelta{
							Content: &chunkContent,
						},
					},
				},
			}

			c.SSEvent("", chunkResponse)
			c.Writer.Flush()
			time.Sleep(50 * time.Millisecond)
		}

		// 发送最后一块
		if len(runes) > 0 {
			lastIndex := (len(runes) / runeChunkSize) * runeChunkSize
			if lastIndex < len(runes) {
				lastChunk := string(runes[lastIndex:])
				lastContent := lastChunk
				finishReason := responseContent.FinishReason

				chunkResponse := api.ChatCompletionChunkResponse{
					ID:      responseID,
					Object:  "chat.completion.chunk",
					Created: now,
					Model:   req.Model,
					Choices: []api.ChatCompletionChunkChoice{
						{
							Index: 0,
							Delta: api.ChatCompletionChunkDelta{
								Content: &lastContent,
							},
							FinishReason: &finishReason,
						},
					},
				}

				c.SSEvent("", chunkResponse)
				c.Writer.Flush()
			} else {
				// 如果没有剩余内容，也需要发送完成状态
				finishReason := responseContent.FinishReason
				emptyContent := ""

				chunkResponse := api.ChatCompletionChunkResponse{
					ID:      responseID,
					Object:  "chat.completion.chunk",
					Created: now,
					Model:   req.Model,
					Choices: []api.ChatCompletionChunkChoice{
						{
							Index: 0,
							Delta: api.ChatCompletionChunkDelta{
								Content: &emptyContent,
							},
							FinishReason: &finishReason,
						},
					},
				}

				c.SSEvent("", chunkResponse)
				c.Writer.Flush()
			}
		}
	} else {
		// 英文内容，按词分割处理
		chunkSize := 2 // 每次发送2个词

		// 发送所有响应内容块，除了最后一块
		for i := 0; i < len(words)-chunkSize; i += chunkSize {
			chunk := strings.Join(words[i:i+chunkSize], " ")
			if i == 0 {
				chunk = strings.TrimSpace(chunk)
			}

			chunkContent := chunk
			chunkResponse := api.ChatCompletionChunkResponse{
				ID:      responseID,
				Object:  "chat.completion.chunk",
				Created: now,
				Model:   req.Model,
				Choices: []api.ChatCompletionChunkChoice{
					{
						Index: 0,
						Delta: api.ChatCompletionChunkDelta{
							Content: &chunkContent,
						},
					},
				},
			}

			c.SSEvent("", chunkResponse)
			c.Writer.Flush()
			time.Sleep(100 * time.Millisecond)
		}

		// 发送最后一块响应内容
		if len(words) > 0 {
			lastChunk := strings.Join(words[max(0, len(words)-chunkSize):], " ")
			lastContent := lastChunk
			finishReason := responseContent.FinishReason

			chunkResponse := api.ChatCompletionChunkResponse{
				ID:      responseID,
				Object:  "chat.completion.chunk",
				Created: now,
				Model:   req.Model,
				Choices: []api.ChatCompletionChunkChoice{
					{
						Index: 0,
						Delta: api.ChatCompletionChunkDelta{
							Content: &lastContent,
						},
						FinishReason: &finishReason,
					},
				},
			}

			c.SSEvent("", chunkResponse)
			c.Writer.Flush()
		}
	}

	// 发送结束事件
	c.SSEvent("", "")
}

// generateChatResponse 生成模拟的Chat回复
func generateChatResponse(req api.ChatCompletionRequest) api.ChatCompletionResponse {
	// 获取最后一条消息内容以便生成相关回复
	var lastContent string
	if len(req.Messages) > 0 {
		lastContent = req.Messages[len(req.Messages)-1].Content
	}

	// 获取响应生成器
	generator := responses.ModelFactory(req.Model)

	// 生成响应内容
	responseContent := generator.GenerateResponse(lastContent, req.Model)

	// 创建模拟的Token使用量
	promptTokens := len(strings.Split(lastContent, " ")) + 10
	completionTokens := len(strings.Split(responseContent.Content, " "))
	totalTokens := promptTokens + completionTokens

	// 构建响应
	now := responses.GetCurrentTimestamp()
	return api.ChatCompletionResponse{
		ID:      responses.GenerateID("chatcmpl"),
		Object:  "chat.completion",
		Created: now,
		Model:   req.Model,
		Choices: []api.ChatCompletionChoice{
			{
				Index: 0,
				Message: api.ChatCompletionMessage{
					Role:             "assistant",
					Content:          responseContent.Content,
					ReasoningContent: responseContent.ReasoningContent,
				},
				FinishReason: responseContent.FinishReason,
			},
		},
		Usage: api.ChatCompletionUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
		},
	}
}
