# 依赖注入指南

## 概述

本项目使用手动依赖注入，在 `main.go` 的 `initDependencies()` 函数中组装所有依赖。

## 依赖注入顺序

DDD 分层架构的依赖注入遵循"从内到外"的原则：

```
1. 基础设施层
   ├── 数据库连接 (DB)
   ├── 缓存连接 (Redis)
   └── RPC 客户端 (UserRPCClient, ReasonConfigClient)

2. 仓储层
   ├── SocialGraphRepository (实现领域接口)
   └── ContentRepository (实现领域接口)

3. 领域服务层
   └── RecommendationGenerator (依赖仓储接口)

4. 应用服务层
   └── RecommendationService (依赖领域服务、仓储、RPC客户端)

5. 接口层
   └── RecommendationHandler (依赖应用服务)
```

## 当前实现（Mock 版本）

```go
func initDependencies() *Dependencies {
    // 1. Mock RPC 客户端
    userRPCClient := repository.NewMockUserRPCClient()

    // 2. 配置服务客户端（可选）
    var reasonConfigClient service.ReasonTextConfigClient = nil

    // 3. Mock 仓储
    socialGraphRepo := repository.NewMockSocialGraphRepository()
    contentRepo := repository.NewMockContentRepository()

    // 4. 领域服务
    generator := domainService.NewRecommendationGenerator(
        socialGraphRepo,
        contentRepo,
    )

    // 5. 应用服务
    recommendationService := service.NewRecommendationService(
        generator,
        socialGraphRepo,
        contentRepo,
        userRPCClient,
        reasonConfigClient,
    )

    return &Dependencies{
        RecommendationService: recommendationService,
    }
}
```


## 生产环境实现

在实际项目中，需要替换 Mock 实现为真实实现：

### 1. 数据库仓储实现

```go
// infrastructure/persistence/mysql_social_graph_repository.go
type MySQLSocialGraphRepository struct {
    db *gorm.DB
}

func NewMySQLSocialGraphRepository(db *gorm.DB) repository.SocialGraphRepository {
    return &MySQLSocialGraphRepository{db: db}
}

func (r *MySQLSocialGraphRepository) GetFollowings(
    ctx context.Context,
    userID valueobject.UserID,
) ([]valueobject.UserID, error) {
    var follows []FollowPO
    err := r.db.WithContext(ctx).
        Where("follower_id = ?", userID.Value()).
        Find(&follows).Error
    if err != nil {
        return nil, err
    }

    result := make([]valueobject.UserID, 0, len(follows))
    for _, f := range follows {
        uid, _ := valueobject.NewUserID(f.FollowingID)
        result = append(result, uid)
    }
    return result, nil
}
```

### 2. RPC 客户端实现

```go
// infrastructure/client/user_rpc_client.go
type UserRPCClientImpl struct {
    client userservice.Client
}

func NewUserRPCClient(addr string) service.UserRPCClient {
    client, err := userservice.NewClient(
        "user-service",
        client.WithHostPorts(addr),
    )
    if err != nil {
        log.Fatal("Failed to create user rpc client:", err)
    }
    return &UserRPCClientImpl{client: client}
}
```

### 3. 配置服务客户端

```go
// 启用配置服务
reasonConfigClient := client.NewReasonTextConfigHTTPClient(
    "http://config-service:8080",
)
```

### 4. 完整的 initDependencies

```go
func initDependencies() *Dependencies {
    // 1. 加载配置
    cfg := config.Load()

    // 2. 初始化数据库
    db := initDB(cfg.Database)

    // 3. 初始化 RPC 客户端
    userRPCClient := client.NewUserRPCClient(cfg.UserService.Addr)

    // 4. 初始化配置服务客户端
    var reasonConfigClient service.ReasonTextConfigClient
    if cfg.Features.UseReasonConfig {
        reasonConfigClient = client.NewReasonTextConfigHTTPClient(
            cfg.ConfigService.URL,
        )
    }

    // 5. 创建仓储
    socialGraphRepo := persistence.NewMySQLSocialGraphRepository(db)
    contentRepo := persistence.NewMySQLContentRepository(db)

    // 6. 创建领域服务
    generator := domainService.NewRecommendationGenerator(
        socialGraphRepo,
        contentRepo,
    )

    // 7. 创建应用服务
    recommendationService := service.NewRecommendationService(
        generator,
        socialGraphRepo,
        contentRepo,
        userRPCClient,
        reasonConfigClient,
    )

    return &Dependencies{
        RecommendationService: recommendationService,
    }
}
```

## 使用依赖注入框架

对于大型项目，推荐使用依赖注入框架：

### Wire (Google)

```go
// wire.go
//go:build wireinject
// +build wireinject

func InitializeApp() (*Dependencies, error) {
    wire.Build(
        // 基础设施
        initDB,
        initRedis,

        // 仓储
        persistence.NewMySQLSocialGraphRepository,
        persistence.NewMySQLContentRepository,

        // 客户端
        client.NewUserRPCClient,
        client.NewReasonTextConfigHTTPClient,

        // 领域服务
        domainService.NewRecommendationGenerator,

        // 应用服务
        service.NewRecommendationService,

        // 组装
        wire.Struct(new(Dependencies), "*"),
    )
    return nil, nil
}
```

### Fx (Uber)

```go
func main() {
    app := fx.New(
        // 提供依赖
        fx.Provide(
            initDB,
            persistence.NewMySQLSocialGraphRepository,
            persistence.NewMySQLContentRepository,
            client.NewUserRPCClient,
            domainService.NewRecommendationGenerator,
            service.NewRecommendationService,
            handler.NewRecommendationHandler,
        ),
        // 启动服务
        fx.Invoke(startServer),
    )
    app.Run()
}
```

## 测试中的依赖注入

单元测试时，使用 Mock 实现：

```go
func TestRecommendationService(t *testing.T) {
    // 使用 Mock 仓储
    mockSocialGraphRepo := repository.NewMockSocialGraphRepository()
    mockContentRepo := repository.NewMockContentRepository()
    mockUserRPC := repository.NewMockUserRPCClient()

    // 创建领域服务
    generator := domainService.NewRecommendationGenerator(
        mockSocialGraphRepo,
        mockContentRepo,
    )

    // 创建应用服务
    svc := service.NewRecommendationService(
        generator,
        mockSocialGraphRepo,
        mockContentRepo,
        mockUserRPC,
        nil, // 不使用配置服务
    )

    // 测试
    result, err := svc.GetFollowingBasedRecommendations(
        context.Background(),
        1,
        10,
    )
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## 最佳实践

1. **接口在领域层定义**：保持依赖倒置
2. **实现在基础设施层**：技术细节不污染领域层
3. **构造函数注入**：明确依赖关系
4. **避免全局变量**：通过参数传递依赖
5. **单一职责**：每个组件只做一件事
6. **可测试性**：所有依赖都可以被 Mock 替换
