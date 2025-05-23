# 模板系统使用说明

## 模板文件

系统使用两个主要的模板文件：

- `templates.json`: 主要的模板存储文件，包含所有注册的模板。
- `default_templates.json`: 默认模板文件，当主文件不存在时，系统会从此文件复制模板。

## 模板格式

模板使用JSON格式存储，每个模板包含以下字段：

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
  "completion_prefix": "补全前缀"
}
```

## 添加新模板

要添加新模板，可以直接编辑`templates.json`文件，添加新的模板条目，或者使用API：

```go
templates.RegisterTemplate(templates.ResponseTemplate{
  ModelID: "new-model",
  Prefix: "[NEW] ",
  Greeting: "Hello!",
  // ... 其他字段
})
```

## 删除模板

要删除模板，可以直接从`templates.json`文件中删除对应条目，或使用API：

```go
templates.DeleteTemplate("model-id")
```

## 自定义默认模板

如果需要更改默认模板，可以编辑`default_templates.json`文件。 