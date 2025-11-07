//go:build wireinject
// +build wireinject

package main

import (
	"service/application/service"
	domainService "service/domain/service"
	"service/infrastructure/repository"
	"service/interface/handler"

	"github.com/google/wire"
)

// Wire 依赖注入配置文件
//
// Wire 是什么？
// Wire 是 Google 开发的依赖注入工具，通过代码生成的方式实现依赖注入。
//
// 为什么用 Wire？
// 1. 编译时检查：依赖错误在编译时发现，不是运行时
// 2. 无反射：生成的代码是普通 Go 代码，性能好
// 3. 易于理解：生成的代码可读性强
// 4. 类型安全：利用 Go 的类型系统
//
// Wire 工作流程：
// 1. 你定义 Provider（如何构造对象）
// 2. 你定义 Injector（需要什么对象）
// 3. Wire 生成代码（自动解决依赖关系）
//
// 对比手动依赖注入：
// ┌─────────────────────────────────────────────────────┐
// │ 手动方式（main.go 的 initDependencies）              │
// │ - 手动创建每个对象                                   │
// │ - 手动传递依赖                                       │
// │ - 依赖顺序容易出错                                   │
// │ - 代码冗长                                           │
// └─────────────────────────────────────────────────────┘
//
// ┌─────────────────────────────────────────────────────┐
// │ Wire 方式（wire.go）                                 │
// │ - 定义 Provider（如何构造）                          │
// │ - Wire 自动生成代码                                  │
// │ - 编译时检查依赖                                     │
// │ - 代码简洁                                           │
// └─────────────────────────────────────────────────────┘

// ProviderSet 定义：Provider 集合
//
// Provider 是什么？
// Provider 是一个函数，告诉 Wire 如何构造某个对象。
//
// ProviderSet 是什么？
// ProviderSet 是一组 Provider 的集合，可以被其他 ProviderSet 引用。
//
// 为什么要分组？
// - 按层分组：基础设施层、仓储层、领域层、应用层
// - 易于管理：每层的依赖清晰
// - 易于复用：可以在不同的 Injector 中复用

// infrastructureSet 基础设施层 Provider
//
// 包含：
// - RPC 客户端（User 服务、Content 服务、配置服务）
// - 数据库连接（实际项目中）
// - Redis 连接（实际项目中）
var infrastructureSet = wire.NewSet(
	// RPC 客户端
	provideUserRPCClient,
	provideContentServiceClient,
	provideReasonConfigClient,

	// 实际项目中还会有：
	// provideDatabase,
	// provideRedis,
	// provideKafka,
)

// repositorySet 仓储层 Provider
//
// 包含：
// - SocialGraphRepository
// - ContentRepository
var repositorySet = wire.NewSet(
	provideSocialGraphRepository,
	provideContentRepository,
)

// domainServiceSet 领域服务层 Provider
//
// 包含：
// - RecommendationGenerator（推荐生成器）
var domainServiceSet = wire.NewSet(
	domainService.NewRecommendationGenerator,
)

// applicationServiceSet 应用服务层 Provider
//
// 包含：
// - RecommendationService（推荐应用服务）
var applicationServiceSet = wire.NewSet(
	service.NewRecommendationService,
)

// handlerSet 接口层 Provider
//
// 包含：
// - RecommendationHandler（RPC Handler）
var handlerSet = wire.NewSet(
	handler.NewRecommendationHandler,
)

// Provider 函数定义
//
// 这些函数告诉 Wire 如何构造每个对象。
// Wire 会分析这些函数的参数和返回值，自动解决依赖关系。

// provideUserRPCClient 提供 User RPC 客户端
//
// 实际项目中，这里会：
// - 读取配置文件
// - 创建真实的 RPC 客户端
// - 配置超时、重试等
//
// 示例：
//
//	func provideUserRPCClient(cfg *Config) service.UserRPCClient {
//	    client, err := userservice.NewClient(
//	        cfg.UserService.Name,
//	        client.WithHostPorts(cfg.UserService.Addr),
//	    )
//	    if err != nil {
//	        panic(err)
//	    }
//	    return client
//	}
func provideUserRPCClient() service.UserRPCClient {
	// 示例：使用 mock 实现
	return repository.NewMockUserRPCClient()
}

// provideContentServiceClient 提供 Content 服务客户端
//
// 这里展示了如何在不同环境使用不同实现：
// - 开发环境：使用 mock
// - 测试环境：使用 HTTP 客户端
// - 生产环境：使用 RPC 客户端
//
// 实际项目中，通过配置文件控制：
//
//	func provideContentServiceClient(cfg *Config) service.ContentServiceClient {
//	    switch cfg.ContentService.Type {
//	    case "rpc":
//	        return client.NewContentServiceRPCClient(...)
//	    case "http":
//	        return client.NewContentServiceHTTPClient(cfg.ContentService.URL)
//	    default:
//	        return nil // 使用本地数据库
//	    }
//	}
func provideContentServiceClient() service.ContentServiceClient {
	// 示例：返回 nil，使用本地数据库
	// 如果需要使用远程服务，可以改为：
	// return client.NewContentServiceHTTPClient("http://content-service:8080")
	// 或：
	// return client.NewContentServiceRPCClient()
	return nil
}

// provideReasonConfigClient 提供推荐理由配置服务客户端
//
// 这是一个可选的依赖（可以为 nil）。
//
// 实际项目中：
//
//	func provideReasonConfigClient(cfg *Config) service.ReasonTextConfigClient {
//	    if !cfg.Features.UseReasonConfig {
//	        return nil // 不使用配置服务
//	    }
//	    return client.NewReasonTextConfigHTTPClient(cfg.ReasonConfigService.URL)
//	}
func provideReasonConfigClient() service.ReasonTextConfigClient {
	// 示例：不使用配置服务
	return nil
}

// provideSocialGraphRepository 提供社交图谱仓储
//
// 实际项目中：
//
//	func provideSocialGraphRepository(db *gorm.DB) repository.SocialGraphRepository {
//	    return persistence.NewMySQLSocialGraphRepository(db)
//	}
func provideSocialGraphRepository() repository.SocialGraphRepository {
	// 示例：使用 mock 实现
	return repository.NewMockSocialGraphRepository()
}

// provideContentRepository 提供内容仓储
//
// 实际项目中：
//
//	func provideContentRepository(db *gorm.DB) repository.ContentRepository {
//	    return persistence.NewMySQLContentRepository(db)
//	}
func provideContentRepository() repository.ContentRepository {
	// 示例：使用 mock 实现
	return repository.NewMockContentRepository()
}

// Injector 函数定义
//
// Injector 是一个函数签名，告诉 Wire 你需要什么对象。
// Wire 会生成这个函数的实现，自动解决所有依赖。

// InitializeRecommendationHandler 初始化推荐 Handler
//
// 这是 Wire 的 Injector 函数。
//
// 工作原理：
// 1. 你定义函数签名（返回 *handler.RecommendationHandler）
// 2. 你指定使用哪些 ProviderSet
// 3. Wire 生成函数实现（在 wire_gen.go）
//
// Wire 会自动：
// - 分析依赖关系
// - 按正确顺序调用 Provider
// - 传递依赖参数
// - 返回最终对象
//
// 依赖链：
// RecommendationHandler
//
//	↓ 依赖
//
// RecommendationService
//
//	↓ 依赖
//
// RecommendationGenerator + SocialGraphRepo + ContentRepo + ContentClient + UserRPCClient + ReasonConfigClient
//
//	↓ 依赖
//
// 基础设施（RPC 客户端、数据库等）
//
// Wire 会自动解决这个依赖链！
func InitializeRecommendationHandler() *handler.RecommendationHandler {
	// 这个函数体会被 Wire 忽略
	// Wire 会生成真实的实现到 wire_gen.go
	wire.Build(
		infrastructureSet,
		repositorySet,
		domainServiceSet,
		applicationServiceSet,
		handlerSet,
	)
	return nil // 占位返回
}

// 实际项目中，可能还需要其他 Injector：

// InitializeRecommendationService 初始化推荐服务（用于测试）
//
// 在测试中，你可能只需要 RecommendationService，不需要 Handler。
// 可以定义一个单独的 Injector：
//
// func InitializeRecommendationService() *service.RecommendationService {
//     wire.Build(
//         infrastructureSet,
//         repositorySet,
//         domainServiceSet,
//         applicationServiceSet,
//     )
//     return nil
// }

// InitializeTestHandler 初始化测试 Handler（使用 mock）
//
// 在测试中，你可能想用 mock 替换某些依赖：
//
// func InitializeTestHandler() *handler.RecommendationHandler {
//     wire.Build(
//         // 使用 mock 的基础设施
//         provideMockUserRPCClient,
//         provideMockContentServiceClient,
//         wire.Value(service.ReasonTextConfigClient(nil)),
//
//         // 使用真实的其他层
//         repositorySet,
//         domainServiceSet,
//         applicationServiceSet,
//         handlerSet,
//     )
//     return nil
// }
