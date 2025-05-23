package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// ApiKeys 存储API密钥的映射
type ApiKeys struct {
	Keys map[string]string // key -> name
	mu   sync.RWMutex
}

// NewApiKeys 创建一个新的API密钥存储
func NewApiKeys() *ApiKeys {
	return &ApiKeys{
		Keys: make(map[string]string),
	}
}

// Global API Keys实例
var GlobalApiKeys = NewApiKeys()

// AddKey 添加一个新的API密钥
func (a *ApiKeys) AddKey(key, name string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Keys[key] = name
}

// RemoveKey 删除一个API密钥
func (a *ApiKeys) RemoveKey(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.Keys, key)
}

// RemoveAllKeys 删除所有API密钥
func (a *ApiKeys) RemoveAllKeys() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Keys = make(map[string]string)
}

// IsValidKey 检查API密钥是否有效
func (a *ApiKeys) IsValidKey(key string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, exists := a.Keys[key]
	return exists
}

// HasKeys 检查是否存在任何API密钥
func (a *ApiKeys) HasKeys() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.Keys) > 0
}

// GetAllKeys 获取所有API密钥
func (a *ApiKeys) GetAllKeys() map[string]string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// 创建一个副本以避免直接返回内部map
	keysCopy := make(map[string]string, len(a.Keys))
	for k, v := range a.Keys {
		keysCopy[k] = v
	}

	return keysCopy
}

// AuthRequired 验证API密钥的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 健康检查路由不需要认证
		if c.Request.URL.Path == "/v1/healthz" {
			c.Next()
			return
		}

		// 如果没有注册任何API密钥，允许自由访问
		if !GlobalApiKeys.HasKeys() {
			c.Next()
			return
		}

		// 从Authorization头中提取API密钥
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "No API key provided",
					"type":    "auth_error",
					"code":    "no_api_key",
				},
			})
			c.Abort()
			return
		}

		// 支持Bearer token格式
		apiKey := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			apiKey = strings.TrimPrefix(strings.ToLower(authHeader), "bearer ")
			apiKey = strings.TrimSpace(apiKey)
		}

		// 验证API密钥
		if !GlobalApiKeys.IsValidKey(apiKey) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Invalid API key",
					"type":    "auth_error",
					"code":    "invalid_api_key",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 我们不再需要默认API密钥，因为我们现在允许在没有API密钥时自由访问
func init() {
	// 不再添加默认测试密钥
}
