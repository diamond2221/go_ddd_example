#!/bin/bash

# Kitex 代码生成脚本
#
# 这个脚本用于根据 Thrift IDL 生成 Kitex 代码
# 在实际项目中，每次修改 IDL 后都需要运行这个脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Kitex Code Generation ===${NC}"

# 检查 kitex 工具是否安装
if ! command -v kitex &> /dev/null; then
    echo -e "${RED}Error: kitex command not found${NC}"
    echo -e "${YELLOW}Please install kitex first:${NC}"
    echo "  go install github.com/cloudwego/kitex/tool/cmd/kitex@latest"
    exit 1
fi

# 检查 thriftgo 是否安装
if ! command -v thriftgo &> /dev/null; then
    echo -e "${RED}Error: thriftgo command not found${NC}"
    echo -e "${YELLOW}Please install thriftgo first:${NC}"
    echo "  go install github.com/cloudwego/thriftgo@latest"
    exit 1
fi

# 项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$PROJECT_ROOT"

echo -e "${GREEN}Project root: $PROJECT_ROOT${NC}"

# IDL 文件路径
IDL_FILE="idl/recommendation.thrift"

if [ ! -f "$IDL_FILE" ]; then
    echo -e "${RED}Error: IDL file not found: $IDL_FILE${NC}"
    exit 1
fi

echo -e "${GREEN}IDL file: $IDL_FILE${NC}"

# 生成代码
echo -e "${YELLOW}Generating Kitex code...${NC}"

# Kitex 生成命令说明：
# -module: Go module 名称（与 go.mod 中的 module 一致）
# -service: 服务名称
# -use: 指定生成代码的输出目录
# IDL_FILE: Thrift IDL 文件路径
#
# 生成的文件：
# - rpc_gen/kitex_gen/recommendation/*.go - Thrift 结构体
# - rpc_gen/kitex_gen/recommendation/recommendationservice/*.go - 服务接口和实现

kitex \
    -module service \
    -service recommendation \
    -use rpc_gen/kitex_gen/recommendation \
    "$IDL_FILE"

echo -e "${GREEN}✓ Code generation completed${NC}"

# 生成的文件列表
echo -e "${YELLOW}Generated files:${NC}"
find rpc_gen/kitex_gen -type f -name "*.go" | while read -r file; do
    echo "  - $file"
done

echo ""
echo -e "${GREEN}=== Generation Summary ===${NC}"
echo -e "IDL file:     ${YELLOW}$IDL_FILE${NC}"
echo -e "Output dir:   ${YELLOW}rpc_gen/kitex_gen/${NC}"
echo -e "Module name:  ${YELLOW}service${NC}"
echo -e "Service name: ${YELLOW}recommendation${NC}"

echo ""
echo -e "${GREEN}Next steps:${NC}"
echo "1. Review generated code in rpc_gen/kitex_gen/"
echo "2. Implement business logic in interface/handler/"
echo "3. Run: go mod tidy"
echo "4. Run: go build"
echo "5. Run: ./service"

echo ""
echo -e "${GREEN}=== Done ===${NC}"
