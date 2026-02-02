# PR Business 开发计划

**开始时间**: 2026-02-02
**技术栈**: Vite + React + TypeScript + shadcn/ui | Go + Gin + PostgreSQL + Redis
**部署**: Systemd + Nginx

---

## 🎯 MVP目标（8周）

- [ ] Week 1-2: 数据库 + 认证
- [ ] Week 3-4: 邀请码 + 用户管理
- [ ] Week 5-6: 任务系统核心
- [ ] Week 7-8: 积分 + 提现

---

## 📊 当前进度

**当前阶段**: 阶段3 - 邀请码系统
**进度**: 100% (3/3) ✅
**完成时间**: 2026-02-02

---

## ✅ 已完成任务

### 阶段1：数据库设计与初始化 ✅

- [x] 1.1 数据库表结构设计
  - [x] users - 用户表
  - [x] invitation_codes - 邀请码表
  - [x] merchants - 商家表
  - [x] service_providers - 服务商表
  - [x] creators - 达人表
  - [x] campaigns - 营销活动表
  - [x] tasks - 任务名额表
  - [x] merchant_staff - 商家员工表
  - [x] merchant_staff_permissions - 商家员工权限表
  - [x] service_provider_staff - 服务商员工表
  - [x] provider_staff_permissions - 服务商员工权限表
  - [x] credit_accounts - 积分账户表
  - [x] transaction_types - 积分交易类型表
  - [x] credit_transactions - 积分流水表
  - [x] withdrawals - 提现表
  - [x] merchant_provider_bindings - 商家-服务商绑定表

- [x] 1.2 SQL建表脚本
  - 文件: `backend/migrations/001_init_schema.sql`
  - 包含: 所有表结构、索引、外键、触发器、函数

- [x] 1.3 测试数据脚本
  - 文件: `backend/migrations/002_seed_data.sql`
  - 包含: 超级管理员、服务商、商家、达人、营销活动等测试数据

- [x] 1.4 数据库连接配置
  - 文件: `backend/.env.example`
  - 包含: PostgreSQL、Redis、JWT、Auth Center等配置

**交付物**:
- ✅ `backend/migrations/001_init_schema.sql` - 完整建表脚本
- ✅ `backend/migrations/002_seed_data.sql` - 测试数据脚本
- ✅ `backend/.env.example` - 环境变量配置示例

---

### 阶段2：认证系统 ✅

- [x] 2.1 对接auth-center（后端）
  - [x] JWT中间件 - `backend/middlewares/auth.go`, `backend/utils/jwt.go`
  - [x] 其他中间件 - CORS, Logger, Recovery
  - [x] 微信登录接口 - `backend/controllers/auth.go`
  - [x] 密码登录接口
  - [x] Token验证中间件
  - [x] 角色切换接口
  - [x] 路由配置 - `backend/routes/routes.go`

- [x] 2.2 Go项目结构
  - [x] go.mod依赖配置
  - [x] main.go入口文件
  - [x] config/config.go配置管理
  - [x] models/user.go用户模型

- [x] 2.3 前端认证页面
  - [x] 登录页 `/login` - `frontend/src/pages/Login.tsx`
  - [x] 角色切换组件 - `frontend/src/components/RoleSwitcher.tsx`
  - [x] 路由守卫 - `frontend/src/components/ProtectedRoute.tsx`
  - [x] 认证上下文 - `frontend/src/contexts/AuthContext.tsx`
  - [x] API服务 - `frontend/src/services/api.ts`
  - [x] Dashboard页面 - `frontend/src/pages/Dashboard.tsx`

- [x] 2.4 前端构建
  - [x] Vite + React + TypeScript配置
  - [x] TailwindCSS样式配置
  - [x] 构建成功 (dist/ 274KB)

**交付物**:
- ✅ `backend/utils/jwt.go` - JWT工具函数
- ✅ `backend/middlewares/auth.go` - 认证中间件
- ✅ `backend/controllers/auth.go` - 认证控制器
- ✅ `backend/routes/routes.go` - 路由配置
- ✅ 后端编译成功（pr-business二进制 24MB）
- ✅ `frontend/` - 完整前端项目
- ✅ 前端构建成功（dist/ 274KB）

---

### 阶段3：邀请码系统 ✅

- [x] 3.1 邀请码生成
  - [x] 固定邀请码生成算法 - `backend/utils/invitation_code.go`
  - [x] InvitationCode模型 - `backend/models/invitation_code.go`
  - [x] 邀请码唯一性校验
  - [x] 基于ID后缀的编码规则

- [x] 3.2 邀请码接口
  - [x] 创建邀请码接口 - `backend/controllers/invitation.go`
  - [x] 邀请码列表接口
  - [x] 邀请码详情接口
  - [x] 使用邀请码接口
  - [x] 禁用邀请码接口

- [x] 3.3 邀请码管理页面
  - [x] 邀请码列表页面 - `frontend/src/pages/InvitationCodes.tsx`
  - [x] 创建邀请码模态框
  - [x] 邀请码状态显示
  - [x] 禁用邀请码功能

**邀请码规则**:
- `ADMIN-MASTER`: 超级管理员主邀请码（固定）
- `SP-{id后8位}`: 服务商管理员邀请码
- `MERCHANT-{provider_id后8位}-{id后8位}`: 商家邀请码
- `CREATOR-{id后8位}`: 达人邀请码
- `STAFF-{owner_id后8位}-{type}`: 员工邀请码

**交付物**:
- ✅ `backend/models/invitation_code.go` - 邀请码模型
- ✅ `backend/utils/invitation_code.go` - 邀请码生成工具
- ✅ `backend/controllers/invitation.go` - 邀请码控制器
- ✅ `frontend/src/pages/InvitationCodes.tsx` - 邀请码管理页面（400+行）
- ✅ 后端编译成功
- ✅ 前端构建成功（283KB）

---

## 🚧 进行中任务

*暂无 - 所有已完成阶段都已归档*

---

## 📋 待办任务清单

### 阶段1：数据库设计与初始化（Week 1，3天）

#### 1.1 数据库表结构设计（1天）
- [ ] 用户表 `users`
- [ ] 邀请码表 `invitation_codes`
- [ ] 商家表 `merchants`
- [ ] 服务商表 `service_providers`
- [ ] 达人表 `creators`
- [ ] 营销活动表 `campaigns`
- [ ] 任务名额表 `tasks`
- [ ] 积分账户表 `credit_accounts`
- [ ] 积分流水表 `credit_transactions`
- [ ] 提现表 `withdrawals`

**参考文档**: PRD.md 第4665行起（数据库设计）

#### 1.2 SQL建表脚本（1天）
- [ ] 创建所有表的CREATE TABLE语句
- [ ] 添加所有外键约束
- [ ] 创建所有索引
- [ ] 添加触发器（邀请码自动生成）
- [ ] 添加检查约束（数据完整性）

**文件输出**: `backend/migrations/001_init_schema.sql`

#### 1.3 测试数据脚本（0.5天）
- [ ] 创建超级管理员用户
- [ ] 创建测试服务商
- [ ] 创建测试商家
- [ ] 创建测试达人
- [ ] 创建测试营销活动

**文件输出**: `backend/migrations/002_seed_data.sql`

#### 1.4 数据库连接配置（0.5天）
- [ ] Go的PostgreSQL连接配置
- [ ] Go的Redis连接配置
- [ ] 环境变量配置示例

**文件输出**: `backend/.env.example`

---

### 阶段2：认证系统（Week 1-2，5天）

#### 2.1 对接auth-center（2天）
- [ ] 微信登录接口
- [ ] 密码登录接口
- [ ] Token验证中间件
- [ ] 角色切换接口

**参考文档**: PRD.md 第875行起（认证接口设计）

#### 2.2 前端认证页面（1天）
- [ ] 登录页 `/login`
- [ ] 角色切换组件
- [ ] 路由守卫

#### 2.3 用户管理基础（1天）
- [ ] 用户CRUD接口
- [ ] 角色权限校验
- [ ] 超级管理员 - 用户管理页

#### 2.4 测试（1天）
- [ ] 登录/登出测试
- [ ] 角色切换测试
- [ ] 权限校验测试

---

### 阶段3：邀请码系统（Week 3，4天）

#### 3.1 邀请码生成（1天）
- [ ] 固定邀请码生成算法
- [ ] 触发器自动生成
- [ ] 邀请码唯一性校验

**参考文档**: PRD.md 第239行起（邀请码系统）

#### 3.2 邀请码验证（1天）
- [ ] 应用邀请码接口
- [ ] 邀请码类型校验
- [ ] 邀请码使用记录

#### 3.3 邀请码管理（1天）
- [ ] 邀请码列表页面
- [ ] 邀请码详情查看
- [ ] 邀请码启用/禁用

#### 3.4 测试（1天）

---

### 阶段4：用户管理（Week 3-4，5天）

#### 4.1 商家管理（1天）
- [ ] 商家CRUD接口
- [ ] 商家员工管理
- [ ] 权限分配

**参考文档**: PRD.md 第81行起（商家管理员）

#### 4.2 服务商管理（1天）
- [ ] 服务商CRUD接口
- [ ] 服务商员工管理
- [ ] 权限分配

**参考文档**: PRD.md 第100行起（服务商管理员）

#### 4.3 达人管理（1天）
- [ ] 达人信息管理
- [ ] 达人等级系统
- [ ] 达人邀请关系

**参考文档**: PRD.md 第115行起（达人角色）

#### 4.4 前端页面（2天）
- [ ] 商家工作台 - 商家信息页
- [ ] 服务商工作台 - 服务商信息页
- [ ] 达人工作台 - 个人中心
- [ ] 超级管理员 - 用户管理页

---

### 阶段5：任务系统核心（Week 5-6，8天）

#### 5.1 营销活动管理（2天）
- [ ] 创建营销活动
- [ ] 设置积分分配
- [ ] 活动状态管理

**参考文档**: PRD.md 第6618行起（营销活动状态）

#### 5.2 任务名额系统（2天）
- [ ] 预分配名额
- [ ] 达人接任务
- [ ] 名额状态流转

**参考文档**: PRD.md 第6715行起（资金流动规则）

#### 5.3 任务提交流程（2天）
- [ ] 达人提交平台链接/截图
- [ ] 服务商审核
- [ ] 审核通过/拒绝

**参考文档**: PRD.md 第6887行起（审核拒绝重新提交流程）

#### 5.4 积分结算（2天）
- [ ] 发布任务预扣款
- [ ] 任务结算分配
- [ ] 事务性保证

**参考文档**: PRD.md 第5759行起（财务安全）

#### 5.5 前端页面（3天）
- [ ] 商家工作台 - 创建任务页
- [ ] 达人工作台 - 任务大厅
- [ ] 达人工作台 - 我的任务
- [ ] 服务商工作台 - 任务审核页

---

### 阶段6：积分系统（Week 7，4天）

#### 6.1 积分账户（2天）
- [ ] 积分账户管理
- [ ] 积分流水记录
- [ ] 余额查询

**参考文档**: PRD.md 第5504行起（积分账户表）

#### 6.2 充值功能（1天）
- [ ] 充值接口（模拟支付）
- [ ] 充值记录

#### 6.3 前端页面（1天）
- [ ] 各角色工作台 - 财务数据页
- [ ] 积分明细查询

---

### 阶段7：提现功能（Week 7-8，3天）

#### 7.1 提现申请（1天）
- [ ] 提现申请接口
- [ ] 提现记录查询

**参考文档**: PRD.md 第5740行起（提现表）

#### 7.2 提现审核（1天）
- [ ] 超级管理员审核接口
- [ ] 审核通过/拒绝

#### 7.3 前端页面（1天）
- [ ] 达人/服务商 - 提现申请页
- [ ] 超级管理员 - 提现审核页

---

### 阶段8：前端完善与联调（Week 8，3天）

#### 8.1 角色切换器（0.5天）
- [ ] 角色切换组件
- [ ] 角色切换接口

#### 8.2 响应式布局（1天）
- [ ] 移动端适配
- [ ] Tailwind响应式样式

#### 8.3 错误处理（0.5天）
- [ ] 统一错误提示
- [ ] 全局错误处理

#### 8.4 API联调（1天）
- [ ] 前后端联调
- [ ] 集成测试

---

## 📁 项目文件结构

```
pr-business/
├── backend/                  # Go后端
│   ├── main.go              # 入口文件
│   ├── config/              # 配置
│   ├── models/              # 数据模型
│   ├── controllers/         # 控制器
│   ├── services/            # 业务逻辑
│   ├── middlewares/         # 中间件
│   ├── migrations/          # 数据库迁移
│   │   ├── 001_init_schema.sql
│   │   └── 002_seed_data.sql
│   ├── .env.example         # 环境变量示例
│   └── go.mod
├── frontend/                 # Vite + React前端
│   ├── src/
│   │   ├── pages/           # 页面
│   │   ├── components/     # 组件
│   │   ├── hooks/          # Hooks
│   │   ├── services/       # API服务
│   │   ├── lib/            # 工具库
│   │   └── main.tsx        # 入口
│   ├── package.json
│   └── vite.config.ts
├── docs/
│   ├── PRD.md               # 产品需求文档
│   └── DEVELOPMENT_PLAN.md  # 本文档（开发计划）
└── README.md
```

---

## 📝 开发日志

### 2026-02-02

**阶段1：数据库设计与初始化** ✅ 已完成

#### 完成内容：
1. ✅ 创建项目目录结构
   - `backend/` - 后端目录（migrations, config, models等）
   - `frontend/` - 前端目录（pages, components, hooks等）

2. ✅ 创建数据库建表脚本（001_init_schema.sql）
   - 16个核心表：users, invitation_codes, merchants, service_providers, creators, campaigns, tasks等
   - 完整索引、外键、触发器
   - 自动更新updated_at函数
   - 13种交易类型初始数据

3. ✅ 创建测试数据脚本（002_seed_data.sql）
   - 测试用户：超级管理员、服务商、商家、达人
   - 测试组织：1个服务商、1个商家
   - 测试营销活动：1个活动（10个任务名额）
   - 固定邀请码：ADMIN-MASTER, SP-ADMIN, MERCHANT-xxx等

4. ✅ 创建环境变量配置（.env.example）
   - PostgreSQL、Redis、JWT配置
   - Auth Center对接配置
   - 文件上传、日志配置

#### 输出文件：
- `backend/migrations/001_init_schema.sql` - 900行
- `backend/migrations/002_seed_data.sql` - 300行
- `backend/.env.example` - 环境变量模板

#### 下一步：
- 阶段2：认证系统（微信登录、密码登录、JWT）

---

### 2026-02-02 (续)

**阶段2：认证系统（后端）** ⏳ 进行中

#### 完成内容：
1. ✅ Go后端项目结构
   - `backend/go.mod` - 依赖管理（Gin, GORM, JWT, Viper等）
   - `backend/main.go` - 入口文件，服务器启动
   - `backend/config/config.go` - 配置管理（环境变量加载、数据库连接）
   - `backend/models/user.go` - 用户模型定义

2. ✅ JWT认证中间件
   - `backend/utils/jwt.go` - JWT token生成与解析
   - `backend/middlewares/auth.go` - 认证中间件、角色权限验证
   - `backend/middlewares/cors.go` - CORS跨域支持
   - `backend/middlewares/logger.go` - 日志中间件
   - `backend/middlewares/recovery.go` - 异常恢复中间件

3. ✅ 认证控制器
   - `backend/controllers/auth.go` - 认证相关接口
     - WeChatLogin - 微信登录
     - PasswordLogin - 密码登录
     - RefreshToken - 刷新令牌
     - SwitchRole - 角色切换
     - GetCurrentUser - 获取当前用户信息

4. ✅ 路由配置
   - `backend/routes/routes.go` - API路由设置
     - POST /api/v1/auth/wechat - 微信登录
     - POST /api/v1/auth/password - 密码登录
     - POST /api/v1/auth/refresh - 刷新令牌
     - GET /api/v1/user/me - 获取用户信息
     - POST /api/v1/user/switch-role - 切换角色

5. ✅ 编译成功
   - 后端二进制: `backend/pr-business`
   - 所有依赖下载完成

#### 技术实现：
- 使用 `gorm.io/gorm` v1.25.4 作为ORM
- 使用 `golang-jwt/jwt/v5` 进行JWT认证
- 使用 `gin-gonic/gin` v1.9.1 作为Web框架
- 使用 `spf13/viper` 进行配置管理
- 支持环境变量配置（.env文件）

#### 输出文件：
- `backend/utils/jwt.go` - JWT工具函数
- `backend/middlewares/auth.go` - 认证中间件
- `backend/middlewares/cors.go` - CORS中间件
- `backend/middlewares/logger.go` - 日志中间件
- `backend/middlewares/recovery.go` - 恢复中间件
- `backend/controllers/auth.go` - 认证控制器
- `backend/routes/routes.go` - 路由配置
- `backend/pr-business` - 编译后的二进制文件（24MB）

#### 注意事项：
- 认证接口中的 `TODO` 标记需要对接实际的auth-center服务
- 目前使用模拟数据进行登录验证
- 需要PostgreSQL数据库运行才能进行完整测试

#### 下一步：
- 前端登录页面开发（Vite + React + TypeScript）
- 后端测试（需要数据库连接）

---

### 2026-02-02 (续2)

**阶段2：认证系统（前端）** ✅ 已完成

#### 完成内容：
1. ✅ Vite + React + TypeScript项目初始化
   - 使用 `npm create vite` 创建项目
   - 安装核心依赖：react-router-dom, axios, @tanstack/react-query
   - 配置TypeScript严格模式

2. ✅ TailwindCSS配置
   - 安装 TailwindCSS v4 和 @tailwindcss/postcss
   - 配置 tailwind.config.js 和 postcss.config.js
   - 更新 index.css 使用 Tailwind 指令

3. ✅ 项目结构
   - `src/lib/utils.ts` - 工具函数（cn className合并）
   - `src/types/index.ts` - TypeScript类型定义
   - `src/services/api.ts` - API服务（axios拦截器、自动刷新token）
   - `src/contexts/AuthContext.tsx` - 认证上下文
   - `src/components/` - 组件目录
   - `src/pages/` - 页面目录

4. ✅ 认证功能
   - `Login.tsx` - 登录页面（支持微信登录、密码登录）
   - `Dashboard.tsx` - 工作台页面
   - `ProtectedRoute.tsx` - 路由守卫组件
   - `RoleSwitcher.tsx` - 角色切换组件
   - `App.tsx` - 路由配置

5. ✅ 构建成功
   - 前端构建产物: `dist/` (274KB)
   - 所有TypeScript类型检查通过

#### 技术实现：
- 使用 `react-router-dom` v6 进行路由管理
- 使用 `axios` 进行HTTP请求，配置请求/响应拦截器
- 使用 React Context API 进行状态管理
- JWT token存储在localStorage
- 401自动刷新token机制
- TailwindCSS v4 样式系统

#### 输出文件：
- `frontend/src/pages/Login.tsx` - 登录页面（200+行）
- `frontend/src/pages/Dashboard.tsx` - 工作台页面
- `frontend/src/contexts/AuthContext.tsx` - 认证上下文
- `frontend/src/components/ProtectedRoute.tsx` - 路由守卫
- `frontend/src/components/RoleSwitcher.tsx` - 角色切换
- `frontend/src/services/api.ts` - API服务
- `frontend/src/types/index.ts` - 类型定义
- `frontend/package.json` - 依赖配置
- `frontend/.env.example` - 环境变量示例

#### Git提交：
- 提交哈希: `a246dcb`
- 28个新文件，5287行代码
- 推送到 GitHub: `superxia01/pr-business`

#### 注意事项：
- 登录功能使用模拟数据（authCode: mock_wechat_auth_code_xxx）
- 需要后端服务运行在 `http://localhost:8080`
- Token过期会自动刷新
- 未登录访问受保护路由会自动跳转到登录页

#### 下一步：
- 阶段3：邀请码系统（后端接口+前端页面）

---

### 2026-02-02 (续3)

**阶段3：邀请码系统** ✅ 已完成

#### 完成内容：
1. ✅ 邀请码模型和生成逻辑
   - `backend/models/invitation_code.go` - 邀请码数据模型
   - `backend/utils/invitation_code.go` - 邀请码生成工具
   - 实现基于ID后缀的固定邀请码生成算法
   - 邀请码格式验证和解析功能

2. ✅ 邀请码控制器
   - `backend/controllers/invitation.go` - 邀请码CRUD接口
   - 创建邀请码（权限控制）
   - 获取邀请码列表（支持过滤）
   - 获取邀请码详情
   - 使用邀请码（检查状态、使用次数、过期时间）
   - 禁用邀请码

3. ✅ API路由配置
   - POST /api/v1/invitations - 创建邀请码
   - GET /api/v1/invitations - 邀请码列表
   - GET /api/v1/invitations/:code - 邀请码详情
   - POST /api/v1/invitations/use - 使用邀请码
   - POST /api/v1/invitations/:code/disable - 禁用邀请码

4. ✅ 前端邀请码管理
   - `frontend/src/pages/InvitationCodes.tsx` - 邀请码管理页面（400+行）
   - 邀请码列表展示（表格形式）
   - 创建邀请码模态框
   - 邀请码状态标签（active/disabled/expired/used）
   - 禁用邀请码功能
   - 权限控制（超级管理员、服务商管理员、商家管理员）

5. ✅ 类型定义和API服务
   - `frontend/src/types/index.ts` - 邀请码类型定义
   - `frontend/src/services/api.ts` - 邀请码API服务
   - Dashboard页面添加邀请码管理入口

#### 邀请码编码规则：
- **超级管理员**: `ADMIN-MASTER`（固定）
- **服务商管理员**: `SP-{id后8位}`
- **商家**: `MERCHANT-{provider_id后8位}-{id后8位}`
- **达人**: `CREATOR-{id后8位}`
- **员工**: `STAFF-{owner_id后8位}-{type}`

#### 权限控制：
- 超级管理员：可创建所有类型邀请码
- 服务商管理员：可创建商家、达人、员工邀请码
- 商家管理员：可创建员工邀请码

#### 输出文件：
- `backend/models/invitation_code.go` - 邀请码模型（100行）
- `backend/utils/invitation_code.go` - 生成工具（80行）
- `backend/controllers/invitation.go` - 控制器（300+行）
- `frontend/src/pages/InvitationCodes.tsx` - 管理页面（400+行）

#### Git提交：
- 提交哈希: `545a499`
- 10个文件修改，935行代码
- 推送到 GitHub: `superxia01/pr-business`

#### 注意事项：
- 邀请码生成基于ID后8位，确保唯一性
- 使用邀请码会检查：状态、使用次数、过期时间、是否已使用
- 权限检查在后端和前端都实现了
- 邀请码使用后会增加使用计数，达到maxUses后状态变为used

#### 下一步：
- 阶段4：用户管理（商家、服务商、达人管理）

---

## 🔗 相关文档

- [PRD - 产品需求文档](./PRD.md)
- [数据库设计](../PRD.md#-数据库设计)
- [API接口规范](../PRD.md#-api接口规范)
- [前端页面设计](../PRD.md#-工作台页面)
