package routes

import (
	"RobinPenn974/OpenAI-mocker/controller"
	"RobinPenn974/OpenAI-mocker/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// 健康检查路由 - 不需要认证
	r.GET("/v1/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 路由组 - 需要认证
	v1 := r.Group("/v1")
	v1.Use(middleware.AuthRequired())
	{
		// Chat Completions API
		v1.POST("/chat/completions", controller.HandleChatCompletions)

		// Completions API
		v1.POST("/completions", controller.HandleCompletions)

		// Embeddings API
		v1.POST("/embeddings", controller.HandleEmbeddings)

		// Rerank API
		v1.POST("/rerank", controller.HandleRerank)

		// Models API
		v1.GET("/models", controller.HandleListModels)
	}

	// 管理员API路由组
	admin := r.Group("/admin")
	{
		// 模型管理
		models := admin.Group("/models")
		models.POST("/load", controller.HandleLoadModel)
		models.POST("/unload", controller.HandleUnloadModel)
		models.POST("/unload_all", controller.HandleUnloadAllModels)

		// 认证管理
		auth := admin.Group("/auth")
		auth.GET("/keys", controller.HandleListApiKeys)
		auth.POST("/keys", controller.HandleCreateApiKey)
		auth.DELETE("/keys/:key_id", controller.HandleDeleteApiKey)
		auth.DELETE("/keys", controller.HandleDeleteAllApiKeys)
	}
}
