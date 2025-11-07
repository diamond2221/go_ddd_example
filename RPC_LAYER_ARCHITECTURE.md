# RPC 层架构说明

## 完整的系统架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                         客户端（其他服务）                            │
│  - Web 服务                                                          │
│  - Mobile 服务                                                       │
│  - 其他微服务                                                        │
└─────────────────────────────────────────────────────────────────────┘
                              ↓ RPC 调用
┌─────────────────────────────────────────────────────────────────────┐
│                    Kitex RPC 框架（网络层）                          │
│  - 序列化/反序列化（Thrift）                                         │
│  - 网络传输（TCP）                                                   │
│  - 服务发现                                                          │
│  - 负载均衡                                                          │
│  - 链路追踪                                                          │
└─────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────┐
│              RPC 生成代码（rpc_gen/kitex_gen）                       │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  recommendation.go - RPC 数据结构                             │  │
│  │  - GetRecommendationsRequest                                  │  │
│  │  - GetRecommendationsResponse                                 │  │
│  │  - UserRecommendation                                         │  │
│  │  - Post                                                       │  │
│  └───────────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  recommendationservice.go - 服务接口                          │  │
│  │  - RecommendationService interface                            │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ↓ 实现接口
┌─────────────────────────────────────────────────────────────────────┐
│                    接口层（Interface Layer）                         │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  RecommendationHandler                                        │  │
│  │  - 实现 RecommendationService 接口                            │  │
│  │  - 协议适配（RPC → 应用层）                                   │  │
│  │  - 参数验证                                                   │  │
│  │  - 错误处理                                                   │  │
│  │  - RPC 对象 ↔ DTO 转换                                        │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────┐
│                   应用层（Application Layer）                        │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  RecommendationService                                        │  │
│  │  - 用例编排                                                   │  │
│  │  - 跨服务调用（RPC 到 user 服务、content 服务）               │  │
│  │  - DTO ↔ 领域对象转换                                         │  │
│  │  - 事务管理                                                   │  │
│  └───────────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  DTO（数据传输对象）                                          │  │
│  │  - RecommendationResponse                                     │  │
│  │  - UserRecommendationDTO                                      │  │
│  │  - PostDTO                                                    │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────┐
│                     领域层（Domain Layer）                           │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  领域服务（Domain Service）                                   │  │
│  │  - RecommendationGenerator                                    │  │
│  │  - 推荐算法核心逻辑                                           │  │
│  └───────────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  聚合（Aggregate）                                            │  │
│  │  - UserRecommendation                                         │  │
│  │  - RecommendationList                                         │  │
│  └───────────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  值对象（Value Object）                                       │  │
│  │  - UserID, RecommendationID, RecommendationReason             │  │
│  └───────────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  仓储接口（Repository Interface）                             │  │
│  │  - SocialGraphRepository                                      │  │
│  │  - ContentRepository                                          │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ↑ 依赖倒置
┌─────────────────────────────────────────────────────────────────────┐
│                 基础设施层（Infrastructure Layer）                   │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  仓储实现（Repository Implementation）                        │  │
│  │  - SocialGraphRepositoryImpl                                  │  │
│  │  - ContentRepositoryImpl                                      │  │
│  │  - 数据库访问（GORM）                                         │  │
│  │  - PO ↔ 领域对象转换                                          │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────┐
│                          数据库（MySQL）                             │
│  - follows 表                                                        │
│  - posts 表                                                          │
└─────────────────────────────────────────────────────────────────────┘
```

## 数据对象的转换流程

### 请求流程（客户端 → 服务端）

```
1. 客户端创建 RPC 请求
   ┌─────────────────────────────────────┐
   │ RPC Request                         │
   │ {                                   │
   │   user_id: 123,                     │
   │   limit: 10                         │
   │ }                                   │
   └─────────────────────────────────────┘
                ↓ Kitex 序列化
   ┌─────────────────────────────────────┐
   │ 二进制数据（Thrift 格式）            │
   └─────────────────────────────────────┘
                ↓ 网络传输
   ┌─────────────────────────────────────┐
   │ 服务端接收                           │
   └─────────────────────────────────────┘
                ↓ Kitex 反序列化
   ┌─────────────────────────────────────┐
   │ GetRecommendationsRequest           │
   │ (RPC 对象)                          │
   └─────────────────────────────────────┘
                ↓ Handler 处理
   ┌─────────────────────────────────────┐
   │ 参数验证                             │
   │ - user_id > 0?                      │
   │ - limit 有默认值                     │
   └─────────────────────────────────────┘
                ↓ 调用应用服务
   ┌─────────────────────────────────────┐
   │ 应用服务处理                         │
   │ - 转换为领域对象                     │
   │ - 调用领域服务                       │
   │ - 跨服务调用                         │
   └─────────────────────────────────────┘
```

### 响应流程（服务端 → 客户端）

```
1. 领域层返回聚合
   ┌─────────────────────────────────────┐
   │ RecommendationList                  │
   │ (领域聚合)                          │
   │ - 包含业务逻辑                       │
   │ - 私有字段                           │
   └─────────────────────────────────────┘
                ↓ 应用服务转换
   ┌─────────────────────────────────────┐
   │ RecommendationResponse              │
   │ (DTO)                               │
   │ - 简单数据结构                       │
   │ - 公开字段                           │
   └─────────────────────────────────────┘
                ↓ Handler 转换
   ┌─────────────────────────────────────┐
   │ GetRecommendationsResponse          │
   │ (RPC 对象)                          │
   │ - 包含 Thrift 标签                   │
   │ - 可序列化                           │
   └─────────────────────────────────────┘
                ↓ Kitex 序列化
   ┌─────────────────────────────────────┐
   │ 二进制数据（Thrift 格式）            │
   └─────────────────────────────────────┘
                ↓ 网络传输
   ┌─────────────────────────────────────┐
   │ 客户端接收                           │
   └─────────────────────────────────────┘
                ↓ Kitex 反序列化
   ┌─────────────────────────────────────┐
   │ GetRecommendationsResponse          │
   │ (客户端使用)                        │
   └─────────────────────────────────────┘
```

## 三种对象的详细对比

### 1. RPC 对象（kitex_gen）

**位置**：`rpc_gen/kitex_gen/recommendation/`

**特点**：
- 由 Kitex 工具自动生成
- 包含 Thrift 序列化标签
- 用于网络传输
- 不包含业务逻辑

**示例**：
```go
type UserRecommendation struct {
    UserId      int64   `thrift:"user_id,1,required"`
    Username    string  `thrift:"username,2,required"`
    Avatar      string  `thrift:"avatar,3,required"`
    Bio         string  `thrift:"bio,4,optional"`
    Reason      string  `thrift:"reason,5,required"`
    Score       int32   `thrift:"score,6,required"`
    RecentPosts []*Post `thrift:"recent_posts,7,required"`
}
```

**使用场景**：
- RPC 客户端调用
- RPC 服务端接收请求
- 网络传输

### 2. DTO（应用层）

**位置**：`application/dto/`

**特点**：
- 应用层定义
- 简单的数据结构
- 可能包含 JSON 标签
- 不依赖 RPC 框架

**示例**：
```go
type UserRecommendationDTO struct {
    UserID      int64      `json:"user_id"`
    Username    string     `json:"username"`
    Avatar      string     `json:"avatar"`
    Bio         string     `json:"bio"`
    Reason      string     `json:"reason"`
    Score       int        `json:"score"`
    RecentPosts []*PostDTO `json:"recent_posts"`
}
```

**使用场景**：
- 应用层和接口层之间传输
- 可能用于 HTTP API
- 内部服务调用

### 3. 领域对象（领域层）

**位置**：`domain/aggregate/`

**特点**：
- 领域层定义
- 包含业务逻辑和行为
- 字段私有，通过方法访问
- 不依赖任何外层

**示例**：
```go
type UserRecommendation struct {
    id              RecommendationID
    targetUserID    UserID
    reason          RecommendationReason
    score           int
    recentPostCount int
    createdAt       time.Time
    expiresAt       time.Time
}

func (r *UserRecommendation) CalculateScore() int {
    // 业务逻辑
}

func (r *UserRecommendation) IsExpired() bool {
    // 业务逻辑
}
```

**使用场景**：
- 领域层业务逻辑
- 封装业务规则
- 保证数据一致性

## 为什么需要三层对象？

### 问题：为什么不直接用一种对象？

#### 方案 1：只用 RPC 对象
```go
// ❌ 问题
type UserRecommendation struct {
    UserId int64 `thrift:"user_id,1,required"`
    Score  int32 `thrift:"score,6,required"`
}

// 业务逻辑混在一起
func (r *UserRecommendation) CalculateScore() int {
    // 领域逻辑依赖 RPC 框架
}
```

**问题**：
- 领域逻辑依赖 RPC 框架
- 无法切换协议（如从 Thrift 到 Protobuf）
- 难以测试（需要 RPC 环境）
- 字段必须公开（序列化需要）

#### 方案 2：只用领域对象
```go
// ❌ 问题
type UserRecommendation struct {
    id     RecommendationID  // 私有字段
    score  int               // 私有字段
}

// 无法序列化
client.Call(recommendation) // 编译错误
```

**问题**：
- 私有字段无法序列化
- 领域对象结构可能很复杂
- 不适合网络传输

### 解决方案：分层对象

```
RPC 对象 ←→ DTO ←→ 领域对象
  ↑          ↑         ↑
网络层    应用层    领域层
```

**好处**：
1. **关注点分离**：每层对象有明确的职责
2. **独立演进**：修改领域模型不影响 RPC 接口
3. **可测试性**：领域对象可以纯内存测试
4. **灵活性**：可以支持多种协议（RPC、HTTP、MQ）

## 实际代码示例

### Handler 中的转换

```go
// interface/handler/recommendation_handler.go

func (h *RecommendationHandler) GetFollowingBasedRecommendations(
    ctx context.Context,
    req *recommendation.GetRecommendationsRequest, // RPC 对象
) (*recommendation.GetRecommendationsResponse, error) {

    // 1. RPC 对象 → 基本类型
    userID := req.UserId
    limit := int(req.Limit)

    // 2. 调用应用服务（返回 DTO）
    dtoResp, err := h.service.GetFollowingBasedRecommendations(
        ctx, userID, limit,
    )
    if err != nil {
        return nil, err
    }

    // 3. DTO → RPC 对象
    rpcResp := &recommendation.GetRecommendationsResponse{
        Recommendations: make([]*recommendation.UserRecommendation, 0),
    }

    for _, dto := range dtoResp.Recommendations {
        rpcRec := &recommendation.UserRecommendation{
            UserId:   dto.UserID,
            Username: dto.Username,
            // ...
        }
        rpcResp.Recommendations = append(rpcResp.Recommendations, rpcRec)
    }

    return rpcResp, nil
}
```

### 应用服务中的转换

```go
// application/service/recommendation_service.go

func (s *RecommendationService) GetFollowingBasedRecommendations(
    ctx context.Context,
    userID int64,
    limit int,
) (*dto.RecommendationResponse, error) {

    // 1. 基本类型 → 领域对象
    domainUserID, _ := valueobject.NewUserID(userID)

    // 2. 调用领域服务（返回领域对象）
    list, err := s.generator.GenerateFollowingBasedRecommendations(
        ctx, domainUserID, 7,
    )
    if err != nil {
        return nil, err
    }

    // 3. 领域对象 → DTO
    dtoResp := &dto.RecommendationResponse{
        Recommendations: make([]*dto.UserRecommendationDTO, 0),
    }

    for _, rec := range list.GetTopN(limit) {
        dtoRec := &dto.UserRecommendationDTO{
            UserID:   rec.TargetUserID().Value(),
            Reason:   rec.Reason().Description(),
            Score:    rec.Score(),
            // ...
        }
        dtoResp.Recommendations = append(dtoResp.Recommendations, dtoRec)
    }

    return dtoResp, nil
}
```

## 总结

### RPC 层在 DDD 中的作用

1. **技术适配**：将 RPC 协议适配到应用层
2. **边界保护**：保护内部领域模型不被外部影响
3. **协议无关**：业务逻辑不依赖具体的通信协议

### 关键设计原则

1. **依赖方向**：外层依赖内层，内层不依赖外层
2. **关注点分离**：每层对象有明确的职责
3. **防腐层**：RPC 层是防腐层，保护领域模型

### 实践建议

1. **不要跳过转换**：即使看起来麻烦，也要做好转换
2. **保持领域纯粹**：领域对象不要依赖 RPC 框架
3. **合理使用 DTO**：DTO 是应用层和接口层的桥梁
4. **自动化测试**：每层都应该有独立的测试

通过理解 RPC 层在 DDD 架构中的位置和作用，可以更好地设计和实现微服务系统。
