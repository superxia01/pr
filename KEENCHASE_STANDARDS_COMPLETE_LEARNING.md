# KeenChase 技术规范完整学习报告

**学习日期**: 2026-02-04
**学习范围**: 全部9份文档（包括2份新增文件）
**标准版本**: V3.7
**学习者**: AI Assistant

---

## 📚 文档概览

### 已学习的9份文档

| # | 文档名称 | 大小 | 更新时间 | 新增 | 学习状态 |
|---|---------|------|----------|------|----------|
| 1 | README.md | 5,675 bytes | 10:04 | | ✅ 完成 |
| 2 | architecture.md | 24,671 bytes | 10:03 | | ✅ 完成 |
| 3 | ssh-setup.md | 11,411 bytes | 10:03 | | ✅ 完成 |
| 4 | database-guide.md | 5,747 bytes | 10:51 | | ✅ 完成 |
| 5 | deployment-and-operations.md | 27,718 bytes | 11:18 | | ✅ 完成 |
| 6 | api.md | 25,424 bytes | 10:03 | | ✅ 完成 |
| 7 | security.md | 45,538 bytes | 10:03 | | ✅ 完成 |
| 8 | deploy-template.sh | 3,548 bytes | 11:19 | ⭐ | ✅ 完成 |
| 9 | env.example.md | 2,075 bytes | 11:18 | ⭐ | ✅ 完成 |

**新增文件重点**: deploy-template.sh 和 env.example.md 提供了标准化的部署脚本和环境变量模板

---

## 🎯 核心概念总结

### 1. 两个重要的"用户"概念（贯穿所有文档）

**操作系统用户** - 用于SSH登录服务器
- 上海服务器: `ubuntu` (普通用户)
- 杭州服务器: `root` (管理员用户)
- 香港服务器: `root` (管理员用户)

**数据库用户** - 用于PostgreSQL连接
- 统一用户: `nexus_user` (超级用户，密码: `hRJ9NSJApfeyFDraaDgkYowY`)
- 2026-02-03密钥泄露事件后，已废弃所有专用数据库用户
- 所有业务系统统一使用 `nexus_user`

### 2. 服务器架构

```
📍 杭州服务器 (47.110.82.96) - 统一数据库服务器
└─ PostgreSQL 15 (端口5432)
   ├─ auth_center_db (账号中心)
   ├─ pr_business_db (PR业务)
   ├─ pixel_business_db (AI生图)
   └─ quote_business_db (报价系统)

📍 上海服务器 (101.35.120.199) - 应用服务器
├─ 前端: Nginx 服务静态文件
└─ 后端: Go API (Systemd 管理)
   ├─ os.crazyaigc.com (账号中心, :8080)
   ├─ pr.crazyaigc.com (PR业务, :8081)
   ├─ pixel.crazyaigc.com (AI生图, :8082)
   └─ quote.crazyaigc.com (报价系统, :8083)
```

### 3. 数据库连接方式（通过SSH隧道）

**连接字符串格式**:
```bash
postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@localhost:5432/{数据库名}?sslmode=disable
```

**强制规则**:
- ✅ 必须使用 `localhost` (通过SSH隧道转发)
- ✅ 必须使用 `nexus_user`
- ✅ 必须使用 `sslmode=disable` (SSH隧道已加密)
- ❌ 禁止直连 47.110.82.96:5432 且不用SSL

---

## 📖 各文档详细学习笔记

### 1. README.md - 文档导航中心

**核心要点**:
- 所有新系统必须采用 V3.0 架构 (Go + Vite + React)
- 快速开始顺序: architecture → ssh-setup → database → deployment
- 已部署4个系统: auth-center, pr-business, pixel-business, quote-business

**强制规范**:
- ✅ 所有新系统必须采用 V3.0 架构
- ✅ 现有系统逐步迁移到 V3.0
- ✅ 数据库统一使用 nexus_user

**命令参考**:
```bash
# SSH连接
ssh shanghai-tencent      # 上海服务器
ssh hangzhou-ali          # 杭州数据库服务器

# 服务管理
sudo systemctl status auth-center-backend
sudo systemctl restart auth-center-backend

# 数据库连接（通过SSH隧道）
psql -h localhost -p 5432 -U nexus_user -d auth_center_db
```

---

### 2. architecture.md - 系统架构与技术标准

**核心要点**:

#### V3.0 技术栈标准
**前端**: Vite 6+ + React 19+ + TypeScript 5+
**后端**: Go 1.21+ + Gin + GORM
**数据库**: PostgreSQL 15+ (UUID主键, snake_case命名)

**强制命名规范**:

数据库 (PostgreSQL):
- ✅ 表名: `snake_case`, 复数形式
- ✅ 列名: `snake_case`
- ✅ 主键: UUID (不是 Auto Increment INT)
- ✅ 时间戳: `{column}_at`
- ✅ 布尔值: `is_{adjective}` 或 `{verb}_ed`

Go 代码:
- ✅ 结构体名: `PascalCase` (单数)
- ✅ 字段名 (JSON): `camelCase`
- ✅ GORM 列映射: 必须使用 `column` 标签指定 snake_case

TypeScript/React:
- ✅ 组件名: `PascalCase`
- ✅ 文件名: 组件用 `PascalCase.tsx`, 工具用 `camelCase.ts`

**API 设计规范**:
- ✅ RESTful: 名词复数 `/api/users`
- ✅ HTTP 方法语义化: GET/POST/PUT/PATCH/DELETE
- ✅ 响应格式: `{success, data, pagination?, error?, errorCode?}`

**数据库设计强制规则**:
- ✅ 主键: `UUID PRIMARY KEY`
- ❌ 禁止: `SERIAL PRIMARY KEY`
- ✅ 外键命名: `{referenced_table}_{referenced_column}`
- ✅ 索引命名: `{table}_{column}_idx`
- ✅ 唯一约束: `{table}_{column}_key`

**常见错误**:
- ❌ 混淆: 用操作系统用户（ubuntu、root）连接数据库
- ✅ 正确: 用数据库用户（nexus_user）连接数据库

---

### 3. ssh-setup.md - SSH配置与密钥管理

**核心要点**:

#### 开发者本地SSH配置 (`~/.ssh/config`)
```bash
Host shanghai-tencent
    HostName 101.35.120.199
    User ubuntu                    # 操作系统用户
    IdentityFile ~/.ssh/xia_mac_shanghai_secure  # ED25519密钥
    StrictHostKeyChecking no
    ServerAliveInterval 60
    ServerAliveCountMax 3

Host hangzhou-ali
    HostName 47.110.82.96
    User root                      # 操作系统用户
    IdentityFile ~/.ssh/xia_mac_hangzhou_secure
    StrictHostKeyChecking no
    ServerAliveInterval 60
    ServerAliveCountMax 3
```

**强制规范**:
- ✅ 使用别名（`shanghai-tencent`, `hangzhou-ali`）
- ✅ 指定操作系统用户（上海用 `ubuntu`，其他用 `root`）
- ✅ 使用密钥认证
- ✅ 启用保活机制

**禁止方式**:
- ❌ 不要直接用IP: `ssh ubuntu@101.35.120.199`
- ❌ 不要每次输入密码
- ❌ 不要用不同的别名

**数据库连接配置**:
```bash
# 所有系统统一使用 nexus_user
postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@localhost:5432/数据库名?sslmode=disable
```

---

### 4. database-guide.md - 数据库连接与使用

**核心要点**:

#### 数据库连接方式对比

| 方式 | 连接字符串 | 优点 | 缺点 | 推荐度 |
|------|-----------|------|------|--------|
| **SSH隧道** | `localhost:5432` | ✅ 加密传输<br>✅ 密钥认证 | ❌ 多一个SSH进程 | ⭐⭐⭐⭐⭐ |
| **直连** | `47.110.82.96:5432` | ✅ 性能略好 | ❌ 需配置SSL<br>❌ 端口暴露 | ⭐⭐⭐ |

**⚠️ 安全警告**:
```bash
# ❌ 危险：直连且不使用SSL
DATABASE_URL=postgresql://nexus_user:pass@47.110.82.96:5432/db?sslmode=disable

# ✅ 正确：通过SSH隧道
DATABASE_URL=postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@localhost:5432/db?sslmode=disable
```

**统一配置**:
```bash
主机:   localhost (通过SSH隧道)
端口:   5432
用户:   nexus_user
密码:   hRJ9NSJApfeyFDraaDgkYowY
SSL模式: disable (SSH隧道已加密)
```

#### SSH隧道管理

**验证隧道状态**:
```bash
ssh shanghai-tencent "sudo systemctl status pg-tunnel"
# 预期: Active: active (running)
```

**故障修复**:
```bash
ssh shanghai-tencent << 'ENDSSH'
sudo systemctl start pg-tunnel
sudo systemctl enable pg-tunnel
sudo systemctl status pg-tunnel
ENDSSH
```

**常见错误排查**:
- `connection timeout` → SSH隧道未启动 → `sudo systemctl start pg-tunnel`
- `password authentication failed` → 密码错误 → 确认是 `hRJ9NSJApfeyFDraaDgkYowY`
- `connection refused` → 端口占用 → 检查 `sudo systemctl status pg-tunnel`
- `server does not support SSL` → 使用了 `sslmode=require` → 改为 `disable`

#### 数据库隔离策略

```
PostgreSQL Server (47.110.82.96:5432)
├─ auth_center_db        -- 账号中心（认证专用）
├─ pr_business_db        -- PR业务系统
├─ pixel_business_db     -- AI生图系统
└─ quote_business_db     -- 报价系统
```

**强制规则**:
- ✅ 每个业务系统使用独立数据库
- ✅ 不允许跨库查询（应用层Join）
- ✅ 通过 `auth_center_user_id` 关联用户身份

---

### 5. deployment-and-operations.md - 部署流程与服务管理

**核心要点**:

#### 统一部署规范（V4.0 - 强制执行）

**核心原则**:
1. **环境变量与代码分离** - 部署不覆盖配置
2. **目录命名统一** - 避免混乱和重复
3. **本地构建，上传产物** - 不在服务器构建
4. **配置模板化管理** - 使用 `.env.example`

**标准目录结构**:
```
/var/www/
├── {system-name}           # 后端目录
│   ├── {binary-name}       # 可执行文件
│   ├── .env                # 环境变量（服务器独立）
│   ├── .env.example        # 环境变量模板
│   └── logs/               # 日志目录
│
└── {system-name}-frontend  # 前端目录
    ├── index.html
    └── assets/
```

**目录命名标准**:

| 系统 | 后端目录 | 前端目录 | 二进制文件名 |
|------|---------|---------|-------------|
| PR业务 | `/var/www/pr-backend` | `/var/www/pr-frontend` | `pr-business` |
| Quote | `/var/www/quote-backend` | `/var/www/quote-frontend` | `quote-api` |
| auth-center | `/var/www/auth-center` | `/var/www/auth-center-frontend` | `auth-center-api` |

**命名规则**:
- ✅ 后端: `{system-name}` 或 `{system-name}-backend`
- ✅ 前端: `{system-name}-frontend`
- ❌ **禁止**: 版本号、日期、随意后缀

#### 环境变量管理规范（强制）

**核心原则**: 环境变量与代码分离，部署不覆盖配置！

**代码仓库**:
```
backend/
├── .env.example          # ✅ 环境变量模板（提交到 Git）
├── .env.local            # 本地开发（不提交）
└── .env.production       # ❌ 不要创建！避免误上传
```

**服务器上**:
```
/var/www/{system-name}/
├── {binary}              # 可执行文件
├── .env                  # ✅ 实际环境变量（首次手动创建）
└── .env.backup           # 自动备份
```

**✅ 正确做法**:
```bash
# 上传二进制文件
scp backend/{system-name}-api shanghai-tencent:/var/www/{system-name}/

# 重启服务
ssh shanghai-tencent "sudo systemctl restart {system-name}"

echo "⚠️ 环境变量未改变（如需修改请登录服务器）"
```

**❌ 错误做法**:
```bash
# ❌ 不要这样做！
scp .env.production shanghai-tencent:/var/www/{system-name}/.env
```

#### 标准环境变量模板

所有业务系统的 `.env.example` 应包含（详见 env.example.md）:
```bash
# 应用配置
APP_ENV=production
APP_PORT=8080
APP_NAME={System Name}
APP_DEBUG=false

# 数据库配置（通过 SSH 隧道）
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME={system_name}_db
DB_SSLMODE=disable

# Auth Center 配置
AUTH_CENTER_URL=https://os.crazyaigc.com
AUTH_CENTER_CALLBACK_URL=https://{domain}.com/api/v1/auth/callback

# 前端地址
FRONTEND_URL=https://{domain}.com

# JWT 配置
JWT_SECRET={CHANGE_THIS_IN_PRODUCTION}
JWT_ACCESS_TOKEN_EXPIRE=24h
JWT_REFRESH_TOKEN_EXPIRE=168h
```

#### 本地构建 vs 服务器构建对比

**为什么禁止在服务器上构建？**

| 对比项 | 本地构建 ✅ | 服务器构建 ❌ |
|--------|------------|-------------|
| **CPU 占用** | 本地 Mac（高性能） | 服务器 CPU（影响线上服务） |
| **磁盘占用** | 只上传构建产物 | 需要安装开发工具 + node_modules |
| **网络带宽** | 只上传必要文件 | 需要下载依赖包 |
| **构建时间** | 快（本地性能好） | 慢（服务器通常不如本地） |
| **环境一致性** | 可控（本地环境） | 不可控（服务器环境变化） |
| **安全性** | 服务器无需开发工具 | 需要安装 Node.js/Go 编译器 |

**具体数据对比**:

Vite 前端构建:
- 本地构建: node_modules ~300MB（本地）, 构建时间~10秒, 上传~2MB
- 服务器构建: node_modules ~300MB（服务器）, 构建时间~30-60秒, 占用~600MB+

Go 后端构建:
- 本地交叉编译: ~5秒, 上传~15-20MB
- 服务器编译: ~10-20秒, 需要Go工具链

**总结**:
- ✅ 本地构建: 快速、节省资源、安全、可控
- ❌ 服务器构建: 占用资源、浪费空间、增加攻击面、构建慢

#### 服务管理规范

**后端服务（Go）- 使用 Systemd**:
```bash
sudo systemctl start <service-name>
sudo systemctl stop <service-name>
sudo systemctl restart <service-name>
sudo systemctl status <service-name>
sudo journalctl -u <service-name> -f
sudo systemctl enable <service-name>
```

**前端服务（静态文件）- 使用 Nginx**:
```bash
# 本地构建
npm run build

# 上传构建产物
rsync -avz dist/ shanghai-tencent:/var/www/<app-name>/

# 测试并重载Nginx
ssh shanghai-tencent "sudo nginx -t"
ssh shanghai-tencent "sudo systemctl reload nginx"
```

**优势**:
- ✅ 性能最佳：Nginx 直接服务静态文件
- ✅ 部署简单：只需要上传构建产物
- ✅ 稳定性高：Nginx 成熟稳定
- ✅ 资源占用低：相比 Node.js 进程，内存占用极小

---

### 6. api.md - API接口说明与认证集成

**核心要点**:

#### 核心架构：三层账号模型

```
第1层: User（用户层）- 真实的人
├─ userId (UUID): 统一用户ID
├─ unionId (VARCHAR): 微信 UnionID，跨应用统一标识
└─ phoneNumber: 手机号（用于密码登录）

第2层: UserAccount（登录入口层）- 各端的 openid
├─ provider: 提供商（如 'wechat'）
├─ appId: 应用 AppID
├─ openId: 该应用下的 openid
└─ type: 登录类型（'web' | 'mp' | 'miniapp' | 'app'）

第3层: Session（会话层）- 登录会话管理
├─ token: JWT token（7天有效）
└─ expiresAt: 过期时间
```

**设计理念**:
```
unionid = 人（同一用户在不同应用）
openid = 登录入口（同一应用不同用户）
```

#### 关键API接口

**1. 发起微信登录**:
- URL: `GET /api/auth/wechat/login?callbackUrl=xxx`
- 响应: 跳转到微信授权页面

**2. 验证Token**:
- URL: `POST /api/auth/verify-token`
- 请求体: `{"token": "xxx"}`
- 响应: `{"success": true, "data": {"valid": true, "userId": "xxx"}}`

**3. 获取用户信息**:
- URL: `GET /api/auth/user-info`
- Header: `Authorization: Bearer <token>`
- 响应: 完整用户信息（包括 accounts 数组）

**4. 密码登录**:
- URL: `POST /api/auth/password/login`
- 请求体: `{"phoneNumber": "xxx", "password": "xxx"}`
- 响应: `{"success": true, "token": "xxx", "userId": "xxx"}`

**5. 获取会话列表**:
- URL: `GET /api/auth/sessions`
- Header: `Authorization: Bearer <token>`
- 响应: 所有活跃会话（设备信息、IP、过期时间）

#### 业务系统集成示例

**前端**:
```typescript
window.location.href = 'https://pr.crazyaigc.com/api/auth/wechat/login'
```

**后端**:
```go
// 1. 发起微信登录（代理接口）
func WechatLogin(c *gin.Context) {
    callbackUrl := "https://pr.crazyaigc.com/api/auth/callback"
    authCenterURL := fmt.Sprintf(
        "https://os.crazyaigc.com/api/auth/wechat/login?callbackUrl=%s",
        url.QueryEscape(callbackUrl),
    )
    c.Redirect(302, authCenterURL)
}

// 2. 接收回调
func AuthCallback(c *gin.Context) {
    userId := c.Query("userId")
    token := c.Query("token")

    // 验证token
    // 创建/获取本地用户
    // 设置session
    // 跳转到首页
    c.Redirect(302, "/dashboard")
}
```

**完整流程**:
1. 用户点击"微信登录" → 跳转到 `/api/auth/wechat/login`
2. 业务系统后端重定向到账号中心
3. 用户看到 `os.crazyaigc.com`（短暂显示，正常现象）
4. 跳转到 `open.weixin.qq.com`（微信授权页面）
5. 用户扫码/授权，微信回调到账号中心
6. 账号中心回调到业务系统 `/api/auth/callback?userId=xxx&token=xxx`
7. 业务系统验证token，创建本地用户，设置session
8. 跳转到 `/dashboard`，登录完成 ✅

#### 业务系统必需字段

```sql
CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY,
  auth_center_user_id UUID UNIQUE NOT NULL,  -- ✅ 关联账号中心
  role VARCHAR(50) DEFAULT 'USER',
  profile JSONB,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX users_auth_center_user_id_idx
  ON users(auth_center_user_id);
```

---

### 7. security.md - 安全规范与防护

**核心要点**:

#### 密钥管理规范（⚠️ 极其重要）

**密钥隔离原则（强制）**:

| 用途 | 密钥命名 | 存放位置 | 权限范围 |
|------|----------|----------|----------|
| 开发者登录 | `xia_mac_<server>_secure` | 本地 Mac | 该服务器 SSH 登录 |
| 服务器间隧道 | `<src>_<dest>_tunnel` | 源服务器仅一份 | 仅用于隧道，不能登录其他服务器 |
| CI/CD 部署 | `deploy_<service>` | CI 服务器 | 仅部署权限 |
| 数据库访问 | **禁止** | **禁止** | 使用应用层连接 |

**❌ 严禁操作**:
- 将本地登录密钥上传到服务器
- 在多个服务器间共享同一个私钥
- 将私钥放在代码仓库、配置文件中
- 使用私钥进行数据库连接

**✅ 正确做法**:
- 每个服务器/用途使用独立的 ED25519 密钥对
- 隧道密钥只留在源服务器，不用于其他目的
- 私钥权限设置为 `400` 或 `600`
- 定期（每季度）审计所有活跃密钥

#### 历史教训：2026-02-03密钥泄露事件

**影响范围**: 上海服务器 SSH 私钥泄露

**根本原因**:
1. 将本地登录密钥上传到服务器（用于 SSH 隧道）
2. 密钥在多个服务器间共享，缺乏隔离
3. 没有定期检查服务器访问日志
4. 缺少高风险操作的确认机制

**响应措施**:
1. 封禁攻击 IP：220.205.242.226, 204.76.203.83
2. 更换所有服务器 SSH 密钥
3. 更换数据库用户名和密码
4. 创建隧道专用密钥对
5. 更新安全策略文档

#### 高风险操作规范

**确认流程（强制）**:
1. **文档记录**：记录操作目的、影响范围、回滚方案
2. **所有者确认**：获得服务器/系统所有者的明确同意
3. **备份**：备份所有将被修改的文件
4. **测试**：在测试环境验证操作步骤
5. **时间窗口**：选择低峰期执行

**高风险操作表**:

| 操作 | 风险等级 | 确认要求 |
|------|----------|----------|
| 上传私钥到服务器 | 🔴 高 | **必须经所有者批准** |
| 删除服务器上的公钥 | 🔴 高 | **确认无影响** |
| 更换服务器密钥 | 🔴 高 | **全面测试连接** |
| 更换数据库密码 | 🔴 高 | **所有服务重启** |

#### 定期安全检查（强制执行）

**每周检查（周一 10:00）**:
```bash
# 检查服务器登录日志
sudo last -n 50
sudo journalctl -u ssh --since "7 days ago" | grep -E "Failed|Invalid"

# 检查活跃 SSH 连接
ss -tunp | grep ssh

# 检查进程异常
ps aux --sort=-%mem | head -20
ps aux --sort=-%cpu | head -20

# 检查磁盘空间
df -h

# 检查服务状态
systemctl list-units --state=failed
```

**每月检查（每月第一个周一）**:
```bash
# 审计所有活跃 SSH 密钥
for server in shanghai hangzhou hongkong; do
  ssh $server "cat ~/.ssh/authorized_keys && ssh-keygen -lf ~/.ssh/authorized_keys"
done

# 检查密钥有效期
for key in ~/.ssh/*_secure; do
  stat -f "%Sm" -t "%Y-%m-%d" $key
done

# 检查服务器访问日志
grep -E "220.205.242.226|204.76.203.83" /var/log/auth.log*
```

#### fail2ban 自动防护（已部署）

**部署状态**:
- 上海: ✅ 已安装，运行中
- 杭州: ✅ 已安装，运行中
- 香港: ✅ 已安装，运行中

**配置**:
```ini
[sshd]
enabled = true
port = 22
filter = sshd
logpath = /var/log/auth.log
maxretry = 5        # 5次失败后封禁
findtime = 600      # 10分钟内
bantime = 3600      # 封禁1小时
```

**防护原理**:
```
攻击者尝试密码登录
  ↓
第1次失败
  ↓
...
第5次失败（10分钟内）
  ↓
fail2ban: 自动封禁IP 1小时 🔒
```

**为什么不会误封禁你**:
1. ✅ 密码登录已禁用 → 不会触发"Failed password"
2. ✅ 密钥认证失败 → 不记录为失败，不触发封禁
3. ✅ 你的IP在白名单 → 永不封禁
4. ✅ 只有尝试**密码暴力破解**才会被封禁

**查看封禁状态**:
```bash
sudo fail2ban-client status sshd
sudo fail2ban-client set sshd unbanip <IP地址>
sudo tail -f /var/log/fail2ban.log
```

#### 安全检查清单

**新建服务器必须执行**:
- [ ] 更新系统：`apt update && apt upgrade`
- [ ] 安装安全工具：fail2ban, rkhunter
- [ ] 配置防火墙：只开放必要端口
- [ ] 禁用 root 远程登录（如适用）
- [ ] 配置 sudo 权限
- [ ] 设置密钥认证（禁用密码登录）
- [ ] 配置日志监控
- [ ] 备份关键配置
- [ ] 测试灾难恢复流程

**每月必须执行**:
- [ ] 审计所有 SSH 密钥
- [ ] 检查登录日志
- [ ] 更新系统补丁
- [ ] 检查磁盘空间
- [ ] 备份数据库
- [ ] 测试备份恢复
- [ ] 审查服务权限

**每季度必须执行**:
- [ ] 全面安全审计
- [ ] 密钥轮换
- [ ] 灾难恢复演练
- [ ] 安全培训
- [ ] 策略更新

---

### 8. deploy-template.sh ⭐ 新增文件

**核心要点**:

这是 KeenChase 统一部署脚本模板，所有业务系统必须使用此模板结构。

**使用说明**:
1. 复制此文件到业务系统根目录
2. 重命名为 `deploy-production.sh`
3. 修改配置项（SYSTEM_NAME, BINARY_NAME, DOMAIN）
4. 添加可执行权限：`chmod +x deploy-production.sh`

**重要**:
- 此脚本不会上传 `.env` 文件（环境变量与代码分离）
- 只上传二进制文件和前端静态文件

**关键配置**:
```bash
SYSTEM_NAME="{system-name}"
BINARY_NAME="{system-name}-api"
DOMAIN="{domain}.com"
SERVER="shanghai-tencent"
REMOTE_DIR="/var/www/${SYSTEM_NAME}"
```

**部署流程**:
1. 前端部署: 本地构建 → 上传静态文件 → 完成
2. 后端部署: 交叉编译 → 上传二进制 → 重启服务
3. 验证部署: 健康检查 → 报告状态

**前端部署部分**:
```bash
cd frontend

# 检查是否已构建
if [ ! -d "dist" ]; then
  npm run build
fi

# 上传静态文件
rsync -avz --delete \
  --exclude '*.map' \
  --exclude '*.html.gz' \
  dist/ \
  ${SERVER}:${REMOTE_DIR}-frontend/
```

**后端部署部分**:
```bash
cd ../backend

# 交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w" \
  -o ${BINARY_NAME} \
  cmd/server/main.go

# 验证二进制文件
file ${BINARY_NAME} | grep -q "ELF 64-bit"

# 上传二进制
scp ${BINARY_NAME} ${SERVER}:${REMOTE_DIR}/

# 重启服务（远程执行）
ssh ${SERVER} << ENDSSH
cd ${REMOTE_DIR}

# 备份旧二进制
if [ -f ${BINARY_NAME} ]; then
  mv ${BINARY_NAME} ${BINARY_NAME}.backup.$(date +%Y%m%d_%H%M%S)
fi

# 启用新版本
mv ${BINARY_NAME}-new ${BINARY_NAME}

# 重启服务
sudo systemctl restart ${SYSTEM_NAME}

# 等待启动
sleep 3

# 检查状态
sudo systemctl status ${SYSTEM_NAME} --no-pager
ENDSSH
```

**健康检查部分**:
```bash
sleep 2

# 健康检查
if curl -f -s https://${DOMAIN}/health > /dev/null; then
  echo "✅ 健康检查通过"
else
  echo "⚠️ 健康检查失败，请手动验证"
fi
```

**关键改进点**（相比旧脚本）:
- ✅ 自动检查 dist 目录是否存在
- ✅ 验证二进制文件格式
- ✅ 自动备份旧版本
- ✅ 详细的错误提示
- ✅ 健康检查验证
- ✅ 环境变量提示

---

### 9. env.example.md ⭐ 新增文件

**核心要点**:

这是 KeenChase 标准环境变量模板，所有业务系统的 `.env.example` 应遵循此格式。

**使用说明**:
1. 复制此文件到项目 backend 目录
2. 重命名为 `.env.example`
3. 修改 `{system-name}`, `{domain}` 等占位符
4. 提交到 Git 仓库
5. 服务器首次部署时，使用此模板创建实际的 `.env` 文件

**标准模板**:
```bash
# ============================================
# 应用配置
# ============================================
APP_ENV=production
APP_PORT=8080
APP_NAME={System Name}
APP_DEBUG=false

# ============================================
# 数据库配置（通过 SSH 隧道）
# ============================================
# ⚠️ 重要配置说明：
# 1. DB_HOST 必须是 localhost（通过 SSH 隧道转发到杭州服务器）
# 2. DB_USER 所有系统统一使用 nexus_user
# 3. DB_PASSWORD 统一使用 hRJ9NSJApfeyFDraaDgkYowY
# 4. DB_NAME 每个系统不同，如 pr_business_db, quote_business_db
# 5. DB_SSLMODE=disable 因为 SSH 隧道已加密
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME={system_name}_db
DB_SSLMODE=disable

# ============================================
# Auth Center 配置
# ============================================
AUTH_CENTER_URL=https://os.crazyaigc.com
AUTH_CENTER_CALLBACK_URL=https://{domain}.com/api/v1/auth/callback

# ============================================
# 前端地址
# ============================================
FRONTEND_URL=https://{domain}.com

# ============================================
# JWT 配置
# ============================================
# ⚠️ 生产环境必须修改为随机字符串！
# 生成命令：openssl rand -base64 32
JWT_SECRET={CHANGE_THIS_IN_PRODUCTION}
JWT_ACCESS_TOKEN_EXPIRE=24h
JWT_REFRESH_TOKEN_EXPIRE=168h

# ============================================
# 日志配置
# ============================================
LOG_LEVEL=info
LOG_FORMAT=json
```

**关键改进点**:
- ✅ 详细的数据配置说明（注释中）
- ✅ 明确指出 SSH 隧道的使用
- ✅ 提供JWT密钥生成命令
- ✅ 清晰的分组结构
- ✅ 占位符替换指导

---

## 🔍 跨文档关键规范汇总

### 强制规范（所有文档共同要求）

#### 1. 数据库连接规范
```bash
# ✅ 唯一正确的方式
DB_HOST=localhost                           # 通过SSH隧道
DB_PORT=5432
DB_USER=nexus_user                         # 统一用户
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY       # 统一密码
DB_NAME={system_name}_db                   # 每个系统不同
DB_SSLMODE=disable                         # SSH隧道已加密

# ❌ 禁止的方式
DB_HOST=47.110.82.96                       # 直连不安全
DB_USER=pr_business_user                   # 已废弃
DB_SSLMODE=require                         # SSH隧道不需要
```

#### 2. SSH配置规范
```bash
# ~/.ssh/config
Host shanghai-tencent
    HostName 101.35.120.199
    User ubuntu                           # 操作系统用户
    IdentityFile ~/.ssh/xia_mac_shanghai_secure
    StrictHostKeyChecking no
    ServerAliveInterval 60

Host hangzhou-ali
    HostName 47.110.82.96
    User root                             # 操作系统用户
    IdentityFile ~/.ssh/xia_mac_hangzhou_secure
    StrictHostKeyChecking no
    ServerAliveInterval 60
```

#### 3. 目录命名规范
```
/var/www/
├── {system-name}           # 后端（不带版本、日期）
└── {system-name}-frontend  # 前端（不带版本、日期）

# 示例
/var/www/pr-backend
/var/www/pr-frontend
/var/www/quote-backend
/var/www/quote-frontend

# ❌ 禁止
/var/www/pr-backend-v2
/var/www/pr-backend-20260204
```

#### 4. 部署流程规范
```bash
# ✅ 正确流程
1. 本地构建（前端: npm run build, 后端: go build）
2. 上传产物（前端: rsync dist/, 后端: scp binary）
3. 重启服务（systemctl restart）
4. 健康检查（curl /health）

# ❌ 禁止流程
- 在服务器上运行 npm install
- 在服务器上运行 go build
- 上传 .env 文件（环境变量与代码分离）
```

#### 5. 环境变量管理规范
```
代码仓库:
├── .env.example          # ✅ 提交到 Git
└── .env.local            # ❌ 不提交

服务器:
└── .env                  # ✅ 首次手动创建，之后不覆盖

# ❌ 禁止
- 创建 .env.production 并上传
- 部署时覆盖服务器 .env
```

### 禁止事项汇总

**❌ 所有文档共同禁止的操作**:

1. **数据库连接**:
   - ❌ 直连 47.110.82.96:5432 且不用SSL
   - ❌ 使用已废弃的专用数据库用户（如 pr_business_user）
   - ❌ 混淆操作系统用户和数据库用户

2. **SSH配置**:
   - ❌ 直接用IP登录（不使用别名）
   - ❌ 将本地登录密钥上传到服务器
   - ❌ 在多个服务器间共享同一个私钥
   - ❌ 将私钥放在代码仓库、配置文件中

3. **部署操作**:
   - ❌ 在服务器上运行构建命令（npm run build, go build）
   - ❌ 上传 .env 文件覆盖服务器配置
   - ❌ 创建带版本号、日期的目录
   - ❌ 使用 .env.production 文件

4. **安全规范**:
   - ❌ 密码认证登录服务器（已禁用）
   - ❌ 将私钥用于数据库连接
   - ❌ 在多个服务器间共享同一个私钥

5. **代码规范**:
   - ❌ 使用自增 INT 作为主键（必须用UUID）
   - ❌ 数据库使用 camelCase（必须用snake_case）
   - ❌ Go 结构体使用 camelCase（必须用PascalCase）
   - ❌ 跨库查询（应用层Join）

### 警告事项汇总

**⚠️ 需要特别注意的事项**:

1. **环境变量分离**:
   - ⚠️ 部署脚本不应上传 .env 文件
   - ⚠️ 服务器 .env 需要手动创建和维护
   - ⚠️ 修改环境变量需要登录服务器手动编辑

2. **SSH隧道依赖**:
   - ⚠️ 数据库连接依赖 SSH 隧道
   - ⚠️ 隧道停止会导致数据库连接失败
   - ⚠️ 需要定期检查隧道状态

3. **密钥管理**:
   - ⚠️ 每个用途使用独立的密钥对
   - ⚠️ 密钥需要定期轮换（登录: 6个月, 数据库: 3个月）
   - ⚠️ 定期审计活跃密钥

4. **微信登录流程**:
   - ⚠️ 用户会短暂看到 os.crazyaigc.com（正常现象）
   - ⚠️ 这是所有第三方登录的标准流程
   - ⚠️ 不要尝试用 iframe 隐藏（微信安全限制）

5. **安全检查**:
   - ⚠️ fail2ban 已部署，密码暴力破解会被自动封禁
   - ⚠️ 密钥认证失败不会触发封禁
   - ⚠️ 你的IP在白名单中不会被误封

---

## 📊 PR Business 项目对比

### 当前状态分析

#### ✅ 符合标准的部分

1. **目录结构**:
   - ✅ 后端目录: `/var/www/pr-backend`
   - ✅ 前端目录: `/var/www/pr-frontend`
   - ✅ 使用标准命名（无版本号、日期）

2. **服务管理**:
   - ✅ 使用 systemd 管理后端服务
   - ✅ 使用 Nginx 服务前端静态文件

3. **技术栈**:
   - ✅ 前端: Vite + React
   - ✅ 后端: Go + Gin

#### ❌ 不符合标准的部分

1. **数据库配置错误**（严重问题）:
   ```bash
   # 当前 .env.production (错误)
   DB_HOST=47.110.82.96                    # ❌ 直连，不安全
   DB_USER=nexus                            # ❌ 应该是 nexus_user
   DB_PASSWORD=nexus123                     # ❌ 密码错误
   DB_SSLMODE=disable                       # ❌ 直连必须用 require

   # 应该改为（正确）
   DB_HOST=localhost                        # ✅ 通过SSH隧道
   DB_USER=nexus_user                       # ✅ 统一用户
   DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY     # ✅ 统一密码
   DB_SSLMODE=disable                       # ✅ SSH隧道已加密
   ```

2. **部署脚本问题**:
   ```bash
   # 当前 deploy-production.sh (第31行)
   scp .env.production shanghai-tencent:/var/www/pr-backend/.env
   # ❌ 违反"环境变量与代码分离"原则

   # 应该删除这行，改为：
   echo "⚠️ 环境变量未改变（如需修改请登录服务器）"
   ```

3. **缺少 .env.example 文件**:
   - ❌ 没有标准的环境变量模板
   - ❌ 不符合 deploy-template.sh 的规范

4. **配置代码默认值问题**:
   ```go
   // config.go 第75行
   viper.SetDefault("DB_USER", "pr_business")  // ❌ 应该是 nexus_user
   ```

### 需要修复的问题清单（按优先级排序）

#### 🔴 优先级1：严重问题（必须立即修复）

1. **数据库连接配置错误**
   - 文件: `backend/.env.production`
   - 问题: 使用直连、错误的用户名和密码
   - 风险: 不安全、连接失败
   - 修复步骤: 见下方详细步骤

2. **部署脚本违反环境变量分离原则**
   - 文件: `deploy-production.sh`
   - 问题: 第31行上传 .env 文件
   - 风险: 可能覆盖服务器配置
   - 修复步骤: 删除第31行

#### 🟡 优先级2：重要问题（应该尽快修复）

3. **缺少 .env.example 文件**
   - 文件: `backend/.env.example`（需要创建）
   - 问题: 没有标准模板
   - 风险: 不符合规范，新部署困难
   - 修复步骤: 复制 env.example.md 模板

4. **配置代码默认值错误**
   - 文件: `backend/config/config.go`
   - 问题: 第75行默认值是 `pr_business`
   - 风险: 可能误导其他开发者
   - 修复步骤: 修改默认值为 `nexus_user`

#### 🟢 优先级3：改进建议（可选）

5. **部署脚本不够完善**
   - 文件: `deploy-production.sh`
   - 问题: 缺少健康检查、二进制验证等
   - 建议: 参考 deploy-template.sh 完善

6. **缺少 SSH 隧道状态检查**
   - 问题: 部署前没有验证 SSH 隧道
   - 建议: 添加隧道状态检查

---

## 🔧 详细修复步骤和命令

### 问题1: 修复数据库连接配置

**当前配置** (`backend/.env.production`):
```bash
DB_HOST=47.110.82.96
DB_USER=nexus
DB_PASSWORD=nexus123
DB_SSLMODE=disable
```

**修复后的配置**:
```bash
DB_HOST=localhost
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME=pr_business_db
DB_SSLMODE=disable
```

**修复步骤**:

1. **备份当前配置**:
   ```bash
   cd /Users/xia/Documents/GitHub/pr-business/backend
   cp .env.production .env.production.backup.$(date +%Y%m%d)
   ```

2. **修改 .env.production 文件**:
   ```bash
   cat > .env.production << 'EOF'
   APP_ENV=production
   APP_PORT=8081
   APP_NAME="PR Business"
   APP_DEBUG=false

   # ============================================
   # 数据库配置（通过 SSH 隧道）
   # ============================================
   # ⚠️ 重要：通过 SSH 隧道连接杭州服务器
   # 隧道服务: pg-tunnel.service
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=nexus_user
   DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
   DB_NAME=pr_business_db
   DB_SSLMODE=disable

   # ============================================
   # Auth Center配置
   # ============================================
   AUTH_CENTER_URL=https://os.crazyaigc.com
   AUTH_CENTER_REDIRECT_URI=https://pr.crazyaigc.com/api/v1/auth/callback

   # ============================================
   # 前端配置
   # ============================================
   FRONTEND_URL=https://pr.crazyaigc.com

   # ============================================
   # JWT配置
   # ============================================
   JWT_SECRET=151jmeLlr7ZSi9L4KXIhrJ/CfTFBY2PV5CezmfUlLzw=
   JWT_ACCESS_TOKEN_EXPIRE=24h
   JWT_REFRESH_TOKEN_EXPIRE=168h
   EOF
   ```

3. **验证修改**:
   ```bash
   cat .env.production | grep -E "DB_HOST|DB_USER|DB_PASSWORD|DB_NAME|DB_SSLMODE"
   ```

4. **更新服务器配置**（需要SSH登录服务器）:
   ```bash
   # 登录上海服务器
   ssh shanghai-tencent

   # 备份服务器配置
   sudo cp /var/www/pr-backend/.env /var/www/pr-backend/.env.backup.$(date +%Y%m%d)

   # 编辑服务器配置
   sudo nano /var/www/pr-backend/.env
   # 修改为:
   # DB_HOST=localhost
   # DB_USER=nexus_user
   # DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
   # DB_NAME=pr_business_db
   # DB_SSLMODE=disable

   # 重启服务
   sudo systemctl restart pr-backend

   # 检查服务状态
   sudo systemctl status pr-backend

   # 查看日志，确认数据库连接成功
   sudo journalctl -u pr-backend -f
   ```

5. **验证数据库连接**:
   ```bash
   # 在上海服务器上测试SSH隧道
   ssh shanghai-tencent << 'ENDSSH'
   # 检查SSH隧道状态
   sudo systemctl status pg-tunnel

   # 测试数据库连接
   PGPASSWORD=hRJ9NSJApfeyFDraaDgkYowY psql -h localhost -p 5432 -U nexus_user -d pr_business_db -c 'SELECT 1;'

   # 应该输出:
   #  ?column?
   # ----------
   #         1
   ENDSSH
   ```

---

### 问题2: 修复部署脚本

**当前问题** (`deploy-production.sh` 第31行):
```bash
scp .env.production shanghai-tencent:/var/www/pr-backend/.env
```

**修复步骤**:

1. **备份当前脚本**:
   ```bash
   cd /Users/xia/Documents/GitHub/pr-business
   cp deploy-production.sh deploy-production.sh.backup.$(date +%Y%m%d)
   ```

2. **删除第31行**:
   ```bash
   # 方式1: 使用 sed
   sed -i '' '31d' deploy-production.sh

   # 方式2: 手动编辑
   nano deploy-production.sh
   # 找到并删除: scp .env.production shanghai-tencent:/var/www/pr-backend/.env
   ```

3. **在脚本末尾添加提示**:
   ```bash
   cat >> deploy-production.sh << 'EOF'

echo ""
echo "⚠️ 环境变量未改变（如需修改请登录服务器）"
echo "修改命令：ssh shanghai-tencent \"sudo nano /var/www/pr-backend/.env\""
echo ""
EOF
   ```

4. **验证修改**:
   ```bash
   # 检查是否已删除 .env 上传行
   grep -n ".env.production" deploy-production.sh
   # 应该没有输出（表示已删除）

   # 检查是否有新提示
   grep "环境变量未改变" deploy-production.sh
   # 应该显示: echo "⚠️ 环境变量未改变（如需修改请登录服务器）"
   ```

---

### 问题3: 创建 .env.example 文件

**创建步骤**:

1. **创建标准模板**:
   ```bash
   cd /Users/xia/Documents/GitHub/pr-business/backend

   cat > .env.example << 'EOF'
   # ============================================
   # 应用配置
   # ============================================
   APP_ENV=production
   APP_PORT=8081
   APP_NAME="PR Business"
   APP_DEBUG=false

   # ============================================
   # 数据库配置（通过 SSH 隧道）
   # ============================================
   # ⚠️ 重要配置说明：
   # 1. DB_HOST 必须是 localhost（通过 SSH 隧道转发到杭州服务器）
   # 2. DB_USER 所有系统统一使用 nexus_user
   # 3. DB_PASSWORD 统一使用 hRJ9NSJApfeyFDraaDgkYowY
   # 4. DB_NAME 每个系统不同，这里是 pr_business_db
   # 5. DB_SSLMODE=disable 因为 SSH 隧道已加密
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=nexus_user
   DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
   DB_NAME=pr_business_db
   DB_SSLMODE=disable

   # ============================================
   # Auth Center 配置
   # ============================================
   AUTH_CENTER_URL=https://os.crazyaigc.com
   AUTH_CENTER_CALLBACK_URL=https://pr.crazyaigc.com/api/v1/auth/callback

   # ============================================
   # 前端地址
   # ============================================
   FRONTEND_URL=https://pr.crazyaigc.com

   # ============================================
   # JWT 配置
   # ============================================
   # ⚠️ 生产环境必须修改为随机字符串！
   # 生成命令：openssl rand -base64 32
   JWT_SECRET={CHANGE_THIS_IN_PRODUCTION}
   JWT_ACCESS_TOKEN_EXPIRE=24h
   JWT_REFRESH_TOKEN_EXPIRE=168h

   # ============================================
   # 日志配置
   # ============================================
   LOG_LEVEL=info
   LOG_FORMAT=json
   EOF
   ```

2. **验证文件**:
   ```bash
   cat .env.example
   ```

3. **提交到 Git**:
   ```bash
   git add .env.example
   git commit -m "chore: 添加标准环境变量模板 (.env.example)"
   ```

---

### 问题4: 修复配置代码默认值

**修复步骤**:

1. **备份当前文件**:
   ```bash
   cd /Users/xia/Documents/GitHub/pr-business/backend
   cp config/config.go config/config.go.backup.$(date +%Y%m%d)
   ```

2. **修改第75行**:
   ```bash
   # 使用 sed 修改
   sed -i '' 's/DB_USER", "pr_business"/DB_USER", "nexus_user"/' config/config.go

   # 或手动编辑
   nano config/config.go
   # 找到第75行，修改为:
   # viper.SetDefault("DB_USER", "nexus_user")
   ```

3. **验证修改**:
   ```bash
   grep -n "DB_USER" config/config.go
   # 应该显示:
   # viper.SetDefault("DB_USER", "nexus_user")
   ```

4. **测试编译**:
   ```bash
   cd /Users/xia/Documents/GitHub/pr-business/backend
   go build -o test-binary ./cmd/server
   rm test-binary
   ```

---

### 问题5: 完善部署脚本（可选）

**建议参考 deploy-template.sh 完善**:

1. **添加前端构建检查**:
   ```bash
   # 检查是否已构建
   if [ ! -d "dist" ]; then
     echo "未发现 dist 目录，开始构建..."
     npm run build
   fi
   ```

2. **添加二进制验证**:
   ```bash
   # 验证二进制文件
   file ${BINARY_NAME} | grep -q "ELF 64-bit" || {
     echo "❌ 编译失败：不是有效的 Linux 二进制文件"
     exit 1
   }
   ```

3. **添加健康检查**:
   ```bash
   # 健康检查
   if curl -f -s https://pr.crazyaigc.com/health > /dev/null; then
     echo "✅ 健康检查通过"
   else
     echo "⚠️ 健康检查失败，请手动验证"
   fi
   ```

---

### 问题6: 添加 SSH 隧道状态检查（可选）

**在部署前添加检查**:

```bash
# 检查SSH隧道状态
echo "🔍 检查 SSH 隧道状态..."
ssh shanghai-tencent "sudo systemctl status pg-tunnel" | grep -q "active (running)" || {
  echo "❌ SSH隧道未运行，请先启动："
  echo "   ssh shanghai-tencent 'sudo systemctl start pg-tunnel'"
  exit 1
}
echo "✅ SSH隧道运行正常"
```

---

## 📋 完整修复清单

### 立即执行（今天完成）

- [ ] 1. 修复数据库连接配置（本地 .env.production）
- [ ] 2. 更新服务器配置（SSH登录服务器修改 /var/www/pr-backend/.env）
- [ ] 3. 删除部署脚本中的 .env 上传行
- [ ] 4. 创建 .env.example 文件
- [ ] 5. 修复配置代码默认值
- [ ] 6. 测试数据库连接（验证 SSH 隧道）
- [ ] 7. 重启后端服务并检查日志
- [ ] 8. 提交代码修改到 Git

### 可选改进（本周完成）

- [ ] 9. 完善部署脚本（参考 deploy-template.sh）
- [ ] 10. 添加 SSH 隧道状态检查
- [ ] 11. 添加部署前的预检查
- [ ] 12. 编写部署文档

---

## 🎓 学习总结

### 关键收获

1. **环境变量与代码分离**的重要性
   - 避免配置覆盖
   - 提高安全性
   - 便于独立管理

2. **SSH隧道**的必要性和正确使用
   - 安全性更高
   - 配置简单
   - 已有成熟方案

3. **标准化部署**的价值
   - 减少错误
   - 提高效率
   - 便于维护

4. **密钥管理**的严格规范
   - 密钥隔离
   - 定期轮换
   - 最小权限

5. **自动化防护**的重要性
   - fail2ban 自动封禁
   - 减少人工干预
   - 提高安全性

### 最容易犯的错误

1. **混淆操作系统用户和数据库用户**
   - 操作系统用户: ubuntu, root（SSH登录）
   - 数据库用户: nexus_user（PostgreSQL连接）

2. **数据库连接配置错误**
   - ❌ 直连 47.110.82.96:5432
   - ✅ 通过SSH隧道: localhost:5432

3. **部署时上传 .env 文件**
   - ❌ scp .env.production server:/path/.env
   - ✅ 只上传二进制，环境变量独立管理

4. **使用已废弃的专用数据库用户**
   - ❌ pr_business_user
   - ✅ nexus_user（统一用户）

5. **在服务器上构建**
   - ❌ ssh server "npm run build"
   - ✅ 本地构建，上传产物

### 最佳实践建议

1. **遵循标准模板**
   - 使用 deploy-template.sh 作为部署脚本基础
   - 使用 env.example.md 作为环境变量模板
   - 严格按照标准目录命名

2. **本地构建，上传产物**
   - 前端: 本地 `npm run build` → `rsync dist/`
   - 后端: 本地 `go build` → `scp binary`

3. **环境变量独立管理**
   - 代码仓库: 只有 `.env.example`
   - 服务器: 手动创建 `.env`，部署不覆盖

4. **定期检查和维护**
   - 每周检查 SSH 隧道状态
   - 每月审计 SSH 密钥
   - 每季度轮换密钥

5. **安全性优先**
   - 使用 SSH 隧道连接数据库
   - fail2ban 自动防护
   - 密钥隔离和定期轮换

---

## 📞 需要帮助的场景

如果遇到问题，可以：

1. **检查文档的"常见问题"章节**
2. **检查服务日志**: `sudo journalctl -u <service-name> -f`
3. **查看系统日志**: `sudo tail -f /var/log/syslog`
4. **验证 SSH 隧道**: `sudo systemctl status pg-tunnel`
5. **测试数据库连接**: `psql -h localhost -p 5432 -U nexus_user -d pr_business_db`

---

## 📚 附录：快速参考

### 常用命令

```bash
# === SSH 连接 ===
ssh shanghai-tencent      # 上海服务器
ssh hangzhou-ali          # 杭州数据库服务器

# === 服务管理 ===
sudo systemctl status pr-backend
sudo systemctl restart pr-backend
sudo journalctl -u pr-backend -f

# === Nginx 管理 ===
sudo nginx -t
sudo systemctl reload nginx

# === 数据库连接（通过 SSH 隧道）===
psql -h localhost -p 5432 -U nexus_user -d pr_business_db

# === fail2ban 状态 ===
sudo fail2ban-client status sshd
```

### 标准配置

```bash
# 数据库连接
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME={system_name}_db
DB_SSLMODE=disable

# Auth Center
AUTH_CENTER_URL=https://os.crazyaigc.com
AUTH_CENTER_CALLBACK_URL=https://{domain}.com/api/v1/auth/callback

# JWT
JWT_SECRET={随机32字符以上}
JWT_ACCESS_TOKEN_EXPIRE=24h
JWT_REFRESH_TOKEN_EXPIRE=168h
```

### 目录结构

```
/var/www/
├── pr-backend              # PR业务后端
│   ├── pr-business         # 可执行文件
│   ├── .env                # 环境变量（不覆盖）
│   └── logs/               # 日志目录
│
└── pr-frontend             # PR业务前端
    ├── index.html
    └── assets/
```

---

**学习完成时间**: 2026-02-04
**下次审查**: 2026-03-04（建议每月重新学习一次）
**维护者**: KeenChase Dev Team
