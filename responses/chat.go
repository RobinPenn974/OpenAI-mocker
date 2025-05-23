package responses

import (
	"strings"
)

// ChatGenerator 普通聊天模型响应生成器
type ChatGenerator struct{}

// NewChatGenerator 创建一个新的聊天响应生成器
func NewChatGenerator() *ChatGenerator {
	return &ChatGenerator{}
}

// GenerateResponse 根据输入生成聊天响应
func (g *ChatGenerator) GenerateResponse(input string, modelID string) ResponseContent {
	prefix := getModelPrefix(modelID)
	var responseText string

	// 根据输入内容生成响应
	if strings.Contains(strings.ToLower(input), "hello") || strings.Contains(strings.ToLower(input), "hi") {
		responseText = prefix + "Hello! I'm a mock GPT model. How can I assist you today?"
	} else if strings.Contains(strings.ToLower(input), "help") {
		responseText = prefix + "I'm here to help! Although I'm just a mock model, I can simulate responses. What do you need assistance with?"
	} else if strings.Contains(strings.ToLower(input), "?") {
		responseText = prefix + "That's an interesting question. As a mock model, I'll respond with this simulated answer. In a real OpenAI API, you would get a more contextual response."
	} else {
		responseText = prefix + "I understand. As a mock GPT model, I'm providing this simulated response to your message. In a real OpenAI API, the response would be generated based on the trained model."
	}

	return ResponseContent{
		Content:          responseText,
		ReasoningContent: nil,
		FinishReason:     "stop",
	}
}

// getModelPrefix 根据模型ID生成标识前缀
func getModelPrefix(modelID string) string {
	// 提取模型名称的简短形式作为前缀
	parts := strings.Split(modelID, "-")
	prefix := "[MOCK]"

	if len(parts) > 1 {
		// 如果是形如 mock-gpt-3.5-turbo 的格式，提取关键部分
		if len(parts) >= 3 && parts[0] == "mock" {
			prefix = "[" + strings.ToUpper(parts[1]) + "-" + parts[2] + "]"
		} else {
			// 否则使用最后一个部分
			prefix = "[" + strings.ToUpper(parts[len(parts)-1]) + "]"
		}
	}

	return prefix + " "
}
