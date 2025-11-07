#!/bin/bash

# 构建脚本
#
# 用于编译 Kitex 微服务

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Building Recommendation Service ===${NC}"

# 项目信息
SERVICE_NAME="recommendation-service"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(go version | awk '{print $3}')

echo -e "${YELLOW}Service:    $SERVICE_NAME${NC}"
echo -e "${YELLOW}Version:    $VERSION${NC}"
echo -e "${YELLOW}Build Time: $BUILD_TIME${NC}"
echo -e "${YELLOW}Go Version: $GO_VERSION${NC}"

# 清理旧的构建产物
echo -e "${YELLOW}Cleaning old builds...${NC}"
rm -f "$SERVICE_NAME"
rm -rf output/

# 下载依赖
echo -e "${YELLOW}Downloading dependencies...${NC}"
go mod download

# 运行测试
echo -e "${YELLOW}Running tests...${NC}"
go test ./... -v

# 构建
echo -e "${YELLOW}Building...${NC}"
go build \
    -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" \
    -o "$SERVICE_NAME" \
    .

# 检查构建结果
if [ -f "$SERVICE_NAME" ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
    echo -e "${GREEN}Binary: ./$SERVICE_NAME${NC}"

    # 显示二进制文件信息
    ls -lh "$SERVICE_NAME"

    echo ""
    echo -e "${GREEN}Run the service:${NC}"
    echo "  ./$SERVICE_NAME"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}=== Build Complete ===${NC}"
