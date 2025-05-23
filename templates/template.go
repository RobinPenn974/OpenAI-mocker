package templates

// ResponseTemplate 定义了一个模型的响应模板
type ResponseTemplate struct {
	// 基本信息
	ModelID string `json:"model_id"` // 模型ID
	Prefix  string `json:"prefix"`   // 模型前缀标识

	// 通用响应模板
	Greeting    string `json:"greeting"`     // 问候语模板
	Question    string `json:"question"`     // 问题回答模板
	HelpRequest string `json:"help_request"` // 帮助请求模板
	Default     string `json:"default"`      // 默认回复模板

	// 推理模型配置
	SupportReasoning  bool   `json:"support_reasoning"`  // 是否支持推理功能
	ReasoningPrefix   string `json:"reasoning_prefix"`   // 推理内容前缀
	ReasoningTemplate string `json:"reasoning_template"` // 推理内容模板

	// 文本补全模型配置
	CompletionPrefix string `json:"completion_prefix"` // 补全前缀
}
