# 跨服务调用指南

## 问题场景

在推荐服务中，我们需要获取用户的帖子数据。这些数据可能来自：
1. **本地数据库**：单体应用，数据在同一个数据库
2. **远程服务**：微服务架构，数据在内容服务

## 解决方案

### 架构设计

```
应用层（Application Layer）
├── 定义接口：ContentServiceClient
└── 使用接口：RecommendationService

基础设施层（Infrastructure Layer）
├── HTTP 实现：ContentServiceHTTPClient
├── RPC 实现：ContentServiceRPCClient
└── 本地实现：ContentRepository
```

### 核心思想

1. **接口在应用层**：定义业务需要什么
2. **实现在基础设施层**：技术细节如何实现
3. **依赖注入**：运行时决定使用哪个实现

## 代码示例

### 1. 定义接口（应用层）

```go
// application/service/recommendation_service.go

// ContentServiceClient 内容服务客户端接口
type ContentServiceClient interface {
    GetRecentPosts(ctx context.Context, userID int64, limit int) ([]*PostInfo, error)
}

// PostInfo 帖子信息（来自远程服务）
type PostInfo struct {
    PostID    int64
    Content   string
    CreatedAt string
}
```

### 2. 实现接口（基础设施层）

#### HTTP 实现

```go
// infrastructure/client/content_service_client.go

type ContentServiceHTTPClient struct {
    baseURL    string
    httpClient *http.Client
}

func (c *ContentServiceHTTPClient) GetRecentPosts(
    ctx context.Context,
    userID int64,
    limit int,
) ([]*service.PostInfo, error) {
    // HTTP 调用实现
    url := fmt.Sprintf("%s/api/v1/users/%d/posts?limit=%d", c.baseURL, userID, limit)
    // ... HTTP 请求逻辑
}
```

#### RPC 实现（Kitex）

```go
// infrastructure/client/content_service_rpc_client.go

type ContentServiceRPCClient struct {
    client contentservice.Client // Kitex 生成的客户端
}

func (c *ContentServiceRPCClient) GetRecentPosts(
    ctx context.Context,
    userID int64,
    limit int,
) ([]*service.PostInfo, error) {
    // RPC 调用实现
    req := &content.GetRecentPostsRequest{
        UserId: userID,
        Limit:  int32(limit),
    }
    resp, err := c.client.GetRecentPosts(ctx, req)
    // ... 转换响应
}
```

### 3. 使用接口（应用层）

```go
// application/service/recommendation_service.go

type RecommendationService struct {
    contentRepo   repository.ContentRepository  // 本地数据库（可选）
    contentClient ContentServiceClient          // 远程服务（可选）
}

// 优先使用远程服务，失败时降级到本地数据库
func (s *RecommendationService) getRecentPosts(
    ctx context.Context,
    userID int64,
    limit int,
) []*dto.PostDTO {
    // 策略1：优先使用远程服务
    if s.contentClient != nil {
        posts, err := s.contentClient.GetRecentPosts(ctx, userID, limit)
        if err == nil {
            return convertToDTO(posts)
        }
    }

    // 策略2：降级到本地数据库
    if s.contentRepo != nil {
        posts, err := s.contentRepo.GetRecentPosts(ctx, userID, limit)
        if err == nil {
            return convertToDTO(posts)
        }
    }

    // 策略3：容错 - 返回空列表
    return []*dto.PostDTO{}
}
```

### 4. 依赖注入（main.go）

```go
// 场景1：单体应用（只用本地数据库）
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    contentRepo,        // 本地数据库
    nil,                // 不使用远程服务
    userRPCClient,
    reasonConfigClient,
)

// 场景2：微服务架构（只用远程服务）
contentClient := client.NewContentServiceHTTPClient("http://content-service:8080")
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    nil,                // 不使用本地数据库
    contentClient,      // 远程服务
    userRPCClient,
    reasonConfigClient,
)

// 场景3：混合架构（优先远程，降级本地）
contentClient := client.NewContentServiceHTTPClient("http://content-service:8080")
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    contentRepo,        // 本地数据库（降级）
    contentClient,      // 远程服务（优先）
    userRPCClient,
    reasonConfigClient,
)
```

## 设计原则

### 1. 接口隔离

```
应用层只关心"需要什么"，不关心"如何实现"
├── 定义：ContentServiceClient 接口
└── 使用：通过接口调用

基础设施层只关心"如何实现"，不关心"为什么需要"
├── HTTP 实现
├── RPC 实现
└── Mock 实现（测试用）
```

### 2. 依赖倒置

```
高层模块（应用层）不依赖低层模块（基础设施层）
两者都依赖抽象（接口）

应用层 → ContentServiceClient 接口 ← 基础设施层
```

### 3. 容错降级

```
优先级：远程服务 > 本地数据库 > 空列表

远程服务成功 → 返回数据
    ↓ 失败
本地数据库成功 → 返回数据
    ↓ 失败
返回空列表 → 不阻塞推荐功能
```

## 对比：Repository vs ServiceClient

| 维度 | Repository | ServiceClient |
|------|-----------|---------------|
| **数据来源** | 本地数据库 | 远程服务 |
| **定义位置** | 领域层（接口） | 应用层（接口） |
| **实现位置** | 基础设施层 | 基础设施层 |
| **返回类型** | 领域对象（Entity） | DTO（PostInfo） |
| **使用场景** | 单体应用 | 微服务架构 |
| **示例** | ContentRepository | ContentServiceClient |

### 为什么返回类型不同？

```go
// Repository：返回领域对象
type ContentRepository interface {
    GetRecentPosts(ctx context.Context, userID valueobject.UserID, limit int) ([]*entity.Post, error)
}

// ServiceClient：返回 DTO
type ContentServiceClient interface {
    GetRecentPosts(ctx context.Context, userID int64, limit int) ([]*PostInfo, error)
}
```

**原因**：
- **Repository**：操作本地数据，可以返回完整的领域对象
- **ServiceClient**：跨服务调用，只能返回序列化的数据（DTO）

## 实际应用场景

### 场景1：从单体迁移到微服务

**阶段1：单体应用**
```go
// 只使用 Repository
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    contentRepo,  // 本地数据库
    nil,          // 不使用远程服务
    userRPCClient,
    reasonConfigClient,
)
```

**阶段2：混合架构（迁移中）**
```go
// 同时使用 Repository 和 ServiceClient
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    contentRepo,      // 本地数据库（降级）
    contentClient,    // 远程服务（优先）
    userRPCClient,
    reasonConfigClient,
)
```

**阶段3：纯微服务**
```go
// 只使用 ServiceClient
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    nil,              // 不使用本地数据库
    contentClient,    // 远程服务
    userRPCClient,
    reasonConfigClient,
)
```

### 场景2：A/B 测试

```go
// 根据配置决定使用哪个实现
var contentClient service.ContentServiceClient

if config.UseNewContentService {
    // 新版本内容服务
    contentClient = client.NewContentServiceHTTPClient("http://content-service-v2:8080")
} else {
    // 旧版本内容服务
    contentClient = client.NewContentServiceHTTPClient("http://content-service-v1:8080")
}
```

### 场景3：测试

```go
// 使用 Mock 实现
type MockContentServiceClient struct{}

func (m *MockContentServiceClient) GetRecentPosts(
    ctx context.Context,
    userID int64,
    limit int,
) ([]*service.PostInfo, error) {
    // 返回测试数据
    return []*service.PostInfo{
        {PostID: 1, Content: "Test Post", CreatedAt: "2024-01-01"},
    }, nil
}

// 测试中使用
recommendationService := service.NewRecommendationService(
    generator,
    socialGraphRepo,
    nil,
    &MockContentServiceClient{}, // Mock 实现
    userRPCClient,
    reasonConfigClient,
)
```

## 总结

### 核心要点

1. **接口在应用层**：定义业务需要什么
2. **实现在基础设施层**：技术细节如何实现
3. **依赖注入**：运行时决定使用哪个实现
4. **容错降级**：远程服务失败不影响核心功能

### 何时使用 Repository vs ServiceClient

- **数据在本地数据库** → 使用 Repository
- **数据在其他服务** → 使用 ServiceClient
- **两者都有** → 同时注入，优先远程服务

### 优势

1. **灵活性**：支持单体和微服务两种架构
2. **可测试性**：可以轻松 Mock
3. **容错性**：远程服务失败不影响核心功能
4. **可维护性**：接口和实现分离
