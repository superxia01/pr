#!/bin/bash
# ============================================
# PR Business 财务系统数据库迁移脚本
# 说明：批量执行迁移 018-022
# 使用：bash migrations/run_migrations.sh
# ============================================

set -e  # 遇到错误立即退出

# 数据库连接参数
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-nexus}"
DB_NAME="${DB_NAME:-pr_business_db}"
export PGPASSWORD="${DB_PASSWORD:-nexus123}"

echo "============================================"
echo "PR Business 财务系统数据库迁移"
echo "============================================"
echo ""
echo "数据库配置："
echo "  主机: $DB_HOST"
echo "  端口: $DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"
echo ""
echo "============================================"
echo ""

# 测试数据库连接
echo "测试数据库连接..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo "❌ 数据库连接失败！"
    echo "请检查："
    echo "  1. PostgreSQL 是否运行"
    echo "  2. 数据库配置是否正确 (.env 文件)"
    echo "  3. 用户名密码是否正确"
    exit 1
fi

echo "✅ 数据库连接成功"
echo ""

# ============================================
# 迁移 018: 创建现金账户表
# ============================================
echo "执行迁移 018: 创建现金账户表..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f backend/migrations/018_create_cash_accounts.sql

if [ $? -eq 0 ]; then
    echo "✅ 迁移 018 完成"
else
    echo "❌ 迁移 018 失败"
    exit 1
fi
echo ""

# ============================================
# 迁移 019: 创建系统账户表
# ============================================
echo "执行迁移 019: 创建系统账户表..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f backend/migrations/019_create_system_accounts.sql

if [ $? -eq 0 ]; then
    echo "✅ 迁移 019 完成"
else
    echo "❌ 迁移 019 失败"
    exit 1
fi
echo ""

# ============================================
# 迁移 020: 添加冻结余额字段
# ============================================
echo "执行迁移 020: 添加冻结余额字段..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f backend/migrations/020_extend_credit_accounts_for_cash.sql

if [ $? -eq 0 ]; then
    echo "✅ 迁移 020 完成"
else
    echo "❌ 迁移 020 失败"
    exit 1
fi
echo ""

# ============================================
# 迁移 021: 创建提现申请表
# ============================================
echo "执行迁移 021: 创建提现申请表..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f backend/migrations/021_create_withdrawal_requests.sql

if [ $? -eq 0 ]; then
    echo "✅ 迁移 021 完成"
else
    echo "❌ 迁移 021 失败"
    exit 1
fi
echo ""

# ============================================
# 迁移 022: 创建财务审计日志表
# ============================================
echo "执行迁移 022: 创建财务审计日志表..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f backend/migrations/022_create_financial_audit_logs.sql

if [ $? -eq 0 ]; then
    echo "✅ 迁移 022 完成"
else
    echo "❌ 迁移 022 失败"
    exit 1
fi
echo ""

# ============================================
# 验证迁移结果
# ============================================
echo "============================================"
echo "验证迁移结果..."
echo "============================================"
echo ""

psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "
-- 检查新创建的表
SELECT
    'cash_accounts' AS table_name,
    COUNT(*) AS row_count
FROM information_schema.tables
WHERE table_schema = 'public' AND table_name = 'cash_accounts'
UNION ALL
SELECT
    'system_accounts' AS table_name,
    COUNT(*) AS row_count
FROM information_schema.tables
WHERE table_schema = 'public' AND table_name = 'system_accounts'
UNION ALL
SELECT
    'withdrawal_requests_enhanced' AS table_name,
    COUNT(*) AS row_count
FROM information_schema.tables
WHERE table_schema = 'public' AND table_name = 'withdrawal_requests_enhanced'
UNION ALL
SELECT
    'financial_audit_logs' AS table_name,
    COUNT(*) AS row_count
FROM information_schema.tables
WHERE table_schema = 'public' AND table_name = 'financial_audit_logs'
UNION ALL
SELECT
    column_name AS column_name,
    data_type,
    is_nullable
FROM information_schema.columns
WHERE table_name = 'credit_accounts'
  AND column_name IN ('balance', 'frozen_balance')
ORDER BY ordinal_position;
"

echo ""
echo "============================================"
echo "迁移完成！"
echo "============================================"
