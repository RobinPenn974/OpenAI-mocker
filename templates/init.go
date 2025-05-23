package templates

// DefaultManager 全局默认模板管理器
var DefaultManager = NewTemplateManager("", "")

// 为了方便访问，提供一些全局方法

// GetTemplate 获取指定模型的模板
func GetTemplate(modelID string) ResponseTemplate {
	return DefaultManager.GetTemplate(modelID)
}

// RegisterTemplate 注册或更新模型的响应模板
func RegisterTemplate(template ResponseTemplate) error {
	return DefaultManager.RegisterTemplate(template)
}

// DeleteTemplate 删除模型的响应模板
func DeleteTemplate(modelID string) error {
	return DefaultManager.DeleteTemplate(modelID)
}

// ListTemplates 获取所有模板的列表
func ListTemplates() []ResponseTemplate {
	return DefaultManager.ListTemplates()
}
