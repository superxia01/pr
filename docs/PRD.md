# PR Business 达人营销平台 - 产品需求文档 (PRD)

**版本**: v1.0 MVP
**日期**: 2026-02-02
**状态**: ✅ 已完成并确认

**技术栈**：
- **前端**：Vite + React + TypeScript + TailwindCSS + shadcn/ui
- **后端**：Go + Gin + PostgreSQL + Redis
- **部署**：Systemd + Nginx

**📌 开发计划**: [DEVELOPMENT_PLAN.md](./DEVELOPMENT_PLAN.md) - 实时更新开发进度
**📌 PRD文档**: 本文档 - 详细的产品需求和技术规范

---

## 🚀 MVP开发顺序与内容索引（8周快速上线）

### 📅 开发时间线

```
Week 1-2: 数据库 + 认证
Week 3-4: 邀请码 + 用户管理
Week 5-6: 任务系统核心
Week 7-8: 积分 + 提现
```

**⚠️ 开发提示**：
1. 先查看 [DEVELOPMENT_PLAN.md](./DEVELOPMENT_PLAN.md) 了解当前进度
2. 按照开发计划逐步完成，每个阶段完成后更新进度
3. 遇到问题时查阅本文档对应章节

---

## 📋 开发顺序与文档索引

### 阶段1：数据库设计与初始化（Week 1，3天）

**目标**：完成数据库表结构设计与初始化脚本

| 开发任务 | 文档位置 | 优先级 |
|---------|---------|-------|
| 核心表结构设计 | [数据库设计](#-数据库设计) | P0 |
| SQL建表脚本 | [各表结构定义](#-2-商家表-) | P0 |
| 外键与索引 | [索引与约束](#-2-商家表-) | P0 |

**关键表（按开发顺序）**：
1. `users` - 用户表
2. `invitation_codes` - 邀请码表
3. `merchants` - 商家表
4. `service_providers` - 服务商表
5. `creators` - 达人表
6. `campaigns` - 营销活动表
7. `tasks` - 任务名额表
8. `credit_accounts` - 积分账户表
9. `credit_transactions` - 积分流水表
10. `withdrawals` - 提现表

**交付物**：
- [ ] PostgreSQL建表SQL脚本
- [ ] 索引与外键约束
- [ ] 测试数据种子脚本

---

### 阶段2：认证系统（Week 1-2，5天）

**目标**：对接auth-center，实现微信登录和密码登录

| 开发任务 | 文档位置 | 优先级 |
|---------|---------|-------|
| 认证流程设计 | [认证系统](#-认证接口设计) | P0 |
| 用户角色管理 | [用户角色定义](#-角色说明) | P0 |
| 角色切换器 | [角色切换器](#-多角色系统) | P0 |

**核心接口**：
- `GET /api/v1/auth/wechat/login` - 微信登录
- `POST /api/v1/auth/password/login` - 密码登录
- `GET /api/v1/auth/me` - 获取当前用户
- `POST /api/v1/auth/logout` - 登出
- `POST /api/v1/auth/switch-role` - 角色切换

**前端页面**：
- [ ] 登录页 `/login`
- [ ] 角色切换组件

**交付物**：
- [ ] 认证中间件（JWT）
- [ ] 角色权限校验
- [ ] 登录/登出功能

---

### 阶段3：邀请码系统（Week 3，4天）

**目标**：实现固定邀请码系统，支持用户注册与角色分配

| 开发任务 | 文档位置 | 优先级 |
|---------|---------|-------|
| 邀请码类型 | [邀请码系统](#-统一邀请码系统) | P0 |
| 固定邀请码生成 | [固定邀请码生成规则](#-固定邀请码生成规则) | P0 |
| 邀请码验证 | [应用邀请码接口](#-3-密码登录代理接口) | P0 |

**核心邀请码类型**：
- `SP-ADMIN` - 创建服务商管理员（超级管理员使用）
- `MERCHANT-{服务商ID后6位}` - 邀请新商家（服务商使用）
- `MSTAFF-{商家ID后6位}` - 邀请商家员工（商家使用）
- `SPSTAFF-{服务商ID后6位}` - 邀请服务商员工（服务商使用）
- `CREATOR-{员工ID后6位}` - 邀请达人（员工使用）

**核心接口**：
- `POST /api/v1/auth/apply-invite-code` - 应用邀请码

**交付物**：
- [ ] 邀请码生成算法
- [ ] 邀请码验证接口
- [ ] 邀请码管理页面

---

### 阶段4：用户管理（Week 3-4，5天）

**目标**：实现超级管理员、服务商、商家、达人的CRUD

| 开发任务 | 文档位置 | 优先级 |
|---------|---------|-------|
| 角色权限矩阵 | [权限矩阵](#-权限矩阵) | P0 |
| 商家管理 | [商家管理员](#-merchant_admin) | P0 |
| 服务商管理 | [服务商管理员](#-service_provider_admin) | P0 |
| 达人管理 | [达人角色](#-creator) | P0 |

**前端页面**：
- [ ] 超级管理员 - 用户管理页
- [ ] 商家工作台 - 商家信息页
- [ ] 服务商工作台 - 服务商信息页
- [ ] 达人工作台 - 个人中心

**交付物**：
- [ ] 用户CRUD接口
- [ ] 角色权限管理
- [ ] 用户信息编辑

---

### 阶段5：任务系统核心（Week 5-6，8天）

**目标**：实现商家创建任务、达人接任务、服务商审核

| 开发任务 | 文档位置 | 优先级 |
|---------|---------|-------|
| 营销活动状态 | [营销活动状态](#-营销活动状态流程) | P0 |
| 任务名额状态 | [活动名额状态](#-活动名额状态) | P0 |
| 资金流动规则 | [资金流动规则](#-资金流动规则) | P0 |

**核心流程**：
1. 商家创建任务 → 服务商设置积分分配并发布
2. 达人接任务 → 达人提交平台链接/截图
3. 服务商审核通过/拒绝

**前端页面**：
- [ ] 商家工作台 - 创建任务页
- [ ] 达人工作台 - 任务大厅
- [ ] 达人工作台 - 我的任务
- [ ] 服务商工作台 - 任务审核页

**交付物**：
- [ ] 任务发布功能
- [ ] 任务大厅（接任务）
- [ ] 任务提交流程
- [ ] 任务审核流程

---

### 阶段6：积分系统（Week 7，4天）

**目标**：实现积分账户、充值、积分流水

| 开发任务 | 文档位置 | 优先级 |
|---------|---------|-------|
| 积分账户设计 | [积分账户表](#-13-积分账户表) | P0 |
| 交易类型 | [交易类型枚举](#-积分交易类型枚举) | P0 |
| 财务安全 | [事务性要求](#-财务安全与一致性保证) | P0 |

**核心功能**：
- 商家充值（`RECHARGE`）
- 发布任务预扣款（`TASK_PUBLISH`）
- 任务结算（`TASK_INCOME`、`STAFF_REFERRAL`、`PROVIDER_INCOME`）

**前端页面**：
- [ ] 商家工作台 - 充值页
- [ ] 各角色工作台 - 积分明细页

**交付物**：
- [ ] 积分账户系统
- [ ] 充值功能（模拟支付）
- [ ] 积分流水记录

---

### 阶段7：提现功能（Week 7-8，3天）

**目标**：实现达人、服务商申请提现，超级管理员审核

| 开发任务 | 文档位置 | 优先级 |
|---------|---------|-------|
| 提现流程 | [提现审核](#-提现审核通过超级管理员) | P0 |
| 提现表设计 | [提现表](#-16-提现表) | P0 |

**核心流程**：
1. 用户申请提现（扣除balance，转入frozen_balance）
2. 超级管理员审核（通过/拒绝）
3. 审核通过后打款（模拟打款）

**前端页面**：
- [ ] 达人/服务商工作台 - 提现申请页
- [ ] 超级管理员 - 提现审核页

**交付物**：
- [ ] 提现申请功能
- [ ] 提现审核功能
- [ ] 提现记录查询

---

### 阶段8：前端完善与联调（Week 8，3天）

**目标**：前端页面完善、API联调、测试

| 开发任务 | 说明 | 优先级 |
|---------|------|-------|
| 角色切换器 | 实现多角色切换 | P0 |
| 响应式布局 | 移动端适配 | P1 |
| 错误处理 | 统一错误提示 | P0 |
| API联调 | 前后端联调 | P0 |

**交付物**：
- [ ] 完整的前端界面
- [ ] API联调完成
- [ ] 基础测试通过

---

## 📚 快速导航

### 核心业务流程
- [角色权限说明](#-角色说明已确认)
- [邀请码系统](#-统一邀请码系统)
- [营销活动状态流程](#-营销活动状态流程)
- [资金流动规则](#-资金流动规则)
- [财务安全机制](#-财务安全与一致性保证)

### 数据库表结构
- [用户表 users](#-1-用户表)
- [邀请码表 invitation_codes](#-7-邀请码表)
- [商家表 merchants](#-2-商家表)
- [服务商表 service_providers](#-3-服务商表)
- [达人表 creators](#-4-达人表)
- [营销活动表 campaigns](#-5-营销活动表)
- [任务名额表 tasks](#-6-任务名额表)
- [积分账户表 credit_accounts](#-13-积分账户表)
- [积分流水表 credit_transactions](#-15-积分流水表)
- [提现表 withdrawals](#-16-提现表)

### API接口
- [认证接口设计](#-认证接口设计)
- [API版本管理](#-api接口规范)
- [统一响应格式](#-统一响应格式)

### 前端页面
- [角色切换器](#-多角色系统)
- [超级管理员工作台](#-超级管理员工作台)
- [商家工作台](#-商家工作台)
- [服务商工作台](#-服务商工作台)
- [达人工作台](#-达人工作台)

---

## 📋 项目概述

### 项目定位
PR Business 是一个连接商家、服务商和达人的任务分发与执行平台。

### 核心流程
```
商家发布任务 → 达人接任务执行 → 服务商审核 → 自动结算积分
```

### 目标用户
- **商家**: 需要达人做营销推广的品牌方
- **服务商**: 管理达人团队、提供审核服务
- **达人**: 执行营销任务、获得收益

---

## 👥 用户角色定义

### 角色列表

| 角色代码 | 角色名称 | 说明 | 默认工作台 |
|---------|---------|------|-----------|
| SUPER_ADMIN | 超级管理员 | 平台运营，管理所有功能 | `/workspace/admin` |
| MERCHANT_ADMIN | 商家管理员 | 商家负责人，拥有所有商家权限 | `/workspace/merchant` |
| MERCHANT_STAFF | 商家员工 | 商家员工，权限由管理员分配 | `/workspace/merchant` |
| SERVICE_PROVIDER_ADMIN | 服务商管理员 | 服务商负责人，审核任务 | `/workspace/service-provider` |
| SERVICE_PROVIDER_STAFF | 服务商员工 | 服务商员工，权限由管理员分配 | `/workspace/service-provider` |
| CREATOR | 达人 | 接任务、执行任务、获得收入 | `/workspace/creator` |

### 多角色系统

**⚠️ 重要**：一个用户可以同时拥有多个角色，可以在不同角色之间切换。

**示例场景**：
- 商家员工：小张既是商家员工，下班后也可以切换到达人工作台接任务
- 服务商员工：小李既是服务商员工（审核任务），也可以到达人工作台接任务
- 商家管理员：小王既是商家管理员，也可以到达人工作台体验产品

**多角色实现**：
1. **用户注册/创建**时，系统根据邀请码或注册方式分配初始角色
2. 用户可以通过邀请码获得新角色：
   - 服务商员工邀请 → 获得 CREATOR 角色
   - 商家管理员邀请 → 获得 MERCHANT_STAFF 角色
   - 服务商管理员邀请 → 获得 SERVICE_PROVIDER_STAFF 角色
3. 用户的 `users.roles` 字段存储所有拥有的角色（JSON数组）
4. 用户可以随时切换 `users.current_role` 来切换工作台

**角色切换**：
- **角色切换器显示规则**：当 `users.roles` 数组长度 > 1 时，在所有工作台的顶部导航显示角色切换器
- 点击切换角色后，跳转到对应的工作台（`/workspace/{role}`）
- 不同角色有独立的权限和数据可见范围
- **积分账户归属**：
  - 商家和服务商的积分账户属于**组织实体**（merchants/service_providers），不属于管理员个人
  - 用户的个人收入（达人接任务、员工返佣）统一进入**个人账户**（USER_PERSONAL）
  - 用户切换角色时，个人余额不变（所有角色共享同一个个人账户）
  - 管理员离职/更换时，组织账户不受影响

---

### 角色说明（已确认）

**SUPER_ADMIN - 超级管理员**
- 职责：平台运营、用户管理、财务审核、创建服务商
- 特殊权限：
  - **可以创建服务商**（为服务商分配管理员）
  - 审核所有提现
  - 查看平台整体数据
- ✅ **已确认**：不需要多级管理员（如运营、财务、技术）
- 权限范围：所有系统权限（可以查看和管理平台所有数据）
- ⚠️ **重要**：超级管理员**不能直接创建商家**，商家必须由服务商通过邀请码创建

**MERCHANT_ADMIN - 商家管理员**
- 职责：管理商家信息、发布任务、充值积分
- 创建流程：
  - **由服务商通过邀请码创建**
  - 服务商生成"商家管理员邀请码" → 商家用邀请码注册 → 成为该商家的管理员
  - 商家自动绑定到生成邀请码的服务商
- 特殊权限：可以添加员工、分配权限
- ✅ **组织结构**：一个商家只有1个管理员，创建者自动成为管理员
- 可用权限：查看任务、创建任务、编辑任务、查看财务、发起充值、管理员工等
- ✅ **权限转移**：管理员可以将管理员权限转移给其他员工

**MERCHANT_STAFF - 商家员工**
- 职责：协助管理商家事务
- 权限管理：由商家管理员通过**最小权限颗粒度**分配
- ✅ 权限系统：采用独立的权限项，管理员可勾选组合（详见"权限管理系统"）
- ✅ **最高权限**：管理员可以勾选所有权限给员工，员工最高权限可以等于管理员
- ✅ **职能名称**：管理员可以给员工设置职能名称（如："运营"、"财务"、"审核专员"），纯展示用途
- 可用权限：查看任务、创建任务、编辑任务、查看财务、发起充值、管理员工等

**SERVICE_PROVIDER_ADMIN - 服务商管理员**
- 职责：管理服务商、邀请商家、审核任务、设置任务积分分配
- 创建流程：
  - **由超级管理员通过邀请码创建**
  - 超级管理员生成"服务商管理员邀请码" → 服务商用邀请码注册 → 成为该服务商的管理员
- 特殊权限：
  - 可以获得任务的"服务商分成"
  - 可以为商家发布的草稿活动设置积分分配并发布
  - 可以添加员工、通过勾选分配权限
  - **可以生成商家邀请码**，邀请新商家注册并绑定
  - 可以生成员工邀请码、达人邀请码
- ✅ **组织结构**：一个服务商只有1个管理员，创建者自动成为管理员
- 可用权限：查看任务、设置积分分配、审核任务、邀请达人、邀请商家、查看财务、发起提现、管理员工等
- ✅ **权限转移**：管理员可以将管理员权限转移给其他员工
- ⚠️ **重要**：服务商负责发展和管理商家，通过邀请码建立绑定关系

**SERVICE_PROVIDER_STAFF - 服务商员工**
- 职责：**邀请达人注册、协助审核任务、管理绑定商家**
- 权限管理：由服务商管理员通过**最小权限颗粒度**分配
- 核心功能：
  - **获得专属邀请码**，邀请新达人注册
  - **获得邀请返佣**：当被邀请的达人完成任务时，员工获得返佣
  - 根据分配的权限执行操作（审核任务、绑定商家、发送任务邀请等）
- ✅ **最高权限**：管理员可以勾选所有权限给员工，员工最高权限可以等于管理员
- ✅ **职能名称**：管理员可以给员工设置职能名称（如："审核专员"、"商务拓展"、"运营"），纯展示用途
- ⚠️ 重要：这是原来"达人组长"角色的核心功能
- ✅ 返佣机制：商家发布任务时，服务商额外设置"员工邀请返佣"作为独立收入项

**CREATOR - 达人**
- 职责：在任务大厅接任务、执行任务、获得收入
- 特殊权限：第一次接任务时自动激活 CREATOR 角色
- ✅ **已确认**：达人可以同时绑定多个服务商（多对多关系）
  - 一个达人可以接受多个服务商的员工邀请
  - 一个服务商的员工可以邀请多个达人
  - 达人接任务时，系统记录是哪个服务商员工的邀请码被使用

---

## 💼 商业模式

---

## 🔐 权限管理系统（最小颗粒度）

### 设计原则

- **最小权限颗粒度**：每个权限都是独立的原子操作
- **灵活组合**：管理员可以勾选任意权限组合分配给员工
- **动态分配**：可以随时修改员工的权限

---

### 商家员工权限列表

| 权限代码 | 权限名称 | 说明 |
|---------|---------|------|
| `task.view` | 查看任务列表 | 可以查看商家的所有任务 |
| `task.create` | 创建任务草稿 | 可以创建任务草稿 |
| `task.edit` | 编辑任务草稿 | 可以编辑草稿状态的任务 |
| `finance.view` | 查看财务数据 | 可以查看余额、收支统计 |
| `finance.recharge` | 发起充值 | 可以发起充值申请 |
| `finance.transaction.view` | 查看积分流水 | 可以查看详细的积分收支记录 |
| `staff.view` | 查看员工列表 | 可以查看商家员工列表 |
| `staff.add` | 添加员工 | 可以邀请新的员工 |
| `staff.remove` | 删除员工 | 可以删除员工 |
| `staff.permission.edit` | 修改员工权限 | 可以分配/修改员工权限 |

**权限分配页面**：`/workspace/merchant/staff/:id/permissions`

---

### 服务商员工权限列表

| 权限代码 | 权限名称 | 说明 |
|---------|---------|------|
| `task.view` | 查看任务列表 | 可以查看绑定商家的任务 |
| `task.allocate` | 设置积分分配 | 可以为草稿活动设置积分并发布 |
| `task.audit` | 审核任务 | 可以审核达人提交的任务 |
| `creator.invite` | 邀请达人 | 可以生成邀请码邀请达人注册 |
| `creator.view` | 查看邀请达人 | 可以查看自己邀请的达人列表 |
| `merchant.bind` | 绑定商家 | 可以绑定新的商家 |
| `merchant.view` | 查看绑定商家 | 可以查看绑定的商家列表 |
| `task.invite` | 任务邀请达人 | 可以生成任务邀请码发送给达人 |
| `finance.view` | 查看财务数据 | 可以查看余额、收支统计 |
| `finance.withdraw` | 发起提现 | 可以发起提现申请 |
| `finance.transaction.view` | 查看积分流水 | 可以查看详细的积分收支记录 |
| `staff.view` | 查看员工列表 | 可以查看服务商员工列表 |
| `staff.add` | 添加员工 | 可以邀请新的员工 |
| `staff.remove` | 删除员工 | 可以删除员工 |
| `staff.permission.edit` | 修改员工权限 | 可以分配/修改员工权限 |

**权限分配页面**：`/workspace/service-provider/staff/:id/permissions`

---

### 权限分配页面设计

**页面布局**（以服务商员工为例）：

```
┌─────────────────────────────────────────────┐
│ 分配权限：张三（服务商员工）                  │
├─────────────────────────────────────────────┤
│ 任务管理                                     │
├─────────────────────────────────────────────┤
│ ☑ 查看任务列表                              │
│ ☑ 设置积分分配并发布任务                    │
│ ☑ 审核任务                                  │
│ ☐ 任务邀请达人                              │
├─────────────────────────────────────────────┤
│ 达人管理                                     │
├─────────────────────────────────────────────┤
│ ☑ 邀请达人注册                              │
│ ☑ 查看邀请达人列表                          │
├─────────────────────────────────────────────┤
│ 商家管理                                     │
├─────────────────────────────────────────────┤
│ ☑ 绑定商家                                  │
│ ☐ 查看绑定商家列表                          │
├─────────────────────────────────────────────┤
│ 财务管理                                     │
├─────────────────────────────────────────────┤
│ ☑ 查看财务数据                              │
│ ☐ 发起提现                                  │
│ ☐ 查看积分流水                              │
├─────────────────────────────────────────────┤
│ 员工管理                                     │
├─────────────────────────────────────────────┤
│ ☐ 查看员工列表                              │
│ ☐ 添加员工                                  │
│ ☐ 删除员工                                  │
│ ☐ 修改员工权限                              │
├─────────────────────────────────────────────┤
│ [保存权限]  [取消]                          │
└─────────────────────────────────────────────┘
```

---

## 🎟️ 统一邀请码系统

### 邀请类型分类

```
A. 超级管理员邀请码（平台级固定邀请码）
├── ADMIN-MASTER              # 创建超级管理员（固定）
└── SP-ADMIN                  # 创建服务商管理员（固定）

B. 商家管理员邀请码（商家级固定邀请码）
├── MSTAFF-{商家ID后6位}      # 邀请商家员工（固定）
└── MADMIN-TRANS-{操作ID后6位} # 转移商家管理员权限（一次性）

C. 服务商管理员邀请码（服务商级固定邀请码）
├── MERCHANT-{服务商ID后6位}  # 邀请新商家（固定，⚠️ 核心功能）
├── SPSTAFF-{服务商ID后6位}   # 邀请服务商员工（固定）
└── SPADMIN-TRANS-{操作ID后6位} # 转移服务商管理员权限（一次性）

D. 服务商员工邀请码（员工级固定邀请码）
├── CREATOR-{员工ID后6位}     # 邀请达人（固定，每个员工1个）
└── TASK-{任务ID后6位}        # 任务邀请码（固定，兼具达人邀请功能）
```

**⚠️ 重要变更**：
- **移除了** `MERCHANT-ADMIN`（超级管理员不再直接创建商家）
- **新增了** `MERCHANT-{服务商ID后6位}`（服务商邀请新商家）

**✅ 设计原则**：
- **固定邀请码**：每个角色/组织/员工都有固定的邀请码，永久使用
- **方便分享**：固定邀请码可以印在名片、海报、邮件签名上
- **便于追踪**：通过邀请码自动追踪邀请关系，用于返佣计算
- **安全性保证**：通过手机号验证、实名认证、风控系统保证安全
- **明确用途**：每个邀请码都有明确的目标角色，不存在万能码

---

### 邀请码数据结构

```typescript
interface InviteCode {
  id: string;                  // 邀请码ID（UUID）
  code: string;                // 固定邀请码（规则见下方）
  type: InviteType;            // 邀请类型（见上方列表）
  targetRole: string;          // 被邀请人将获得的角色

  // 生成者信息
  generatorId: string;         // 生成者ID
  generatorType: string;       // 生成者类型 (super_admin/merchant_admin/sp_admin/sp_staff/system)

  // 关联组织（可选）
  organizationId?: string;     // 商家ID或服务商ID
  organizationType?: string;   // merchant/service_provider

  // 状态管理
  isActive: boolean;           // 是否启用
  isOneTime: boolean;          // 是否一次性（如权限转移码）
  maxUses?: number;            // 最大使用次数（NULL=无限）
  useCount: number;            // 已使用次数

  // 时间戳
  createdAt: Date;
  updatedAt: Date;
}
```

---

### 固定邀请码生成规则

**设计理念**：每个用户/组织拥有一个固定邀请码，基于ID生成，永久不变。

#### 平台级邀请码（超级管理员）

**格式**：固定字符串（平台初始化时生成）

| 邀请码类型 | 固定码格式 | 目标角色 | 说明 |
|-----------|-----------|---------|------|
| 创建超级管理员 | `ADMIN-MASTER` | `SUPER_ADMIN` | 平台初始化时创建，永久固定 |
| 创建服务商管理员 | `SP-ADMIN` | `SERVICE_PROVIDER_ADMIN` | 平台初始化时创建，永久固定 |

**⚠️ 重要说明**：
- 超级管理员**不能直接创建商家管理员**
- 商家管理员必须由服务商通过邀请码创建（见"服务商级邀请码"）

#### 商家级邀请码（商家管理员）

**格式**：`{前缀}-{商家ID后6位}`

| 邀请码类型 | 格式 | 目标角色 | 生成时机 |
|-----------|------|---------|----------|
| 邀请商家员工 | `MSTAFF-{商家ID后6位}` | `MERCHANT_STAFF` | 商家注册时自动生成，固定不变 |
| 转移商家管理员 | `MADMIN-TRANS-{操作ID后6位}` | `MERCHANT_ADMIN` | 管理员手动生成，一次性 |

**示例**：
- 商家ID: `clj1234567890abcdef0123456789`
- 商家员工邀请码: `MSTAFF-456789`（固定，永久不变）

#### 服务商级邀请码（服务商管理员）

**格式**：`{前缀}-{服务商ID后6位}`

| 邀请码类型 | 格式 | 目标角色 | 生成时机 |
|-----------|------|---------|----------|
| 邀请商家 | `MERCHANT-{服务商ID后6位}` | `MERCHANT_ADMIN` | 服务商注册时自动生成，固定不变 |
| 邀请服务商员工 | `SPSTAFF-{服务商ID后6位}` | `SERVICE_PROVIDER_STAFF` | 服务商注册时自动生成，固定不变 |
| 转移服务商管理员 | `SPADMIN-TRANS-{操作ID后6位}` | `SERVICE_PROVIDER_ADMIN` | 管理员手动生成，一次性 |

**示例**：
- 服务商ID: `clj9876543210fedcba9876543210`
- 商家邀请码: `MERCHANT-987654`（固定，永久不变，用于邀请新商家）
- 服务商员工邀请码: `SPSTAFF-987654`（固定，永久不变）

**⚠️ 重要说明**：
- **商家邀请码（MERCHANT-xxx）**：用于服务商发展商家
  - 新商家使用此邀请码注册后，成为该商家的管理员
  - 商家自动绑定到生成邀请码的服务商
  - 一个服务商可以邀请多个商家

#### 员工级邀请码（服务商员工）

**格式**：`{前缀}-{员工ID后6位}`

| 邀请码类型 | 格式 | 目标角色 | 生成时机 |
|-----------|------|---------|----------|
| 邀请达人 | `CREATOR-{员工ID后6位}` | `CREATOR` | 员工注册时自动生成，固定不变 |
| 任务邀请码 | `TASK-{任务ID后6位}` | `CREATOR` + 绑定任务 | 任务创建时自动生成 |

**示例**：
- 员工ID: `cljabcdefghijklmnopqrstuvwx`
- 达人邀请码: `CREATOR-qrstuv`（固定，永久不变）

#### 任务邀请码

**格式**：`TASK-{任务ID后6位}`

**说明**：
- **统一格式**：所有任务邀请码均使用任务ID后6位生成固定码
- **使用次数限制**：通过数据库字段 `max_uses` 和 `used_count` 控制使用次数，而非改变码的格式

**示例**：
- 任务ID: `task-1234567890abcdef`
- 任务邀请码: `TASK-789012`（基于任务ID的固定码）

---

### 邀请码生成算法

#### **方案：基于ID后6位（推荐）**

```sql
-- 用户表增加邀请码字段
ALTER TABLE users ADD COLUMN invite_code VARCHAR(20) UNIQUE;

-- 商家表增加邀请码字段
ALTER TABLE merchants ADD COLUMN invite_code VARCHAR(20) UNIQUE;

-- 服务商表增加邀请码字段
ALTER TABLE service_providers ADD COLUMN invite_code VARCHAR(20) UNIQUE;

-- 生成函数：基于ID后6位（⚠️ 改进版：确保唯一性）
CREATE OR REPLACE FUNCTION generate_invite_code_from_id(table_id VARCHAR, prefix VARCHAR)
RETURNS VARCHAR(20) AS $$
DECLARE
  new_code VARCHAR(20);
  code_exists BOOLEAN;
  attempt INT := 0;
BEGIN
  -- 尝试生成，最多10次
  WHILE attempt < 10 LOOP
    attempt := attempt + 1;

    -- 第一次尝试：使用ID后6位
    IF attempt = 1 THEN
      new_code := prefix || '-' || UPPER(SUBSTRING(table_id FROM LENGTH(table_id) - 5 FOR 6));
    ELSE
      -- 如果重复，添加随机后缀（3位）
      new_code := prefix || '-' || UPPER(SUBSTRING(table_id FROM LENGTH(table_id) - 5 FOR 6)) || '-' || LPAD(FLOOR(RANDOM() * 1000)::TEXT, 3, '0');
    END IF;

    -- 检查是否已存在
    SELECT EXISTS(
      SELECT 1 FROM invitation_codes WHERE code = new_code
      UNION
      SELECT 1 FROM users WHERE invite_code = new_code
      UNION
      SELECT 1 FROM merchants WHERE invite_code = new_code
      UNION
      SELECT 1 FROM service_providers WHERE invite_code = new_code
    ) INTO code_exists;

    IF NOT code_exists THEN
      RETURN new_code;
    END IF;
  END LOOP;

  -- 如果10次都失败，抛出异常
  RAISE EXCEPTION 'Failed to generate unique invite code after 10 attempts for prefix: %', prefix;
END;
$$ LANGUAGE plpgsql;

-- 使用示例
SELECT generate_invite_code_from_id('clj1234567890abcdef0123456789', 'CREATOR');
-- 返回: CREATOR-456789

SELECT generate_invite_code_from_id('clj9876543210fedcba9876543210', 'MSTAFF');
-- 返回: MSTAFF-987654

-- 触发器：自动生成邀请码
CREATE OR REPLACE FUNCTION auto_generate_invite_code()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_TABLE_NAME = 'users' THEN
    NEW.invite_code := generate_invite_code_from_id(NEW.id, 'CREATOR');
  ELSIF TG_TABLE_NAME = 'merchants' THEN
    NEW.invite_code := generate_invite_code_from_id(NEW.id, 'MSTAFF');
  ELSIF TG_TABLE_NAME = 'service_providers' THEN
    -- ⚠️ 服务商需要生成两个邀请码：
    -- 1. MERCHANT-xxx 用于邀请新商家
    -- 2. SPSTAFF-xxx 用于邀请员工
    -- 这里只存储员工邀请码，商家邀请码需要单独的字段或表
    NEW.invite_code := generate_invite_code_from_id(NEW.id, 'SPSTAFF');
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ⚠️ 商家邀请码生成说明：
-- 服务商注册时，除了生成SPSTAFF-xxx（员工邀请码）外，
-- 还应该生成MERCHANT-xxx（商家邀请码），用于邀请新商家
-- 这可能需要：
-- 1. 在service_providers表中添加merchant_invite_code字段
-- 2. 或者在invitation_codes表中预生成MERCHANT-xxx类型的邀请码
-- 3. 前缀为'MERCHANT'，organization_id=服务商ID，organization_type='service_provider'

-- 创建触发器
CREATE TRIGGER users_invite_code_trigger
  BEFORE INSERT ON users
  FOR EACH ROW
  EXECUTE FUNCTION auto_generate_invite_code();

CREATE TRIGGER merchants_invite_code_trigger
  BEFORE INSERT ON merchants
  FOR EACH ROW
  EXECUTE FUNCTION auto_generate_invite_code();

CREATE TRIGGER service_providers_invite_code_trigger
  BEFORE INSERT ON service_providers
  FOR EACH ROW
  EXECUTE FUNCTION auto_generate_invite_code();
```

**优势**：
- ✅ **真正固定**：每个用户/组织只有一个邀请码，永久不变
- ✅ **SaaS主流**：Stripe、Notion、Fiber等平台都使用类似方案
- ✅ **简洁**：长度适中（12-15字符），如 `CREATOR-456789`
- ✅ **可读**：容易记忆和传播
- ✅ **反查快**：通过ID索引快速查询邀请人
- ✅ **安全性**：不暴露完整ID，只暴露后6位

---

### 邀请码规则表

| 邀请码类型 | 格式示例 | 生成者 | 被邀请人角色 | 邀请码数量 | 生成时机 |
|-----------|---------|-------|------------|----------|----------|
| `ADMIN-MASTER` | `ADMIN-MASTER` | 系统 | `SUPER_ADMIN` | 平台1个 | 平台初始化 |
| `SP-ADMIN` | `SP-ADMIN` | 系统 | `SERVICE_PROVIDER_ADMIN` | 平台1个 | 平台初始化 |
| `MERCHANT-{服务商ID后6位}` | `MERCHANT-987654` | 系统（基于服务商ID） | `MERCHANT_ADMIN` | 每个服务商1个 | 服务商注册时 |
| `MSTAFF-{商家ID后6位}` | `MSTAFF-456789` | 系统（基于商家ID） | `MERCHANT_STAFF` | 每个商家1个 | 商家注册时 |
| `MADMIN-TRANS-{操作ID后6位}` | `MADMIN-TRANS-A1B2C3` | 商家管理员 | `MERCHANT_ADMIN` | 按需生成 | 手动生成（一次性） |
| `SPSTAFF-{服务商ID后6位}` | `SPSTAFF-987654` | 系统（基于服务商ID） | `SERVICE_PROVIDER_STAFF` | 每个服务商1个 | 服务商注册时 |
| `SPADMIN-TRANS-{操作ID后6位}` | `SPADMIN-TRANS-X9Y8Z7` | 服务商管理员 | `SERVICE_PROVIDER_ADMIN` | 按需生成 | 手动生成（一次性） |
| `CREATOR-{员工ID后6位}` | `CREATOR-QRSTUV` | 系统（基于员工ID） | `CREATOR` | 每个员工1个 | 员工注册时 |
| `TASK-{任务ID后6位}` | `TASK-789012` | 系统（基于任务ID） | `CREATOR` + 绑定任务 | 每个任务1个 | 任务创建时 |

**⚠️ 重要变更**：
- **移除了** `MERCHANT-ADMIN`（平台级固定码）
- **新增了** `MERCHANT-{服务商ID后6位}`（服务商级固定码，用于邀请商家）

**邀请码特性**：
- ✅ **固定码**：平台级、组织级、员工级邀请码都是固定码，永久不变
- ✅ **自动生成**：注册时通过触发器自动生成
- ✅ **唯一性**：通过数据库UNIQUE约束保证唯一
- ✅ **可追溯**：通过邀请码查询数据库可知邀请人信息

### 邀请码与角色对应关系

| 邀请码格式示例 | 生成者类型 | 被邀请人角色 | 说明 |
|--------------|-----------|------------|------|
| `ADMIN-MASTER` | 系统 | `SUPER_ADMIN` | 创建超级管理员 |
| `SP-ADMIN` | 系统 | `SERVICE_PROVIDER_ADMIN` | 创建服务商管理员 |
| `MERCHANT-987654` | 系统（基于服务商ID） | `MERCHANT_ADMIN` | 服务商邀请新商家（固定） |
| `MSTAFF-456789` | 系统（基于商家ID） | `MERCHANT_STAFF` | 商家员工邀请码（固定） |
| `MADMIN-TRANS-A1B2` | 商家管理员 | `MERCHANT_ADMIN` | 商家管理员转移码（一次性） |
| `SPSTAFF-987654` | 系统（基于服务商ID） | `SERVICE_PROVIDER_STAFF` | 服务商员工邀请码（固定） |
| `SPADMIN-TRANS-X9Y8` | 服务商管理员 | `SERVICE_PROVIDER_ADMIN` | 服务商管理员转移码（一次性） |
| `CREATOR-QRSTUV` | 系统（基于员工ID） | `CREATOR` | 达人邀请码（固定） |
| `TASK-789012` | 系统（基于任务ID） | `CREATOR` | 任务邀请码（固定，兼具达人邀请） |

---

### 邀请码安全机制

**⚠️ 重要**：固定邀请码不代表不安全，安全通过以下方式保证：

1. **手机号验证**：
   - 注册时必须验证手机号
   - 一个手机号只能注册一个账号

2. **实名认证**：
   - 提现时需要进行实名认证
   - 绑定银行卡信息

3. **风控系统**：
   - 检测异常注册行为（同IP大量注册、同设备多账号）
   - 限制注册频率

4. **邀请码禁用**：
   - 员工离职时，管理员可以禁用其邀请码
   - 管理员可以随时禁用某个邀请码

5. **审批机制**（可选）：
   - 商家/服务商注册时，可能需要审核
   - 审核通过后才能正常使用

---

### 邀请码使用流程

---

### 邀请码使用流程

#### **流程1：超级管理员创建商家/服务商管理员**

```
┌─────────────────────────────────────────────┐
│ 1. 超级管理员使用平台固定邀请码              │
├─────────────────────────────────────────────┤
│  - 创建商家管理员：MERCHANT-ADMIN            │
│  - 创建服务商管理员：SP-ADMIN                │
│  - 超级管理员分享固定码给负责人              │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 2. 负责人注册时填写邀请码                    │
├─────────────────────────────────────────────┤
│  - 输入：MERCHANT-ADMIN 或 SP-ADMIN          │
│  - 系统验证邀请码有效性                      │
│  - 手机号验证                                │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 3. 注册成功                                  │
├─────────────────────────────────────────────┤
│  - 创建账户                                  │
│  - 获得对应管理员角色                        │
│  - 自动创建组织（商家/服务商）               │
│  - 自动生成组织专属邀请码：                  │
│    • 商家 → MSTAFF-456789                   │
│    • 服务商 → SPSTAFF-987654                │
│  - 记录邀请关系                              │
└─────────────────────────────────────────────┘
```

#### **流程2：商家/服务商管理员邀请员工**

```
┌─────────────────────────────────────────────┐
│ 1. 商家/服务商注册时自动生成固定邀请码       │
├─────────────────────────────────────────────┤
│  - 商家：MSTAFF-{商家ID后6位}               │
│  - 服务商：SPSTAFF-{服务商ID后6位}           │
│  - 管理员可在"邀请码管理"页面查看            │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 2. 管理员分享固定邀请码给员工                │
├─────────────────────────────────────────────┤
│  - 复制固定邀请码分享给被邀请人              │
│  - 可生成二维码打印在名片/海报上             │
│  - 邀请码永久有效（除非手动禁用）            │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 3. 被邀请人注册时填写邀请码                  │
├─────────────────────────────────────────────┤
│  - 输入固定邀请码                            │
│  - 系统验证：                                │
│    ✅ 邀请码是否存在                         │
│    ✅ 邀请码是否处于启用状态                 │
│    ✅ 手机号验证（防刷）                     │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 4. 注册成功                                  │
├─────────────────────────────────────────────┤
│  - 创建账户                                  │
│  - 获得对应员工角色                          │
│  - 绑定到商家/服务商                         │
│  - 记录邀请关系                              │
│  - 如果是服务商员工，自动生成达人邀请码：     │
│    CREATOR-{员工ID后6位}                     │
└─────────────────────────────────────────────┘
```

#### **流程3：服务商员工邀请达人**

```
┌─────────────────────────────────────────────┐
│ 1. 服务商员工注册时自动生成固定达人邀请码     │
├─────────────────────────────────────────────┤
│  - 固定码：CREATOR-{员工ID后6位}              │
│  - 员工永久使用此码邀请达人                  │
│  - 可在"个人中心"查看                        │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 2. 员工分享固定邀请码给达人                  │
├─────────────────────────────────────────────┤
│  - 通过微信、QQ等分享给达人                  │
│  - 可生成二维码                              │
│  - 邀请码永久有效                            │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 3. 达人注册时填写邀请码                      │
├─────────────────────────────────────────────┤
│  - 输入：CREATOR-{员工ID后6位}              │
│  - 系统验证邀请码                            │
│  - 手机号验证                                │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 4. 注册成功                                  │
├─────────────────────────────────────────────┤
│  - 创建账户                                  │
│  - 获得 CREATOR 角色                         │
│  - 建立邀请关系（用于返佣统计）              │
│  - 绑定到该员工所属的服务商                  │
└─────────────────────────────────────────────┘
```

#### **流程4：任务邀请码（兼具达人邀请）**

```
┌─────────────────────────────────────────────┐
│ 1. 服务商员工为任务生成临时邀请码            │
├─────────────────────────────────────────────┤
│  - 格式：TASK-{任务ID后6位}                 │
│  - 设置最大使用次数和过期时间                │
│  - 关联任务和员工（用于返佣）                │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 2. 员工分享任务邀请码给达人                  │
├─────────────────────────────────────────────┤
│  - 通过微信、QQ等分享给达人                  │
│  - 或生成二维码                              │
└─────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────┐
│ 3. 达人使用任务邀请码                        │
├─────────────────────────────────────────────┤
│  - 情况A：达人未注册                         │
│    • 注册并获得 CREATOR 角色                 │
│    • 自动绑定任务                            │
│    • 记录邀请关系（返佣）                    │
│  - 情况B：达人已注册                         │
│    • 仅绑定任务                              │
│    • 记录邀请关系（返佣）                    │
└─────────────────────────────────────────────┘
```

---

### 注册登录流程（集成 auth-center）

**⚠️ 重要**：PR Business 通过 **auth-center** (os.crazyaigc.com) 进行统一的用户认证。

#### **认证架构**

```
┌─────────────────────────────────────────────────────────────┐
│                    认证架构                                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  PR Business (pr.crazyaigc.com)                             │
│      │                                                      │
│      ├─ 前端：Vite + React                                  │
│      ├─ 后端：Go API                                        │
│      └─ 数据库：pr_business_db (PostgreSQL)                │
│                                                             │
│  Auth Center (os.crazyaigc.com)                             │
│      │                                                      │
│      ├─ 用户认证：微信登录 + 密码登录                        │
│      ├─ 用户管理：统一用户ID (auth_center_user_id)          │
│      └─ 数据库：auth_center_db (PostgreSQL)                │
│                                                             │
│  关联方式：                                                  │
│  pr_business_db.users.auth_center_user_id                  │
│      └─> auth_center_db.users.user_id                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### **登录方式**

| 登录方式 | 说明 | 使用场景 |
|---------|------|----------|
| **微信登录** | 通过微信授权登录，支持开放平台扫码和公众号授权 | 推荐方式 |
| **密码登录** | 手机号 + 密码登录 | 管理员设置密码后可用 |

#### **两种邀请码填写方式（冗余设计）**

**方式1：扫码自动传递（推荐）**

```
┌─────────────────────────────────────────────────────────────┐
│ 流程：扫码自动填充邀请码                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ 1. 邀请人分享二维码                                          │
│    └─ 生成二维码URL:                                        │
│       https://pr.crazyaigc.com/register?code=MSTAFF-456789 │
│                                                             │
│ 2. 被邀请人扫码                                              │
│    ├─ 微信/浏览器自动跳转到注册页面                          │
│    └─ 前端自动获取URL中的邀请码参数                         │
│                                                             │
│ 3. 前端保存邀请码（双重保险）                                │
│    ├─ sessionStorage.setItem('pendingInviteCode', code)    │
│    └─ document.cookie 设置（兜底方案）                      │
│                                                             │
│ 4. 被邀请人点击"微信登录"                                    │
│    └─ 跳转到微信授权页面                                    │
│                                                             │
│ 5. 微信授权成功，回调到 PR Business                          │
│    └─ 邀请码已保存在前端，不会丢失 ✅                       │
│                                                             │
│ 6. 后端自动应用邀请码                                        │
│    ├─ 前端发送 POST /api/v1/auth/apply-invite-code            │
│    ├─ 后端验证邀请码，分配角色                              │
│    └─ 跳转到对应工作台                                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**方式2：登录后手动填写（兜底）**

```
┌─────────────────────────────────────────────────────────────┐
│ 流程：登录后填写邀请码                                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ 1. 被邀请人直接访问 pr.crazyaigc.com                        │
│    └─ 正常微信登录 / 密码登录                              │
│                                                             │
│ 2. 登录成功，检查用户角色                                    │
│    ├─ 如果有角色 → 跳转到工作台                             │
│    └─ 如果无角色 → 强制跳转到输入邀请码页面                 │
│                                                             │
│ 3. 被邀请人手动输入邀请码                                    │
│    └─ 输入框：MSTAFF-456789                                 │
│                                                             │
│ 4. 提交验证                                                  │
│    ├─ 前端发送 POST /api/v1/auth/apply-invite-code            │
│    ├─ 后端验证邀请码，分配角色                              │
│    └─ 跳转到对应工作台                                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### **邀请码存储机制（双重保险）**

**⚠️ 重要**：由于微信登录流程会经过多个跳转，URL 参数会丢失，因此采用前端存储方案。

```typescript
// 邀请码存储实现（双重保险）
class InviteCodeStorage {
  // 保存邀请码
  save(code: string) {
    // 1. 保存到 sessionStorage（页面会话级别）
    sessionStorage.setItem('pendingInviteCode', code)

    // 2. 保存到 Cookie（更持久，兜底方案）
    document.cookie = `pendingInviteCode=${code}; path=/; max-age=3600`
  }

  // 读取邀请码
  load(): string | null {
    // 1. 优先从 sessionStorage 读取
    const sessionCode = sessionStorage.getItem('pendingInviteCode')
    if (sessionCode) return sessionCode

    // 2. 兜底：从 Cookie 读取
    const cookies = document.cookie.split(';')
    const cookieCode = cookies.find(c => c.trim().startsWith('pendingInviteCode='))
    if (cookieCode) {
      const code = cookieCode.split('=')[1]
      sessionStorage.setItem('pendingInviteCode', code)
      return code
    }

    return null
  }

  // 清除邀请码
  clear() {
    sessionStorage.removeItem('pendingInviteCode')
    document.cookie = 'pendingInviteCode=; path=/; max-age=0'
  }
}
```

**兼容性保证**：

| 浏览器类型 | sessionStorage | Cookie | 邀请码不丢失 |
|-----------|---------------|--------|------------|
| 微信内置浏览器 | ✅ 完全支持 | ✅ 完全支持 | ✅ 100% |
| Safari/Chrome | ✅ 完全支持 | ✅ 完全支持 | ✅ 100% |
| 其他现代浏览器 | ✅ 完全支持 | ✅ 完全支持 | ✅ 100% |

#### **完整注册流程（扫码方式）**

```
1. 用户扫码
   URL: https://pr.crazyaigc.com/register?code=MSTAFF-456789

2. 前端页面加载
   ├─ 获取 URL 参数：code = MSTAFF-456789
   ├─ 保存到 sessionStorage 和 Cookie
   └─ 显示："您已使用邀请码 MSTAFF-456789"

3. 用户点击"微信登录"
   └─ 跳转：/api/v1/auth/wechat/login

4. PR Business 后端
   └─ 重定向到：https://os.crazyaigc.com/api/v1/auth/wechat/login?callbackUrl=...

5. Auth Center
   ├─ 重定向到微信授权页面
   └─ 用户扫码/点击确认

6. 微信回调到 Auth Center
   └─ 处理授权，获取用户信息

7. Auth Center 回调到 PR Business
   URL: /api/v1/auth/callback?userId=xxx&token=xxx

8. PR Business 后端
   ├─ 验证 token，获取用户信息
   ├─ 查找/创建本地用户（通过 auth_center_user_id 关联）
   └─ 设置 session，跳转到前端

9. 前端自动应用邀请码
   ├─ 从 sessionStorage 读取邀请码
   ├─ 发送 POST /api/v1/auth/apply-invite-code
   ├─ 后端验证邀请码，分配角色
   └─ 清除已使用的邀请码

10. 跳转到工作台
    URL: /workspace/merchant
```

#### **页面路由设计**

| 路由 | 说明 | 需要登录 |
|------|------|---------|
| `/register?code=xxx` | 注册页面（带邀请码） | ❌ |
| `/login` | 登录页面 | ❌ |
| `/login/password` | 密码登录页面 | ❌ |
| `/api/v1/auth/wechat/login` | 微信登录（后端） | ❌ |
| `/api/v1/auth/callback` | 登录回调（后端） | ❌ |
| `/api/v1/auth/password/login` | 密码登录（后端） | ❌ |
| `/api/v1/auth/apply-invite-code` | 应用邀请码（后端） | ✅ |
| `/invite-code-required` | 邀请码必填页面 | ✅ |
| `/workspace/*` | 工作台 | ✅ |

---

### 认证 API 接口设计

**⚠️ 重要**：PR Business 通过 auth-center 进行统一认证，本系统的认证API主要是代理接口和邀请码相关接口。

#### **1. 微信登录（代理接口）**

**接口**：`GET /api/v1/auth/wechat/login`

**说明**：代理到 auth-center 发起微信登录

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| callbackUrl | string | 否 | 登录后的回调URL（默认：https://pr.crazyaigc.com/api/v1/auth/callback） |

**流程**：
```
1. 前端请求：GET /api/v1/auth/wechat/login
2. 后端重定向到：https://os.crazyaigc.com/api/v1/auth/wechat/login?callbackUrl=xxx
3. 用户在微信授权页面扫码/点击确认
4. 微信回调到 auth-center
5. auth-center 回调到：https://pr.crazyaigc.com/api/v1/auth/callback?userId=xxx&token=xxx
```

**示例**：
```typescript
// 前端调用
window.location.href = '/api/v1/auth/wechat/login'
```

---

#### **2. 微信登录回调（代理接口）**

**接口**：`GET /api/v1/auth/callback`

**说明**：接收 auth-center 的登录回调，创建/获取本地用户

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| userId | string | 是 | auth-center 的用户ID |
| token | string | 是 | auth-center 的 JWT token |
| unionId | string | 否 | 微信 unionid |
| phoneNumber | string | 否 | 手机号 |

**响应**：
```json
{
  "success": true,
  "data": {
    "userId": "clxxxxx",
    "authCenterUserId": "550e8400-e29b-41d4-a716-446655440000",
    "roles": ["CREATOR"],
    "currentRole": null,
    "hasRole": false
  }
}
```

**流程**：
```go
// 1. 验证 auth-center 的 token
verifyTokenResp := callAuthCenterVerifyToken(token)

// 2. 查找或创建本地用户
user := findOrCreateUserByAuthCenterID(userId)

// 3. 设置 session
setSessionCookie(c, user)

// 4. 检查用户角色
if len(user.Roles) == 0 {
  // 无角色，跳转到邀请码必填页面
  c.Redirect(302, "/invite-code-required")
} else {
  // 有角色，跳转到工作台
  c.Redirect(302, "/workspace")
}
```

---

#### **3. 密码登录（代理接口）**

**接口**：`POST /api/v1/auth/password/login`

**说明**：代理到 auth-center 进行密码登录

**请求体**：
```json
{
  "phoneNumber": "13800138000",
  "password": "password123"
}
```

**响应**：
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "userId": "clxxxxx",
  "roles": ["MERCHANT_ADMIN"],
  "currentRole": "MERCHANT_ADMIN"
}
```

**流程**：
```go
// 1. 调用 auth-center 密码登录API
authCenterResp := callAuthCenterPasswordLogin(phoneNumber, password)

// 2. 查找或创建本地用户
user := findOrCreateUserByAuthCenterID(authCenterResp.UserID)

// 3. 设置 session
setSessionCookie(c, user)

// 4. 返回用户信息
c.JSON(200, gin.H{
  "success": true,
  "userId": user.ID,
  "roles": user.Roles,
  "currentRole": user.CurrentRole,
})
```

---

#### **4. 应用邀请码**

**接口**：`POST /api/v1/auth/apply-invite-code`

**说明**：登录后应用邀请码，分配角色

**权限**：需要登录

**请求体**：
```json
{
  "inviteCode": "MSTAFF-456789"
}
```

**响应**：
```json
{
  "success": true,
  "data": {
    "currentRole": "MERCHANT_STAFF",
    "roles": ["MERCHANT_STAFF", "CREATOR"]
  }
}
```

**流程**：
```go
// 1. 获取当前用户
user := getCurrentUser(c)

// 2. 验证邀请码
invitationCode := validateInviteCode(req.InviteCode)
if !invitationCode.IsActive {
  return c.JSON(400, gin.H{"error": "邀请码无效或已过期"})
}

// 3. 检查用户是否已有该角色
if contains(user.Roles, invitationCode.TargetRole) {
  return c.JSON(400, gin.H{"error": "您已拥有该角色"})
}

// 4. 应用邀请码，添加角色
user.Roles = append(user.Roles, invitationCode.TargetRole)
user.CurrentRole = invitationCode.TargetRole
user.InvitationCodeID = &invitationCode.ID
db.Save(&user)

// 5. 记录邀请关系
record := InvitationRecord{
  InvitationCodeID: invitationCode.ID,
  InviteeID:        &user.ID,
  InviteeRole:      invitationCode.TargetRole,
  OrganizationID:   invitationCode.OrganizationID,
  Status:           "completed",
  UsedAt:           time.Now(),
}
db.Create(&record)

// 6. 更新邀请码使用次数
invitationCode.UseCount++
db.Save(&invitationCode)

c.JSON(200, gin.H{
  "success": true,
  "data": gin.H{
    "currentRole": user.CurrentRole,
    "roles":       user.Roles,
  },
})
```

---

#### **5. 获取当前用户信息**

**接口**：`GET /api/v1/auth/me`

**权限**：需要登录

**响应**：
```json
{
  "success": true,
  "data": {
    "id": "clxxxxx",
    "authCenterUserId": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "张三",
    "avatarUrl": "https://wx.qlogo.cn/xxx",
    "roles": ["SERVICE_PROVIDER_STAFF", "CREATOR"],
    "currentRole": "SERVICE_PROVIDER_STAFF",
    "invitedBy": "clyyyyy",
    "status": "active"
  }
}
```

---

#### **6. 登出**

**接口**：`POST /api/v1/auth/logout`

**权限**：需要登录

**响应**：
```json
{
  "success": true,
  "message": "登出成功"
}
```

**流程**：
```go
// 1. 清除 session cookie
clearSessionCookie(c)

// 2. 可选：调用 auth-center 登出接口
// callAuthCenterLogout(token)

c.JSON(200, gin.H{"success": true})
```

---

#### **7. 切换角色**

**接口**：`POST /api/v1/auth/switch-role`

**权限**：需要登录

**请求体**：
```json
{
  "role": "CREATOR"
}
```

**响应**：
```json
{
  "success": true,
  "data": {
    "currentRole": "CREATOR"
  }
}
```

**流程**：
```go
// 1. 获取当前用户
user := getCurrentUser(c)

// 2. 检查用户是否拥有该角色
if !contains(user.Roles, req.Role) {
  return c.JSON(403, gin.H{"error": "您没有该角色"})
}

// 3. 切换角色
user.CurrentRole = req.Role
db.Save(&user)

c.JSON(200, gin.H{
  "success": true,
  "data": gin.H{
    "currentRole": user.CurrentRole,
  },
})
```

---

### API接口规范

#### API版本管理

**URL格式**：
```
/api/v1/{resource}
```

**示例**：
- `/api/v1/auth/wechat/login` - 微信登录
- `/api/v1/auth/password/login` - 密码登录
- `/api/v1/auth/apply-invite-code` - 应用邀请码
- `/api/v1/tasks` - 获取任务列表
- `/api/v1/tasks/{id}` - 获取任务详情

**版本管理策略**：
- 当前版本：v1
- 向后兼容：保证至少3个主版本的兼容性
- 废弃版本：提前6个月通知，标记为deprecated
- 新版本测试：使用 `/api/v2/` 作为测试环境

**⚠️ API版本统一要求**：
- **所有接口必须包含版本号前缀**（包括认证接口）
- **错误示例**：`/api/auth/login` ❌
- **正确示例**：`/api/v1/auth/login` ✅

#### 统一响应格式

**成功响应**：
```json
{
  "success": true,
  "data": { ... },
  "message": "操作成功"  // 可选
}
```

**错误响应**：
```json
{
  "success": false,
  "error": {
    "code": "INVITE_CODE_INVALID",
    "message": "邀请码无效或已过期",
    "details": {
      "field": "inviteCode",
      "value": "INVALID_CODE",
      "timestamp": "2026-02-02T10:30:00Z"
    }
  },
  "request_id": "req_123456789"
}
```

**错误码规范**：

| 错误码 | 说明 | HTTP状态码 | 示例场景 |
|-------|------|-----------|---------|
| `UNAUTHORIZED` | 未登录 | 401 | token无效或过期 |
| `FORBIDDEN` | 无权限 | 403 | 无权限访问资源 |
| `INVALID_PARAMS` | 参数错误 | 400 | 必填参数缺失或格式错误 |
| `INVITE_CODE_INVALID` | 邀请码无效 | 400 | 邀请码不存在或已过期 |
| `INVITE_CODE_USED` | 邀请码已使用 | 400 | 一次性邀请码已使用 |
| `INSUFFICIENT_BALANCE` | 余额不足 | 400 | 账户余额不足以支付 |
| `TASK_ALREADY_ASSIGNED` | 任务已接 | 400 | 任务已被其他人接单 |
| `TASK_NOT_FOUND` | 任务不存在 | 404 | 任务ID无效 |
| `RATE_LIMIT_EXCEEDED` | 请求过于频繁 | 429 | 触发限流 |
| `INTERNAL_ERROR` | 服务器错误 | 500 | 服务器内部错误 |

**⚠️ 错误处理要求**：
- 所有API错误必须返回统一格式
- 必须包含request_id用于日志追踪
- 敏感信息不暴露在错误消息中（如手机号、密码）
- 前端根据错误码显示对应的友好提示

---

### 邀请码管理页面（固定邀请码模式）

#### **超级管理员 → 平台管理**

```
┌─────────────────────────────────────────────────────┐
│ 平台邀请码管理                                       │
├─────────────────────────────────────────────────────┤
│                                                     │
│ 创建超级管理员                                       │
│ ┌─────────────────────────────────────────────────┐ │
│ │ ADMIN-MASTER                    [复制] [二维码]  │ │
│ │ 已使用：3次  │  最后使用：2026-02-01            │ │
│ └─────────────────────────────────────────────────┘ │
│                                                     │
│ 创建商家管理员                                       │
│ ┌─────────────────────────────────────────────────┐ │
│ │ MERCHANT-ADMIN                  [复制] [二维码]  │ │
│ │ 已使用：15次  │  最后使用：2026-02-02            │ │
│ └─────────────────────────────────────────────────┘ │
│                                                     │
│ 创建服务商管理员                                     │
│ ┌─────────────────────────────────────────────────┐ │
│ │ SP-ADMIN                        [复制] [二维码]  │ │
│ │ 已使用：8次  │  最后使用：2026-02-01            │ │
│ └─────────────────────────────────────────────────┘ │
│                                                     │
│ [查看使用记录]                                      │
└─────────────────────────────────────────────────────┘
```

#### **商家/服务商管理员 → 邀请码管理**

**页面路由**：
- 商家：`/workspace/merchant/invite-codes`
- 服务商：`/workspace/service-provider/invite-codes`

```
┌─────────────────────────────────────────────────────┐
│ 我的邀请码                                           │
├─────────────────────────────────────────────────────┤
│                                                     │
│ 员工邀请码                                           │
│ ┌─────────────────────────────────────────────────┐ │
│ │ MSTAFF-456789                   [复制] [二维码]  │ │
│ │ 已使用：5次  │  最后使用：2026-02-01            │ │
│ └─────────────────────────────────────────────────┘ │
│                                                     │
│ [查看使用记录]                                      │
│                                                     │
│ ──────────────────────────────────────────────────│
│                                                     │
│ 管理员权限转移                                       │
│ ┌─────────────────────────────────────────────────┐ │
│ │ [生成转移码]                                     │ │
│ └─────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

#### **服务商员工 → 个人中心**

```
┌─────────────────────────────────────────────────────┐
│ 我的达人邀请码                                       │
├─────────────────────────────────────────────────────┤
│                                                     │
│ ┌─────────────────────────────────────────────────┐ │
│ │ CREATOR-XYZ789                  [复制] [二维码]  │ │
│ └─────────────────────────────────────────────────┘ │
│                                                     │
│ 已邀请达人：15人                                     │
│ 累计返佣：￥1,230                                    │
│                                                     │
│ [查看邀请达人列表]                                   │
└─────────────────────────────────────────────────────┘
```

---

### 邀请码使用记录

**记录详情**：

| 使用时间 | 被邀请人 | 手机号 | 获得角色 | 状态 |
|---------|--------|------|----------|------|
| 2026-02-01 10:23 | 张三 | 138****1234 | 商家员工 | ✅ 已激活 |
| 2026-02-02 14:30 | 李四 | 139****5678 | 商家员工 | ✅ 已激活 |
| 2026-02-03 09:15 | 王五 | 136****9012 | 达人 | ✅ 已激活 |
| ... | ... | ... | ... | ... |

---

## 📚 数据字典

### 数据字典说明

本文档包含两类数据字典：
1. **术语表**：业务术语的定义、说明和示例
2. **字段说明**：数据库表字段的详细说明（类型、约束、默认值等）

---

### 一、术语表

#### 1. 核心业务术语

| 中文术语 | 英文术语 | 数据库表/字段 | 说明 | 示例 |
|---------|---------|---------------|------|------|
| **营销活动** | Campaign | `campaigns` 表 | 商家发布的活动 | "新品推广活动" |
| **活动ID** | Campaign ID | `campaigns.id` | 活动唯一标识 | UUID |
| **活动名称** | Title | `campaigns.title` | 活动标题 | "新品体验推广" |
| **活动状态** | Status | `campaigns.status` | DRAFT/OPEN/CLOSED | "开放中" |
| **任务金额** | Task Amount | `campaigns.task_amount` | 每个名额的积分费用 | 100积分 |
| **活动总费用** | Campaign Amount | `campaigns.campaign_amount` | task_amount × quota | 1000积分 |
| **活动名额** | Quota | `campaigns.quota` | 需要的达人总数 | 10人 |
| **已接名额数** | Enrolled Count | `COUNT(tasks.id)` WHERE status != 'OPEN' | 已接的名额数（查询得出） | 3人 |
| **剩余名额** | Remaining Slots | `COUNT(tasks.id)` WHERE status = 'OPEN' | 还可以接的名额数 | 7人 |
| **活动要求** | Requirements | `campaigns.requirements` | 活动要求详情 | 富文本 |
| **平台类型** | Platforms | `campaigns.platforms` | 任务平台列表 | ["xiaohongshu"] |
| **达人收入** | Creator Amount | `campaigns.creator_amount` | 每个任务达人获得 | 80积分 |
| **员工返佣** | Staff Referral Amount | `campaigns.staff_referral_amount` | 每个任务员工返佣 | 10积分 |
| **服务商收入** | Provider Amount | `campaigns.provider_amount` | 每个任务服务商收入 | 10积分 |
| **接任务截止时间** | Task Deadline | `campaigns.task_deadline` | 达人接任务的最后期限 | "2024-01-25 18:00" |
| **提交截止时间** | Submission Deadline | `campaigns.submission_deadline` | 已接达人提交的最后期限 | "2024-02-01 18:00" |
| **任务** | Task | `tasks` 表 | 营销活动的单个名额（从达人视角看是"一个任务"） | "赏金100积分的任务" |
| **任务名额** | Task Slot | `tasks` 表 | 营销活动的一个名额（预分配模式） | "#1号名额" |
| **名额ID** | Task ID | `tasks.id` | 名额唯一标识 | UUID |
| **名额状态** | Task Status | `tasks.status` | OPEN/ASSIGNED/SUBMITTED/APPROVED/REJECTED | "开放中" |
| **达人ID** | Creator ID | `tasks.creator_id` | 接任务的达人ID | 关联 creators.id |
| **平台链接** | Platform URL | `tasks.platform_url` | 达人提交的链接 | "https://..." |
| **截图凭证** | Screenshots | `tasks.screenshots` | 达人提交的截图 | JSON数组 |
| **积分** | Credits | （通用概念） | 平台虚拟货币，1元=1积分 | "充值100元=100积分" |
| **可用余额** | Balance | `credit_accounts.balance` | 可随时使用的积分 | 500积分 |
| **冻结余额** | Frozen Balance | `credit_accounts.frozen_balance` | 暂时冻结的积分 | 200积分 |
| **充值** | Recharge | `RECHARGE` 交易类型 | 商家充值 | +1000积分 |
| **任务收入** | Task Income | `TASK_INCOME` 交易类型 | 达人完成任务收入 | +80积分 |
| **提现** | Withdraw | `WITHDRAW` 交易类型 | 用户提现 | -100积分 |
| **退款** | Refund | `TASK_REFUND` 交易类型 | 任务关闭退款 | +200积分 |
| **员工返佣** | Staff Referral | `STAFF_REFERRAL` 交易类型 | 员工邀请返佣 | +10积分 |
| **积分流水** | Transaction | `credit_transactions` 表 | 所有积分变动记录 | - |
| **提现申请** | Withdrawal | `withdrawals` 表 | 提现申请记录 | - |

---

### 角色术语

| 中文术语 | 英文代码 | 数据库表 | 说明 | 工作台路由 |
|---------|---------|---------|------|----------|
| **超级管理员** | SUPER_ADMIN | `users` 表（roles字段） | 平台运营，管理所有功能 | `/workspace/admin` |
| **商家管理员** | MERCHANT_ADMIN | `merchants.admin_id` | 商家负责人，唯一管理员 | `/workspace/merchant` |
| **商家员工** | MERCHANT_STAFF | `merchant_staff` 表 | 商家员工，权限可配置 | `/workspace/merchant` |
| **服务商管理员** | SERVICE_PROVIDER_ADMIN | `service_providers.admin_id` | 服务商负责人，唯一管理员 | `/workspace/service-provider` |
| **服务商员工** | SERVICE_PROVIDER_STAFF | `service_provider_staff` 表 | 服务商员工，权限可配置 | `/workspace/service-provider` |
| **达人** | CREATOR | `creators` 表 | 接任务、执行任务、获得收入 | `/workspace/creator` |

---

### 营销活动状态术语

| 中文术语 | 状态代码 | 数据库值 | 说明 | 下一个状态 |
|---------|---------|---------|------|----------|
| **草稿** | DRAFT | `'DRAFT'` | 商家创建的草稿活动 | OPEN |
| **开放中** | OPEN | `'OPEN'` | 已发布，等待达人接任务 | CLOSED |
| **已关闭** | CLOSED | `'CLOSED'` | 名额已满或手动关闭（截止时间到不自动关闭） | - |

**⚠️ 重要**：
- 活动有**2个截止时间**，可独立延长
- **接任务截止时间到**：活动标记为"已到期"，但**不会自动关闭**，服务商可延长截止时间继续招募
- **提交截止时间到**：所有未提交的`ASSIGNED`状态名额自动拒绝
- 只有**手动关闭**或**名额已满**时才退款

---

### 任务名额状态术语

| 中文术语 | 状态代码 | 数据库值 | 说明 | 下一个状态 |
|---------|---------|---------|------|----------|
| **进行中** | ASSIGNED | `'ASSIGNED'` | 达人已接任务，等待提交 | SUBMITTED |
| **待审核** | SUBMITTED | `'SUBMITTED'` | 达人已提交，等待审核 | APPROVED/REJECTED |
| **已完成** | APPROVED | `'APPROVED'` | 审核通过，积分已分配 | - |
| **已拒绝** | REJECTED | `'REJECTED'` | 审核拒绝，可重新提交 | SUBMITTED |

---

### 达人等级术语

| 中文术语 | 英文代码 | 数据库字段 | 说明 | 粉丝数要求 |
|---------|---------|-----------|------|-----------|
| **普通用户** | UGC | `creators.level = 'UGC'` | 普通用户 | 无要求 |
| **关键意见消费者** | KOC | `creators.level = 'KOC'` | 关键意见消费者 | 1K-10K |
| **达人** | INF | `creators.level = 'INF'` | 达人 | 10K-100K |
| **关键意见领袖** | KOL | `creators.level = 'KOL'` | 关键意见领袖 | 100K+ |

---

### 邀请码类型术语（固定邀请码模式）

| 中文术语 | 邀请码格式 | 被邀请人角色 | 说明 | 使用场景 |
|---------|-----------|------------|------|----------|
| **创建超级管理员** | `ADMIN-MASTER` | `SUPER_ADMIN` | 平台级固定码 | 超级管理员邀请 |
| **创建服务商管理员** | `SP-ADMIN` | `SERVICE_PROVIDER_ADMIN` | 平台级固定码 | 超级管理员邀请 |
| **邀请新商家** | `MERCHANT-{服务商ID后6位}` | `MERCHANT_ADMIN` | 服务商级固定码 | ⚠️ 服务商发展商家 |
| **邀请商家员工** | `MSTAFF-{商家ID后6位}` | `MERCHANT_STAFF` | 商家级固定码 | 商家管理员邀请 |
| **转移商家管理员** | `MADMIN-TRANS-{操作ID后6位}` | `MERCHANT_ADMIN` | 一次性码 | 商家管理员转移权限 |
| **邀请服务商员工** | `SPSTAFF-{服务商ID后6位}` | `SERVICE_PROVIDER_STAFF` | 服务商级固定码 | 服务商管理员邀请 |
| **转移服务商管理员** | `SPADMIN-TRANS-{操作ID后6位}` | `SERVICE_PROVIDER_ADMIN` | 一次性码 | 服务商管理员转移权限 |
| **邀请达人** | `CREATOR-{员工ID后6位}` | `CREATOR` | 员工级固定码 | 服务商员工邀请 |
| **任务邀请** | `TASK-{任务ID后6位}` | `CREATOR` + 绑定任务 | 固定码，兼具达人邀请 | 任务创建时自动生成 |

---

### 积分交易类型术语（含字段）

| 中文术语 | 类型代码 | 适用账户类型 | 金额方向 | 说明 |
|---------|---------|-------------|---------|------|
| **充值** | RECHARGE | ORG_MERCHANT | +（正数） | 商家充值积分 |
| **任务收入** | TASK_INCOME | USER_PERSONAL | +（正数） | 达人任务审核通过，frozen_balance→balance |
| **任务提交** | TASK_SUBMIT | USER_PERSONAL | +（正数） | 达人提交任务，收入冻结到frozen_balance |
| **员工返佣** | STAFF_REFERRAL | USER_PERSONAL | +（正数） | 员工邀请的达人完成任务时的返佣 |
| **服务商收入** | PROVIDER_INCOME | ORG_PROVIDER | +（正数） | 服务商完成任务审核通过获得的收入 |
| **发布活动** | TASK_PUBLISH | ORG_MERCHANT | -（负数） | 商家发布活动，balance→frozen_balance |
| **接任务扣除** | TASK_ACCEPT | ORG_MERCHANT | -（负数） | 达人接任务时从商家frozen_balance扣除 |
| **审核拒绝** | TASK_REJECT | USER_PERSONAL | -（负数） | 达人任务被拒绝，释放frozen_balance |
| **超时拒绝** | TASK_ESCALATE | ORG_MERCHANT | +（正数） | 提交截止时间到，自动拒绝，费用退回商家 |
| **活动退款** | TASK_REFUND | ORG_MERCHANT | +（正数） | 活动关闭时，未完成名额费用退回balance |
| **提现** | WITHDRAW | 所有账户类型 | -（负数） | 用户申请提现，扣除balance |
| **提现冻结** | WITHDRAW_FREEZE | 所有账户类型 | -（负数） | 提现审核通过，balance→frozen_balance |
| **提现拒绝** | WITHDRAW_REFUND | 所有账户类型 | +（正数） | 提现被拒绝，frozen_balance→balance |
| **系统赠送** | BONUS_GIFT | USER_PERSONAL | +（正数） | 平台赠送的积分（如新用户奖励） |

---

### 权限术语

> **说明**：权限分为两类：
> - **固定权限**：角色天生拥有，不可修改（超级管理员、商家管理员、服务商管理员、达人）
> - **可配置权限**：由管理员分配给员工（商家员工、服务商员工）

---

#### SUPER_ADMIN - 超级管理员（固定权限）

| 权限代码 | 权限名称 | 说明 |
|---------|---------|------|
| - | **所有系统权限** | 可以查看和管理平台所有数据 |

**⚠️ 权限边界和安全约束**：
1. **数据查看权限**：
   - ✅ 可以查看所有账户的余额和交易记录
   - ✅ 可以查看所有商家、服务商、达人的数据
   - ✅ 可以查看系统配置和日志

2. **数据操作权限**：
   - ✅ 可以调整任何账户的balance和frozen_balance（需要填写原因和记录操作日志）
   - ✅ 可以手动调整达人等级（用于特殊合作）
   - ✅ 可以禁用/启用任何邀请码
   - ✅ 可以封禁/解封任何用户
   - ❌ 不能直接删除账户（只能软删除，数据保留180天）
   - ❌ 不能修改自己的角色（防止误操作）

3. **操作审计要求**：
   - 所有敏感操作必须记录：操作人、操作时间、操作类型、影响对象、原因
   - 敏感操作包括：余额调整、角色变更、账户封禁、数据删除
   - 审计日志永久保存，不可删除
   - 超级管理员之间的操作也需要审计

4. **二次确认机制**：
   - 调整账户余额需要二次确认（输入密码或验证码）
   - 封禁用户需要填写原因并等待30秒冷却期
   - 删除数据需要"超级管理员+技术总监"双重授权（⚠️ 建议但不强制）

---

#### MERCHANT_ADMIN - 商家管理员（固定权限）

| 权限代码 | 权限名称 | 说明 |
|---------|---------|------|
| `task.view` | 查看任务列表 | 可以查看商家的所有任务 |
| `task.create` | 创建任务草稿 | 可以创建任务草稿 |
| `task.edit` | 编辑任务草稿 | 可以编辑草稿状态的任务 |
| `task.delete` | 删除任务 | 可以删除草稿或已关闭的任务 |
| `task.publish` | 发布任务 | 提交给服务商设置积分分配（实际发布由服务商操作） |
| `finance.view` | 查看财务数据 | 可以查看余额、收支统计 |
| `finance.recharge` | 发起充值 | 可以发起充值申请 |
| `finance.transaction.view` | 查看积分流水 | 可以查看详细的积分收支记录 |
| `staff.view` | 查看员工列表 | 可以查看商家员工列表 |
| `staff.add` | 添加员工 | 可以邀请新的员工 |
| `staff.remove` | 删除员工 | 可以删除员工 |
| `staff.permission.edit` | 修改员工权限 | 可以分配/修改员工权限 |
| `admin.transfer` | 管理员权限转移 | 可以将管理员权限转移给其他员工 |

---

#### MERCHANT_STAFF - 商家员工（可配置权限）

> ⚠️ **重要**：这些权限由商家管理员通过勾选分配给员工

| 权限代码 | 权限名称 | 说明 |
|---------|---------|------|
| `task.view` | 查看任务列表 | 可以查看商家的所有任务 |
| `task.create` | 创建任务草稿 | 可以创建任务草稿 |
| `task.edit` | 编辑任务草稿 | 可以编辑草稿状态的任务 |
| `finance.view` | 查看财务数据 | 可以查看余额、收支统计 |
| `finance.recharge` | 发起充值 | 可以发起充值申请 |
| `finance.transaction.view` | 查看积分流水 | 可以查看详细的积分收支记录 |
| `staff.view` | 查看员工列表 | 可以查看商家员工列表 |
| `staff.add` | 添加员工 | 可以邀请新的员工 |
| `staff.remove` | 删除员工 | 可以删除员工 |
| `staff.permission.edit` | 修改员工权限 | 可以分配/修改员工权限 |

---

#### SERVICE_PROVIDER_ADMIN - 服务商管理员（固定权限）

| 权限代码 | 权限名称 | 说明 |
|---------|---------|------|
| `task.view` | 查看任务列表 | 可以查看绑定商家的任务 |
| `task.allocate` | 设置积分分配 | 可以为草稿活动设置积分并发布 |
| `task.audit` | 审核任务 | 可以审核达人提交的任务 |
| `creator.invite` | 邀请达人 | 可以生成邀请码邀请达人注册 |
| `creator.view` | 查看邀请达人 | 可以查看自己邀请的达人列表 |
| `merchant.bind` | 绑定商家 | 可以绑定新的商家 |
| `merchant.view` | 查看绑定商家 | 可以查看绑定的商家列表 |
| `task.invite` | 任务邀请达人 | 可以生成任务邀请码发送给达人 |
| `finance.view` | 查看财务数据 | 可以查看余额、收支统计 |
| `finance.withdraw` | 发起提现 | 可以发起提现申请 |
| `finance.transaction.view` | 查看积分流水 | 可以查看详细的积分收支记录 |
| `staff.view` | 查看员工列表 | 可以查看服务商员工列表 |
| `staff.add` | 添加员工 | 可以邀请新的员工 |
| `staff.remove` | 删除员工 | 可以删除员工 |
| `staff.permission.edit` | 修改员工权限 | 可以分配/修改员工权限 |

---

#### SERVICE_PROVIDER_STAFF - 服务商员工（可配置权限）

> ⚠️ **重要**：这些权限由服务商管理员通过勾选分配给员工

| 权限代码 | 权限名称 | 说明 |
|---------|---------|------|
| `task.view` | 查看任务列表 | 可以查看绑定商家的任务 |
| `task.allocate` | 设置积分分配 | 可以为草稿活动设置积分并发布 |
| `task.audit` | 审核任务 | 可以审核达人提交的任务 |
| `creator.invite` | 邀请达人 | 可以生成邀请码邀请达人注册 |
| `creator.view` | 查看邀请达人 | 可以查看自己邀请的达人列表 |
| `merchant.bind` | 绑定商家 | 可以绑定新的商家 |
| `merchant.view` | 查看绑定商家 | 可以查看绑定的商家列表 |
| `task.invite` | 任务邀请达人 | 可以生成任务邀请码发送给达人 |
| `finance.view` | 查看财务数据 | 可以查看余额、收支统计 |
| `finance.withdraw` | 发起提现 | 可以发起提现申请 |
| `finance.transaction.view` | 查看积分流水 | 可以查看详细的积分收支记录 |
| `staff.view` | 查看员工列表 | 可以查看服务商员工列表 |
| `staff.add` | 添加员工 | 可以邀请新的员工 |
| `staff.remove` | 删除员工 | 可以删除员工 |
| `staff.permission.edit` | 修改员工权限 | 可以分配/修改员工权限 |

---

#### CREATOR - 达人（固定权限）

| 权限代码 | 权限名称 | 说明 |
|---------|---------|------|
| `task.browse` | 浏览任务大厅 | 可以浏览所有可接任务的任务 |
| `task.accept` | 接任务 | 可以接受任务（占用一个名额） |
| `task.submit` | 提交任务 | 可以提交平台链接/截图 |
| `task.view` | 查看我的任务 | 可以查看所有已接任务 |
| `finance.view` | 查看财务数据 | 可以查看余额、收入统计 |
| `finance.withdraw` | 发起提现 | 可以发起提现申请 |
| `finance.transaction.view` | 查看积分流水 | 可以查看详细的积分收支记录 |
| `profile.edit` | 编辑个人资料 | 可以修改自己的达人信息 |

---

### 数据库表术语

| 术语 | 表名 | 说明 | 关联 |
|-----|------|------|------|
| **用户表** | users | 存储所有用户的基本信息，支持多角色 | - |
| **商家表** | merchants | 存储商家信息 | → users |
| **服务商表** | service_providers | 存储服务商信息 | → users |
| **达人表** | creators | 存储达人信息 | → users |
| **商家员工表** | merchant_staff | 存储商家员工信息 | → users, → merchants |
| **服务商员工表** | service_provider_staff | 存储服务商员工信息 | → users, → service_providers |
| **营销活动表** | campaigns | 存储营销活动信息 | → merchants |
| **任务名额表** | tasks | 存储所有任务名额（预分配模式） | → campaigns, → creators |
| **邀请码表** | invite_codes | 存储所有类型的邀请码（支持员工专属） | → users, → campaigns, → service_provider_staff |
| **邀请码使用记录表** | invite_code_usages | 记录邀请码的使用情况 | → invite_codes, → users, → creators |
| **商家员工权限表** | merchant_staff_permissions | 存储商家员工的权限 | → merchant_staff |
| **服务商员工权限表** | provider_staff_permissions | 存储服务商员工的权限 | → service_provider_staff |
| **积分账户表** | credit_accounts | 存储各角色和组织实体的积分账户 | → merchants, → service_providers, → users |
| **积分流水表** | credit_transactions | 记录所有积分变动 | → credit_accounts |
| **提现表** | withdrawals | 存储提现申请 | → credit_accounts |
| **商家-服务商绑定表** | merchant_provider_bindings | 记录商家和服务商的绑定关系 | → merchants, → service_providers |

---

### 界面元素术语

| 术语 | 说明 | 使用位置 |
|-----|------|---------|
| **角色切换器** | 顶部导航栏左侧的下拉菜单，用于切换当前激活的角色（当用户拥有多个角色时显示） | 所有工作台的顶部导航（条件显示） |
| **积分余额** | 当前用户的积分余额，显示在顶部导航 | 所有工作台的顶部导航 |
| **任务大厅** | 浏览所有可接的"任务"的页面 | 达人工作台 |
| **我的任务** | 查看我所有的"活动名额"的页面 | 达人工作台 |
| **活动名额** | 活动招募的人数限制 | 活动详情页 |
| **已接名额数** | 已经接的名额数 | 活动详情页 |
| **剩余名额** | 还可以接的名额数 | 活动详情页 |
| **职能名称** | 管理员给员工设置的职能标识，如"运营"、"财务"等 | 员工管理页面 |

---

### API路径术语

| 术语 | 路径 | 说明 |
|-----|------|------|
| **创建活动草稿** | POST /api/campaigns | 商家创建营销活动草稿 |
| **设置积分分配并发布** | POST /api/campaigns/:id/allocate | 服务商为草稿活动设置积分分配并发布 |
| **任务大厅** | GET /api/creator/tasks | 达人浏览所有可接的任务（OPEN状态的名额） |
| **接任务** | POST /api/campaigns/:id/tasks/:taskId/accept | 达人接任务（通过邀请码，将名额从OPEN改为ASSIGNED） |
| **我的任务** | GET /api/creator/my-tasks | 达人查看所有已接任务 |
| **提交任务** | POST /api/tasks/:id/submit | 达人提交平台链接/截图 |
| **审核任务名额** | POST /api/tasks/:id/review | 服务商审核任务名额 |
| **邀请码生成** | POST /api/invite-codes | 生成各种类型的邀请码（支持员工专属） |

---

### 二、字段说明

#### 1. 通用字段类型说明

| 数据库类型 | 说明 | 示例 | 默认值 |
|-----------|------|------|--------|
| `UUID` | 主键唯一标识，使用 PostgreSQL `gen_random_uuid()` | `'550e8400-e29b-41d4-a716-446655440000'` | 自动生成 |
| `VARCHAR(N)` | 变长字符串，N为最大长度 | `'小红书'` | `NULL` |
| `TEXT` | 不限长度的文本 | `'详细描述...'` | `NULL` |
| `INTEGER` | 整数 | `100` | `0` 或 `NULL` |
| `BIGINT` | 大整数 | `9223372036854775807` | `0` 或 `NULL` |
| `DECIMAL(M,N)` | 精确小数，M总位数，N小数位数 | `100.00` | `0.00` |
| `BOOLEAN` | 布尔值 | `true` / `false` | `false` |
| `TIMESTAMP` | 时间戳（时区：Asia/Shanghai） | `'2024-01-01 12:00:00'` | `CURRENT_TIMESTAMP` |
| `JSONB` | 二进制JSON，支持索引查询 | `{"key": "value"}` | `'{}'` |
| `ENUM` | 枚举类型，限定值域 | `'OPEN'` | 见具体字段 |

#### 2. 通用系统字段

| 字段名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `id` | UUID | ✅ | 自动生成 | 主键唯一标识 |
| `created_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | 创建时间 |
| `updated_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | 更新时间（自动触发更新） |
| `deleted_at` | TIMESTAMP | ❌ | `NULL` | 软删除时间（`NULL`=未删除） |
| `version` | INTEGER | ❌ | `0` | 乐观锁版本号（防止并发冲突） |

**索引策略**：
- 所有表的主键 `id` 自动创建索引
- 所有表的 `deleted_at` 添加部分索引 `WHERE deleted_at IS NULL`（软删除表）
- 高频查询字段（如 `status`, `created_at`）添加普通索引
- 外键字段自动创建索引（如 `user_id`, `campaign_id`）

#### 3. 核心表字段说明

##### 3.1 用户表（users）

| 字段名 | 类型 | 必填 | 默认值 | 约束 | 说明 |
|--------|------|------|--------|------|------|
| `id` | UUID | ✅ | 自动生成 | PRIMARY KEY | 用户唯一标识 |
| `phone` | VARCHAR(20) | ❌ | `NULL` | UNIQUE | 手机号（通过 auth-center 统一登录） |
| `wechat_openid` | VARCHAR(100) | ❌ | `NULL` | UNIQUE | 微信 OpenID（通过 auth-center 统一登录） |
| `roles` | JSONB | ✅ | `'[]'` | CHECK (jsonb_array_length(roles) > 0) | 用户角色列表，如 `['CREATOR', 'MERCHANT_STAFF']` |
| `created_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 注册时间 |
| `updated_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 更新时间 |
| `deleted_at` | TIMESTAMP | ❌ | `NULL` | - | 软删除时间 |

**索引**：
```sql
CREATE UNIQUE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_wechat ON users(wechat_openid) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_roles ON users USING GIN(roles);
```

##### 3.2 商家表（merchants）

| 字段名 | 类型 | 必填 | 默认值 | 约束 | 说明 |
|--------|------|------|--------|------|------|
| `id` | UUID | ✅ | 自动生成 | PRIMARY KEY | 商家唯一标识 |
| `admin_id` | UUID | ✅ | - | FOREIGN KEY → users.id | 商家管理员用户ID |
| `name` | VARCHAR(100) | ✅ | - | NOT NULL | 商家名称 |
| `industry` | VARCHAR(50) | ❌ | `NULL` | - | 所属行业 |
| `description` | TEXT | ❌ | `NULL` | - | 商家简介 |
| `logo_url` | VARCHAR(500) | ❌ | `NULL` | - | 商家 Logo URL |
| `status` | ENUM | ✅ | `'ACTIVE'` | CHECK (status IN ('ACTIVE', 'INACTIVE')) | 商家状态 |
| `created_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 创建时间 |
| `updated_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 更新时间 |
| `deleted_at` | TIMESTAMP | ❌ | `NULL` | - | 软删除时间 |

**约束**：
```sql
-- 每个商家必须有唯一的管理员
CREATE UNIQUE INDEX idx_merchants_admin ON merchants(admin_id) WHERE deleted_at IS NULL;

-- 商家名称唯一性（同软删除）
CREATE UNIQUE INDEX idx_merchants_name ON merchants(name) WHERE deleted_at IS NULL;
```

##### 3.3 营销活动表（campaigns）

| 字段名 | 类型 | 必填 | 默认值 | 约束 | 说明 |
|--------|------|------|--------|------|------|
| `id` | UUID | ✅ | 自动生成 | PRIMARY KEY | 活动唯一标识 |
| `merchant_id` | UUID | ✅ | - | FOREIGN KEY → merchants.id | 商家ID |
| `title` | VARCHAR(200) | ✅ | - | NOT NULL | 活动标题 |
| `description` | TEXT | ✅ | - | NOT NULL | 活动要求详情（富文本） |
| `platforms` | JSONB | ✅ | `'["xiaohongshu"]'` | CHECK (jsonb_array_length(platforms) > 0) | 任务平台列表 |
| `quota` | INTEGER | ✅ | - | CHECK (quota > 0) | 活动名额总数 |
| `task_amount` | DECIMAL(10,2) | ✅ | - | CHECK (task_amount > 0) | 每个名额的积分费用 |
| `campaign_amount` | DECIMAL(10,2) | ✅ | - | CHECK (campaign_amount = task_amount * quota) | 活动总费用 |
| `creator_amount` | DECIMAL(10,2) | ✅ | - | CHECK (creator_amount > 0) | 达人每任务收入 |
| `staff_referral_amount` | DECIMAL(10,2) | ✅ | - | CHECK (staff_referral_amount >= 0) | 员工每任务返佣 |
| `provider_amount` | DECIMAL(10,2) | ✅ | - | CHECK (provider_amount >= 0) | 服务商每任务收入 |
| `task_deadline` | TIMESTAMP | ✅ | - | NOT NULL | 接任务截止时间 |
| `submission_deadline` | TIMESTAMP | ✅ | - | CHECK (submission_deadline > task_deadline) | 提交截止时间 |
| `status` | ENUM | ✅ | `'DRAFT'` | CHECK (status IN ('DRAFT', 'OPEN', 'CLOSED')) | 活动状态 |
| `created_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 创建时间 |
| `updated_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 更新时间 |
| `deleted_at` | TIMESTAMP | ❌ | `NULL` | - | 软删除时间 |

**索引**：
```sql
-- 商家查看自己的活动
CREATE INDEX idx_campaigns_merchant ON campaigns(merchant_id) WHERE deleted_at IS NULL;

-- 达人浏览开放活动
CREATE INDEX idx_campaigns_status ON campaigns(status, created_at DESC) WHERE deleted_at IS NULL AND status = 'OPEN';

-- 服务商查看待设置活动
CREATE INDEX idx_campaigns_draft ON campaigns(merchant_id, status) WHERE deleted_at IS NULL AND status = 'DRAFT';
```

##### 3.4 任务名额表（tasks）

| 字段名 | 类型 | 必填 | 默认值 | 约束 | 说明 |
|--------|------|------|--------|------|------|
| `id` | UUID | ✅ | 自动生成 | PRIMARY KEY | 名额唯一标识 |
| `campaign_id` | UUID | ✅ | - | FOREIGN KEY → campaigns.id | 活动ID |
| `creator_id` | UUID | ❌ | `NULL` | FOREIGN KEY → creators.id | 接任务的达人ID |
| `status` | ENUM | ✅ | `'OPEN'` | CHECK (status IN ('OPEN', 'ASSIGNED', 'SUBMITTED', 'APPROVED', 'REJECTED')) | 名额状态 |
| `invited_by_staff_id` | UUID | ❌ | `NULL` | FOREIGN KEY → users.id | 邀请的员工用户ID |
| `platform_url` | VARCHAR(500) | ❌ | `NULL` | - | 达人提交的平台链接 |
| `screenshots` | JSONB | ❌ | `'[]'` | - | 达人提交的截图数组 |
| `submitted_at` | TIMESTAMP | ❌ | `NULL` | - | 提交时间 |
| `reviewed_at` | TIMESTAMP | ❌ | `NULL` | - | 审核时间 |
| `reviewer_id` | UUID | ❌ | `NULL` | FOREIGN KEY → users.id | 审核人用户ID |
| `rejection_reason` | TEXT | ❌ | `NULL` | - | 拒绝原因 |
| `version` | INTEGER | ✅ | `0` | CHECK (version >= 0) | 乐观锁版本号（防止并发审核重复结算） |
| `created_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 创建时间 |
| `updated_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 更新时间 |

**约束**：
```sql
-- 防止重复接任务（同一达人不能接同一活动的多个名额）
CREATE UNIQUE INDEX idx_tasks_creator_campaign ON tasks(creator_id, campaign_id) WHERE deleted_at IS NULL AND status != 'OPEN';

-- 防止同一名额被多次接取
CREATE UNIQUE INDEX idx_tasks_status_open ON tasks(id, status) WHERE status = 'OPEN';
```

**索引**：
```sql
-- 达人查看自己的任务
CREATE INDEX idx_tasks_creator ON tasks(creator_id, created_at DESC) WHERE creator_id IS NOT NULL;

-- 服务商按活动筛选待审核任务
CREATE INDEX idx_tasks_campaign_status ON tasks(campaign_id, status) WHERE status IN ('SUBMITTED', 'ASSIGNED');

-- 按邀请人筛选任务
CREATE INDEX idx_tasks_inviter ON tasks(invited_by_staff_id) WHERE invited_by_staff_id IS NOT NULL;
```

##### 3.5 积分账户表（credit_accounts）

| 字段名 | 类型 | 必填 | 默认值 | 约束 | 说明 |
|--------|------|------|--------|------|------|
| `id` | UUID | ✅ | 自动生成 | PRIMARY KEY | 账户唯一标识 |
| `account_type` | ENUM | ✅ | - | CHECK (account_type IN ('USER_PERSONAL', 'ORG_MERCHANT', 'ORG_PROVIDER')) | 账户类型 |
| `owner_id` | UUID | ✅ | - | - | 账户所有者ID（user_id / merchant_id / provider_id） |
| `balance` | DECIMAL(12,2) | ✅ | `0.00` | CHECK (balance >= 0) | 可用余额 |
| `frozen_balance` | DECIMAL(12,2) | ✅ | `0.00` | CHECK (frozen_balance >= 0) | 冻结余额 |
| `created_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 创建时间 |
| `updated_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 更新时间 |
| `deleted_at` | TIMESTAMP | ❌ | `NULL` | - | 软删除时间 |

**约束**：
```sql
-- 防止余额为负数（财务安全保护）
ALTER TABLE credit_accounts
ADD CONSTRAINT check_balance_non_negative
CHECK (balance >= 0);

ALTER TABLE credit_accounts
ADD CONSTRAINT check_frozen_balance_non_negative
CHECK (frozen_balance >= 0);

-- 每个所有者每种账户类型只能有一个账户
CREATE UNIQUE INDEX idx_credit_accounts_owner_type ON credit_accounts(owner_id, account_type) WHERE deleted_at IS NULL;
```

##### 3.6 积分流水表（credit_transactions）

| 字段名 | 类型 | 必填 | 默认值 | 约束 | 说明 |
|--------|------|------|--------|------|------|
| `id` | UUID | ✅ | 自动生成 | PRIMARY KEY | 流水唯一标识 |
| `account_id` | UUID | ✅ | - | FOREIGN KEY → credit_accounts.id | 账户ID |
| `transaction_type` | VARCHAR(50) | ✅ | - | FOREIGN KEY → transaction_types.code | 交易类型代码 |
| `amount` | DECIMAL(12,2) | ✅ | - | CHECK (amount != 0) | 交易金额（正数=收入，负数=支出） |
| `balance_after` | DECIMAL(12,2) | ✅ | - | - | 交易后余额 |
| `frozen_balance_after` | DECIMAL(12,2) | ✅ | - | - | 交易后冻结余额 |
| `transaction_group_id` | UUID | ❌ | `NULL` | - | 交易组ID（关联相关交易） |
| `related_task_id` | UUID | ❌ | `NULL` | FOREIGN KEY → tasks.id | 关联任务ID |
| `related_withdrawal_id` | UUID | ❌ | `NULL` | FOREIGN KEY → withdrawals.id | 关联提现ID |
| `description` | VARCHAR(500) | ❌ | `NULL` | - | 交易说明 |
| `metadata` | JSONB | ❌ | `'{}'` | - | 额外元数据（如对手方信息） |
| `created_at` | TIMESTAMP | ✅ | `CURRENT_TIMESTAMP` | NOT NULL | 交易时间 |

**索引**：
```sql
-- 按账户查询交易记录
CREATE INDEX idx_transactions_account ON credit_transactions(account_id, created_at DESC);

-- 按交易组查询关联交易
CREATE INDEX idx_transactions_group ON credit_transactions(transaction_group_id) WHERE transaction_group_id IS NOT NULL;

-- 按关联任务查询
CREATE INDEX idx_transactions_task ON credit_transactions(related_task_id) WHERE related_task_id IS NOT NULL;

-- 按交易类型统计
CREATE INDEX idx_transactions_type ON credit_transactions(transaction_type, created_at DESC);
```

#### 4. 字段约束与校验规则

##### 4.1 必填字段校验

所有标记为 `✅` 必填的字段在数据库层有 `NOT NULL` 约束，应用层也需校验：

```javascript
// 应用层校验示例
if (!title || title.trim().length === 0) {
  throw new Error('活动标题不能为空');
}
if (title.length > 200) {
  throw new Error('活动标题不能超过200字符');
}
```

##### 4.2 枚举值校验

所有 `ENUM` 类型字段在数据库层有 `CHECK` 约束，应用层需使用常量：

```javascript
// 错误示例：硬编码字符串
campaign.status = 'OPEN';

// 正确示例：使用常量
import { CampaignStatus } from '@/constants/campaign';
campaign.status = CampaignStatus.OPEN;
```

##### 4.3 数值范围校验

```javascript
// 金额必须大于0
if (task_amount <= 0) {
  throw new Error('任务金额必须大于0');
}

// 时间先后关系
if (submission_deadline <= task_deadline) {
  throw new Error('提交截止时间必须晚于接任务截止时间');
}
```

##### 4.4 JSONB 字段校验

```javascript
// platforms 必须是非空数组
if (!Array.isArray(platforms) || platforms.length === 0) {
  throw new Error('任务平台不能为空');
}

// screenshots 格式校验
if (!Array.isArray(screenshots)) {
  throw new Error('截图必须是数组');
}
```

---

## 🔐 权限矩阵

### 功能权限对照表

**请填写每个角色是否可以执行该操作（✅ = 可以，❌ = 不可以，⚠️ = 需要审批）**

| 功能模块 | 具体操作 | 超级管理员 | 商家管理员 | 商家员工 | 服务商管理员 | 服务商员工 | 达人 |
|---------|---------|-----------|-----------|---------|-------------|-----------|---------|------|
| **用户管理** | | | | | | | | |
| | 查看自己信息 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| | 修改自己信息 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| | 查看所有用户 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 修改用户角色 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **商家管理** | | | | | | | | |
| | 查看商家列表 | ✅ | 自己 | 自己 | 绑定的 | 绑定的 | ❌ | ❌ |
| | 创建商家 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 编辑商家信息 | ✅ | 自己 | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 添加商家员工 | ✅ | 自己 | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 删除商家员工 | ✅ | 自己 | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 停用/启用商家 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **服务商管理** | | | | | | | | |
| | 查看服务商列表 | ✅ | ❌ | ❌ | 自己 | 自己 | ❌ | ❌ |
| | 创建服务商 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 编辑服务商信息 | ✅ | ❌ | ❌ | 自己 | ❌ | ❌ | ❌ |
| | 绑定商家 | ✅ | ❌ | ❌ | 自己 | ❌ | ❌ | ❌ |
| | 解绑商家 | ✅ | ❌ | ❌ | 自己 | ❌ | ❌ | ❌ |
| | 绑定达人 | ✅ | ❌ | ❌ | 自己 | ❌ | ❌ | ❌ |
| | 解绑达人 | ✅ | ❌ | ❌ | 自己 | ❌ | ❌ | ❌ |
| | 添加服务商员工 | ✅ | ❌ | ❌ | 自己 | ❌ | ❌ | ❌ |
| **达人管理** | | | | | | | | |
| | 查看达人列表 | ✅ | ❌ | ❌ | 绑定的 | 绑定的 | ❌ |
| | 查看达人档案 | ✅ | ❌ | ❌ | 绑定的 | 绑定的 | 自己 |
| | 修改达人等级 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 停用/启用达人 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **任务管理** | | | | | | | | |
| | 查看任务列表 | ✅ | 自己的 | 自己的 | 所有绑定的 | 所有绑定的 | 所有可接的 |
| | 创建任务 | ✅ | ✅ | ⚠️ | ✅ | ⚠️ | ❌ | ❌ |
| | 编辑任务 | ✅ | 自己的 | ⚠️ | 自己的 | ⚠️ | ❌ | ❌ |
| | 删除任务 | ✅ | 自己的 | ❌ | 自己的 | ❌ | ❌ | ❌ |
| | 发布任务 | ✅ | ✅ | ⚠️ | ✅ | ⚠️ | ❌ | ❌ |
| | 查看任务详情 | ✅ | 自己的 | 自己的 | 所有绑定的 | 所有绑定的 | 下属的 | 自己的 |
| | 接任务（抢单） | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| | 提交任务 | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| | 审核任务-通过 | ✅ | ❌ | ❌ | ✅ | ⚠️ | ❌ | ❌ |
| | 审核任务-拒绝 | ✅ | ❌ | ❌ | ✅ | ⚠️ | ❌ | ❌ |
| **积分系统** | | | | | | | | |
| | 查看自己的积分 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| | 查看自己的积分历史 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| | 充值积分 | ✅ | ✅ | ⚠️ | ❌ | ❌ | ❌ | ❌ |
| | 申请提现 | ❌ | ✅ | ⚠️ | ✅ | ⚠️ | ❌ | ✅ |
| | 审核提现 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 手动调整积分 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| | 查看平台积分流水 | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |

**请标注需要修改的权限**：
```
📝 待确认：
1. 商家员工是否可以发布任务？（默认：需要审批）
2. 服务商员工是否可以审核任务？（默认：需要审批）
3. 商家是否可以申请提现？（默认：可以）

请填写你的修改：_________________________________
```

---

## 📊 多维表界面设计

### 概述

任务管理的核心界面是**多维表（类似Airtable/Notion Database）**，支持灵活的筛选、排序、批量操作。每个任务名额都是独立的一行，完美对应数据库中的tasks记录。

**设计理念**：
- **预分配模式**：发布活动时创建所有名额记录（如：10人×100元 = 10行）
- **差异化视图**：不同角色看到不同的列和操作权限
- **状态驱动**：所有操作围绕状态流转（OPEN → ASSIGNED → SUBMITTED → APPROVED/REJECTED）
- **批量友好**：支持批量审核、批量分配、批量设置

---

### 核心列设计

#### 基础列（所有角色可见）

| 列名 | 字段 | 说明 | 示例 |
|------|------|------|------|
| 名额编号 | `task_slot_number` | 活动内的编号，从1开始 | #1, #2, #3... |
| 状态 | `status` | OPEN/ASSIGNED/SUBMITTED/APPROVED/REJECTED | 🟢开放中 🟡已分配 🔵待审核 ✅已完成 ❌已拒绝 |
| 接单达人 | `creator_id` | 关联达人信息 | 张三（完成率95%） |

#### 审核列（商家、服务商、管理员可见）

| 列名 | 字段 | 说明 | 示例 |
|------|------|------|------|
| 邀请人 | `invited_by_staff_id` | 哪个员工的邀请码 | 李四（客服专员） |
| 接单时间 | `assigned_at` | 达人接任务的时间 | 2024-02-01 10:30 |
| 提交时间 | `submitted_at` | 达人提交任务的时间 | 2024-02-01 15:20 |
| 审核状态 | `audit_status` | 待审核/已通过/已拒绝 | ⏳待审核 |
| 审核时间 | `audited_at` | 服务商审核的时间 | 2024-02-01 16:00 |
| 审核备注 | `audit_note` | 审核意见 | "图片不够清晰" |

#### 扩展列（可选）

| 列名 | 字段 | 类型 | 说明 |
|------|------|------|------|
| 优先级 | `priority` | ENUM | 高/中/低 |
| 标签 | `tags` | TEXT[] | 紧急、VIP、测试任务 |

#### 任务内容列（根据提交情况显示）

| 列名 | 字段 | 说明 | 示例 |
|------|------|------|------|
| 提交链接 | `platform_url` | 达人提交的平台链接 | https://xiaohongshu.com/... |
| 提交截图 | `screenshots` | 达人提交的截图 | [图片1, 图片2] |
| 备注 | `notes` | 达人或服务商的备注 | "达人反馈：流程顺畅" |

---

### 差异化视图设计

#### 商家视图

**路由**：`/workspace/merchant/campaigns/:id/tasks`

**可见列**：名额编号、状态、接单达人、邀请人、提交时间、审核状态、备注

```
┌──────┬──────────┬─────────┬──────────┬────────┬──────────┬────────┐
│ 编号 │ 状态     │ 接单达人│ 邀请人   │ 提交时间│ 审核状态 │ 备注   │
├──────┼──────────┼─────────┼──────────┼────────┼──────────┼────────┤
│ #1   │ ✅已完成 │ 张三    │ 李四     │ 2/1 15:│ ✅通过   │ -      │
│ #2   │ 🔵待审核 │ 李四    │ 王五     │ 2/1 16:│ ⏳待审核 │ -      │
│ #3   │ 🟡已分配 │ 王五    │ 李四     │ -      │ -        │ -      │
│ #4   │ 🟢开放中 │ -       │ -        │ -      │ -        │ -      │
└──────┴──────────┴─────────┴──────────┴────────┴──────────┴────────┘
```

**权限**：
- ✅ 查看所有名额
- ❌ 不能审核（由服务商审核）
- ✅ 可以查看审核详情
- ✅ 可以导出数据

---

#### 服务商视图

**路由**：`/workspace/service-provider/campaigns/:id/tasks`

**可见列**：名额编号、状态、接单达人、邀请人、提交内容、提交时间、审核状态、操作

```
┌──────┬──────────┬─────────┬──────────┬────────┬──────────┬──────────┬────────┐
│ 编号 │ 状态     │ 接单达人│ 邀请人   │ 提交内容│ 提交时间 │ 审核状态 │ 操作   │
├──────┼──────────┼─────────┼──────────┼────────┼──────────┼──────────┼────────┤
│ #1   │ 🔵待审核 │ 张三    │ 李四     │ [查看] │ 2/1 15:00│ ⏳待审核 │[通过][拒绝]│
│ #2   │ 🔵待审核 │ 李四    │ 王五     │ [查看] │ 2/1 16:00│ ⏳待审核 │[通过][拒绝]│
│ #3   │ 🟡已分配 │ 王五    │ 李四     │ -      │ -        │ -        │[提醒]  │
└──────┴──────────┴─────────┴──────────┴────────┴──────────┴──────────┴────────┘
```

**权限**：
- ✅ 查看所有名额
- ✅ 审核待提交名额（通过/拒绝）
- ✅ 批量审核
- ✅ 发送提醒给未提交的达人
- ✅ 查看每个员工邀请的名额统计

---

#### 达人视图

**路由**：`/workspace/creator/my-tasks`

**可见列**：活动名称、状态、接单时间、提交时间、审核状态、备注

```
┌──────────────────┬──────────┬──────────┬──────────┬────────┬──────────┐
│ 活动名称         │ 状态     │ 接单时间 │ 提交时间 │ 审核状态│ 备注     │
├──────────────────┼──────────┼──────────┼──────────┼────────┼──────────┤
│ 新品体验推广     │ ✅已完成 │ 2/1 10:00│ 2/1 15:00│ ✅通过 │ -        │
│ 双11大促         │ 🟡进行中 │ 2/1 11:00│ -        │ -      │ 请尽快提交│
└──────────────────┴──────────┴──────────┴──────────┴────────┴──────────┘
```

**权限**：
- ✅ 只能看到自己接的任务
- ✅ 提交任务内容
- ❌ 不能看到其他达人的任务
- ❌ 不能审核

---

#### 超级管理员视图

**路由**：`/workspace/admin/campaigns/:id/tasks`

**可见列**：所有列（包括系统字段）

**权限**：
- ✅ 全局查看所有活动、所有名额
- ✅ 查看所有角色的操作日志
- ✅ 紧急情况下可以干预审核

---

### 核心交互功能

#### 1. 状态筛选（优先级最高）

**筛选器位置**：表格上方

**支持的筛选**：
- **按状态筛选**：
  - 只看未分配（status = 'OPEN'）
  - 只看已分配（status = 'ASSIGNED'）
  - 只看待审核（status = 'SUBMITTED'）
  - 只看已完成（status = 'APPROVED'）
  - 只看已拒绝（status = 'REJECTED'）

- **按邀请人筛选**：
  - 查看某个员工邀请的所有达人（`invited_by_staff_id = ?`）
  - 服务商管理员查看下属员工的绩效

- **按时间筛选**：
  - 今天接单的
  - 逾期未提交的（`submission_deadline < NOW()`）
  - 本周完成的

- **组合筛选**：
  - 待审核 + 我邀请的 + 今天提交的

**实现示例**：
```sql
-- 查询待审核 + 李四邀请的 + 今天提交的
SELECT * FROM tasks
WHERE status = 'SUBMITTED'
  AND invited_by_staff_id = '李四ID'
  AND DATE(submitted_at) = CURRENT_DATE
  AND campaign_id = ?;
```

---

#### 2. 批量操作

**批量选择**：
- 复选框选择多行
- 全选/取消全选

**批量操作按钮**：
- **批量审核通过**：
  ```sql
  UPDATE tasks
  SET status = 'APPROVED',
      audited_by = ?,
      audited_at = NOW()
  WHERE id IN (?, ?, ?)
    AND status = 'SUBMITTED';
  ```

- **批量审核拒绝**：
  ```sql
  UPDATE tasks
  SET status = 'REJECTED',
      audited_by = ?,
      audited_at = NOW(),
      audit_note = ?
  WHERE id IN (?, ?, ?)
    AND status = 'SUBMITTED';
  ```

- **批量设置优先级**：
  ```sql
  UPDATE tasks
  SET priority = 'HIGH'
  WHERE id IN (?, ?, ?);
  ```

- **批量添加标签**：
  ```sql
  UPDATE tasks
  SET tags = array_append(tags, '紧急')
  WHERE id IN (?, ?, ?);
  ```

- **批量发送提醒**：
  - 给未提交的达人发送站内消息
  - 触发推送通知

---

#### 3. 排序分组

**排序**：
- 按名额编号排序（默认）
- 按接单时间排序（最新接单在前）
- 按提交时间排序（最新提交在前）
- 按优先级排序（HIGH > MEDIUM > LOW）

**分组**：
- 按状态分组：
  ```
  🔵 待审核 (5)
    ├─ #1 张三
    ├─ #2 李四
    └─ ...

  🟡 进行中 (3)
    ├─ #3 王五
    └─ ...

  ✅ 已完成 (2)
    ├─ #4 赵六
    └─ ...
  ```

- 按邀请人分组：
  ```
  👤 李四邀请的 (5)
    ├─ #1 张三
    ├─ #3 王五
    └─ ...

  👤 王五邀请的 (3)
    ├─ #2 李四
    └─ ...
  ```

---

#### 4. 单元格编辑

**可编辑单元格**：
- **备注**（`notes`）：双击单元格直接编辑
- **优先级**：点击单元格弹出下拉选择（HIGH/MEDIUM/LOW）
- **标签**：点击单元格弹出标签选择器

**点击查看详情**：
- 点击名额编号 → 右侧滑出详情抽屉
- 详情抽屉显示：
  - 所有字段信息
  - 达人历史数据
  - 操作日志
  - 审核历史（如果有）

---

### 多维表API需求

#### 1. 查询任务名额（支持筛选、排序、分页）

**端点**：`GET /api/campaigns/:id/tasks`

**查询参数**：
```
?status=SUBMITTED              # 按状态筛选
&invited_by_staff_id=xxx       # 按邀请人筛选
&sort_by=submitted_at          # 排序字段
&sort_order=desc               # 排序方向
&page=1                        # 页码
&page_size=20                  # 每页数量
```

**响应示例**：
```json
{
  "data": [
    {
      "id": "task-1",
      "task_slot_number": 1,
      "status": "SUBMITTED",
      "creator": {
        "id": "creator-1",
        "name": "张三",
        "completion_rate": 0.95
      },
      "invited_by_staff": {
        "id": "staff-1",
        "name": "李四"
      },
      "submitted_at": "2024-02-01T15:00:00Z",
      "platform_url": "https://xiaohongshu.com/...",
      "screenshots": ["url1", "url2"]
    }
  ],
  "pagination": {
    "total": 50,
    "page": 1,
    "page_size": 20,
    "total_pages": 3
  }
}
```

---

#### 2. 批量审核

**端点**：`POST /api/tasks/batch-review`

**请求体**：
```json
{
  "task_ids": ["task-1", "task-2", "task-3"],
  "action": "approve",  // or "reject"
  "audit_note": "质量很好"
}
```

**响应**：
```json
{
  "success": true,
  "updated_count": 3,
  "failed_ids": []
}
```

---

#### 3. 批量更新

**端点**：`POST /api/tasks/batch-update`

**请求体**：
```json
{
  "task_ids": ["task-1", "task-2"],
  "updates": {
    "priority": "HIGH",
    "tags": ["紧急", "VIP"]
  }
}
```

---

### 多维表组件实现建议

**前端技术选型**：
- **React** + **TanStack Table**（原React Table）- 强大的表格功能
- 或 **AG Grid** - 企业级表格组件（更强大但更重）

**核心功能模块**：
1. **数据获取**：`useTasks` Hook（支持筛选、排序、分页）
2. **筛选器组件**：`TaskFilters`（状态、邀请人、时间范围）
3. **批量操作组件**：`BatchActions`（审核、更新、提醒）
4. **单元格编辑**：`EditableCell`（备注、优先级、标签）
5. **详情抽屉**：`TaskDetailDrawer`（完整信息、操作日志）

---

## 📊 达人管理多维表

### 概述

达人管理多维表用于管理所有达人用户，支持查看达人基础信息、数据分析、标签管理和分级管理。

**设计理念**：
- **数据驱动**：直观展示达人完成率、收入、任务数等核心指标
- **灵活筛选**：支持按等级、状态、标签等多维度筛选
- **批量管理**：支持批量打标签、批量修改状态
- **详情快捷**：点击达人姓名快速查看完整档案

---

### 核心列设计

#### 基础信息列

| 列名 | 字段 | 说明 | 示例 |
|------|------|------|------|
| 微信头像 | `wechat_avatar` | 达人微信头像 | [图片] |
| 微信昵称 | `wechat_nickname` | 达人微信昵称（可点击查看详情） | 张三 |
| 达人等级 | `level` | UGC/KOC/INF/KOL | KOC |
| 粉丝数 | `followers_count` | 达人粉丝数量 | 10,000 |
| 注册时间 | `created_at` | 达人注册时间 | 2024-01-15 |
| 状态 | `status` | active/banned/inactive | ✅正常 / 🔴封禁 |

#### 数据统计列（计算字段）

| 列名 | 计算逻辑 | 说明 | 示例 |
|------|---------|------|------|
| 完成任务数 | `COUNT(tasks.id)` WHERE status='APPROVED' | 已完成的任务数量 | 25 |
| 进行中任务 | `COUNT(tasks.id)` WHERE status='ASSIGNED/SUBMITTED' | 进行中的任务数量 | 3 |
| 完成率 | `完成任务数 / 总接任务数` | 任务完成率 | 89.3% |
| 总收入 | `SUM(creator_amount)` | 从积分账户统计 | 2,500积分 |
| 平均完成时间 | `AVG(完成时间 - 接单时间)` | 平均完成时长 | 2.5天 |
| 邀请人 | `inviter_id` + `inviter_type` | 谁邀请的该达人 | 李四（员工） |

#### 扩展字段列

| 列名 | 字段 | 类型 | 说明 |
|------|------|------|------|
| 标签 | `tags` | TEXT[] | 自定义标签，如：优质达人、新人、活跃 |
| 备注 | `notes` | TEXT | 管理员备注 |

---

### 差异化视图设计

#### 超级管理员视图

**路由**：`/workspace/admin/creators`

**可见列**：所有列

**权限**：
- ✅ 查看所有达人
- ✅ 封禁/解封达人
- ✅ 修改达人等级
- ✅ 批量打标签
- ✅ 查看完整数据

```
┌──────────┬──────────┬────────┬────────┬──────────┬──────────┬──────────┬────────┐
│ 头像     │ 昵称     │ 等级   │ 粉丝数  │ 完成任务数 │ 完成率   │ 总收入   │ 状态   │
├──────────┼──────────┼────────┼────────┼──────────┼──────────┼──────────┼────────┤
│ [头像]   │ 张三     │ KOC    │ 10,000 │ 25       │ 89.3%    │ 2,500    │ ✅正常 │
│ [头像]   │ 李四     │ INF    │ 50,000 │ 48       │ 95.2%    │ 4,800    │ ✅正常 │
│ [头像]   │ 王五     │ UGC    │ 1,000  │ 2        │ 66.7%    │ 200      | 🔴封禁│
└──────────┴──────────┴────────┴────────┴──────────┴──────────┴──────────┴────────┘
```

---

#### 服务商视图

**路由**：`/workspace/service-provider/creators`

**可见列**：微信头像、微信昵称、等级、粉丝数、完成任务数、完成率、邀请人

**权限**：
- ✅ 查看绑定的达人（通过自己员工邀请的）
- ✅ 查看达人任务数据
- ✅ 给达人打标签（仅自己可见）
- ❌ 不能修改达人等级
- ❌ 不能封禁达人

**筛选默认条件**：`inviter_id IN (本服务商员工)`

```
┌──────────┬──────────┬────────┬────────┬──────────┬──────────┬──────────┐
│ 头像     │ 昵称     │ 等级   │ 粉丝数  │ 完成任务数 │ 完成率   │ 邀请人   │
├──────────┼──────────┼────────┼────────┼──────────┼──────────┼──────────┤
│ [头像]   │ 张三     │ KOC    │ 10,000 │ 15       │ 89.3%    │ 李四     │
│ [头像]   │ 赵六     │ INF    │ 30,000 │ 22       │ 91.5%    │ 王五     │
└──────────┴──────────┴────────┴────────┴──────────┴──────────┴──────────┘
```

---

#### 商家视图

**路由**：`/workspace/merchant/creators`

**说明**：商家一般不直接管理达人，此视图可能不展示或仅提供只读访问。

---

### 核心交互功能

#### 1. 状态筛选

**支持的筛选**：
- **按等级筛选**：UGC / KOC / INF / KOL
- **按状态筛选**：正常 / 封禁 / 未激活
- **按完成率筛选**：高完成率（>90%） / 中等（70-90%） / 低（<70%）
- **按标签筛选**：优质达人 / 新人 / 活跃
- **按邀请人筛选**：查看某个员工邀请的所有达人

**实现示例**：
```sql
-- 查询KOC等级 + 完成率>90% + 正常状态
SELECT c.*,
  COUNT(t.id) FILTER (WHERE t.status = 'APPROVED') as completed_count,
  COUNT(t.id) as total_count,
  CASE WHEN COUNT(t.id) > 0 THEN
    ROUND(100.0 * COUNT(t.id) FILTER (WHERE t.status = 'APPROVED') / COUNT(t.id), 2)
  ELSE 0 END as completion_rate
FROM creators c
LEFT JOIN tasks t ON t.creator_id = c.id
WHERE c.level = 'KOC'
  AND c.status = 'active'
GROUP BY c.id
HAVING COUNT(t.id) > 0 AND
  (COUNT(t.id) FILTER (WHERE t.status = 'APPROVED'))::FLOAT / COUNT(t.id) > 0.9
ORDER BY completion_rate DESC;
```

---

#### 2. 批量操作

**批量操作按钮**：
- **批量打标签**：
  ```sql
  UPDATE creators
  SET tags = array_append(tags, '优质达人')
  WHERE id IN (?, ?, ?);
  ```

- **批量修改等级**：
  ```sql
  UPDATE creators
  SET level = 'KOL'
  WHERE id IN (?, ?, ?);
  ```

- **批量封禁**：
  ```sql
  UPDATE creators
  SET status = 'banned'
  WHERE id IN (?, ?, ?);
  ```

---

#### 3. 排序分组

**排序**：
- 按完成率排序（高→低）
- 按完成任务数排序（多→少）
- 按总收入排序（高→低）
- 按粉丝数排序（多→少）
- 按注册时间排序（最新→最早）

**分组**：
- 按等级分组：
  ```
  🌟 KOL (5)
  💎 INF (12)
  ⭐ KOC (28)
  👤 UGC (45)
  ```

- 按状态分组：
  ```
  ✅ 正常 (85)
  🔴 封禁 (3)
  ⚪ 未激活 (2)
  ```

---

#### 4. 单元格编辑与详情查看

**可编辑单元格**：
- **标签**：点击单元格弹出标签选择器
- **备注**：双击单元格直接编辑

**点击查看详情**：
- 点击达人微信昵称 → 右侧滑出详情抽屉
- 详情抽屉显示：
  - 所有基础信息
  - 完整数据统计（任务列表、收入曲线、完成率趋势）
  - 历史任务记录
  - 标签管理
  - 操作日志

---

### 达人管理API需求

#### 1. 查询达人列表（支持筛选、排序、分页）

**端点**：`GET /api/admin/creators`

**查询参数**：
```
?level=KOC                  # 按等级筛选
&status=active              # 按状态筛选
&min_completion_rate=90     # 最低完成率
&inviter_id=xxx             # 按邀请人筛选
&sort_by=completion_rate    # 排序字段
&sort_order=desc            # 排序方向
&page=1                     # 页码
&page_size=20               # 每页数量
```

**响应示例**：
```json
{
  "data": [
    {
      "id": "creator-1",
      "wechat_nickname": "张三",
      "wechat_avatar": "https://...",
      "level": "KOC",
      "followers_count": 10000,
      "status": "active",
      "statistics": {
        "completed_count": 25,
        "in_progress_count": 3,
        "completion_rate": 89.3,
        "total_income": 2500,
        "avg_completion_days": 2.5
      },
      "inviter": {
        "id": "staff-1",
        "name": "李四",
        "type": "PROVIDER_STAFF"
      },
      "tags": ["优质达人", "活跃"],
      "created_at": "2024-01-15T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 90,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

---

#### 2. 批量更新达人

**端点**：`POST /api/admin/creators/batch-update`

**请求体**：
```json
{
  "creator_ids": ["creator-1", "creator-2"],
  "updates": {
    "level": "KOL",
    "tags": ["VIP", "优质达人"]
  }
}
```

---

## 📊 员工管理多维表

### 概述

员工管理多维表用于管理商家员工和服务商员工，支持查看员工基础信息、权限管理、绩效追踪和操作审计。

**设计理念**：
- **权限透明**：直观展示员工权限列表
- **绩效可视**：追踪员工邀请达人数量、完成任务情况
- **操作审计**：记录员工关键操作
- **灵活授权**：批量授予/撤销权限

---

### 核心列设计

#### 基础信息列

| 列名 | 字段 | 说明 | 示例 |
|------|------|------|------|
| 员工姓名 | `users.name` | 关联users表获取真实姓名 | 张三 |
| 微信头像 | `users.wechat_avatar` | 微信头像 | [图片] |
| 微信昵称 | `users.wechat_nickname` | 微信昵称 | 张三 |
| 职能名称 | `title` | 员工职位 | 运营专员 |
| 所属组织 | `merchant_id/provider_id` | 关联商家/服务商名称 | XX商家 |
| 状态 | `status` | active/inactive | ✅正常 / ⚪离职 |

#### 权限概览列（计算字段）

| 列名 | 计算逻辑 | 说明 | 示例 |
|------|---------|------|------|
| 权限数量 | `COUNT(permission_code)` | 拥有的权限数量 | 5 |
| 权限列表 | `ARRAY_AGG(permission_code)` | 权限代码列表 | [task.audit, task.view, ...] |

#### 绩效指标列（计算字段）

| 列名 | 计算逻辑 | 说明 | 示例 |
|------|---------|------|------|
| 邀请达人数量 | `COUNT(creators.id)` WHERE inviter_id=staff_id | 邀请的达人总数 | 15 |
| 活跃达人数量 | `COUNT(DISTINCT tasks.creator_id)` WHERE tasks.invited_by_staff_id=staff_id | 实际接任务的达人 | 12 |
| 完成任务总数 | `COUNT(tasks.id)` WHERE invited_by_staff_id=staff_id AND status='APPROVED' | 通过审核的任务数 | 35 |
| 产生的GMV | `SUM(creator_amount + staff_referral_amount)` | 产生的交易总额 | 3,500积分 |

#### 活动状态列

| 列名 | 字段 | 说明 | 示例 |
|------|------|------|------|
| 最近登录时间 | `users.last_login_at` | 最后登录时间 | 2小时前 |
| 最近操作时间 | `MAX(操作日志时间)` | 最后操作时间 | 30分钟前 |
| 添加时间 | `created_at` | 员工添加时间 | 2024-01-15 |

---

### 差异化视图设计

#### 商家员工管理视图

**路由**：`/workspace/merchant/staff`

**可见列**：员工姓名、职能名称、权限数量、最近登录时间、状态

**权限**：
- ✅ 查看所有员工（仅限管理员）
- ✅ 添加/编辑员工
- ✅ 分配/撤销权限
- ✅ 停用/启用员工
- ❌ 不能删除员工（软删除）

```
┌──────────┬──────────┬──────────┬──────────┬────────┐
│ 员工姓名 │ 职能名称 │ 权限数量 │ 最近登录 │ 状态   │
├──────────┼──────────┼──────────┼──────────┼────────┤
│ 张三     │ 运营     │ 5        │ 2小时前  │ ✅正常 │
│ 李四     │ 财务     │ 3        │ 1天前    │ ✅正常 │
│ 王五     │ 审核专员 │ 4        │ 5天前    │ ⚪离职 │
└──────────┴──────────┴──────────┴──────────┴────────┘
```

---

#### 服务商员工管理视图

**路由**：`/workspace/service-provider/staff`

**可见列**：员工姓名、职能名称、权限数量、邀请达人数量、完成任务总数、产生的GMV、最近登录时间、状态

**权限**：
- ✅ 查看所有员工（仅限管理员）
- ✅ 添加/编辑员工
- ✅ 分配/撤销权限
- ✅ 查看员工绩效数据
- ✅ 停用/启用员工

```
┌──────────┬──────────┬──────────┬──────────┬──────────┬──────────┬────────┐
│ 员工姓名 │ 职能名称 │ 权限数量 │ 邀请达人 │ 完成任务 │ 产生GMV │ 状态   │
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼────────┤
│ 张三     │ 商务     │ 5        │ 15       │ 35       │ 3,500    │ ✅正常 │
│ 李四     │ 运营     │ 7        │ 22       │ 48       │ 4,800    │ ✅正常 │
│ 王五     │ 审核     │ 4        │ 8        │ 18       │ 1,800    │ ⚪离职 │
└──────────┴──────────┴──────────┴──────────┴──────────┴──────────┴────────┘
```

---

#### 超级管理员视图

**路由**：`/workspace/admin/staff`

**可见列**：所有列 + 操作审计

**权限**：
- ✅ 查看所有商家和服务商的员工
- ✅ 查看操作日志
- ✅ 全局视角

---

### 核心交互功能

#### 1. 状态筛选

**支持的筛选**：
- **按状态筛选**：正常 / 离职
- **按职能筛选**：运营 / 财务 / 审核 / 商务
- **按权限筛选**：拥有某权限的员工
- **按绩效筛选**：高绩效（GMV>5000） / 中等 / 低绩效

**实现示例**：
```sql
-- 查询服务商员工 + 绩效数据
SELECT s.*, u.name, u.wechat_nickname, u.wechat_avatar,
  COUNT(DISTINCT c.id) as invited_creators,
  COUNT(DISTINCT t.id) FILTER (WHERE t.status = 'APPROVED') as completed_tasks,
  COALESCE(SUM(
    (ca.creator_amount + ca.staff_referral_amount)
  ), 0) as generated_gmv
FROM service_provider_staff s
JOIN users u ON u.id = s.user_id
LEFT JOIN creators c ON c.inviter_id = s.id
LEFT JOIN tasks t ON t.invited_by_staff_id = s.id
LEFT JOIN credit_transactions ct ON ct.related_task_slot_id = t.id
LEFT JOIN credit_accounts ca ON ca.id = ct.account_id
WHERE s.provider_id = ?
  AND s.status = 'active'
GROUP BY s.id, u.id
ORDER BY generated_gmv DESC;
```

---

#### 2. 批量操作

**批量操作按钮**：
- **批量授予权限**：
  ```sql
  INSERT INTO merchant_staff_permissions (staff_id, permission_code, granted_by)
  SELECT unnest(ARRAY['staff_id1', 'staff_id2']),
         'task.view',
         'admin_id'
  ON CONFLICT (staff_id, permission_code) DO NOTHING;
  ```

- **批量撤销权限**：
  ```sql
  DELETE FROM merchant_staff_permissions
  WHERE staff_id = ANY(ARRAY['staff_id1', 'staff_id2'])
    AND permission_code = 'task.view';
  ```

- **批量停用**：
  ```sql
  UPDATE merchant_staff
  SET status = 'inactive'
  WHERE id = ANY(ARRAY['staff_id1', 'staff_id2']);
  ```

---

#### 3. 权限管理

**权限抽屉**：
- 点击员工姓名 → 右侧滑出权限管理抽屉
- 显示所有可用权限（按模块分组）
- 勾选/取消勾选权限
- 实时保存

**权限分组**：
```
任务管理
  ☑ task.view     查看任务
  ☑ task.audit    审核任务
  ☐ task.assign   分配任务

财务管理
  ☑ finance.view  查看财务
  ☐ finance.approve 审批提现

达人管理
  ☐ creator.invite 邀请达人
  ☐ creator.manage 管理达人
```

---

#### 4. 绩效追踪

**绩效图表**：
- 点击员工 → 查看绩效详情
- 显示：
  - 邀请达人数量趋势（折线图）
  - 完成任务数量趋势（柱状图）
  - 产生的GMV趋势（面积图）

**绩效对比**：
- 员工之间绩效对比（排行榜）
- 平均完成时间对比
- 达人活跃度对比

---

#### 5. 操作审计

**操作日志**：
- 记录员工的关键操作
  - 登录/登出
  - 审核任务（通过/拒绝）
  - 修改达人信息
  - 分配权限
  - 导出数据

**日志查看**：
- 点击员工 → 操作审计标签
- 显示时间线：
  ```
  2024-02-01 16:30  审核通过任务 #123
  2024-02-01 15:20  邀请达人"张三"
  2024-02-01 14:10  登录系统
  ```

---

### 员工管理API需求

#### 1. 查询员工列表（支持筛选、排序、分页）

**端点**：`GET /api/service-provider/staff`

**查询参数**：
```
?status=active               # 按状态筛选
&title=运营                  # 按职能筛选
&has_permission=task.audit   # 拥有某权限的员工
&min_gmv=5000               # 最低GMV
&sort_by=generated_gmv      # 排序字段
&sort_order=desc            # 排序方向
&page=1                     # 页码
&page_size=20               # 每页数量
```

**响应示例**：
```json
{
  "data": [
    {
      "id": "staff-1",
      "name": "张三",
      "wechat_nickname": "张三",
      "wechat_avatar": "https://...",
      "title": "商务",
      "status": "active",
      "permissions": {
        "count": 5,
        "list": ["task.view", "task.audit", "creator.invite", ...]
      },
      "performance": {
        "invited_creators": 15,
        "completed_tasks": 35,
        "generated_gmv": 3500
      },
      "last_login_at": "2024-02-01T14:00:00Z",
      "created_at": "2024-01-15T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 10,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
```

---

#### 2. 批量更新权限

**端点**：`POST /api/service-provider/staff/batch-grant-permissions`

**请求体**：
```json
{
  "staff_ids": ["staff-1", "staff-2"],
  "permission_codes": ["task.view", "task.audit"]
}
```

---

#### 3. 获取员工操作日志

**端点**：`GET /api/service-provider/staff/:id/audit-logs`

**查询参数**：
```
?page=1
&page_size=50
```

---

### 管理多维表组件实现建议

**前端技术选型**：
- **React** + **TanStack Table** - 灵活的表格组件
- **Recharts** - 数据可视化（绩效图表）

**核心功能模块**：
1. **达人管理模块**：
   - `useCreators` Hook（查询达人列表）
   - `CreatorFilters`（筛选器）
   - `CreatorBatchActions`（批量操作）
   - `CreatorDetailDrawer`（达人详情）

2. **员工管理模块**：
   - `useStaff` Hook（查询员工列表）
   - `StaffFilters`（筛选器）
   - `PermissionManager`（权限管理抽屉）
   - `StaffPerformanceChart`（绩效图表）
   - `AuditLogsTimeline`（操作日志）

---

## 📱 页面路由设计

### 路由总览

```
/                                    # 首页
/login                              # 登录页
/register                           # 注册页（暂不开放）
/workspace                          # 工作区首页（根据角色重定向）
/workspace/admin                     # 超级管理员后台
/workspace/merchant                 # 商家工作台
/workspace/service-provider         # 服务商工作台
/workspace/creator                  # 达人工作台
```

### 1. 公共页面

#### `/` - 首页
**页面元素**：
- [ ] 平台介绍 Banner
- [ ] 角色选择入口（商家、服务商、达人）
- [ ] 成功案例展示
- [ ] 联系方式

**功能**：
- 未登录用户展示
- 点击角色选择跳转到登录页

**交互**：
```
点击"我是商家" → /login?redirect=/workspace/merchant
点击"我是服务商" → /login?redirect=/workspace/service-provider
点击"我是达人" → /login?redirect=/workspace/creator
```

---

#### `/login` - 登录页
**页面元素**：
- [ ] 微信登录按钮（主要方式）
- [ ] 手机号 + 密码登录（备用方式）
- [ ] 注册引导文案

**功能**：
- **微信扫码登录**（通过auth-center）
- **手机号+密码登录**（通过auth-center）
- 登录成功后跳转到对应工作台

**✅ 登录方式说明**：
- auth-center提供两种登录方式，PR Business都支持
- 用户可选择任一方式登录，系统通过auth_center_user_id关联账户

---

### 2. 超级管理员后台 `/workspace/admin`

#### `/workspace/admin` - 管理员首页 Dashboard
**页面布局**：
```
┌─────────────────────────────────────────────┐
│ 顶部导航：[角色切换器] Logo | 管理员后台 | 用户菜单 │
├─────────────────────────────────────────────┤
│ 左侧菜单 │ 主内容区                       │
│         │                                 │
│ 首页    │ [统计卡片]                      │
│ 用户    │ - 总用户数                      │
│ 商家    │ - 总商家数                      │
│ 服务商  │ - 总服务商数                    │
│ 达人    │ - 总任务数                      │
│ 任务    │ - 总交易额                      │
│ 财务    │                                 │
│         │ [待处理事项]                    │
│         │ - 待审核提现: N 条 → 点击查看  │
│         │ - 待处理申诉: N 条 → 点击查看  │
│         │                                 │
└─────────────────────────────────────────────┘
```

**数据展示**（分角色展示不同数据）：

**系统级别数据（超级管理员后台）**：
- ✅ 用户总数：____ 人
- ✅ 商家总数：____ 家
- ✅ 服务商总数：____ 家
- ✅ 达人总数：____ 人
- ✅ 任务总数：____ 个
- ✅ 今日任务：____ 个
- ✅ 今日交易额：____ 元
- ✅ 平台总收入：____ 元
- ✅ 待处理提现：____ 条（点击查看详情）
- ✅ 待处理申诉：____ 条（点击查看详情）
- ✅ 活跃用户数（7日）：____ 人
- ✅ 新增商家数（本周）：____ 家

**服务商级别数据（服务商工作台）**：
- ✅ 绑定商家数：____ 家
- ✅ 进行中任务数：____ 个
- ✅ 本月交易额：____ 元
- ✅ 本月收入：____ 元（返佣+服务费）
- ✅ 待审核任务：____ 个
- ✅ 达人总数：____ 人
- ✅ 新增达人数（本周）：____ 人

**商家级别数据（商家工作台）**：
- ✅ 发布任务数：____ 个
- ✅ 进行中任务：____ 个
- ✅ 已完成任务：____ 个
- ✅ 总花费：____ 元
- ✅ 账户余额：____ 积分
- ✅ 冻结余额：____ 积分
- ✅ 待提交任务：____ 个（已接但未提交）

---

#### `/workspace/admin/users` - 用户管理
**页面功能**：
- [ ] 用户列表（分页）
- [ ] 搜索（按手机号、角色）
- [ ] 筛选（按角色）
- [ ] 查看用户详情
- [ ] 修改用户角色
- [ ] 停用/启用用户

**列表字段**（请确认）：
- [ ] 用户 ID
- [ ] 手机号
- [ ] 角色
- [ ] 积分余额
- [ ] 钻石余额
- [ ] 注册时间
- [ ] 状态（正常/停用）

**操作**：
- [ ] 查看详情
- [ ] 修改角色
- [ ] 停用/启用

---

#### `/workspace/admin/merchants` - 商家管理
**页面功能**：
- [ ] 商家列表（分页）
- [ ] 创建商家（弹出表单）
- [ ] 搜索（按商家名称）
- [ ] 查看商家详情
- [ ] 停用/启用商家

**创建商家表单字段**（请确认）：
- [ ] 商家名称 *（必填）
- [ ] 行业 *（必填，下拉选择）
- [ ] 联系人 *（必填）
- [ ] 联系电话 *（必填）
- [ ] 营业执照（图片上传）
- [ ] 商家地址
- [ ] 备注

**行业选项**（请填写）：
```
📝 请列出你的行业分类：
1. 美妆
2. 服装
3. 食品
4. 电子产品
5. ____（请补充）
6. ____
```

---

#### `/workspace/admin/service-providers` - 服务商管理
**页面功能**：
- [ ] 服务商列表（分页）
- [ ] 创建服务商（弹出表单）
- [ ] 搜索（按服务商名称）
- [ ] 查看服务商详情
- [ ] 停用/启用服务商

**创建服务商表单字段**（请确认）：
- [ ] 服务商名称 *（必填）
- [ ] 联系人 *（必填）
- [ ] 联系电话 *（必填）
- [ ] 营业执照（图片上传）
- [ ] 服务地址
- [ ] 备注

---

#### `/workspace/admin/creators` - 达人管理（多维表）

**页面类型**：达人管理多维表

**核心功能**：
- ✅ 达人多维表展示（头像、昵称、等级、粉丝数、完成任务数、完成率、总收入、状态）
- ✅ 多维度筛选（按等级、状态、完成率、标签、邀请人）
- ✅ 批量操作（批量打标签、批量修改等级、批量封禁）
- ✅ 排序分组（按完成率、任务数、收入排序）
- ✅ 详情抽屉（点击达人姓名查看完整档案、任务列表、收入曲线）
- ✅ 搜索（按达人昵称、手机号）

**实现参考**：详见 PRD.md 第 1188 行"达人管理多维表"章节

---

#### `/workspace/admin/staff` - 员工管理（多维表）

**页面类型**：员工管理多维表

**核心功能**：
- ✅ 员工多维表展示（员工姓名、职能、所属组织、权限数量、绩效指标、状态）
- ✅ 权限管理抽屉（按模块分组显示权限、勾选授予/撤销）
- ✅ 绩效追踪图表（邀请达人趋势、完成任务趋势、GMV趋势）
- ✅ 操作审计日志（登录/登出、审核任务、修改达人、导出数据）
- ✅ 批量操作（批量授予权限、批量撤销权限、批量停用）

**实现参考**：详见 PRD.md 第 1471 行"员工管理多维表"章节

---

#### `/workspace/admin/tasks` - 任务管理
**页面功能**：
- [ ] 所有任务列表（分页）
- [ ] 搜索（按任务标题、商家）
- [ ] 筛选（按状态、时间范围）
- [ ] 查看任务详情

**列表字段**：
- [ ] 任务 ID
- [ ] 任务标题
- [ ] 商家名称
- [ ] 服务商名称（如有）
- [ ] 任务赏金
- [ ] 名额总数
- [ ] 已接名额数
- [ ] 状态
- [ ] 发布时间
- [ ] 截止时间

---

#### `/workspace/admin/finance` - 财务管理
**页面功能**：
- [ ] 财务概览
  - [ ] 平台总余额
  - [ ] 商家充值总额
  - [ ] 达人提现总额
  - [ ] 平台收入
- [ ] 充值记录列表
- [ ] 提现审核列表
- [ ] 积分调整记录

**子页面**：
- [ ] `/workspace/admin/finance/overview` - 财务概览
- [ ] `/workspace/admin/finance/recharge` - 充值记录
- [ ] `/workspace/admin/finance/withdrawals` - 提现审核
- [ ] `/workspace/admin/finance/adjustments` - 积分调整
- [ ] `/workspace/admin/finance/reconciliation` - 财务对账

---

### 3. 商家工作台 `/workspace/merchant`

#### `/workspace/merchant` - 商家首页 Dashboard
**页面布局**：
```
┌─────────────────────────────────────────────┐
│ 顶部导航：                                  │
│ [角色切换器] 商家 | 积分余额 | 充值 | 用户菜单│
├─────────────────────────────────────────────┤
│ 左侧菜单 │ 主内容区                       │
│         │                                 │
│ 首页    │ [统计卡片]                      │
│ 我的任务│ - 总任务数                      │
│ 发布任务│ - 进行中任务                    │
│ 财务    │ - 已完成任务                    │
│ 员工    │ - 总支出                        │
│         │                                 │
│         │ [快捷操作]                      │
│         │ + 发布新任务                    │
│         │ + 充值积分                      │
│         │                                 │
└─────────────────────────────────────────────┘
```

---

#### `/workspace/merchant/tasks` - 我的任务
**页面功能**：
- [ ] 任务列表（Tab 切换）
  - [ ] 全部任务
  - [ ] 草稿中
  - [ ] 已发布（进行中）
  - [ ] 已完成
- [ ] 搜索任务
- [ ] 操作按钮
  - [ ] 创建任务
  - [ ] 编辑任务（仅草稿/已发布）
  - [ ] 查看详情

**列表字段**：
- [ ] 任务标题
- [ ] 赏金
- [ ] 名额（已接/总数）
- [ ] 状态
- [ ] 发布时间
- [ ] 截止时间
- [ ] 操作（查看/编辑）

---

#### `/workspace/merchant/tasks/create` - 创建任务（草稿）
**页面说明**：
- 商家在此页面创建任务草稿
- **不设置积分分配**（由绑定的服务商设置）
- **不扣除积分**（发布时才扣除）
- 创建后的营销活动状态为"草稿（DRAFT）"
- 服务商可以在后台设置积分分配并发布任务

**页面表单**：

---

**1. 基本信息**：

- **任务标题** *（必填）
  - 类型：文本输入
  - 限制：2-50字
  - 占位符："例如：新品体验推广"
  - 验证：不能为空，不能包含特殊字符

- **活动要求** *（必填）
  - 类型：富文本编辑器
  - 功能：支持文字、图片、链接、列表
  - 限制：10-5000字
  - 占位符："详细描述活动要求、发布规范、注意事项等..."
  - 提示：可以上传示例图片供达人参考

- **平台类型** *（必填，多选）
  - 类型：复选框组
  - 选项：
    - [ ] 小红书
    - [ ] 抖音
    - [ ] 微信朋友圈
    - [ ] 公众号
    - [ ] 微博
    - [ ] B站
    - [ ] 快手
  - 验证：至少选择1个
  - 提示：达人可以选择任意一个平台完成

---

**2. 赏金设置**（单位：元）

- **任务金额** *（必填）
  - 类型：数字输入
  - 范围：1-10000元
  - 占位符："例如：100"
  - 验证：必须 > 0
  - 说明：这是每条任务的预算总额，由服务商分配给达人、员工返佣、服务商

- **提示**：
  - 💡 服务商发布任务时会将此金额分配为：达人收入 + 员工邀请返佣 + 服务商收入
  - 💡 建议：单个任务预算100元，达人获得70-90元，员工返佣5-10元，服务商10-20元

---

**3. 任务配置**

- **活动名额** *（必填）
  - 类型：数字输入
  - 范围：1-1000人
  - 占位符："例如：10"
  - 默认值：10
  - 验证：必须 ≥ 1

- **接任务截止时间** *（必填）
  - 类型：日期时间选择器
  - 限制：最早选择当前时间+24小时
  - 默认值：7天后
  - 格式：YYYY-MM-DD HH:mm
  - 说明：达人可以接任务的最后期限（可延长）

- **提交截止时间** *（必填）
  - 类型：日期时间选择器
  - 限制：必须在"接任务截止时间"之后或同时
  - 默认值：接任务截止时间 + 7天
  - 格式：YYYY-MM-DD HH:mm
  - 说明：已接任务达人必须提交任务内容的最后期限（可延长）

- **任务有效期**（自动显示）
  - 类型：只读文本
  - 计算：接任务截止时间 - 当前时间
  - 示例："7天（可接任务），提交截止时间：2024-02-01 18:00"

---

**4. 高级设置**（默认折叠，点击"展开高级设置"显示）

- **指定服务商**（可选）
  - 类型：下拉选择
  - 选项：已绑定的服务商列表
  - 默认值：不指定（所有服务商可见）
  - 说明：指定后只有该服务商可以审核

- **达人等级要求**（可选）
  - 类型：复选框组
  - 选项：
    - [ ] UGC（普通用户）
    - [ ] KOC（关键意见消费者）
    - [ ] INF（达人）
    - [ ] KOL（关键意见领袖）
  - 默认值：全部选中（不限制）
  - 说明：未选中的等级无法接任务

- **备注**（可选）
  - 类型：多行文本框
  - 限制：0-500字
  - 占位符："其他需要说明的信息..."

---

**5. 表单操作**

- [ ] 保存草稿按钮
  - 行为：保存任务为草稿状态（DRAFT）
  - 验证：验证所有必填项
  - 成功提示："任务草稿已创建，等待服务商设置积分分配后发布"
  - 后续：跳转到任务列表页面

- [ ] 预览按钮
  - 行为：打开新标签页，展示任务在达人端的显示效果

- [ ] 取消按钮
  - 行为：放弃当前编辑，返回任务列表

---

**6. 表单验证规则**

- 前端实时验证：
  - 必填项不能为空
  - 数字范围验证（任务金额：1-10000元）
  - 日期逻辑验证（发布截止时间 < 任务截止时间）

- 后端验证：
  - 商家账号状态是否正常
  - 服务商是否存在（如果指定了）
  - **不验证积分余额**（发布时才验证）

---

**7. 创建成功后的流程**

创建成功后，营销活动状态为"草稿（DRAFT）"：
- 任务出现在商家的"草稿活动"列表中
- 商家可以编辑草稿活动
- 绑定的服务商可以在后台看到草稿活动
- 服务商可以为草稿活动设置积分分配并发布

---

## 📌 第二步：服务商设置积分分配并发布任务

**页面路由**：`/workspace/service-provider/tasks/:id/allocate`

**触发条件**：服务商在任务列表中看到商家创建的草稿活动，点击"设置积分并发布"

**页面说明**：
- 服务商在此页面为草稿活动设置积分分配
- 三部分收入总和必须等于商家填写的"任务金额"
- 设置后扣除商家积分（预扣款模式）：`balance -= 总费用`，`frozen_balance += 总费用`
- 营销活动状态变为"开放中（OPEN）"
- 生成任务分享链接供达人接任务

**⚠️ 预扣款模式说明**：
- 发布时立即扣除全部费用到 frozen_balance（如：10人×100元=1000元）
- 达人每次接任务时，frozen_balance 减少任务金额
- 任务手动关闭或名额已满时，未完成名额的费用自动退回 balance

---

**页面布局**：

```
┌─────────────────────────────────────────────┐
│ 任务基本信息（只读显示）                      │
├─────────────────────────────────────────────┤
│  任务标题：新品体验推广                       │
│  任务金额：100元                       │
│  活动名额：10人                              │
│  总费用：1000元                              │
├─────────────────────────────────────────────┤
│ 积分分配设置                                 │
├─────────────────────────────────────────────┤
│  达人收入 *（必填）                           │
│  [  80  ] 元                                 │
│                                             │
│  员工邀请返佣 *（必填）                       │
│  [  10  ] 元                                 │
│  💡 当被邀请的达人完成任务时，邀请者获得此返佣 │
│                                             │
│  服务商收入 *（必填）                         │
│  [  10  ] 元（自动计算）                     │
│                                             │
│  分成比例：达人 80% | 员工 10% | 服务商 10%  │
├─────────────────────────────────────────────┤
│ 验证提示                                     │
│  ✅ 总和 = 100元（符合任务金额）        │
│  ⚠️ 商家当前余额：5000元，足够发布            │
├─────────────────────────────────────────────┤
│ [取消]  [设置并发布任务]                     │
└─────────────────────────────────────────────┘
```

---

**表单字段**：

**1. 活动信息预览**（只读）

- **任务标题**
  - 显示商家填写的任务标题

- **平台类型**
  - 显示商家选择的平台列表

- **任务金额**
  - 显示商家填写的金额（不可修改）

- **活动名额**
  - 显示招募人数

- **总费用**
  - 自动计算：任务金额 × 活动名额

---

**2. 积分分配设置**（必填）

- **达人收入** *（必填）
  - 类型：数字输入
  - 范围：1-任务金额-2
  - 占位符："例如：80"
  - 验证：必须 > 0
  - 说明：达人完成任务后获得的收入

- **员工邀请返佣** *（必填）
  - 类型：数字输入
  - 范围：0-任务金额-1
  - 占位符："例如：10"
  - 默认值：0
  - 验证：必须 ≥ 0
  - 说明：
    - 当被邀请的达人完成任务时，邀请该达人注册的服务商员工获得此返佣
    - 如果达人是自己注册（无邀请者），此部分收入归服务商所有
    - 可以设置为0（不提供邀请返佣）

- **服务商收入** *（必填，自动计算）
  - 类型：数字输入（只读）
  - 计算：任务金额 - 达人收入 - 员工邀请返佣
  - 验证：必须 ≥ 0
  - 说明：
    - 固定归服务商的收入
    - 如果达人无邀请者，员工邀请返佣部分也归服务商

- **分成比例**（自动显示）
  - 类型：百分比显示
  - 示例：达人 80% | 员工 10% | 服务商 10%
  - 实时更新

---

**3. 实时验证提示**

- **总和验证**
  - ✅ 达人收入 + 员工返佣 + 服务商收入 = 任务金额
  - ❌ 总和不等于任务金额，请调整

- **余额验证**
  - ✅ 商家余额充足（XXX元 > 总费用）
  - ❌ 商家余额不足（XXX元 < 总费用），请联系商家充值

---

**4. 表单操作**

- [ ] 取消按钮
  - 行为：返回任务列表，不保存

- [ ] 设置并发布按钮
  - 行为：
    1. 验证积分分配总和是否正确
    2. 验证商家积分余额是否充足
    3. 扣除商家积分（任务金额 × 活动名额）
    4. 将扣除的积分存入商家的 frozen_balance（预扣款）
    5. 营销活动状态从 DRAFT 变为 OPEN
    6. 生成任务分享链接
    7. 通知商家任务已发布
  - 成功提示："任务已发布！分享链接已生成"

---

**5. 发布成功后的流程**

- 任务出现在达人任务大厅
- 服务商可以复制分享链接推广
- 达人可以点击链接或在大厅中接任务
- 商家可以在后台查看任务进度

---

**6. 积分分配示例**

**示例1：标准分配（100元任务）**
- 达人收入：80元
- 员工邀请返佣：10元
- 服务商收入：10元
- 分成：达人 80% | 员工 10% | 服务商 10%

**示例2：无员工返佣（100元任务）**
- 达人收入：85元
- 员工邀请返佣：0元
- 服务商收入：15元
- 分成：达人 85% | 员工 0% | 服务商 15%

**示例3：高返佣激励（100元任务）**
- 达人收入：70元
- 员工邀请返佣：20元
- 服务商收入：10元
- 分成：达人 70% | 员工 20% | 服务商 10%

---

**7. 发布成功后的操作**

**创建任务名额记录**：
- 发布成功后，系统自动创建 quota 个任务名额记录
- 每个名额的状态为 OPEN（开放中）
- 名额编号从1开始递增（task_slot_number = 1, 2, 3...）

**示例**：
```sql
-- 发布10人任务时，自动执行：
INSERT INTO tasks (campaign_id, task_slot_number, status) VALUES
('campaign-1', 1, 'OPEN'),
('campaign-1', 2, 'OPEN'),
('campaign-1', 3, 'OPEN'),
...
('campaign-1', 10, 'OPEN');
```

**多维表视图**：
- 商家/服务商可以在多维表中看到所有10个名额
- 所有名额初始状态都是 🟢开放中（OPEN）
- 达人可以通过邀请码接任务，将名额状态改为 🟡已分配（ASSIGNED）

---

**8. 生成员工专属邀请码**

**功能说明**：
- 营销活动发布后，服务商可以为每个员工生成专属邀请码
- 员工将邀请码分享给达人，达人通过邀请码接任务
- 系统追踪是哪个员工的邀请码被使用，用于绩效考核

**邀请码管理页面**：`/workspace/service-provider/campaigns/:id/invite-codes`

**页面功能**：
- 查看所有邀请码列表（固定邀请码）
- 为新员工自动生成专属邀请码
- 查看每个邀请码的使用情况
- 禁用/启用邀请码

**邀请码列表（固定邀请码模式）**：
```
┌─────────────────────┬──────────┬──────────┬──────────┬────────┐
│ 邀请码              │ 所属员工 │ 已使用   │ 状态     │ 操作   │
├─────────────────────┼──────────┼──────────┼──────────┼────────┤
│ CR_LI_SI_123456     │ 李四     │ 3        │ ✅启用   │[禁用] │
│ CR_WANG_WU_789012   │ 王五     │ 5        │ ✅启用   │[禁用] │
│ CR_ZHAO_LIU_345678  │ 赵六     │ 8        │ ❌禁用   │[启用] │
└─────────────────────┴──────────┴──────────┴──────────┴────────┘
```

**操作流程（固定邀请码模式）**：
1. 服务商管理员添加新员工到系统
2. 系统自动为该员工生成专属固定邀请码
3. 邀请码格式：`CR_{员工姓名拼音}_{6位随机码}`
4. 员工可以在自己的工作台看到专属邀请码
5. 员工复制邀请码分享给达人（可打印在名片上）
6. 邀请码永久有效，除非管理员手动禁用

---

**9. 截止时间延长**

**功能说明**：
- 服务商可以延长任务的截止时间
- 延长后任务继续开放接任务
- 已接任务的达人不受影响
- 由于采用预扣款模式，延长截止时间无需额外扣款

**使用场景**：
1. 任务原定截止时间到，但招募人数未满
2. 服务商希望继续招募更多达人
3. 延长截止时间后，任务恢复为"进行中"状态

**示例流程**：
- 活动名额：10人，任务金额100积分，总计1000元
- 原定截止时间：2月1日 18:00
- 2月1日到期时：任务标记为"已到期"但不会关闭
- 当前已接任务：3人
- 服务商在2月2日延长截止时间到2月5日 18:00
- 任务恢复为"开放中（OPEN）"状态，达人可以继续接任务
- frozen_balance 中的预扣款（700元）继续冻结

**操作按钮**：
- 位置：任务详情页 / 服务商任务列表
- 条件：营销活动状态为 OPEN 或已过期
- 点击后：弹出延长截止时间选择器

---

## 📌 第三步：达人通过邀请码接任务

**触发条件**：
- 员工将邀请码分享给达人（如微信群、朋友圈）
- 达人输入邀请码或点击邀请链接

**页面路由**：`/creator/accept-task?code=TASK_A1B2C3`

**页面说明**：
- 达人输入或点击邀请码后，显示任务详情
- 达人确认接任务
- 系统将一个OPEN状态的名额分配给该达人（状态变为ASSIGNED）
- 记录邀请追踪信息（哪个员工的邀请码）

**接任务流程**：

1. **验证邀请码**（固定邀请码模式）：
   ```sql
   SELECT * FROM invite_codes
   WHERE code = 'TASK_A1B2C3'
     AND type = 'CAMPAIGN_TASK'
     AND status = 'active';
   ```

2. **查找可用的OPEN名额**：
   ```sql
   SELECT * FROM tasks
   WHERE campaign_id = ?
     AND status = 'OPEN'
   ORDER BY task_slot_number
   LIMIT 1;
   ```

3. **更新任务名额状态**：
   ```sql
   UPDATE tasks
   SET status = 'ASSIGNED',
       creator_id = ?,        -- 达人ID
       assigned_at = NOW(),
       invite_code_id = ?,    -- 邀请码ID
       invited_by_staff_id = ? -- 员工ID
   WHERE id = ?;              -- 任务名额ID
   ```

4. **更新邀请码使用次数**：
   ```sql
   UPDATE invite_codes
   SET used_count = used_count + 1
   WHERE id = ?;
   ```

5. **记录邀请码使用**：
   ```sql
   INSERT INTO invite_code_usages (invite_code_id, creator_id, used_at)
   VALUES (?, ?, NOW());
   ```

**多维表视图变化**：

接任务前：
```
┌──────┬──────────┬─────────┬──────────┐
│ 编号 │ 状态     │ 接单达人│ 邀请人   │
├──────┼──────────┼─────────┼──────────┤
│ #1   │ 🟢开放中 │ -       │ -        │
│ #2   │ 🟢开放中 │ -       │ -        │
│ #3   │ 🟢开放中 │ -       │ -        │
└──────┴──────────┴─────────┴──────────┘
```

张三通过李四的邀请码接任务后：
```
┌──────┬──────────┬─────────┬──────────┐
│ 编号 │ 状态     │ 接单达人│ 邀请人   │
├──────┼──────────┼─────────┼──────────┤
│ #1   │ 🟡已分配 │ 张三    │ 李四     │
│ #2   │ 🟢开放中 │ -       │ -        │
│ #3   │ 🟢开放中 │ -       │ -        │
└──────┴──────────┴─────────┴──────────┘
```

**成功提示**：
- "任务已接单！请按时提交任务内容"
- 显示提交截止时间
- 添加到"我的任务"列表

---

**错误处理**：

- ❌ 邀请码无效："邀请码不存在或已过期"
- ❌ 邀请码已用完："邀请码使用次数已达上限"
- ❌ 活动已满："活动名额已满，无法接任务"
- ❌ 已接过该任务："您已经接过了这个任务"

---

#### `/workspace/merchant/tasks/:id` - 任务详情
**页面信息展示**：
- [ ] 任务基本信息
- [ ] 赏金分配
- [ ] 任务进度（已接名额数/总名额）
- [ ] 接任务达人列表
  - 达人昵称
  - 接任务时间
  - 提交状态
  - 审核状态
  - 操作（查看提交内容）

---

#### `/workspace/merchant/finance` - 财务管理
**页面功能**：
- [ ] 账户余额显示
  - 积分余额
  - 冻结积分（任务托管中）
  - 可用余额
- [ ] 充值按钮
- [ ] 交易记录列表
  - 时间
  - 类型（充值/发布任务/任务退款）
  - 金额
  - 余额
  - 备注

---

#### `/workspace/merchant/finance/recharge` - 充值
**页面表单**：
- [ ] 充值金额 *（必填）
  - 预设选项：100元、500元、1000元、5000元、自定义
- [ ] 支付方式 *（必填）
  - [ ] 微信支付
  - [ ] 支付宝
- [ ] 确认充值按钮

**✅ 充值流程**：
- 点击充值 → 选择金额（100/500/1000/5000/自定义）→ 跳转第三方支付 → 支付完成 → 回调更新余额
- **充值比例**：1元 = 1积分（无充值套餐、无优惠活动）

---

#### `/workspace/merchant/staff` - 员工管理（多维表）

**页面类型**：员工管理多维表（商家视图）

**核心功能**：
- ✅ 员工多维表展示（员工姓名、职能名称、权限数量、最近登录时间、状态）
- ✅ 添加员工（弹出表单，输入手机号）
- ✅ 权限管理抽屉（按模块分组显示权限、勾选授予/撤销）
- ✅ 批量操作（批量授予权限、批量撤销权限、批量停用）
- ✅ 权限组管理（查看权限组、创建权限组、编辑权限组）

**权限说明**：
- 仅限管理员可以查看和操作
- 可以给员工分配权限（不能超过管理员自己的权限）

**实现参考**：详见 PRD.md 第 1471 行"员工管理多维表"章节

**✅ 权限管理方式**：
- **不使用权限组**，采用最小颗粒度权限管理
- 管理员可以勾选任意权限组合给员工
- 权限按模块分类显示（任务、财务、用户、数据等）
- 员工最高权限可以等于管理员（可勾选所有权限）

**权限分类**：
- 任务权限：查看、创建、编辑、删除、审核
- 财务权限：查看、充值、提现审核
- 用户权限：查看员工、添加员工、编辑权限
- 数据权限：查看报表、导出数据

---

### 4. 服务商工作台 `/workspace/service-provider`

#### `/workspace/service-provider` - 服务商首页 Dashboard
**页面布局**：
```
┌─────────────────────────────────────────────┐
│ 顶部导航：                                  │
│ [角色切换器] 服务商 | 积分余额 | 提现 | 用户菜单│
├─────────────────────────────────────────────┤
│ 左侧菜单 │ 主内容区                       │
│         │                                 │
│ 首页    │ [统计卡片]                      │
│ 商家管理│ - 绑定商家数                    │
│ 达人管理│ - 绑定达人数                    │
│ 任务审核│ - 待审核任务：N 条 → 立即审核  │
│ 财务    │ - 本月收入                      │
│         │                                 │
└─────────────────────────────────────────────┘
```

---

#### `/workspace/service-provider/merchants` - 商家管理
**页面功能**：
- [ ] 绑定商家列表
- [ ] 查看商家详情
  - 商家基本信息
  - 合作任务数
  - 总交易额

**✅ 商家-服务商关系说明**：
- 商家和服务商是**多对多绑定关系**
- **不提供解绑功能**（如需终止合作，联系超级管理员处理）
- 一个服务商可以绑定多个商家
- 一个商家可以绑定多个服务商

---

#### `/workspace/service-provider/creators` - 达人管理（多维表）

**页面类型**：达人管理多维表（服务商视图）

**核心功能**：
- ✅ 达人多维表展示（头像、昵称、等级、粉丝数、完成任务数、完成率、邀请人）
- ✅ 筛选默认条件：只显示由本服务商员工邀请的达人（`inviter_id IN (本服务商员工)`）
- ✅ 多维度筛选（按等级、状态、完成率、标签、邀请人）
- ✅ 批量打标签（仅自己可见）
- ✅ 详情抽屉（查看达人完整档案、任务列表）
- ✅ 邀请统计（累计邀请、已激活数量、邀请返佣总额）

**限制**：
- ❌ 不能修改达人等级
- ❌ 不能封禁达人

**实现参考**：详见 PRD.md 第 1188 行"达人管理多维表"章节

---

#### `/workspace/service-provider/staff` - 员工管理（多维表）

**页面类型**：员工管理多维表（服务商视图）

**核心功能**：
- ✅ 员工多维表展示（姓名、职能、权限数量、邀请达人数量、完成任务总数、产生GMV、状态）
- ✅ 权限管理抽屉（按模块分组显示权限、勾选授予/撤销）
- ✅ 绩效追踪图表（邀请达人趋势、完成任务趋势、GMV趋势）
- ✅ 批量操作（批量授予权限、批量撤销权限、批量停用）
- ✅ 绩效排行榜（员工之间绩效对比）

**实现参考**：详见 PRD.md 第 1471 行"员工管理多维表"章节

---

#### `/workspace/service-provider/creators/invite` - 邀请达人（核心功能）

**⚠️ 重要**：这是服务商员工的核心功能，类似于原系统的"达人组长"角色
**✅ 使用统一邀请码系统**

**页面布局**：

```
┌─────────────────────────────────────────────┐
│ 邀请达人统计                                 │
├─────────────────────────────────────────────┤
│  累计邀请：15人                             │
│  已激活：12人（完成过任务）                 │
│  本月邀请：3人                              │
│  返佣收入：¥1,250.00                        │
├─────────────────────────────────────────────┤
│ 我的邀请码（固定）                           │
├─────────────────────────────────────────────┤
│  CR_ZHANG_SAN_123456                        │
│                                             │
│  [复制邀请码]  [生成二维码]  [查看详情]      │
├─────────────────────────────────────────────┤
│ 返佣规则说明                                 │
├─────────────────────────────────────────────┤
│  1. 分享邀请码给潜在达人                     │
│  2. 达人注册时填写邀请码                     │
│  3. 达人通过邀请码注册后成为"你的达人"       │
│  4. 当"你的达人"完成任务时，你获得返佣       │
│  5. 返佣金额由服务商在发布任务时设置         │
│  6. 返佣实时到账，可随时提现                 │
├─────────────────────────────────────────────┤
│ 邀请记录                                     │
├─────────────────────────────────────────────┤
│  达人昵称    | 等级  | 注册时间   | 返佣总额  │
│  小明        | KOC   | 2024-01-15 | ¥120.00  │
│  小红        | INF   | 2024-01-20 | ¥350.00  │
│  小刚        | UGC   | 2024-02-01 | ¥0.00    │
│  ...                                       │
└─────────────────────────────────────────────┘
```

---

**功能说明**：

**1. 我的邀请码（固定邀请码）**
- 每个服务商员工自动生成一个专属固定邀请码
- 邀请码格式：`CR_{员工姓名拼音}_{6位随机码}`（如：`CR_ZHANG_SAN_123456`）
- 达人注册时填写此邀请码后，员工成为该达人的"邀请者"
- 关联关系永久保存
- 邀请码类型：`INVITE_CREATOR`
- 邀请码永久有效，可重复使用（除非管理员禁用）
- 可打印在名片、海报等宣传材料上

**2. 分享邀请码**
- [ ] **复制邀请码**（必选）
  - 一键复制邀请码
  - 可分享到微信、QQ、短信等

- [ ] **生成二维码**（必选）
  - 生成包含邀请码的二维码图片
  - 可下载保存
  - 达人扫码后自动填写邀请码

- [ ] **邀请海报**（可选）
  - 自动生成带邀请码/二维码的宣传海报
  - 包含平台介绍、返佣说明
  - 可分享到社交媒体

**3. 邀请码管理**
- [ ] 查看邀请码详情
  - 邀请码
  - 使用次数统计
  - 创建时间
  - 状态（有效/禁用）
- [ ] 禁用/启用邀请码
- [ ] 查看使用记录（谁使用了这个邀请码）

**4. 返佣机制**

- **返佣来源**：
  - 服务商在发布任务时设置"员工邀请返佣"
  - 从"任务金额"中分配

- **返佣触发条件**：
  - 被邀请的达人完成任务
  - 任务通过审核
  - 系统自动结算返佣给员工

- **返佣到账**：
  - 实时到账
  - 记录在员工的"邀请返佣"收入类型中
  - 可以提现

- **无邀请者情况**：
  - 如果达人是自己注册（无邀请者）
  - "员工邀请返佣"部分归服务商所有

**4. 邀请记录列表**

**列表字段**：
- [ ] 达人昵称
- [ ] 达人等级
- [ ] 注册时间
- [ ] 完成任务数
- [ ] 返佣总额
- [ ] 操作（查看详情、解除关联）

**筛选和排序**：
- [ ] 按注册时间排序
- [ ] 按返佣金额排序
- [ ] 筛选已激活/未激活达人

---

**5. 邀请返佣示例**

**场景**：服务商发布任务，设置积分分配为：
- 达人收入：80元
- 员工邀请返佣：10元
- 服务商收入：10元
- 总计：100元

**达人A完成任务后**：
- 达人A获得：80元
- 如果达人A是通过员工小明的链接注册的：
  - 小明获得：10元（邀请返佣）
  - 服务商获得：10元
- 如果达人A是自己注册的：
  - 员工返佣：0元
  - 服务商获得：20元（10+10）

---

**6. 邀请统计面板**

**统计卡片**：
- [ ] 累计邀请达人数量
- [ ] 已激活达人数量（完成过任务的）
- [ ] 本月邀请达人数量
- [ ] 邀请返佣总收入
- [ ] 本月返佣收入

**图表**（可选）：
- [ ] 邀请趋势图（按日期/月份）
- [ ] 返佣收入趋势图
- [ ] 达人等级分布饼图

---

**7. 权限说明**

- **服务商员工**：
  - 只能看到"自己邀请的达人"
  - 只能看到"自己的返佣收入"
  - 不能查看其他员工的邀请数据

- **服务商管理员**：
  - 可以查看所有员工的邀请数据
  - 可以看到总的邀请统计
  - 可以调整返佣比例

---

**⚠️ 邀请关系规则**（已确认）：

1. **邀请关系**：永久关联
   - 达人永久属于邀请者
   - 邀请者持续获得该达人完成任务后的返佣

2. **解除邀请关系**：允许解除
   - 员工可以解除与达人的关联
   - 解除后不再获得该达人的返佣
   - 达人可以重新被其他员工邀请
   - 解除操作需要二次确认

3. **邀请码机制**：单独讨论
   - 邀请码/邀请链接的设计将单独确认

---

#### `/workspace/service-provider/tasks` - 任务审核
**页面功能**：
- [ ] 待审核任务列表（Tab 切换）
  - [ ] 待审核
  - [ ] 已通过
  - [ ] 已拒绝
- [ ] 筛选（按商家、达人）
- [ ] 审核操作

**列表字段**：
- [ ] 任务标题
- [ ] 商家名称
- [ ] 达人昵称
- [ ] 提交时间
- [ ] 任务赏金
- [ ] 提交内容（预览）
- [ ] 操作（通过/拒绝）

**审核弹窗**：
- [ ] 显示活动要求
- [ ] 显示达人提交的内容
  - 文字描述
  - 图片
  - 链接
- [ ] 审核结果选择
  - [ ] 通过（必填审核意见）
  - [ ] 拒绝（必填拒绝理由）
- [ ] 提交按钮

---

#### `/workspace/service-provider/finance` - 财务管理
**页面功能**：
- [ ] 账户余额
- [ ] 收入统计
  - 今日收入
  - 本月收入
  - 总收入
- [ ] 收入明细列表
  - 活动ID
  - 商家名称
  - 达人名称
  - 收入金额
  - 收入时间

---

#### `/workspace/service-provider/finance/withdraw` - 提现
**页面表单**：
- [ ] 提现金额 *（必填，不能超过可用余额）
- [ ] 提现方式 *（必填）
  - [ ] 支付宝
  - [ ] 银行卡
- [ ] 收款信息 *（必填）
  - 支付宝账号 / 银行卡号
  - 收款人姓名
- [ ] 提交申请

**✅ 提现规则**：
- 最低提现金额：100 元
- 提现手续费：免费
- 到账时间：1-3 个工作日
- 提现审核：超级管理员审核后打款
- 每日提现次数：无限制

---

### 5. 达人工作台 `/workspace/creator`

#### `/workspace/creator` - 达人首页 Dashboard
**页面布局**：
```
┌─────────────────────────────────────────────┐
│ 顶部导航：                                  │
│ [角色切换器] 达人 | 积分余额 | 提现 | 用户菜单│
├─────────────────────────────────────────────┤
│ 左侧菜单 │ 主内容区                       │
│         │                                 │
│ 首页    │ [统计卡片]                      │
│ 任务大厅│ - 我的收入                      │
│ 我的任务│ - 进行中：N 个                    │
│ 财务    │ - 待审核：N 个                    │
│         │ - 已完成：N 个                    │
│         │                                 │
│         │ [快捷操作]                      │
│         │ 前往任务大厅 →                  │
│         │                                 │
└─────────────────────────────────────────────┘
```

**⚠️ 重要**：
- **任务大厅**：浏览所有可接任务的"任务"
- **我的任务**：查看我的所有"活动名额"

---

#### `/workspace/creator/tasks/hall` - 任务大厅
**页面功能**：
- [ ] 任务列表（卡片展示）
- [ ] 筛选器
  - [ ] 任务赏金范围（滑块）
  - [ ] 达人等级要求
  - [ ] 营销活动状态（可接任务）
- [ ] 搜索任务
- [ ] 抢单按钮（点击直接接任务）

**任务卡片信息**：
- [ ] 商家 Logo 和名称
- [ ] 任务标题
- [ ] 活动要求（预览）
- [ ] 任务赏金（高亮显示）
- [ ] 剩余名额（如：剩 5/10 个名额）
- [ ] 截止时间
- [ ] 抢单按钮

**✅ 任务列表加载方式**：
- 采用**无限滚动**（下拉加载更多）
- 首次加载20个任务，滚动到底部自动加载下一批
- 提供筛选和排序功能优化查找

---

#### `/workspace/creator/tasks/my` - 我的任务
**页面功能**：
- [ ] 活动名额列表（Tab 切换）
  - [ ] 进行中（ASSIGNED）- 已接任务，待提交
  - [ ] 待审核（SUBMITTED）- 已提交，等待审核
  - [ ] 已完成（APPROVED）- 审核通过，已获得收入
  - [ ] 需重新提交（REJECTED）- 审核拒绝，需要修改后重新提交

**⚠️ 重要说明**：
- 此页面显示的是**活动名额**，不是"任务"
- 一个任务可以有多个活动名额（如果活动名额是10人，最多有10条活动名额）
- Tab显示的是活动名额的状态，不是任务的状态

**活动名额信息**：
- [ ] 任务标题（来自营销活动表）
- [ ] 任务赏金（来自营销活动表）
- [ ] 商家名称（来自营销活动表）
- [ ] 活动名额状态（ASSIGNED/SUBMITTED/APPROVED/REJECTED）
- [ ] 接任务时间
- [ ] 任务截止时间（来自营销活动表）
- [ ] 操作（提交/查看详情）

**操作按钮**：
- [ ] 提交平台链接/截图（仅进行中状态的活动名额）
- [ ] 重新提交（仅REJECTED状态的活动名额）
- [ ] 查看详情

**列表字段**：
- [ ] 任务标题
- [ ] 任务赏金
- [ ] 商家名称
- [ ] 状态
- [ ] 接任务时间
- [ ] 截止时间
- [ ] 操作（提交/查看）

---

#### `/workspace/creator/tasks/:id` - 活动名额详情
**页面信息**：
- [ ] 任务基本信息（只读，来自营销活动表）
  - 任务标题
  - 活动要求
  - 平台类型
  - 赏金明细
    - 达人收入
    - 服务商收入
- [ ] 活动名额信息
  - 接任务时间
  - 活动名额状态
  - 提交历史（如果有多次提交）
- [ ] 提交表单（仅ASSIGNED状态显示）
  - 平台链接 *（必填，URL 输入框）
  - 截图凭证 *（必填，图片上传）
  - 备注说明 *（可选）
  - 提交按钮

**⚠️ 重要说明**：
- 此页面显示的是**活动名额详情**，不是"任务详情"
- URL路由`/workspace/creator/tasks/:id`中的`:id`实际上是**活动名额ID**
- 一个任务可以有多个活动名额，每个活动名额都有自己的详情页
- 达人在任务大厅点击"接任务"后，会跳转到该活动名额的详情页

---

#### `/workspace/creator/tasks/:id/submit` - 提交平台链接/截图
**页面表单**：
- [ ] 平台链接 *（必填，URL 输入框）
  - 示例：小红书笔记链接、抖音视频链接、朋友圈分享截图等
  - 根据平台类型填写对应平台的发布链接
- [ ] 截图凭证 *（必填，图片上传，最多9张）
  - 必须包含发布内容的完整截图
  - 截图需清晰显示账号信息和发布内容
- [ ] 备注说明 *（可选，文本框）
  - 可补充说明发布情况或特殊要求
- [ ] 提交按钮

**提交规则**：
- 提交后服务商进行审核
- 审核通过后积分到账
- 审核不通过可重新提交

---

#### `/workspace/creator/finance` - 财务管理
**页面功能**：
- [ ] 账户余额显示
  - 积分余额
  - 冻结积分（审核中）
  - 可用余额
- [ ] 收入统计
  - 今日收入
  - 本月收入
  - 总收入
- [ ] 交易记录列表
  - 时间
  - 类型（任务收入/提现）
  - 金额
  - 备注

---

#### `/workspace/creator/finance/withdraw` - 提现
**页面表单**：
- [ ] 提现金额 *（必填）
- [ ] 提现方式 *（必填）
- [ ] 收款信息 *（必填）
- [ ] 提交申请

---

### 6. 通用页面

#### `/profile` - 个人中心
**页面功能**：
- [ ] 用户信息展示
  - 头像
  - 昵称
  - 手机号
  - **拥有的角色**（显示所有角色列表）
  - **当前激活的角色**（高亮显示）
  - 等级（如果是达人）
- [ ] **角色切换**（核心功能）
  - 显示所有拥有的角色
  - 点击角色后切换工作台
  - 示例：
    ```
    ┌────────────────────────────────┐
    │ 我的角色                        │
    ├────────────────────────────────┤
    │ ✅ 商家员工（当前）            │
    │   达人                        │
    │   服务商员工                  │
    └────────────────────────────────┘
    ```
- [ ] 修改个人信息
- [ ] 退出登录

**✅ 个人中心功能**：
- ✅ 绑定/解绑微信（通过auth-center）
- ✅ 通知设置（推送通知、邮件通知开关）
- ❌ 修改密码（通过auth-center，不在本系统中）
- ❌ 隐私设置（暂不需要）

---

## 🎨 页面交互设计

### 角色切换器

**位置**：顶部导航栏左侧

**功能**：显示当前角色，点击后展示所有拥有的角色，可以快速切换

**设计**：
```
┌─────────────────────────────────────┐
│ [▼ 商家员工] | 积分余额 | 充值 | ... │
└─────────────────────────────────────┘

点击后展开：
┌─────────────────────────────────────┐
│ ✅ 商家员工（当前）                  │
│    达人                             │
│    服务商员工                       │
└─────────────────────────────────────┘
```

**交互流程**：
1. 用户点击角色切换器
2. 展开角色列表（所有拥有的角色）
3. 当前角色高亮显示（带✅）
4. 点击其他角色后：
   - 调用API切换 `current_role`
   - 跳转到新角色的工作台
   - 显示切换成功提示

**API设计**：
```
POST /api/user/switch-role
Request: { "role": "CREATOR" }
Response: { "success": true, "workspace": "/workspace/creator" }
```

**样式**：
- 当前角色：高亮显示，左侧带 ✅
- 其他角色：普通样式
- Hover效果：背景色变化
- 图标：▼ 表示可展开

---

### 反馈提示

**请确认需要的提示类型**：

1. **加载状态**
- [ ] 页面加载动画
- [ ] 按钮加载状态（提交时）
- [ ] 列表加载骨架屏

2. **成功提示**
- [ ] Toast 提示（如"任务创建成功"）
- [ ] Modal 弹窗（如"充值成功"）

3. **错误提示**
- [ ] 表单验证错误
- [ ] API 错误提示
- [ ] 网络错误提示

4. **确认对话框**
- [ ] 删除确认
- [ ] 解绑确认

**✅ 提示类型规范**：

**1. 成功提示（Toast提示，3秒自动消失）**：
- ✅ 任务创建成功
- ✅ 任务提交成功
- ✅ 充值成功
- ✅ 权限更新成功
- ✅ 绑定成功

**2. 错误提示（Toast提示，5秒自动消失）**：
- ✅ 表单验证错误（如：必填项未填、格式错误）
- ✅ API错误提示（如：余额不足、邀请码无效）
- ✅ 网络错误提示（如：网络连接失败、请重试）
- ✅ 权限错误（如：无权限执行此操作）

**3. 确认对话框（Modal弹窗）**：
- ✅ 删除确认（删除任务、删除员工等）
- ✅ 重要操作确认（批量审核、批量删除）
- ✅ 二次确认（余额调整、账户封禁等敏感操作）

**4. 状态通知（Inbox通知列表）**：
- ✅ 任务审核结果（通过/拒绝，带审核意见）
- ✅ 提现审核结果（已打款/已拒绝）
- ✅ 任务截止提醒（即将到期）
- ✅ 新任务通知（有新的可接任务）

**5. 加载状态**：
- ✅ 列表加载骨架屏
- ✅ 按钮Loading状态
- ✅ 全屏Loading（页面切换时）

---

## 📝 问卷调查

### 关于任务

**1. 平台类型**（请确认/补充）：
- [ ] 小红书种草
- [ ] 抖音视频
- [ ] 微信朋友圈
- [ ] 快手视频
- [ ] 公众号推广
- [ ] 其他：______

**2. 任务素材**（请确认）：
- [ ] 图片（最多几张）：____ 张
- [ ] 视频（是否需要）：是 / 否
- [ ] 文案（是否需要）：是 / 否

**3. 任务审核标准**（请填写）：
- ⚠️ 达人提交的内容，审核通过/拒绝的标准是什么？
```
通过标准：
1. ____
2. ____
3. ____

拒绝标准：
1. ____
2. ____
3. ____
```

---

### 关于积分

**1. 充值金额**（请确认）：
- [ ] 100 元
- [ ] 500 元
- [ ] 1000 元
- [ ] 5000 元
- [ ] 自定义金额

**2. 充值优惠**（是否需要）：
- [ ] 充100送10
- [ ] 充500送60
- [ ] 充1000送150
- [ ] 不需要优惠

**3. 提现规则**（请确认）：
- 最低提现：______ 元
- 手续费：______ %
- 到账时间：______ 天

---

### 关于达人

**1. 达人等级标准**（请填写）：
```
UGC（普通用户）：
- 粉丝数：不限
- 要求：注册即可

KOC（关键意见消费者）：
- 粉丝数：______
- 要求：______

INF（达人）：
- 粉丝数：______
- 要求：______

KOL（关键意见领袖）：
- 粉丝数：______
- 要求：______
```

**2. 达人升级**（请确认）：
- [ ] 用户自主申请 → 审核 → 升级
- [ ] 系统自动根据粉丝数升级
- [ ] 管理员手动升级
- 📝 你的选择：______

---

## ✅ 待确认清单

### 请逐项确认以下问题：

**角色相关**：
- [ ] 6 种角色是否满足需求？是否需要增加/减少？
- [ ] 商家员工、服务商员工是否需要？还是只用管理员就够了？

**权限相关**：
- [ ] 权限矩阵是否符合预期？
- [ ] 是否需要更细的权限控制（如某个员工只能发布任务不能充值）？
- [ ] 权限组预设是否符合需求？

**页面相关**：
- [ ] 页面路由是否完整？是否有遗漏的页面？
- [ ] 每个页面的功能是否符合预期？
- [ ] 是否需要额外的页面？

**业务逻辑**：
- [ ] 任务流程是否符合预期？是否需要调整？
- [ ] 积分结算规则是否合理？
- [ ] 审核流程是否符合预期？

**其他**：
- [ ] 是否需要移动端适配？（手机浏览）
- [ ] 是否需要消息通知（站内信、短信）？
- [ ] 是否需要数据报表？

---

## 🗄️ 数据库设计

### 数据库规范（⚠️ 统一标准）

**1. 表命名规范**：
- ✅ 统一使用**复数形式**（如：users, tasks, campaigns, invitations）
- ✅ 使用下划线分隔小写单词（snake_case）
- ✅ 表名应简洁明了，表达实体含义

**2. 软删除策略**：
- ✅ **核心业务表**使用软删除（deleted_at字段）：
  - users, merchants, service_providers, creators, campaigns, tasks
- ✅ **关联表**使用硬删除：
  - permission_*（权限表）、invitation_records（邀请记录）
  - credit_transactions（交易流水，永久保留）
- ✅ **配置表**使用is_active字段：
  - invitation_codes（邀请码表）

**3. 时间字段规范**：
- ✅ 统一使用 `TIMESTAMP WITH TIME ZONE`（存储为UTC时区）
- ✅ 所有时间字段命名：*_at（如created_at, updated_at, deleted_at）
- ✅ 前端展示时转换时区，数据库统一使用UTC

**4. JSONB字段索引**：
- ✅ 所有JSONB字段必须添加GIN索引
- ✅ 已添加索引的表：users.roles, users.profile, campaigns.platforms, credit_transactions.metadata

**5. 字符编码**：
- ✅ 统一使用 UTF-8 编码
- ✅ VARCHAR字段默认长度：50（代码）、100（名称）、200（说明）、500（URL）

**6. 复合索引策略**（⚠️ 性能优化）：

**常用查询场景及对应的复合索引**：
```sql
-- tasks表：按活动+状态+邀请人筛选
CREATE INDEX idx_tasks_campaign_status_inviter
ON tasks(campaign_id, status, invited_by_staff_id);

-- tasks表：按活动+状态+创建时间排序
CREATE INDEX idx_tasks_campaign_status_created
ON tasks(campaign_id, status, created_at DESC);

-- credit_transactions表：按账户+时间范围查询
CREATE INDEX idx_transactions_account_time
ON credit_transactions(account_id, created_at DESC);

-- invitation_codes表：按码+状态查询
CREATE INDEX idx_invite_codes_code_active
ON invitation_codes(code, is_active);

-- invitation_records表：按邀请人+时间查询
CREATE INDEX idx_invitation_records_inviter_time
ON invitation_records(inviter_id, created_at DESC);

-- campaigns表：按状态+创建时间查询
CREATE INDEX idx_campaigns_status_created
ON campaigns(status, created_at DESC);
```

**7. 性能优化配置**（⚠️ 生产环境配置）：

**数据库连接池**：
- 最大连接数：100（根据实际负载调整）
- 最小空闲连接数：10
- 连接超时时间：30秒
- 空闲连接超时：600秒（10分钟）

**查询优化**：
- 慢查询阈值：1秒（超过1秒的查询记录到慢查询日志）
- 查询超时设置：30秒（默认）
- 批量操作大小限制：每次最多1000条

**缓存策略**：
- 用户session缓存：Redis，过期时间24小时
- 邀请码缓存：Redis，过期时间1小时
- 统计数据缓存：Redis，过期时间5分钟
- 热点数据缓存：Redis，过期时间1小时

**监控告警**：
- 数据库CPU使用率 > 80%：告警
- 数据库连接数 > 80%：告警
- 慢查询数量 > 10/分钟：告警
- API响应时间P99 > 1秒：告警

---

### 设计原则

- **表名**：使用复数形式，小写，单词间用下划线分隔（如：`users`, `campaigns`, `tasks`, `invite_codes`）
- **字段名**：小写，单词间用下划线分隔（如：`created_at`, `creator_id`）
- **主键**：统一使用 `id` 字段，类型为 `UUID`
- **时间戳**：所有表都包含 `created_at` 和 `updated_at`
- **软删除**：重要表使用 `deleted_at` 实现软删除

---

### 核心表结构

#### 1. 用户表 (users)

存储所有用户的基本信息，支持一个用户拥有多个角色。

**⚠️ 重要**：用户通过 **auth_center** (os.crazyaigc.com) 进行统一认证。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | VARCHAR(255) | ✅ | - | （内部使用） | 用户唯一标识（本地主键，CUID） |
| auth_center_user_id | UUID | ✅ | - | （内部使用） | **关联 auth_center_db.users.user_id（核心字段）** |
| nickname | VARCHAR(50) | ❌ | - | 昵称 | 用户显示名称（本地存储） |
| avatar_url | VARCHAR(500) | ❌ | - | 头像 | 头像URL（本地存储） |
| profile | JSONB | ❌ | - | 个人资料 | 用户配置信息 {nickname, avatarUrl, bio} |
| roles | JSONB | ✅ | '[]' | 拥有角色 | 用户拥有的所有角色，例如：["MERCHANT_ADMIN", "CREATOR"] |
| current_role | VARCHAR(50) | ❌ | - | 当前角色 | 用户当前激活的角色，用于工作台切换 |
| last_used_role | VARCHAR(50) | ❌ | - | 最后使用角色 | 记录用户最后使用的角色，用于登录时默认选择 |
| invited_by | VARCHAR(255) | ❌ | - | 邀请人ID | 邀请人的 users.id（用于返佣） |
| invitation_code_id | UUID | ❌ | - | 邀请码ID | 使用的邀请码ID（关联 invitation_codes.id） |
| status | VARCHAR(20) | ✅ | 'active' | 状态 | active: 正常, banned: 封禁, inactive: 未激活 |
| last_login_at | TIMESTAMP | ❌ | - | 最后登录时间 | - |
| last_login_ip | VARCHAR(50) | ❌ | - | 最后登录IP | - |
| created_at | TIMESTAMP | ✅ | NOW() | 注册时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |
| deleted_at | TIMESTAMP | ❌ | - | - | 软删除时间 |

**索引**：
- PRIMARY KEY (id)
- **UNIQUE INDEX (auth_center_user_id)** - 关联账号中心的核心索引
- INDEX (status)
- INDEX (current_role)
- INDEX (invited_by)
- INDEX (invitation_code_id)
- INDEX (created_at)
- **INDEX (roles)** - GIN索引用于JSONB查询
- **INDEX (profile)** - GIN索引用于JSONB查询（⚠️ 性能优化）

**外键**：
- FOREIGN KEY (auth_center_user_id) REFERENCES auth_center_db.users(user_id)
- FOREIGN KEY (invitation_code_id) REFERENCES invitation_codes(id)

**多角色用户的默认角色选择逻辑**：
1. **登录时默认角色选择**：
   - 优先使用`last_used_role`（记录用户最后使用的角色）
   - 如果`last_used_role`为空，按优先级选择：SUPER_ADMIN > MERCHANT_ADMIN > SERVICE_PROVIDER_ADMIN > MERCHANT_STAFF > SERVICE_PROVIDER_STAFF > CREATOR
2. **角色切换**：
   - 用户可以手动切换角色
   - 切换后更新`current_role`和`last_used_role`
3. **角色权限继承**：
   - 管理员权限包含所有员工权限
   - 员工最高权限不能超过管理员权限

**用户信息同步机制**：

**同步策略**：
1. **实时同步（Webhook回调）**：
   - auth-center在用户信息变更时主动通知PR Business
   - 适用场景：手机号变更、密码重置、账号封禁等关键操作
   - 实现方式：auth-center调用PR Business的Webhook接口

2. **定时同步（每小时一次）**：
   - PR Business定时从auth-center拉取用户信息
   - 适用场景：微信unionid更新、头像昵称同步等
   - 实现方式：后台定时任务

3. **登录时同步**：
   - 用户登录时，从auth-center获取最新用户信息
   - 适用场景：头像、昵称等非关键信息
   - 实现方式：在登录回调接口中同步

**Webhook接口定义**（PR Business提供）：
```http
POST /api/webhook/auth-center/user-update
Content-Type: application/json
X-Auth-Center-Signature: SHA256(请求体 + 密钥)

Request Body:
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_type": "phone_number_changed",  // phone_number_changed | password_changed | account_banned | account_deleted
  "old_value": "+86-13800138000",
  "new_value": "+86-13900139000",
  "timestamp": "2026-02-02T10:30:00Z"
}

Response:
{
  "success": true,
  "synced": true
}
```

**⚠️ 重要说明**：
- 手机号、密码等敏感信息只存储在auth-center，PR Business不存储
- PR Business通过`auth_center_user_id`外键关联，但不强依赖外键约束（考虑跨库）
- 建议使用应用层逻辑确保数据一致性，而非数据库外键

**与 auth-center 的关系**：
```sql
-- auth_center_db.users（账号中心）
CREATE TABLE users (
  user_id UUID PRIMARY KEY,
  union_id VARCHAR(255) UNIQUE,  -- 微信 unionid
  phone_number VARCHAR(255) UNIQUE,  -- 手机号（用于密码登录）
  password_hash VARCHAR(255),
  email VARCHAR(255) UNIQUE
);

-- pr_business_db.users（本地用户）
CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY,  -- 本地主键
  auth_center_user_id UUID UNIQUE NOT NULL,  -- 关联账号中心
  roles JSONB NOT NULL DEFAULT '[]',
  current_role VARCHAR(50)
);

-- 通过 auth_center_user_id 关联
SELECT
  pr_user.id AS local_user_id,
  pr_user.roles,
  pr_user.current_role,
  auth_user.phone_number,
  auth_user.union_id
FROM pr_business_db.users pr_user
JOIN auth_center_db.users auth_user
  ON pr_user.auth_center_user_id = auth_user.user_id
WHERE pr_user.id = 'xxx';
```

**角色说明**：
- `roles` 字段存储用户拥有的所有角色（JSON数组）
- `current_role` 字段存储用户当前激活的角色
- 用户可以在多个角色之间切换
- 切换角色后，工作台和权限根据 `current_role` 变化

**示例**：
```json
{
  "id": "clxxxxx",
  "auth_centerUserId": "550e8400-e29b-41d4-a716-446655440000",
  "nickname": "张三",
  "avatarUrl": "https://wx.qlogo.cn/xxx",
  "roles": ["SERVICE_PROVIDER_STAFF", "CREATOR"],
  "currentRole": "SERVICE_PROVIDER_STAFF",
  "invitedBy": "clyyyyy",
  "invitationCodeId": "inv-zzz"
}
```
表示张三既是服务商员工，也是达人，当前激活的是服务商员工角色，是通过邀请码邀请进来的。

**⚠️ 重要变更**：
- **不再存储手机号和密码**：手机号和密码存储在 auth_center
- **通过 auth_center_user_id 关联**：确保跨系统用户统一
- **profile 字段**：本地存储用户配置信息（昵称、头像等）
- **invited_by 和 invitation_code_id**：记录邀请关系，用于返佣计算

---

#### 2. 商家表 (merchants)

存储商家信息，一个用户可以同时是商家、达人等多个角色。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 商家唯一标识 |
| admin_id | UUID | ✅ | - | 管理员ID | 关联 users.id（商家管理员的用户ID） |
| provider_id | UUID | ✅ | - | 服务商ID | ⚠️ 关联 service_providers.id（创建该商家的服务商） |
| user_id | UUID | ✅ | - | 关联用户ID | 关联 users.id（一个用户可以有多个商家身份） |
| name | VARCHAR(100) | ✅ | - | 商家名称 | - |
| description | TEXT | ❌ | - | 商家简介 | - |
| logo_url | VARCHAR(500) | ❌ | - | 商家LOGO | - |
| industry | VARCHAR(50) | ❌ | - | 所属行业 | - |
| status | VARCHAR(20) | ✅ | 'active' | 状态 | active: 正常, suspended: 暂停, inactive: 未激活 |
| created_at | TIMESTAMP | ✅ | NOW() | 创建时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |
| deleted_at | TIMESTAMP | ❌ | - | - | 软删除时间 |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (admin_id)  -- 一个商家只有一个管理员
- INDEX (provider_id)  -- ⚠️ 关联服务商，用于查询服务商管理的所有商家
- INDEX (user_id)  -- 去掉 UNIQUE，一个用户可以关联多个商家
- INDEX (status)
- INDEX (name)

**外键**：
- FOREIGN KEY (admin_id) REFERENCES users(id)
- FOREIGN KEY (provider_id) REFERENCES service_providers(id)
- FOREIGN KEY (user_id) REFERENCES users(id)

**⚠️ 重要说明**：
- **创建流程**：商家必须由服务商通过邀请码创建（MERCHANT-{服务商ID后6位}）
- **绑定关系**：provider_id记录商家归属的服务商，绑定后永久保存
- **一个服务商可以管理多个商家**：通过provider_id关联查询
- **一个商家只能绑定一个服务商**：provider_id是单一值，不是多对多关系
- 创建商家时，`admin_id`自动设置为当前登录用户ID
- 管理员权限转移：未来可通过更新`admin_id`字段实现（需添加权限转移功能）

---

#### 3. 服务商表 (service_providers)

存储服务商信息，一个用户可以同时是服务商、达人等多个角色。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 服务商唯一标识 |
| admin_id | UUID | ✅ | - | 管理员ID | 关联 users.id（服务商管理员的用户ID） |
| user_id | UUID | ✅ | - | 关联用户ID | 关联 users.id（一个用户可以有多个服务商身份） |
| name | VARCHAR(100) | ✅ | - | 服务商名称 | - |
| description | TEXT | ❌ | - | 服务商简介 | - |
| logo_url | VARCHAR(500) | ❌ | - | 服务商LOGO | - |
| status | VARCHAR(20) | ✅ | 'active' | 状态 | active: 正常, suspended: 暂停, inactive: 未激活 |
| created_at | TIMESTAMP | ✅ | NOW() | 创建时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |
| deleted_at | TIMESTAMP | ❌ | - | - | 软删除时间 |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (admin_id)  -- 一个服务商只有一个管理员
- INDEX (user_id)  -- 去掉 UNIQUE，一个用户可以关联多个服务商
- INDEX (status)
- INDEX (name)

**外键**：
- FOREIGN KEY (admin_id) REFERENCES users(id)
- FOREIGN KEY (user_id) REFERENCES users(id)

**⚠️ 重要说明**：
- 创建服务商时，`admin_id`自动设置为当前登录用户ID
- 管理员权限转移：未来可通过更新`admin_id`字段实现（需添加权限转移功能）

---

#### 4. 达人表 (creators)

存储达人信息，一个用户可以同时是达人、商家、服务商等多个角色。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 达人唯一标识 |
| user_id | UUID | ✅ | - | 关联用户ID | 关联 users.id（一个用户可以有多个达人身份） |
| is_primary | BOOLEAN | ✅ | TRUE | 是否主账号 | ⚠️ 标识主达人账号（一个用户的第一个达人身份自动设为主账号） |
| level | VARCHAR(20) | ✅ | 'UGC' | 达人等级 | UGC: 普通用户, KOC: 关键意见消费者, INF: 达人, KOL: 关键意见领袖 |
| followers_count | INT | ✅ | 0 | 粉丝数 | 自行填写或系统获取 |
| wechat_openid | VARCHAR(100) | ❌ | - | 微信OpenID | 微信登录关联 |
| wechat_nickname | VARCHAR(100) | ❌ | - | 微信昵称 | - |
| wechat_avatar | VARCHAR(500) | ❌ | - | 微信头像 | - |
| inviter_id | UUID | ❌ | - | 邀请人ID | 邀请该达人的用户ID（可能是服务商员工、服务商管理员等） |
| inviter_type | VARCHAR(50) | ❌ | - | 邀请人类型 | PROVIDER_STAFF: 服务商员工, PROVIDER_ADMIN: 服务商管理员, OTHER: 其他 |
| inviter_relationship_broken | BOOLEAN | ✅ | FALSE | 邀请关系已解除 | TRUE: 已解除, FALSE: 有效 |
| status | VARCHAR(20) | ✅ | 'active' | 状态 | active: 正常, banned: 封禁, inactive: 未激活 |
| created_at | TIMESTAMP | ✅ | NOW() | 注册时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |
| deleted_at | TIMESTAMP | ❌ | - | - | 软删除时间 |

**索引**：
- PRIMARY KEY (id)
- INDEX (user_id)  -- 去掉 UNIQUE，一个用户可以有多个达人身份
- INDEX (level)
- INDEX (inviter_id)
- INDEX (status)
- INDEX (is_primary)  -- ⚠️ 用于查询用户的主达人账号
- UNIQUE INDEX (user_id, is_primary) WHERE is_primary = true  -- ⚠️ 确保每个用户只有一个主达人账号

**外键**：
- FOREIGN KEY (user_id) REFERENCES users(id)
- FOREIGN KEY (inviter_id) REFERENCES users(id)  -- 改为关联 users 表

**⚠️ 重要变更**：
- 一个用户可以作为达人身份注册多次（例如：不同微信号对应不同达人账号）
- `inviter_id` 改为关联 `users.id`，因为邀请人可能是任何角色（服务商员工、服务商管理员等）
- 添加 `inviter_type` 字段标识邀请人类型
- 添加 `is_primary` 字段标识主达人账号，默认第一个达人身份为主账号

**多达人身份的业务场景**：
- **场景1**：不同微信账号 - 合理（如个人号、工作号）
- **场景2**：不同平台专长 - 合理（小红书达人、抖音达人）
- **场景3**：测试账号 - 限制（最多3个达人账号，需要管理员审批）
- **主账号标识**：
  - 第一个达人身份自动设为主账号（is_primary=true）
  - 可手动切换主账号（确保同一时间只有一个is_primary=true）
  - 系统默认显示主账号的数据和统计

**达人等级标准**：
| 等级 | 粉丝数要求 | 任务完成率要求 | 最低完成任务数 | 说明 |
|------|-----------|---------------|--------------|------|
| UGC | 无要求 | 无要求 | 0 | 普通用户，注册即可 |
| KOC | 1K-10K | ≥60% | 3 | 关键意见消费者，有基本影响力 |
| INF | 10K-100K | ≥70% | 10 | 达人，有稳定粉丝基础 |
| KOL | 100K+ | ≥80% | 20 | 关键意见领袖，行业头部 |

**⚠️ 等级升级规则**：
- 等级由系统自动计算（每日凌晨更新）
- 同时满足粉丝数、完成率、任务数三个条件才升级
- 降级：完成率连续30天低于要求，降一级
- 手动调整：超级管理员可手动调整达人等级（用于特殊合作）

---

#### 5. 营销活动表 (campaigns)

存储营销活动信息（商家发布的活动）。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 营销活动唯一标识 |
| merchant_id | UUID | ✅ | - | 商家ID | 关联 merchants.id |
| provider_id | UUID | ❌ | - | 服务商ID | 关联 service_providers.id（指定服务商时） |
| title | VARCHAR(100) | ✅ | - | 活动名称 | - |
| requirements | TEXT | ✅ | - | 活动要求 | 富文本内容 |
| platforms | JSONB | ✅ | - | 平台类型 | ["xiaohongshu", "douyin", "weixin_moment", ...] |
| task_amount | INT | ✅ | - | 任务金额 | 每个活动名额的费用（单位：积分），例如 100 |
| campaign_amount | INT | ✅ | - | 活动总费用 | task_amount × quota（单位：积分），例如 1000 |
| creator_amount | INT | ❌ | - | 达人收入 | 服务商发布时设置，例如 80 |
| staff_referral_amount | INT | ❌ | - | 员工邀请返佣 | 服务商发布时设置，例如 10 |
| provider_amount | INT | ❌ | - | 服务商收入 | 服务商发布时设置，例如 10 |
| quota | INT | ✅ | - | 活动名额 | 招募人数，例如 10 |
| task_deadline | TIMESTAMP | ✅ | - | 接任务截止时间 | 达人可以接任务的最后期限（可延长） |
| submission_deadline | TIMESTAMP | ✅ | - | 提交截止时间 | 已接达人必须提交任务内容的最后期限（可延长） |
| status | VARCHAR(20) | ✅ | 'DRAFT' | 活动状态 | DRAFT: 草稿, OPEN: 开放中, CLOSED: 已关闭 |
| created_at | TIMESTAMP | ✅ | NOW() | 创建时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |
| deleted_at | TIMESTAMP | ❌ | - | - | 软删除时间 |

**索引**：
- PRIMARY KEY (id)
- INDEX (merchant_id)
- INDEX (provider_id)
- INDEX (status)
- INDEX (created_at)
- **INDEX (platforms)** - GIN索引用于JSONB查询（⚠️ 性能优化）

**外键**：
- FOREIGN KEY (merchant_id) REFERENCES merchants(id)
- FOREIGN KEY (provider_id) REFERENCES service_providers(id)

**营销活动状态机**：
- **DRAFT（草稿）**：商家创建的草稿活动，未发布
  - 所有任务名额状态必须为OPEN
  - 可操作：编辑、删除、发布（服务商设置积分分配后）
- **OPEN（开放中）**：服务商设置积分分配并发布，等待达人接任务
  - 允许达人接任务（ASSIGNED操作）
  - 可操作：延长截止时间、手动关闭
- **CLOSED（已关闭）**：名额已满或手动关闭（截止时间到不自动关闭）
  - 不允许新的ASSIGNED操作
  - 已ASSIGNED/SUBMITTED的名额可继续操作

**⚠️ 状态转换规则**：
1. DRAFT → OPEN：服务商设置积分分配并发布
2. OPEN → CLOSED：名额已满（已接名额数=quota）或商家手动关闭
3. CLOSED无法重新OPEN（如需重新开放，创建新活动）

**自动拒绝逻辑（定时任务）**：
```sql
-- 定时任务：每小时执行一次
-- 检查所有超过submission_deadline的ASSIGNED状态名额，自动拒绝

UPDATE tasks
SET status = 'REJECTED',
    audited_by = 'SYSTEM',
    audited_at = NOW(),
    reject_reason = '提交截止时间已到，自动拒绝',
    version = version + 1
WHERE status = 'ASSIGNED'
  AND submission_deadline < NOW()
  AND audited_by IS NULL;

-- ⚠️ 自动拒绝的业务规则：
-- 1. 只拒绝ASSIGNED状态（已接任务但未提交）
-- 2. 不影响SUBMITTED状态（已提交待审核的）
-- 3. 商家frozen_balance退回balance
-- 4. 不需要通知达人（可选）
-- 5. 达人无法申诉（系统自动拒绝）
-- 6. 服务商可以延长submission_deadline，避免自动拒绝
```

---

#### 6. 任务名额表 (tasks)

存储营销活动的所有名额记录（采用方案1：预分配所有名额）。

**⚠️ 重要设计决策**：
- **预分配模式**：营销活动发布时，创建 quota 个任务名额记录（如：10人×100元 = 创建10条tasks记录）
- **每个名额都是独立行**：便于多维表展示和操作，支持按状态筛选、批量操作
- **状态流转**：OPEN → ASSIGNED → SUBMITTED → APPROVED/REJECTED

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 任务名额唯一标识 |
| campaign_id | UUID | ✅ | - | 营销活动ID | 关联 campaigns.id |
| task_slot_number | INT | ✅ | - | 名额编号 | 活动内的编号（1, 2, 3...），方便定位和讨论 |

| **状态字段** ||||||
| status | VARCHAR(20) | ✅ | 'OPEN' | 状态 | OPEN: 开放中, ASSIGNED: 已分配, SUBMITTED: 待审核, APPROVED: 已完成, REJECTED: 已拒绝 |

| **达人信息** ||||||
| creator_id | UUID | ❌ | - | 达人ID | 接任务的达人ID（关联 creators.id），OPEN状态时为NULL |
| assigned_at | TIMESTAMP | ❌ | - | 接单时间 | 达人接任务的时间 |

| **提交信息** ||||||
| platform | VARCHAR(50) | ❌ | - | 发布平台 | 达人选择的平台（xiaohongshu, douyin等） |
| platform_url | VARCHAR(500) | ❌ | - | 平台链接 | 达人提交的发布链接 |
| screenshots | JSONB | ❌ | - | 截图凭证 | ["url1", "url2", ...] |
| submitted_at | TIMESTAMP | ❌ | - | 提交时间 | 达人提交任务的时间 |
| notes | TEXT | ❌ | - | 备注说明 | 达人或服务商的备注 |

| **审核信息** ||||||
| audited_by | UUID | ❌ | - | 审核人ID | 审核此任务名额的用户ID（服务商管理员或员工） |
| audited_at | TIMESTAMP | ❌ | - | 审核时间 | - |
| audit_note | TEXT | ❌ | - | 审核备注 | 服务商审核时填写的备注 |

| **邀请追踪** ||||||
| inviter_id | UUID | ❌ | - | 邀请人ID | 邀请该达人的用户ID（可能是服务商员工、商家员工等） |
| inviter_type | VARCHAR(50) | ❌ | - | 邀请人类型 | PROVIDER_STAFF: 服务商员工, MERCHANT_STAFF: 商家员工, OTHER: 其他 |

| **扩展字段** ||||||
| priority | VARCHAR(10) | ❌ | 'MEDIUM' | 优先级 | HIGH: 高, MEDIUM: 中, LOW: 低 |
| tags | TEXT[] | ❌ | - | 标签 | 自定义标签数组，如：["紧急", "VIP", "测试任务"] |

| **系统字段** ||||||
| version | INT | ✅ | 0 | 版本号 | ⚠️ 乐观锁，防止并发审核重复结算 |
| created_at | TIMESTAMP | ✅ | NOW() | 创建时间 | 名额创建时间（发布活动时创建所有名额） |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |

**索引**：
- PRIMARY KEY (id)
- INDEX (campaign_id)
- INDEX (status)
- INDEX (creator_id)
- INDEX (inviter_id)
- INDEX (audited_by)
- UNIQUE INDEX (campaign_id, task_slot_number) -- 确保每个活动内的名额编号唯一
- UNIQUE INDEX (campaign_id, creator_id) WHERE creator_id IS NOT NULL -- ⚠️ 防止同一达人重复接同一个活动的任务

**外键**：
- FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
- FOREIGN KEY (creator_id) REFERENCES creators(id)
- FOREIGN KEY (audited_by) REFERENCES users(id)
- FOREIGN KEY (inviter_id) REFERENCES users(id) -- 统一关联users表

**⚠️ 多维表核心**：
- 每个名额都是独立的一行，完美对应多维表的行
- 支持按状态筛选：`WHERE status = 'OPEN'` → 查看未分配名额
- 支持按邀请人筛选：`WHERE inviter_id = ? AND inviter_type = 'PROVIDER_STAFF'` → 查看某员工邀请的名额
- 支持批量操作：`UPDATE tasks SET status = 'APPROVED' WHERE id IN (...)`
- **一致性**：与creators表的邀请人字段设计保持一致（inviter_id + inviter_type）

---

#### 7. 邀请码表 (invitation_codes)

存储所有固定邀请码和临时邀请码。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 邀请码唯一标识 |
| code | VARCHAR(30) | ✅ | - | 邀请码 | 固定邀请码格式见上方"邀请码生成规则" |
| type | VARCHAR(50) | ✅ | - | 邀请码类型 | 见邀请码规则表 |
| target_role | VARCHAR(50) | ✅ | - | 目标角色 | 被邀请人将获得的角色 |
| generator_id | UUID | ✅ | - | 生成者ID | 生成此邀请码的用户ID |
| generator_type | VARCHAR(50) | ✅ | - | 生成者类型 | super_admin/merchant_admin/sp_admin/sp_staff/system |
| organization_id | UUID | ❌ | - | 关联组织ID | 商家ID或服务商ID |
| organization_type | VARCHAR(50) | ❌ | - | 组织类型 | merchant/service_provider |
| is_active | BOOLEAN | ✅ | true | 是否启用 | true: 有效, false: 已禁用 |
| is_one_time | BOOLEAN | ✅ | false | 是否一次性 | true: 一次性码（如权限转移） |
| max_uses | INT | ❌ | - | 最大使用次数 | null表示无限制 |
| use_count | INT | ✅ | 0 | 已使用次数 | 用于统计，不限制使用 |
| created_at | TIMESTAMP | ✅ | NOW() | 创建时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (code)
- INDEX (type)
- INDEX (generator_id)
- INDEX (organization_id)

**外键**：
- FOREIGN KEY (generator_id) REFERENCES users(id)
- FOREIGN KEY (organization_id) REFERENCES merchants(id) OR service_providers(id)

**数据完整性约束**：
```sql
-- 确保generator_type和generator_id的一致性
ALTER TABLE invitation_codes
ADD CONSTRAINT check_generator_consistency
CHECK (
  (generator_type = 'system' AND generator_id IS NULL) OR
  (generator_type IN ('super_admin', 'merchant_admin', 'sp_admin', 'sp_staff') AND generator_id IS NOT NULL)
);

-- 确保organization_type和organization_id的一致性
ALTER TABLE invitation_codes
ADD CONSTRAINT check_organization_consistency
CHECK (
  (organization_type IS NULL AND organization_id IS NULL) OR
  (organization_type = 'merchant' AND organization_id IN (SELECT id FROM merchants)) OR
  (organization_type = 'service_provider' AND organization_id IN (SELECT id FROM service_providers))
);
```

**max_uses 使用规则**：
- **固定邀请码**（商家员工、服务商员工、达人邀请）：max_uses = NULL（无限制，永久使用）
- **任务邀请码**：max_uses = 任务名额数（quota，如：50人任务 = max_uses=50）
- **管理员权限转移码**：max_uses = 1 + is_one_time = true（一次性，使用后自动失效）

**应用层验证逻辑**：
```sql
-- 使用邀请码时检查
SELECT * FROM invitation_codes
WHERE code = ?
  AND is_active = true
  AND (max_uses IS NULL OR use_count < max_uses);

-- 如果max_uses不为NULL且use_count >= max_uses，返回"邀请码已达使用上限"
```

**邀请码失效条件**：
- ✅ **一次性码使用后**：is_one_time=true 的邀请码，使用后自动设置 is_active=false
- ✅ **管理员手动禁用**：管理员可以手动设置 is_active=false
- ✅ **员工离职**：generator_type='sp_staff' 的员工离职时，自动设置其邀请码 is_active=false
- ✅ **商家/服务商禁用**：关联的商家/服务商被禁用时，其邀请码自动设置 is_active=false
- ✅ **达到使用上限**：max_uses 不为 NULL 且 use_count >= max_uses 时，自动设置 is_active=false

---

#### 8. 邀请记录表 (invitation_records)

记录所有邀请码的使用情况，用于返佣计算。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 记录唯一标识 |
| invitation_code_id | UUID | ✅ | - | 邀请码ID | 关联 invitation_codes.id |
| inviter_id | UUID | ✅ | - | 邀请人ID | 邀请人的用户ID |
| invitee_id | UUID | ❌ | - | 被邀请人ID | 被邀请人的用户ID（注册后填充） |
| invitee_phone | VARCHAR(20) | ❌ | - | 被邀请人手机号 | 被邀请人的手机号 |
| invitee_role | VARCHAR(50) | ❌ | - | 被邀请人角色 | 被邀请人获得的角色 |
| organization_id | UUID | ❌ | - | 关联组织ID | 商家/服务商ID |
| task_id | UUID | ❌ | - | 任务ID | 如果是任务邀请，记录任务ID |
| status | VARCHAR(20) | ✅ | 'pending' | 状态 | pending/registered/completed |
| used_at | TIMESTAMP | ❌ | - | 使用时间 | - |
| created_at | TIMESTAMP | ✅ | NOW() | 创建时间 | - |

**索引**：
- PRIMARY KEY (id)
- INDEX (invitation_code_id)
- INDEX (inviter_id)
- INDEX (invitee_id)
- INDEX (task_id)

**外键**：
- FOREIGN KEY (invitation_code_id) REFERENCES invitation_codes(id)
- FOREIGN KEY (inviter_id) REFERENCES users(id)
- FOREIGN KEY (invitee_id) REFERENCES users(id)
- FOREIGN KEY (task_id) REFERENCES tasks(id)

---

**说明**：
- ✅ 任务邀请码使用统一的 `invitation_codes` 表
- ✅ 通过 `target_type='task'` 和 `target_id` 关联任务
- ✅ 通过 `max_uses` 控制使用次数（任务名额数）
- ✅ 通过 `expires_at` 设置过期时间（任务截止时间）
- ✅ 代码格式：`TASK-{任务ID后6位}`（固定，基于任务ID）

**示例**：
- 任务ID: `task-1234567890abcdef`
- 任务邀请码: `TASK-789012`（存储在invitation_codes表）
- max_uses: 50（活动名额数）
- expires_at: 2026-03-01 23:59:59（提交截止时间）

---

#### 9. 商家员工表 (merchant_staff)

存储商家员工信息及其权限。

**⚠️ 组织结构**：
- 一个商家只有**1个商家管理员**（MERCHANT_ADMIN）
- 管理员创建者自动成为管理员
- 商家可以有**多个商家员工**（MERCHANT_STAFF）
- 管理员可以给员工勾选所有权限（员工最高权限可以等于管理员）

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 商家员工唯一标识 |
| user_id | UUID | ✅ | - | 关联用户ID | 关联 users.id |
| merchant_id | UUID | ✅ | - | 商家ID | 关联 merchants.id |
| title | VARCHAR(50) | ❌ | - | 职能名称 | 例如："运营"、"财务"、"审核专员"，纯展示用途 |
| status | VARCHAR(20) | ✅ | 'active' | 状态 | active: 正常, inactive: 离职 |
| created_at | TIMESTAMP | ✅ | NOW() | 添加时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (user_id)
- INDEX (merchant_id)
- INDEX (status)

**外键**：
- FOREIGN KEY (user_id) REFERENCES users(id)
- FOREIGN KEY (merchant_id) REFERENCES merchants(id)

---

#### 10. 商家员工权限表 (merchant_staff_permissions)

存储商家员工的权限。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 权限记录唯一标识 |
| staff_id | UUID | ✅ | - | 员工ID | 关联 merchant_staff.id |
| permission_code | VARCHAR(50) | ✅ | - | 权限代码 | 例如: task.view, task.create, finance.view |
| granted_at | TIMESTAMP | ✅ | NOW() | 授权时间 | - |
| granted_by | UUID | ✅ | - | 授权人ID | 授予权限的管理员ID |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (staff_id, permission_code)
- INDEX (staff_id)
- INDEX (permission_code)

**外键**：
- FOREIGN KEY (staff_id) REFERENCES merchant_staff(id)
- FOREIGN KEY (granted_by) REFERENCES merchant_staff(id)

---

#### 11. 服务商员工表 (service_provider_staff)

存储服务商员工信息及其权限。

**⚠️ 组织结构**：
- 一个服务商只有**1个服务商管理员**（SERVICE_PROVIDER_ADMIN）
- 管理员创建者自动成为管理员
- 服务商可以有**多个服务商员工**（SERVICE_PROVIDER_STAFF）
- 管理员可以给员工勾选所有权限（员工最高权限可以等于管理员）

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 服务商员工唯一标识 |
| user_id | UUID | ✅ | - | 关联用户ID | 关联 users.id |
| provider_id | UUID | ✅ | - | 服务商ID | 关联 service_providers.id |
| title | VARCHAR(50) | ❌ | - | 职能名称 | 例如："审核专员"、"商务拓展"、"运营"，纯展示用途 |
| status | VARCHAR(20) | ✅ | 'active' | 状态 | active: 正常, inactive: 离职 |
| created_at | TIMESTAMP | ✅ | NOW() | 添加时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (user_id)
- INDEX (provider_id)
- INDEX (status)

**外键**：
- FOREIGN KEY (user_id) REFERENCES users(id)
- FOREIGN KEY (provider_id) REFERENCES service_providers(id)

---

#### 12. 服务商员工权限表 (provider_staff_permissions)

存储服务商员工的权限。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 权限记录唯一标识 |
| staff_id | UUID | ✅ | - | 员工ID | 关联 service_provider_staff.id |
| permission_code | VARCHAR(50) | ✅ | - | 权限代码 | 例如: task.audit, creator.invite, merchant.bind |
| granted_at | TIMESTAMP | ✅ | NOW() | 授权时间 | - |
| granted_by | UUID | ✅ | - | 授权人ID | 授予权限的管理员ID |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (staff_id, permission_code)
- INDEX (staff_id)
- INDEX (permission_code)

**外键**：
- FOREIGN KEY (staff_id) REFERENCES service_provider_staff(id)
- FOREIGN KEY (granted_by) REFERENCES service_provider_staff(id)

---

#### 13. 积分账户表 (credit_accounts)

存储各角色和组织实体的积分账户。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 账户唯一标识 |
| owner_id | UUID | ✅ | - | 账户持有者ID | 根据 owner_type 关联不同表 |
| owner_type | VARCHAR(20) | ✅ | - | 账户类型 | ORG_MERCHANT: 商家组织, ORG_PROVIDER: 服务商组织, USER_PERSONAL: 用户个人 |
| balance | INT | ✅ | 0 | 可用余额 | 单位：积分（1元=1积分） |
| frozen_balance | INT | ✅ | 0 | 冻结余额 | ⚠️ **双重语义（按账户类型区分）**：
**商家/服务商（ORG_MERCHANT/ORG_PROVIDER）**：
- 预扣未分配的活动费用（如：发布10人×100元任务，冻结1000元）
- 达人接任务时减少（如：-100元）
- 任务关闭时退回balance（如：未完成2人，退回200元）

**达人/员工个人（USER_PERSONAL）**：
- **类型A**：已提交但待审核的任务收入（如：提交任务后+80元，审核通过后转入balance）
- **类型B**：提现申请通过后待打款的金额（如：提现100元通过后，从balance转入100元，打款后清除）

⚠️ **重要**：商家/服务商的frozen_balance存储"预扣费用"，达人/员工的frozen_balance存储"待确认收入+待打款金额"，语义不同！
| created_at | TIMESTAMP | ✅ | NOW() | 创建时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (owner_id, owner_type)
- INDEX (owner_type)

**数据完整性约束**：
```sql
-- ⚠️ 防止余额为负数（财务安全保护）
ALTER TABLE credit_accounts
ADD CONSTRAINT check_balance_non_negative
CHECK (balance >= 0);

ALTER TABLE credit_accounts
ADD CONSTRAINT check_frozen_balance_non_negative
CHECK (frozen_balance >= 0);

-- 应用层检查：所有UPDATE语句前必须先验证余额是否足够
-- 示例：
-- IF 账户.balance < 扣除金额 THEN RETURN "余额不足";

**⚠️ 账户归属说明**：
- **ORG_MERCHANT**: 商家组织账户，属于商家实体（关联 `merchants.id`），不属于管理员个人
- **ORG_PROVIDER**: 服务商组织账户，属于服务商实体（关联 `service_providers.id`），不属于管理员个人
- **USER_PERSONAL**: 用户个人账户，属于用户个人（关联 `users.id`），所有个人收入统一到此账户

**多角色用户的账户管理**：
- **场景1**：用户既是"商家员工"又是"达人"
  - 商家员工收入：无（商家员工没有个人收入，所有费用归商家组织）
  - 达人收入：进入USER_PERSONAL账户（owner_id=用户ID）
- **场景2**：用户创建多个商家
  - 商家A的账户：ORG_MERCHANT（owner_id=商家A的ID）
  - 商家B的账户：ORG_MERCHANT（owner_id=商家B的ID）
  - 个人账户：USER_PERSONAL（owner_id=用户ID）
- **场景3**：用户是"服务商员工"+"达人"
  - 服务商员工返佣：进入USER_PERSONAL账户
  - 达人任务收入：进入USER_PERSONAL账户
  - ✅ **统一管理**：所有个人收入进入同一个USER_PERSONAL账户

**⚠️ 重要设计决策**：
1. 商家和服务商的积分属于**组织实体**，不属于管理员
2. 管理员离职/更换时，组织账户不受影响
3. 管理员只是账户的操作者，不是所有者
4. **用户的个人收入统一管理**：
   - 达人接任务收入 → 进入 USER_PERSONAL 账户
   - 服务商员工邀请返佣 → 进入 USER_PERSONAL 账户
   - 用户切换角色后，个人余额不变（同一个账户）

**外键关联逻辑**：
```sql
-- 应用层根据 owner_type 动态关联
-- ORG_MERCHANT → merchants.id
-- ORG_PROVIDER → service_providers.id
-- USER_PERSONAL → users.id
```

**frozen_balance 使用规则**：

**商家账户（ORG_MERCHANT）- 预扣活动费用**：
- 增加：发布任务时 `+总费用`（如：10人×100元=1000元）
- 减少：达人接任务时 `-任务金额`（如：-100元）；达人提交时 `-达人预计收入`（如：-80元）；审核通过时 `-员工返佣-服务商收入`（如：-20元）
- 退款：任务手动关闭或名额已满时，未完成名额费用 `frozen_balance → balance`

**服务商账户（ORG_PROVIDER）- 预扣活动费用**：
- 增加：任务审核通过时 `+provider_amount`（收入直接进balance，通常不使用frozen_balance）
- 减少：提现申请通过时，`balance → frozen_balance`（待打款金额）

**达人/员工个人账户（USER_PERSONAL）- 待确认收入 + 待打款金额**：
- **场景A-任务待审核**：
  - 增加：提交任务时 `+预计收入`（等待审核，如：+80元）
  - 转账：审核通过时，`frozen_balance → balance`（+80元转入可用余额）
  - 减少：审核拒绝时，`frozen_balance` 减少（释放冻结，如：-80元）
- **场景B-提现待打款**：
  - 增加：提现申请通过时，`balance → frozen_balance`（如：提现100元通过，frozen_balance+100元）
  - 减少：财务打款后，`frozen_balance` 减少（如：打款100元后，frozen_balance-100元）
  - 退款：提现被拒绝时，`frozen_balance → balance`（转回可用余额）

⚠️ **财务对账时注意**：商家frozen_balance是"负债"（待支付费用），达人frozen_balance是"资产"（待确认收入），方向相反！

---

#### 14. 积分交易类型表 (transaction_types)

**⚠️ 新增表**：将交易类型从硬编码枚举改为数据库表存储，支持版本管理和动态扩展。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| code | VARCHAR(50) | ✅ | - | 类型代码 | 唯一标识，如 RECHARGE, TASK_INCOME |
| name | VARCHAR(100) | ✅ | - | 类型名称 | 中文描述，如"商家充值" |
| description | TEXT | ❌ | - | 类型说明 | 详细说明 |
| account_types | VARCHAR(20)[] | ✅ | - | 适用账户类型 | 如：{ORG_MERCHANT, USER_PERSONAL} |
| amount_direction | VARCHAR(10) | ✅ | - | 金额方向 | positive（正数）或negative（负数） |
| is_active | BOOLEAN | ✅ | true | 是否启用 | true：有效，false：已废弃 |
| version | INT | ✅ | 1 | 版本号 | 用于版本管理 |
| created_at | TIMESTAMP | ✅ | NOW() | 创建时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |
| deprecated_at | TIMESTAMP | ❌ | - | 废弃时间 | 标记为废弃的时间 |

**索引**：
- PRIMARY KEY (code)
- INDEX (is_active)
- INDEX (version)

**初始数据**：
```sql
INSERT INTO transaction_types (code, name, description, account_types, amount_direction) VALUES
('RECHARGE', '商家充值', '商家充值积分', ARRAY['ORG_MERCHANT'], 'positive'),
('TASK_INCOME', '任务收入', '达人任务审核通过，frozen_balance转入balance', ARRAY['USER_PERSONAL'], 'positive'),
('TASK_SUBMIT', '任务提交', '达人提交任务，预计收入冻结到frozen_balance', ARRAY['USER_PERSONAL'], 'positive'),
('STAFF_REFERRAL', '员工返佣', '员工邀请的达人完成任务时的返佣', ARRAY['USER_PERSONAL'], 'positive'),
('PROVIDER_INCOME', '服务商收入', '服务商完成任务审核通过获得的收入', ARRAY['ORG_PROVIDER'], 'positive'),
('TASK_PUBLISH', '发布任务', '商家发布任务扣除的积分（balance→frozen_balance）', ARRAY['ORG_MERCHANT'], 'negative'),
('TASK_ACCEPT', '接任务扣除', '达人接任务时从商家frozen_balance扣除', ARRAY['ORG_MERCHANT'], 'negative'),
('TASK_REJECT', '审核拒绝', '达人任务被拒绝，释放frozen_balance', ARRAY['USER_PERSONAL'], 'negative'),
('TASK_ESCALATE', '超时拒绝', '提交截止时间到，自动拒绝', ARRAY['ORG_MERCHANT'], 'positive'),
('TASK_REFUND', '任务退款', '任务关闭时，未完成名额费用退回balance', ARRAY['ORG_MERCHANT'], 'positive'),
('WITHDRAW', '提现', '用户申请提现（扣除balance）', ARRAY['USER_PERSONAL', 'ORG_MERCHANT', 'ORG_PROVIDER'], 'negative'),
('WITHDRAW_FREEZE', '提现冻结', '提现申请通过，从balance转入frozen_balance', ARRAY['USER_PERSONAL', 'ORG_MERCHANT', 'ORG_PROVIDER'], 'negative'),
('WITHDRAW_REFUND', '提现拒绝', '提现被拒绝，frozen_balance转回balance', ARRAY['USER_PERSONAL', 'ORG_MERCHANT', 'ORG_PROVIDER'], 'positive'),
('BONUS_GIFT', '系统赠送', '平台赠送的积分（如新用户奖励）', ARRAY['USER_PERSONAL'], 'positive');
```

**⚠️ 版本管理规则**：
- 新增交易类型：INSERT新记录，version=1
- 修改交易类型：更新version+=1，保留旧记录（is_active=false）
- 废弃交易类型：设置is_active=false，记录deprecated_at时间
- 外键引用：credit_transactions.type REFERENCES transaction_types(code)

---

#### 15. 积分流水表 (credit_transactions)

**积分交易类型枚举**（⚠️ 已废弃，改用transaction_types表）：

| 类型代码 | 类型名称 | 说明 | 适用账户类型 | 金额 |
|---------|---------|------|-------------|------|
| `RECHARGE` | 商家充值 | 商家充值积分 | ORG_MERCHANT | 正数 |
| `TASK_INCOME` | 任务收入 | 达人任务审核通过，frozen_balance 转入 balance | USER_PERSONAL | 正数 |
| `TASK_SUBMIT` | 任务提交 | 达人提交任务，预计收入冻结到 frozen_balance | USER_PERSONAL | 正数 |
| `STAFF_REFERRAL` | 员工返佣 | 员工邀请的达人完成任务时的返佣 | USER_PERSONAL | 正数 |
| `PROVIDER_INCOME` | 服务商收入 | 服务商完成任务审核通过获得的收入 | ORG_PROVIDER | 正数 |
| `TASK_PUBLISH` | 发布任务 | 商家发布任务扣除的积分（balance → frozen_balance） | ORG_MERCHANT | 负数 |
| `TASK_ACCEPT` | 接任务扣除 | 达人接任务时从商家 frozen_balance 扣除（记账，减少冻结） | ORG_MERCHANT | 负数 |
| `TASK_REJECT` | 审核拒绝 | 达人任务被拒绝，释放 frozen_balance | USER_PERSONAL | 负数 |
| `TASK_ESCALATE` | 超时拒绝 | 提交截止时间到，自动拒绝，费用退回商家 frozen_balance | ORG_MERCHANT | 正数 |
| `TASK_REFUND` | 任务退款 | 任务关闭时，未完成名额费用从 frozen_balance 退回 balance | ORG_MERCHANT | 正数 |
| `WITHDRAW` | 提现 | 用户申请提现（扣除 balance） | 所有账户类型 | 负数 |
| `WITHDRAW_FREEZE` | 提现冻结 | 提现申请通过，从 balance 转入 frozen_balance（待打款） | 所有账户类型 | 负数 |
| `WITHDRAW_REFUND` | 提现拒绝 | 提现被拒绝，frozen_balance 转回 balance | 所有账户类型 | 正数 |
| `BONUS_GIFT` | 系统赠送 | 平台赠送的积分（如新用户奖励） | USER_PERSONAL | 正数 |

**⚠️ 重要设计原则**：
- 所有积分变动都必须记录交易流水
- 每笔交易必须记录：交易前余额、交易后余额、交易类型、关联对象
- 支持财务对账：来源总额 - 去向总额 = 当前总余额

---

记录所有积分变动。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 交易唯一标识 |
| account_id | UUID | ✅ | - | 账户ID | 关联 credit_accounts.id |
| type | VARCHAR(50) | ✅ | - | 交易类型 | 见上方"积分交易类型枚举"表格 |
| amount | INT | ✅ | - | 金额 | 正数表示增加，负数表示减少 |
| balance_before | INT | ✅ | - | 交易前余额 | - |
| balance_after | INT | ✅ | - | 交易后余额 | - |
| transaction_group_id | UUID | ❌ | - | 交易组ID | ⚠️ 关联的多个交易共享同一个group_id（如任务审核通过的4方交易） |
| group_sequence | INT | ❌ | - | 组内序号 | 同一交易组内的执行顺序（1,2,3,4） |
| related_campaign_id | UUID | ❌ | - | 关联活动ID | 关联 campaigns.id（如果与任务相关） |
| related_task_id | UUID | ❌ | - | 关联任务名额ID | 关联 tasks.id（如果与任务名额相关） |
| description | VARCHAR(200) | ❌ | - | 说明 | 交易说明 |
| created_at | TIMESTAMP | ✅ | NOW() | 交易时间 | - |

**索引**：
- PRIMARY KEY (id)
- INDEX (account_id)
- INDEX (type)
- INDEX (transaction_group_id) -- ⚠️ 用于查询同一交易组的所有流水
- INDEX (created_at)

**外键**：
- FOREIGN KEY (account_id) REFERENCES credit_accounts(id)
- FOREIGN KEY (related_task_id) REFERENCES campaigns(id)
- FOREIGN KEY (related_task_slot_id) REFERENCES tasks(id)

---

#### 15. 提现表 (withdrawals)

存储提现申请。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 提现唯一标识 |
| account_id | UUID | ✅ | - | 账户ID | 关联 credit_accounts.id |
| amount | INT | ✅ | - | 提现金额 | 单位：积分（1元=1积分） |
| fee | INT | ✅ | 0 | 手续费 | 单位：积分 |
| actual_amount | INT | ✅ | - | 实际到账金额 | amount - fee |
| method | VARCHAR(20) | ✅ | - | 提现方式 | ALIPAY: 支付宝, WECHAT: 微信, BANK: 银行转账 |
| account_info | JSONB | ✅ | - | 收款信息 | ⚠️ **AES-256加密存储**，{"name": "张三", "account": "xxx", "bank": "xxx"} |
| account_info_hash | VARCHAR(64) | ❌ | - | 收款信息哈希 | ⚠️ SHA-256哈希，用于查询和去重，不存储敏感信息 |
| status | VARCHAR(20) | ✅ | 'pending' | 状态 | pending: 待审核, approved: 已通过, rejected: 已拒绝, completed: 已完成 |
| audit_note | VARCHAR(200) | ❌ | - | 审核备注 | - |
| audited_by | UUID | ❌ | - | 审核人ID | 审核此提现的超级管理员ID（所有提现都需要审核，包括服务商、达人、员工） |
| audited_at | TIMESTAMP | ❌ | - | 审核时间 | - |
| completed_at | TIMESTAMP | ❌ | - | 完成时间 | - |
| created_at | TIMESTAMP | ✅ | NOW() | 申请时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |

**索引**：
- PRIMARY KEY (id)
- INDEX (account_id)
- INDEX (status)
- INDEX (created_at)

**外键**：
- FOREIGN KEY (account_id) REFERENCES credit_accounts(id)
- FOREIGN KEY (audited_by) REFERENCES users(id)

**⚠️ 数据安全说明**：
- account_info使用AES-256加密存储（密钥存储在环境变量中）
- account_info_hash使用SHA-256单向哈希，用于快速查询和去重
- 前端传输时必须使用HTTPS协议
- 审核人员查看明文需要二次权限验证

---

#### 16. 商家-服务商绑定表 (merchant_provider_bindings)

记录商家和服务商的绑定关系。

| 字段名 | 类型 | 是否必填 | 默认值 | 前端显示名 | 说明 |
|-------|------|---------|--------|-----------|------|
| id | UUID | ✅ | - | （内部使用） | 绑定唯一标识 |
| merchant_id | UUID | ✅ | - | 商家ID | 关联 merchants.id |
| provider_id | UUID | ✅ | - | 服务商ID | 关联 service_providers.id |
| status | VARCHAR(20) | ✅ | 'active' | 状态 | active: 正常, inactive: 已解绑 |
| created_at | TIMESTAMP | ✅ | NOW() | 绑定时间 | - |
| updated_at | TIMESTAMP | ✅ | NOW() | 更新时间 | - |

**索引**：
- PRIMARY KEY (id)
- UNIQUE INDEX (merchant_id, provider_id)
- INDEX (merchant_id)
- INDEX (provider_id)
- INDEX (status)

**外键**：
- FOREIGN KEY (merchant_id) REFERENCES merchants(id)
- FOREIGN KEY (provider_id) REFERENCES service_providers(id)

**⚠️ 双向关联优化**：
- 绑定表保留（用于多对多关系）
- merchants表添加冗余字段`default_provider_id`（默认服务商）
- service_providers表添加冗余字段`merchant_ids JSONB`（关联的商家ID列表）
- 查询优化：常用查询使用冗余字段，复杂查询使用绑定表

---

### 表关系图

```
users (1) ----< (1) merchants (admin_id)
users (1) ----< (1) service_providers (admin_id)
users (1) ----< (1) creators
users (1) ----< (1) merchant_staff
users (1) ----< (1) service_provider_staff

merchants (1) ----< (N) tasks
merchants (1) ----< (N) merchant_provider_bindings
merchants (1) ----< (1) credit_accounts (商家组织账户)

service_providers (1) ----< (N) merchant_provider_bindings
service_providers (1) ----< (N) tasks (optional)
service_providers (1) ----< (1) credit_accounts (服务商组织账户)

users (1) ----< (1) credit_accounts (用户个人账户，所有个人收入统一)

campaigns (1) ----< (N) tasks (一个营销活动有多个任务名额)
campaigns (1) ----< (N) invite_codes (活动专属邀请码)

tasks (N) ----< (1) campaigns (每个任务名额属于一个营销活动)
tasks (N) ----< (1) creators (每个任务名额属于一个达人，OPEN状态时为NULL)

creators (1) ----< (N) tasks (一个达人可以接多个任务名额)
creators (N) ----< (1) users (inviter - 邀请人)

invite_codes (1) ----< (N) invite_code_usages

merchant_staff (1) ----< (N) merchant_staff_permissions
service_provider_staff (1) ----< (N) provider_staff_permissions

credit_accounts (1) ----< (N) credit_transactions
credit_accounts (1) ----< (N) withdrawals
```

---

## 🔐 财务安全与一致性保证

### 1. 关键操作的事务性要求

**⚠️ 原则**：所有涉及资金变动的操作必须使用数据库事务，确保原子性。

#### 必须使用事务的操作

**发布任务（预扣款）**：
```sql
BEGIN TRANSACTION;

1. 检查商家 balance 是否足够
   IF 商家 balance < 总费用 THEN
     ROLLBACK;
     RETURN "余额不足";
   END IF;

2. 扣除商家 balance
   UPDATE credit_accounts
   SET balance = balance - 总费用,
       frozen_balance = frozen_balance + 总费用
   WHERE owner_id = 商家ID AND owner_type = 'ORG_MERCHANT';

3. 创建任务记录
   INSERT INTO tasks (id, merchant_id, total_amount, quota, ...)
   VALUES (...);

4. 记录交易流水
   INSERT INTO credit_transactions (account_id, type, amount, ...)
   VALUES (商家账户ID, 'TASK_PUBLISH', -总费用, ...);

COMMIT;
```

**任何步骤失败 → 全部回滚，确保资金安全**

---

**达人接任务（分配名额费用）**：
```sql
BEGIN TRANSACTION;

1. 检查营销活动状态 = OPEN
   IF 营销活动状态 ≠ 'OPEN' THEN
     ROLLBACK;
     RETURN "任务未开放";
   END IF;

2. 检查是否已接任务
   IF 已存在活动名额 WHERE creator_id = 达人ID AND task_id = 活动ID THEN
     ROLLBACK;
     RETURN "已接任务";
   END IF;

3. 检查商家 frozen_balance 是否足够
   IF 商家 frozen_balance < 任务金额 THEN
     ROLLBACK;
     RETURN "商家冻结余额不足";
   END IF;

4. 扣除商家 frozen_balance（记账，标记为"已接任务，待分配"）
   UPDATE credit_accounts
   SET frozen_balance = frozen_balance - 任务金额
   WHERE owner_id = 商家ID AND owner_type = 'ORG_MERCHANT';

5. 更新任务名额状态
   UPDATE tasks
   SET status = 'ASSIGNED',
       creator_id = 达人ID,
       assigned_at = NOW(),
       invite_code_id = 邀请码ID,
       invited_by_staff_id = 员工ID
   WHERE id = 任务名额ID
     AND campaign_id = 活动ID
     AND status = 'OPEN';

6. 记录交易流水
   INSERT INTO credit_transactions (account_id, type, amount, ...)
   VALUES (商家账户ID, 'TASK_ACCEPT', -任务金额, ...);

COMMIT;
```

---

**达人提交任务（冻结收入到 frozen_balance）**：
```sql
BEGIN TRANSACTION;

1. 检查活动名额状态 = ASSIGNED
   IF 活动名额状态 ≠ 'ASSIGNED' THEN
     ROLLBACK;
     RETURN "活动名额状态错误";
   END IF;

2. 检查提交截止时间
   IF NOW() > 提交截止时间 THEN
     ROLLBACK;
     RETURN "已超过提交截止时间";
   END IF;

3. 冻结达人预计收入（从商家 frozen_balance 转移）
   UPDATE credit_accounts
   SET frozen_balance = frozen_balance - 达人预计收入
   WHERE owner_id = 商家ID AND owner_type = 'ORG_MERCHANT';

4. 增加达人 frozen_balance（预计收入）
   UPDATE credit_accounts
   SET frozen_balance = frozen_balance + 达人预计收入
   WHERE owner_id = 达人ID AND owner_type = 'USER_PERSONAL';

5. 更新任务名额状态
   UPDATE tasks
   SET status = 'SUBMITTED',
       platform_url = 平台链接,
       screenshots = 截图数组,
       submitted_at = NOW(),
       note = 备注,
       updated_at = NOW()
   WHERE id = 活动名额ID;

6. 记录双方交易流水
   INSERT INTO credit_transactions (account_id, type, amount, ...)
   VALUES
     (商家账户ID, 'TASK_SUBMIT', -达人预计收入, ...),
     (达人账户ID, 'TASK_SUBMIT', 达人预计收入, ...);

COMMIT;
```

**任何步骤失败 → 全部回滚，确保资金安全**

---

**任务审核通过（积分分配）**：
```sql
BEGIN TRANSACTION;

1. 检查活动名额状态 = SUBMITTED
   IF 活动名额状态 ≠ 'SUBMITTED' THEN
     ROLLBACK;
     RETURN "活动名额状态错误";
   END IF;

2. 检查是否已审核（防止重复结算）
   IF 活动名额.audited_by IS NOT NULL THEN
     ROLLBACK;
     RETURN "活动名额已审核";
   END IF;

3. 检查商家 frozen_balance 是否足够（员工返佣 + 服务商收入）
   IF 商家 frozen_balance < 员工返佣 + 服务商收入 THEN
     ROLLBACK;
     RETURN "商家冻结余额不足";
   END IF;

4. 检查达人 frozen_balance 是否足够
   IF 达人 frozen_balance < 达人金额 THEN
     ROLLBACK;
     RETURN "达人冻结余额不足";
   END IF;

5. 更新任务名额状态（⚠️ 使用乐观锁防止并发审核）
   UPDATE tasks
   SET status = 'APPROVED',
       audited_by = 审核人ID,
       audited_at = NOW(),
       completed_at = NOW(),
       version = version + 1
   WHERE id = 活动名额ID
     AND status = 'SUBMITTED'
     AND audited_by IS NULL
     AND version = ?;  -- 前端传入的当前版本号

   -- 检查是否更新成功（affected_rows = 1）
   IF affected_rows = 0 THEN
     ROLLBACK;
     RETURN "活动名额已被审核或版本不匹配，请刷新后重试";
   END IF;

6. 分配达人积分（frozen_balance → balance）
   UPDATE credit_accounts
   SET frozen_balance = frozen_balance - 达人金额,
       balance = balance + 达人金额
   WHERE owner_id = 达人ID AND owner_type = 'USER_PERSONAL';

7. 分配员工返佣（从商家 frozen_balance 转入员工 frozen_balance，再转入 balance）
   UPDATE credit_accounts
   SET frozen_balance = frozen_balance + 员工返佣,
       balance = balance + 员工返佣
   WHERE owner_id = 员工ID AND owner_type = 'USER_PERSONAL';

8. 分配服务商收入（从商家 frozen_balance 转入服务商 balance）
   UPDATE credit_accounts
   SET balance = balance + 服务商收入
   WHERE owner_id = 服务商ID AND owner_type = 'ORG_PROVIDER';

9. 扣除商家 frozen_balance（员工返佣 + 服务商收入）
   UPDATE credit_accounts
   SET frozen_balance = frozen_balance - (员工返佣 + 服务商收入)
   WHERE owner_id = 商家ID AND owner_type = 'ORG_MERCHANT';

10. 记录多方交易流水
    INSERT INTO credit_transactions (account_id, type, amount, ...)
    VALUES
      (达人账户ID, 'TASK_INCOME', 达人金额, ...),
      (员工账户ID, 'STAFF_REFERRAL', 员工返佣, ...),
      (服务商账户ID, 'PROVIDER_INCOME', 服务商收入, ...),
      (商家账户ID, 'TASK_ACCEPT', -(员工返佣 + 服务商收入), ...);

COMMIT;
```

**关键安全检查**：
- ✅ 状态检查：确保活动名额是 SUBMITTED 状态
- ✅ 防重复审核：检查 audited_by 是否为空
- ✅ 余额检查：确保商家和达人的 frozen_balance 足够
- ✅ 任何检查失败 → 回滚

---

**提现申请（扣除 balance，转入 frozen_balance）**：
```sql
BEGIN TRANSACTION;

1. 检查用户 balance 是否足够
   IF 用户 balance < 提现金额 + 手续费 THEN
     ROLLBACK;
     RETURN "余额不足";
   END IF;

2. 扣除用户 balance，转入 frozen_balance
   UPDATE credit_accounts
   SET balance = balance - (提现金额 + 手续费),
       frozen_balance = frozen_balance + (提现金额 + 手续费)
   WHERE id = 账户ID;

3. 创建提现记录
   INSERT INTO withdrawals (id, account_id, amount, fee, status, ...)
   VALUES (提现ID, 账户ID, 提现金额, 手续费, 'pending', ...);

4. 记录交易流水
   INSERT INTO credit_transactions (account_id, type, amount, ...)
   VALUES (账户ID, 'WITHDRAW', -(提现金额 + 手续费), ...);

COMMIT;
```

---

**提现审核通过（超级管理员）**：
```sql
BEGIN TRANSACTION;

1. 检查提现记录状态 = 'pending'
   IF 状态 ≠ 'pending' THEN
     ROLLBACK;
     RETURN "提现记录状态错误";
   END IF;

2. 更新提现记录状态
   UPDATE withdrawals
   SET status = 'approved',
       audited_by = 超级管理员ID,
       audited_at = NOW()
   WHERE id = 提现ID;

3. 记录交易流水（从 balance 转入 frozen_balance）
   INSERT INTO credit_transactions (account_id, type, amount, ...)
   VALUES (账户ID, 'WITHDRAW_FREEZE', -(提现金额 + 手续费), ...);

COMMIT;
```

---

**提现审核拒绝（超级管理员）**：
```sql
BEGIN TRANSACTION;

1. 检查提现记录状态 = 'pending'
   IF 状态 ≠ 'pending' THEN
     ROLLBACK;
     RETURN "提现记录状态错误";
   END IF;

2. 更新提现记录状态
   UPDATE withdrawals
   SET status = 'rejected',
       audited_by = 超级管理员ID,
       audited_at = NOW(),
       reject_reason = 拒绝原因
   WHERE id = 提现ID;

3. 释放用户 frozen_balance，转回 balance
   UPDATE credit_accounts
   SET balance = balance + (提现金额 + 手续费),
       frozen_balance = frozen_balance - (提现金额 + 手续费)
   WHERE id = 账户ID;

4. 记录交易流水
   INSERT INTO credit_transactions (account_id, type, amount, ...)
   VALUES (账户ID, 'WITHDRAW_REFUND', 提现金额 + 手续费, ...);

COMMIT;
```

---

**服务商拒绝任务（TASK_REJECT）**：
```sql
BEGIN TRANSACTION;

1. 检查活动名额状态 = SUBMITTED
   IF 活动名额状态 ≠ 'SUBMITTED' THEN
     ROLLBACK;
     RETURN "活动名额状态错误";
   END IF;

2. 释放达人 frozen_balance（预计收入）
   UPDATE credit_accounts
   SET frozen_balance = frozen_balance - 达人预计收入
   WHERE owner_id = 达人ID AND owner_type = 'USER_PERSONAL';

3. ⚠️ 商家 frozen_balance 保持冻结（不返还，等待达人重新提交）
   -- 无需SQL操作，商家frozen_balance中的费用继续冻结

4. 更新任务名额状态（⚠️ 使用乐观锁防止并发审核）
   UPDATE tasks
   SET status = 'REJECTED',
       audited_by = 审核人ID,
       audited_at = NOW(),
       reject_reason = 拒绝原因,
       version = version + 1
   WHERE id = 活动名额ID
     AND status = 'SUBMITTED'
     AND audited_by IS NULL
     AND version = ?;  -- 前端传入的当前版本号

   -- 检查是否更新成功（affected_rows = 1）
   IF affected_rows = 0 THEN
     ROLLBACK;
     RETURN "活动名额已被审核或版本不匹配，请刷新后重试";
   END IF;

5. 记录交易流水（TASK_REJECT）
   INSERT INTO credit_transactions (account_id, type, amount, ...)
   VALUES
     (达人账户ID, 'TASK_REJECT', -达人预计收入, ...);

   -- ⚠️ 商家frozen_balance保持冻结，无需记录商家交易流水

COMMIT;
```

---

**超时自动拒绝（TASK_ESCALATE）**：
```sql
-- 定时任务：每小时检查提交截止时间到期的活动名额

BEGIN TRANSACTION;

1. 查询所有 ASSIGNED 状态且超过提交截止时间的任务名额
   FOR EACH 任务名额 IN
     SELECT * FROM tasks
     WHERE status = 'ASSIGNED'
       AND campaign_id IN (SELECT id FROM campaigns WHERE submission_deadline < NOW())
   DO
     -- 获取活动名额的活动费用信息
     SELECT 任务金额, 商家ID, 达人ID INTO 费用信息
     FROM tasks WHERE id = 活动名额.task_id;

     -- 2. 返还商家 frozen_balance（全部费用）
     UPDATE credit_accounts
     SET frozen_balance = frozen_balance + 任务金额
     WHERE owner_id = 商家ID AND owner_type = 'ORG_MERCHANT';

     -- 3. 更新任务名额状态
     UPDATE tasks
     SET status = 'REJECTED',
         audit_note = '提交截止时间已到，自动拒绝',
         updated_at = NOW()
     WHERE id = 活动名额ID;

     -- 4. 记录交易流水
     INSERT INTO credit_transactions (account_id, type, amount, ...)
     VALUES (商家账户ID, 'TASK_ESCALATE', 任务金额, ...);
   END FOR;

COMMIT;
```

---

**任务关闭退款（TASK_REFUND）**：
```sql
-- 任务手动关闭或名额已满时，退还未接活动名额的费用

BEGIN TRANSACTION;

1. 获取活动信息和已接名额数
   SELECT 总名额, 任务金额, 已接名额数 INTO 活动信息
   FROM tasks
   WHERE id = 活动ID;


2. 计算未完成名额
   未完成名额 = 总名额 - 已接名额数;

3. 检查是否有未完成名额
   IF 未完成名额 <= 0 THEN
     ROLLBACK;
     RETURN "无未完成名额，无需退款";
   END IF;

4. 计算退款金额
   退款金额 = 未完成名额 × 任务金额;

5. 释放商家 frozen_balance，转回 balance
   UPDATE credit_accounts
   SET frozen_balance = frozen_balance - 退款金额,
       balance = balance + 退款金额
   WHERE owner_id = 商家ID AND owner_type = 'ORG_MERCHANT';

6. 更新营销活动状态
   UPDATE tasks
   SET status = 'CLOSED',
       closed_at = NOW()
   WHERE id = 活动ID;

7. 记录交易流水
   INSERT INTO credit_transactions (account_id, type, amount, ...)
   VALUES (商家账户ID, 'TASK_REFUND', 退款金额, ...);

COMMIT;
```

**⚠️ 重要说明**：
- 只有未接任务的名额才退款
- 已接任务但未完成的记录（ASSIGNED/SUBMITTED）不参与退款

**任何步骤失败 → 全部回滚，确保资金安全**

---

### 2. 防重复结算机制

**问题**：在高并发情况下，同一个活动名额可能被多个审核员同时审核通过，导致重复分配积分。

**解决方案**：使用数据库级别的条件更新（Compare-and-Set）

```sql
-- 方案1：使用 WHERE 条件更新
UPDATE tasks
SET status = 'APPROVED',
    audited_by = 审核人ID,
    audited_at = NOW()
WHERE id = 活动名额ID
  AND status = 'SUBMITTED'           -- 只有 SUBMITTED 状态才能更新
  AND audited_by IS NULL;          -- 未被审核过才能更新

-- 如果 affected_rows = 0，说明已被其他审核员审核，拒绝操作
```

**好处**：
- ✅ 数据库层面保证原子性
- ✅ 即使并发审核，也只会有一人成功
- ✅ 失败的审核操作可以明确告知用户"该接任务已被其他审核员审核"

---

### 3. 财务对账机制

**对账目的**：确保账实相符，及时发现财务异常。

#### 对账类型

**1. 积分平衡对账**

**⚠️ 重要说明**：
- 由于我们使用 `frozen_balance` 字段而非独立的托管账户表
- 财务平衡需要验证 **总资产（balance + frozen_balance）** 而非仅 balance

**验证公式（总资产对账）**：
```
Σ (balance + frozen_balance) = 所有 RECHARGE 交易总和
                              + 所有 TASK_SUBMIT 交易总和
                              + 所有 STAFF_REFERRAL 交易总和
                              + 所有 PROVIDER_INCOME 交易总和
                              + 所有 BONUS_GIFT 交易总和
                              + 所有 TASK_ESCALATE 交易总和
                              + 所有 TASK_REFUND 交易总和
                              + 所有 WITHDRAW_REFUND 交易总和
```

**验证公式（可用余额对账）**：
```
Σ balance = 所有 RECHARGE 交易总和
         + 所有 TASK_INCOME 交易总和
         + 所有 STAFF_REFERRAL 交易总和
         + 所有 PROVIDER_INCOME 交易总和
         + 所有 BONUS_GIFT 交易总和
         + 所有 TASK_ESCALATE 交易总和
         + 所有 TASK_REFUND 交易总和
         + 所有 WITHDRAW_REFUND 交易总和
         - 所有 TASK_PUBLISH 交易总和
         - 所有 TASK_ACCEPT 交易总和
         - 所有 WITHDRAW 交易总和
         - 所有 WITHDRAW_FREEZE 交易总和
         + 所有 TASK_REJECT 交易总和
```

**说明**：
- **总资产对账**：验证账户实际拥有的总积分数（可用+冻结）
  - 所有内部转账（balance ↔ frozen_balance）不影响总资产
  - TASK_PUBLISH, TASK_INCOME, TASK_ACCEPT 等是内部转账，不计入总资产变化
- **可用余额对账**：验证可随时使用的余额
  - 需要考虑 frozen_balance 中待分配的资金

**交易类型对总资产的影响**：
| 交易类型 | 对总资产影响 | 说明 |
|---------|-------------|------|
| RECHARGE | + | 外部充值，增加总资产 |
| TASK_SUBMIT | 0 | 达人内部转账（系统记账），不影响总资产 |
| TASK_INCOME | 0 | 达人内部转账，不影响总资产 |
| STAFF_REFERRAL | 0 | 员工内部转账，不影响总资产 |
| PROVIDER_INCOME | 0 | 服务商内部转账，不影响总资产 |
| TASK_PUBLISH | 0 | 商家内部转账，不影响总资产 |
| TASK_ACCEPT | 0 | 商家内部记账，不影响总资产 |
| TASK_REJECT | 0 | 达人内部转账，不影响总资产 |
| TASK_ESCALATE | + | 退还商家费用，增加总资产 |
| TASK_REFUND | + | 退还商家费用，增加总资产 |
| WITHDRAW | - | 外部提现，减少总资产 |
| WITHDRAW_FREEZE | 0 | 内部转账，不影响总资产 |
| WITHDRAW_REFUND | + | 提现拒绝，返还费用，增加总资产 |
| BONUS_GIFT | + | 系统赠送，增加总资产 |

**2. frozen_balance 对账**
```
验证公式（按账户类型）：

商家 frozen_balance  = Σ (所有 OPEN 任务的未接活动名额费用)
                    + Σ (所有 ASSIGNED 状态的活动名额费用)
                    + Σ (所有 SUBMITTED 状态的活动名额费用)

达人 frozen_balance = Σ (所有 SUBMITTED 状态的活动名额预计达人收入)
                    + Σ (所有提现申请通过待打款的金额)

员工 frozen_balance = Σ (所有 SUBMITTED 状态的活动名额预计员工返佣)
                    + Σ (所有提现申请通过待打款的金额)

服务商 frozen_balance = Σ (所有提现申请通过待打款的金额)
```

**说明**：
- **商家 frozen_balance**：预扣但未完成分配的活动费用（**负债性质**）
  - OPEN 任务：未接活动名额 × 任务金额
  - ASSIGNED 接任务：已接任务但未提交的费用
  - SUBMITTED 接任务：已提交但未审核通过的费用（包含达人、员工、服务商三部分）

- **达人/员工 frozen_balance**：已提交但待审核的任务收入（**资产性质**）
  - SUBMITTED 接任务的预计收入
  - 提现申请通过后待打款的金额

- **服务商 frozen_balance**：提现申请通过后待打款的金额
  - 服务商收入直接进入 balance，不经过 frozen_balance

**⚠️ 重要：frozen_balance的财务语义差异**
- **商家/服务商**：frozen_balance 是**负债**（预扣待支付）
  - 对账时验证：`frozen_balance` 应该等于待支付总额
  - 符号处理：在财务报表中显示为"应付账款"（负数）

- **达人/员工**：frozen_balance 是**资产**（待确认收入）
  - 对账时验证：`frozen_balance` 应该等于待收入总额
  - 符号处理：在财务报表中显示为"应收账款"（正数）

- **系统内部转账平衡**：
  ```
  商家 frozen_balance（负债）+ 达人 frozen_balance（资产）+ 员工 frozen_balance（资产） = 0
  ```
  如果不为0，说明存在资金流转异常，需要告警

**3. 现金对账**
```
验证公式：
系统现金余额 = 充值总额 - 提现总额 - 退款总额
```

#### 对账频率
- **自动对账**：每日凌晨自动执行
- **手动对账**：超级管理员可随时触发
- **异常告警**：对账发现差异时，发送告警通知

#### 对账报告内容
- 总体账户余额 vs 流水计算余额
- 差异明细列表
- 异常交易记录
- 建议：根据差异类型生成处理建议

---

### 4. 财务对账页面设计

#### 页面路由
- **路径**：`/workspace/admin/finance/reconciliation`
- **权限**：仅超级管理员可访问

#### 页面布局
```
┌──────────────────────────────────────────────────────────────────────┐
│ 顶部导航栏                                                            │
│ [财务对账]                                                            │
├──────────────────────────────────────────────────────────────────────┤
│ 操作栏                                                                │
│ [开始对账] [下载报告] 对账日期范围：[2024-01-01] 至 [2024-01-31]      │
├──────────────────────────────────────────────────────────────────────┤
│ 对账状态卡片                                                          │
│ ┌─────────────┬─────────────┬─────────────┬─────────────┐           │
│ │ 积分平衡对账 │ frozen对账  │ 现金对账    │ 总体状态    │           │
│ │ ✅ 已通过    │ ⚠️ 2个差异  │ ✅ 已通过    │ ⚠️ 有异常    │           │
│ └─────────────┴─────────────┴─────────────┴─────────────┘           │
├──────────────────────────────────────────────────────────────────────┤
│                                                                        │
│ 【积分平衡对账】                                                       │
│ ┌─────────────────────────────────────────────────────────────────┐  │
│ │ 对账类型：积分平衡对账                                            │  │
│ │ 对账时间：2024-01-31 03:00:00                                    │  │
│ │ 对账结果：✅ 通过                                                 │  │
│ │                                                                  │  │
│ │ 实际总余额：1,234,567 积分                                        │  │
│ │ 计算总余额：1,234,567 积分                                        │  │
│ │ 差异：0 积分                                                      │  │
│ │                                                                  │  │
│ │ 流水汇总：                                                        │  │
│ │ • RECHARGE（充值）：+500,000 积分                                │  │
│ │ • TASK_INCOME（任务收入）：+300,000 积分                         │  │
│ │ • STAFF_REFERRAL（员工返佣）：+50,000 积分                       │  │
│ │ • PROVIDER_INCOME（服务商收入）：+80,000 积分                    │  │
│ │ • BONUS_GIFT（系统赠送）：+20,000 积分                           │  │
│ │ • TASK_PUBLISH（发布任务）：-400,000 积分                        │  │
│ │ • TASK_ACCEPT（接任务扣除）：-250,000 积分                         │  │
│ │ • TASK_ESCALATE（超时拒绝）：-5,000 积分                         │  │
│ │ • TASK_REFUND（任务退款）：+40,000 积分                          │  │
│ │ • WITHDRAW（提现）：-100,000 积分                                 │  │
│ │ • WITHDRAW_REFUND（提现拒绝）：+500 积分                          │  │
│ └─────────────────────────────────────────────────────────────────┘  │
│                                                                        │
│ 【frozen_balance 对账】                                                │
│ ┌─────────────────────────────────────────────────────────────────┐  │
│ │ 对账类型：frozen_balance 对账                                     │  │
│ │ 对账时间：2024-01-31 03:00:00                                    │  │
│ │ 对账结果：⚠️ 发现 2 个差异                                         │  │
│ │                                                                  │  │
│ │ 差异明细：                                                        │  │
│ │ ┌──────┬───────────────┬─────────────┬────────────┬──────────┐ │  │
│ │ │账户ID│ 账户类型      │ 实际冻结    │ 计算冻结   │ 差异     │ │  │
│ │ ├──────┼───────────────┼─────────────┼────────────┼──────────┤ │  │
│ │ │12345 │ ORG_MERCHANT  │ 10,000      │ 9,000      │ +1,000   │ │  │
│ │ │67890 │ USER_PERSONAL │ 5,000       │ 5,500      │ -500     │ │  │
│ │ └──────┴───────────────┴─────────────┴────────────┴──────────┘ │  │
│ │                                                                  │  │
│ │ [查看账户12345详情] [查看账户67890详情]                            │  │
│ └─────────────────────────────────────────────────────────────────┘  │
│                                                                        │
│ 【异常交易记录】                                                       │
│ ┌─────────────────────────────────────────────────────────────────┐  │
│ │ 发现 3 条异常交易：                                               │  │
│ │ ┌──────┬──────────┬────────────┬──────────────┬────────┐       │  │
│ │ │交易ID│ 交易类型 │ 金额       │ 时间         │ 异常   │       │  │
│ │ ├──────┼──────────┼────────────┼──────────────┼────────┤       │  │
│ │ │10001 │ TASK_... │ -100       │ 01-30 15:23  │ 重复   │       │  │
│ │ │10002 │ WITHDRAW │ -1,000     │ 01-30 16:45  │ 余额不 │       │  │
│ │ │      │          │            │              │ 足     │       │  │
│ │ │10003 │ TASK_... │ +50        │ 01-31 09:12  │ 金额负 │       │  │
│ │ │      │          │            │              │ 数异常 │       │  │
│ │ └──────┴──────────┴────────────┴──────────────┴────────┘       │  │
│ │                                                                  │  │
│ │ [导出异常交易] [批量修复]                                         │  │
│ └─────────────────────────────────────────────────────────────────┘  │
│                                                                        │
│ 【处理建议】                                                           │
│ ┌─────────────────────────────────────────────────────────────────┐  │
│ │ 1. 账户12345（ORG_MERCHANT）frozen_balance差异 +1,000 积分       │  │
│ │    建议操作：手动调整 frozen_balance -= 1,000                    │  │
│ │    可能原因：任务退款时未正确释放 frozen_balance                  │  │
│ │    [立即修复]                                                     │  │
│ │                                                                  │  │
│ │ 2. 交易10001（TASK_ACCEPT）疑似重复结算                          │  │
│ │    建议操作：检查任务 #12345 的活动名额 #67890                   │  │
│ │    可能原因：并发审核导致重复分配积分                             │  │
│ │    [查看详情]                                                     │  │
│ └─────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────┘
```

#### 页面功能详解

**1. 操作栏**
- **开始对账**：手动触发对账任务
  - 点击后显示确认对话框："确定要开始财务对账吗？"
  - 对账期间显示进度条
  - 对账完成后自动刷新页面显示结果
- **下载报告**：导出当前对账结果的 PDF 或 Excel 报告
- **对账日期范围**：默认显示本月，可选择历史日期范围

**2. 对账状态卡片**
显示三种对账类型的简要结果：
- ✅ 已通过：绿色勾号，表示无差异
- ⚠️ 有差异：黄色警告，显示差异数量
- ❌ 失败：红色叉号，表示对账过程出错

点击卡片可滚动到对应的详细结果区域。

**3. 积分平衡对账详情**
**显示内容**：
- 对账基本信息（时间、结果）
- 实际总余额 vs 计算总余额
- 差异金额
- 流水汇总表格（11种交易类型的金额汇总）

**交互**：
- 点击交易类型可展开显示该类型的所有交易记录
- 如果发现差异，显示红色警告和差异金额

**4. frozen_balance 对账详情**
**显示内容**：
- 对账基本信息
- 差异账户列表（表格）
  - 账户ID、账户类型
  - 实际冻结余额、计算冻结余额
  - 差异金额

**交互**：
- 点击"查看账户详情"跳转到账户详情页面
- 显示该账户所有相关的 frozen_balance 变动记录

**5. 现金对账详情**
**显示内容**：
- 对账基本信息
- 充值金额 vs 实际到账金额对比
- 提现金额 vs 实际打款金额对比

**6. 异常交易记录**
**异常类型**：
- **重复交易**：同一任务/提现有多条结算记录
- **余额不足**：交易金额超过账户余额
- **金额异常**：负数充值、正数提现等
- **状态异常**：已取消的任务有交易记录

**交互**：
- 点击交易ID查看交易详情
- 支持批量选择异常交易进行修复
- "导出异常交易"导出 CSV 文件

**7. 处理建议**
**自动生成建议**：
- 根据差异类型自动生成修复建议
- 显示可能的原因分析
- 提供"立即修复"或"查看详情"操作

**修复操作类型**：
- **手动调整**：创建一笔调整交易修正余额
- **查看详情**：跳转到相关任务/提现/账户详情页面
- **标记忽略**：标记该差异为已忽略（需填写原因）

#### 对账触发方式

**1. 自动对账**
- **触发时间**：每日凌晨 3:00
- **执行方式**：后台定时任务
- **完成通知**：
  - 对账通过：发送通知给超级管理员
  - 发现差异：发送告警邮件 + 站内通知

**2. 手动对账**
- **触发入口**：
  - 财务对账页面的"开始对账"按钮
  - 超级管理员可随时触发
- **对账范围**：
  - 全量对账：检查所有账户和交易
  - 增量对账：只检查上次对账后的新交易

**3. 异常告警**
- **告警触发条件**：
  - 任何对账类型发现差异
  - 异常交易数量 > 0
  - 对账执行失败

- **告警渠道**：
  - 站内通知
  - 邮件通知（发送给所有超级管理员）
  - 可选：钉钉/企业微信机器人

- **告警内容**：
  - 对账类型
  - 差异金额/数量
  - 快速跳转链接

#### 页面权限

**访问权限**：
- 仅超级管理员（SUPER_ADMIN）可访问

**操作权限**：
- **查看对账结果**：所有超级管理员
- **手动触发对账**：所有超级管理员
- **下载报告**：所有超级管理员
- **手动调整余额**：所有超级管理员（需二次确认）
- **标记忽略差异**：所有超级管理员（需填写原因）

#### API 接口设计

**1. 获取对账结果**
```
GET /api/admin/reconciliation/result
Query: ?date_from=2024-01-01&date_to=2024-01-31
Response:
{
  "balance_check": {
    "status": "passed",
    "actual_balance": 1234567,
    "calculated_balance": 1234567,
    "difference": 0,
    "transaction_summary": { ... }
  },
  "frozen_check": {
    "status": "warning",
    "diff_count": 2,
    "diff_accounts": [ ... ]
  },
  "cash_check": {
    "status": "passed",
    ...
  }
}
```

**2. 触发对账**
```
POST /api/admin/reconciliation/trigger
Body: {
  "type": "full",  // full | incremental
  "date_from": "2024-01-01",
  "date_to": "2024-01-31"
}
Response: {
  "job_id": "recon_20240131_030000",
  "status": "running",
  "message": "对账任务已开始，请稍后查看结果"
}
```

**3. 获取异常交易**
```
GET /api/admin/reconciliation/anomalies
Query: ?date_from=2024-01-01&date_to=2024-01-31&type=DUPLICATE
Response: {
  "total": 3,
  "transactions": [ ... ]
}
```

**4. 修复差异**
```
POST /api/admin/reconciliation/fix
Body: {
  "account_id": "12345",
  "adjustment_amount": -1000,
  "reason": "对账发现frozen_balance差异，手动调整"
}
Response: {
  "status": "success",
  "adjustment_id": "adj_20240131_001"
}
```

---

### 🔄 营销活动状态流程

**⚠️ 重要概念区分**：
- **任务**：商家/服务商发布的任务
- **活动名额**：达人接任务后创建的记录

---

#### 营销活动状态

| 状态代码 | 状态名称 | 说明 | 可操作角色 |
|---------|---------|------|-----------|
| **DRAFT** | 草稿 | 商家创建的草稿活动，未发布 | 商家、服务商 |
| **OPEN** | 开放中 | 服务商设置积分分配并发布，等待达人接任务 | 达人、服务商 |
| **CLOSED** | 已关闭 | 名额已满或手动关闭（截止时间到不自动关闭） | - |

---

#### 活动名额状态

| 状态代码 | 状态名称 | 说明 | 可操作角色 |
|---------|---------|------|-----------|
| **ASSIGNED** | 进行中 | 达人已接任务，等待提交平台链接/截图 | 达人、服务商 |
| **SUBMITTED** | 待审核 | 达人已提交平台链接/截图，等待服务商审核 | 服务商 |
| **APPROVED** | 已完成 | 服务商审核通过，积分已分配 | - |
| **REJECTED** | 已拒绝 | 服务商审核拒绝，达人可重新提交 | 达人 |

---

#### 完整状态流转图

```
商家创建任务
    ↓
[1. DRAFT - 草稿]
    创建者：商家
    可见：商家、绑定的服务商
    操作：编辑、删除
    ↓
服务商设置积分分配并发布
    ↓
[2. OPEN - 开放中]
    发布者：服务商
    可见：商家、服务商、所有达人
    活动名额：10人
    已接任务：3人
    剩余名额：7人
    操作：达人接任务
    ↓
达人接任务（创建活动名额）

**⚠️ 名额已满自动关闭**：
- 当已接名额数 = 总名额时（如：10人已满）
- 系统自动将营销活动状态设置为 CLOSED
- 未接任务达人无法再接任务
    ↓
    ↓
[3. ASSIGNED - 进行中]
    接任务者：达人
    活动名额状态：ASSIGNED
    可见：接任务达人、商家、服务商
    操作：达人提交平台链接/截图
    ↓
达人提交平台链接/截图
    ↓
[4. SUBMITTED - 待审核]
    提交者：达人
    活动名额状态：SUBMITTED
    可见：接任务达人、商家、服务商
    操作：服务商审核
    ↓
服务商审核
    ├─ 通过
    │   ↓
    │ [5. APPROVED - 已完成]
    │   活动名额状态：APPROVED
    │   操作：分配积分
    │   - 达人积分账户 + creator_amount（如：80元）
    │   - 邀请人积分账户 + staff_referral_amount（如：10元，如果有邀请人）
    │   - 服务商积分账户 + provider_amount（如：10元）
    │   - 任务已接名额数 +1
    │
    └─ 拒绝
        ↓
        [6. REJECTED - 已拒绝]
        活动名额状态：REJECTED
        操作：填写拒绝原因
        资金处理：
        - 达人 frozen_balance -= 预计收入（释放冻结）
        - 商家 frozen_balance 继续保持冻结（该笔费用不退回 balance）
        - 达人可重新提交 → 回到 [4. SUBMITTED]
        - 重新提交后，资金再次从商家 frozen_balance → 达人 frozen_balance
        - 审核通过后才完成最终分配
```

---

## 💰 资金流动规则

### 发布任务（预扣款模式）
**服务商设置积分分配并发布任务**：
- 商家 `balance -= 总费用`（如：10人×100元=1000元）
- 商家 `frozen_balance += 总费用`（如：+1000元）
- 营销活动状态：DRAFT → OPEN

### 达人接任务
**达人接任务**：
- 商家 `frozen_balance -= 任务金额`（如：-100元）
- 标记为"已接任务，待分配"
- 活动名额状态：ASSIGNED

### 达人提交（冻结待审核）
**达人提交任务**：
- 商家 `frozen_balance -= 达人预计收入`（如：-80元，转移到达人）
- 达人 `frozen_balance += 预计收入`（如：+80元，从商家转入）
- 活动名额状态：ASSIGNED → SUBMITTED

**⚠️ 说明**：达人提交时，预计收入从商家的 frozen_balance 转移到达人的 frozen_balance

### 服务商审核通过
**审核通过**：
- 达人 `frozen_balance -= 80元`，`balance += 80元`（内部转账）
- 员工 `frozen_balance += 10元`，`balance += 10元`（从商家 frozen_balance 转入）
- 服务商 `balance += 10元`（从商家 frozen_balance 转入）
- 商家 `frozen_balance -= 20元`（员工返佣 + 服务商收入）
- 活动名额状态：SUBMITTED → APPROVED

**⚠️ 说明**：员工返佣和服务商收入都从商家的 frozen_balance 转入

### 服务商拒绝
**审核拒绝**：
- 达人 `frozen_balance -= 预计收入`（如：-80元，释放冻结）
- 商家 `frozen_balance` 继续保持冻结（等待重新提交）
- 活动名额状态：SUBMITTED → REJECTED
- 达人可修改后重新提交 → 回到 SUBMITTED

**⚠️ 说明**：
- 拒绝时，达人 frozen_balance 中的预计收入被释放
- 该笔费用继续冻结在商家的 frozen_balance 中（不退回到 balance）
- 达人重新提交后，资金再次到达人 frozen_balance
- 只有审核通过后，才完成最终分配（达人 balance+、员工 balance+、服务商 balance+）

### 任务手动关闭或名额已满（自动退款）
**任务关闭**：
- 计算未完成名额：总名额 - 已接名额数
- 退款金额 = 未完成名额 × 任务金额
- 商家 `frozen_balance -= 退款金额`
- 商家 `balance += 退款金额`
- 创建 REFUND 类型交易记录

**示例**：
- 活动名额：10人，任务金额100积分，总计1000元
- 已接任务：8人（3人已完成APPROVED，5人进行中ASSIGNED/SUBMITTED）
- 未接任务：2人
- 退款：200元（`frozen_balance → balance`）

**⚠️ 重要说明**：
- **未完成名额 = 总名额 - 已接名额数**（不是已完成数）
- 已接任务但未完成的活动名额（ASSIGNED/SUBMITTED状态），其费用仍在 frozen_balance 中
- 只有未接任务的名额费用才退还给商家
- 已接任务的名额，无论是否完成，都不参与退款（已接任务的费用已经在达人接任务和提交时冻结/转移）

---

#### 截止时间规则

**⚠️ 重要：任务有2个独立的截止时间**

### 1. 接任务截止时间（task_deadline）

**作用**：达人可以接任务的最后期限

**到期处理**：
- ✅ 接任务截止时间到达前，达人可以接任务
- ✅ 接任务截止时间到达后，任务标记为"已到期"状态，但**不会自动关闭**
- ✅ 已接任务的达人不受影响，可以继续提交任务
- ✅ 未接任务的名额可以继续招募
- ✅ 服务商可以延长接任务截止时间，任务恢复为"开放中（OPEN）"状态

**服务商操作**：
- **延长接任务截止时间**：服务商可以延长
  - 营销活动状态从"已到期"恢复为"开放中（OPEN）"
  - 设置新的 `task_deadline`
  - 达人可以继续接任务

---

### 2. 提交截止时间（submission_deadline）

**作用**：已接任务达人必须提交任务内容的最后期限

**到期处理**：
- ✅ 提交截止时间到达前，ASSIGNED 状态的达人可以提交任务
- ⚠️ **提交截止时间到达后，所有未提交的活动名额（ASSIGNED）自动拒绝**
- ⚠️ 自动拒绝的活动名额状态变为 REJECTED
- ⚠️ 冻结在达人 frozen_balance 中的预计收入释放
- ⚠️ 对应的商家 frozen_balance 中的费用退回商家 balance
- ✅ 服务商可以延长提交截止时间，给已接任务达人更多时间

**服务商操作**：
- **延长提交截止时间**：服务商可以延长
  - 设置新的 `submission_deadline`
  - 已接任务达人可以继续提交
  - 未接任务达人不受影响

---

### 3. 两个截止时间的关系

```
时间轴示例：
1月1日   ---------> 1月25日 18:00 ---------> 2月1日 18:00
任务发布            接任务截止             提交截止

1月1日-1月25日：达人可以接任务
1月25日之后：   标记"已到期"，但已接任务的可继续提交，未接任务的不能接任务
1月1日-2月1日： 已接任务达人可以提交任务
2月1日之后：   未提交的活动名额自动拒绝
```

---

### 4. 提交截止时间到后的退款逻辑

**场景示例**：
- 活动名额：10人，任务金额100积分，总计1000元
- 已接任务：8人（3人已完成APPROVED，5人进行中ASSIGNED）
- 提交截止时间到：5个ASSIGNED未提交

**自动拒绝处理（ASSIGNED状态）**：
- 5个ASSIGNED活动名额 → 自动变为REJECTED
- 商家 frozen_balance：退回全部费用（5 × 100 = 500元）
- 达人 frozen_balance：无需操作（ASSIGNED状态未提交，frozen_balance中没有这笔钱）

**自动拒绝处理（SUBMITTED状态）**：
- 如果是SUBMITTED状态（已提交待审核），超时不会自动拒绝
- 服务商可以继续审核，通过或拒绝都可以

**任务关闭时的退款**：
- 假设任务手动关闭，未接任务：2人
- 退款金额 = 2人 × 100元 = 200元
- 总退款 = 自动拒绝退回的500元 + 未接任务退回的200元 = 700元

---

### 5. 系统提示

**接任务截止时间提醒**：
```
┌─────────────────────────────────────────┐
│ ⚠️ 接任务截止时间提醒                     │
├─────────────────────────────────────────┤
│ 任务：新品体验推广                       │
│ 接任务截止：2024-01-25 18:00             │
│ 剩余时间：2天                           │
│ 已接任务：6/10                             │
│ 剩余名额：4/10                           │
│                                         │
│ [延长接任务截止时间]                      │
└─────────────────────────────────────────┘
```

**提交截止时间提醒**：
```
┌─────────────────────────────────────────┐
│ ⚠️ 提交截止时间提醒                     │
├─────────────────────────────────────────┤
│ 任务：新品体验推广                       │
│ 提交截止：2024-02-01 18:00             │
│ 剩余时间：2天                           │
│ 已提交：3/6                              │
│ 未提交：3/6（即将自动拒绝）             │
│                                         │
│ [延长提交截止时间]                      │
│                                         │
│ ⚠️ 提交截止后未提交的将自动拒绝          │
└─────────────────────────────────────────┘
```

---

**⚠️ 重要设计原则**：
1. 两个截止时间**独立设置、独立延长**
2. **接任务截止时间到**：不影响已接任务达人
3. **提交截止时间到**：自动拒绝未提交的活动名额，释放资金
4. 避免达人无限期不提交导致的资金冻结问题

---

#### 审核拒绝重新提交流程

```
[SUBMITTED - 待审核]
    达人提交：平台链接 + 截图
    ↓
服务商审核：拒绝
    填写原因："平台链接无法访问"
    ↓
[REJECTED - 已拒绝]
    活动名额状态保持 REJECTED
    系统通知达人审核失败
    达人收到拒绝原因
    ↓
达人查看拒绝原因并修改内容：
    - 检查并修改平台链接
    - 重新上传截图
    - 添加备注说明
    ↓
[SUBMITTED - 待审核]（重新提交）
    活动名额状态从 REJECTED 变为 SUBMITTED
    达人可以无限次重新提交
    但必须在任务截止时间之前
    ↓
服务商重新审核
    ↓
[APPROVED - 已完成] 或 [REJECTED - 已拒绝]
```

**⚠️ 重要规则**：
- ✅ 达人可以**无限次重新提交**
- ✅ **截止时间到达后不能提交**
- ✅ 服务商可以手动处理未完成的活动名额

---

### 字段类型说明

- **UUID**: 通用唯一标识符，PostgreSQL 使用 `UUID` 类型
- **VARCHAR(n)**: 变长字符串，最大长度 n
- **TEXT**: 长文本，无长度限制
- **INT**: 整数，PostgreSQL 使用 `INTEGER`
- **BOOLEAN**: 布尔值，TRUE/FALSE
- **TIMESTAMP**: 时间戳，包含日期和时间
- **JSONB**: JSON 二进制格式，支持高效查询

---

### ⚠️ 待确认问题

1. **表名和字段名**是否符合你的习惯？
2. **前端显示名**是否准确？
3. **字段类型**是否合理？
4. 是否需要**增加字段**或**删除字段**？
5. **索引设计**是否合理？
6. **外键约束**是否需要？

---

## 🎯 P0-1: 非功能性需求（注意非功能性需求暂时不用考虑）

### 1. 性能要求

#### 1.1 并发用户数

| 场景 | 并发用户数 | 说明 |
|------|-----------|------|
| **正常负载** | 500 DAU | 日常同时在线用户数 |
| **高峰负载** | 2000 DAU | 活动推广期间峰值 |
| **压力测试** | 5000 DAU | 系统稳定性测试目标 |

#### 1.2 响应时间要求

| 操作类型 | 响应时间要求 | 说明 |
|---------|------------|------|
| **页面加载** | < 2秒 | 首屏渲染完成 |
| **API接口** | < 500ms | 普通查询接口 |
| **复杂查询** | < 1秒 | 涉及多表关联、聚合的查询 |
| **文件上传** | < 3秒 | 单个图片/截图上传（< 5MB） |
| **任务审核** | < 1秒 | 审核操作提交到完成 |

#### 1.3 数据量要求

| 数据类型 | 预估数据量（1年） | 存储空间 | 说明 |
|---------|-----------------|---------|------|
| **用户数** | 10,000 用户 | ~50 MB | 含角色、关联关系 |
| **商家数** | 500 商家 | ~10 MB | 含员工、权限 |
| **服务商数** | 50 服务商 | ~1 MB | 含员工、权限 |
| **营销活动** | 5,000 活动 | ~100 MB | 含任务名额 |
| **任务名额** | 50,000 名额 | ~500 MB | 含提交数据 |
| **积分流水** | 200,000 条 | ~200 MB | 含交易记录 |
| **提现记录** | 5,000 条 | ~10 MB | 含审核记录 |
| **总计** | - | ~1 GB | 不含文件存储 |

### 2. 可用性要求

#### 2.1 系统可用性

| 指标 | 目标值 | 说明 |
|------|-------|------|
| **服务可用性（SLA）** | 99.9% | 年度停机时间 < 8.76 小时 |
| **数据备份频率** | 每日 | 凌晨 2:00 自动备份 |
| **故障恢复时间（RTO）** | < 1 小时 | 从故障到恢复服务的时间 |
| **数据恢复点（RPO）** | < 24 小时 | 最多丢失一天的数据 |
| **灾难恢复（DR）** | 72 小时 | 完整灾难恢复时间 |

#### 2.2 容错机制

- **数据库主从复制**：1主2从，自动故障转移
- **应用服务多实例**：至少2个实例，负载均衡
- **Redis 集群**：主从 + 哨兵模式
- **CDN 加速**：静态资源（图片、JS、CSS）使用 CDN
- **健康检查**：每30秒检查服务健康状态，异常自动重启

### 3. 安全要求

#### 3.1 认证与授权

| 安全措施 | 实现方式 | 说明 |
|---------|---------|------|
| **统一认证** | auth-center (os.crazyaigc.com) | 微信登录 + 手机/密码登录 |
| **Token 管理** | JWT + Redis 黑名单 | Access Token 24小时，Refresh Token 7天 |
| **权限校验** | RBAC + 资源级权限 | 每个接口校验用户权限 |
| **密码策略** | 8位以上，含字母数字 | 仅密码登录用户适用 |
| **会话超时** | 24小时无操作自动登出 | 前端 + 后端双重校验 |

#### 3.2 数据安全

| 安全措施 | 实现方式 | 说明 |
|---------|---------|------|
| **HTTPS 强制** | 全站 HTTPS (TLS 1.3) | 防止中间人攻击 |
| **敏感数据加密** | 银行卡号 AES-256 加密 | 提现账户信息加密存储 |
| **SQL 注入防护** | 参数化查询 + ORM | 所有数据库操作使用 Sequelize ORM |
| **XSS 防护** | 输入过滤 + 输出转义 | 前端 React 自动转义 |
| **CSRF 防护** | CSRF Token + SameSite Cookie | 所有修改操作需验证 Token |
| **文件上传限制** | 文件类型白名单 + 大小限制 | 仅允许 jpg/png，单文件 < 5MB |

#### 3.3 审计与监控

- **操作日志**：所有敏感操作记录（余额变动、角色变更、数据删除）
- **审计日志永久保存**：不可删除，不可修改
- **异常登录监控**：异地登录、多次失败密码尝试触发告警
- **财务操作二次确认**：调整余额、封禁用户需要二次验证

### 4. 兼容性要求

#### 4.1 浏览器兼容性

| 浏览器 | 最低版本 | 说明 |
|-------|---------|------|
| **Chrome** | ≥ 90 | 主要支持浏览器 |
| **Safari** | ≥ 14 | iOS/macOS 支持 |
| **Edge** | ≥ 90 | Chromium 内核 |
| **Firefox** | ≥ 88 | 基本支持 |
| **微信内置浏览器** | 最新版 | 微信登录必需 |

**不支持**：IE 浏览器（任何版本）

#### 4.2 移动端支持

| 设备类型 | 支持情况 | 说明 |
|---------|---------|------|
| **iOS** | ✅ 完全支持 | iOS 14+ |
| **Android** | ✅ 完全支持 | Android 10+ |
| **响应式设计** | ✅ 自适应 | 手机、平板、桌面 |
| **PWA** | ⚠️ 可选 | 后续版本支持离线使用 |

### 5. 可维护性要求

#### 5.1 代码质量

| 指标 | 目标值 | 说明 |
|------|-------|------|
| **代码覆盖率** | ≥ 80% | 单元测试覆盖率 |
| **ESLint 规则** | 0 错误，< 10 警告 | 代码风格检查 |
| **TypeScript 严格模式** | 启用 | 类型安全 |
| **API 文档** | 自动生成 | 使用 Swagger/OpenAPI |

#### 5.2 日志与监控

- **应用日志**：结构化日志（JSON 格式），按日期滚动
- **错误追踪**：Sentry 集成，实时错误告警
- **性能监控**：APM 工具（如 New Relic）监控接口性能
- **数据库慢查询**：记录超过 1秒 的查询，优化索引

---

## 🏗️ P0-2: 系统架构设计

### 1. 技术栈选型

#### 1.1 前端技术栈

| 技术 | 版本 | 用途 | 说明 |
|------|------|------|------|
| **React** | 18.x | 前端框架 | 组件化开发 |
| **TypeScript** | 5.x | 类型系统 | 类型安全 |
| **Next.js** | 14.x | SSR 框架 | SEO 优化、性能优化 |
| **TailwindCSS** | 3.x | CSS 框架 | 快速 UI 开发 |
| **Zustand** | 4.x | 状态管理 | 轻量级全局状态 |
| **React Query** | 5.x | 数据获取 | 缓存、自动重试 |
| **React Hook Form** | 7.x | 表单管理 | 高性能表单验证 |
| **Ant Design** | 5.x | UI 组件库 | 企业级组件 |

#### 1.2 后端技术栈

| 技术 | 版本 | 用途 | 说明 |
|------|------|------|------|
| **Node.js** | 20.x | 运行时 | LTS 长期支持版本 |
| **TypeScript** | 5.x | 类型系统 | 类型安全 |
| **NestJS** | 10.x | 后端框架 | 模块化架构 |
| **Sequelize** | 6.x | ORM | 数据库操作 |
| **PostgreSQL** | 15.x | 关系数据库 | 主数据库 |
| **Redis** | 7.x | 缓存数据库 | 会话、缓存 |
| **Passport** | 0.6.x | 认证中间件 | OAuth 2.0 集成 |

#### 1.3 基础设施

| 技术 | 版本 | 用途 | 说明 |
|------|------|------|------|
| **Docker** | 24.x | 容器化 | 应用打包部署 |
| **Nginx** | 1.24.x | 反向代理 | 负载均衡、静态资源 |
| **PM2** | 5.x | 进程管理 | Node.js 进程守护 |
| **GitHub Actions** | - | CI/CD | 自动化测试、部署 |

### 2. 部署架构

#### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         用户浏览器                            │
│                    (React + Next.js)                         │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTPS
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                      Nginx 反向代理                          │
│              (负载均衡 + 静态资源 + SSL)                      │
└─────────┬─────────────────┬─────────────────────────────────┘
          │                 │
          ↓                 ↓
┌──────────────────┐  ┌──────────────────┐
│  前端服务器 #1    │  │  前端服务器 #2   │
│  (Next.js SSR)   │  │  (Next.js SSR)   │
└────────┬─────────┘  └────────┬─────────┘
         │                    │
         └────────┬───────────┘
                  ↓
┌─────────────────────────────────────────────────────────────┐
│                      API 网关层                              │
│                  (NestJS + Gateway)                         │
└─────────┬───────────────────────────────────────────────────┘
          │
          ↓
┌─────────┬─────────┬─────────┬──────────────────────────────┐
          │         │         │                               │
          ↓         ↓         ↓                               ↓
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  Auth Service│ │ Business API │ │  File Service│ │ Cron Service │
│  (认证模块)  │ │  (业务API)   │ │  (文件上传)  │ │  (定时任务)  │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │                │
       └────────────────┼────────────────┼────────────────┘
                        ↓                ↓
              ┌─────────────────┐  ┌─────────────────┐
              │   PostgreSQL    │  │     Redis       │
              │   (主从复制)     │  │   (主从 + 哨兵)  │
              │   1主2从         │  │   1主2从         │
              └─────────────────┘  └─────────────────┘
                        ↓                ↓
              ┌─────────────────┐  ┌─────────────────┐
              │   数据卷备份     │  │    Redis 持久化  │
              │   (每日备份)     │  │    (RDB + AOF)  │
              └─────────────────┘  └─────────────────┘
```

#### 2.2 服务器配置

**生产环境配置**（最小配置）：

| 服务 | 数量 | CPU | 内存 | 存储 | 说明 |
|------|------|-----|------|------|------|
| **Nginx** | 1 | 2核 | 4GB | 40GB SSD | 反向代理 |
| **前端服务器** | 2 | 4核 | 8GB | 40GB SSD | Next.js SSR |
| **后端服务器** | 2 | 4核 | 8GB | 40GB SSD | NestJS API |
| **PostgreSQL** | 3 | 8核 | 32GB | 500GB SSD | 1主2从 |
| **Redis** | 3 | 4核 | 16GB | 100GB SSD | 1主2从 |

**总成本估算**：约 ¥5000-8000/月（云服务器）

### 3. 服务模块设计

#### 3.1 后端模块划分

```
src/
├── main.ts                    # 应用入口
├── app.module.ts              # 根模块
├── common/                    # 通用模块
│   ├── guards/                # 守卫（权限校验）
│   ├── interceptors/          # 拦截器（日志、转换）
│   ├── pipes/                 # 管道（验证、转换）
│   ├── filters/               # 过滤器（异常处理）
│   └── decorators/            # 装饰器
├── auth/                      # 认证模块
│   ├── auth.module.ts
│   ├── auth.controller.ts     # 登录、登出、Token刷新
│   ├── auth.service.ts
│   └── strategies/            # Passport 策略（微信、密码）
├── users/                     # 用户模块
├── merchants/                 # 商家模块
├── service-providers/         # 服务商模块
├── creators/                  # 达人模块
├── campaigns/                 # 营销活动模块
├── tasks/                     # 任务模块
├── credits/                   # 积分模块
├── withdrawals/               # 提现模块
├── invite-codes/              # 邀请码模块
├── files/                     # 文件上传模块
└── database/                  # 数据库模块
    ├── migrations/            # 数据库迁移
    └── seeds/                 # 种子数据
```

#### 3.2 前端模块划分

```
src/
├── pages/                     # Next.js 页面
│   ├── workspace/             # 工作台页面
│   │   ├── admin/             # 超级管理员工作台
│   │   ├── merchant/          # 商家工作台
│   │   ├── service-provider/  # 服务商工作台
│   │   └── creator/           # 达人工作台
│   ├── auth/                  # 认证页面（登录）
│   └── _app.tsx               # 应用入口
├── components/                # 组件
│   ├── common/                # 通用组件
│   ├── tables/                # 多维表组件
│   └── forms/                 # 表单组件
├── hooks/                     # 自定义 Hooks
├── stores/                    # Zustand 状态管理
├── services/                  # API 服务
├── utils/                     # 工具函数
├── constants/                 # 常量（枚举、配置）
└── types/                     # TypeScript 类型定义
```

### 4. 数据库设计要点

#### 4.1 索引策略

```sql
-- 1. 主键索引（自动创建）
PRIMARY KEY (id)

-- 2. 唯一索引（防止重复）
UNIQUE INDEX idx_table_name ON table(column) WHERE deleted_at IS NULL;

-- 3. 普通索引（加速查询）
INDEX idx_table_column ON table(column);

-- 4. 复合索引（多条件查询）
INDEX idx_table_col1_col2 ON table(col1, col2);

-- 5. JSONB 索引（支持 JSON 查询）
INDEX idx_table_json ON table USING GIN(jsonb_column);

-- 6. 部分索引（减少索引大小）
INDEX idx_table_active ON table(status) WHERE status = 'ACTIVE';
```

#### 4.2 查询优化

- **避免 SELECT \***：仅查询需要的字段
- **使用 LIMIT**：分页查询，避免一次性加载大量数据
- **JOIN 优化**：优先使用 INNER JOIN，减少关联表数量
- **慢查询监控**：记录超过 1秒 的查询，分析并优化
- **连接池配置**：
  ```javascript
  {
    max: 100,              // 最大连接数
    min: 10,               // 最小空闲连接数
    acquireTimeoutMillis: 30000,  // 获取连接超时时间
    idle: 10000,           // 连接空闲时间
  }
  ```

### 5. 缓存策略

#### 5.1 缓存层次

| 缓存类型 | 存储位置 | 过期时间 | 用途 |
|---------|---------|---------|------|
| **浏览器缓存** | 用户浏览器 | 静态资源1年 | 静态资源（JS、CSS、图片） |
| **CDN 缓存** | CDN 节点 | 静态资源30天 | 加速静态资源访问 |
| **Redis 缓存** | 服务端 Redis | 见下表 | 热点数据、会话数据 |

#### 5.2 Redis 缓存策略

| 数据类型 | Key 格式 | 过期时间 | 说明 |
|---------|---------|---------|------|
| **用户会话** | `session:{userId}` | 24小时 | 用户登录状态 |
| **邀请码** | `invite_code:{code}` | 1小时 | 邀请码验证结果 |
| **统计数据** | `stats:{type}:{date}` | 5分钟 | 财务统计、任务统计 |
| **用户信息** | `user:{userId}` | 30分钟 | 用户基础信息 |
| **活动列表** | `campaigns:list:{page}` | 1分钟 | 达人任务大厅 |
| **配置数据** | `config:{key}` | 永久 | 系统配置、枚举值 |

---

## 📅 P0-3: 开发计划

### 第一阶段：MVP 核心功能（8周）

#### Week 1-2: 基础设施与认证

| 任务 | 工作量 | 交付物 | 负责人 |
|------|-------|-------|-------|
| **项目初始化** | 2天 | 前后端项目脚手架、Docker 配置 | 后端 |
| **数据库设计** | 3天 | 完整数据库表结构、索引、约束 | 后端 |
| **认证模块开发** | 5天 | 微信登录、手机/密码登录、JWT Token | 后端 |
| **前端路由与布局** | 3天 | 工作台布局、角色切换器、路由守卫 | 前端 |
| **前端通用组件** | 2天 | 按钮、表单、表格、弹窗 | 前端 |

**里程碑**：用户可以登录系统，看到对应的工作台

#### Week 3-4: 用户与角色管理

| 任务 | 工作量 | 交付物 | 负责人 |
|------|-------|-------|-------|
| **商家管理** | 3天 | 商家 CRUD、商家员工管理、权限分配 | 后端 |
| **服务商管理** | 3天 | 服务商 CRUD、员工管理、权限分配 | 后端 |
| **达人管理** | 2天 | 达人资料、达人等级、邀请达人 | 后端 |
| **超级管理员功能** | 2天 | 用户管理、角色管理、数据查看 | 后端 |
| **用户管理页面** | 4天 | 用户列表、角色管理、权限管理 | 前端 |
| **多维表组件** | 3天 | 可筛选、排序、分页的数据表格 | 前端 |

**里程碑**：超级管理员可以管理所有用户，商家/服务商可以管理员工

#### Week 5-6: 营销活动与任务系统

| 任务 | 工作量 | 交付物 | 负责人 |
|------|-------|-------|-------|
| **邀请码系统** | 3天 | 邀请码生成、验证、使用记录 | 后端 |
| **营销活动模块** | 4天 | 活动 CRUD、状态流转、名额管理 | 后端 |
| **任务模块** | 4天 | 任务接取、提交、审核流程 | 后端 |
| **活动创建页面** | 3天 | 活动表单、富文本编辑、平台选择 | 前端 |
| **任务大厅** | 2天 | 达人浏览任务、接任务流程 | 前端 |
| **任务审核页面** | 3天 | 服务商审核任务、查看截图 | 前端 |

**里程碑**：商家可以创建活动，达人可以接任务，服务商可以审核任务

#### Week 7-8: 积分与财务系统

| 任务 | 工作量 | 交付物 | 负责人 |
|------|-------|-------|-------|
| **积分账户模块** | 3天 | 账户管理、余额冻结/解冻 | 后端 |
| **积分流水模块** | 3天 | 交易记录、资金流水、对账 | 后端 |
| **充值功能** | 2天 | 充值接口、支付回调（模拟） | 后端 |
| **提现功能** | 3天 | 提现申请、审核、转账（模拟） | 后端 |
| **财务数据面板** | 2天 | 余额展示、收支统计、交易记录 | 前端 |
| **充值提现页面** | 2天 | 充值表单、提现申请、提现记录 | 前端 |

**里程碑**：完整的财务流转，商家充值、达人提现功能完成

### 第二阶段：增强功能（6周）

#### Week 9-10: 通知与优化

| 任务 | 工作量 | 说明 |
|------|-------|------|
| **站内通知系统** | 5天 | 任务状态变更通知、系统通知 |
| **邮件通知** | 3天 | 提现审核结果、重要事件邮件 |
| **性能优化** | 4天 | 数据库查询优化、接口性能优化 |
| **前端性能优化** | 2天 | 代码分割、懒加载、缓存策略 |

#### Week 11-12: 数据统计与报表

| 任务 | 工作量 | 说明 |
|------|-------|------|
| **商家数据看板** | 4天 | 活动统计、费用统计、达人分析 |
| **服务商数据看板** | 3天 | 审核统计、收入统计、达人增长 |
| **达人数据看板** | 3天 | 收入统计、任务完成率、等级提升 |
| **超级管理员报表** | 4天 | 平台总览、财务总览、用户增长 |

#### Week 13-14: 安全与合规

| 任务 | 工作量 | 说明 |
|------|-------|------|
| **操作审计日志** | 3天 | 所有敏感操作记录、审计日志查询 |
| **财务对账系统** | 4天 | 自动对账、异常检测、对账报表 |
| **数据备份与恢复** | 2天 | 自动备份脚本、灾难恢复演练 |
| **安全加固** | 3天 | HTTPS 强制、CSRF 防护、文件上传限制 |

### 第三阶段：扩展与优化（4周）

#### Week 15-16: 移动端适配

| 任务 | 工作量 | 说明 |
|------|-------|------|
| **响应式设计优化** | 4天 | 适配手机、平板屏幕 |
| **移动端交互优化** | 3天 | 触摸操作、手势支持 |
| **移动端性能优化** | 3天 | 图片压缩、资源懒加载 |

#### Week 17-18: 压力测试与上线准备

| 任务 | 工作量 | 说明 |
|------|-------|------|
| **压力测试** | 3天 | 并发测试、性能瓶颈分析 |
| **Bug 修复** | 4天 | 测试发现的问题修复 |
| **文档编写** | 3天 | API 文档、部署文档、用户手册 |
| **上线准备** | 2天 | 域名备案、服务器配置、监控告警 |

---

## 📞 下一步

**请完成以下操作**：

1. ✅ 通读整个 PRD 文档
2. 📝 在标记 `📝` 或 `⚠️` 的地方填写你的答案
3. ✅ 在"待确认清单"中逐项确认
4. 💬 提出任何需要修改或补充的部分
5. ✅ 确认后，我将基于这个 PRD 开始开发

**预计开发时间**：确认 PRD 后，7-8 天完成 MVP

**有问题随时问我！**
