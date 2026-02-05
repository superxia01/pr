# PR Business 部署配置

**最后更新**: 2026-02-05
**符合标准**: KeenChase V4.0 部署规范

---

## 📚 快速导航

| 文档 | 说明 | 适用场景 |
|------|------|---------|
| **[部署指南](./DEPLOYMENT_GUIDE.md)** | 完整的部署流程和故障排查 | 首次部署、详细了解 |
| **[部署检查清单](./DEPLOYMENT-CHECKLIST.md)** | 部署检查清单 | 每次部署必看 |
| `deploy.sh` | 自动化部署脚本 | 日常部署 |

---

## 🚀 快速开始

### 日常部署（推荐）

```bash
# 1. 确保在项目根目录
cd /path/to/pr-business

# 2. 执行部署脚本
bash deploy/deploy.sh all          # 部署全部
bash deploy/deploy.sh frontend     # 仅前端
bash deploy/deploy.sh backend      # 仅后端
```

### 首次部署

**请参考**：[部署指南 - 首次部署章节](./DEPLOYMENT_GUIDE.md#-首次部署)

---

## 📋 部署前必读

### ⚠️ 核心原则（强制执行）

1. **环境变量与代码分离**
   - ❌ 不要上传 `.env` 文件
   - ✅ 只上传二进制文件和静态资源

2. **目录命名统一**
   - 后端：`/var/www/pr-backend`
   - 前端：`/var/www/pr-frontend`
   - ❌ 不要使用 `pr-business`、`pr-business-frontend`

3. **本地构建，上传产物**
   - ✅ 前端：本地 `npm run build`，上传 `dist/`
   - ✅ 后端：本地交叉编译，上传二进制
   - ❌ 不要在服务器上构建

### 部署前检查清单

使用 **[部署检查清单](./DEPLOYMENT-CHECKLIST.md)** 确保以下项：

- [ ] 本地构建成功
- [ ] 代码已提交
- [ ] 确认部署范围（前端/后端/全部）
- [ ] 服务器连接正常

---

## 🔧 部署脚本说明

### deploy.sh

自动化部署脚本，支持：

```bash
bash deploy/deploy.sh [frontend|backend|all]
```

**功能**：
- ✅ 检查构建产物
- ✅ 上传到正确目录
- ✅ 自动备份旧版本
- ✅ 重启服务
- ✅ 验证部署状态

**配置**：
```bash
SSH_HOST="shanghai-tencent"
BACKEND_DIR="/var/www/pr-backend"       # ✅ 标准目录
FRONTEND_DIR="/var/www/pr-frontend"     # ✅ 标准目录
SERVICE_NAME="pr-business-backend"      # ✅ 正确服务名
```

---

## 📊 服务器信息

### 服务器列表

| 服务器 | IP | SSH 别名 | 用途 |
|--------|-----|---------|------|
| 上海服务器 | 101.35.120.199 | `shanghai-tencent` | 应用服务器 |
| 杭州服务器 | 47.110.82.96 | `hangzhou-ali` | 数据库服务器 |

### 目录结构（标准）

```
/var/www/
├── pr-backend/              # 后端目录
│   ├── pr-business          # 可执行文件
│   ├── .env                 # 环境变量（不提交）
│   └── server.log           # 服务日志
│
└── pr-frontend/             # 前端目录
    ├── index.html
    └── assets/
```

### 服务配置

| 组件 | 服务名/配置 | 状态命令 |
|------|-----------|---------|
| 后端 | `pr-business-backend.service` | `sudo systemctl status pr-business-backend` |
| 前端 | Nginx 静态文件 | `sudo systemctl status nginx` |
| 域名 | pr.crazyaigc.com | - |

---

## 🔍 服务管理

### 常用命令

```bash
# === 后端服务 ===

# 查看状态
ssh shanghai-tencent "sudo systemctl status pr-business-backend"

# 重启服务
ssh shanghai-tencent "sudo systemctl restart pr-business-backend"

# 查看日志
ssh shanghai-tencent "sudo journalctl -u pr-business-backend -f"

# === 前端服务 ===

# 测试配置
ssh shanghai-tencent "sudo nginx -t"

# 重载配置
ssh shanghai-tencent "sudo systemctl reload nginx"

# 查看错误日志
ssh shanghai-tencent "sudo tail -f /var/log/nginx/error.log"
```

### 健康检查

```bash
# 测试网站访问
curl -I https://pr.crazyaigc.com

# 测试后端 API
curl https://pr.crazyaigc.com/api/v1/service-providers

# 检查服务状态（一键）
ssh shanghai-tencent << 'ENDSSH'
sudo systemctl status pr-business-backend --no-pager | head -5
sudo nginx -t
ls -la /var/www/pr-frontend/
ENDSSH
```

---

## 🚨 常见问题

### Q1: 部署后服务无法启动？

**检查**：
```bash
ssh shanghai-tencent "sudo journalctl -u pr-business-backend -n 50"
```

**常见原因**：
- 环境变量错误
- 数据库连接失败
- 端口被占用

**解决**：参考 [部署指南 - 故障排查](./DEPLOYMENT_GUIDE.md#-故障排查)

### Q2: 前端显示旧版本？

**原因**：浏览器缓存

**解决**：
- Mac: `Cmd + Shift + R`
- Windows: `Ctrl + Shift + R`
- 或使用无痕模式

### Q3: API 请求 404？

**检查**：
```bash
# 确认后端服务运行
ssh shanghai-tencent "sudo systemctl status pr-business-backend"

# 确认 Nginx 配置
ssh shanghai-tencent "sudo nginx -t"
```

### Q4: 目录不符合标准？

**检查当前目录**：
```bash
ssh shanghai-tencent "ls -la /var/www/ | grep pr"
```

**正确输出**：
```
pr-backend
pr-frontend
```

**如果出现错误目录**（`pr-business`, `pr-business-frontend`）：
```bash
ssh shanghai-tencent "sudo rm -rf /var/www/pr-business /var/www/pr-business-frontend"
```

---

## 📝 部署记录

每次部署后请记录：

- 部署时间
- 部署内容
- 遇到的问题
- 部署结果

**使用**：[部署检查清单 - 部署记录模板](./DEPLOYMENT-CHECKLIST.md#-部署记录模板)

---

## 🔗 相关文档

- **[KeenChase 部署标准](../../keenchase-standards/deployment-and-operations.md)** - 通用部署规范
- **[部署指南](./DEPLOYMENT_GUIDE.md)** - 详细部署流程
- **[部署检查清单](./DEPLOYMENT-CHECKLIST.md)** - 部署检查清单
- **[SSH 配置指南](../../keenchase-standards/ssh-setup.md)** - SSH 密钥配置
- **[数据库使用指南](../../keenchase-standards/database-guide.md)** - 数据库连接

---

## 📞 联系方式

**技术问题**：
- 查看 [部署指南](./DEPLOYMENT_GUIDE.md)
- 查看 [故障排查](./DEPLOYMENT_GUIDE.md#-故障排查)

**紧急联系**：
- [联系方式待补充]

---

**文档维护**: DevOps Team
**最后更新**: 2026-02-05
