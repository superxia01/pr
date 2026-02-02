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

**当前阶段**: 阶段2 - 认证系统
**进度**: 60% (3/5) ⏳
**更新时间**: 2026-02-02

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

### 阶段2：认证系统 ⏳

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

- [ ] 2.3 前端认证页面
  - [ ] 登录页 `/login`
  - [ ] 角色切换组件
  - [ ] 路由守卫

- [ ] 2.4 测试
  - [ ] 登录/登出测试（需要数据库）
  - [ ] 角色切换测试
  - [ ] 权限校验测试

**交付物**:
- ✅ `backend/utils/jwt.go` - JWT工具函数
- ✅ `backend/middlewares/auth.go` - 认证中间件
- ✅ `backend/controllers/auth.go` - 认证控制器
- ✅ `backend/routes/routes.go` - 路由配置
- ✅ 后端编译成功（pr-business二进制）

---

## 🚧 进行中任务

**阶段2.3：前端认证页面** - 准备开始
- 创建前端登录页面
- 实现角色切换组件
- 添加路由守卫

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

## 🔗 相关文档

- [PRD - 产品需求文档](./PRD.md)
- [数据库设计](../PRD.md#-数据库设计)
- [API接口规范](../PRD.md#-api接口规范)
- [前端页面设计](../PRD.md#-工作台页面)
