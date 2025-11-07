# Kitex 微服务项目指南

## 项目概述

这是一个基于 **Kitex** 框架和 **DDD** 架构的完整微服务项目示例。

### 技术栈

- **RPC 框架**: Kitex (字节跳动开源)
- **IDL**: Thrift
- **数据库**: MySQL (GORM)
- **架构模式**: DDD (领域驱动设计)
- **语言**: Go 1.22+

## 完整项目结构

```
recommendation/
├── main.go                      # 服务启动入口 ⭐
├── Makefile                     # 构建命令 ⭐
├── build.sh                     # 构建脚本 ⭐
├── Dockerfile                   # Docker 镜像 ⭐
├── .gitignore                   # Git 忽略文件 ⭐
├── go.mod                       # Go 模块定义 ⭐
├── go.sum                       # Go 依赖锁定
│
├── config/                      # 配置文件 ⭐
│   └── config.yaml              # 服务配置
│
├── script/                      # 脚本 ⭐
│   └── bootstrap.sh             # Kitex 代码生成脚本
│
├── idl/                         # IDL 定义
│   └── recommendation.thrift    # Thrift IDL
│
├── rpc_gen/                     # RPC 生成代码
│   └── kitex_gen/               # Kitex 生成的代码
│       └── recommendation/
│           ├── recommendation.go
│           └── recommendationservice/
│
├── domain/                      # 领域层
│   ├── aggregate/               # 聚合
│   ├── entity/                  # 实体
│   ├── valueobject/             # 值对象
│   ├── service/                 # 领域服务
│   └── repository/              # 仓储接口
│
├── application/                 # 应用层
│   ├── service/                 # 应用服务
│   └── dto/                     # DTO
│
├── infrastructure/              # 基础设施层
│   └── persistence/             # 持久化
│
├── interface/                   # 接口层
│   └── handler/                 # RPC 处理器
│
└── tests/                       # 测试
    ├── unit/                    # 单元测试
    └── integration/             # 集成测试
```

## 快速开始

### 1. 安装依赖工具

```bash
# 安装 Kitex 工具
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

# 安装 Thriftgo
go install github.com/cloudwego/thriftgo@latest

# 验证安装
kitex --version
thriftgo --version
```

### 2. 初始化项目

```bash
# 方式一：使用 Makefile（推荐）
make init

# 方式二：手动执行
go mod download
bash script/bootstrap.sh
```

### 3. 生成 Kitex 代码

```bash
# 使用 Makefile
make gen

# 或直接运行脚本
bash script/bootstrap.sh
```

这会根据 `idl/recommendation.thrift` 生成：
- `rpc_gen/kitex_gen/recommendation/*.go` - Thrift 结构体
- `rpc_gen/kitex_gen/recommendation/recommendationservice/*.go` - 服务接口

### 4. 编译项目

```bash
# 使用 Makefile
make build

# 或使用构建脚本
bash build.sh

# 或直接编译
go build -o recommendation-service .
```

### 5. 运行服务

```bash
# 方式一：直接运行
./recommendation-service

# 方式二：使用 go run
make run

# 方式三：开发模式（需要安装 air）
make dev
```

服务将在 `:8888` 端口启动。

## Makefile 命令

```bash
make help              # 显示所有可用命令
make gen               # 生成 Kitex 代码
make build             # 编译服务
make run               # 运行服务
make test              # 运行所有测试
make test-unit         # 运行单元测试
make test-integration  # 运行集成测试
make test-coverage     # 生成测试覆盖率报告
make lint              # 代码检查
make fmt               # 格式化代码
make clean             # 清理构建产物
make docker-build      # 构建 Docker 镜像
make docker-run        # 运行 Docker 容器
make install-tools     # 安装开发工具
make init              # 初始化项目
```

## 配置说明

### config/config.yaml

配置文件包含：
- 服务配置（端口、注册中心）
- 数据库配置（MySQL）
- Redis 配置
- RPC 客户端配置
- 业务配置
- 日志配置
- 监控配置
- 限流和熔断配置

在实际项目中，通常会：
1. 使用配置中心（Apollo、Nacos）
2. 支持环境变量覆盖
3. 支持多环境配置（dev、test、prod）

## 依赖注入

### main.go 中的依赖注入

```go
func main() {
    // 1. 初始化依赖
    deps := initDependencies()

    // 2. 创建 Handler
    handler := handler.NewRecommendationHandler(deps.Service)

    // 3. 创建 Kitex Server
    svr := recommendationservice.NewServer(handler)

    // 4. 启动服务
    svr.Run()
}

func initDependencies() *Dependencies {
    // 初始化数据库
    db := initDB()

    // 创建仓储（基础设施层）
    socialGraphRepo := persistence.NewSocialGraphRepository(db)
    contentRepo := persistence.NewContentRepository(db)

    // 创建领域服务（领域层）
    generator := domainservice.NewRecommendationGenerator(
        socialGraphRepo,
        contentRepo,
    )

    // 创建应用服务（应用层）
    service := appservice.NewRecommendationService(
        generator,
        socialGraphRepo,
        contentRepo,
        userRPCClient,
    )

    return &Dependencies{Service: service}
}
```

### 依赖注入的好处

1. **控制反转**：依赖由外部注入
2. **易于测试**：可以注入 mock 对象
3. **解耦**：各层不直接依赖具体实现
4. **灵活配置**：可以根据环境注入不同实现

## Kitex 服务配置

### 基本配置

```go
svr := recommendationservice.NewServer(
    handler,
    server.WithServiceAddr(&net.TCPAddr{Port: 8888}),
)
```

### 完整配置示例

```go
import (
    "github.com/cloudwego/kitex/server"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/pkg/limit"
    etcd "github.com/kitex-contrib/registry-etcd"
    prometheus "github.com/kitex-contrib/monitor-prometheus"
)

func main() {
    // 服务注册
    r, _ := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})

    // 监控
    tracer := prometheus.NewServerTracer(":9091", "/metrics")

    svr := recommendationservice.NewServer(
        handler,
        // 服务地址
        server.WithServiceAddr(&net.TCPAddr{Port: 8888}),

        // 服务注册
        server.WithRegistry(r),

        // 服务信息
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: "recommendation-service",
        }),

        // 监控
        server.WithTracer(tracer),

        // 限流
        server.WithLimit(&limit.Option{
            MaxConnections: 10000,
            MaxQPS:         1000,
        }),

        // 中间件
        server.WithMiddleware(LogMiddleware),
        server.WithMiddleware(RecoveryMiddleware),
    )

    svr.Run()
}
```

## 中间件示例

### 日志中间件

```go
func LogMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) error {
        start := time.Now()
        err := next(ctx, req, resp)
        duration := time.Since(start)

        log.Printf("Method: %s, Duration: %v, Error: %v",
            rpcinfo.GetRPCInfo(ctx).To().Method(),
            duration,
            err,
        )

        return err
    }
}
```

### 恢复中间件

```go
func RecoveryMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) (err error) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic recovered: %v", r)
                err = fmt.Errorf("internal server error")
            }
        }()
        return next(ctx, req, resp)
    }
}
```

## 客户端调用示例

### 创建客户端

```go
import (
    "service/rpc_gen/kitex_gen/recommendation"
    "service/rpc_gen/kitex_gen/recommendation/recommendationservice"
    "github.com/cloudwego/kitex/client"
)

func main() {
    // 创建客户端
    client, err := recommendationservice.NewClient(
        "recommendation-service",
        client.WithHostPorts("127.0.0.1:8888"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 调用服务
    req := &recommendation.GetRecommendationsRequest{
        UserId: 123,
        Limit:  10,
    }

    resp, err := client.GetFollowingBasedRecommendations(
        context.Background(),
        req,
    )
    if err != nil {
        log.Fatal(err)
    }

    // 处理响应
    for _, rec := range resp.Recommendations {
        fmt.Printf("User: %s, Reason: %s\n",
            rec.Username, rec.Reason)
    }
}
```

### 客户端配置

```go
client, err := recommendationservice.NewClient(
    "recommendation-service",
    // 服务发现
    client.WithResolver(r),

    // 超时配置
    client.WithRPCTimeout(3 * time.Second),
    client.WithConnectTimeout(1 * time.Second),

    // 重试配置
    client.WithFailureRetry(retry.NewFailurePolicy()),

    // 负载均衡
    client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),

    // 熔断
    client.WithCircuitBreaker(circuitbreak.NewCBSuite(...)),
)
```

## Docker 部署

### 构建镜像

```bash
# 使用 Makefile
make docker-build

# 或直接使用 Docker
docker build -t recommendation-service:latest .
```

### 运行容器

```bash
# 使用 Makefile
make docker-run

# 或直接使用 Docker
docker run -p 8888:8888 recommendation-service:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  recommendation:
    build: .
    ports:
      - "8888:8888"
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: recommendation
    ports:
      - "3306:3306"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

## 测试

### 单元测试

```bash
# 运行所有单元测试
make test-unit

# 运行特定包的测试
go test ./domain/aggregate/... -v

# 运行特定测试
go test ./domain/aggregate -run TestUserRecommendation_CalculateScore -v
```

### 集成测试

```bash
# 运行集成测试
make test-integration

# 生成覆盖率报告
make test-coverage
```

### 测试示例

```go
// 单元测试（领域层）
func TestUserRecommendation_CalculateScore(t *testing.T) {
    reason := valueobject.NewFollowedByFollowingReason(
        []valueobject.UserID{user1, user2, user3},
    )

    rec, err := aggregate.NewUserRecommendation(
        targetUser,
        reason,
        5, // 5个帖子
    )

    assert.NoError(t, err)
    assert.Equal(t, 40, rec.Score()) // 3*10 + 5*2 = 40
}

// 集成测试（应用层）
func TestRecommendationService_GetRecommendations(t *testing.T) {
    // 使用 mock 仓储
    mockRepo := &MockSocialGraphRepository{
        followings: []UserID{user1, user2},
    }

    service := NewRecommendationService(mockRepo, ...)

    resp, err := service.GetFollowingBasedRecommendations(
        context.Background(),
        123,
        10,
    )

    assert.NoError(t, err)
    assert.NotEmpty(t, resp.Recommendations)
}
```

## 监控和日志

### Prometheus 监控

```go
import prometheus "github.com/kitex-contrib/monitor-prometheus"

tracer := prometheus.NewServerTracer(":9091", "/metrics")
svr := recommendationservice.NewServer(
    handler,
    server.WithTracer(tracer),
)
```

访问 `http://localhost:9091/metrics` 查看指标。

### 日志

```go
import "github.com/cloudwego/kitex/pkg/klog"

klog.Info("Service started")
klog.Errorf("Error: %v", err)
```

## 性能优化

### 1. 连接池配置

```go
// 数据库连接池
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 2. 批量查询

```go
// 批量获取用户信息，避免 N+1 查询
userInfos, err := userRPCClient.GetUserInfoBatch(ctx, userIDs)
```

### 3. 缓存

```go
// Redis 缓存推荐结果
func (s *Service) GetRecommendations(ctx context.Context, userID int64) {
    // 先查缓存
    cached, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        return unmarshal(cached)
    }

    // 缓存未命中，查询数据库
    result := s.queryDB(ctx, userID)

    // 写入缓存
    s.redis.Set(ctx, cacheKey, marshal(result), 5*time.Minute)

    return result
}
```

## 常见问题

### Q1: 如何修改 IDL？

1. 修改 `idl/recommendation.thrift`
2. 运行 `make gen` 重新生成代码
3. 更新 Handler 实现
4. 运行 `go mod tidy`

### Q2: 如何添加新的 RPC 方法？

1. 在 Thrift IDL 中添加方法定义
2. 运行 `make gen`
3. 在 Handler 中实现新方法
4. 更新应用服务和领域服务

### Q3: 如何连接服务注册中心？

```go
import etcd "github.com/kitex-contrib/registry-etcd"

r, _ := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
svr := recommendationservice.NewServer(
    handler,
    server.WithRegistry(r),
)
```

### Q4: 如何实现链路追踪？

```go
import (
    "github.com/kitex-contrib/tracer-opentracing"
    "github.com/uber/jaeger-client-go"
)

tracer, closer := jaeger.NewTracer(...)
defer closer.Close()

svr := recommendationservice.NewServer(
    handler,
    server.WithSuite(tracer_opentracing.NewDefaultServerSuite()),
)
```

## 最佳实践

### 1. 项目组织
- ✅ 按 DDD 分层组织代码
- ✅ 使用 Makefile 管理命令
- ✅ 配置文件与代码分离

### 2. 代码质量
- ✅ 编写单元测试和集成测试
- ✅ 使用 golangci-lint 检查代码
- ✅ 保持测试覆盖率 > 80%

### 3. 性能优化
- ✅ 使用连接池
- ✅ 批量查询避免 N+1
- ✅ 合理使用缓存

### 4. 可观测性
- ✅ 添加日志
- ✅ 添加监控指标
- ✅ 添加链路追踪

### 5. 安全性
- ✅ 参数验证
- ✅ 错误处理
- ✅ 限流和熔断

## 参考资源

- [Kitex 官方文档](https://www.cloudwego.io/zh/docs/kitex/)
- [Thrift IDL 语法](https://thrift.apache.org/docs/idl)
- [DDD 实践指南](./QUICK_REFERENCE.md)
- [项目结构说明](./PROJECT_STRUCTURE.md)

## 总结

这是一个完整的 Kitex 微服务项目，包含：
- ✅ 完整的项目结构
- ✅ Makefile 和构建脚本
- ✅ 配置文件
- ✅ Docker 支持
- ✅ DDD 架构
- ✅ 详细的文档

通过这个项目，你可以学习到：
- 如何使用 Kitex 框架
- 如何实践 DDD 架构
- 如何组织微服务项目
- 如何进行依赖注入
- 如何部署和监控服务
