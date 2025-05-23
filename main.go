package main

import (
	"fmt"
	"log"

	"RobinPenn974/OpenAI-mocker/models"
	"RobinPenn974/OpenAI-mocker/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化默认模型
	models.InitDefaultModels()

	// 创建默认的gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := 8080
	fmt.Printf("OpenAI-mocker服务启动在 :%d\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
