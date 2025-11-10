package handler

import (
	"context"
	"errors"

	"service/application/service"

	"service/application/dto"

	"service/rpc_gen/kitex_gen/recommendation"
)

// RecommendationHandler 接口层：RPC 处理器
//
// 什么是接口层？
// 接口层（也叫适配器层、表现层）是 DDD 架构的最外层，
// 负责处理外部世界与应用之间的交互。
//
// 接口层的职责（记住：适配，不是实现）：
// 1. 协议适配：将外部协议（RPC、HTTP、MQ）转换为应用层调用
// 2. 参数验证：验证请求参数的格式和有效性
// 3. 调用应用服务：将请求委托给应用层处理
// 4. 响应转换：将应用层返回的 DTO 转换为协议响应
// 5. 错误处理：捕获异常并转换为合适的错误响应
// 6. 日志记录：记录请求日志、性能指标等
//
// 接口层不应该包含：
// - 业务逻辑：应该在领域层
// - 用例编排：应该在应用层
// - 数据访问：应该通过仓储
//
// 为什么需要接口层？
// 1. 协议无关：业务逻辑不依赖具体的通信协议
//   - 可以同时支持 RPC、HTTP、GraphQL
//   - 切换协议不影响业务代码
//
// 2. 多端适配：不同客户端可以有不同的接口
//   - Web 端可能用 HTTP
//   - 内部服务可能用 RPC
//   - 消息队列可能用 MQ
//
// 3. 版本管理：可以同时维护多个 API 版本
//   - v1、v2 接口共存
//   - 底层业务逻辑不变
//
// 实际业务场景：
// 客户端发起 RPC 请求 →
//
//	Handler 接收并验证参数 →
//	调用应用服务 →
//	转换响应 →
//	返回给客户端
//
// 对比传统方式：
// 传统方式：Controller 直接调用 Service，协议和业务耦合
// DDD 方式：Handler 只负责协议适配，业务逻辑在内层
type RecommendationHandler struct {
	recommendationService *service.RecommendationService
}

// NewRecommendationHandler 构造函数
func NewRecommendationHandler(
	recommendationService *service.RecommendationService,
) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationService: recommendationService,
	}
}

// GetFollowingBasedRecommendations RPC 方法实现
func (h *RecommendationHandler) GetFollowingBasedRecommendations(
	ctx context.Context,
	req *recommendation.GetRecommendationsRequest,
) (*recommendation.GetRecommendationsResponse, error) {

	// 参数验证
	if req.UserId <= 0 {
		return nil, ErrInvalidUserID
	}
	if req.Limit <= 0 {
		req.Limit = 10 // 默认值
	}

	// 调用应用服务
	result, err := h.recommendationService.GetFollowingBasedRecommendations(
		ctx,
		req.UserId,
		int(req.Limit),
	)
	if err != nil {
		return nil, err
	}

	// 转换为 RPC 响应
	res := h.convertToRPCResponse(result)
	return res, nil
}

// convertToRPCResponse 辅助方法：DTO -> RPC 响应转换
func (h *RecommendationHandler) convertToRPCResponse(
	dto *dto.RecommendationResponse,
) *recommendation.GetRecommendationsResponse {
	resp := &recommendation.GetRecommendationsResponse{
		Recommendations: make([]*recommendation.UserRecommendation, 0, len(dto.Recommendations)),
	}

	for _, rec := range dto.Recommendations {
		rpcRec := &recommendation.UserRecommendation{
			UserId:      rec.UserID,
			Username:    rec.Username,
			Avatar:      rec.Avatar,
			Bio:         rec.Bio,
			Reason:      rec.Reason,
			Score:       int32(rec.Score),
			RecentPosts: h.convertPostsToRPC(rec.RecentPosts),
		}
		resp.Recommendations = append(resp.Recommendations, rpcRec)
	}

	return resp
}

// convertPostsToRPC 辅助方法：PostDTO -> RPC Post 转换
func (h *RecommendationHandler) convertPostsToRPC(
	posts []*dto.PostDTO,
) []*recommendation.Post {
	result := make([]*recommendation.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, &recommendation.Post{
			PostId:    post.PostID,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
		})
	}
	return result
}

var (
	ErrInvalidUserID = errors.New("invalid user id")
)
