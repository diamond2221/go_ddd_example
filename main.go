package main

import (
	"log"
	"net"

	"service/rpc_gen/kitex_gen/recommendation/recommendationservice"

	"github.com/cloudwego/kitex/server"
)

// main 服务启动入口（使用 Wire 依赖注入）
//
// Kitex 微服务的标准启动流程：
// 1. 初始化依赖（使用 Wire 自动生成）
// 2. 创建 Kitex Server
// 3. 启动服务监听
//
// 依赖注入方式：
// - 旧方式：手动在 initDependencies() 中创建所有对象（已移除）
// - 新方式：使用 Wire 自动生成依赖注入代码
//
// Wire 使用步骤：
// 1. 定义 wire.go（Provider 和 Injector）
// 2. 运行 wire 命令生成 wire_gen.go
// 3. 使用生成的 InitializeRecommendationHandler() 函数
//
// 命令：
//
//	go install github.com/google/wire/cmd/wire@latest
//	wire  # 生成 wire_gen.go
//
// 对比：
// ┌─────────────────────────────────────────────────────┐
// │ 手动方式（旧）                                       │
// │ - initDependencies() 手动创建所有对象（100+ 行）     │
// │ - 依赖顺序容易出错                                   │
// │ - 运行时才发现依赖错误                               │
// └─────────────────────────────────────────────────────┘
//
// ┌─────────────────────────────────────────────────────┐
// │ Wire 方式（新）                                      │
// │ - InitializeRecommendationHandler() 自动生成         │
// │ - Wire 自动解决依赖顺序                              │
// │ - 编译时检查依赖错误                                 │
// └─────────────────────────────────────────────────────┘
func main() {
	// 1. 使用 Wire 生成的函数初始化依赖
	// 这一行代码替代了之前的整个 initDependencies() 函数！
	// Wire 会自动：
	// - 创建所有依赖对象
	// - 按正确顺序注入依赖
	// - 返回最终的 Handler
	recommendationHandler := InitializeRecommendationHandler()

	// 2. 创建 Kitex Server
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

	// 3. 启动服务
	log.Println("Recommendation Service starting on :8888 (using Wire)")
	err := svr.Run()
	if err != nil {
		log.Fatal("Server run failed:", err)
	}
}

// Wire 依赖注入说明
//
// 之前的手动依赖注入代码（initDependencies 函数）已经移除。
// 现在使用 Wire 自动生成依赖注入代码。
//
// Wire 配置文件：
// - wire.go：定义 Provider（如何构造对象）和 Injector（需要什么对象）
// - wire_gen.go：Wire 自动生成的依赖注入代码（不要手动编辑）
//
// 使用步骤：
// 1. 安装 Wire：go install github.com/google/wire/cmd/wire@latest
// 2. 运行 Wire：wire（在项目根目录）
// 3. Wire 会生成 wire_gen.go 文件
// 4. 使用生成的 InitializeRecommendationHandler() 函数
//
// 依赖注入流程（由 Wire 自动完成）：
// 1. 基础设施层：创建 RPC 客户端、数据库连接等
// 2. 仓储层：创建仓储实现
// 3. 领域服务层：创建领域服务（依赖仓储）
// 4. 应用服务层：创建应用服务（依赖领域服务、仓储、RPC 客户端）
// 5. 接口层：创建 Handler（依赖应用服务）
//
// Wire 的优势：
// 1. 编译时检查：依赖错误在编译时发现，不是运行时
// 2. 自动解决依赖顺序：不需要手动管理依赖顺序
// 3. 代码简洁：不需要写冗长的初始化代码
// 4. 易于维护：添加新依赖只需添加 Provider
//
// 详细文档：
// - docs/WIRE_GUIDE.md：Wire 完整使用指南
// - docs/WIRE_COMPARISON.md：手动 vs Wire 的详细对比
