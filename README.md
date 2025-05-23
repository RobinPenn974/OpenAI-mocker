# OpenAI-mocker

## 项目简介

OpenAI-mocker 是一个模拟 OpenAI API 接口的服务，使用 Gin 框架开发的单体应用。本项目不依赖独立的数据库系统，设计保持轻量和简洁。

## 功能特性

- 提供与 OpenAI API 兼容的接口
- 支持模型加载和卸载功能
- 支持简单的 API 密钥鉴权
- 内置三种类型的虚拟模型（LLM、嵌入、重排序）
- 轻量级设计，易于部署和使用

## 支持的 API 接口

本项目实现了以下 OpenAI 兼容的 API 接口：

### Chat Completions API
- **POST /v1/chat/completions**: 生成对话完成内容
- 默认虚拟模型: `mock-gpt-3.5-turbo`

### Completions API
- **POST /v1/completions**: 生成文本完成内容
- 默认虚拟模型: `mock-davinci-002`

### Embeddings API
- **POST /v1/embeddings**: 生成文本嵌入向量
- 默认虚拟模型: `mock-embedding-ada-002`

### Rerank API
- **POST /v1/rerank**: 文本重排序
- 默认虚拟模型: `mock-rerank-v1`

### Models API
- **GET /v1/models**: 获取可用模型列表
  - 返回所有内置虚拟模型和已加载的自定义模型

### 健康检查
- **GET /v1/healthz**: API 服务健康状态检查

## 内置虚拟模型

本项目预设了以下虚拟模型，无需额外加载即可使用：

1. **LLM 模型**
   - `mock-gpt-3.5-turbo`: 用于对话完成任务
   - `mock-davinci-002`: 用于文本完成任务

2. **嵌入模型**
   - `mock-embedding-ada-002`: 生成标准化的文本嵌入向量

3. **重排序模型**
   - `mock-rerank-v1`: 提供基础的文本重排序功能

## 模型管理

本项目提供以下模型管理功能和对应接口：

### 模型加载
- **POST /admin/models/load**: 加载指定模型到内存中
- **POST /admin/models/preload**: 预加载配置文件中的默认模型

### 模型卸载
- **POST /admin/models/unload**: 从内存中卸载指定模型
- **POST /admin/models/unload_all**: 卸载所有已加载模型

### 鉴权配置
- **POST /admin/auth/keys**: 创建新的 API 密钥
- **DELETE /admin/auth/keys/{key_id}**: 删除指定 API 密钥

## 快速开始

(安装和使用说明将在项目完成后补充)

## 技术栈

- 后端框架：Gin
- 语言：Go
- 存储：内存存储（无需外部数据库）

## 注意事项

- 本项目仅用于开发和测试环境，不建议在生产环境使用
- API 设计与 OpenAI 官方 API 保持兼容，便于切换
- 虚拟模型仅返回模拟数据，不具备实际的 AI 能力
