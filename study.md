# DDD 模式特点：

## domain 领域层（核心业务逻辑）
- #### valueobject 值对象
```zsh
什么需要定义成值对象？
· 有明确的业务规则（长度限制、字符限制）
· 验证逻辑复杂，需要集中管理
· 昵称在多处使用，避免到处重复验证
· 保证数据有效性，创建后就是合法的昵称
```
```go
type Nickname struct {
    value string
}

func NewNickname(value string) (*Nickname, error) {
    if len(value) > 10 {
        return nil, errors.New("nickname too long")
    }
    return &Nickname{value: value}, nil
}
```
- #### entity 实体
```zsh
什么是实体？
· 实体是有唯一标识的领域对象，即使属性相同，只要 ID 不同就是不同的对象。
```
```zsh
跟数据库实体的区别？
· 位于领域层（Domain Layer），代表业务概念
· 包含业务逻辑和行为方法
· 关注业务规则和不变性约束
· 使用值对象（Value Object）来封装、表达业务概念
· 独立于持久化技术
· 字段私有，通过方法访问（封装性）
· 不关心数据库细节
```
```go
type User struct {
    id       int64
    nickname string
}

func NewUser(id int64, nickname string) *User {
    return &User{id: id, nickname: nickname}
}
func (u *User) GetId() int64 {
    return u.id
}
func (u *User) GetNickname() string {
    return u.nickname
}
```

- #### aggregate 聚合根
```zsh
什么是聚合根？
聚合根是实体的集合，聚合根是聚合的入口，聚合根的聚合，聚合根的根。
· 1. 完整的生命周期：推荐有创建、使用、过期的完整生命周期
· 2. 业务规则封装：推荐分数计算、过期判断等规则都在这里
· 3. 一致性边界：推荐的所有数据必须保持一致（分数、理由、时间等）
· 4. 事务边界：对推荐的修改应该在一个事务内完成
```
```go
type Recommend struct {
    id      int64
    score   int64
    reason  string
    expired int64
}

func NewRecommend(id int64, score int64, reason string, expired int64) *Recommend {
    return &Recommend{id: id, score: score, reason: reason, expired: expired}
}
func (r *Recommend) IsValidAndScore() bool {
    return r.score > 90 && r.expired > time.format("2026-01-01").Unix()
}
func (r *Recommend) GetId() int64 {
    return r.id
}

```

- #### service 服务
```zsh
什么是服务？
· 服务是领域服务，领域服务是领域模型的业务逻辑，领域服务是领域模型的业务逻辑，领域服务是领域模型的业务逻辑。
```
- #### repository 仓储
```zsh
什么是仓储？
· 仓储是领域模型的持久化，仓储是领域模型的持久化，仓储是领域模型的持久化。
```

## application 应用层（用例编排）

