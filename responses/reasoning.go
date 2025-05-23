package responses

import (
	"os"
	"strings"

	"RobinPenn974/OpenAI-mocker/templates"
)

// ReasoningGenerator 推理模型响应生成器
type ReasoningGenerator struct{}

// NewReasoningGenerator 创建一个新的推理响应生成器
func NewReasoningGenerator() *ReasoningGenerator {
	return &ReasoningGenerator{}
}

// GenerateResponse 根据输入生成推理模型的响应
func (g *ReasoningGenerator) GenerateResponse(input string, modelID string) ResponseContent {
	// 获取模型的响应模板
	template := templates.GetTemplate(modelID)

	// 生成推理内容
	reasoningContent := g.GenerateReasoningContent(input, modelID)
	reasoningContent = template.Prefix + template.ReasoningPrefix + reasoningContent

	// 生成常规回复内容
	var responseText string
	if strings.Contains(strings.ToLower(input), "hello") || strings.Contains(strings.ToLower(input), "hi") {
		responseText = template.Prefix + template.Greeting
	} else if strings.Contains(strings.ToLower(input), "help") {
		responseText = template.Prefix + template.HelpRequest
	} else if strings.Contains(strings.ToLower(input), "?") {
		responseText = template.Prefix + template.Question
	} else {
		responseText = template.Prefix + template.Default
	}

	// 根据环境变量决定如何处理推理内容
	if g.ShouldUseReasoningField() {
		// 启用了推理功能，单独返回推理内容和回复内容
		reasoningPtr := &reasoningContent
		return ResponseContent{
			Content:          responseText,
			ReasoningContent: reasoningPtr,
			FinishReason:     "stop",
		}
	} else {
		// 未启用推理功能，将推理内容合并到回复内容中
		combinedContent := "<think>" + reasoningContent + "</think>\n\n" + responseText
		return ResponseContent{
			Content:          combinedContent,
			ReasoningContent: nil,
			FinishReason:     "stop",
		}
	}
}

// GenerateReasoningContent 生成推理内容
func (g *ReasoningGenerator) GenerateReasoningContent(question string, modelID string) string {
	// 根据不同的问题输入生成不同的推理内容
	return "Let me think step by step about this question.\n\n" +
		"First, I need to understand what is being asked:\n" +
		"The user asked about: " + question + "\n\n" +
		"Now I will analyze this by breaking it down:\n" +
		"1. Identify key information\n" +
		"2. Apply relevant knowledge\n" +
		"3. Consider different angles\n" +
		"4. Form a logical conclusion\n\n" +
		"Based on my analysis, I can now provide a comprehensive response."
}

// ShouldUseReasoningField 判断是否应该使用专门的reasoning_content字段
func (g *ReasoningGenerator) ShouldUseReasoningField() bool {
	val := os.Getenv("ENABLE_REASONING")
	return strings.ToLower(val) == "true" || val == "1"
}
