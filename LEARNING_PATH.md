# DDD 学习路径

## 建议学习顺序

### 第一步：理解核心概念（从 README.md 开始）
1. 阅读 `README.md` 了解整体架构
2. 理解四层架构的职责划分
3. 理解依赖方向：`Interface → Application → Domain ← Infrastructure`

### 第二步：从领域层开始（自底向上）

#### 1. 值对象 (Value Object)
**文件顺序**：
- `domain/valueobject/user_id.go` - 最简单的值对象
- `domain/valueobject/recommendation_id.go` - UUID 类型的值对象
- `domain/valueobject/recommendation_reason.go` - 复杂值对象

**关键点**：
- 不可变性（所有字段私有）
- 通过值比较相等性
- 封装验证规则
- 工厂方法创建

#### 2. 实体 (Entity)
**文件**：
- `domain/entity/post.go`

**关键点**：
- 有唯一标识（ID）
- 通过 ID 比较相等性
- 与值对象的区别

#### 3. 聚合 (Aggregate)
**文件顺序**：
- `domain/aggregate/user_recommendation.go` - 单个推荐聚合
- `domain/aggregate/recommendation_list.go` - 推荐列表聚合

**关键点**：
- 聚合根的概念
- 业务规则封装
- 事务边界
- 对外暴露有限接口

#### 4. 仓储接口 (Repository Interface)
**文件**：
- `domain/repository/social_graph_repository.go`
- `domain/repository/content_repository.go`

**关键点**：
- 接口定义在领域层
- 使用领域对象而非数据库模型
- 依赖倒置原则

#### 5. 领域服务 (Domain Service)
**文件**：
- `domain/service/_generator.go`

**关键点**：
- 跨聚合的业务逻辑
- 纯业务逻辑，不涉及技术细节
- 与应用服务的区别

### 第三步：基础设施层（技术实现）

**文件顺序**：
- `infrastructure/persistence/social_graph_repository_impl.go`
- `infrastructure/persistence/content_repository_impl.go`

**关键点**：
- 实现领域层定义的接口
- PO（持久化对象）与领域对象的转换
- 数据库访问细节

### 第四步：应用层（用例编排）

**文件顺序**：
- `application/dto/recommendation_dto.go` - 数据传输对象
- `application/service/_service.go` - 应用服务

**关键点**：
- 用例编排
- 跨服务调用
- DTO 转换
- 与领域服务的区别

### 第五步：接口层（协议适配）

**文件**：
- `interface/handler/recommendation_handler.go`
- `idl/recommendation.thrift`

**关键点**：
- 协议适配（RPC 请求/响应转换）
- 参数验证
- 调用应用服务

### 第六步：测试（验证理解）

**文件顺序**：
- `tests/unit/user_recommendation_test.go` - 聚合单元测试
- `tests/unit/recommendation_list_test.go` - 聚合单元测试
- `tests/integration/recommendation_service_test.go` - 集成测试

### 第七步：架构图示（整体理解）

**文件**：
- `architecture-diagram.md`

**关键点**：
- 整体架构图
- 数据流转
- 聚合边界
- 设计模式

## 学习检查清单

### 理解值对象
- [ ] 为什么值对象是不可变的？
- [ ] 值对象和实体的区别是什么？
- [ ] 为什么要用值对象包装基本类型（如 int64）？

### 理解聚合
- [ ] 什么是聚合根？
- [ ] 聚合的边界如何划分？
- [ ] 为什么要通过聚合根修改内部状态？
- [ ] UserRecommendation 和 RecommendationList 为什么是两个独立聚合？

### 理解领域服务
- [ ] 领域服务和应用服务的区别是什么？
- [ ] 什么样的逻辑应该放在领域服务中？
- [ ] RecommendationGenerator 为什么是领域服务？

### 理解仓储模式
- [ ] 为什么仓储接口定义在领域层？
- [ ] PO 和领域对象的区别是什么？
- [ ] 为什么要分离 PO 和领域对象？

### 理解依赖方向
- [ ] 为什么领域层不依赖任何外层？
- [ ] 依赖倒置原则是什么？
- [ ] 如何实现依赖倒置？

### 理解分层职责
- [ ] 接口层的职责是什么？
- [ ] 应用层的职责是什么？
- [ ] 领域层的职责是什么？
- [ ] 基础设施层的职责是什么？

## 实践建议

### 1. 动手实践
- 尝试添加新的推荐策略（如基于热度的推荐）
- 尝试添加新的值对象（如 ContentType）
- 尝试修改推荐分数计算规则

### 2. 对比传统架构
- 思考如果用传统三层架构如何实现
- 对比两种架构的优缺点
- 理解 DDD 解决了什么问题

### 3. 阅读源码
- 阅读你们项目的 `app/engagement/` 目录
- 对比示例代码和实际项目
- 思考如何改进现有代码

### 4. 画图理解
- 画出聚合边界图
- 画出依赖关系图
- 画出数据流转图

## 常见问题

### Q1: 什么时候用值对象，什么时候用实体？
**A**: 如果对象有唯一标识且生命周期重要，用实体；否则用值对象。

### Q2: 聚合应该多大？
**A**: 尽可能小，只包含必须保持一致性的对象。

### Q3: 什么逻辑放在领域层，什么放在应用层？
**A**: 纯业务规则放领域层，用例编排和技术细节放应用层。

### Q4: 是否所有项目都适合 DDD？
**A**: 不是。简单的 CRUD 项目用传统架构更合适。DDD 适合复杂业务领域。

### Q5: DDD 会增加代码量吗？
**A**: 会。但换来的是更清晰的结构和更好的可维护性。

## 进阶阅读

### 书籍
- 《领域驱动设计》（Eric Evans）- DDD 经典
- 《实现领域驱动设计》（Vaughn Vernon）- 实践指南

### 在线资源
- Martin Fowler 的博客（关于 DDD 的文章）
- DDD Community 网站

### 项目文档
- `docs/refactor/README.md` - 项目重构文档
- `AGENTS.md` - 开发规范
