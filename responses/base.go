package responses

import (
	"strings"
	"time"

	"RobinPenn974/OpenAI-mocker/api"
)

// ResponseGenerator 定义了响应生成器的通用接口
type ResponseGenerator interface {
	// GenerateResponse 根据输入生成响应内容
	GenerateResponse(input string, modelID string) ResponseContent
}

// ReasoningSupport 定义了支持推理功能的响应生成器接口
type ReasoningSupport interface {
	// GenerateReasoningContent 生成推理内容
	GenerateReasoningContent(input string, modelID string) string

	// ShouldUseReasoningField 判断是否应该使用专门的reasoning_content字段
	ShouldUseReasoningField() bool
}

// ResponseContent 包含生成的响应内容
type ResponseContent struct {
	Content          string  // 主要内容
	ReasoningContent *string // 可选的推理内容，如果不支持则为nil
	FinishReason     string  // 结束原因
}

// GenerateID 生成唯一的响应ID
func GenerateID(prefix string) string {
	return prefix + "-" + api.GenerateShortUUID()
}

// GetCurrentTimestamp 获取当前Unix时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// ModelFactory 根据模型ID和类型返回合适的响应生成器
func ModelFactory(modelID string) ResponseGenerator {
	// 特殊处理推理模型
	if modelID == "deepseek-reasoner" {
		return NewReasoningGenerator()
	}

	// 处理文本补全模型
	if strings.Contains(modelID, "davinci") {
		return NewCompletionGenerator()
	}

	// 默认使用普通聊天模型
	return NewChatGenerator()
}
