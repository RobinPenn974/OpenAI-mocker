package responses

import (
	"strings"

	"RobinPenn974/OpenAI-mocker/templates"
)

// ChatGenerator 普通聊天模型响应生成器
type ChatGenerator struct{}

// NewChatGenerator 创建一个新的聊天响应生成器
func NewChatGenerator() *ChatGenerator {
	return &ChatGenerator{}
}

// GenerateResponse 根据输入生成聊天响应
func (g *ChatGenerator) GenerateResponse(input string, modelID string) ResponseContent {
	// 获取模型的响应模板
	template := templates.GetTemplate(modelID)
	var responseText string

	// 根据输入内容生成响应
	if strings.Contains(strings.ToLower(input), "hello") || strings.Contains(strings.ToLower(input), "hi") {
		responseText = template.Prefix + template.Greeting
	} else if strings.Contains(strings.ToLower(input), "help") {
		responseText = template.Prefix + template.HelpRequest
	} else if strings.Contains(strings.ToLower(input), "?") {
		responseText = template.Prefix + template.Question
	} else {
		responseText = template.Prefix + template.Default
	}

	return ResponseContent{
		Content:          responseText,
		ReasoningContent: nil,
		FinishReason:     "stop",
	}
}
