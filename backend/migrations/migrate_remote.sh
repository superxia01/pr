#!/bin/bash
# ============================================
# PR Business 财务系统数据库迁移（杭州服务器）
# 说明：在杭州服务器 Docker 容器内执行迁移 018-022
# 使用：bash migrate_remote.sh
# ============================================

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "============================================"
echo "PR Business 财务系统数据库迁移"
echo "============================================"
echo ""

# 检查是否在 Docker 容器内
if [ ! -f /.dockerenv ]; then
    echo -e "${RED}❌ 错误：此脚本必须在 PostgreSQL Docker 容器内执行${NC}"
    echo ""
    echo "请先执行以下命令进入容器："
    echo "  1. SSH 登录杭州服务器："
    echo "     ssh hangzhou-aliyun"
    echo ""
    echo "  2. 查找 PostgreSQL 容器："
    echo "     docker ps | grep postgres"
    echo ""
    echo "  3. 进入容器并执行迁移："
    echo "     docker exec -it <容器名> bash"
    echo "     cd /tmp"
    echo "     bash migrate_remote.sh"
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

# 检查迁移文件是否存在
if [ ! -f "/tmp/018_022_financial_system.sql" ]; then
    echo -e "${RED}❌ 迁移文件不存在：/tmp/018_022_financial_system.sql${NC}"
    echo ""
    echo "请先上传迁移文件："
    echo "  scp backend/migrations/018_022_financial_system.sql hangzhou-aliyun:/tmp/"
    exit 1
fi

# 执行迁移
echo -e "${YELLOW}开始执行迁移...${NC}"
echo ""

psql -U "$DB_USER" -d "$DB_NAME" -f /tmp/018_022_financial_system.sql

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}============================================"
    echo "✅ 迁移完成！"
    echo "============================================${NC}"
    echo ""
    echo "新创建的表："
    psql -U "$DB_USER" -d "$DB_NAME" -c "
        SELECT
            tablename AS \"表名\",
            schemaname AS \"模式\"
        FROM pg_tables
        WHERE tablename IN ('cash_accounts', 'system_accounts', 'withdrawal_requests_enhanced', 'financial_audit_logs')
        ORDER BY tablename;
    "
    echo ""
    echo "credit_accounts 新增字段："
    psql -U "$DB_USER" -d "$DB_NAME" -c "
        SELECT
            column_name AS \"字段名\",
            data_type AS \"数据类型\",
            is_nullable AS \"可空\"
        FROM information_schema.columns
        WHERE table_name = 'credit_accounts'
          AND column_name = 'frozen_balance';
    "
else
    echo ""
    echo -e "${RED}❌ 迁移失败，请检查错误信息${NC}"
    exit 1
fi
