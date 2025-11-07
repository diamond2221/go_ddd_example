# 推荐理由文案配置服务集成指南

## 概述

支持从后端配置服务动态获取推荐理由文案，同时保留本地降级逻辑。

## 架构设计

```
┌─────────────┐
│  接口层      │
└──────┬──────┘
       │
┌──────▼──────────────────────────────────┐
│  应用层 (RecommendationService)          │
│  - 调用配置服务获取文案                    │
│  - 降级到本地逻辑                        │
└──────┬──────────────────────────────────┘
       │
┌──────▼──────────────────────────────────┐
│  领域层 (RecommendationReason)           │
│  - 本地文案生成逻辑（降级）                 │
└─────────────────────────────────────────┘
       │
┌──────▼──────────────────────────────────┐
│  基础设施层 (ReasonTextConfigHTTPClient) │
│  - HTTP 调用配置服务                      │
└─────────────────────────────────────────┘
```

## 使用方式

### 1. 启用配置服务

```go
// 创建配置服务客户端
reasonConfigClient := client.NewReasonTextConfigHTTPClient(
    "http://config-service:8080",
)

// 创建推荐服务
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    contentRepo,
    userRPCClient,
    reasonConfigClient, // 传入配置客户端
)
```

### 2. 不使用配置服务（降级到本地逻辑）

```go
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    contentRepo,
    userRPCClient,
    nil, // 传 nil，使用本地逻辑
)
```


## 配置服务 API 规范

### 请求

```
GET /api/v1/recommendation/reason-text?type={reasonType}&count={count}
```

参数：
- `type`: 推荐理由类型
  - `followed_by_following`: 关注的人关注了TA
  - `popular_in_network`: 网络中受欢迎
- `count`: 相关用户数量（用于生成文案）

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "text": "你的 3 位好友也关注了TA"
  }
}
```

## 降级策略

配置服务不可用时，自动降级到本地逻辑：

1. `reasonConfigClient` 为 `nil` → 使用本地逻辑
2. HTTP 请求失败 → 使用本地逻辑
3. 返回空字符串 → 使用本地逻辑

本地逻辑示例：
- `followed_by_following` + count=3 → "3 位你关注的人也关注了TA"
- `popular_in_network` → "在你的社交网络中很受欢迎"

## 渐进式迁移

### 阶段1：部署前端（保留降级）
```go
// 支持配置服务，但可以传 nil
reasonConfigClient := nil // 暂时不启用
```

### 阶段2：灰度配置服务
```go
// 通过特性开关控制
if featureFlag.IsEnabled("use_reason_config") {
    reasonConfigClient = client.NewReasonTextConfigHTTPClient(...)
}
```

### 阶段3：全量上线
```go
// 所有用户使用配置服务
reasonConfigClient := client.NewReasonTextConfigHTTPClient(...)
```

## 扩展建议

1. **添加缓存**：减少 HTTP 调用
2. **A/B 测试**：不同用户看到不同文案
3. **多语言支持**：根据用户语言返回对应文案
4. **监控告警**：配置服务异常时告警
