package responses

import (
	"strings"
)

// CompletionGenerator 文本补全模型响应生成器
type CompletionGenerator struct{}

// NewCompletionGenerator 创建一个新的补全响应生成器
func NewCompletionGenerator() *CompletionGenerator {
	return &CompletionGenerator{}
}

// GenerateResponse 根据输入生成文本补全响应
func (g *CompletionGenerator) GenerateResponse(prompt string, modelID string) ResponseContent {
	prefix := getModelPrefix(modelID)
	var responseText string

	// 根据输入内容构造简单回复
	if strings.Contains(strings.ToLower(prompt), "hello") || strings.Contains(strings.ToLower(prompt), "hi") {
		responseText = prefix + "there! I'm a mock AI model. How can I assist you today?"
	} else if strings.Contains(strings.ToLower(prompt), "help") {
		responseText = prefix + "is on the way! This is a simulated response from the mock completions API."
	} else if strings.Contains(strings.ToLower(prompt), "?") {
		responseText = prefix + "That's an interesting question. As a mock model, I'll provide this simulated answer."
	} else {
		responseText = prefix + "As a mock AI model, I'm continuing your text with this simulated response. In a real OpenAI API, this would be generated based on the trained model."
	}

	return ResponseContent{
		Content:          responseText,
		ReasoningContent: nil,
		FinishReason:     "stop",
	}
}
