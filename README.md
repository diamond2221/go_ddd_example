# DDD 示例：基于关注的用户推荐功能

> 这是一个完整的 **Kitex 微服务项目**，使用 **DDD（领域驱动设计）** 架构实现推荐功能。

## 项目特点

- ✅ **完整的 Kitex 微服务** - 包含 main.go、Makefile、配置文件等
- ✅ **DDD 架构** - 清晰的四层架构（接口层、应用层、领域层、基础设施层）
- ✅ **详细的中文注释** - 每个文件都有详细的业务和技术说明
- ✅ **丰富的文档** - 10+ 个文档，涵盖各个方面
- ✅ **开箱即用** - 提供 Makefile、构建脚本、Docker 支持

## 快速开始

```bash
# 1. 安装工具
make install-tools

# 2. 初始化项目
make init

# 3. 编译
make build

# 4. 运行
./recommendation-service
```

详细说明请查看：[KITEX_PROJECT_GUIDE.md](KITEX_PROJECT_GUIDE.md)

## 功能描述

**业务需求**：推荐我关注的人最近关注的用户及其帖子

- 用户 A 关注了 B、C、D
- B 最近关注了 E、F
- C 最近关注了 G
- 推荐给 A：E、F、G 以及他们的帖子

## DDD 架构分层

```
recommendation/
├── domain/              # 领域层（核心业务逻辑）
│   ├── aggregate/      # 聚合根
│   ├── valueobject/    # 值对象
│   ├── service/        # 领域服务
│   └── repository/     # 仓储接口（定义）
│
├── application/        # 应用层（用例编排）
│   ├── service/       # 应用服务
│   └── dto/           # 数据传输对象
│
├── infrastructure/     # 基础设施层（技术实现）
│   └── persistence/   # 仓储实现
│
└── interface/         # 接口层（对外暴露）
    └── handler/      # RPC 处理器
```

## 核心概念

### 1. 值对象 (Value Object)
- **特点**：不可变、无唯一标识、通过值比较相等性
- **示例**：`UserID`、`RecommendationReason`
- **作用**：封装业务概念和验证规则

### 2. 聚合 (Aggregate)
- **特点**：一组相关对象的集合，有聚合根
- **示例**：`UserRecommendation`（单个推荐）、`RecommendationList`（推荐列表）
- **作用**：定义事务边界、封装业务规则

### 3. 领域服务 (Domain Service)
- **特点**：不属于任何单一实体/聚合的业务逻辑
- **示例**：`RecommendationGenerator`（推荐生成逻辑）
- **作用**：协调多个聚合完成复杂业务逻辑

### 4. 仓储 (Repository)
- **特点**：接口定义在领域层，实现在基础设施层
- **示例**：`SocialGraphRepository`、`ContentRepository`
- **作用**：提供类似集合的接口来访问聚合

### 5. 应用服务 (Application Service)
- **特点**：编排用例，协调领域对象和基础设施
- **示例**：`RecommendationService`
- **作用**：处理跨服务调用、DTO 转换、事务管理

## 依赖方向

```
Interface Layer (接口层)
    ↓
Application Layer (应用层)
    ↓
Domain Layer (领域层) ← 核心，不依赖外层
    ↑
Infrastructure Layer (基础设施层) - 实现领域层定义的接口
```

## 关键设计决策

### 1. 为什么 UserRecommendation 是聚合根？
- 它有完整的生命周期（创建、过期）
- 它封装了推荐分数计算规则
- 它保证了推荐数据的一致性

### 2. 为什么需要 RecommendationList 聚合？
- 管理推荐集合的业务规则（去重、排序、过滤）
- 定义推荐列表级别的操作（获取 Top N、移除过期）
- 保证推荐列表的完整性

### 3. 为什么 RecommendationGenerator 是领域服务？
- 推荐生成逻辑涉及多个聚合（用户、关注关系、内容）
- 不属于任何单一聚合的职责
- 包含核心业务规则，应该在领域层

### 4. 为什么仓储接口定义在领域层？
- 领域层定义"需要什么数据"（业务语言）
- 基础设施层实现"如何获取数据"（技术细节）
- 遵循依赖倒置原则（DIP）

## 业务规则示例

### 推荐分数计算
```
分数 = 关注者数量 × 10 + 帖子数量 × 2
```

### 推荐有效期
- 推荐生成后 7 天过期
- 过期推荐自动从列表中移除

### 推荐约束
- 不能推荐自己
- 不能重复推荐同一用户
- 至少要有 1 个关注者才能生成推荐

## 测试策略

### 单元测试（领域层）
```go
// 测试聚合的业务规则
func TestUserRecommendation_CalculateScore(t *testing.T) {
    // 3 个关注者，5 个帖子
    // 预期分数：3 × 10 + 5 × 2 = 40
}

func TestRecommendationList_CannotRecommendSelf(t *testing.T) {
    // 验证不能推荐自己的规则
}
```

### 集成测试（应用层）
```go
// 测试完整用例
func TestRecommendationService_GetFollowingBasedRecommendations(t *testing.T) {
    // 使用 mock 仓储测试完整流程
}
```

## 扩展点

### 1. 新增推荐策略
- 创建新的领域服务（如 `PopularityBasedGenerator`）
- 应用层可以组合多个推荐策略

### 2. 新增推荐理由类型
- 在 `RecommendationReason` 值对象中添加新的 `ReasonType`
- 实现对应的描述逻辑

### 3. 更换数据库
- 只需实现新的仓储实现类
- 领域层和应用层代码无需修改

## 与传统架构对比

### 传统三层架构
```
Controller → Service → DAO → Database
```
- 以数据库为中心
- 业务逻辑分散在 Service 层
- 难以测试和维护

### DDD 架构
```
Interface → Application → Domain ← Infrastructure
```
- 以领域模型为中心
- 业务逻辑集中在 Domain 层
- 易于测试和扩展

## 📚 文档导航

### 快速开始
- **README.md**（本文件）- 项目概览和架构说明
- **FINAL_SUMMARY.md** - 项目完成总结 🎉
- **KITEX_PROJECT_GUIDE.md** - Kitex 微服务项目指南 ⭐⭐⭐
- **QUICK_REFERENCE.md** - DDD 快速参考卡片 ⭐
- **PROJECT_STRUCTURE.md** - 项目结构详解 ⭐
- **LEARNING_PATH.md** - 详细的学习路径和检查清单
- **CHANGELOG.md** - 更新日志

### 架构文档
- **architecture-diagram.md** - DDD 架构图示
- **RPC_LAYER_ARCHITECTURE.md** - RPC 层架构详解

### 代码说明
- **CODE_COMMENTS_GUIDE.md** - 代码注释完善指南
- **KITEX_GEN_README.md** - Kitex 生成代码说明

### 项目文件
- **Makefile** - 构建命令（make help 查看所有命令）
- **main.go** - 服务启动入口
- **config/config.yaml** - 服务配置文件
- **script/bootstrap.sh** - Kitex 代码生成脚本

## Makefile 命令

```bash
make help              # 显示所有可用命令
make init              # 初始化项目（安装工具、依赖、生成代码）
make gen               # 生成 Kitex 代码
make build             # 编译服务
make run               # 运行服务
make test              # 运行所有测试
make clean             # 清理构建产物
make docker-build      # 构建 Docker 镜像
```

## 学习建议

### 学习 Kitex 微服务
1. **阅读 KITEX_PROJECT_GUIDE.md** - 了解如何使用 Kitex 框架
2. **查看 main.go** - 理解服务启动和依赖注入
3. **运行 make gen** - 了解代码生成流程
4. **修改 IDL** - 尝试添加新的 RPC 方法

### 学习 DDD 架构
1. **从领域层开始**：先理解业务规则和领域模型
2. **关注聚合边界**：理解为什么这样划分聚合
3. **理解依赖方向**：领域层不依赖任何外层
4. **实践值对象**：用值对象封装业务概念
5. **区分领域服务和应用服务**：理解它们的不同职责
6. **理解对象转换**：RPC 对象 → DTO → 领域对象的转换流程

## RPC 层说明

本项目使用 Kitex 作为 RPC 框架。关于 RPC 层的详细说明：

- **KITEX_GEN_README.md** - Kitex 生成代码的使用说明
- **RPC_LAYER_ARCHITECTURE.md** - RPC 层在 DDD 架构中的位置和作用

### RPC 对象 vs DTO vs 领域对象

| 对象类型 | 位置 | 职责 | 依赖 |
|---------|------|------|------|
| RPC 对象 | rpc_gen/kitex_gen | 网络传输 | 依赖 RPC 框架 |
| DTO | application/dto | 应用层数据传输 | 不依赖框架 |
| 领域对象 | domain | 业务逻辑 | 不依赖外层 |

## 参考资料

- 《领域驱动设计》（Eric Evans）
- 《实现领域驱动设计》（Vaughn Vernon）
- [Kitex 官方文档](https://www.cloudwego.io/zh/docs/kitex/)
- 项目重构文档：`docs/refactor/README.md`
