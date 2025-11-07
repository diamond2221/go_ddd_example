# Kitex 生成代码说明

## 关于 kitex_gen 目录

### 什么是 Kitex？
Kitex 是字节跳动开源的 Go RPC 框架，支持 Thrift 和 Protobuf。

### 生成代码的作用
`rpc_gen/kitex_gen/` 目录包含根据 Thrift IDL 自动生成的代码：
- RPC 请求/响应结构体
- 服务接口定义
- 序列化/反序列化代码
- 客户端/服务端桩代码

### 实际项目中如何生成？

在实际项目中，应该使用 Kitex 命令行工具生成这些代码：

```bash
# 安装 Kitex 工具
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

# 生成代码
kitex -module service -service recommendation idl/recommendation.thrift
```

生成的文件会包含：
- `kitex_gen/recommendation/` - Thrift 结构体定义
- `kitex_gen/recommendation/recommendationservice/` - 服务接口和实现

### 本示例项目的处理

为了让示例项目完整可运行，我手动创建了简化版的生成代码：
- `rpc_gen/kitex_gen/recommendation/recommendation.go` - RPC 数据结构
- `rpc_gen/kitex_gen/recommendation/recommendationservice.go` - 服务接口

这些文件包含了详细的中文注释，解释了：
- RPC 对象 vs 领域对象的区别
- 为什么需要分离
- 在 DDD 架构中的位置

## RPC 层在 DDD 中的位置

```
┌─────────────────────────────────────────────────────────┐
│  RPC 客户端（其他服务）                                   │
└─────────────────────────────────────────────────────────┘
                        ↓ RPC 调用
┌─────────────────────────────────────────────────────────┐
│  Kitex 框架（序列化、网络传输）                           │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│  Interface Layer - RecommendationHandler                 │
│  - 协议适配（RPC → 应用层）                               │
│  - 参数验证                                               │
│  - 错误处理                                               │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│  Application Layer - RecommendationService               │
│  - 用例编排                                               │
│  - DTO 转换                                               │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│  Domain Layer - 领域服务、聚合、值对象                    │
│  - 核心业务逻辑                                           │
└─────────────────────────────────────────────────────────┘
```

## RPC 对象 vs 领域对象 vs DTO

### 三种对象的区别

| 对象类型 | 位置 | 职责 | 示例 |
|---------|------|------|------|
| RPC 对象 | kitex_gen | 网络传输 | recommendation.UserRecommendation |
| DTO | application/dto | 应用层数据传输 | dto.UserRecommendationDTO |
| 领域对象 | domain | 业务逻辑 | aggregate.UserRecommendation |

### 为什么需要三层对象？

#### 1. RPC 对象（kitex_gen）
```go
// 由 Kitex 生成，包含序列化标签
type UserRecommendation struct {
    UserId   int64  `thrift:"user_id,1,required"`
    Username string `thrift:"username,2,required"`
    // ...
}
```
- 用途：网络传输
- 特点：包含 Thrift 标签，由工具生成
- 依赖：依赖 RPC 框架

#### 2. DTO（应用层）
```go
// 应用层定义，用于内部传输
type UserRecommendationDTO struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    // ...
}
```
- 用途：应用层和接口层之间传输数据
- 特点：简单的数据结构，可能包含 JSON 标签
- 依赖：不依赖任何框架

#### 3. 领域对象（领域层）
```go
// 领域层定义，包含业务逻辑
type UserRecommendation struct {
    id           RecommendationID
    targetUserID UserID
    score        int
    // ...
}

func (r *UserRecommendation) CalculateScore() int {
    // 业务逻辑
}
```
- 用途：表达业务概念和规则
- 特点：包含业务行为，字段私有
- 依赖：不依赖任何外层

### 数据流转过程

```
客户端请求
    ↓
RPC 对象（反序列化）
    ↓
Handler（RPC → DTO）
    ↓
DTO
    ↓
应用服务（DTO → 领域对象）
    ↓
领域对象（业务处理）
    ↓
应用服务（领域对象 → DTO）
    ↓
DTO
    ↓
Handler（DTO → RPC）
    ↓
RPC 对象（序列化）
    ↓
客户端响应
```

## 实际使用示例

### 1. 定义 Thrift IDL
```thrift
// idl/recommendation.thrift
struct UserRecommendation {
    1: required i64 user_id,
    2: required string username,
}

service RecommendationService {
    GetRecommendationsResponse GetFollowingBasedRecommendations(
        1: GetRecommendationsRequest req
    )
}
```

### 2. 生成代码
```bash
kitex -module service idl/recommendation.thrift
```

### 3. 实现 Handler
```go
// interface/handler/recommendation_handler.go
type RecommendationHandler struct {
    service *service.RecommendationService
}

func (h *RecommendationHandler) GetFollowingBasedRecommendations(
    ctx context.Context,
    req *recommendation.GetRecommendationsRequest,
) (*recommendation.GetRecommendationsResponse, error) {
    // 1. 参数验证
    // 2. 调用应用服务
    // 3. 转换响应
}
```

### 4. 启动服务
```go
// main.go
func main() {
    handler := handler.NewRecommendationHandler(...)
    svr := recommendation.NewServer(handler)
    svr.Run()
}
```

## 注意事项

### 1. 不要手动修改生成的代码
生成的代码会在重新生成时被覆盖，所以：
- ✅ 在 Handler 中实现业务逻辑
- ❌ 不要修改 kitex_gen 中的代码

### 2. 版本管理
- 生成的代码应该提交到 Git
- IDL 文件是源文件，必须提交
- 团队成员使用相同版本的 Kitex 工具

### 3. 向后兼容
修改 Thrift IDL 时注意：
- 不要删除已有字段
- 新增字段使用 optional
- 不要修改字段编号

## 扩展阅读

- [Kitex 官方文档](https://www.cloudwego.io/zh/docs/kitex/)
- [Thrift IDL 语法](https://thrift.apache.org/docs/idl)
- [DDD 中的防腐层](https://martinfowler.com/bliki/AnticorruptionLayer.html)

## 总结

Kitex 生成的代码是 RPC 层的技术实现，在 DDD 架构中：
- 属于接口层的一部分
- 负责协议适配和序列化
- 不包含业务逻辑
- 通过 Handler 连接到应用层

理解 RPC 对象、DTO、领域对象的区别和职责，是实践 DDD 的关键。
