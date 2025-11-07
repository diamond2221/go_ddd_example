package main

import (
	"errors"
	"fmt"
)

// ============================================================
// 场景对比：昵称只需要唯一性，没有格式限制
// ============================================================

// ❌ 错误做法：在值对象中验证唯一性
type WrongNickname struct {
	value string
}

// 这样做是错误的！值对象不应该依赖仓储
// func NewWrongNickname(value string, repo UserRepository) (WrongNickname, error) {
//     exists, _ := repo.ExistsByNickname(value)
//     if exists {
//         return WrongNickname{}, errors.New("昵称已存在")
//     }
//     return WrongNickname{value: value}, nil
// }

// ✅ 正确做法：直接用 string，在领域服务中验证唯一性

// 1. 定义仓储接口
type UserRepository interface {
	ExistsByNickname(nickname string) (bool, error)
	Save(user *User) error
}

// 2. 定义用户实体（昵称直接用 string）
type User struct {
	id       string
	nickname string // 直接用 string，因为没有格式限制
	email    string
}

func NewUser(id, nickname, email string) *User {
	return &User{
		id:       id,
		nickname: nickname,
		email:    email,
	}
}

// 3. 定义领域服务（负责唯一性验证）
type UserService struct {
	userRepo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

// 在领域服务中验证唯一性
func (s *UserService) CreateUser(nickname, email string) (*User, error) {
	// 验证昵称唯一性
	exists, err := s.userRepo.ExistsByNickname(nickname)
	if err != nil {
		return nil, fmt.Errorf("检查昵称唯一性失败: %w", err)
	}
	if exists {
		return nil, errors.New("昵称已被使用")
	}

	// 创建用户
	user := NewUser(generateID(), nickname, email)

	// 保存到数据库
	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("保存用户失败: %w", err)
	}

	return user, nil
}

// ============================================================
// 模拟实现（用于演示）
// ============================================================

type InMemoryUserRepository struct {
	users map[string]*User // key: nickname
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryUserRepository) ExistsByNickname(nickname string) (bool, error) {
	_, exists := r.users[nickname]
	return exists, nil
}

func (r *InMemoryUserRepository) Save(user *User) error {
	r.users[user.nickname] = user
	return nil
}

func generateID() string {
	return "user_123" // 简化实现
}

// ============================================================
// 使用示例
// ============================================================

func main() {
	// 创建仓储和服务
	repo := NewInMemoryUserRepository()
	service := NewUserService(repo)

	// 场景1：创建第一个用户（成功）
	user1, err := service.CreateUser("张三", "zhangsan@example.com")
	if err != nil {
		fmt.Printf("❌ 创建用户失败: %v\n", err)
	} else {
		fmt.Printf("✅ 创建用户成功: %s\n", user1.nickname)
	}

	// 场景2：尝试创建相同昵称的用户（失败）
	user2, err := service.CreateUser("张三", "lisi@example.com")
	if err != nil {
		fmt.Printf("❌ 创建用户失败: %v\n", err) // 预期会失败
	} else {
		fmt.Printf("✅ 创建用户成功: %s\n", user2.nickname)
	}

	// 场景3：创建不同昵称的用户（成功）
	user3, err := service.CreateUser("李四", "lisi@example.com")
	if err != nil {
		fmt.Printf("❌ 创建用户失败: %v\n", err)
	} else {
		fmt.Printf("✅ 创建用户成功: %s\n", user3.nickname)
	}

	fmt.Println("\n============================================================")
	fmt.Println("总结：")
	fmt.Println("- 昵称只需要唯一性 → 直接用 string，在领域服务中验证")
	fmt.Println("- 昵称有格式限制 → 定义值对象，在值对象中验证格式")
	fmt.Println("- 昵称既有格式限制又需要唯一性 → 值对象验证格式 + 领域服务验证唯一性")
	fmt.Println("============================================================")
}
