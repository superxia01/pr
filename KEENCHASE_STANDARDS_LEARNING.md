# KeenChase 标准规范学习总结报告

**学习时间**: 2026-02-04
**学习者**: xia
**标准版本**: V3.7
**学习目标**: 系统学习 KeenChase 技术规范并对比 PR Business 项目现状

---

## 📚 学习文档清单

本次学习按照推荐顺序完成了所有标准文档的学习：

1. ✅ README.md - 总览与导航
2. ✅ architecture.md - 系统架构与技术标准
3. ✅ ssh-setup.md - SSH配置指南
4. ✅ database-guide.md - 数据库使用指南
5. ✅ deployment-and-operations.md - 部署与运维
6. ✅ api.md - API接口说明
7. ✅ security.md - 安全规范

---

## 📖 文档1: README.md - 总览

### 核心概念

**KeenChase 技术规范**是所有业务系统的"宪法"级别的技术标准，**强制执行**。

#### 关键架构概念

1. **V3.0 架构**（当前标准）
   - 前后端分离架构
   - 前端：Vite + React（推荐）/ Next.js（仅SSR场景）
   - 后端：Go + Gin
   - 数据库：PostgreSQL 15+（统一数据层）

2. **分布式部署架构**
   - **上海服务器** (101.35.120.199): 应用服务器
     - 运行所有 Go 后端服务
     - Nginx 服务前端静态文件
     - 用户: ubuntu（操作系统用户）

   - **杭州服务器** (47.110.82.96): 数据库服务器
     - 运行 PostgreSQL 15 (Docker)
     - 统一数据层
     - 用户: root（操作系统用户）

3. **已部署系统**
   - 账号中心: os.crazyaigc.com ✅
   - PR业务: pr.crazyaigc.com ✅
   - AI生图: pixel.crazyaigc.com ✅
   - 报价系统: quote.crazyaigc.com ✅
   - 官网: www.crazyaigc.com (Vercel) ✅

### 关键规范

**强制要求**：
- 所有新系统必须采用 V3.0 架构
- 所有系统必须遵循统一的技术栈标准
- 所有系统必须集成账号中心认证

### 常用命令

```bash
# SSH 连接
ssh shanghai-tencent      # 上海服务器
ssh hangzhou-ali          # 杭州数据库服务器

# 服务管理
sudo systemctl status auth-center-backend
sudo systemctl restart auth-center-backend
sudo journalctl -u auth-center-backend -f

# Nginx 管理
sudo nginx -t
sudo systemctl reload nginx

# 数据库连接（通过 SSH 隧道）
psql -h localhost -p 5432 -U nexus_user -d auth_center_db

# fail2ban 状态
sudo fail2ban-client status sshd
```

### 注意事项

⚠️ **特别重要**：
- 必须区分**操作系统用户**（ubuntu/root）和**数据库用户**（nexus_user）
- 上海服务器用 ubuntu 用户，杭州/香港用 root 用户
- 数据库连接必须通过 SSH 隧道

---

## 📖 文档2: architecture.md - 系统架构与技术标准

### 核心概念

#### 1. 用户类型区分（极其重要）

**操作系统用户**（OS User）：
- 用于 SSH 登录服务器
- 上海: `ubuntu`（普通用户）
- 杭州: `root`（管理员用户）

**数据库用户**（Database User）：
- 用于 PostgreSQL 连接
- 统一使用: `nexus_user`（超级用户）
- 密码: `hRJ9NSJApfeyFDraaDgkYowY`

**常见错误**：
- ❌ 用操作系统用户（ubuntu、root）连接数据库
- ❌ 用数据库用户（nexus_user）SSH 登录服务器

#### 2. V3.0 技术栈标准

**前端技术栈**（Vite + React 推荐）：
```
✅ 构建工具: Vite 6+
✅ 框架: React 19+
✅ 语言: TypeScript 5+
✅ 路由: React Router 6+
✅ 样式: Tailwind CSS
✅ 状态管理: Zustand / React Context
✅ HTTP 客户端: Axios / Fetch API
✅ 组件库: Radix UI / shadcn/ui / Material-UI
✅ 表单处理: React Hook Form
```

**后端技术栈**：
```
✅ 语言: Go 1.21+
✅ 框架: Gin
✅ ORM: GORM
✅ 数据库驱动: PostgreSQL
✅ 认证: JWT (golang-jwt/jwt/v5)
✅ 密码加密: bcrypt
✅ 配置管理: godotenv
✅ 日志: Zap（推荐）
```

**数据库标准**：
```
✅ 数据库: PostgreSQL 15+
✅ 主键类型: UUID（不是 Auto Increment INT）
✅ 列命名: snake_case（不是 camelCase）
✅ 时间戳: timestamp with time zone
✅ JSON 字段: JSONB
```

### 关键规范

#### 1. 数据库命名规范（强制）

**表名**: `snake_case`，复数形式
```sql
users          -- ✅ 正确
user_accounts  -- ✅ 正确
userAccounts   -- ❌ 错误
User           -- ❌ 错误
```

**列名**: `snake_case`，全部小写
```sql
user_id        -- ✅ 正确
created_at     -- ✅ 正确
phone_number   -- ✅ 正确
userId         -- ❌ 错误
CreatedAt      -- ❌ 错误
```

**主键**: UUID
```sql
-- 方式1: 表名_id（推荐用于外键）
users.user_id          UUID PRIMARY KEY
user_accounts.id       UUID PRIMARY KEY

-- 方式2: id（推荐用于主表）
users.id               UUID PRIMARY KEY
user_accounts.user_id  UUID REFERENCES users(id)
```

**时间戳**: `{column}_at`
```sql
created_at     TIMESTAMP WITH TIME ZONE
updated_at     TIMESTAMP WITH TIME ZONE
deleted_at     TIMESTAMP WITH TIME ZONE
expires_at     TIMESTAMP WITH TIME ZONE
```

**布尔值**: `is_{adjective}` 或 `{verb}_ed`
```sql
is_active      BOOLEAN
is_verified    BOOLEAN
is_deleted     BOOLEAN
published      BOOLEAN
```

#### 2. Go 代码命名规范（强制）

**结构体名**: `PascalCase`（单数）
```go
type User struct { }         // ✅ 正确
type UserAccount struct { }  // ✅ 正确
type user struct { }         // ❌ 错误
```

**字段名（JSON）**: `camelCase`（导出字段）
```go
type User struct {
    UserID       string    `json:"userId"`        // ✅ 正确
    PhoneNumber  string    `json:"phoneNumber"`   // ✅ 正确
    CreatedAt    time.Time `json:"createdAt"`     // ✅ 正确
}
```

**GORM 列映射**: **必须**使用 `column` 标签指定 snake_case
```go
type User struct {
    UserID       string    `gorm:"primaryKey;column:user_id;type:uuid" json:"userId"`
    UnionID      string    `gorm:"uniqueIndex;column:union_id;type:varchar(255)" json:"unionId"`
    PhoneNumber  string    `gorm:"uniqueIndex;column:phone_number;type:varchar(255)" json:"phoneNumber"`
    CreatedAt    time.Time `gorm:"column:created_at;type:timestamp with time zone" json:"createdAt"`
}
```

#### 3. API 设计规范（强制）

**RESTful API 标准**：
```
# 用户资源
GET    /api/users              - 获取用户列表 (分页)
GET    /api/users/:id          - 获取单个用户
POST   /api/users              - 创建用户
PUT    /api/users/:id          - 完整更新用户
PATCH  /api/users/:id          - 部分更新用户
DELETE /api/users/:id          - 删除用户

# 特殊操作
POST   /api/auth/login         - 登录
POST   /api/auth/logout        - 登出
POST   /api/users/:id/verify   - 验证用户
```

**响应格式标准**：
```json
// 成功响应
{
  "success": true,
  "data": {
    "userId": "uuid-xxx",
    "userName": "张三"
  }
}

// 列表响应
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 100
  }
}

// 错误响应
{
  "success": false,
  "error": "错误信息（用户可读）",
  "errorCode": "USER_NOT_FOUND",
  "details": {}
}
```

**HTTP 状态码使用**：
```
200 OK          - 查询成功
201 Created     - 创建成功
204 No Content  - 删除成功

400 Bad Request         - 请求参数错误
401 Unauthorized        - 未认证
403 Forbidden           - 无权限
404 Not Found           - 资源不存在
422 Unprocessable Entity - 参数验证失败

500 Internal Server Error - 服务器错误
```

### 常用命令

```bash
# 数据库连接
psql -h localhost -p 5432 -U nexus_user -d pr_business_db

# 检查表结构
\d users
\d+ user_accounts

# 检查索引
\di

# 退出
\q
```

### 与 PR Business 项目对比

| 规范项 | 标准要求 | PR Business 现状 | 符合度 |
|--------|---------|-----------------|--------|
| 前端框架 | Vite + React | ✅ Vite + React | ✅ 符合 |
| 后端框架 | Go + Gin | ✅ Go + Gin | ✅ 符合 |
| 数据库 | PostgreSQL | ✅ PostgreSQL | ✅ 符合 |
| 主键类型 | UUID | ✅ 使用 UUID | ✅ 符合 |
| 表命名 | snake_case | ✅ snake_case | ✅ 符合 |
| 列命名 | snake_case | ⚠️ 部分使用 camelCase | ⚠️ 部分不符合 |
| API 响应格式 | 标准格式 | ⚠️ 不统一 | ⚠️ 需改进 |

**需要改进的地方**：
1. ⚠️ 确保 GORM 模型中所有列都有明确的 `column` 标签
2. ⚠️ 统一 API 响应格式
3. ⚠️ 检查所有表的列命名是否符合 snake_case

---

## 📖 文档3: ssh-setup.md - SSH配置指南

### 核心概念

#### SSH 配置规范（强制）

**本地 SSH 配置**（`~/.ssh/config`）：
```bash
# ===== 上海应用服务器 =====
Host shanghai-tencent
    HostName 101.35.120.199
    User ubuntu                    # 操作系统用户
    IdentityFile ~/.ssh/xia_mac_shanghai_secure  # ED25519密钥
    StrictHostKeyChecking no
    ServerAliveInterval 60
    ServerAliveCountMax 3

# ===== 杭州数据库服务器 =====
Host hangzhou-ali
    HostName 47.110.82.96
    User root                      # 操作系统用户（管理员）
    IdentityFile ~/.ssh/xia_mac_hangzhou_secure  # ED25519密钥
    StrictHostKeyChecking no
    ServerAliveInterval 60
    ServerAliveCountMax 3

# ===== 香港服务器 =====
Host hongkong-tencent
    HostName 150.109.157.61
    User root
    IdentityFile ~/.ssh/xia_mac_hongkong_secure  # ED25519密钥
    StrictHostKeyChecking no
    ServerAliveInterval 60
    ServerAliveCountMax 3
```

### 关键规范

#### 1. 密钥权限（必须设置）
```bash
chmod 600 ~/.ssh/config
chmod 600 ~/.ssh/xia_mac_shanghai_secure          # 上海服务器私钥
chmod 644 ~/.ssh/xia_mac_shanghai_secure.pub      # 上海服务器公钥
chmod 600 ~/.ssh/xia_mac_hangzhou_secure          # 杭州服务器私钥
chmod 644 ~/.ssh/xia_mac_hangzhou_secure.pub      # 杭州服务器公钥
```

#### 2. 使用方式（强制）

**✅ 正确方式**：
```bash
# 连接上海服务器
ssh shanghai-tencent

# 连接杭州服务器
ssh hangzhou-ali

# 上传文件
scp bin/server shanghai-tencent:/var/www/pr-business-backend/

# 执行远程命令
ssh shanghai-tencent "sudo systemctl restart pr-business-backend"
```

**❌ 禁止的方式**：
```bash
# ❌ 不要直接用IP
ssh ubuntu@101.35.120.199

# ❌ 不要每次输入密码（应该配置密钥认证）
ssh ubuntu@101.35.120.199

# ❌ 不要用不同的别名
ssh shanghai
ssh prod
```

#### 3. 数据库连接配置（重要）

**统一配置**（所有 V3.0 系统）：
```bash
# PostgreSQL数据库连接配置
主机:   localhost (通过SSH隧道转发)
端口:   5432
数据库用户: nexus_user (PostgreSQL超级用户)
数据库密码: hRJ9NSJApfeyFDraaDgkYowY
SSL模式: disable (SSH隧道已加密，数据库层可disable)
连接字符串: postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@localhost:5432/数据库名?sslmode=disable
```

**错误示例**：
```bash
# ❌ 错误1：直连且不用SSL（不安全）
DATABASE_URL=postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@47.110.82.96:5432/db?sslmode=disable

# ❌ 错误2：使用专用数据库用户（已废弃）
DATABASE_URL=postgresql://pr_business_user:pass@localhost:5432/db?sslmode=disable
```

**正确配置**：
```bash
# ✅ 正确：通过SSH隧道 + nexus_user用户
DATABASE_URL=postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@localhost:5432/数据库名?sslmode=disable
```

### 常用命令

```bash
# 测试SSH连接
ssh shanghai-tencent "hostname && echo '连接成功'"
ssh hangzhou-ali "hostname && echo '连接成功'"

# 测试数据库连接（通过SSH隧道）
PGPASSWORD=hRJ9NSJApfeyFDraaDgkYowY psql -h localhost -p 5432 -U nexus_user -d postgres -c 'SELECT 1;'
```

### 常见错误排查

**错误1：连接超时**
```
connection timeout
```
原因：SSH隧道未启动
解决：`sudo systemctl start pg-tunnel`

**错误2：密码认证失败**
```
password authentication failed for user "nexus_user"
```
原因：配置文件中密码错误
解决：确保密码是 `hRJ9NSJApfeyFDraaDgkYowY`

**错误3：连接被拒绝**
```
connection refused
```
原因：SSH隧道未启动或端口占用
解决：检查 `sudo systemctl status pg-tunnel`

**错误4：SSL错误**
```
server does not support SSL
```
原因：使用了 `sslmode=require` 或 `prefer`
解决：使用 `sslmode=disable`

### 与 PR Business 项目对比

| 配置项 | 标准要求 | PR Business 现状 | 符合度 |
|--------|---------|-----------------|--------|
| 数据库用户 | nexus_user | ❌ nexus | ⚠️ 不符合 |
| 数据库密码 | hRJ9NSJApfeyFDraaDgkYowY | ❌ nexus123 | ⚠️ 不符合 |
| 连接方式 | SSH隧道 localhost:5432 | ⚠️ 直连 47.110.82.96:5432 | ⚠️ 不符合 |
| SSL模式 | disable | ✅ disable | ✅ 符合 |

**需要改进的地方**：
1. 🔴 **紧急**: 更新数据库用户为 `nexus_user`
2. 🔴 **紧急**: 更新数据库密码为 `hRJ9NSJApfeyFDraaDgkYowY`
3. 🔴 **紧急**: 配置 SSH 隧道服务（pg-tunnel.service）
4. 🔴 **紧急**: 修改连接方式为通过 SSH 隧道（localhost:5432）

**当前配置文件 (.env.production)**：
```bash
# ❌ 当前配置（不符合标准）
DB_HOST=47.110.82.96
DB_PORT=5432
DB_USER=nexus
DB_PASSWORD=nexus123
DB_NAME=pr_business_db
DB_SSLMODE=disable

# ✅ 应该改为
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME=pr_business_db
DB_SSLMODE=disable
```

---

## 📖 文档4: database-guide.md - 数据库使用指南

### 核心概念

#### 数据库架构规范

**统一数据库服务器**（杭州）：
```
IP: 47.110.82.96
数据库: PostgreSQL 15 (Docker)
端口: 5432
用途: 统一数据存储中心
```

**⚠️ 重要：数据库连接配置规范**

**网络架构**：
- 上海服务器**可以**直连杭州服务器的 `47.110.82.96:5432`（延迟约13ms）
- **但强烈推荐使用 SSH 隧道**（更安全、已配置好）

**两种连接方式对比**：

| 方式 | 连接字符串 | 优点 | 缺点 | 推荐度 |
|------|-----------|------|------|--------|
| **SSH隧道** | `postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@localhost:5432/db?sslmode=disable` | ✅ 加密传输<br>✅ 密钥认证<br>✅ 端口不暴露 | ❌ 多一个SSH进程 | ⭐⭐⭐⭐⭐ 强烈推荐 |
| **直连** | `postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@47.110.82.96:5432/db?sslmode=require` | ✅ 性能略好(~3%)<br>✅ 配置简单 | ❌ 需配置SSL<br>❌ 端口暴露 | ⭐⭐⭐ 可接受 |

**⚠️ 安全警告**：
```bash
# ❌ 危险：直连且不使用SSL
DATABASE_URL=postgresql://nexus_user:hRJ9NSJApfeyFDraaDgkYowY@47.110.82.96:5432/db?sslmode=disable
# 密码和数据都是明文传输！
```

### 关键规范

#### SSH隧道设置（必须执行）

**在上海服务器上执行**（首次部署时执行一次）：

```bash
# 1. 创建systemd服务
sudo tee /etc/systemd/system/pg-tunnel.service <<EOF
[Unit]
Description=PostgreSQL SSH Tunnel to Hangzhou
After=network.target

[Service]
User=ubuntu
ExecStart=/usr/bin/ssh -N -o StrictHostKeyChecking=no -o ServerAliveInterval=60 -o ServerAliveCountMax=3 -L 5432:localhost:5432 hangzhou-ali
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 2. 启动服务
sudo systemctl daemon-reload
sudo systemctl enable pg-tunnel
sudo systemctl start pg-tunnel

# 3. 验证隧道状态
sudo systemctl status pg-tunnel
# 应该看到: Active: active (running)

# 4. 测试连接
PGPASSWORD=hRJ9NSJApfeyFDraaDgkYowY psql -h localhost -p 5432 -U nexus_user -d postgres -c 'SELECT 1;'
# 应该输出: ?column?
#              ----------
#                      1
```

**SSH配置**（`~/.ssh/config`）：
```bash
Host hangzhou-ali
    HostName 47.110.82.96
    User root
    IdentityFile ~/.ssh/xia_mac_alicloud_local
    StrictHostKeyChecking no
    ServerAliveInterval 60
    ServerAliveCountMax 3
```

#### 数据库隔离策略

```
PostgreSQL Server (47.110.82.96:5432)
│
├─ auth_center_db        -- 账号中心（认证专用）
├─ pr_business_db        -- PR业务系统
├─ pixel_business_db     -- AI生图系统
├─ quote_business_db     -- 报价系统
├─ study_business_db     -- 知识库系统
└─ crm_business_db       -- 客户管理系统
```

**强制规则**：
- ✅ 每个业务系统使用**独立数据库**
- ✅ 不允许跨库查询（应用层Join）
- ✅ 通过 `auth_center_user_id` 关联用户身份

#### 用户关联规范

**auth_center_db.users**（统一身份）：
```sql
CREATE TABLE users (
  user_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  union_id      VARCHAR(255) UNIQUE,      -- 微信 unionid（跨应用）
  phone_number  VARCHAR(255) UNIQUE,
  password_hash VARCHAR(255),
  email         VARCHAR(255) UNIQUE,
  created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

**业务数据库.users**（本地用户）：
```sql
-- pr_business_db.users
CREATE TABLE users (
  id                     VARCHAR(255) PRIMARY KEY,  -- 本地 ID (CUID)
  auth_center_user_id    UUID UNIQUE,              -- 关联账号中心 ✅ 强制
  union_id               VARCHAR(255) UNIQUE,
  role                   VARCHAR(50),              -- 业务角色
  profile                JSONB,
  created_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ✅ 关键索引
CREATE INDEX users_auth_center_user_id_idx
  ON users(auth_center_user_id);
```

**跨系统用户识别**：
```sql
-- 通过 auth_center_user_id 关联
SELECT
  u.*,
  a.union_id,
  a.phone_number
FROM pr_business_db.users u
JOIN auth_center_db.users a
  ON u.auth_center_user_id = a.user_id
WHERE u.id = 'xxx';
```

### 常用命令

```bash
# === 数据库连接 ===
# 通过SSH隧道连接
PGPASSWORD=hRJ9NSJApfeyFDraaDgkYowY psql -h localhost -p 5432 -U nexus_user -d pr_business_db

# 查看所有数据库
\l

# 查看当前数据库的所有表
\dt

# 查看表结构
\d users
\d+ user_accounts

# 查看索引
\di

# 执行SQL查询
SELECT * FROM users LIMIT 10;

# 退出
\q

# === SSH隧道管理 ===
# 启动隧道
sudo systemctl start pg-tunnel

# 停止隧道
sudo systemctl stop pg-tunnel

# 重启隧道
sudo systemctl restart pg-tunnel

# 查看隧道状态
sudo systemctl status pg-tunnel

# 查看隧道日志
sudo journalctl -u pg-tunnel -f
```

### 与 PR Business 项目对比

| 配置项 | 标准要求 | PR Business 现状 | 符合度 |
|--------|---------|-----------------|--------|
| 数据库用户 | nexus_user | ❌ nexus | ⚠️ 不符合 |
| 数据库密码 | hRJ9NSJApfeyFDraaDgkYowY | ❌ nexus123 | ⚠️ 不符合 |
| 连接方式 | SSH隧道 localhost | ❌ 直连 47.110.82.96 | ⚠️ 不符合 |
| 数据库隔离 | 独立数据库 | ✅ pr_business_db | ✅ 符合 |
| 用户关联 | auth_center_user_id | ⚠️ 需确认 | ⚠️ 待检查 |

**需要改进的地方**：
1. 🔴 **紧急**: 更新数据库连接配置
2. 🔴 **紧急**: 配置并启动 SSH 隧道服务
3. ⚠️ 检查用户表是否有 `auth_center_user_id` 字段
4. ⚠️ 如有，添加索引以提升性能

---

## 📖 文档5: deployment-and-operations.md - 部署与运维

### 核心概念

#### V3.0 标准部署架构

**架构 A: Vite + React（推荐）**：
```
┌─────────────────────────────────────────────────────────────┐
│              Vite + React 系统部署架构                      │
└─────────────────────────────────────────────────────────────┘

业务系统 (example.com)
│
├── Nginx (443/80)
│   ├── SSL 终止
│   ├── 静态资源服务 (Vite 构建产物)
│   │   ├── /          → /var/www/example-frontend/index.html
│   │   └── /assets/   → 静态资源 (1年缓存)
│   └── 反向代理
│       └── /api       → Backend (Go :8080)
│
└── Backend (Go)
    ├── 端口: 8080
    ├── 运行: Systemd
    ├── 功能: RESTful API
    └── 连接: PostgreSQL (通过SSH隧道)
```

**架构 B: Next.js（仅用于 SSR 场景）**：
```
┌─────────────────────────────────────────────────────────────┐
│               Next.js 系统部署架构                          │
└─────────────────────────────────────────────────────────────┘

业务系统 (example.com)
│
├── Nginx (443/80)
│   ├── SSL 终止
│   ├── 静态资源服务
│   └── 反向代理
│       ├── /          → Frontend (Next.js :3000)
│       └── /api       → Backend (Go :8080)
│
├── Frontend (Next.js)
│   ├── 端口: 3000
│   ├── 运行: PM2
│   └── 功能: SSR + 静态页面
│
└── Backend (Go)
    ├── 端口: 8080
    ├── 运行: Systemd
    ├── 功能: RESTful API
    └── 连接: PostgreSQL (通过SSH隧道)
```

**注意**：所有新系统应使用 **架构 A (Vite + React)**。

### 关键规范

#### ⚠️ 核心原则：本地构建，上传产物

**强制要求**：
- ✅ **前端和后端都必须在本地构建**
- ✅ **只上传构建产物到服务器**
- ❌ **禁止在服务器上运行构建命令**

**为什么必须本地构建？**

1. **不占用服务器资源**
   - 编译非常消耗 CPU 和内存
   - 服务器资源宝贵，应该用于运行服务
   - 本地 Mac 性能通常比服务器强

2. **构建产物可重现**
   - 本地环境可控（依赖版本、编译器版本）
   - 避免服务器环境差异导致的问题
   - Go 交叉编译已验证完全可靠

3. **安全性更好**
   - 服务器不需要安装开发工具（Node.js、Go 编译器等）
   - 减少攻击面
   - 服务器只保留必要的运行时文件

4. **部署更快**
   - 本地构建完成后，只需要上传文件
   - 避免服务器编译耗时长
   - 减少服务中断时间

#### 1. 前端部署（Vite + React）

```bash
# === 本地开发 ===
cd frontend/

# 1. 安装依赖
npm install

# 2. 环境配置
cat > .env.production << EOF
VITE_API_URL=https://pr.crazyaigc.com/api
VITE_APP_URL=https://pr.crazyaigc.com
EOF

# 3. 开发（可选）
npm run dev

# 4. 类型检查 + 构建
npm run build  # tsc -b && vite build

# === 部署到服务器 ===

# 5. 上传构建产物（静态文件）
rsync -avz dist/ shanghai-tencent:/var/www/pr-business-frontend/

# 6. Nginx 配置（直接服务静态文件）
# sudo nginx -t && sudo systemctl reload nginx
```

**Nginx 配置（Vite 静态文件）**:
```nginx
server {
    listen 443 ssl http2;
    server_name pr.crazyaigc.com;

    # 前端静态文件（Vite 构建）
    location / {
        root /var/www/pr-business-frontend;
        try_files $uri $uri/ /index.html;

        # 静态资源缓存
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # 后端 API
    location /api {
        rewrite ^/api/?(.*) /$1 break;
        proxy_pass http://localhost:8081;

        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

#### 2. 后端部署（Go）

**重要**：
- ✅ **本地交叉编译**（Mac → Linux）
- ✅ **禁用 cgo**（`CGO_ENABLED=0`）
- ✅ **生成静态链接二进制**
- ✅ **上传二进制文件到服务器**

```bash
# === 本地开发 ===
cd backend/

# 1. 下载依赖
go mod download

# 2. 本地运行（可选）
go run main.go

# === 交叉编译（本地 Mac → Linux 服务器） ===

# 3. 编译 Linux 二进制（静态链接）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -o bin/pr-business-linux \
  main.go

# 4. 验证二进制文件
file bin/pr-business-linux
# 输出：bin/pr-business-linux: ELF 64-bit LSB executable, x86-64, ...

# 5. 上传二进制和配置
scp bin/pr-business-linux shanghai-tencent:/var/www/pr-business-backend/
scp .env.production shanghai-tencent:/var/www/pr-business-backend/.env

# 6. 重启服务（在服务器上）
ssh shanghai-tencent "sudo systemctl restart pr-business-backend"

# 7. 检查服务状态
ssh shanghai-tencent "sudo systemctl status pr-business-backend"
```

**Systemd 服务配置**：
```ini
# /etc/systemd/system/pr-business-backend.service
[Unit]
Description=PR Business Backend API
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/var/www/pr-business-backend
ExecStart=/var/www/pr-business-backend/pr-business-linux
Restart=always
RestartSec=5
Environment="PORT=8081"
EnvironmentFile=/var/www/pr-business-backend/.env

[Install]
WantedBy=multi-user.target
```

### 服务管理规范

#### 后端服务（Go）- 使用 Systemd

**启动服务**：
```bash
sudo systemctl start pr-business-backend
```

**停止服务**：
```bash
sudo systemctl stop pr-business-backend
```

**重启服务**：
```bash
sudo systemctl restart pr-business-backend
```

**查看状态**：
```bash
sudo systemctl status pr-business-backend
```

**查看日志**：
```bash
sudo journalctl -u pr-business-backend -f
```

**开机自启**：
```bash
sudo systemctl enable pr-business-backend
```

#### 前端服务（静态文件）- 使用 Nginx

**部署流程**：
```bash
# 1. 本地构建
npm run build

# 2. 上传到服务器
rsync -avz dist/ shanghai-tencent:/var/www/pr-business-frontend/

# 3. 测试 Nginx 配置
ssh shanghai-tencent "sudo nginx -t"

# 4. 重载 Nginx（无需重启）
ssh shanghai-tencent "sudo systemctl reload nginx"
```

**Nginx 管理命令**：
```bash
sudo nginx -t                    # 测试配置文件
sudo systemctl reload nginx      # 重载配置（不中断服务）
sudo systemctl restart nginx     # 重启服务
sudo systemctl status nginx      # 查看状态
```

### 常用命令

```bash
# === 前端部署 ===
npm run build                                          # 本地构建
rsync -avz dist/ shanghai-tencent:/var/www/pr-business-frontend/  # 上传
ssh shanghai-tencent "sudo nginx -t"                   # 测试配置
ssh shanghai-tencent "sudo systemctl reload nginx"     # 重载Nginx

# === 后端部署 ===
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/pr-business-linux main.go  # 本地编译
scp bin/pr-business-linux shanghai-tencent:/var/www/pr-business-backend/  # 上传
ssh shanghai-tencent "sudo systemctl restart pr-business-backend"  # 重启服务
ssh shanghai-tencent "sudo systemctl status pr-business-backend"   # 查看状态

# === 日志查看 ===
sudo journalctl -u pr-business-backend -f             # 后端日志
sudo tail -f /var/log/nginx/access.log                # Nginx访问日志
sudo tail -f /var/log/nginx/error.log                 # Nginx错误日志

# === 服务管理 ===
sudo systemctl status pr-business-backend             # 查看状态
sudo systemctl restart pr-business-backend            # 重启服务
sudo systemctl start pg-tunnel                        # 启动SSH隧道
sudo systemctl status pg-tunnel                       # 查看隧道状态
```

### 与 PR Business 项目对比

| 配置项 | 标准要求 | PR Business 现状 | 符合度 |
|--------|---------|-----------------|--------|
| 前端构建 | 本地构建 | ✅ 本地构建 | ✅ 符合 |
| 后端构建 | 本地交叉编译 | ✅ 本地编译 | ✅ 符合 |
| 服务管理 | Systemd | ✅ Systemd | ✅ 符合 |
| 前端部署 | Nginx静态文件 | ✅ Nginx | ✅ 符合 |
| 环境变量 | .env.production | ✅ 有配置文件 | ✅ 符合 |
| SSH隧道 | pg-tunnel.service | ⚠️ 需确认 | ⚠️ 待检查 |

**需要改进的地方**：
1. ⚠️ 确认 SSH 隧道服务是否已配置和运行
2. ⚠️ 确认 Nginx 配置是否符合标准
3. ⚠️ 确认服务名称是否一致

---

## 📖 文档6: api.md - API接口说明

### 核心概念

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

**设计理念**：
```
unionid = 人（同一用户在不同应用）
openid = 登录入口（同一应用不同用户）
```

### 关键规范

#### 1. 认证规范

**所有 API 请求需要在 Header 中携带 Token**:
```
Authorization: Bearer <token>
```

#### 2. 账号中心 API 接口

**发起微信登录（智能检测）**：
- **接口**: `GET /api/auth/wechat/login`
- **参数**: `callbackUrl`（URL编码）
- **响应**: 重定向到微信授权页面

**验证Token**：
- **接口**: `POST /api/auth/verify-token`
- **请求体**: `{"token": "xxx"}`
- **响应**: `{"success": true, "data": {"valid": true, "userId": "xxx"}}`

**获取用户信息**：
- **接口**: `GET /api/auth/user-info`
- **请求头**: `Authorization: Bearer <token>`
- **响应**: 用户详细信息

**密码登录**：
- **接口**: `POST /api/auth/password/login`
- **请求体**: `{"phoneNumber": "13800138000", "password": "xxx"}`
- **响应**: `{"success": true, "token": "xxx", "userId": "xxx"}`

**登出**：
- **接口**: `POST /api/auth/signout`
- **请求头**: `Authorization: Bearer <token>`

**获取会话列表**：
- **接口**: `GET /api/auth/sessions`
- **请求头**: `Authorization: Bearer <token>`
- **响应**: 所有登录设备列表

#### 3. 业务系统集成

**快速开始（5分钟完成对接）**：

**前端部分（TypeScript/React）**：
```typescript
// 引导用户跳转到业务系统的后端接口
window.location.href = 'https://pr.crazyaigc.com/api/auth/wechat/login'
```

**后端部分（Go 示例）**：

**发起微信登录**：
```go
// GET /api/auth/wechat/login
func WechatLogin(c *gin.Context) {
    callbackUrl := "https://pr.crazyaigc.com/api/auth/callback"
    authCenterURL := fmt.Sprintf(
        "https://os.crazyaigc.com/api/auth/wechat/login?callbackUrl=%s",
        url.QueryEscape(callbackUrl),
    )
    c.Redirect(302, authCenterURL)
}
```

**接收微信授权回调**：
```go
// GET /api/auth/callback?code=xxx&type=open
func AuthCallback(c *gin.Context) {
    code := c.Query("code")
    loginType := c.Query("type")

    // 调用账号中心的微信登录API
    loginResp, _ := http.Post(
        "https://os.crazyaigc.com/api/auth/wechat/login",
        "application/json",
        strings.NewReader(fmt.Sprintf(`{"code":"%s","type":"%s"}`, code, loginType)),
    )

    var result struct {
        Success bool `json:"success"`
        Data struct {
            UserID    string `json:"userId"`
            Token     string `json:"token"`
            UnionID   string `json:"unionId"`
        } `json:"data"`
    }
    json.NewDecoder(loginResp.Body).Decode(&result)

    // 创建/获取本地用户
    user := findOrCreateUser(result.Data.UserID)

    // 设置 session
    setSession(c, user, result.Data.Token)

    // 跳转到首页
    c.Redirect(302, "/dashboard")
}
```

#### 4. 业务系统必需字段

```sql
-- 业务系统用户表示例
CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY,  -- 本地主键（CUID/UUID）
  auth_center_user_id UUID UNIQUE NOT NULL,  -- ✅ 关联账号中心
  role VARCHAR(50) DEFAULT 'USER',
  profile JSONB,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 关键索引
CREATE UNIQUE INDEX users_auth_center_user_id_idx
  ON users(auth_center_user_id);
```

### 常用命令

```bash
# 测试账号中心API
curl -X POST https://os.crazyaigc.com/api/auth/verify-token \
  -H "Content-Type: application/json" \
  -d '{"token": "your-token-here"}'

# 获取用户信息
curl -X GET https://os.crazyaigc.com/api/auth/user-info \
  -H "Authorization: Bearer your-token-here"
```

### 与 PR Business 项目对比

| 配置项 | 标准要求 | PR Business 现状 | 符合度 |
|--------|---------|-----------------|--------|
| 认证集成 | 账号中心 | ✅ 已集成 | ✅ 符合 |
| 用户关联 | auth_center_user_id | ⚠️ 需确认 | ⚠️ 待检查 |
| 会话管理 | JWT Token | ✅ JWT | ✅ 符合 |
| Token有效期 | 7天 | ⚠️ 24h访问token | ⚠️ 不完全符合 |

**需要改进的地方**：
1. ⚠️ 检查用户表是否有 `auth_center_user_id` 字段
2. ⚠️ 检查 Token 有效期设置（标准是 7 天）
3. ⚠️ 确认是否实现了账号中心的回调接口

---

## 📖 文档7: security.md - 安全规范

### 核心概念

#### 2026-02-03 密钥泄露事件教训

**事件概述**：
- 攻击者通过泄露的 SSH 私钥登录上海服务器
- 尝试访问数据库
- 暴露了密钥管理和权限控制的严重问题

**根本原因**：
1. 将本地登录密钥上传到服务器（用于 SSH 隧道）
2. 密钥在多个服务器间共享，缺乏隔离
3. 没有定期检查服务器访问日志
4. 缺少高风险操作的确认机制

### 关键规范

#### 1. 密钥管理规范（⚠️ 极其重要）

##### 密钥隔离原则（强制）

**每个用途使用独立的密钥对**：

| 用途 | 密钥命名 | 存放位置 | 权限范围 |
|------|----------|----------|----------|
| 开发者登录 | `xia_mac_<server>_secure` | 本地 Mac | 该服务器 SSH 登录 |
| 服务器间隧道 | `<src>_<dest>_tunnel` | 源服务器仅一份 | 仅用于隧道，不能登录其他服务器 |
| CI/CD 部署 | `deploy_<service>` | CI 服务器 | 仅部署权限 |
| 数据库访问 | **禁止** | **禁止** | 使用应用层连接 |

**❌ 严禁操作**：
- 将本地登录密钥上传到服务器
- 在多个服务器间共享同一个私钥
- 将私钥放在代码仓库、配置文件中
- 使用私钥进行数据库连接（应用层应使用用户名密码）

**✅ 正确做法**：
- 每个服务器/用途使用独立的 ED25519 密钥对
- 隧道密钥只留在源服务器，不用于其他目的
- 私钥权限设置为 `400` 或 `600`
- 定期（每季度）审计所有活跃密钥

##### 密钥生成标准（强制）

**所有新密钥必须使用 ED25519 算法**：
```bash
# 生成 ED25519 密钥
ssh-keygen -t ed25519 -C "<用途>-$(date +%Y%m%d)" -f ~/.ssh/<密钥名> -N ""
```

**密钥命名规范**：
- 本地登录：`<服务器>_secure` (例: `shanghai_secure`)
- 隧道专用：`<源>_<目标>_tunnel` (例: `shanghai_hangzhou_tunnel`)
- CI/CD：`deploy_<service>` (例: `deploy_auth_center`)

##### 密钥轮换周期（强制）

| 密钥类型 | 轮换周期 | 触发条件 |
|---------|----------|----------|
| 登录密钥 | 6 个月 | 怀疑泄露、人员离职 |
| 隧道密钥 | 12 个月 | 服务器迁移、架构变更 |
| 数据库密码 | 3 个月 | 怀疑泄露、安全事件 |
| JWT Secret | 3 个月 | 系统部署、密钥泄露 |

#### 2. 高风险操作规范（⚠️ 需要确认）

以下操作**必须获得人工确认**后才能执行：

##### 密钥相关操作（🔴 极高风险）

| 操作 | 风险等级 | 确认要求 | 审批流程 |
|------|----------|----------|----------|
| 生成新密钥对 | 🟡 中 | 记录用途和指纹 | 开发者自行确认 |
| 上传私钥到服务器 | 🔴 高 | **必须经所有者批准** | **必须所有者确认** |
| 添加公钥到服务器 | 🟡 中 | 记录来源和用途 | 文档记录 |
| 删除服务器上的公钥 | 🔴 高 | **确认无影响** | **必须所有者确认** |
| 更换服务器密钥 | 🔴 高 | **全面测试连接** | **必须所有者确认** |
| 更换数据库密码 | 🔴 高 | **所有服务重启** | **必须所有者确认** |

##### 确认流程（强制）

**高风险操作执行前必须**：
1. **文档记录**：在项目文档中记录操作目的、影响范围、回滚方案
2. **所有者确认**：获得服务器/系统所有者的明确同意（文字记录）
3. **备份**：备份所有将被修改的文件
4. **测试**：在测试环境验证操作步骤
5. **时间窗口**：选择低峰期执行，减少用户影响

**示例确认记录**：
```
操作：更换上海服务器 SSH 密钥
原因：2026-02-03 密钥泄露事件
影响范围：所有连接上海服务器的开发者
执行时间：2026-02-04 00:00 CST
备份位置：/home/ubuntu/.ssh/authorized_keys.backup.20260204
测试结果：✅ 已在测试环境验证
所有者确认：[所有者姓名/签名] - 2026-02-03
执行结果：✅ 成功
回滚方案：恢复旧公钥文件
```

#### 3. 认证与授权

##### JWT Token 规范
```go
// ✅ Token 结构
type Claims struct {
    UserID string `json:"userId"`
    jwt.RegisteredClaims
}

// ✅ 标准配置
- 算法: HS256
- 有效期: 7天 (168小时)
- 签名密钥: 最少32字符
- 存储: PostgreSQL sessions 表
```

##### 密码安全
```go
// ✅ 密码哈希（强制 bcrypt）
import "golang.org/x/crypto/bcrypt"

// 哈希密码
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// 验证密码
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

#### 4. CORS 配置（带白名单验证）

**⚠️ 安全要求**：账号中心必须实现 CORS 白名单，防止钓鱼网站非法调用。

```go
// ✅ CORS 中间件（带白名单验证）
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
    // 解析白名单
    allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")
    originMap := make(map[string]bool)
    for _, origin := range allowedOrigins {
        originMap[strings.TrimSpace(origin)] = true
    }

    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // 验证 Origin 是否在白名单中
        if origin != "" {
            if !originMap[origin] {
                c.JSON(403, gin.H{
                    "success": false,
                    "error":   "域名未在白名单中",
                })
                c.Abort()
                return
            }
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
        }

        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

**环境变量配置**:
```bash
# .env
# CORS 白名单（逗号分隔，只允许这些域名调用账号中心 API）
ALLOWED_ORIGINS=https://os.crazyaigc.com,https://pr.crazyaigc.com,https://pixel.crazyaigc.com,https://quote.crazyaigc.com

# 回调域名白名单（逗号分隔，只允许这些域名接收登录回调）
ALLOWED_CALLBACK_DOMAINS=pr.crazyaigc.com,www.crazyaigc.com,os.crazyaigc.com,pixel.crazyaigc.com,quote.crazyaigc.com,localhost
```

#### 5. 定期安全检查（强制执行）

##### 每周检查（周一 10:00）

**检查内容**：
```bash
# 1. 检查服务器登录日志
sudo last -n 50
sudo journalctl -u ssh --since "7 days ago" | grep -E "Failed|Invalid"

# 2. 检查活跃 SSH 连接
ss -tunp | grep ssh

# 3. 检查进程异常
ps aux --sort=-%mem | head -20
ps aux --sort=-%cpu | head -20

# 4. 检查磁盘空间
df -h

# 5. 检查服务状态
systemctl list-units --state=failed
```

##### 每月检查（每月第一个周一）

**检查内容**：
```bash
# 1. 审计所有活跃 SSH 密钥
for server in shanghai hangzhou hongkong; do
  ssh $server "cat ~/.ssh/authorized_keys && ssh-keygen -lf ~/.ssh/authorized_keys"
done

# 2. 检查密钥有效期
for key in ~/.ssh/*_secure; do
  stat -f "%Sm" -t "%Y-%m-%d" $key
done

# 3. 检查服务器访问日志
grep -E "220.205.242.226|204.76.203.83" /var/log/auth.log*

# 4. 检查数据库连接
psql -U nexus_user -d postgres -c "\l"

# 5. 检查系统更新
sudo apt list --upgradable
```

##### 每季度检查（1月/4月/7月/10月）

**检查内容**：
- 全面安全审计（包括渗透测试）
- 密钥轮换计划执行
- 安全策略更新
- 灾难恢复演练
- 访问控制审查

#### 6. fail2ban 自动防护（已部署）

**2026-02-04 部署状态**：

| 服务器 | fail2ban | 状态 | 配置 |
|--------|----------|------|------|
| 上海 | ✅ 已安装 | ✅ 运行中 | 5次失败/1小时封禁 |
| 杭州 | ✅ 已安装 | ✅ 运行中 | 5次失败/1小时封禁 |
| 香港 | ✅ 已安装 | ✅ 运行中 | 5次失败/1小时封禁 |

**防护原理**：
```
攻击者尝试密码登录
  ↓
第1次失败
  ↓
...
  ↓
第5次失败（10分钟内）
  ↓
fail2ban: 自动封禁IP 1小时 🔒
```

**配置文件：`/etc/fail2ban/jail.local`**
```ini
[sshd]
enabled = true
port = 22
filter = sshd
logpath = /var/log/auth.log
maxretry = 5        # 5次失败后封禁
findtime = 600      # 10分钟内
bantime = 3600      # 封禁1小时
# 白名单IP（不会被封禁）
ignoreip = 127.0.0.1 ::1 39.185.202.1 39.185.203.92
```

**查看封禁状态**：
```bash
# 查看当前被封禁的IP
sudo fail2ban-client status sshd

# 手动解封（如果误封）
sudo fail2ban-client set sshd unbanip <IP地址>

# 查看fail2ban日志
sudo tail -f /var/log/fail2ban.log
```

### 常用命令

```bash
# === 安全检查 ===
# 查看登录日志
sudo last -n 50
sudo journalctl -u ssh --since "7 days ago" | grep -E "Failed|Invalid"

# 查看活跃SSH连接
ss -tunp | grep ssh

# 检查进程异常
ps aux --sort=-%mem | head -20
ps aux --sort=-%cpu | head -20

# 检查磁盘空间
df -h

# 检查服务状态
systemctl list-units --state=failed

# === fail2ban 管理 ===
sudo fail2ban-client status sshd                    # 查看被封禁的IP
sudo fail2ban-client set sshd unbanip <IP>         # 手动解封
sudo tail -f /var/log/fail2ban.log                 # 查看日志
```

### 与 PR Business 项目对比

| 配置项 | 标准要求 | PR Business 现状 | 符合度 |
|--------|---------|-----------------|--------|
| 密码加密 | bcrypt | ⚠️ 需确认 | ⚠️ 待检查 |
| JWT有效期 | 7天 | ⚠️ 24h | ⚠️ 不完全符合 |
| CORS白名单 | 必须实现 | ⚠️ 需确认 | ⚠️ 待检查 |
| 数据库连接 | SSH隧道 | ❌ 直连 | ⚠️ 不符合 |
| 密钥隔离 | 独立密钥 | ⚠️ 需确认 | ⚠️ 待检查 |

**需要改进的地方**：
1. 🔴 **紧急**: 更新数据库连接方式为 SSH 隧道
2. 🔴 **紧急**: 更新数据库用户和密码
3. ⚠️ 检查密码加密方式是否使用 bcrypt
4. ⚠️ 检查 JWT Token 有效期设置
5. ⚠️ 实现 CORS 白名单验证

---

## 📊 总体对比分析

### 符合度总结

| 分类 | 符合度 | 关键问题 |
|------|--------|----------|
| **架构设计** | ✅ 90% | 基本符合 V3.0 标准 |
| **数据库配置** | ⚠️ 40% | 🔴 数据库用户和密码不符合标准<br>🔴 连接方式未使用SSH隧道 |
| **命名规范** | ⚠️ 70% | 部分列命名需要检查 |
| **API设计** | ⚠️ 60% | 响应格式需要统一 |
| **安全规范** | ⚠️ 50% | 🔴 密钥管理需改进<br>⚠️ CORS白名单待实现 |
| **部署规范** | ✅ 85% | 基本符合标准 |

### 🔴 紧急需要修复的问题

#### 1. 数据库连接配置（最高优先级）

**当前配置** (.env.production)：
```bash
# ❌ 错误配置
DB_HOST=47.110.82.96
DB_PORT=5432
DB_USER=nexus
DB_PASSWORD=nexus123
DB_NAME=pr_business_db
DB_SSLMODE=disable
```

**应该改为**：
```bash
# ✅ 正确配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME=pr_business_db
DB_SSLMODE=disable
```

#### 2. SSH隧道配置（最高优先级）

**在上海服务器上执行**：
```bash
# 1. 创建systemd服务
sudo tee /etc/systemd/system/pg-tunnel.service <<EOF
[Unit]
Description=PostgreSQL SSH Tunnel to Hangzhou
After=network.target

[Service]
User=ubuntu
ExecStart=/usr/bin/ssh -N -o StrictHostKeyChecking=no -o ServerAliveInterval=60 -o ServerAliveCountMax=3 -L 5432:localhost:5432 hangzhou-ali
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 2. 启动服务
sudo systemctl daemon-reload
sudo systemctl enable pg-tunnel
sudo systemctl start pg-tunnel

# 3. 验证隧道状态
sudo systemctl status pg-tunnel
```

#### 3. 部署流程

```bash
# === 本地操作 ===
cd /Users/xia/Documents/GitHub/pr-business

# 1. 更新 .env.production 配置
# 编辑 backend/.env.production，修改数据库配置

# 2. 重新编译后端
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/pr-business-linux main.go

# 3. 上传到服务器
scp bin/pr-business-linux shanghai-tencent:/var/www/pr-business-backend/
scp .env.production shanghai-tencent:/var/www/pr-business-backend/.env

# === 服务器操作 ===
ssh shanghai-tencent

# 4. 配置SSH隧道（首次）
sudo tee /etc/systemd/system/pg-tunnel.service <<EOF
[Unit]
Description=PostgreSQL SSH Tunnel to Hangzhou
After=network.target

[Service]
User=ubuntu
ExecStart=/usr/bin/ssh -N -o StrictHostKeyChecking=no -o ServerAliveInterval=60 -o ServerAliveCountMax=3 -L 5432:localhost:5432 hangzhou-ali
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable pg-tunnel
sudo systemctl start pg-tunnel
sudo systemctl status pg-tunnel

# 5. 重启后端服务
sudo systemctl restart pr-business-backend
sudo systemctl status pr-business-backend

# 6. 查看日志
sudo journalctl -u pr-business-backend -f
```

### ⚠️ 需要检查和改进的地方

#### 1. 代码规范检查

**需要检查的文件**：
- `/Users/xia/Documents/GitHub/pr-business/backend/models/user.go`
- `/Users/xia/Documents/GitHub/pr-business/backend/models/*.go`

**检查要点**：
```go
// ✅ 确保所有GORM模型都有正确的column标签
type User struct {
    ID       string    `gorm:"primaryKey;column:id;type:varchar(255)" json:"id"`
    AuthCenterUserID string `gorm:"uniqueIndex;column:auth_center_user_id;type:uuid" json:"authCenterUserId"`
    // ... 其他字段
}
```

#### 2. API响应格式统一

**标准格式**：
```json
{
  "success": true,
  "data": { ... },
  "pagination": { ... }
}
```

需要检查所有API响应是否符合此格式。

#### 3. 安全配置检查

**需要检查的项目**：
- [ ] 密码加密是否使用 bcrypt
- [ ] JWT Token 有效期是否为 7 天
- [ ] 是否实现了 CORS 白名单验证
- [ ] 是否实现了 JWT 认证中间件
- [ ] 敏感配置是否不在代码仓库中

#### 4. 用户关联检查

**需要确认**：
```sql
-- 检查用户表是否有 auth_center_user_id 字段
\d users

-- 如果有，检查是否有索引
\di users_auth_center_user_id_idx
```

---

## 🎯 改进计划

### 阶段1：紧急修复（立即执行）

**优先级：🔴 最高**

1. ✅ 更新 `.env.production` 数据库配置
   - 修改 DB_USER=nexus_user
   - 修改 DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
   - 修改 DB_HOST=localhost

2. ✅ 配置并启动 SSH 隧道服务
   - 创建 pg-tunnel.service
   - 启动并验证隧道状态

3. ✅ 重新编译并部署后端
   - 本地交叉编译
   - 上传到服务器
   - 重启服务

4. ✅ 验证数据库连接
   - 测试连接
   - 查看日志确认无错误

### 阶段2：规范检查（本周完成）

**优先级：🟡 高**

1. ⚠️ 检查所有 GORM 模型的列命名
   - 确保所有列都有明确的 `column` 标签
   - 确保使用 snake_case

2. ⚠️ 检查 API 响应格式
   - 统一响应格式
   - 实现标准的错误处理

3. ⚠️ 检查用户表结构
   - 确认有 `auth_center_user_id` 字段
   - 确认有索引

4. ⚠️ 检查安全配置
   - 密码加密方式
   - JWT 有效期
   - CORS 配置

### 阶段3：优化完善（下周完成）

**优先级：🟢 中**

1. 🔧 实现 CORS 白名单验证
2. 🔧 完善日志记录
3. 🔧 添加健康检查端点
4. 🔧 编写API文档
5. 🔧 添加单元测试

### 阶段4：持续改进（长期）

**优先级：🔵 低**

1. 📚 定期安全检查（每周）
2. 📚 密钥轮换（每季度）
3. 📚 依赖更新（每月）
4. 📚 性能优化（按需）

---

## 📝 学习总结

### 核心收获

1. **理解了 KeenChase 技术规范的重要性**
   - 这是"宪法"级别的标准，必须严格遵守
   - 所有系统必须遵循统一的技术栈和命名规范

2. **掌握了 V3.0 架构标准**
   - 前后端分离架构
   - 本地构建，上传产物的部署方式
   - 统一的数据层设计

3. **理解了数据库连接的正确方式**
   - 必须通过 SSH 隧道连接
   - 统一使用 nexus_user 用户
   - 独立数据库隔离策略

4. **认识到了安全规范的重要性**
   - 密钥管理的严格规范
   - 定期安全检查的必要性
   - fail2ban 自动防护的价值

5. **学会了服务管理规范**
   - Systemd 管理 Go 后端
   - Nginx 管理前端静态文件
   - PM2 管理 Next.js SSR（如需要）

### 关键注意事项

⚠️ **最需要记住的几点**：

1. **永远区分操作系统用户和数据库用户**
   - 操作系统用户：ubuntu（上海）、root（杭州）
   - 数据库用户：nexus_user（所有系统统一）

2. **数据库连接必须通过 SSH 隧道**
   - 不要直连 47.110.82.96:5432
   - 使用 localhost:5432 通过隧道

3. **本地构建，上传产物**
   - 不要在服务器上编译
   - 本地 Mac 性能更好

4. **命名规范必须严格遵守**
   - 数据库：snake_case
   - Go代码：PascalCase（结构体）、camelCase（JSON）
   - 必须使用 GORM 的 `column` 标签

5. **安全第一**
   - 密钥隔离
   - 定期检查
   - fail2ban 防护

### 下一步行动

**立即执行**（今天）：
1. ✅ 更新 `.env.production` 配置
2. ✅ 配置 SSH 隧道服务
3. ✅ 重新编译并部署

**本周完成**：
1. ⚠️ 代码规范检查
2. ⚠️ 安全配置检查
3. ⚠️ API 响应格式统一

**下周完成**：
1. 🔧 实现高级安全功能
2. 🔧 完善文档
3. 🔧 添加测试

---

## 📞 技术支持

**遇到问题时**：

1. 先查看对应文档的"常见问题"章节
2. 检查服务日志：`sudo journalctl -u <service-name> -f`
3. 查看系统日志：`sudo tail -f /var/log/syslog`
4. 联系团队支持

---

**学习完成时间**: 2026-02-04
**下次审查时间**: 2026-03-04（建议每月复习一次）

---

**本文档保存位置**: `/Users/xia/Documents/GitHub/pr-business/KEENCHASE_STANDARDS_LEARNING.md`
