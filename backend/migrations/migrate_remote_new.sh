#!/bin/bash
# ============================================
# PR Business 财务系统数据库迁移（杭州服务器）
# 说明：在杭州服务器 Docker 容器内执行迁移
# 使用：
#   1. SSH 到服务器：ssh hangzhou-aliyun
#   2. 找到 PostgreSQL 容器：docker ps | grep postgres
#   3. 进入容器：docker exec -it <容器名> bash
#   4. 在容器内执行：cd /tmp && bash migrate_remote_new.sh
# ============================================

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "============================================"
echo "PR Business 数据库迁移"
echo "============================================"
echo ""

# 检查是否在 Docker 容器内
if [ ! -f /.dockerenv ]; then
    echo -e "${RED}❌ 错误：此脚本必须在 PostgreSQL Docker 容器内执行${NC}"
    echo ""
    echo "请先执行以下命令进入容器："
    echo " 1. SSH 登录杭州服务器："
    echo "     ssh hangzhou-aliyun"
    echo ""
    echo " 2. 查找 PostgreSQL 容器："
    echo "     docker ps | grep postgres"
    echo ""
    echo " 3. 进入容器并执行迁移："
    echo "     docker exec -it <容器名> bash"
    echo ""
    echo "     cd /tmp"
    echo "     bash migrate_remote_new.sh"
    exit 1
fi

# 数据库连接参数（容器内环境）
DB_USER="${DB_USER:-nexus_user}"
DB_NAME="${DB_NAME:-pr_business_db}"

echo "数据库配置："
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"
echo ""
echo "============================================"
echo ""

# 迁移文件列表（按顺序执行）
MIGRATION_FILES=(
    "023_create_invitation_relationships.sql"
    "024_fix_creator_inviter_type_constraint.sql"
    "025_create_recharge_orders.sql"
)

# 检查并执行每个迁移
for migration_file in "${MIGRATION_FILES[@]}"; do
    local_path="/tmp/${migration_file}"

    echo -e "${YELLOW}处理迁移文件: ${migration_file}${NC}"

    # 检查文件是否存在
    if [ ! -f "$local_path" ]; then
        echo -e "${RED}❌ 迁移文件不存在：${local_path}${NC}"
        echo ""
        echo "请先上传迁移文件到服务器 /tmp/ 目录："
        echo "  scp backend/migrations/${migration_file} hangzhou-aliyun:/tmp/"
        exit 1
    fi

    # 执行迁移
    echo -e "执行迁移: ${migration_file}..."
    psql -U "$DB_USER" -d "$DB_NAME" -f "$local_path" > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ ${migration_file} 完成${NC}"
    else
        echo -e "${RED}❌ ${migration_file} 失败${NC}"
        exit 1
    fi

    echo ""
done

echo ""
echo -e "${GREEN}============================================"
echo "✅ 所有迁移完成！"
echo "============================================${NC}"
echo ""
echo "新创建的表："
psql -U "$DB_USER" -d "$DB_NAME" -c "
    SELECT
        tablename AS \"表名\",
        schemaname AS \"模式\"
    FROM pg_tables
    WHERE tablename IN ('invitation_relationships', 'recharge_orders')
    ORDER BY tablename;
"
echo ""
echo "修复的约束："
echo "  - creators.inviter_type 约束（PROVIDER_STAFF -> SERVICE_PROVIDER_STAFF）"
echo ""
