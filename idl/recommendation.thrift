namespace go recommendation

// 推荐请求
struct GetRecommendationsRequest {
    1: required i64 user_id,  // 用户ID
    2: optional i32 limit = 10,  // 返回数量限制
    3: optional i32 day = 7, // 时间范围 (7 天)
}

// 推荐响应
struct GetRecommendationsResponse {
    1: required list<UserRecommendation> recommendations,
}

// 用户推荐
struct UserRecommendation {
    1: required i64 user_id,
    2: required string username,
    3: required string avatar,
    4: optional string bio,
    5: required string reason,  // 推荐理由
    6: required i32 score,  // 推荐分数
    7: required list<Post> recent_posts,  // 最近的帖子
}

// 帖子
struct Post {
    1: required i64 post_id,
    2: required string content,
    3: required string created_at,
}

// 推荐服务
service RecommendationService {
    // 获取基于关注的推荐
    GetRecommendationsResponse GetFollowingBasedRecommendations(
        1: GetRecommendationsRequest req
    )
}
