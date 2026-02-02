# PR Business - 技术架构文档

**版本**: v1.0
**日期**: 2026-02-01
**状态**: 设计阶段

---

## 🏗️ 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                    用户层 (浏览器)                        │
└─────────────────────────────────────────────────────────┘
                            │
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Nginx (443/80) - 上海服务器                  │
│  ├── SSL 终止                                           │
│  ├── 静态文件服务 (前端)                                 │
│  └── 反向代理 /api → Go后端 (8081)                      │
└─────────────────────────────────────────────────────────┘
                            │
                            ↓
┌─────────────────────────────────────────────────────────┐
│            Go Backend (端口 8081) - systemd              │
│  ├── Gin Web 框架                                        │
│  ├── JWT 认证                                            │
│  ├── GORM ORM                                            │
│  └── 业务逻辑                                             │
└─────────────────────────────────────────────────────────┘
                            │
                            ↓
┌─────────────────────────────────────────────────────────┐
│       PostgreSQL (端口 5433) - 杭州服务器               │
│  └── pr_business_db 数据库                               │
└─────────────────────────────────────────────────────────┘
```

---

## 📁 项目结构

```
/Users/xia/Documents/GitHub/pr-business/
├── backend/                    # Go 后端
│   ├── cmd/server/
│   │   └── main.go            # 入口文件
│   ├── internal/
│   │   ├── config/            # 配置
│   │   ├── handler/           # HTTP 处理器
│   │   ├── service/           # 业务逻辑
│   │   ├── repository/        # 数据访问
│   │   ├── models/            # 数据模型
│   │   └── middleware/        # 中间件
│   ├── go.mod
│   └── .env                   # 环境变量
│
├── frontend/                   # React 前端
│   ├── src/
│   │   ├── pages/             # 页面组件
│   │   ├── components/        # 通用组件
│   │   ├── lib/               # 工具函数
│   │   ├── contexts/          # React Context
│   │   └── router/            # 路由配置
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
│
├── prisma/                     # 数据库模型
│   └── schema.prisma          # Prisma Schema
│
└── docs/                       # 文档
    ├── PRD.md                 # 产品需求文档
    └── ARCHITECTURE.md        # 本文档
```

---

## 💾 数据库设计

### 表关系图

```
User (用户)
  ├─ 1:1 → CreditAccount (积分账户)
  ├─ 1:1 → Creator (达人)
  ├─ 1:N → MerchantStaff (商家员工)
  └─ 1:N → ServiceProviderStaff (服务商员工)

Merchant (商家)
  ├─ N:1 → User (管理员)
  ├─ 1:N → MerchantStaff (员工)
  └─ 1:N → Task (任务)

ServiceProvider (服务商)
  ├─ N:1 → User (管理员)
  ├─ 1:N → ServiceProviderStaff (员工)
  ├─ 1:N → ProviderMerchantBinding (商家绑定)
  ├─ 1:N → ProviderCreatorBinding (达人绑定)
  └─ 1:N → Task (任务)

Creator (达人)
  ├─ N:1 → User (基础用户)
  ├─ N:1 → Creator (邀请人)
  ├─ 1:N → Creator (邀请的达人)
  ├─ 1:N → ProviderCreatorBinding (服务商绑定)
  └─ 1:N → TaskAssignment (任务分配)

Task (任务)
  ├─ N:1 → Merchant (商家)
  ├─ N:1 → ServiceProvider (服务商)
  └─ 1:N → TaskAssignment (任务分配)

TaskAssignment (任务分配)
  ├─ N:1 → Task (任务)
  ├─ N:1 → Creator (达人)
  └─ N:1 → Creator (组长)
```

### 完整 Prisma Schema

```prisma
// ==================== 环境配置 ====================
generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}

// ==================== 用户系统 ====================

model User {
  id                    String    @id @default(cuid())
  authCenterUserId      String?   @unique
  phoneNumber           String?   @unique
  role                  UserRole  @default(USER)

  // 关系
  creditAccount         CreditAccount?
  creator               Creator?
  merchantStaff         MerchantStaff[]
  serviceProviderStaff  ServiceProviderStaff[]

  createdAt             DateTime  @default(now())
  updatedAt             DateTime  @updatedAt
}

enum UserRole {
  USER
  SUPER_ADMIN
  MERCHANT_ADMIN
  MERCHANT_STAFF
  SERVICE_PROVIDER_ADMIN
  SERVICE_PROVIDER_STAFF
  CREATOR_LEADER
  CREATOR
}

// ==================== 积分系统 ====================

model CreditAccount {
  id              String   @id @default(cuid())
  userId          String   @unique
  goldCoins       BigInt   @default(0)
  diamondCredits  BigInt   @default(0)
  frozenGoldCoins BigInt   @default(0)

  user            User     @relation(fields: [userId], references: [id])

  createdAt       DateTime @default(now())
  updatedAt       DateTime @updatedAt
}

model GoldCoinCreditHistory {
  id              String   @id @default(cuid())
  userId          String
  amount          BigInt
  balance         BigInt
  type            String
  description     String?
  relatedUserId   String?
  relatedTaskId   String?

  createdAt       DateTime @default(now())
}

model DiamondCreditHistory {
  id              String   @id @default(cuid())
  userId          String
  amount          BigInt
  balance         BigInt
  type            String
  description     String?
  relatedUserId   String?
  relatedTaskId   String?

  createdAt       DateTime @default(now())
}

// ==================== 商家系统 ====================

model Merchant {
  id               String         @id @default(cuid())
  name             String
  logoUrl          String?
  industry         String?
  contactInfo      String?
  businessLicense  String?
  address          String?
  status           MerchantStatus @default(ACTIVE)
  adminId          String?

  // 关系
  admin            User?           @relation("MerchantAdmin")
  staff            MerchantStaff[]
  tasks            Task[]

  createdAt        DateTime        @default(now())
  updatedAt        DateTime        @updatedAt
}

enum MerchantStatus {
  ACTIVE
  SUSPENDED
  DELETED
}

model MerchantStaff {
  id           String       @id @default(cuid())
  merchantId   String
  userId       String
  position     String?
  status       StaffStatus  @default(ACTIVE)

  // 关系
  merchant     Merchant     @relation(fields: [merchantId], references: [id])
  user         User         @relation(fields: [userId], references: [id])

  createdAt    DateTime     @default(now())
  updatedAt    DateTime     @updatedAt
}

enum StaffStatus {
  ACTIVE
  SUSPENDED
  REMOVED
}

// ==================== 服务商系统 ====================

model ServiceProvider {
  id               String                @id @default(cuid())
  name             String
  logoUrl          String?
  businessLicense  String?
  address          String?
  status           ServiceProviderStatus @default(ACTIVE)
  adminId          String?

  // 关系
  admin            User?                 @relation("ProviderAdmin")
  staff            ServiceProviderStaff[]
  merchantBindings ProviderMerchantBinding[]
  creatorBindings  ProviderCreatorBinding[]
  tasks            Task[]

  createdAt        DateTime             @default(now())
  updatedAt        DateTime             @updatedAt
}

enum ServiceProviderStatus {
  ACTIVE
  SUSPENDED
  DELETED
}

model ProviderMerchantBinding {
  id             String         @id @default(cuid())
  providerId     String
  merchantId     String
  status         BindingStatus  @default(ACTIVE)
  contactPerson  String?
  contactPhone   String?

  // 关系
  provider       ServiceProvider @relation(fields: [providerId], references: [id])
  merchant       Merchant        @relation(fields: [merchantId], references: [id])

  createdAt      DateTime        @default(now())
  updatedAt      DateTime        @updatedAt
}

model ProviderCreatorBinding {
  id         String         @id @default(cuid())
  providerId String
  creatorId  String
  status     BindingStatus  @default(ACTIVE)
  isDirect   Boolean        @default(false)

  // 关系
  provider   ServiceProvider @relation(fields: [providerId], references: [id])
  creator    Creator         @relation(fields: [creatorId], references: [id])

  createdAt  DateTime        @default(now())
  updatedAt  DateTime        @updatedAt
}

enum BindingStatus {
  ACTIVE
  INACTIVE
}

model ServiceProviderStaff {
  id           String       @id @default(cuid())
  providerId   String
  userId       String
  position     String?
  status       StaffStatus  @default(ACTIVE)

  // 关系
  provider     ServiceProvider @relation(fields: [providerId], references: [id])
  user         User            @relation(fields: [userId], references: [id])

  createdAt    DateTime        @default(now())
  updatedAt    DateTime        @updatedAt
}

// ==================== 达人系统 ====================

model Creator {
  id         String         @id @default(cuid())
  userId     String         @unique
  level      CreatorLevel   @default(UGC)
  platformIds String?       // Prisma Json, 但用 String 存储
  tags       String?        // Prisma Json, 但用 String 存储
  status     CreatorStatus  @default(NORMAL)
  inviterId  String?

  // 关系
  user       User           @relation(fields: [userId], references: [id])
  inviter    Creator?       @relation("CreatorInviter", fields: [inviterId], references: [id])
  invitees   Creator[]      @relation("CreatorInviter")
  bindings   ProviderCreatorBinding[]
  assignments TaskAssignment[]

  createdAt  DateTime       @default(now())
  updatedAt  DateTime       @updatedAt
}

enum CreatorLevel {
  UGC
  KOC
  INF
  KOL
  LEADER
}

enum CreatorStatus {
  NORMAL
  SUSPENDED
}

// ==================== 任务系统 ====================

model Task {
  id                   String       @id @default(cuid())
  merchantId           String
  providerId           String?

  // 任务信息
  title                String
  requirements         String       @db.Text
  materials            String?      @db.Text

  // 财务信息（单位：分）
  taskCommission       BigInt       @default(0)
  creatorEarnings       BigInt       @default(0)
  leaderEarnings        BigInt       @default(0)
  providerEarnings      BigInt       @default(0)

  // 任务配置
  quota                Int

  // 时间
  deadline             DateTime?
  publishDeadline      DateTime?

  // 状态
  status               TaskStatus   @default(DRAFT)

  // 关系
  merchant             Merchant     @relation(fields: [merchantId], references: [id])
  provider             ServiceProvider? @relation(fields: [providerId], references: [id])
  assignments          TaskAssignment[]

  createdAt            DateTime     @default(now())
  updatedAt            DateTime     @updatedAt
}

enum TaskStatus {
  DRAFT
  OPEN
  IN_PROGRESS
  COMPLETED
  CANCELLED
}

model TaskAssignment {
  id             String             @id @default(cuid())
  taskId         String
  creatorId      String
  leaderId       String?

  // 提交信息
  submissionUrl  String?
  screenshotUrl  String?

  // 状态
  status         AssignmentStatus  @default(SUBMITTED)

  // 审核信息
  auditTime      DateTime?
  feedback       String?            @db.Text

  // 关系
  task           Task               @relation(fields: [taskId], references: [id])
  creator        Creator            @relation(fields: [creatorId], references: [id])
  leader         Creator?           @relation("AssignmentLeader", fields: [leaderId], references: [id])

  createdAt      DateTime           @default(now())
  updatedAt      DateTime           @updatedAt
}

enum AssignmentStatus {
  SUBMITTED
  APPROVED
  REJECTED
}
```

---

## 🔐 安全设计

### 1. 认证 (Authentication)

**JWT Token 结构**:
```json
{
  "userId": "xxx",
  "role": "CREATOR",
  "exp": 1234567890,
  "iat": 1234567890
}
```

**Token 生成**:
```go
func GenerateToken(userID, role string) (string, error) {
  claims := jwt.MapClaims{
    "userId": userID,
    "role":   role,
    "exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
    "iat":    time.Now().Unix(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString([]byte(JWT_SECRET))
}
```

### 2. 授权 (Authorization)

**权限中间件**:
```go
func RequireRole(roles ...string) gin.HandlerFunc {
  return func(c *gin.Context) {
    userRole := c.GetString("role")
    for _, role := range roles {
      if userRole == role {
        c.Next()
        return
      }
    }
    c.JSON(403, gin.H{"error": "权限不足"})
    c.Abort()
  }
}
```

**使用示例**:
```go
// 只有超级管理员可以访问
router.GET("/admin/users", middleware.RequireRole("SUPER_ADMIN"), handler.GetUsers)

// 商家管理员和服务商管理员都可以访问
router.GET("/tasks", middleware.RequireRole("MERCHANT_ADMIN", "SERVICE_PROVIDER_ADMIN"), handler.GetTasks)
```

### 3. 并发控制

**抢单逻辑（悲观锁）**:
```go
func AcceptTask(db *gorm.DB, taskID, creatorID string) error {
  return db.Transaction(func(tx *gorm.DB) error {
    // 1. 查询任务并加锁（FOR UPDATE）
    var task Task
    if err := tx.Where("id = ? AND status = ?", taskID, "OPEN").
      Clauses(clause.Locking{Strength: "UPDATE"}).
      First(&task).Error; err != nil {
      return err
    }

    // 2. 检查名额
    if task.Quota <= 0 {
      return errors.New("名额已满")
    }

    // 3. 创建任务分配
    assignment := TaskAssignment{
      TaskID:    taskID,
      CreatorID: creatorID,
      Status:    "SUBMITTED",
    }
    if err := tx.Create(&assignment).Error; err != nil {
      return err
    }

    // 4. 减少名额
    if err := tx.Model(&task).
      Update("quota", gorm.Expr("quota - ?", 1)).Error; err != nil {
      return err
    }

    // 5. 如果名额满，更新任务状态
    if task.Quota-1 == 0 {
      tx.Model(&task).Update("status", "COMPLETED")
    }

    return nil
  })
}
```

### 4. 积分安全

**原子操作（防止并发问题）**:
```go
// 错误示例（不安全）
account.GoldCoins += amount
db.Save(&account)

// 正确示例（使用原子更新）
db.Model(&CreditAccount{}).
  Where("user_id = ?", userID).
  Update("gold_coins", gorm.Expr("gold_coins + ?", amount))
```

**条件更新（防止重复结算）**:
```go
result := db.Model(&TaskAssignment{}).
  Where("id = ? AND status = ?", assignmentID, "SUBMITTED").
  Updates(map[string]interface{}{
    "status":    "APPROVED",
    "auditTime": time.Now(),
  })

if result.RowsAffected == 0 {
  return errors.New("任务已审核或不存在")
}
```

---

## 📡 API 规范

### RESTful 设计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/tasks | 获取任务列表 |
| POST | /api/tasks | 创建任务 |
| GET | /api/tasks/:id | 获取任务详情 |
| PUT | /api/tasks/:id | 更新任务 |
| DELETE | /api/tasks/:id | 删除任务 |
| POST | /api/tasks/:id/accept | 接单 |
| PUT | /api/tasks/:id/submit | 提交任务 |

### 响应格式

**成功响应**:
```json
{
  "success": true,
  "data": {
    "taskId": "xxx",
    "title": "任务标题"
  }
}
```

**错误响应**:
```json
{
  "success": false,
  "error": "错误描述",
  "errorCode": "TASK_NOT_FOUND"
}
```

**列表响应（带分页）**:
```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

### HTTP 状态码

| 状态码 | 场景 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如重复创建） |
| 500 | 服务器错误 |

---

## 🚀 部署配置

### 环境变量

**后端 (.env)**:
```bash
# 服务器
PORT=8081
GIN_MODE=release

# 数据库
DATABASE_URL=postgresql://nexus:nexus123@localhost:5433/pr_business_db?sslmode=disable

# JWT
JWT_SECRET=your-secret-key-min-32-chars-long

# 微信
WECHAT_OPEN_APP_ID=xxx
WECHAT_OPEN_APP_SECRET=xxx
WECHAT_MP_APPID=xxx
WECHAT_MP_SECRET=xxx

# CORS
ALLOWED_ORIGINS=https://pr.crazyaigc.com
```

**前端 (.env.production)**:
```bash
VITE_API_URL=https://pr.crazyaigc.com
VITE_APP_URL=https://pr.crazyaigc.com
```

### Nginx 配置

```nginx
server {
  listen 443 ssl;
  server_name pr.crazyaigc.com;

  # SSL 证书
  ssl_certificate /etc/letsencrypt/live/pr.crazyaigc.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/pr.crazyaigc.com/privkey.pem;

  # 前端静态文件
  location / {
    root /var/www/pr-business-frontend;
    try_files $uri $uri/ /index.html;

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
      expires 1y;
      add_header Cache-Control "public, immutable";
    }
  }

  # API 代理
  location /api {
    rewrite ^/api/?(.*) /$1 break;
    proxy_pass http://127.0.0.1:8081;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # CORS
    add_header Access-Control-Allow-Origin https://pr.crazyaigc.com always;
    add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
    add_header Access-Control-Allow-Headers "Authorization, Content-Type" always;
    add_header Access-Control-Allow-Credentials true always;

    if ($request_method = 'OPTIONS') {
      return 204;
    }
  }

  # 健康检查
  location /health {
    proxy_pass http://127.0.0.1:8081/health;
    access_log off;
  }
}
```

### Systemd 服务

```ini
[Unit]
Description=PR Business Backend API
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/var/www/pr-business-backend
ExecStart=/var/www/pr-business-backend/server
Restart=always
RestartSec=5
Environment="PORT=8081"
EnvironmentFile=/var/www/pr-business-backend/.env

[Install]
WantedBy=multi-user.target
```

---

## 📊 性能优化

### 1. 数据库索引

```sql
-- 用户表
CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_auth_center ON users(auth_center_id);
CREATE INDEX idx_users_role ON users(role);

-- 任务表
CREATE INDEX idx_tasks_merchant ON tasks(merchant_id);
CREATE INDEX idx_tasks_provider ON tasks(provider_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_created ON tasks(created_at);

-- 任务分配表
CREATE INDEX idx_assignments_task ON task_assignments(task_id);
CREATE INDEX idx_assignments_creator ON task_assignments(creator_id);
CREATE INDEX idx_assignments_status ON task_assignments(status);

-- 积分历史表
CREATE INDEX idx_gold_history_user ON gold_coin_credit_histories(user_id);
CREATE INDEX idx_gold_history_created ON gold_coin_credit_histories(created_at);
```

### 2. 缓存策略

**不使用 Redis（简化部署）**:
- 利用 PostgreSQL 的查询缓存
- 前端使用 TanStack Query 缓存 API 响应

### 3. 分页查询

```go
func GetTasks(page, limit int) ([]Task, int64, error) {
  var tasks []Task
  var total int64

  db.Model(&Task{}).Count(&total)

  offset := (page - 1) * limit
  err := db.Limit(limit).Offset(offset).Find(&tasks).Error

  return tasks, total, err
}
```

---

## 🧪 测试策略

### 1. 单元测试

**积分系统测试**:
```go
func TestAddCredits(t *testing.T) {
  // 初始化测试数据库
  db := setupTestDB()
  defer db.Close()

  // 创建测试用户
  user := createTestUser(db)

  // 添加积分
  err := AddCredits(db, user.ID, 1000, "RECHARGE", "充值")

  // 验证
  assert.NoError(t, err)

  var account CreditAccount
  db.Where("user_id = ?", user.ID).First(&account)
  assert.Equal(t, int64(1000), account.GoldCoins)
}
```

### 2. 集成测试

**抢单并发测试**:
```go
func TestAcceptTaskConcurrency(t *testing.T) {
  db := setupTestDB()
  defer db.Close()

  task := createTestTask(db, 10) // 10个名额

  // 并发抢单
  var wg sync.WaitGroup
  successCount := 0
  mutex := sync.Mutex{}

  for i := 0; i < 20; i++ { // 20人抢10个名额
    wg.Add(1)
    go func(creatorID string) {
      defer wg.Done()
      err := AcceptTask(db, task.ID, creatorID)
      if err == nil {
        mutex.Lock()
        successCount++
        mutex.Unlock()
      }
    }(fmt.Sprintf("creator-%d", i))
  }

  wg.Wait()

  // 验证：只有10人成功
  assert.Equal(t, 10, successCount)
}
```

### 3. API 测试

```go
func TestGetTasks(t *testing.T) {
  router := setupRouter()

  w := httptest.NewRecorder()
  req := httptest.NewRequest("GET", "/api/tasks?page=1&limit=20", nil)

  router.ServeHTTP(w, req)

  assert.Equal(t, 200, w.Code)

  var response map[string]interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  assert.True(t, response["success"].(bool))
}
```

---

## 📝 开发规范

### Go 代码规范

1. **包命名**: 小写单词，不使用下划线
2. **文件命名**: snake_case
3. **结构体命名**: PascalCase
4. **常量命名**: PascalCase 或 UPPER_SNAKE_CASE
5. **错误处理**: 永远不要忽略错误

### React 代码规范

1. **组件命名**: PascalCase
2. **文件命名**:
   - 组件: PascalCase.tsx
   - 工具: camelCase.ts
   - 类型: camelCase.types.ts
3. **Hooks**: 必须以 `use` 开头
4. **状态管理**: 优先使用 TanStack Query

### Git 提交规范

```
feat(task): 添加任务创建API
fix(credits): 修复积分并发问题
docs(readme): 更新部署文档
refactor(auth): 重构认证中间件
test(task): 添加任务测试用例
```

---

## ✅ 验收标准

### 功能完整性

- [ ] 所有 MVP 功能已实现
- [ ] 所有 API 已测试通过
- [ ] 所有页面可正常访问
- [ ] 核心业务流程可闭环

### 性能指标

- [ ] API 平均响应时间 < 500ms
- [ ] 抢单并发无超卖
- [ ] 前端首屏加载 < 2s
- [ ] 数据库查询有索引

### 安全性

- [ ] 所有 API 有认证
- [ ] 敏感操作有权限控制
- [ ] SQL 注入防护
- [ ] XSS 攻击防护

### 代码质量

- [ ] 所有函数有注释
- [ ] 错误处理完善
- [ ] 没有硬编码的配置
- [ ] Git 提交规范

---

## 📚 参考资料

- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Prisma 文档](https://www.prisma.io/docs/)
- [React Router 文档](https://reactrouter.com/)
- [TanStack Query 文档](https://tanstack.com/query/latest)
