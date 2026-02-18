#!/bin/bash
# ============================================
# PR Business 迁移部署和执行脚本
# 说明：上传迁移文件到杭州服务器并执行
# ============================================

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 服务器配置
SERVER="hangzhou-aliyun"
SERVER_USER="xia"  # 根据实际情况修改
REMOTE_TMP="/tmp"
CONTAINER_NAME="postgres"  # 根据实际容器名修改

# 新迁移文件列表
MIGRATION_FILES=(
    "backend/migrations/023_create_invitation_relationships.sql"
    "backend/migrations/024_fix_creator_inviter_type_constraint.sql"
    "backend/migrations/025_create_recharge_orders.sql"
)

echo "============================================"
echo "PR Business 迁移部署"
echo "============================================"
echo ""

# 1. 上传迁移文件到服务器
echo -e "${YELLOW}步骤 1/3: 上传迁移文件到服务器${NC}"
echo ""

for migration_file in "${MIGRATION_FILES[@]}"; do
    echo "上传: ${migration_file}"
    scp "${migration_file}" "${SERVER_USER}@${SERVER}:${REMOTE_TMP}/"

    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ 上传失败: ${migration_file}${NC}"
        exit 1
    fi
done

echo -e "${GREEN}✅ 所有文件上传完成${NC}"
echo ""

# 2. 在服务器上执行迁移
echo -e "${YELLOW}步骤 2/3: 在服务器上执行迁移${NC}"
echo ""

# 生成远程执行命令
REMOTE_CMD="cd ${REMOTE_TMP} && bash migrate_remote_new.sh"

echo "执行命令: ssh ${SERVER_USER}@${SERVER} \"${REMOTE_CMD}\""
echo ""

ssh "${SERVER_USER}@${SERVER}" "${REMOTE_CMD}"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}============================================"
    echo "✅ 迁移部署完成！"
    echo "============================================${NC}"
else
    echo ""
    echo -e "${RED}❌ 迁移执行失败，请检查错误信息${NC}"
    exit 1
fi

echo ""
echo "后续步骤："
echo "  1. SSH 登录服务器验证表已创建："
echo "     ssh ${SERVER_USER}@${SERVER}"
echo "  2. 进入 PostgreSQL 容器："
echo "     docker exec -it ${CONTAINER_NAME} bash"
echo "  3. 连接数据库验证："
echo "     psql -U nexus_user -d pr_business_db"
echo "  4. 查询新表："
echo "     \\dt invitation_relationships recharge_orders"
echo ""
