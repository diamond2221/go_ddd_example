# 项目结构说明

## 目录结构

```
recommendation/
├── domain/                          # 领域层（核心业务逻辑）
│   ├── aggregate/                   # 聚合根
│   │   ├── user_recommendation.go   # 用户推荐聚合
│   │   └── recommendation_list.go   # 推荐列表聚合
│   ├── entity/                      # 实体
│   │   └── post.go                  # 帖子实体
│   ├── valueobject/                 # 值对象
│   │   ├── user_id.go               # 用户ID值对象
│   │   ├── post_id.go               # 帖子ID值对象
│   │   ├── recommendation_id.go     # 推荐ID值对象
│   │   └── recommendation_reason.go # 推荐理由值对象
│   ├── service/                     # 领域服务
│   │   └── recommendation_generator.go  # 推荐生成服务
│   └── repository/                  # 仓储接口（定义）
│       ├── social_graph_repository.go   # 社交关系仓储接口
│       └── content_repository.go        # 内容仓储接口
│
├── application/                     # 应用层（用例编排）
│   ├── service/                     # 应用服务
│   │   └── recommendation_service.go    # 推荐应用服务
│   └── dto/                         # 数据传输对象
│       └── recommendation_dto.go        # 推荐DTO
│
├── infrastructure/                  # 基础设施层（技术实现）
│   └── persistence/                 # 持久化
│       ├── social_graph_repository_impl.go  # 社交关系仓储实现
│       └── content_repository_impl.go       # 内容仓储实现
│
├── interface/                       # 接口层（对外暴露）
│   └── handler/                     # RPC处理器
│       └── recommendation_handler.go    # 推荐RPC处理器
│
├── rpc_gen/                         # RPC生成代码
│   └── kitex_gen/                   # Kitex生成的代码
│       └── recommendation/          # 推荐服务
│           ├── recommendation.go        # RPC数据结构
│           └── recommendationservice.go # RPC服务接口
│
├── idl/                             # 接口定义语言
│   └── recommendation.thrift        # Thrift IDL定义
│
├── tests/                           # 测试
│   ├── unit/                        # 单元测试
│   │   ├── user_recommendation_test.go
│   │   └── recommendation_list_test.go
│   └── integration/                 # 集成测试
│       └── recommendation_service_test.go
│
├── docs/                            # 文档（本项目的文档在根目录）
│
├── go.mod                           # Go模块定义
├── go.sum                           # Go依赖锁定
│
└── 文档文件
    ├── README.md                    # 项目说明
    ├── QUICK_REFERENCE.md           # 快速参考
    ├── LEARNING_PATH.md             # 学习路径
    ├── SUMMARY.md                   # 项目总结
    ├── architecture-diagram.md      # 架构图示
    ├── RPC_LAYER_ARCHITECTURE.md    # RPC层架构
    ├── CODE_COMMENTS_GUIDE.md       # 代码注释指南
    ├── KITEX_GEN_README.md          # Kitex说明
    └── PROJECT_STRUCTURE.md         # 本文件
```

## 分层说明

### 1. 领域层（domain/）

**职责**：核心业务逻辑，不依赖任何外层

#### aggregate/ - 聚合根
- 定义事务边界和一致性边界
- 封装业务规则
- 对外暴露有限接口

**文件**：
- `user_recommendation.go` - 单个推荐的聚合根
- `recommendation_list.go` - 推荐列表的聚合根

#### entity/ - 实体
- 有唯一标识
- 有生命周期
- 通过ID比较相等性

**文件**：
- `post.go` - 帖子实体

#### valueobject/ - 值对象
- 不可变
- 无唯一标识
- 通过值比较相等性
- 封装验证规则

**文件**：
- `user_id.go` - 用户ID
- `post_id.go` - 帖子ID
- `recommendation_id.go` - 推荐ID
- `recommendation_reason.go` - 推荐理由

#### service/ - 领域服务
- 跨聚合的业务逻辑
- 不属于任何单一实体
- 纯业务规则

**文件**：
- `recommendation_generator.go` - 推荐生成逻辑

#### repository/ - 仓储接口
- 接口定义在领域层
- 实现在基础设施层
- 依赖倒置原则

**文件**：
- `social_graph_repository.go` - 社交关系仓储接口
- `content_repository.go` - 内容仓储接口

### 2. 应用层（application/）

**职责**：用例编排，协调领域对象和基础设施

#### service/ - 应用服务
- 编排用例
- 跨服务调用
- DTO转换
- 事务管理

**文件**：
- `recommendation_service.go` - 推荐应用服务

#### dto/ - 数据传输对象
- 用于数据传输
- 简单的数据结构
- 不包含业务逻辑

**文件**：
- `recommendation_dto.go` - 推荐DTO

### 3. 基础设施层（infrastructure/）

**职责**：技术实现，实现领域层定义的接口

#### persistence/ - 持久化
- 实现仓储接口
- 数据库访问
- PO ↔ 领域对象转换

**文件**：
- `social_graph_repository_impl.go` - 社交关系仓储实现
- `content_repository_impl.go` - 内容仓储实现

### 4. 接口层（interface/）

**职责**：协议适配，对外暴露接口

#### handler/ - RPC处理器
- 协议适配
- 参数验证
- 调用应用服务
- RPC对象 ↔ DTO转换

**文件**：
- `recommendation_handler.go` - 推荐RPC处理器

### 5. RPC生成代码（rpc_gen/）

**职责**：RPC框架生成的代码

#### kitex_gen/recommendation/
- RPC数据结构
- 服务接口定义
- 序列化/反序列化代码

**文件**：
- `recommendation.go` - RPC数据结构
- `recommendationservice.go` - RPC服务接口

### 6. IDL定义（idl/）

**职责**：接口定义语言

**文件**：
- `recommendation.thrift` - Thrift IDL定义

### 7. 测试（tests/）

**职责**：测试代码

#### unit/ - 单元测试
- 测试领域对象的业务规则
- 不依赖外部资源

#### integration/ - 集成测试
- 测试完整用例
- 使用mock仓储

## 依赖关系图

```
┌─────────────────────────────────────────────────────────┐
│                    interface/                            │
│                  (接口层)                                 │
│  - handler/recommendation_handler.go                     │
└─────────────────────────────────────────────────────────┘
                        ↓ 依赖
┌─────────────────────────────────────────────────────────┐
│                   application/                           │
│                  (应用层)                                 │
│  - service/recommendation_service.go                     │
│  - dto/recommendation_dto.go                             │
└─────────────────────────────────────────────────────────┘
                        ↓ 依赖
┌─────────────────────────────────────────────────────────┐
│                     domain/                              │
│                   (领域层 - 核心)                         │
│  - aggregate/                                            │
│  - entity/                                               │
│  - valueobject/                                          │
│  - service/                                              │
│  - repository/ (接口定义)                                │
└─────────────────────────────────────────────────────────┘
                        ↑ 依赖倒置
┌─────────────────────────────────────────────────────────┐
│                 infrastructure/                          │
│                 (基础设施层)                              │
│  - persistence/ (仓储实现)                               │
└─────────────────────────────────────────────────────────┘
```

## 文件职责详解

### 领域层文件

#### user_recommendation.go
```go
// 职责：
// - 定义用户推荐聚合根
// - 封装推荐分数计算规则
// - 管理推荐的生命周期（创建、过期）
// - 保证推荐数据的一致性

// 关键方法：
// - NewUserRecommendation() - 工厂方法
// - CalculateScore() - 分数计算
// - IsExpired() - 过期判断
// - Refresh() - 刷新推荐
```

#### recommendation_list.go
```go
// 职责：
// - 管理推荐集合
// - 封装列表级别的业务规则
// - 去重、排序、过滤

// 关键方法：
// - AddRecommendation() - 添加推荐（含验证）
// - GetTopN() - 获取Top N
// - RemoveExpired() - 移除过期
// - FilterByMinScore() - 过滤低分
```

#### recommendation_generator.go
```go
// 职责：
// - 实现推荐算法
// - 协调多个聚合
// - 生成推荐列表

// 关键方法：
// - GenerateFollowingBasedRecommendations() - 基于关注的推荐
// - GeneratePopularityBasedRecommendations() - 基于热度的推荐
```

### 应用层文件

#### recommendation_service.go
```go
// 职责：
// - 编排推荐用例
// - 调用领域服务
// - 跨服务调用（user服务、content服务）
// - DTO转换

// 关键方法：
// - GetFollowingBasedRecommendations() - 获取推荐用例
// - getUserInfoMap() - 批量获取用户信息
// - convertPostsToDTO() - 转换帖子为DTO
```

### 基础设施层文件

#### social_graph_repository_impl.go
```go
// 职责：
// - 实现社交关系仓储接口
// - 数据库访问（GORM）
// - PO ↔ 领域对象转换

// 关键方法：
// - GetFollowings() - 获取关注列表
// - GetRecentFollowings() - 获取最近关注
// - IsFollowing() - 检查关注关系
```

### 接口层文件

#### recommendation_handler.go
```go
// 职责：
// - 实现RPC服务接口
// - 协议适配（RPC ↔ 应用层）
// - 参数验证
// - 错误处理

// 关键方法：
// - GetFollowingBasedRecommendations() - RPC方法实现
// - convertToRPCResponse() - DTO → RPC对象
// - convertPostsToRPC() - 帖子DTO → RPC对象
```

## 数据流转示例

### 完整的请求-响应流程

```
1. 客户端发起RPC调用
   ↓
2. Kitex框架接收请求
   ↓
3. recommendation_handler.go
   - 接收RPC请求对象
   - 参数验证
   - 调用应用服务
   ↓
4. recommendation_service.go
   - 转换为领域对象
   - 调用领域服务
   - 跨服务调用（user服务）
   - 组装响应
   - 转换为DTO
   ↓
5. recommendation_generator.go
   - 执行推荐算法
   - 调用仓储获取数据
   - 创建推荐聚合
   - 返回推荐列表
   ↓
6. social_graph_repository_impl.go
   - 查询数据库
   - PO → 领域对象
   - 返回数据
   ↓
7. 逐层返回
   - 领域对象 → DTO
   - DTO → RPC对象
   - 序列化
   - 返回客户端
```

## 命名规范

### 包命名
- 小写，单数形式
- 例如：`aggregate`, `service`, `repository`

### 文件命名
- 小写，下划线分隔
- 例如：`user_recommendation.go`, `recommendation_service.go`

### 类型命名
- 大驼峰（PascalCase）
- 例如：`UserRecommendation`, `RecommendationService`

### 方法命名
- 大驼峰（公开方法）
- 小驼峰（私有方法）
- 例如：`GetFollowings()`, `calculateScore()`

### 接口命名
- 名词 + 功能
- 例如：`SocialGraphRepository`, `ContentRepository`

### 实现命名
- 接口名 + Impl
- 例如：`SocialGraphRepositoryImpl`

## 扩展指南

### 添加新的推荐策略

1. 在 `domain/service/recommendation_generator.go` 添加新方法
2. 实现推荐算法逻辑
3. 在 `application/service/recommendation_service.go` 添加对应的用例
4. 在 `interface/handler/recommendation_handler.go` 添加RPC方法

### 添加新的值对象

1. 在 `domain/valueobject/` 创建新文件
2. 定义值对象结构
3. 实现工厂方法和验证规则
4. 添加访问器方法

### 添加新的聚合

1. 在 `domain/aggregate/` 创建新文件
2. 定义聚合根结构
3. 实现工厂方法
4. 封装业务规则
5. 添加业务行为方法

### 添加新的仓储

1. 在 `domain/repository/` 定义接口
2. 在 `infrastructure/persistence/` 实现接口
3. 定义PO结构
4. 实现数据访问逻辑
5. 实现PO ↔ 领域对象转换

## 最佳实践

### 1. 保持领域层纯粹
- ✅ 只包含业务逻辑
- ❌ 不依赖外部框架
- ❌ 不包含技术细节

### 2. 明确分层职责
- ✅ 每层有明确的职责
- ❌ 不跨层调用
- ❌ 不在错误的层实现逻辑

### 3. 正确的依赖方向
- ✅ 外层依赖内层
- ✅ 领域层不依赖外层
- ✅ 使用依赖倒置

### 4. 合理的对象转换
- ✅ 在正确的层做转换
- ✅ 不暴露内部实现
- ✅ 保持对象职责单一

### 5. 充分的测试覆盖
- ✅ 单元测试领域逻辑
- ✅ 集成测试完整用例
- ✅ 使用mock隔离依赖

## 总结

这个项目结构遵循了DDD的核心原则：
- ✅ 清晰的分层架构
- ✅ 正确的依赖方向
- ✅ 领域模型为中心
- ✅ 业务逻辑集中在领域层
- ✅ 技术细节隔离在基础设施层

通过这样的结构，项目具有：
- 高内聚低耦合
- 易于测试
- 易于维护
- 易于扩展
