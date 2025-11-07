package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ReasonTextConfigHTTPClient HTTP 客户端：调用配置服务获取推荐理由文案
//
// 这是基础设施层的实现，负责与外部 HTTP 服务通信。
//
// 为什么在基础设施层？
// - HTTP 调用是技术细节，不是业务逻辑
// - 依赖具体的网络协议和序列化方式
// - 可以随时替换实现（如改用 gRPC）而不影响业务层
//
// 实际业务场景：
// 运营人员在配置后台修改推荐理由文案 →
//
//	配置服务提供 HTTP API →
//	  这个客户端调用 API 获取文案 →
//	    应用服务使用文案展示给用户
//
// 容错设计：
// - 超时控制：避免配置服务慢影响主流程
// - 错误返回：让上层决定如何降级
// - 不缓存：保证文案实时性（可以在上层添加缓存）
type ReasonTextConfigHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewReasonTextConfigHTTPClient 构造函数
func NewReasonTextConfigHTTPClient(baseURL string) *ReasonTextConfigHTTPClient {
	return &ReasonTextConfigHTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 2 * time.Second, // 2秒超时，避免影响主流程
		},
	}
}

// GetReasonText 实现接口：获取推荐理由文案
//
// API 设计示例：
// GET /api/v1/recommendation/reason-text?type=followed_by_following&count=3
//
// 响应示例：
//
//	{
//	  "code": 0,
//	  "message": "success",
//	  "data": {
//	    "text": "你的 3 位好友也关注了TA"
//	  }
//	}
//
// 容错处理：
// - HTTP 请求失败：返回错误，上层降级
// - 响应解析失败：返回错误，上层降级
// - 返回空文案：返回空字符串，上层降级
func (c *ReasonTextConfigHTTPClient) GetReasonText(
	ctx context.Context,
	reasonType string,
	count int,
) (string, error) {
	// 构造请求 URL
	url := fmt.Sprintf(
		"%s/api/v1/recommendation/reason-text?type=%s&count=%d",
		c.baseURL,
		reasonType,
		count,
	)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %w", err)
	}

	// 解析响应
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Text string `json:"text"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("parse response failed: %w", err)
	}

	// 检查业务状态码
	if response.Code != 0 {
		return "", fmt.Errorf("api error: code=%d, message=%s", response.Code, response.Message)
	}

	return response.Data.Text, nil
}
