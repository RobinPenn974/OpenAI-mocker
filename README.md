# OpenAI-mocker

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat&logo=gin)](https://gin-gonic.com/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

一个轻量级的 OpenAI API 接口模拟服务，基于 Gin 框架构建，无需依赖外部数据库。

## 目录

- [项目概述](#项目概述)
- [功能特性](#功能特性)
- [快速开始](#快速开始)
  - [Docker 部署](#docker-部署)
  - [API 密钥管理](#api-密钥管理)
- [API 接口文档](#api-接口文档)
  - [模型 API](#模型-api)
  - [聊天完成 API](#聊天完成-api)
  - [文本完成 API](#文本完成-api)
  - [嵌入 API](#嵌入-api)
  - [重排序 API](#重排序-api)
  - [推理模型 API](#推理模型-api)
- [模型管理](#模型管理)
  - [预设模型](#预设模型)
  - [模型管理接口](#模型管理接口)
- [高级功能](#高级功能)
  - [推理模型功能](#推理模型功能)
- [技术栈](#技术栈)
- [注意事项](#注意事项)

## 项目概述

OpenAI-mocker 是一个模拟 OpenAI API 接口的服务，专为开发和测试环境设计。它提供与官方 API 兼容的接口，使开发者能够在不消耗 API 配额的情况下进行应用开发和测试。

## 功能特性

- ✅ 提供与 OpenAI API 完全兼容的接口
- ✅ 支持模型动态加载和卸载
- ✅ 内置三种类型的虚拟模型（LLM、嵌入、重排序）
- ✅ 简单高效的 API 密钥鉴权系统
- ✅ Docker 支持，便于快速部署
- ✅ 轻量级设计，无需外部数据库依赖
- ✅ 支持 DeepSeek Reasoner 等推理模型的模拟

## 快速开始

### Docker 部署

最简单的方式是使用 Docker Compose：

```bash
# 克隆仓库
git clone https://github.com/yourusername/OpenAI-mocker.git
cd OpenAI-mocker

# 启动服务
docker-compose up -d
```

服务将在后台启动，并在 `http://localhost:8080` 提供 API 访问。

也可以手动构建和运行 Docker 镜像：

```bash
# 构建镜像
docker build -t openai-mocker .

# 运行容器
docker run -d -p 8080:8080 --name openai-mocker openai-mocker
```

### API 密钥管理

系统默认采用**零配置启动模式**：

- 初始状态下，无需 API 密钥即可访问所有 API 接口
- 仅当通过管理接口添加 API 密钥后，系统才会开启鉴权
- 删除所有 API 密钥后，系统将恢复到无需鉴权的状态

#### 获取所有 API 密钥

```bash
curl -X GET http://localhost:8080/admin/auth/keys
```

#### 创建新的 API 密钥

```bash
curl -X POST http://localhost:8080/admin/auth/keys \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-application"
  }'
```

#### 删除指定 API 密钥

```bash
curl -X DELETE http://localhost:8080/admin/auth/keys/sk-mock-xxxx
```

#### 删除所有 API 密钥

```bash
curl -X DELETE http://localhost:8080/admin/auth/keys
```

## API 接口文档

### 模型 API

获取所有可用模型列表：

```bash
curl http://localhost:8080/v1/models \
  -H "Authorization: Bearer sk-mock-xxxx"
```

### 聊天完成 API

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-mock-xxxx" \
  -d '{
    "model": "mock-gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "你好"}]
  }'
```

### 文本完成 API

```bash
curl -X POST http://localhost:8080/v1/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-mock-xxxx" \
  -d '{
    "model": "mock-davinci-002",
    "prompt": "你好"
  }'
```

### 嵌入 API

```bash
curl -X POST http://localhost:8080/v1/embeddings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-mock-xxxx" \
  -d '{
    "model": "mock-embedding-ada-002",
    "input": "你好"
  }'
```

### 重排序 API

```bash
curl -X POST http://localhost:8080/v1/rerank \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-mock-xxxx" \
  -d '{
    "model": "mock-rerank-v1",
    "query": "搜索查询",
    "documents": ["文档1", "文档2", "文档3"]
  }'
```

### 推理模型 API

支持像 DeepSeek Reasoner 这样带有思维链的模型接口：

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-mock-xxxx" \
  -d '{
    "model": "deepseek-reasoner",
    "messages": [{"role": "user", "content": "9.11 and 9.8, which is greater?"}]
  }'
```

## 模型管理

### 预设模型

系统预设了以下虚拟模型，无需额外加载即可使用：

| 模型类型 | 模型名称 | 描述 |
|---------|---------|------|
| LLM | `deepseek-reasoner` | 用于带思维链的对话完成任务 |
| 嵌入 | `mock-embedding-ada-002` | 生成文本嵌入向量 |
| 重排序 | `mock-rerank-v1` | 提供文本重排序功能 |

### 模型管理接口

您可以通过以下 API 动态管理模型：

#### 加载自定义模型

```bash
curl -X POST http://localhost:8080/admin/models/load \
  -H "Content-Type: application/json" \
  -d '{
    "model_id": "custom-gpt-4",
    "model_type": "chat",
    "owned_by": "my-organization"
  }'
```

#### 卸载指定模型

```bash
curl -X POST http://localhost:8080/admin/models/unload \
  -H "Content-Type: application/json" \
  -d '{
    "model_id": "custom-gpt-4"
  }'
```

#### 卸载所有模型

```bash
curl -X POST http://localhost:8080/admin/models/unload_all
```

## 高级功能

### 推理模型功能

系统支持模拟推理模型（如 DeepSeek Reasoner、QwQ-32B 等）的思维链和回答内容输出。

#### 环境变量配置

可通过环境变量 `ENABLE_REASONING` 控制推理模型的行为：

```bash
# 启用推理输出（使用reasoning_content字段）
ENABLE_REASONING=true ./openai-mocker

# 禁用推理输出（将推理过程包含在content字段）
ENABLE_REASONING=false ./openai-mocker
```

#### 推理模式

1. **启用推理模式** (ENABLE_REASONING=true)
   - 思维链内容会通过专门的 `reasoning_content` 字段返回
   - 最终回答通过常规的 `content` 字段返回
   - 适用于支持 reasoning_content 字段的客户端

2. **禁用推理模式** (ENABLE_REASONING=false)
   - 思维链内容会包装在 `<think>...</think>` 标记中
   - 思维链和最终回答都通过 `content` 字段返回
   - 适用于不支持专门推理字段的常规客户端

#### 流式输出支持

推理模型在流式响应中也能正确处理思维链内容：

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-mock-xxxx" \
  -d '{
    "model": "deepseek-reasoner",
    "messages": [{"role": "user", "content": "9.11 and 9.8, which is greater?"}],
    "stream": true
  }'
```

## 技术栈

- **后端框架**：Gin
- **编程语言**：Go
- **存储方式**：内存存储（无需外部数据库）
- **容器化**：Docker & Docker Compose

## 注意事项

- 本项目仅用于开发和测试环境，不建议在生产环境使用
- API 设计与 OpenAI 官方 API 保持兼容，便于从测试环境切换到生产环境
- 虚拟模型仅返回模拟数据，不具备实际的 AI 能力
- 推理模型功能设计参考了 DeepSeek API 和 vLLM 的实现，但仅提供模拟响应
