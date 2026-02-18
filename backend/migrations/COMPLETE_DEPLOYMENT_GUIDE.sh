#!/bin/bash
# ============================================
# PR Business 完整部署流程
# 步骤：1. 数据库迁移 → 2. 部署后端
# ============================================

set -e

GREEN='\033[0;32m'
BLUE='\033[1;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}============================================"
echo "PR Business 完整部署流程"
echo "============================================${NC}"
echo ""

# ============================================
# 步骤 1: 数据库迁移
# ============================================
echo -e "${BLUE}步骤 1/3: 数据库迁移${NC}"
echo ""
echo "执行以下命令："
echo ""
echo -e "${YELLOW}# 1. 进入 PostgreSQL 容器${NC}"
echo "docker exec -it postgres bash"
echo ""
echo -e "${YELLOW}# 2. 在容器内执行迁移脚本${NC}"
echo "cd /tmp"
echo "bash execute_migrations.sh"
echo ""
echo "或者直接执行 SQL："
echo ""
echo "psql -U nexus_user -d pr_business_db << 'SQLEOF'"
echo ""
echo "  -- 创建邀请关系表"
echo "  CREATE TABLE IF NOT EXISTS invitation_relationships (...);"
echo ""
echo "  -- 修复约束"
echo "  ALTER TABLE creators DROP CONSTRAINT IF EXISTS creators_inviter_type_check;"
echo "  ALTER TABLE creators ADD CONSTRAINT creators_inviter_type_check"
echo "    CHECK (inviter_type IS NULL OR inviter_type IN ('SERVICE_PROVIDER_STAFF', 'SERVICE_PROVIDER_ADMIN', 'OTHER'));"
echo ""
echo "  -- 创建充值订单表"
echo "  CREATE TABLE IF NOT EXISTS recharge_orders (...);"
echo ""
echo "SQLEOF"
echo ""
echo -e "${GREEN}✅ 迁移完成标志：invitation_relationships 和 recharge_orders 表存在${NC}"
echo ""
read -p "按 Enter 继续部署步骤..."
echo ""

# ============================================
# 步骤 2: 上传迁移文件（如果尚未上传）
# ============================================
echo -e "${BLUE}步骤 2/3: 上传迁移文件到服务器${NC}"
echo ""
echo "在本地 Mac 上执行："
echo ""
echo -e "${YELLOW}cd /Users/xia/Documents/GitHub/pr/backend/migrations${NC}"
echo ""
echo -e "${YELLOW}scp execute_migrations.sh hangzhou-aliyun:/tmp/${NC}"
echo ""
echo -e "${YELLOW}scp 023_create_invitation_relationships.sql hangzhou-aliyun:/tmp/${NC}"
echo -e "${YELLOW}scp 024_fix_creator_inviter_type_constraint.sql hangzhou-aliyun:/tmp/${NC}"
echo -e "${YELLOW}scp 025_create_recharge_orders.sql hangzhou-aliyun:/tmp/${NC}"
echo ""
read -p "文件上传完成后按 Enter 继续..."
echo ""

# ============================================
# 步骤 3: 准备部署
# ============================================
echo -e "${BLUE}步骤 3/3: 准备 skill 部署${NC}"
echo ""
echo "确保项目配置正确："
echo "  - 项目名: pr"
echo "  - 项目目录: /Users/xia/Documents/GitHub/pr"
echo ""
echo "然后执行："
echo -e "${YELLOW}/keenchase-deploy${NC}"
echo ""
echo "skill 会自动完成："
echo "  1. ✅ 前置检查"
echo "  2. ✅ 后端编译（CGO_ENABLED=0 GOOS=linux go build）"
echo "  3. ✅ 上传到服务器"
echo "  4. ✅ 重启服务"
echo "  5. ✅ 健康检查"
echo ""
echo -e "${GREEN}============================================"
echo "部署完成后验证："
echo "============================================${NC}"
echo ""
echo "1. 数据库表已创建："
echo "   - invitation_relationships"
echo "   - recharge_orders"
echo ""
echo "2. 新增 API 端点："
echo "   - POST /api/v1/invitations/my (邀请列表）"
echo "   - POST /api/v1/recharge-orders (提交充值订单）"
echo "   - POST /api/v1/recharge-orders/:id/audit (审核充值订单）"
echo ""
echo "3. 验证修复："
echo "   - 任务审核后自动执行结算"
echo "   - Creator 约束已修复"
echo ""
