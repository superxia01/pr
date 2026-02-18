#!/bin/bash
# ============================================
# PR Business 数据库迁移 - 一键执行版
# 说明：在杭州服务器 Docker 容器内直接执行
# 使用：在容器内执行 bash execute_migrations.sh
# ============================================

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

clear
echo "============================================"
echo "PR Business 数据库迁移"
echo "============================================"
echo ""

# 数据库配置
DB_USER="${DB_USER:-nexus_user}"
DB_NAME="${DB_NAME:-pr_business_db}"

echo "数据库配置："
echo -e "  用户: ${BLUE}${DB_USER}${NC}"
echo -e "  数据库: ${BLUE}${DB_NAME}${NC}"
echo ""
echo "============================================"
echo ""

# 定义迁移函数
execute_migration() {
    local migration_file="$1"
    local migration_name="$2"
    local file_path="/tmp/${migration_file}"

    echo -e "${YELLOW}[1/3]${NC} 执行迁移: ${migration_name}"
    echo -e "  文件: ${migration_file}"

    # 检查文件
    if [ ! -f "$file_path" ]; then
        echo -e "${RED}❌ 文件不存在: ${file_path}${NC}"
        echo ""
        echo "请先上传文件到容器 /tmp/ 目录"
        return 1
    fi

    # 执行迁移
    if psql -U "$DB_USER" -d "$DB_NAME" -f "$file_path" 2>&1; then
        echo -e "${RED}❌ 迁移失败，查看上方错误信息${NC}"
        return 1
    else
        echo -e "${GREEN}✅ 成功${NC}"
    fi

    echo ""
}

# 执行迁移
echo -e "${BLUE}开始执行迁移...${NC}"
echo ""

# 1. 创建 invitation_relationships 表
execute_migration "023_create_invitation_relationships.sql" "创建邀请关系表"

# 2. 修复 creators 表约束
execute_migration "024_fix_creator_inviter_type_constraint.sql" "修复Creator约束拼写"

# 3. 创建 recharge_orders 表
execute_migration "025_create_recharge_orders.sql" "创建充值订单表"

echo -e "${GREEN}============================================"
echo "✅ 所有迁移执行完成！"
echo "============================================${NC}"
echo ""

# 验证结果
echo -e "${BLUE}验证迁移结果...${NC}"
echo ""

echo "检查新创建的表："
psql -U "$DB_USER" -d "$DB_NAME" -c "
    SELECT
        tablename AS \"表名\",
        CASE
            WHEN tablename = 'invitation_relationships' THEN '✅ 邀请关系表'
            WHEN tablename = 'recharge_orders' THEN '✅ 充值订单表'
            ELSE tablename
        END AS \"状态\"
    FROM pg_tables
    WHERE schemaname = 'public'
      AND tablename IN ('invitation_relationships', 'recharge_orders')
    ORDER BY tablename;
"
echo ""

echo "检查约束修复："
CONSTRAINT_CHECK=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c "
    SELECT pg_get_constraintdef(oid)
    FROM pg_constraint
    WHERE conname = 'creators_inviter_type_check'
    LIMIT 1;
")

if echo "$CONSTRAINT_CHECK" | grep -q "SERVICE_PROVIDER_STAFF"; then
    echo -e "${GREEN}✅ creators.inviter_type 约束已正确修复${NC}"
else
    echo -e "${RED}❌ creators.inviter_type 约束未正确修复${NC}"
fi

echo ""
echo -e "${GREEN}============================================"
echo "🎉 迁移流程全部完成！"
echo "============================================${NC}"
echo ""

echo "后续步骤："
echo "  1. 部署新后端代码到服务器"
echo "  2. 重启后端服务"
echo "  3. 测试新功能："
echo "     - 邀请码使用（/api/v1/invitations/use）"
echo "     - 邀请列表查询（/api/v1/invitations/my）"
echo "     - 充值订单流程（/api/v1/recharge-orders）"
echo "     - 任务结算（审核通过后自动执行）"
echo ""
