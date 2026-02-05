#!/bin/bash
# PR Business 生产环境部署脚本
# 遵循 KeenChase V3.0 技术规范
#
# 服务器:
#   - 上海服务器 (101.35.120.199): 应用服务器，运行 Go 后端和 Nginx
#   - 杭州服务器 (47.110.82.96): 数据库服务器
#
# 部署目录:
#   - 后端: /var/www/pr-backend
#   - 前端: /var/www/pr-frontend
#
# 使用方法:
#   ./deploy-production.sh [frontend|backend|all]
#

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 服务器配置
SSH_HOST="shanghai-tencent"
BACKEND_DIR="/var/www/pr-backend"
FRONTEND_DIR="/var/www/pr-frontend"
SERVICE_NAME="pr-business-backend"

# 部署目标（默认全部）
DEPLOY_TARGET="${1:-all}"

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# ============================================
# 前置检查
# ============================================
check_prerequisites() {
    log_info "🔍 检查部署前置条件..."

    # 检查 SSH 连接
    if ! ssh -o ConnectTimeout=5 ${SSH_HOST} "echo 'SSH连接成功'" > /dev/null 2>&1; then
        log_error "无法连接到服务器 ${SSH_HOST}"
        log_error "请检查 SSH 配置: ~/.ssh/config"
        exit 1
    fi

    # 检查本地构建产物
    if [ "$DEPLOY_TARGET" = "frontend" ] || [ "$DEPLOY_TARGET" = "all" ]; then
        if [ ! -d "frontend/dist" ]; then
            log_error "前端构建产物不存在，请先运行: cd frontend && npm run build"
            exit 1
        fi
    fi

    if [ "$DEPLOY_TARGET" = "backend" ] || [ "$DEPLOY_TARGET" = "all" ]; then
        if [ ! -f "backend/pr-business-linux" ]; then
            log_error "后端二进制文件不存在，请先运行: cd backend && go build -o pr-business-linux"
            exit 1
        fi
    fi

    log_info "✅ 前置检查通过"
}

# ============================================
# 部署前端
# ============================================
deploy_frontend() {
    log_info "📦 开始部署前端..."

    cd "$(dirname "$0")/.."

    # 上传静态文件（删除地图文件和 .DS_Store）
    rsync -avz --delete \
        --exclude '*.map' \
        --exclude '.DS_Store' \
        dist/ \
        ${SSH_HOST}:${FRONTEND_DIR}/

    # 重载 Nginx
    ssh ${SSH_HOST} "sudo systemctl reload nginx"

    log_info "✅ 前端部署完成"
}

# ============================================
# 部署后端
# ============================================
deploy_backend() {
    log_info "📦 开始部署后端..."

    cd "$(dirname "$0")/.."

    # 上传二进制文件
    # ⚠️ 注意：不上传 .env 文件（环境变量与代码分离）
    scp pr-business-linux ${SSH_HOST}:${BACKEND_DIR}/

    # 在服务器上执行部署操作
    ssh ${SSH_HOST} <<ENDSSH
set -e

cd ${BACKEND_DIR}

# 备份旧版本
if [ -f pr-business ]; then
    BACKUP_FILE="pr-business.backup.\$(date +%Y%m%d_%H%M%S)"
    mv pr-business \${BACKUP_FILE}
    echo "✅ 已备份旧版本: \${BACKUP_FILE}"
fi

# 重命名新版本
mv pr-business-linux pr-business

# 设置可执行权限
chmod +x pr-business

# 重启服务
sudo systemctl restart ${SERVICE_NAME}

# 等待服务启动
sleep 3

# 检查服务状态
if sudo systemctl is-active --quiet ${SERVICE_NAME}; then
    echo "✅ 服务启动成功"
else
    echo "❌ 服务启动失败"
    sudo systemctl status ${SERVICE_NAME} --no-pager
    exit 1
fi
ENDSSH

    log_info "✅ 后端部署完成"
}

# ============================================
# 验证部署
# ============================================
verify_deployment() {
    log_info "🔍 验证部署..."

    # 测试后端健康检查
    sleep 2

    if curl -sf https://pr.crazyaigc.com/health > /dev/null; then
        log_info "✅ 后端健康检查通过"
    else
        log_warn "⚠️ 后端健康检查失败，请手动检查"
    fi

    # 显示服务状态
    ssh ${SSH_HOST} "sudo systemctl status ${SERVICE_NAME} --no-pager"
}

# ============================================
# 主流程
# ============================================
main() {
    echo ""
    echo "=========================================="
    echo "  PR Business 生产环境部署"
    echo "  目标: ${DEPLOY_TARGET}"
    echo "=========================================="
    echo ""

    check_prerequisites

    if [ "$DEPLOY_TARGET" = "frontend" ] || [ "$DEPLOY_TARGET" = "all" ]; then
        deploy_frontend
        echo ""
    fi

    if [ "$DEPLOY_TARGET" = "backend" ] || [ "$DEPLOY_TARGET" = "all" ]; then
        deploy_backend
        echo ""
    fi

    if [ "$DEPLOY_TARGET" = "all" ]; then
        verify_deployment
    fi

    echo ""
    log_info "🎉 部署完成！"
    echo ""
    echo "📍 访问地址："
    echo "  前端: https://pr.crazyaigc.com"
    echo "  后端: https://pr.crazyaigc.com/api/v1"
    echo ""
    echo "📊 管理命令："
    echo "  查看日志: ssh ${SSH_HOST} 'sudo journalctl -u ${SERVICE_NAME} -f'"
    echo "  查看状态: ssh ${SSH_HOST} 'sudo systemctl status ${SERVICE_NAME}'"
    echo "  重启服务: ssh ${SSH_HOST} 'sudo systemctl restart ${SERVICE_NAME}'"
    echo ""
}

# 执行主流程
main
