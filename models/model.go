package models

import (
	"errors"
	"sync"
)

// 模型类型常量
const (
	ModelTypeLLM       = "llm"
	ModelTypeEmbedding = "embedding"
	ModelTypeRerank    = "rerank"
)

// ModelInfo 表示模型的基本信息
type ModelInfo struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Created   int64  `json:"created"`
	OwnedBy   string `json:"owned_by"`
	ModelType string `json:"model_type"` // llm, embedding, rerank
}

// 全局模型存储
var (
	models     = make(map[string]ModelInfo)
	modelMutex sync.RWMutex
)

// RegisterModel 注册一个模型到系统中
func RegisterModel(model ModelInfo) {
	modelMutex.Lock()
	defer modelMutex.Unlock()
	models[model.ID] = model
}

// GetModel 获取指定ID的模型信息
func GetModel(id string) (ModelInfo, error) {
	modelMutex.RLock()
	defer modelMutex.RUnlock()

	model, exists := models[id]
	if !exists {
		return ModelInfo{}, errors.New("model not found")
	}
	return model, nil
}

// ListModels 列出所有已注册的模型
func ListModels() []ModelInfo {
	modelMutex.RLock()
	defer modelMutex.RUnlock()

	result := make([]ModelInfo, 0, len(models))
	for _, model := range models {
		result = append(result, model)
	}
	return result
}

// UnloadModel 卸载指定ID的模型
func UnloadModel(id string) error {
	modelMutex.Lock()
	defer modelMutex.Unlock()

	// 检查模型是否存在
	if _, exists := models[id]; !exists {
		return errors.New("model not found")
	}

	delete(models, id)
	return nil
}

// UnloadAllModels 卸载所有模型
func UnloadAllModels() {
	modelMutex.Lock()
	defer modelMutex.Unlock()

	// 清空模型映射
	models = make(map[string]ModelInfo)
}

// InitDefaultModels 初始化默认内置模型
func InitDefaultModels() {
	// 注册LLM模型
	RegisterModel(ModelInfo{
		ID:        "mock-gpt-3.5-turbo",
		Object:    "model",
		Created:   1677610602,
		OwnedBy:   "openai-mocker",
		ModelType: ModelTypeLLM,
	})

	RegisterModel(ModelInfo{
		ID:        "mock-davinci-002",
		Object:    "model",
		Created:   1649880484,
		OwnedBy:   "openai-mocker",
		ModelType: ModelTypeLLM,
	})

	// 注册推理模型
	RegisterModel(ModelInfo{
		ID:        "deepseek-reasoner",
		Object:    "model",
		Created:   1714207996,
		OwnedBy:   "openai-mocker",
		ModelType: ModelTypeLLM,
	})

	// 注册Embedding模型
	RegisterModel(ModelInfo{
		ID:        "mock-embedding-ada-002",
		Object:    "model",
		Created:   1671217299,
		OwnedBy:   "openai-mocker",
		ModelType: ModelTypeEmbedding,
	})

	// 注册Rerank模型
	RegisterModel(ModelInfo{
		ID:        "mock-rerank-v1",
		Object:    "model",
		Created:   1709486145,
		OwnedBy:   "openai-mocker",
		ModelType: ModelTypeRerank,
	})
}
