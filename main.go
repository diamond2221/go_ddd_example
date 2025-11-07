package main

import (
	"log"
	"net"

	"service/application/service"
	domainService "service/domain/service"
	"service/infrastructure/client"
	"service/infrastructure/repository"
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
	RecommendationService *service.RecommendationService

	// 领域服务
	// RecommendationGenerator *domainservice.RecommendationGenerator

	// 仓储
	// SocialGraphRepo repository.SocialGraphRepository
	// ContentRepo     repository.ContentRepository

	// 基础设施
	// DB          *gorm.DB
	// RedisClient *redis.Client
	// UserRPC     service.UserRPCClient
	// ReasonConfigClient service.ReasonTextConfigClient
}

// initDependencies 初始化依赖
//
// 这是依赖注入的核心函数，负责创建和组装所有依赖。
//
// DDD 分层依赖注入顺序（从内到外）：
// 1. 基础设施层：数据库、Redis、RPC 客户端
// 2. 仓储层：实现领域仓储接口
// 3. 领域服务层：实现核心业务逻辑
// 4. 应用服务层：编排用例
// 5. 接口层：处理外部请求
//
// 依赖注入的好处：
// 1. 控制反转：依赖由外部注入，不在内部创建
// 2. 易于测试：可以注入 mock 对象
// 3. 解耦：各层不直接依赖具体实现
// 4. 灵活配置：可以根据环境注入不同实现
func initDependencies() *Dependencies {
	log.Println("Initializing dependencies...")

	// 1. 初始化基础设施（数据库、缓存、RPC 客户端等）
	// 在实际项目中，这里会：
	// - 加载配置文件
	// - 初始化数据库连接
	// - 初始化 Redis 连接
	// db := initDB(cfg.Database)
	// redis := initRedis(cfg.Redis)

	// 2. 初始化 RPC 客户端
	// userRPCClient := initUserRPCClient(cfg.UserService)
	// 示例：使用 mock 实现
	userRPCClient := repository.NewMockUserRPCClient()

	// 3. 初始化推荐理由配置服务客户端（可选）
	// 方式1：使用真实的配置服务
	// reasonConfigClient := client.NewReasonTextConfigHTTPClient("http://config-service:8080")
	//
	// 方式2：不使用配置服务（传 nil，会降级到本地逻辑）
	var reasonConfigClient service.ReasonTextConfigClient = nil
	//
	// 方式3：通过环境变量或配置文件控制
	// if os.Getenv("USE_REASON_CONFIG") == "true" {
	//     reasonConfigClient = client.NewReasonTextConfigHTTPClient(os.Getenv("CONFIG_SERVICE_URL"))
	// }
	//
	// 注意：client.NewReasonTextConfigHTTPClient 已经在 infrastructure/client 包中实现
	// 如果需要使用，取消注释上面的代码即可
	_ = client.NewReasonTextConfigHTTPClient // 避免 unused import 警告

	// 4. 创建仓储实现（基础设施层 → 领域层接口）
	// 在实际项目中，这里会创建真实的数据库仓储
	// socialGraphRepo := persistence.NewMySQLSocialGraphRepository(db)
	// contentRepo := persistence.NewMySQLContentRepository(db)
	//
	// 示例：使用 mock 实现
	socialGraphRepo := repository.NewMockSocialGraphRepository()
	contentRepo := repository.NewMockContentRepository()

	// 5. 创建领域服务（领域层）
	// 领域服务依赖仓储接口，不依赖具体实现
	generator := domainService.NewRecommendationGenerator(
		socialGraphRepo,
		contentRepo,
	)

	// 6. 创建应用服务（应用层）
	// 应用服务依赖领域服务、仓储、RPC 客户端
	recommendationService := service.NewRecommendationService(
		generator,
		socialGraphRepo,
		contentRepo,
		userRPCClient,
		reasonConfigClient, // 可以为 nil
	)

	log.Println("Dependencies initialized successfully")

	return &Dependencies{
		RecommendationService: recommendationService,
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
