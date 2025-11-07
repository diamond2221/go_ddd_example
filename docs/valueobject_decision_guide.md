# 值对象决策指南

## 快速决策树

```
昵称需要验证吗？
│
├─ 只需要唯一性（查数据库）
│  └─ ❌ 不要定义值对象
│     ✅ 直接用 string
│     ✅ 在领域服务中验证唯一性
│
├─ 只需要格式验证（长度、字符类型等）
│  └─ ✅ 定义值对象
│     ✅ 在值对象中验证格式
│
└─ 既需要格式验证，又需要唯一性
   └─ ✅ 定义值对象（验证格式）
      ✅ 在领域服务中验证唯一性
```

## 三种场景对比

### 场景1：只需要唯一性

**需求：** 昵称没有格式限制，但必须唯一

**❌ 错误做法：**
```go
// 不要这样做！值对象不应该依赖仓储
type Nickname struct {
    value string
}

func NewNickname(value string, repo UserRepository) (Nickname, error) {
    exists, _ := repo.ExistsByNickname(value)
    if exists {
        return Nickname{}, errors.New("昵称已存在")
    }
    return Nickname{value: value}, nil
}
```

**✅ 正确做法：**
```go
// 直接用 string
type User struct {
    nickname string  // 直接用 string
}

// 在领域服务中验证唯一性
type UserService struct {
    userRepo UserRepository
}

func (s *UserService) CreateUser(nickname string) (*User, error) {
    // 验证唯一性
    exists, _ := s.userRepo.ExistsByNickname(nickname)
    if exists {
        return nil, errors.New("昵称已被使用")
    }

    user := NewUser(nickname)
    return user, s.userRepo.Save(user)
}
```

---

### 场景2：只需要格式验证

**需求：** 昵称必须 3-16 个字符，只能是中文/英文/数字

**✅ 正确做法：**
```go
// 定义值对象
type Nickname struct {
    value string
}

func NewNickname(value string) (Nickname, error) {
    // 验证长度
    if len(value) < 3 || len(value) > 16 {
        return Nickname{}, errors.New("长度不合法")
    }

    // 验证字符格式
    if !nicknamePattern.MatchString(value) {
        return Nickname{}, errors.New("格式不合法")
    }

    return Nickname{value: value}, nil
}

// 使用
type User struct {
    nickname Nickname  // 使用值对象
}

func CreateUser(nicknameStr string) (*User, error) {
    // 格式验证在这里完成
    nickname, err := NewNickname(nicknameStr)
    if err != nil {
        return nil, err
    }

    user := NewUser(nickname)
    return user, nil
}
```

---

### 场景3：既需要格式验证，又需要唯一性（最常见）

**需求：** 昵称必须 3-16 个字符，只能是中文/英文/数字，且必须唯一

**✅ 正确做法：分层处理**
```go
// 1. 值对象：负责格式验证
type Nickname struct {
    value string
}

func NewNickname(value string) (Nickname, error) {
    // 只验证格式，不验证唯一性
    if len(value) < 3 || len(value) > 16 {
        return Nickname{}, errors.New("长度不合法")
    }
    if !nicknamePattern.MatchString(value) {
        return Nickname{}, errors.New("格式不合法")
    }
    return Nickname{value: value}, nil
}

// 2. 领域服务：负责唯一性验证
type UserService struct {
    userRepo UserRepository
}

func (s *UserService) CreateUser(nicknameStr string) (*User, error) {
    // 步骤1：格式验证（值对象负责）
    nickname, err := NewNickname(nicknameStr)
    if err != nil {
        return nil, fmt.Errorf("昵称格式不合法: %w", err)
    }

    // 步骤2：唯一性验证（领域服务负责）
    exists, _ := s.userRepo.ExistsByNickname(nickname.Value())
    if exists {
        return nil, errors.New("昵称已被使用")
    }

    // 步骤3：创建用户
    user := NewUser(nickname)
    return user, s.userRepo.Save(user)
}
```

---

## 核心原则

### 值对象的职责边界

**✅ 值对象应该负责：**
- 格式验证（长度、字符类型、正则匹配）
- 业务计算（如 RecommendationReason 的 Weight()）
- 类型安全（如 UserID vs PostID）
- 不变性保证

**❌ 值对象不应该负责：**
- 唯一性验证（需要查数据库）
- 与其他对象的关系验证
- 依赖外部服务的验证
- 任何需要 I/O 操作的验证

### 判断标准

问自己一个问题：

> **"这个验证规则可以在不访问数据库的情况下完成吗？"**

- **可以** → 放在值对象中 ✅
- **不可以** → 放在领域服务或仓储中 ✅

---

## 实际例子

| 验证规则 | 是否定义值对象 | 原因 |
|---------|--------------|------|
| 昵称长度 3-16 个字符 | ✅ 是 | 不需要访问数据库 |
| 昵称只能是中文/英文/数字 | ✅ 是 | 不需要访问数据库 |
| 昵称必须唯一 | ❌ 否 | 需要查询数据库 |
| 邮箱格式验证 | ✅ 是 | 不需要访问数据库 |
| 邮箱必须唯一 | ❌ 否 | 需要查询数据库 |
| 手机号格式验证 | ✅ 是 | 不需要访问数据库 |
| 用户ID必须是正数 | ✅ 是 | 不需要访问数据库 |
| 密码强度验证 | ✅ 是 | 不需要访问数据库 |
| 年龄必须 18-100 岁 | ✅ 是 | 不需要访问数据库 |

---

## 为什么要这样分层？

### 1. 单一职责原则
- **值对象**：负责自身的有效性
- **领域服务**：负责业务规则和对象间的关系
- **仓储**：负责数据持久化

### 2. 可测试性
```go
// 值对象的测试：不需要数据库
func TestNewNickname(t *testing.T) {
    nickname, err := NewNickname("张三123")
    assert.NoError(t, err)
}

// 领域服务的测试：可以 mock 仓储
func TestCreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    // ...
}
```

### 3. 性能优化
```go
// 快速失败：格式错误的数据不会查询数据库
nickname, err := NewNickname("a")  // 立即失败，不查数据库
if err != nil {
    return err  // 格式不合法，快速返回
}

// 只有格式正确的数据才会查询数据库
exists, _ := repo.ExistsByNickname(nickname.Value())
```

### 4. 可维护性
- 修改格式规则 → 只需改值对象
- 修改唯一性验证逻辑 → 只需改领域服务
- 职责清晰，不会混乱

---

## 总结

**记住这个简单的规则：**

```
需要查数据库？
├─ 是 → 不要放在值对象中，放在领域服务中
└─ 否 → 可以放在值对象中
```

**唯一性验证永远不应该在值对象中！**
