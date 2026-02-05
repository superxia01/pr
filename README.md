# PR 业务系统

**域名**: pr.crazyaigc.com
**架构**: V3.0 (Go + Vite + React)

---

## 📚 技术规范

本项目遵循 **KeenChase 通用技术规范**：

- **[系统架构与技术标准](../keenchase-standards/architecture.md)** - V3.0 架构、技术栈、代码规范
- **[部署与服务管理](../keenchase-standards/deployment-and-operations.md)** - 部署流程、服务管理
- **[SSH 配置指南](../keenchase-standards/ssh-setup.md)** - SSH 密钥配置、服务器连接
- **[数据库使用指南](../keenchase-standards/database-guide.md)** - 数据库连接、用户权限
- **[安全规范](../keenchase-standards/security.md)** - fail2ban、密钥管理

**认证集成**：
- **[API 接口说明](../keenchase-standards/api.md)** - auth-center 认证接口

---

## 🚀 快速开始

### 1. 本地开发

**前端**：
```bash
cd frontend
npm install
npm run dev
```

**后端**：
```bash
cd backend
go mod download
go run main.go
```

### 2. 部署

详细的部署步骤和配置请查看 **[部署指南](./deploy/README.md)**。

```bash
# 快速部署
cd deploy
./deploy.sh
```

---

## 🏗️ 技术栈

### 前端
- **框架**: React 18 + Vite
- **UI**: shadcn/ui + Tailwind CSS
- **状态管理**: React Context
- **路由**: React Router
- **HTTP**: Axios

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL 15+
- **认证**: JWT (集成 auth-center)

### 部署
- **前端**: Nginx (静态文件服务)
- **后端**: systemd (Go 进程管理)
- **数据库**: PostgreSQL (独立数据库服务器)

---

## 📖 项目文档

- **[PRD](./docs/PRD.md)** - 产品需求文档
- **[架构设计](./docs/ARCHITECTURE.md)** - 系统架构设计
- **[开发计划](./docs/DEVELOPMENT_PLAN.md)** - 开发计划
- **[测试清单](./docs/TESTING_CHECKLIST.md)** - 测试检查清单

---

## 🔗 相关链接

- **auth-center**: [github.com/xxx/auth-center](../auth-center)
- **KeenChase 技术规范**: [github.com/xxx/keenchase-standards](../keenchase-standards)

---

**最后更新**: 2026-02-04
