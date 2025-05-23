package responses

import (
	"strings"

	"RobinPenn974/OpenAI-mocker/templates"
)

// CompletionGenerator 文本补全模型响应生成器
type CompletionGenerator struct{}

// NewCompletionGenerator 创建一个新的补全响应生成器
func NewCompletionGenerator() *CompletionGenerator {
	return &CompletionGenerator{}
}

// GenerateResponse 根据输入生成文本补全响应
func (g *CompletionGenerator) GenerateResponse(prompt string, modelID string) ResponseContent {
	// 获取模型的响应模板
	template := templates.GetTemplate(modelID)
	var responseText string

	// 根据输入内容构造简单回复
	if strings.Contains(strings.ToLower(prompt), "hello") || strings.Contains(strings.ToLower(prompt), "hi") {
		responseText = template.Prefix + template.Greeting
	} else if strings.Contains(strings.ToLower(prompt), "help") {
		responseText = template.Prefix + template.HelpRequest
	} else if strings.Contains(strings.ToLower(prompt), "?") {
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
