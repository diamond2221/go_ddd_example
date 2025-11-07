# Wire 依赖注入指南

## Wire 是什么？

Wire 是 Google 开发的依赖注入工具，通过**代码生成**的方式实现依赖注入。

### 核心概念

```
┌─────────────────────────────────────────────────────┐
│ 你定义：                                             │
│ 1. Provider：如何构造对象                            │
│ 2. Injector：需要什么对象                            │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│ Wire 生成：                                          │
│ - 分析依赖关系                                       │
│ - 生成依赖注入代码                                   │
│ - 编译时检查依赖                                     │
└─────────────────────────────────────────────────────┘
```

## 为什么用 Wire？

### 对比：手动依赖注入 vs Wire

#### 手动方式（main.go）

```go
func initDependencies() *Dependencies {
    // 1. 创建基础设施
    userRPCClient := repository.NewMockUserRPCClient()
    var reasonConfigClient service.ReasonTextConfigClient = nil

    // 2. 创建仓储
    socialGraphRepo := repository.NewMockSocialGraphRepository()
    contentRepo := repository.NewMockContentRepository()

    // 3. 创建领域服务
    generator := domainService.NewRecommendationGenerator(
        socialGraphRepo,
        contentRepo,
    )

    // 4. 创建应用服务
    recommendationService := service.NewRecommendationService(
        generator,
        socialGraphRepo,
        contentRepo,
        client.NewContentServiceRPCClient(),
        userRPCClient,
        reasonConfigClient,
    )

    return &Dependencies{
        RecommendationService: recommendationService,
    }
}
```

**问题**：
- ❌ 代码冗长（100+ 行）
- ❌ 依赖顺序容易出错
- ❌ 运行时才发现依赖错误
- ❌ 添加新依赖需要修改多处

#### Wire 方式（wire.go）

```go
// 1. 定义 Provider（如何构造）
func provideUserRPCClient() service.UserRPCClient {
    return repository.NewMockUserRPCClient()
}

func provideSocialGraphRepository() repository.SocialGraphRepository {
    return repository.NewMockSocialGraphRepository()
}

// ... 其他 Provider

// 2. 定义 Injector（需要什么）
func InitializeRecommendationHandler() *handler.RecommendationHandler {
    wire.Build(
        infrastructureSet,
        repositorySet,
        domainServiceSet,
        applicationServiceSet,
        handlerSet,
    )
    return nil
}
```

**优势**：
- ✅ 代码简洁（20 行）
- ✅ Wire 自动解决依赖顺序
- ✅ 编译时检查依赖
- ✅ 添加新依赖只需添加 Provider

## 快速开始

### 1. 安装 Wire

```bash
go install github.com/google/wire/cmd/wire@latest
```

### 2. 创建 wire.go

```go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
)

// 定义 Provider
func provideUserRPCClient() service.UserRPCClient {
    return repository.NewMockUserRPCClient()
}

// 定义 ProviderSet
var infrastructureSet = wire.NewSet(
    provideUserRPCClient,
    // ... 其他 Provider
)

// 定义 Injector
func InitializeRecommendationHandler() *handler.RecommendationHandler {
    wire.Build(
        infrastructureSet,
        repositorySet,
        domainServiceSet,
        applicationServiceSet,
        handlerSet,
    )
    return nil
}
```

### 3. 生成代码

```bash
# 在项目根目录运行
wire

# 输出：
# wire: service: wrote wire_gen.go
```

### 4. 使用生成的代码

```go
// main.go
func main() {
    // 使用 Wire 生成的函数
    handler := InitializeRecommendationHandler()

    // 启动服务
    svr := recommendationservice.NewServer(handler)
    svr.Run()
}
```

## 核心概念详解

### 1. Provider

Provider 是一个函数，告诉 Wire 如何构造某个对象。

```go
// 简单 Provider：无依赖
func provideUserRPCClient() service.UserRPCClient {
    return repository.NewMockUserRPCClient()
}

// 复杂 Provider：有依赖
func provideRecommendationService(
    generator *domainService.RecommendationGenerator,
    socialGraphRepo repository.SocialGraphRepository,
    contentRepo repository.ContentRepository,
    contentClient service.ContentServiceClient,
    userRPCClient service.UserRPCClient,
    reasonConfigClient service.ReasonTextConfigClient,
) *service.RecommendationService {
    return service.NewRecommendationService(
        generator,
        socialGraphRepo,
        contentRepo,
        contentClient,
        userRPCClient,
        reasonConfigClient,
    )
}

// 实际上，如果构造函数签名和 Provider 一样，可以直接使用构造函数：
var applicationServiceSet = wire.NewSet(
    service.NewRecommendationService, // 直接使用构造函数
)
```

### 2. ProviderSet

ProviderSet 是一组 Provider 的集合。

```go
// 按层分组
var infrastructureSet = wire.NewSet(
    provideUserRPCClient,
    provideContentServiceClient,
    provideReasonConfigClient,
)

var repositorySet = wire.NewSet(
    provideSocialGraphRepository,
    provideContentRepository,
)

var domainServiceSet = wire.NewSet(
    domainService.NewRecommendationGenerator,
)

var applicationServiceSet = wire.NewSet(
    service.NewRecommendationService,
)

var handlerSet = wire.NewSet(
    handler.NewRecommendationHandler,
)
```

**为什么要分组？**
- 按 DDD 分层组织
- 易于管理和复用
- 可以在不同 Injector 中组合使用

### 3. Injector

Injector 是一个函数签名，告诉 Wire 你需要什么对象。

```go
// Injector 函数
func InitializeRecommendationHandler() *handler.RecommendationHandler {
    wire.Build(
        infrastructureSet,
        repositorySet,
        domainServiceSet,
        applicationServiceSet,
        handlerSet,
    )
    return nil // 占位返回，Wire 会生成真实实现
}
```

**Wire 生成的代码**（wire_gen.go）：

```go
func InitializeRecommendationHandler() *handler.RecommendationHandler {
    // 1. 基础设施层
    userRPCClient := provideUserRPCClient()
    contentServiceClient := provideContentServiceClient()
    reasonTextConfigClient := provideReasonConfigClient()

    // 2. 仓储层
    socialGraphRepository := provideSocialGraphRepository()
    contentRepository := provideContentRepository()

    // 3. 领域服务层
    recommendationGenerator := domainService.NewRecommendationGenerator(
        socialGraphRepository,
        contentRepository,
    )

    // 4. 应用服务层
    recommendationService := service.NewRecommendationService(
        recommendationGenerator,
        socialGraphRepository,
        contentRepository,
        contentServiceClient,
        userRPCClient,
        reasonTextConfigClient,
    )

    // 5. 接口层
    recommendationHandler := handler.NewRecommendationHandler(
        recommendationService,
    )

    return recommendationHandler
}
```

## 高级用法

### 1. 条件注入

根据配置决定使用哪个实现：

```go
// 定义配置
type Config struct {
    UseRPC bool
}

// 根据配置提供不同实现
func provideContentServiceClient(cfg *Config) service.ContentServiceClient {
    if cfg.UseRPC {
        return client.NewContentServiceRPCClient()
    }
    return client.NewContentServiceHTTPClient("http://content-service:8080")
}

// Injector
func InitializeRecommendationHandler(cfg *Config) *handler.RecommendationHandler {
    wire.Build(
        // 传入配置
        wire.Value(cfg),

        // 其他 Provider
        infrastructureSet,
        repositorySet,
        domainServiceSet,
        applicationServiceSet,
        handlerSet,
    )
    return nil
}

// 使用
func main() {
    cfg := &Config{UseRPC: true}
    handler := InitializeRecommendationHandler(cfg)
    // ...
}
```

### 2. 可选依赖

某些依赖可以为 nil：

```go
// Provider 返回 nil
func provideReasonConfigClient() service.ReasonTextConfigClient {
    // 不使用配置服务
    return nil
}

// 或者根据配置决定
func provideReasonConfigClient(cfg *Config) service.ReasonTextConfigClient {
    if !cfg.UseReasonConfig {
        return nil
    }
    return client.NewReasonTextConfigHTTPClient(cfg.ReasonConfigURL)
}
```

### 3. 接口绑定

当 Provider 返回具体类型，但需要接口时：

```go
// 具体类型
type MySQLSocialGraphRepository struct {
    db *gorm.DB
}

func (r *MySQLSocialGraphRepository) GetFollowing(...) { ... }

// Provider 返回具体类型
func provideMySQLSocialGraphRepository(db *gorm.DB) *MySQLSocialGraphRepository {
    return &MySQLSocialGraphRepository{db: db}
}

// 绑定到接口
var repositorySet = wire.NewSet(
    provideMySQLSocialGraphRepository,
    wire.Bind(new(repository.SocialGraphRepository), new(*MySQLSocialGraphRepository)),
)
```

### 4. 清理资源

某些对象需要清理（如数据库连接）：

```go
// 返回对象和清理函数
func provideDatabase(cfg *Config) (*gorm.DB, func(), error) {
    db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
    if err != nil {
        return nil, nil, err
    }

    cleanup := func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }

    return db, cleanup, nil
}

// Injector 返回清理函数
func InitializeRecommendationHandler() (*handler.RecommendationHandler, func(), error) {
    wire.Build(
        infrastructureSet,
        repositorySet,
        domainServiceSet,
        applicationServiceSet,
        handlerSet,
    )
    return nil, nil, nil
}

// 使用
func main() {
    handler, cleanup, err := InitializeRecommendationHandler()
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup() // 清理资源

    // 启动服务
    svr := recommendationservice.NewServer(handler)
    svr.Run()
}
```

### 5. 测试中使用 Wire

定义测试专用的 Injector：

```go
// wire_test.go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
)

// 测试用的 ProviderSet（使用 mock）
var testInfrastructureSet = wire.NewSet(
    provideMockUserRPCClient,
    provideMockContentServiceClient,
    wire.Value(service.ReasonTextConfigClient(nil)),
)

// 测试用的 Injector
func InitializeTestHandler() *handler.RecommendationHandler {
    wire.Build(
        testInfrastructureSet, // 使用 mock
        repositorySet,
        domainServiceSet,
        applicationServiceSet,
        handlerSet,
    )
    return nil
}
```

```go
// handler_test.go
func TestRecommendationHandler(t *testing.T) {
    // 使用测试专用的 Injector
    handler := InitializeTestHandler()

    // 测试逻辑
    // ...
}
```

## 实际项目示例

### 项目结构

```
service/
├── wire.go                    # Wire 配置
├── wire_gen.go                # Wire 生成的代码（自动生成）
├── main.go                    # 使用 Wire 的启动入口
├── config/
│   └── config.go              # 配置定义
├── application/
│   └── service/
│       └── recommendation_service.go
├── domain/
│   └── service/
│       └── recommendation_generator.go
├── infrastructure/
│   ├── client/
│   │   ├── content_service_client.go
│   │   └── user_rpc_client.go
│   └── repository/
│       ├── social_graph_repository.go
│       └── content_repository.go
└── interface/
    └── handler/
        └── recommendation_handler.go
```

### wire.go（完整示例）

```go
//go:build wireinject
// +build wireinject

package main

import (
    "service/application/service"
    "service/config"
    domainService "service/domain/service"
    "service/infrastructure/client"
    "service/infrastructure/persistence"
    "service/interface/handler"

    "github.com/google/wire"
)

// 配置 Provider
func provideConfig() (*config.Config, error) {
    return config.Load("config.yaml")
}

// 数据库 Provider
func provideDatabase(cfg *config.Config) (*gorm.DB, func(), error) {
    db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
    if err != nil {
        return nil, nil, err
    }

    cleanup := func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }

    return db, cleanup, nil
}

// RPC 客户端 Provider
func provideUserRPCClient(cfg *config.Config) (service.UserRPCClient, error) {
    return client.NewUserRPCClient(cfg.UserService.Addr)
}

func provideContentServiceClient(cfg *config.Config) service.ContentServiceClient {
    if cfg.ContentService.UseRPC {
        return client.NewContentServiceRPCClient(cfg.ContentService.Addr)
    }
    return client.NewContentServiceHTTPClient(cfg.ContentService.URL)
}

func provideReasonConfigClient(cfg *config.Config) service.ReasonTextConfigClient {
    if !cfg.Features.UseReasonConfig {
        return nil
    }
    return client.NewReasonTextConfigHTTPClient(cfg.ReasonConfigService.URL)
}

// 仓储 Provider
func provideSocialGraphRepository(db *gorm.DB) repository.SocialGraphRepository {
    return persistence.NewMySQLSocialGraphRepository(db)
}

func provideContentRepository(db *gorm.DB) repository.ContentRepository {
    return persistence.NewMySQLContentRepository(db)
}

// ProviderSet
var infrastructureSet = wire.NewSet(
    provideConfig,
    provideDatabase,
    provideUserRPCClient,
    provideContentServiceClient,
    provideReasonConfigClient,
)

var repositorySet = wire.NewSet(
    provideSocialGraphRepository,
    provideContentRepository,
)

var domainServiceSet = wire.NewSet(
    domainService.NewRecommendationGenerator,
)

var applicationServiceSet = wire.NewSet(
    service.NewRecommendationService,
)

var handlerSet = wire.NewSet(
    handler.NewRecommendationHandler,
)

// Injector
func InitializeRecommendationHandler() (*handler.RecommendationHandler, func(), error) {
    wire.Build(
        infrastructureSet,
        repositorySet,
        domainServiceSet,
        applicationServiceSet,
        handlerSet,
    )
    return nil, nil, nil
}
```

### main.go（使用 Wire）

```go
package main

import (
    "log"
    "net"

    "service/rpc_gen/kitex_gen/recommendation/recommendationservice"
    "github.com/cloudwego/kitex/server"
)

func main() {
    // 使用 Wire 初始化
    handler, cleanup, err := InitializeRecommendationHandler()
    if err != nil {
        log.Fatal("Initialize failed:", err)
    }
    defer cleanup()

    // 创建服务
    svr := recommendationservice.NewServer(
        handler,
        server.WithServiceAddr(&net.TCPAddr{
            IP:   net.IPv4(0, 0, 0, 0),
            Port: 8888,
        }),
    )

    // 启动服务
    log.Println("Recommendation Service starting on :8888")
    if err := svr.Run(); err != nil {
        log.Fatal("Server run failed:", err)
    }
}
```

## 常见问题

### 1. Wire 生成失败

```bash
wire: service: no provider found for service.UserRPCClient
```

**原因**：缺少 Provider

**解决**：添加 Provider 函数

```go
func provideUserRPCClient() service.UserRPCClient {
    return repository.NewMockUserRPCClient()
}

var infrastructureSet = wire.NewSet(
    provideUserRPCClient, // 添加到 ProviderSet
)
```

### 2. 循环依赖

```bash
wire: service: cycle detected in provider graph
```

**原因**：A 依赖 B，B 依赖 A

**解决**：重新设计依赖关系，避免循环

### 3. 类型不匹配

```bash
wire: service: provider returns *MySQLRepository, but *Repository is needed
```

**原因**：Provider 返回具体类型，但需要接口

**解决**：使用 wire.Bind

```go
var repositorySet = wire.NewSet(
    provideMySQLRepository,
    wire.Bind(new(repository.Repository), new(*MySQLRepository)),
)
```

### 4. 多个 Provider 返回同一类型

```bash
wire: service: multiple providers for string
```

**原因**：多个 Provider 返回相同类型（如 string）

**解决**：使用结构体包装

```go
type DatabaseDSN string
type RedisDSN string

func provideDatabaseDSN() DatabaseDSN {
    return "mysql://..."
}

func provideRedisDSN() RedisDSN {
    return "redis://..."
}
```

## 总结

### Wire vs 手动依赖注入

| 维度 | 手动方式 | Wire |
|------|---------|------|
| **代码量** | 多（100+ 行） | 少（20 行） |
| **错误检查** | 运行时 | 编译时 |
| **依赖顺序** | 手动管理 | 自动解决 |
| **可维护性** | 低 | 高 |
| **学习成本** | 低 | 中等 |
| **性能** | 好 | 好（无反射） |

### 何时使用 Wire？

✅ **推荐使用**：
- 中大型项目（依赖关系复杂）
- 微服务架构（多个服务）
- 需要编译时检查
- 团队协作（依赖关系清晰）

❌ **不推荐使用**：
- 小型项目（依赖简单）
- 快速原型（手动注入更快）
- 团队不熟悉 Wire

### 核心要点

1. **Provider**：告诉 Wire 如何构造对象
2. **ProviderSet**：按层组织 Provider
3. **Injector**：告诉 Wire 需要什么对象
4. **编译时检查**：依赖错误在编译时发现
5. **代码生成**：生成的代码是普通 Go 代码

### 下一步

1. 安装 Wire：`go install github.com/google/wire/cmd/wire@latest`
2. 创建 wire.go
3. 运行 `wire` 生成代码
4. 在 main.go 中使用生成的函数
5. 享受自动依赖注入的便利！
