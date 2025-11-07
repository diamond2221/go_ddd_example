# Wire 依赖注入设置指南

## 当前状态

项目已经配置好 Wire 依赖注入，但还需要运行 Wire 命令生成代码。

## 快速开始

### 1. 安装 Wire

```bash
go install github.com/google/wire/cmd/wire@latest
```

### 2. 生成依赖注入代码

```bash
# 在项目根目录运行
wire
```

成功后会看到：

```
wire: service: wrote wire_gen.go
```

### 3. 运行服务

```bash
go run .
```

## 文件说明

### 已创建的文件

1. **wire.go** - Wire 配置文件
   - 定义 Provider（如何构造对象）
   - 定义 ProviderSet（按层分组）
   - 定义 Injector（需要什么对象）

2. **main.go** - 服务启动入口
   - 使用 Wire 生成的 `InitializeRecommendationHandler()` 函数
   - 简洁的启动代码（只有几行）

3. **wire_gen.go.example** - Wire 生成代码示例
   - 展示 Wire 会生成什么代码
   - 实际的 wire_gen.go 会在运行 wire 命令后生成

### 需要生成的文件

- **wire_gen.go** - Wire 自动生成的依赖注入代码
  - 运行 `wire` 命令后自动生成
  - 不要手动编辑这个文件
  - 如果修改了 wire.go，重新运行 wire 命令

## 项目结构

```
service/
├── wire.go                    # Wire 配置（已创建）
├── wire_gen.go                # Wire 生成的代码（需要运行 wire 生成）
├── wire_gen.go.example        # 生成代码示例（参考用）
├── main.go                    # 启动入口（已更新为使用 Wire）
├── WIRE_SETUP.md              # 本文件
├── docs/
│   ├── WIRE_GUIDE.md          # Wire 完整使用指南
│   ├── WIRE_COMPARISON.md     # 手动 vs Wire 对比
│   ├── CROSS_SERVICE_CALL.md  # 跨服务调用指南
│   └── HTTP_VS_RPC.md         # HTTP vs RPC 对比
├── application/
│   └── service/
│       └── recommendation_service.go
├── domain/
│   └── service/
│       └── recommendation_generator.go
├── infrastructure/
│   ├── client/
│   │   ├── content_service_client.go      # HTTP 客户端实现
│   │   └── content_service_rpc_client.go  # RPC 客户端实现
│   └── repository/
│       └── mock_repository.go
└── interface/
    └── handler/
        └── recommendation_handler.go
```

## Wire 工作流程

```
1. 你定义 wire.go
   ├── Provider：如何构造对象
   ├── ProviderSet：按层分组
   └── Injector：需要什么对象

2. 运行 wire 命令
   └── Wire 分析依赖关系

3. Wire 生成 wire_gen.go
   ├── 自动解决依赖顺序
   ├── 生成创建和注入代码
   └── 编译时检查依赖

4. 在 main.go 中使用
   └── InitializeRecommendationHandler()
```

## 依赖注入流程

Wire 会按以下顺序创建对象：

```
1. 基础设施层
   ├── UserRPCClient
   ├── ContentServiceClient
   └── ReasonTextConfigClient

2. 仓储层
   ├── SocialGraphRepository
   └── ContentRepository

3. 领域服务层
   └── RecommendationGenerator
       ├── 依赖 SocialGraphRepository
       └── 依赖 ContentRepository

4. 应用服务层
   └── RecommendationService
       ├── 依赖 RecommendationGenerator
       ├── 依赖 SocialGraphRepository
       ├── 依赖 ContentRepository
       ├── 依赖 ContentServiceClient
       ├── 依赖 UserRPCClient
       └── 依赖 ReasonTextConfigClient

5. 接口层
   └── RecommendationHandler
       └── 依赖 RecommendationService
```

## 修改依赖

### 添加新依赖

假设要添加一个 `CacheClient`：

1. **在 wire.go 中添加 Provider**：

```go
func provideCacheClient() service.CacheClient {
    return client.NewCacheClient()
}

var infrastructureSet = wire.NewSet(
    provideUserRPCClient,
    provideContentServiceClient,
    provideReasonConfigClient,
    provideCacheClient, // 新增
)
```

2. **修改构造函数签名**（如果需要）：

```go
func NewRecommendationService(
    generator *service.RecommendationGenerator,
    socialGraphRepo repository.SocialGraphRepository,
    contentRepo repository.ContentRepository,
    contentClient ContentServiceClient,
    userRPCClient UserRPCClient,
    reasonConfigClient ReasonTextConfigClient,
    cacheClient CacheClient, // 新增
) *RecommendationService {
    // ...
}
```

3. **重新运行 Wire**：

```bash
wire
```

Wire 会自动更新 wire_gen.go，添加 `cacheClient` 的创建和传递！

### 删除依赖

1. 从 wire.go 中删除 Provider
2. 修改构造函数签名
3. 重新运行 `wire`

### 修改实现

假设要从 HTTP 客户端改为 RPC 客户端：

```go
// 修改 Provider 实现
func provideContentServiceClient() service.ContentServiceClient {
    // 旧：return client.NewContentServiceHTTPClient("http://...")
    // 新：
    return client.NewContentServiceRPCClient()
}
```

重新运行 `wire`，完成！

## 常见问题

### 1. wire 命令找不到

```bash
# 确保 GOPATH/bin 在 PATH 中
export PATH=$PATH:$(go env GOPATH)/bin

# 或者重新安装
go install github.com/google/wire/cmd/wire@latest
```

### 2. Wire 生成失败

```bash
wire: service: no provider found for service.UserRPCClient
```

**原因**：缺少 Provider

**解决**：在 wire.go 中添加对应的 Provider 函数

### 3. 循环依赖

```bash
wire: service: cycle detected in provider graph
```

**原因**：A 依赖 B，B 依赖 A

**解决**：重新设计依赖关系，避免循环

### 4. 类型不匹配

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

## 对比：手动 vs Wire

### 手动方式（已移除）

```go
func initDependencies() *Dependencies {
    // 50+ 行手动创建和注入代码
    userRPCClient := repository.NewMockUserRPCClient()
    socialGraphRepo := repository.NewMockSocialGraphRepository()
    // ...
    return &Dependencies{...}
}
```

**问题**：
- ❌ 代码冗长
- ❌ 依赖顺序容易出错
- ❌ 运行时才发现错误

### Wire 方式（当前）

```go
// wire.go
func provideUserRPCClient() service.UserRPCClient {
    return repository.NewMockUserRPCClient()
}

func InitializeRecommendationHandler() *handler.RecommendationHandler {
    wire.Build(infrastructureSet, repositorySet, ...)
    return nil
}

// main.go
func main() {
    handler := InitializeRecommendationHandler()
    // ...
}
```

**优势**：
- ✅ 代码简洁
- ✅ 自动解决依赖顺序
- ✅ 编译时检查错误

## 下一步

1. **安装 Wire**：`go install github.com/google/wire/cmd/wire@latest`
2. **生成代码**：`wire`
3. **运行服务**：`go run .`
4. **阅读文档**：
   - `docs/WIRE_GUIDE.md` - 完整使用指南
   - `docs/WIRE_COMPARISON.md` - 详细对比
   - `docs/CROSS_SERVICE_CALL.md` - 跨服务调用
   - `docs/HTTP_VS_RPC.md` - HTTP vs RPC

## 参考资料

- [Wire 官方文档](https://github.com/google/wire)
- [Wire 用户指南](https://github.com/google/wire/blob/main/docs/guide.md)
- [Wire 最佳实践](https://github.com/google/wire/blob/main/docs/best-practices.md)
