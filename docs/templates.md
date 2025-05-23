# 模型响应模板管理

OpenAI-mocker 支持通过模板管理每个模型的响应内容，并将模板持久化存储到本地文件。本文档介绍如何使用模板功能。

## 模板结构

每个模型的响应模板包含以下字段：

```json
{
  "model_id": "模型ID",
  "prefix": "模型前缀，如 [GPT-3.5] ",
  "greeting": "问候语模板",
  "question": "问题回答模板",
  "help_request": "帮助请求模板",
  "default": "默认回复模板",
  "support_reasoning": false,
  "reasoning_prefix": "推理内容前缀",
  "completion_prefix": "补全前缀"
}
```

## 在加载模型时指定模板

您可以在加载模型时同时指定模板：

```bash
curl -X POST http://localhost:8080/admin/models/load \
  -H "Content-Type: application/json" \
  -d '{
    "model_id": "custom-gpt-4",
    "model_type": "llm",
    "owned_by": "my-organization",
    "template": {
      "prefix": "[自定义GPT-4] ",
      "greeting": "您好！我是一个自定义的GPT-4模型。请问我能帮您做什么？",
      "question": "这是一个很好的问题。作为GPT-4的模拟模型，我的回答是...",
      "help_request": "我很乐意提供帮助！请详细描述您需要什么样的协助。",
      "default": "我理解您的意思。作为一个模拟GPT-4的模型，我的回复是...",
      "support_reasoning": true,
      "reasoning_prefix": "思考过程: "
    }
  }'
```

如果未指定某些字段，系统会自动填充合理的默认值。

## 管理模板

### 列出所有模板

```bash
curl -X GET http://localhost:8080/admin/templates
```

### 获取特定模型的模板

```bash
curl -X GET http://localhost:8080/admin/templates/custom-gpt-4
```

### 更新模板

```bash
curl -X PUT http://localhost:8080/admin/templates/custom-gpt-4 \
  -H "Content-Type: application/json" \
  -d '{
    "prefix": "[升级版GPT-4] ",
    "greeting": "您好！我是升级版GPT-4模型。有什么可以帮助您的？",
    "question": "这个问题很有挑战性。我的分析如下...",
    "help_request": "我随时准备提供专业帮助！",
    "default": "我已经理解您的需求。我的回应是...",
    "support_reasoning": true,
    "reasoning_prefix": "分析思路: "
  }'
```

### 删除模板

```bash
curl -X DELETE http://localhost:8080/admin/templates/custom-gpt-4
```

## 模板持久化

所有模板会自动持久化存储到 `template_data/templates.json` 文件，服务重启后会自动加载。

## 内置模型模板

系统预设了以下内置模型的默认模板：

1. `mock-gpt-3.5-turbo` - 普通聊天模型
2. `mock-davinci-002` - 文本补全模型
3. `deepseek-reasoner` - 支持推理功能的模型

这些模板也会被持久化，允许您根据需要修改它们。

## 高级用法

### 推理模型

对于支持推理的模型，将 `support_reasoning` 设置为 `true`，并通过 `reasoning_prefix` 指定推理内容的前缀。

根据环境变量 `ENABLE_REASONING` 的值，系统会决定：

- 当设置为 `true` 时：通过专门的 `reasoning_content` 字段返回推理内容
- 当设置为 `false` 或未设置时：将推理内容作为 `<think>...</think>` 标记包装在 `content` 字段中

### 模板替换变量

在未来版本中，我们计划支持在模板中使用变量，例如：

```
"greeting": "你好，{user_name}！我是 {model_name}，很高兴为您服务。"
``` 