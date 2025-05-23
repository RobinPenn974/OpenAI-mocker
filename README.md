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
- [响应模板系统](#响应模板系统)
  - [模板文件](#模板文件)
  - [模板管理接口](#模板管理接口)
  - [自定义模板](#自定义模板)
- [Hoppscotch 测试工具](#hoppscotch-测试工具)
  - [导入配置](#导入配置)
  - [使用方法](#使用方法)
  - [环境变量](#环境变量)
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
- ✅ 可配置的响应模板系统
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

## 响应模板系统

系统提供了可自定义的响应模板功能，以控制不同模型的回复格式和内容风格。所有模板均以JSON格式存储，便于编辑和管理。

### 模板文件

系统使用两个主要的模板文件：

- `template_data/templates.json`: 主要的模板存储文件
- `template_data/default_templates.json`: 默认模板文件，当主文件不存在时使用

### 模板管理接口

您可以通过以下API管理响应模板：

#### 获取所有模板

```bash
curl -X GET http://localhost:8080/admin/templates
```

#### 获取指定模型的模板

```bash
curl -X GET http://localhost:8080/admin/templates/mock-gpt-3.5-turbo
```

#### 更新指定模型的模板

```bash
curl -X PUT http://localhost:8080/admin/templates/mock-gpt-3.5-turbo \
  -H "Content-Type: application/json" \
  -d '{
    "model_id": "mock-gpt-3.5-turbo",
    "prefix": "[自定义前缀] ",
    "greeting": "您好！我是一个模拟的GPT模型。",
    "question": "这是一个很好的问题。作为模拟模型，我将提供以下回答...",
    "help_request": "我很乐意帮助！尽管我只是一个模拟模型，但我可以提供回答。",
    "default": "理解了。作为模拟GPT模型，我正在提供这个模拟回复。",
    "support_reasoning": false,
    "reasoning_prefix": "",
    "reasoning_template": "",
    "completion_prefix": ""
  }'
```

#### 删除指定模型的模板

```bash
curl -X DELETE http://localhost:8080/admin/templates/mock-gpt-3.5-turbo
```

### 自定义模板

您可以通过以下两种方式自定义模板：

1. **通过API接口**：使用上述PUT接口更新指定模型的模板
2. **直接编辑文件**：修改`template_data/templates.json`或`template_data/default_templates.json`文件

模板格式如下：

```json
{
  "model_id": "模型ID",
  "prefix": "模型回复前缀",
  "greeting": "问候语模板",
  "question": "问题回答模板",
  "help_request": "帮助请求模板",
  "default": "默认回复模板",
  "support_reasoning": false,
  "reasoning_prefix": "推理前缀",
  "reasoning_template": "推理内容模板，支持{question}占位符",
  "completion_prefix": "补全前缀"
}
```

当配置支持推理的模型时，您可以通过以下字段自定义推理行为：

- `support_reasoning`: 设置为 `true` 启用推理功能
- `reasoning_prefix`: 在推理内容前添加的前缀
- `reasoning_template`: 推理内容的模板文本，支持使用 `{question}` 占位符引用用户提问

例如，以下是一个支持推理的模板配置：

```json
{
  "model_id": "custom-reasoning-model",
  "prefix": "[思考] ",
  "greeting": "您好！我是一个支持推理的模型。",
  "default": "根据我的分析，我提供如下回答。",
  "support_reasoning": true,
  "reasoning_prefix": "我的思考过程：",
  "reasoning_template": "我需要思考这个问题：{question}\n\n首先，我会分析关键信息。\n然后，我会应用相关知识来解决问题。\n最后，我会得出合理的结论。"
}
```

### Docker卷挂载修改模板

通过Docker卷挂载是修改模板最方便的方式，无需进入容器内部：

#### Docker命令方式

```bash
docker run -d -p 8080:8080 \
  -v $(pwd)/custom_templates:/app/template_data \
  --name openai-mocker openai-mocker
```

#### Docker Compose方式

```yaml
version: '3'
services:
  openai-mocker:
    image: openai-mocker
    ports:
      - "8080:8080"
    volumes:
      - ./custom_templates:/app/template_data
```

#### 使用步骤

1. 在宿主机创建模板目录

```bash
mkdir -p custom_templates
```

2. 创建或编辑default_templates.json文件

```bash
cat > custom_templates/default_templates.json << 'EOF'
{
  "mock-gpt-3.5-turbo": {
    "model_id": "mock-gpt-3.5-turbo",
    "prefix": "[GPT-3.5] ",
    "greeting": "您好！我是模拟模型。有什么可以帮助您？",
    "question": "这是个好问题。作为模拟模型，我给出以下回答。",
    "help_request": "我很乐意帮忙！请告诉我您需要什么。",
    "default": "明白了。这是一个模拟回复。",
    "support_reasoning": false,
    "reasoning_prefix": "",
    "reasoning_template": "",
    "completion_prefix": ""
  }
}
EOF
```

3. 启动服务后，所有修改会立即生效

## Hoppscotch 测试工具

[Hoppscotch](https://hoppscotch.io/) 是一个轻量级、开源的 API 开发和测试工具，可以方便地用于测试 OpenAI-mocker 服务。项目提供了完整的 Hoppscotch 配置文件，包含所有 API 接口。

### 导入配置

1. 下载项目中的 `OpenAI.json` 配置文件
2. 打开 [Hoppscotch](https://hoppscotch.io/) 网站或应用
3. 点击左侧导航栏的 "Collections"
4. 点击 "Import" 按钮，选择 "Import from JSON"
5. 上传下载好的 `OpenAI.json` 文件
6. 导入后，将在左侧看到 "OpenAI" 集合，包含所有预配置的接口

### 使用方法

Hoppscotch 配置中包含两个主要文件夹：

1. **API** - 包含所有 OpenAI 兼容的 API 接口：
   - 模型列表
   - 聊天完成
   - 文本完成
   - 嵌入接口
   - 重排序接口
   - 推理模型接口
   - 健康检查

2. **admin** - 包含所有管理接口，分为三个子类别：
   - 模型管理 - 加载、卸载模型
   - 模板管理 - 管理响应模板
   - 认证管理 - 管理 API 密钥

### 环境变量

使用前需要设置以下环境变量：

1. 点击 Hoppscotch 顶部的 "Environments" 按钮
2. 创建一个新的环境（如 "OpenAI-mocker"）
3. 添加以下变量：
   - `ip`：服务器 IP 地址（例如 localhost 或 192.168.1.100）
   - `port`：服务器端口（默认 8080）
   - `api_key`：您创建的 API 密钥（可选，仅当启用身份验证时需要）

设置完成后，选择该环境即可开始测试 API。

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
- **存储方式**：内存存储 + JSON文件 (无需外部数据库)
- **容器化**：Docker & Docker Compose

## 注意事项

- 本项目仅用于开发和测试环境，不建议在生产环境使用
- API 设计与 OpenAI 官方 API 保持兼容，便于从测试环境切换到生产环境
- 虚拟模型仅返回模拟数据，不具备实际的 AI 能力
- 推理模型功能设计参考了 DeepSeek API 和 vLLM 的实现，但仅提供模拟响应
