package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const (
	// 默认模板存储目录
	defaultTemplateDir = "template_data"
	// 默认模板文件名
	defaultTemplateFile = "templates.json"
	// 默认模板数据文件
	defaultTemplateDataFile = "default_templates.json"
)

// TemplateManager 管理模型响应模板，提供持久化存储能力
type TemplateManager struct {
	templates  map[string]ResponseTemplate // 模板存储
	mu         sync.RWMutex                // 读写锁
	storageDir string                      // 存储目录
	filename   string                      // 存储文件名
}

// NewTemplateManager 创建一个新的模板管理器
func NewTemplateManager(storageDir, filename string) *TemplateManager {
	// 如果未指定存储目录，使用默认目录
	if storageDir == "" {
		storageDir = defaultTemplateDir
	}

	// 如果未指定文件名，使用默认文件名
	if filename == "" {
		filename = defaultTemplateFile
	}

	manager := &TemplateManager{
		templates:  make(map[string]ResponseTemplate),
		storageDir: storageDir,
		filename:   filename,
	}

	// 确保存储目录存在
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		fmt.Printf("Error creating template directory: %v\n", err)
	}

	// 加载模板数据
	manager.loadTemplates()

	return manager
}

// GetTemplate 获取指定模型的模板，如果不存在则返回默认模板
func (tm *TemplateManager) GetTemplate(modelID string) ResponseTemplate {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	template, exists := tm.templates[modelID]
	if !exists {
		// 返回默认模板
		return getDefaultTemplate(modelID)
	}
	return template
}

// RegisterTemplate 注册或更新模型的响应模板
func (tm *TemplateManager) RegisterTemplate(template ResponseTemplate) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 从文件中读取当前所有模板
	templates, err := tm.readTemplatesFromFile()
	if err != nil {
		return err
	}

	// 更新或添加模板
	templates[template.ModelID] = template

	// 保存所有模板回文件
	err = tm.writeTemplatesToFile(templates)
	if err != nil {
		return err
	}

	// 更新内存中的模板
	tm.templates[template.ModelID] = template

	return nil
}

// DeleteTemplate 删除模型的响应模板
func (tm *TemplateManager) DeleteTemplate(modelID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 从文件中读取当前所有模板
	templates, err := tm.readTemplatesFromFile()
	if err != nil {
		return err
	}

	// 检查模板是否存在
	if _, exists := templates[modelID]; !exists {
		return fmt.Errorf("template for model %s not found", modelID)
	}

	// 删除模板
	delete(templates, modelID)

	// 保存更改回文件
	err = tm.writeTemplatesToFile(templates)
	if err != nil {
		return err
	}

	// 更新内存中的模板
	delete(tm.templates, modelID)

	return nil
}

// ListTemplates 获取所有模板的列表
func (tm *TemplateManager) ListTemplates() []ResponseTemplate {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	templates := make([]ResponseTemplate, 0, len(tm.templates))
	for _, template := range tm.templates {
		templates = append(templates, template)
	}

	return templates
}

// loadTemplates 从文件加载模板数据
func (tm *TemplateManager) loadTemplates() error {
	filePath := filepath.Join(tm.storageDir, tm.filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，加载默认模板
		return tm.loadDefaultTemplates()
	}

	// 读取文件内容
	templates, err := tm.readTemplatesFromFile()
	if err != nil {
		return err
	}

	// 更新模板存储
	tm.templates = templates

	return nil
}

// readTemplatesFromFile 从文件中读取模板数据
func (tm *TemplateManager) readTemplatesFromFile() (map[string]ResponseTemplate, error) {
	filePath := filepath.Join(tm.storageDir, tm.filename)

	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading template file: %v", err)
	}

	// 解析JSON
	var templates map[string]ResponseTemplate
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, fmt.Errorf("error parsing template file: %v", err)
	}

	return templates, nil
}

// writeTemplatesToFile 将模板数据写入文件
func (tm *TemplateManager) writeTemplatesToFile(templates map[string]ResponseTemplate) error {
	// 将模板数据编码为JSON
	data, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding templates: %v", err)
	}

	// 保存到文件
	filePath := filepath.Join(tm.storageDir, tm.filename)
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing template file: %v", err)
	}

	return nil
}

// loadDefaultTemplates 加载默认模板，从默认模板文件加载，文件不存在则创建默认模板文件
func (tm *TemplateManager) loadDefaultTemplates() error {
	defaultTemplatesPath := filepath.Join(tm.storageDir, defaultTemplateDataFile)

	// 检查默认模板文件是否存在
	if _, err := os.Stat(defaultTemplatesPath); os.IsNotExist(err) {
		// 默认模板文件不存在，创建默认模板文件
		if err := createDefaultTemplatesFile(defaultTemplatesPath); err != nil {
			return fmt.Errorf("error creating default templates file: %v", err)
		}
	}

	// 从默认模板文件加载
	data, err := ioutil.ReadFile(defaultTemplatesPath)
	if err != nil {
		return fmt.Errorf("error reading default template file: %v", err)
	}

	// 解析JSON
	var templates map[string]ResponseTemplate
	if err := json.Unmarshal(data, &templates); err != nil {
		return fmt.Errorf("error parsing default template file: %v", err)
	}

	// 更新模板存储
	tm.templates = templates

	// 将默认模板复制到主模板文件
	return tm.writeTemplatesToFile(templates)
}

// createDefaultTemplatesFile 创建默认模板文件
func createDefaultTemplatesFile(filePath string) error {
	// 创建默认模板
	templates := map[string]ResponseTemplate{
		"mock-gpt-3.5-turbo": {
			ModelID:     "mock-gpt-3.5-turbo",
			Prefix:      "[GPT-3.5] ",
			Greeting:    "Hello! I'm a mock GPT model. How can I assist you today?",
			Question:    "That's an interesting question. As a mock model, I'll respond with this simulated answer. In a real OpenAI API, you would get a more contextual response.",
			HelpRequest: "I'm here to help! Although I'm just a mock model, I can simulate responses. What do you need assistance with?",
			Default:     "I understand. As a mock GPT model, I'm providing this simulated response to your message. In a real OpenAI API, the response would be generated based on the trained model.",
		},
		"mock-davinci-002": {
			ModelID:          "mock-davinci-002",
			Prefix:           "[DAVINCI] ",
			Greeting:         "there! I'm a mock AI model. How can I assist you today?",
			Question:         "That's an interesting question. As a mock model, I'll provide this simulated answer.",
			HelpRequest:      "is on the way! This is a simulated response from the mock completions API.",
			Default:          "As a mock AI model, I'm continuing your text with this simulated response. In a real OpenAI API, this would be generated based on the trained model.",
			CompletionPrefix: "",
		},
		"deepseek-reasoner": {
			ModelID:          "deepseek-reasoner",
			Prefix:           "[DEEPSEEK] ",
			Greeting:         "Hello! I'm a mock reasoning model. How can I assist you today?",
			Question:         "That's an interesting question. After analyzing the problem, I've arrived at this answer. In a real reasoning model, the response would be more contextual.",
			HelpRequest:      "I'm here to help! After careful reasoning, I can provide this simulated response. What do you need assistance with?",
			Default:          "I understand. After careful reasoning, I'm providing this simulated response. In a real reasoning model, both the reasoning process and final answer would be more sophisticated.",
			SupportReasoning: true,
			ReasoningPrefix:  "REASONING: ",
		},
	}

	// 将模板数据编码为JSON
	data, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding default templates: %v", err)
	}

	// 保存到文件
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory for default templates: %v", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing default template file: %v", err)
	}

	fmt.Printf("Created default templates file: %s\n", filePath)
	return nil
}

// getDefaultTemplate 获取默认模板
func getDefaultTemplate(modelID string) ResponseTemplate {
	return ResponseTemplate{
		ModelID:     modelID,
		Prefix:      "[MOCK] ",
		Greeting:    "Hello! I'm a mock AI model. How can I assist you today?",
		Question:    "That's an interesting question. As a mock model, I'll provide this simulated answer.",
		HelpRequest: "I'm here to help! As a mock model, I can simulate responses.",
		Default:     "I understand. I'm providing this simulated response to your message.",
	}
}
