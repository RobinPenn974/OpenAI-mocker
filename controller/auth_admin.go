package controller

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"RobinPenn974/OpenAI-mocker/middleware"

	"github.com/gin-gonic/gin"
)

// ApiKeyRequest 创建API密钥的请求结构
type ApiKeyRequest struct {
	Name string `json:"name" binding:"required"`
}

// ApiKeyResponse API密钥响应结构
type ApiKeyResponse struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// ApiKeyListResponse API密钥列表响应结构
type ApiKeyListResponse struct {
	Keys []ApiKeyInfo `json:"keys"`
}

// ApiKeyInfo API密钥信息
type ApiKeyInfo struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// GenerateAPIKey 生成一个新的API密钥
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk-mock-" + hex.EncodeToString(bytes), nil
}

// HandleCreateApiKey 处理创建API密钥的请求
func HandleCreateApiKey(c *gin.Context) {
	var req ApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid request: " + err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 生成新的API密钥
	apiKey, err := GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to generate API key: " + err.Error(),
				"type":    "server_error",
			},
		})
		return
	}

	// 添加API密钥到全局存储
	middleware.GlobalApiKeys.AddKey(apiKey, req.Name)

	c.JSON(http.StatusOK, ApiKeyResponse{
		Key:  apiKey,
		Name: req.Name,
	})
}

// HandleDeleteApiKey 处理删除API密钥的请求
func HandleDeleteApiKey(c *gin.Context) {
	keyID := c.Param("key_id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "API key ID is required",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 删除API密钥
	middleware.GlobalApiKeys.RemoveKey(keyID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "API key deleted successfully",
	})
}

// HandleDeleteAllApiKeys 处理删除所有API密钥的请求
func HandleDeleteAllApiKeys(c *gin.Context) {
	// 删除所有API密钥
	middleware.GlobalApiKeys.RemoveAllKeys()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All API keys deleted successfully",
	})
}

// HandleListApiKeys 处理列出所有API密钥的请求
func HandleListApiKeys(c *gin.Context) {
	keys := middleware.GlobalApiKeys.GetAllKeys()

	var keyInfos []ApiKeyInfo
	for key, name := range keys {
		keyInfos = append(keyInfos, ApiKeyInfo{
			Key:  key,
			Name: name,
		})
	}

	c.JSON(http.StatusOK, ApiKeyListResponse{
		Keys: keyInfos,
	})
}
