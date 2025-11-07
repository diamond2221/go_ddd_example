package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"service/application/service"
)

// ContentServiceHTTPClient 内容服务HTTP客户端实现
//
// 这是 ContentServiceClient 接口的具体实现，负责通过 HTTP 调用内容服务。
//
// 为什么在基础设施层？
// - HTTP 调用是技术细节
// - 序列化/反序列化是实现细节
// - 错误处理、超时、重试等都是基础设施关注点
//
// 实际场景：
// 推荐服务需要获取用户帖子 →
//
//	应用层调用 ContentServiceClient 接口 →
//	  基础设施层通过 HTTP 调用内容服务 →
//	    内容服务返回帖子数据 →
//	  转换为应用层的 PostInfo →
//	返回给应用层
//
// 对比：
// - ContentRepository：查询本地数据库（SQL）
// - ContentServiceClient：调用远程服务（HTTP/RPC）
type ContentServiceHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewContentServiceHTTPClient 构造函数
func NewContentServiceHTTPClient(baseURL string) *ContentServiceHTTPClient {
	return &ContentServiceHTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Second, // 3秒超时
		},
	}
}

// GetRecentPosts 获取用户最近的帖子
//
// HTTP 调用示例：
// GET /api/v1/users/{userID}/posts?limit=3
//
// 响应示例：
//
//	{
//	  "posts": [
//	    {
//	      "post_id": 123,
//	      "content": "Hello World",
//	      "created_at": "2024-01-01 12:00:00"
//	    }
//	  ]
//	}
//
// 错误处理：
// - 网络错误：返回错误
// - 超时：返回错误
// - 4xx/5xx：返回错误
// - 解析失败：返回错误
func (c *ContentServiceHTTPClient) GetRecentPosts(
	ctx context.Context,
	userID int64,
	limit int,
) ([]*service.PostInfo, error) {
	// 构造请求 URL
	url := fmt.Sprintf("%s/api/v1/users/%d/posts?limit=%d", c.baseURL, userID, limit)

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("http status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response struct {
		Posts []struct {
			PostID    int64  `json:"post_id"`
			Content   string `json:"content"`
			CreatedAt string `json:"created_at"`
		} `json:"posts"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	// 转换为应用层的 PostInfo
	result := make([]*service.PostInfo, 0, len(response.Posts))
	for _, post := range response.Posts {
		result = append(result, &service.PostInfo{
			PostID:    post.PostID,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
		})
	}

	return result, nil
}
