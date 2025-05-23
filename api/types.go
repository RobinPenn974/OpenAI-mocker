package api

import "time"

// 通用错误响应格式
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}

// Chat相关类型定义
type ChatCompletionMessage struct {
	Role             string  `json:"role"`
	Content          string  `json:"content"`
	Name             *string `json:"name,omitempty"`
	ReasoningContent *string `json:"reasoning_content,omitempty"`
}

type ChatCompletionRequest struct {
	Model       string                  `json:"model"`
	Messages    []ChatCompletionMessage `json:"messages"`
	Temperature float64                 `json:"temperature,omitempty"`
	MaxTokens   int                     `json:"max_tokens,omitempty"`
	Stream      bool                    `json:"stream,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   ChatCompletionUsage    `json:"usage"`
}

// 添加Chat流式响应类型
type ChatCompletionChunkDelta struct {
	Role             *string `json:"role,omitempty"`
	Content          *string `json:"content,omitempty"`
	ReasoningContent *string `json:"reasoning_content,omitempty"`
}

type ChatCompletionChunkChoice struct {
	Index        int                      `json:"index"`
	Delta        ChatCompletionChunkDelta `json:"delta"`
	FinishReason *string                  `json:"finish_reason,omitempty"`
}

type ChatCompletionChunkResponse struct {
	ID      string                      `json:"id"`
	Object  string                      `json:"object"`
	Created int64                       `json:"created"`
	Model   string                      `json:"model"`
	Choices []ChatCompletionChunkChoice `json:"choices"`
}

// Completion相关类型定义
type CompletionRequest struct {
	Model       string   `json:"model"`
	Prompt      string   `json:"prompt"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
	Temperature float64  `json:"temperature,omitempty"`
	Stop        []string `json:"stop,omitempty"`
	Stream      bool     `json:"stream,omitempty"`
}

type CompletionChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type CompletionResponse struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created int64               `json:"created"`
	Model   string              `json:"model"`
	Choices []CompletionChoice  `json:"choices"`
	Usage   ChatCompletionUsage `json:"usage"`
}

// 添加流式响应类型
type CompletionChunkChoice struct {
	Index        int     `json:"index"`
	Text         string  `json:"text"`
	FinishReason *string `json:"finish_reason,omitempty"`
}

type CompletionChunkResponse struct {
	ID      string                  `json:"id"`
	Object  string                  `json:"object"`
	Created int64                   `json:"created"`
	Model   string                  `json:"model"`
	Choices []CompletionChunkChoice `json:"choices"`
}

// Embedding相关类型定义
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Rerank相关类型定义
type RerankRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopN      int      `json:"top_n,omitempty"`
}

type RerankResult struct {
	Index          int     `json:"index"`
	Document       string  `json:"document"`
	Score          float64 `json:"score"`
	RelevanceScore float64 `json:"relevance_score"`
}

type RerankResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Results []RerankResult `json:"results"`
	Model   string         `json:"model"`
	Created int64          `json:"created"`
}

// Models API相关类型定义
type ModelListResponse struct {
	Object string      `json:"object"`
	Data   []ModelData `json:"data"`
}

type ModelData struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// API密钥管理相关
type ApiKey struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// 模型管理API的请求类型
type LoadModelRequest struct {
	ModelID string `json:"model_id" binding:"required"`
	Path    string `json:"path,omitempty"`
}

type CreateKeyRequest struct {
	Name string `json:"name" binding:"required"`
}
