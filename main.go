package main

import (
	"log"
	"net"

	"service/interface/handler"
	"service/rpc_gen/kitex_gen/recommendation/recommendationservice"

	"github.com/cloudwego/kitex/server"
)

// main 服务启动入口
//
// Kitex 微服务的标准启动流程：
// 1. 初始化依赖（数据库、仓储、服务等）
// 2. 创建 Handler（实现 RPC 接口）
// 3. 创建 Kitex Server
// 4. 启动服务监听
//
// 在 DDD 架构中的位置：
// main.go 是基础设施层的一部分，负责：
// - 依赖注入（Dependency Injection）
// - 服务启动和配置
// - 不包含业务逻辑
func main() {
	// 1. 初始化依赖
	// 在实际项目中，这里会：
	// - 加载配置文件
	// - 初始化数据库连接
	// - 初始化 Redis、MQ 等
	// - 创建仓储实现
	// - 创建领域服务
	// - 创建应用服务
	deps := initDependencies()

	// 2. 创建 Handler
	// Handler 实现了 RPC 服务接口
	recommendationHandler := handler.NewRecommendationHandler(
		deps.RecommendationService,
	)

	// 3. 创建 Kitex Server
	// 配置服务选项：
	// - 服务地址和端口
	// - 中间件（日志、监控、限流等）
	// - 服务注册与发现
	// - 链路追踪
	svr := recommendationservice.NewServer(
		recommendationHandler,
		server.WithServiceAddr(&net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 8888,
		}),
		// 在实际项目中，还会添加：
		// server.WithMiddleware(...),      // 中间件
		// server.WithRegistry(...),        // 服务注册
		// server.WithSuite(...),           // 链路追踪
		// server.WithLimit(...),           // 限流配置
	)

	// 4. 启动服务
	log.Println("Recommendation Service starting on :8888")
	err := svr.Run()
	if err != nil {
		log.Fatal("Server run failed:", err)
	}
}

// Dependencies 依赖容器
//
// 这是一个简单的依赖注入容器，在实际项目中可能会使用：
// - Wire (Google 的依赖注入工具)
// - Fx (Uber 的依赖注入框架)
// - 或者自己实现的 IoC 容器
type Dependencies struct {
	// 应用服务
	// RecommendationService *service.RecommendationService

	// 领域服务
	// RecommendationGenerator *domainservice.RecommendationGenerator

	// 仓储
	// SocialGraphRepo repository.SocialGraphRepository
	// ContentRepo     repository.ContentRepository

	// 基础设施
	// DB          *gorm.DB
	// RedisClient *redis.Client
	// UserRPC     UserRPCClient

	// 为了示例简化，这里用 interface{} 代替
	RecommendationService interface{}
}

// initDependencies 初始化依赖
//
// 这是依赖注入的核心函数，负责创建和组装所有依赖。
//
// 实际项目中的完整实现：
//
//	func initDependencies() *Dependencies {
//	    // 1. 加载配置
//	    cfg := config.Load()
//
//	    // 2. 初始化数据库
//	    db := initDB(cfg.Database)
//
//	    // 3. 初始化 Redis
//	    redis := initRedis(cfg.Redis)
//
//	    // 4. 初始化 RPC 客户端
//	    userRPC := initUserRPCClient(cfg.UserService)
//
//	    // 5. 创建仓储实现（基础设施层）
//	    socialGraphRepo := persistence.NewSocialGraphRepository(db)
//	    contentRepo := persistence.NewContentRepository(db)
//
//	    // 6. 创建领域服务（领域层）
//	    generator := domainservice.NewRecommendationGenerator(
//	        socialGraphRepo,
//	        contentRepo,
//	    )
//
//	    // 7. 创建应用服务（应用层）
//	    recommendationService := service.NewRecommendationService(
//	        generator,
//	        socialGraphRepo,
//	        contentRepo,
//	        userRPC,
//	    )
//
//	    return &Dependencies{
//	        RecommendationService: recommendationService,
//	    }
//	}
//
// 依赖注入的好处：
// 1. 控制反转：依赖由外部注入，不在内部创建
// 2. 易于测试：可以注入 mock 对象
// 3. 解耦：各层不直接依赖具体实现
// 4. 灵活配置：可以根据环境注入不同实现
func initDependencies() *Dependencies {
	// 这里是简化版本，实际项目中需要完整实现
	// 参考上面的注释

	log.Println("Initializing dependencies...")

	// TODO: 实际项目中在这里初始化所有依赖
	// - 数据库连接
	// - 仓储实现
	// - 领域服务
	// - 应用服务
	// - RPC 客户端

	return &Dependencies{
		RecommendationService: nil, // 实际项目中返回真实的服务实例
	}
}

// 实际项目中还需要的辅助函数：

// initDB 初始化数据库连接
// func initDB(cfg DatabaseConfig) *gorm.DB {
//     dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
//         cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
//     db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//     if err != nil {
//         log.Fatal("Failed to connect database:", err)
//     }
//     return db
// }

// initRedis 初始化 Redis 连接
// func initRedis(cfg RedisConfig) *redis.Client {
//     return redis.NewClient(&redis.Options{
//         Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
//         Password: cfg.Password,
//         DB:       cfg.DB,
//     })
// }

// initUserRPCClient 初始化 User 服务 RPC 客户端
// func initUserRPCClient(cfg ServiceConfig) UserRPCClient {
//     client, err := userservice.NewClient(
//         cfg.ServiceName,
//         client.WithHostPorts(cfg.Addr),
//     )
//     if err != nil {
//         log.Fatal("Failed to create user rpc client:", err)
//     }
//     return client
// }
