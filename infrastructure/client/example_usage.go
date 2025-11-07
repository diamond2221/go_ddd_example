package client

import (
	"service/application/service"
	domainService "service/domain/service"
	"service/infrastructure/repository"
)

// ExampleWireRecommendationService 示例：如何组装推荐服务
//
// 这个示例展示了如何在实际项目中组装依赖。
//
// 在真实项目中，通常使用依赖注入框架（如 Wire、Dig）来自动化这个过程。
func ExampleWireRecommendationService() *service.RecommendationService {
	// 1. 创建仓储实现（假设已经实现）
	socialGraphRepo := repository.NewMySQLSocialGraphRepository( /* db */ )
	contentRepo := repository.NewMySQLContentRepository( /* db */ )

	// 2. 创建领域服务
	generator := domainService.NewRecommendationGenerator(
		socialGraphRepo,
		contentRepo,
	)

	// 3. 创建 RPC 客户端（假设已经实现）
	userRPCClient := NewUserRPCClient( /* config */ )

	// 4. 创建配置服务客户端（可选）
	// 方式1：使用配置服务
	reasonConfigClient := NewReasonTextConfigHTTPClient("http://config-service:8080")

	// 方式2：不使用配置服务（传 nil，会降级到本地逻辑）
	// var reasonConfigClient service.ReasonTextConfigClient = nil

	// 5. 创建应用服务
	recommendationService := service.NewRecommendationService(
		generator,
		socialGraphRepo,
		contentRepo,
		userRPCClient,
		reasonConfigClient, // 可以传 nil
	)

	return recommendationService
}

// ExampleGradualMigration 示例：渐进式迁移策略
//
// 展示如何从不使用配置服务逐步迁移到使用配置服务。
func ExampleGradualMigration() {
	// 阶段1：不使用配置服务（当前状态）
	// 所有文案使用本地逻辑生成
	_ = service.NewRecommendationService(
		nil, nil, nil, nil,
		nil, // reasonConfigClient = nil
	)

	// 阶段2：灰度发布配置服务
	// 部分用户使用配置服务，部分用户使用本地逻辑
	// 通过特性开关（Feature Flag）控制
	var reasonConfigClient service.ReasonTextConfigClient
	if isFeatureEnabled("use_reason_config_service") {
		reasonConfigClient = NewReasonTextConfigHTTPClient("http://config-service:8080")
	} else {
		reasonConfigClient = nil
	}
	_ = service.NewRecommendationService(
		nil, nil, nil, nil,
		reasonConfigClient,
	)

	// 阶段3：全量使用配置服务
	// 所有用户都使用配置服务，但保留降级逻辑
	reasonConfigClient = NewReasonTextConfigHTTPClient("http://config-service:8080")
	_ = service.NewRecommendationService(
		nil, nil, nil, nil,
		reasonConfigClient,
	)

	// 阶段4（可选）：移除本地逻辑
	// 如果配置服务足够稳定，可以考虑移除 RecommendationReason.Description() 中的降级逻辑
	// 但通常建议保留降级逻辑，以应对配置服务异常
}

// 辅助函数（示例）
func isFeatureEnabled(feature string) bool {
	// 实际项目中，这里会查询特性开关服务
	return false
}

// NewUserRPCClient 示例：创建用户 RPC 客户端（需要实际实现）
func NewUserRPCClient() service.UserRPCClient {
	// TODO: 实现用户 RPC 客户端
	return nil
}
