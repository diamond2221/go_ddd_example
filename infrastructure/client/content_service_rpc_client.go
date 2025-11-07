package client

import (
	"context"
	"fmt"

	"service/application/service"
	// 假设你有 content 服务的 Kitex 生成代码
	// "service/rpc_gen/kitex_gen/content"
	// "service/rpc_gen/kitex_gen/content/contentservice"
)

// ContentServiceRPCClient 内容服务RPC客户端实现（使用 Kitex）
//
// 这是使用 Kitex RPC 框架的实现版本。
//
// 对比 HTTP 版本：
// - HTTP：通用、跨语言、易调试
// - RPC：高性能、类型安全、代码生成
//
// 使用场景：
// - 内部微服务：推荐使用 RPC（性能更好）
// - 跨团队/跨语言：推荐使用 HTTP（兼容性更好）
//
// 实际使用：
// 1. 定义 content.thrift（IDL）
// 2. 使用 Kitex 生成客户端代码
// 3. 实现这个适配器（将 RPC 响应转换为应用层的 PostInfo）
type ContentServiceRPCClient struct {
	// client contentservice.Client // Kitex 生成的客户端
}

// NewContentServiceRPCClient 构造函数
//
// 实际使用示例：
//
//	client, err := contentservice.NewClient(
//	    "content-service",
//	    client.WithHostPorts("127.0.0.1:8889"),
//	)
//	if err != nil {
//	    panic(err)
//	}
//	return &ContentServiceRPCClient{client: client}
func NewContentServiceRPCClient( /* client contentservice.Client */ ) *ContentServiceRPCClient {
	return &ContentServiceRPCClient{
		// client: client,
	}
}

// GetRecentPosts 获取用户最近的帖子（RPC 版本）
//
// RPC 调用示例：
//
//	req := &content.GetRecentPostsRequest{
//	    UserId: userID,
//	    Limit:  int32(limit),
//	}
//	resp, err := c.client.GetRecentPosts(ctx, req)
//
// 优势：
// - 类型安全：编译时检查
// - 高性能：二进制序列化
// - 代码生成：自动生成客户端代码
func (c *ContentServiceRPCClient) GetRecentPosts(
	ctx context.Context,
	userID int64,
	limit int,
) ([]*service.PostInfo, error) {
	// 实际实现示例（需要 Kitex 生成代码）：
	//
	// req := &content.GetRecentPostsRequest{
	//     UserId: userID,
	//     Limit:  int32(limit),
	// }
	//
	// resp, err := c.client.GetRecentPosts(ctx, req)
	// if err != nil {
	//     return nil, fmt.Errorf("rpc call failed: %w", err)
	// }
	//
	// // 转换 RPC 响应 → 应用层 PostInfo
	// result := make([]*service.PostInfo, 0, len(resp.Posts))
	// for _, post := range resp.Posts {
	//     result = append(result, &service.PostInfo{
	//         PostID:    post.PostId,
	//         Content:   post.Content,
	//         CreatedAt: post.CreatedAt,
	//     })
	// }
	//
	// return result, nil

	// 占位实现
	return nil, fmt.Errorf("not implemented: need Kitex generated code")
}
