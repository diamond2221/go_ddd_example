# DDD 快速参考

## 核心概念速查

### 值对象（Value Object）
```go
// 特点：不可变、无标识、类型安全
type UserID struct {
    value int64  // 私有字段
}

func NewUserID(value int64) (UserID, error) {
    if value <= 0 {
        return UserID{}, ErrInvalidUserID
    }
    return UserID{value: value}, nil
}
```

**何时使用**：表示业务概念，需要验证规则

### 实体（Entity）
```go
// 特点：有唯一标识、有生命周期
type Post struct {
    id        PostID
    authorID  UserID
    content   string
    createdAt time.Time
}

func (p *Post) ID() PostID {
    return p.id
}
```

**何时使用**：有唯一标识，需要追踪生命周期

### 聚合（Aggregate）
```go
// 特点：定义事务边界、封装业务规则
type UserRecommendation struct {
    id           RecommendationID
    targetUserID UserID
    score        int
}

func NewUserRecommendation(...) (*UserRecommendation, error) {
    // 创建时验证业务规则
}

func (r *UserRecommendation) CalculateScore() int {
    // 业务逻辑
}
```

**何时使用**：需要保证一致性的对象集合

### 领域服务（Domain Service）
```go
// 特点：跨聚合的业务逻辑
type RecommendationGenerator struct {
    socialGraphRepo SocialGraphRepository
    contentRepo     ContentRepository
}

func (g *RecommendationGenerator) GenerateRecommendations(...) {
    // 协调多个聚合的业务逻辑
}
```

**何时使用**：业务逻辑不属于任何单一聚合

### 应用服务（Application Service）
```go
// 特点：用例编排、DTO 转换
type RecommendationService struct {
    generator   *RecommendationGenerator
    userClient  UserRPCClient
}

func (s *RecommendationService) GetRecommendations(...) {
    // 1. 调用领域服务
    // 2. 跨服务调用
    // 3. DTO 转换
}
```

**何时使用**：编排用例，处理技术细节

### 仓储（Repository）
```go
// 接口定义在领域层
type SocialGraphRepository interface {
    GetFollowings(ctx context.Context, userID UserID) ([]UserID, error)
}

// 实现在基础设施层
type SocialGraphRepositoryImpl struct {
    db *gorm.DB
}
```

**何时使用**：访问聚合，隔离数据访问

## 分层职责速查

| 层 | 职责 | 不应该包含 |
|---|------|-----------|
| 接口层 | 协议适配、参数验证 | 业务逻辑、用例编排 |
| 应用层 | 用例编排、DTO转换 | 业务规则、数据访问 |
| 领域层 | 业务规则、领域模型 | 技术细节、外部依赖 |
| 基础设施层 | 数据访问、技术实现 | 业务逻辑 |

## 依赖方向

```
接口层 → 应用层 → 领域层 ← 基础设施层
```

**原则**：所有层都依赖领域层，领域层不依赖任何外层

## 对象转换速查

### RPC 对象
```go
// 位置：rpc_gen/kitex_gen
// 用途：网络传输
type UserRecommendation struct {
    UserId   int64  `thrift:"user_id,1,required"`
    Username string `thrift:"username,2,required"`
}
```

### DTO
```go
// 位置：application/dto
// 用途：应用层数据传输
type UserRecommendationDTO struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
}
```

### 领域对象
```go
// 位置：domain/aggregate
// 用途：业务逻辑
type UserRecommendation struct {
    id           RecommendationID
    targetUserID UserID
}

func (r *UserRecommendation) CalculateScore() int {
    // 业务逻辑
}
```

## 常见模式速查

### 工厂方法
```go
// 用于创建聚合，保证业务规则
func NewUserRecommendation(
    targetUserID UserID,
    reason RecommendationReason,
    postCount int,
) (*UserRecommendation, error) {
    // 验证业务规则
    if len(reason.RelatedUsers()) == 0 {
        return nil, ErrNoReason
    }

    // 计算分数
    score := calculateScore(reason, postCount)

    return &UserRecommendation{
        id:           NewRecommendationID(),
        targetUserID: targetUserID,
        score:        score,
    }, nil
}
```

### 业务规则封装
```go
// 在聚合中封装业务规则
func (r *UserRecommendation) IsExpired() bool {
    return time.Now().After(r.expiresAt)
}

func (l *RecommendationList) AddRecommendation(rec *UserRecommendation) error {
    // 业务规则：不能推荐自己
    if rec.TargetUserID().Equals(l.forUserID) {
        return ErrCannotRecommendSelf
    }

    // 业务规则：不能重复推荐
    for _, existing := range l.recommendations {
        if existing.TargetUserID().Equals(rec.TargetUserID()) {
            return ErrDuplicateRecommendation
        }
    }

    l.recommendations = append(l.recommendations, rec)
    return nil
}
```

### 依赖倒置
```go
// 领域层定义接口
package repository

type SocialGraphRepository interface {
    GetFollowings(ctx context.Context, userID UserID) ([]UserID, error)
}

// 基础设施层实现接口
package persistence

type SocialGraphRepositoryImpl struct {
    db *gorm.DB
}

func (r *SocialGraphRepositoryImpl) GetFollowings(...) ([]UserID, error) {
    // 数据库访问实现
}

// 构造函数返回接口类型
func NewSocialGraphRepository(db *gorm.DB) repository.SocialGraphRepository {
    return &SocialGraphRepositoryImpl{db: db}
}
```

## 命名约定

### 值对象
- 名词：`UserID`, `RecommendationReason`
- 工厂方法：`NewXxx()`
- 访问器：`Value()`, `Equals()`

### 实体
- 名词：`Post`, `User`
- 工厂方法：`NewXxx()`
- 访问器：`ID()`, `GetXxx()`

### 聚合
- 名词：`UserRecommendation`, `RecommendationList`
- 工厂方法：`NewXxx()`
- 业务方法：`CalculateScore()`, `IsExpired()`

### 领域服务
- 名词+动词：`RecommendationGenerator`
- 方法：`GenerateXxx()`, `CalculateXxx()`

### 应用服务
- 名词+Service：`RecommendationService`
- 方法：`GetXxx()`, `CreateXxx()`

### 仓储
- 名词+Repository：`SocialGraphRepository`
- 方法：`Get()`, `Find()`, `Save()`

## 测试策略

### 单元测试（领域层）
```go
func TestUserRecommendation_CalculateScore(t *testing.T) {
    // 测试业务规则
    reason := NewFollowedByFollowingReason([]UserID{u1, u2, u3})
    rec, _ := NewUserRecommendation(targetUser, reason, 5)

    expected := 3*10 + 5*2  // 3个关注者 + 5个帖子
    assert.Equal(t, expected, rec.Score())
}
```

### 集成测试（应用层）
```go
func TestRecommendationService_GetRecommendations(t *testing.T) {
    // 使用 mock 仓储
    mockRepo := &MockSocialGraphRepository{}
    service := NewRecommendationService(mockRepo, ...)

    // 测试完整用例
    result, err := service.GetRecommendations(ctx, userID, 10)
    assert.NoError(t, err)
    assert.NotEmpty(t, result.Recommendations)
}
```

## 常见问题

### Q: 什么时候用值对象，什么时候用实体？
**A**: 如果对象有唯一标识且生命周期重要，用实体；否则用值对象。

### Q: 聚合应该多大？
**A**: 尽可能小，只包含必须保持一致性的对象。

### Q: 什么逻辑放在领域层，什么放在应用层？
**A**: 纯业务规则放领域层，用例编排和技术细节放应用层。

### Q: 为什么需要三层对象（RPC、DTO、领域）？
**A**: 关注点分离，每层对象有明确的职责，便于独立演进。

### Q: 是否所有项目都适合 DDD？
**A**: 不是。简单的 CRUD 项目用传统架构更合适。DDD 适合复杂业务领域。

## 设计检查清单

### 值对象
- [ ] 字段是否私有？
- [ ] 是否不可变？
- [ ] 是否有工厂方法？
- [ ] 是否有验证规则？

### 聚合
- [ ] 是否有明确的边界？
- [ ] 是否通过工厂方法创建？
- [ ] 是否封装了业务规则？
- [ ] 是否只暴露必要的方法？

### 领域服务
- [ ] 是否只包含业务逻辑？
- [ ] 是否不依赖技术框架？
- [ ] 是否协调多个聚合？

### 应用服务
- [ ] 是否只负责编排？
- [ ] 是否处理了 DTO 转换？
- [ ] 是否处理了跨服务调用？

### 仓储
- [ ] 接口是否在领域层？
- [ ] 实现是否在基础设施层？
- [ ] 是否使用领域对象而非 PO？

## 重构建议

### 从传统架构迁移到 DDD

1. **识别聚合**
   - 找出需要保持一致性的对象集合
   - 定义聚合边界

2. **提取值对象**
   - 将基本类型包装为值对象
   - 添加验证规则

3. **分离领域逻辑**
   - 将业务规则从 Service 移到聚合
   - 创建领域服务处理跨聚合逻辑

4. **应用依赖倒置**
   - 在领域层定义仓储接口
   - 在基础设施层实现接口

5. **引入 DTO**
   - 分离领域对象和传输对象
   - 在应用层处理转换

## 参考资源

- [领域驱动设计](https://book.douban.com/subject/26819666/) - Eric Evans
- [实现领域驱动设计](https://book.douban.com/subject/25844633/) - Vaughn Vernon
- [Kitex 官方文档](https://www.cloudwego.io/zh/docs/kitex/)
- [Martin Fowler 的 DDD 文章](https://martinfowler.com/tags/domain%20driven%20design.html)

---

**提示**：这是一个快速参考，详细说明请查看各个文档文件。
