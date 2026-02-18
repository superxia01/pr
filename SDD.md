# PR Business System - Software Design Document

**版本**: 1.0
**最后更新**: 2026-02-17
**系统名称**: PR Business System (PR业务系统)
**文档类型**: 软件设计文档 (SDD)

---

## 目录

1. [系统概述](#1-系统概述)
2. [角色定义](#2-角色定义)
3. [权限系统](#3-权限系统)
4. [功能模块](#4-功能模块)
5. [积分系统](#5-积分系统)
6. [活动和任务系统](#6-活动和任务系统)
7. [页面路由](#7-页面路由)
8. [API端点](#8-api端点)
9. [数据模型](#9-数据模型)
10. [技术架构](#10-技术架构)

---

## 1. 系统概述

### 1.1 系统简介

PR Business System是一个基于角色的多级营销任务管理平台，连接商家、服务商和内容创作者（达人），实现任务发布、执行、审核和结算的完整闭环。

### 1.2 核心功能

- **用户管理**: 多角色用户体系，支持7种用户角色
- **组织管理**: 服务商管理商家，组织结构清晰
- **活动管理**: 营销活动创建、审核、发布、关闭全流程
- **任务管理**: 任务接取、提交、审核、结算完整生命周期
- **积分系统**: 双余额设计（可用余额+冻结余额），完整的交易流水
- **财务管理**: 充值、提现、佣金分配、财务审计
- **邀请系统**: 固定邀请码、任务邀请码、多级返佣

### 1.3 业务流程图

```
商家/服务商创建活动 → 提交审核 → 审核通过 → 活动发布（冻结积分）
                                                    ↓
达人查看任务大厅 ← 活动开放中 ← 积分冻结成功
    ↓
达人接取任务（ASSIGNED）
    ↓
达人执行任务
    ↓
达人提交作品（SUBMITTED）
    ↓
服务商/商家审核
    ↓
    ├─ 通过 → 结算完成（佣金分配）
    └─ 拒绝 → 重新开放任务
```

---

## 2. 角色定义

系统共定义了**7种角色**，分为3类：

### 2.1 管理员角色（3种）

| 角色代码 | 显示名称 | 短名称 | 英文描述 | 特权 |
|---------|---------|--------|---------|------|
| `SUPER_ADMIN` | 超级管理员 | 超管 | Super Administrator | 系统最高权限，拥有所有权限 |
| `SERVICE_PROVIDER_ADMIN` | 服务商管理员 | 服务商管理员 | Service Provider Administrator | 管理服务商和员工，默认拥有所有权限 |
| `MERCHANT_ADMIN` | 商家管理员 | 商家管理员 | Merchant Administrator | 管理商家和员工，默认拥有所有权限 |

**特权说明**:
- 管理员角色默认拥有所有权限，无需单独配置
- 可以查看本组织所有数据
- 可以管理本组织员工

### 2.2 员工角色（2种）

| 角色代码 | 显示名称 | 短名称 | 英文描述 | 权限特点 |
|---------|---------|--------|---------|----------|
| `SERVICE_PROVIDER_STAFF` | 服务商员工 | 服务商员工 | Service Provider Staff | 需要从权限表细粒度控制 |
| `MERCHANT_STAFF` | 商家员工 | 商家员工 | Merchant Staff | 需要从权限表细粒度控制 |

**权限特点**:
- 员工角色没有默认权限
- 需要管理员授予具体权限（15个权限可选）
- 权限存储在 `provider_staff_permissions` 或 `merchant_staff_permissions` 表

### 2.3 业务角色（2种）

| 角色代码 | 显示名称 | 短名称 | 英文描述 | 核心功能 |
|---------|---------|--------|---------|----------|
| `CREATOR` | 达人 | 达人 | Creator | 接取任务、提交作品、获得佣金 |
| `BASIC_USER` | 基础用户 | 基础用户 | Basic User | 新注册用户默认角色 |

### 2.4 角色层级关系

```
SUPER_ADMIN (系统级)
    ├─ SERVICE_PROVIDER_ADMIN (组织级)
    │   └─ SERVICE_PROVIDER_STAFF (员工级)
    │       └─ MERCHANT_ADMIN (下级组织)
    │           └─ MERCHANT_STAFF (员工级)
    │               └─ CREATOR (业务执行)
    └─ BASIC_USER (默认角色)
```

---

## 3. 权限系统

### 3.1 权限分组

系统共定义了**15个细粒度权限**，分为5个权限组：

#### 3.1.1 达人管理权限组（3个）

| 权限代码 | 权限名称 | 描述 | 适用角色 |
|---------|---------|------|----------|
| `MANAGE_CREATORS` | 管理达人 | 查看和管理本组织所有达人 | SERVICE_PROVIDER_STAFF |
| `EDIT_CREATOR_INFO` | 编辑达人资料 | 编辑达人昵称、头像等基本信息 | SERVICE_PROVIDER_STAFF |
| `DELETE_CREATOR` | 删除达人 | 从本组织移除达人 | SERVICE_PROVIDER_STAFF |

#### 3.1.2 活动管理权限组（4个）

| 权限代码 | 权限名称 | 描述 | 适用角色 |
|---------|---------|------|----------|
| `PUBLISH_CAMPAIGN` | 发布活动 | 将活动状态改为 OPEN（公开接受任务） | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |
| `EDIT_CAMPAIGN_COMMISSION` | 编辑佣金分配 | 设置和修改活动的佣金分配方案 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |
| `EDIT_CAMPAIGN_INFO` | 编辑活动信息 | 编辑活动的标题、要求、平台、截止时间 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |
| `DELETE_CAMPAIGN` | 删除活动 | 删除活动 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |

#### 3.1.3 任务管理权限组（2个）

| 权限代码 | 权限名称 | 描述 | 适用角色 |
|---------|---------|------|----------|
| `REVIEW_TASK` | 审核任务 | 审核达人提交的任务作品 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |
| `VIEW_ALL_TASKS` | 查看所有任务 | 查看本组织所有活动的任务情况 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |

#### 3.1.4 财务管理权限组（3个）

| 权限代码 | 权限名称 | 描述 | 适用角色 |
|---------|---------|------|----------|
| `RECHARGE` | 充值 | 为积分账户充值 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |
| `APPROVE_WITHDRAWAL` | 审核提现 | 审批用户的提现申请 | 仅超级管理员 |
| `VIEW_FINANCIAL_REPORTS` | 查看财务报表 | 查看财务统计和报表 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |

#### 3.1.5 组织管理权限组（3个）

| 权限代码 | 权限名称 | 描述 | 适用角色 |
|---------|---------|------|----------|
| `MANAGE_MERCHANT` | 管理商家 | 创建、编辑、删除商家 | SERVICE_PROVIDER_STAFF |
| `MANAGE_STAFF` | 管理员工 | 添加、删除员工，管理员工权限 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |
| `CREATE_INVITATION_CODE` | 创建邀请码 | 生成各类角色的邀请码 | SERVICE_PROVIDER_STAFF, MERCHANT_STAFF |

### 3.2 权限检查机制

#### 3.2.1 权限检查函数

```go
// HasPermission 检查用户是否拥有指定权限
func HasPermission(db *gorm.DB, user *models.User, permissionCode string) bool

// HasAnyPermission 检查用户是否拥有任一指定权限
func HasAnyPermission(db *gorm.DB, user *models.User, permissionCodes ...string) bool

// HasAllPermissions 检查用户是否拥有所有指定权限
func HasAllPermissions(db *gorm.DB, user *models.User, permissionCodes ...string) bool

// RequirePermission 权限检查中间件
func RequirePermission(db *gorm.DB, permissionCode string) func(*gin.Context)

// GetCurrentUserPermissions 获取用户的所有权限列表
func GetCurrentUserPermissions(db *gorm.DB, user *models.User) []string
```

#### 3.2.2 权限检查流程

```
用户请求 → 获取当前用户 → 角色判断
                                    ↓
                    ┌───────────────┴───────────────┐
                    ↓                               ↓
              管理员                           员工
            (自动拥有所有权限)              (查询权限表)
                    ↓                               ↓
              继续处理 ← 权限匹配? ← → 有权限: 继续
                                        ↓
                                    无权限: 返回403
```

#### 3.2.3 权限使用示例

```go
// 任务审核权限检查
if utils.IsServiceProviderStaff(user) || utils.IsMerchantStaff(user) {
    if !utils.HasPermission(ctrl.db, user, constants.PermissionReviewTask) {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "无审核任务权限",
            "requiredPermission": constants.PermissionReviewTask,
        })
        return
    }
}
```

### 3.3 已实现权限检查的端点

| 端点 | 检查权限 | 实现位置 |
|------|---------|----------|
| `POST /api/v1/tasks/:id/audit` | REVIEW_TASK | task.go:307 |
| `PUT /api/v1/creators/:id` | EDIT_CREATOR_INFO | creator.go:205 |

---

## 4. 功能模块

### 4.1 用户认证模块

**功能**:
- 微信登录（OAuth2.0）
- 统一认证中心（auth-center）集成
- JWT Token管理

**API端点**:
- `GET /api/v1/auth/wechat/login` - 发起微信登录

**流程**:
1. 用户点击微信登录
2. 重定向到微信OAuth2.0授权页面
3. 微信回调返回code
4. 后端通过code换取用户信息
5. 生成JWT Token并返回

### 4.2 用户管理模块

**功能**:
- 用户信息查看和编辑
- 用户角色管理
- 用户列表查询（超级管理员）

**API端点**:
- `GET /api/v1/user/me` - 获取当前用户信息
- `GET /api/v1/users` - 获取用户列表（超级管理员）

**相关页面**:
- UserManagement.tsx - 用户管理后台

### 4.3 组织管理模块

#### 4.3.1 服务商管理

**功能**:
- 服务商注册和信息管理
- 服务商员工管理
- 服务商权限配置

**API端点**:
- `POST /api/v1/service-providers` - 创建服务商
- `GET /api/v1/service-providers` - 获取服务商列表
- `GET /api/v1/service-providers/:id` - 获取服务商详情
- `PUT /api/v1/service-providers/:id` - 更新服务商信息
- `DELETE /api/v1/service-providers/:id` - 删除服务商
- `GET /api/v1/service-provider/me` - 获取我的服务商信息
- `POST /api/v1/service-providers/:id/staff` - 添加员工
- `GET /api/v1/service-providers/:id/staff` - 获取员工列表
- `PUT /api/v1/service-providers/:id/staff/:staff_id/permissions` - 更新员工权限
- `DELETE /api/v1/service-providers/:id/staff/:staff_id` - 删除员工
- `GET /api/v1/service-providers/permissions` - 获取权限配置

**相关页面**:
- ServiceProviders.tsx - 服务商列表
- ServiceProviderDetail.tsx - 服务商详情
- ServiceProviderInfo.tsx - 服务商信息管理

#### 4.3.2 商家管理

**功能**:
- 商家注册和信息管理
- 商家员工管理
- 商家权限配置

**API端点**:
- `POST /api/v1/merchants` - 创建商家
- `GET /api/v1/merchants` - 获取商家列表
- `GET /api/v1/merchants/:id` - 获取商家详情
- `PUT /api/v1/merchants/:id` - 更新商家信息
- `DELETE /api/v1/merchants/:id` - 删除商家
- `GET /api/v1/merchant/me` - 获取我的商家信息
- `POST /api/v1/merchants/:id/staff` - 添加员工
- `GET /api/v1/merchants/:id/staff` - 获取员工列表
- `PUT /api/v1/merchants/:id/staff/:staff_id/permissions` - 更新员工权限
- `DELETE /api/v1/merchants/:id/staff/:staff_id` - 删除员工
- `GET /api/v1/merchants/permissions` - 获取权限配置

**相关页面**:
- Merchants.tsx - 商家列表
- MerchantInfo.tsx - 商家信息管理

### 4.4 达人管理模块

**功能**:
- 达人资料查看和编辑
- 达人等级管理（UGC、KOC、INF、KOL）
- 达人邀请关系管理
- 达人等级统计

**API端点**:
- `GET /api/v1/creators` - 获取达人列表
- `GET /api/v1/creators/stats/level` - 获取达人等级统计
- `GET /api/v1/creators/:id` - 获取达人详情
- `PUT /api/v1/creators/:id` - 更新达人信息
- `GET /api/v1/creators/:id/inviter` - 获取达人邀请关系
- `POST /api/v1/creators/:id/break-relationship` - 断开邀请关系
- `GET /api/v1/creator/me` - 获取我的达人资料
- `PUT /api/v1/creator/me` - 更新我的达人资料

**相关页面**:
- CreatorProfile.tsx - 达人资料页面

**达人等级定义**:
- `UGC` - 基础内容创作者
- `KOC` - 关键意见消费者
- `INF` - 普通达人
- `KOL` - 关键意见领袖

### 4.5 邀请码管理模块

**功能**:
- 固定邀请码生成和管理
- 邀请码使用记录
- 邀请码禁用

**邀请码类型**:
- `SERVICE_PROVIDER_STAFF` - 服务商员工邀请码
- `SERVICE_PROVIDER_ADMIN` - 服务商管理员邀请码
- `MERCHANT_STAFF` - 商家员工邀请码
- `MERCHANT_ADMIN` - 商家管理员邀请码
- `CREATOR` - 达人邀请码

**API端点**:
- `GET /api/v1/invitations/fixed-codes` - 获取我的固定邀请码
- `POST /api/v1/invitations/use` - 使用邀请码
- `GET /api/v1/invitations` - 获取邀请码列表
- `GET /api/v1/invitations/my` - 获取我的邀请列表
- `GET /api/v1/invitations/:code` - 获取邀请码详情
- `POST /api/v1/invitations/:code/disable` - 禁用邀请码

**相关页面**:
- InvitationCodes.tsx - 邀请码管理

### 4.6 任务邀请模块

**功能**:
- 任务邀请码生成
- 任务邀请码验证
- 任务邀请使用

**API端点**:
- `POST /api/v1/task-invitations/generate` - 生成任务邀请码
- `POST /api/v1/task-invitations/use` - 使用任务邀请码
- `GET /api/v1/task-invitations/validate/:code` - 验证任务邀请码

**相关页面**:
- TaskInvite.tsx - 任务邀请页面

---

## 5. 积分系统

### 5.1 积分账户模型

**核心字段**:
```go
type CreditAccount struct {
    ID            uuid.UUID  // 主键
    OwnerID       uuid.UUID  // 所有者ID
    OwnerType     OwnerType  // 所有者类型（ORG_MERCHANT/ORG_PROVIDER/USER_PERSONAL）
    Balance       int        // 可用余额
    FrozenBalance int        // 冻结余额
    CreatedAt     time.Time  // 创建时间
    UpdatedAt     time.Time  // 更新时间
}
```

**账户类型**:
- `ORG_MERCHANT` - 商家组织账户（用于充值、发布活动冻结、任务结算扣除）
- `ORG_PROVIDER` - 服务商组织账户（用于分成收入）
- `USER_PERSONAL` - 个人用户账户（达人收入、员工返佣）

### 5.2 双余额设计

**可用余额（Balance）**:
- 可自由使用的积分
- 用于发布活动、支付任务等

**冻结余额（FrozenBalance）**:
- 被特定用途占用的积分
- 包括：活动冻结、提现冻结

**总积分** = Balance + FrozenBalance

### 5.3 积分交易类型

| 交易代码 | 交易名称 | 描述 | 余额影响 |
|---------|---------|------|----------|
| `RECHARGE` | 充值 | 商家充值入账 | +Balance |
| `TASK_INCOME` | 任务收入 | 达人完成任务收入 | +Balance |
| `TASK_SUBMIT` | 任务提交 | 达人提交任务 | 冻结/解冻 |
| `STAFF_REFERRAL` | 员工返佣 | 员工推荐返佣 | +Balance |
| `PROVIDER_INCOME` | 服务商收入 | 服务商分成收入 | +Balance |
| `TASK_PUBLISH` | 发布任务 | 商家发布活动 | -Balance, +FrozenBalance |
| `TASK_ACCEPT` | 接任务 | 达人接取任务 | 记录 |
| `TASK_REJECT` | 审核拒绝 | 任务审核拒绝 | +FrozenBalance, -Balance |
| `TASK_ESCALATE` | 超时拒绝 | 任务超时拒绝 | +FrozenBalance, -Balance |
| `TASK_REFUND` | 任务退款 | 任务退款 | +Balance |
| `WITHDRAW` | 提现 | 用户提现成功 | -Balance |
| `WITHDRAW_FREEZE` | 提现冻结 | 提现申请冻结 | -Balance, +FrozenBalance |
| `WITHDRAW_REFUND` | 提现拒绝 | 提现申请拒绝 | +FrozenBalance, -Balance |
| `BONUS_GIFT` | 系统赠送 | 系统赠送积分 | +Balance |
| `CAMPAIGN_FREEZE` | 活动冻结 | 活动发布冻结 | -Balance, +FrozenBalance |
| `CAMPAIGN_REFUND` | 活动退还 | 活动关闭退还 | +FrozenBalance, -Balance |

### 5.4 充值流程

**流程图**:
```
商家创建充值订单
    ↓
填写充值信息（金额、支付方式、凭证）
    ↓
提交订单（状态：pending）
    ↓
超级管理员审核
    ├─ 通过：approved → completed
    │   ↓
    │   积分入账（+Balance）
    │   记录RECHARGE交易
    └─ 拒绝：rejected
        ↓
        记录拒绝原因
```

**API端点**:
- `POST /api/v1/recharge-orders` - 创建充值订单
- `GET /api/v1/recharge-orders` - 获取充值订单列表
- `POST /api/v1/recharge-orders/:id/audit` - 审核充值订单

**相关页面**:
- Recharge.tsx - 充值页面

### 5.5 提现流程

**流程图**:
```
用户创建提现申请
    ↓
填写提现信息（积分、金额、账户信息）
    ↓
系统冻结积分（-Balance, +FrozenBalance）
    ↓
创建提现申请（状态：pending）
    ↓
超级管理员审核
    ├─ 通过：
    │   ↓
    │   扣除冻结积分（-FrozenBalance）
    │   扣除现金（现金账户）
    │   状态：completed
    │   记录WITHDRAW交易
    └─ 拒绝：
        ↓
        解冻积分（+Balance, -FrozenBalance）
        状态：rejected
        记录拒绝原因
```

**API端点**:
- `POST /api/v1/withdrawals/enhanced` - 创建提现申请
- `GET /api/v1/withdrawals/enhanced` - 获取提现申请列表
- `GET /api/v1/withdrawals/enhanced/:id` - 获取提现申请详情
- `POST /api/v1/withdrawals/enhanced/:id/approve` - 批准提现
- `POST /api/v1/withdrawals/enhanced/:id/reject` - 拒绝提现

**相关页面**:
- Withdrawal.tsx - 提现申请
- Withdrawals.tsx - 提现记录
- WithdrawalReview.tsx - 提现审核

### 5.6 积分账户API

**API端点**:
- `GET /api/v1/credit/accounts` - 获取积分账户列表
- `GET /api/v1/credit/balance` - 获取账户余额
- `GET /api/v1/credit/transactions` - 获取交易记录
- `POST /api/v1/credit/recharge` - 积分充值

**相关页面**:
- CreditTransactions.tsx - 交易记录

---

## 6. 活动和任务系统

### 6.1 活动（Campaign）生命周期

**活动状态**:
```
DRAFT（草稿）
    ↓ 商家创建
PENDING_APPROVAL（待审核）
    ↓ 服务商审核
OPEN（开放中）
    ↓ 达人完成任务
CLOSED（已关闭）
```

**状态说明**:
- `DRAFT` - 活动草稿，可编辑
- `PENDING_APPROVAL` - 待服务商审核
- `OPEN` - 活动已发布，达人可接任务
- `CLOSED` - 活动已结束，结算完成

**活动字段**:
```go
type Campaign struct {
    ID                    uuid.UUID      // 主键
    ProviderID            *uuid.UUID     // 服务商ID
    MerchantID            uuid.UUID      // 商家ID
    Title                 string         // 活动标题
    Requirements          string         // 任务要求
    Platforms             string         // 允许的平台
    TaskAmount            int            // 任务单价（积分）
    Quota                 int            // 任务名额
    CreatorAmount         *int           // 达人佣金
    StaffReferralAmount   *int           // 员工返佣
    ProviderAmount        *int           // 服务商分成
    CampaignAmount        int            // 活动总金额
    TaskDeadline          time.Time      // 任务截止时间
    SubmissionDeadline    time.Time      // 提交截止时间
    Status                CampaignStatus // 活动状态
    CreatedAt             time.Time      // 创建时间
    UpdatedAt             time.Time      // 更新时间
}
```

### 6.2 任务（Task）生命周期

**任务状态**:
```
OPEN（开放中）
    ↓ 达人接取
ASSIGNED（已分配）
    ↓ 达人提交
SUBMITTED（已提交）
    ↓ 服务商审核
    ├─ APPROVED（已通过） → 结算完成
    └─ REJECTED（已拒绝） → 重新开放
```

**状态说明**:
- `OPEN` - 任务开放，可接取
- `ASSIGNED` - 已被达人接取
- `SUBMITTED` - 达人已提交作品
- `APPROVED` - 审核通过，完成结算
- `REJECTED` - 审核拒绝，重新开放

**任务字段**:
```go
type Task struct {
    ID              uuid.UUID       // 主键
    CampaignID      uuid.UUID       // 活动ID
    CreatorID       *uuid.UUID      // 达人ID
    InviterID       *string         // 邀请人ID
    InviterType     *string         // 邀请人类型
    Platform        string          // 执行平台
    PlatformURL     string          // 作品链接
    Screenshots     string          // 截图
    Notes           string          // 备注
    Status          TaskStatus      // 任务状态
    AssignedAt      *time.Time      // 接取时间
    SubmittedAt     *time.Time      // 提交时间
    AuditedBy       *uuid.UUID      // 审核人ID
    AuditedAt       *time.Time      // 审核时间
    AuditNote       string          // 审核备注
    Version         int             // 版本号
    CreatedAt       time.Time       // 创建时间
    UpdatedAt       time.Time       // 更新时间
}
```

### 6.3 活动创建流程

**商家创建活动**:
```
商家创建活动（状态：DRAFT）
    ↓
填写活动信息
    ├─ 标题、要求、平台
    ├─ 任务单价、名额
    └─ 截止时间
    ↓
提交审核（状态：PENDING_APPROVAL）
    ↓
服务商审核
    ├─ 通过：发布活动（OPEN）
    │   ↓
    │   冻结商家积分
    │   创建任务名额
    └─ 拒绝：退回草稿（DRAFT）
```

**服务商创建活动**:
```
服务商创建活动
    ↓
填写活动信息 + 佣金分配
    ├─ 达人佣金
    ├─ 员工返佣
    └─ 服务商分成
    ↓
直接发布（状态：OPEN）
    ↓
冻结商家积分
创建任务名额
```

**佣金分配验证**:
```
任务单价 = 达人佣金 + 员工返佣 + 服务商分成
```

### 6.4 任务执行流程

**接取任务**:
```
前置检查：
├── 任务状态为OPEN
├── 活动状态为OPEN
├── 未过TaskDeadline
└── 达人未接取过该活动的任务

执行操作：
├── 使用行锁防止并发
├── 更新任务状态为ASSIGNED
├── 记录达人和平台信息
└── 记录邀请人信息
```

**提交任务**:
```
前置检查：
├── 任务状态为ASSIGNED
├── 提交者是任务所属达人
└── 未过SubmissionDeadline

执行操作：
├── 更新任务状态为SUBMITTED
├── 记录提交信息（URL、截图、备注）
└── 更新提交时间
```

**审核任务**:
```
审核通过：
├── 异步执行结算
├── 从商家冻结账户扣除活动总金额
├── 给达人账户增加收入
├── 给员工增加返佣（如有）
├── 给服务商增加分成（如有）
└── 记录交易流水（同一TransactionGroupID）

审核拒绝：
├── 释放任务状态回OPEN
├── 清空达人提交信息
└── 允许达人重新接取
```

### 6.5 活动关闭结算

**流程**:
```
活动关闭（OPEN → CLOSED）
    ↓
统计已完成任务数（APPROVED状态）
    ↓
计算未完成任务数
    ↓
退还金额 = 未完成任务数 × 任务单价
    ↓
从冻结余额转移到可用余额
    ↓
记录CAMPAIGN_REFUND交易
```

### 6.6 API端点

**活动管理**:
- `POST /api/v1/campaigns` - 创建活动
- `GET /api/v1/campaigns` - 获取活动列表
- `GET /api/v1/campaigns/:id` - 获取活动详情
- `POST /api/v1/campaigns/:id/approve` - 审核活动
- `PUT /api/v1/campaigns/:id` - 更新活动
- `DELETE /api/v1/campaigns/:id` - 删除活动
- `GET /api/v1/campaigns/my` - 获取我的活动

**任务管理**:
- `GET /api/v1/tasks` - 获取任务列表
- `GET /api/v1/tasks/:id` - 获取任务详情
- `GET /api/v1/tasks/hall` - 任务大厅
- `GET /api/v1/tasks/my` - 我的任务
- `GET /api/v1/tasks/pending-review` - 待审核任务
- `POST /api/v1/tasks/:id/accept` - 接取任务
- `POST /api/v1/tasks/:id/submit` - 提交任务
- `POST /api/v1/tasks/:id/audit` - 审核任务

**相关页面**:
- CreateCampaign.tsx - 创建活动
- CampaignApproval.tsx - 活动审核
- TaskHall.tsx - 任务大厅
- MyTasks.tsx - 我的任务

---

## 7. 页面路由

### 7.1 公开页面

| 路由 | 页面组件 | 功能 |
|------|---------|------|
| `/` | Navigate | 根路径重定向到 /dashboard |
| `/login` | Login | 登录页面 |
| `/invite/:code` | TaskInvite | 任务邀请码页面 |

### 7.2 受保护页面（需要登录）

#### 7.2.1 主页和控制台

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/dashboard` | Dashboard | 主控制台 | 所有角色 |

#### 7.2.2 用户管理

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/user-management` | UserManagement | 用户管理后台 | 超级管理员 |

#### 7.2.3 商家管理

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/merchant` | MerchantInfo | 商家信息管理 | 商家管理员 |
| `/merchants` | Merchants | 商家列表 | 服务商/超管 |

#### 7.2.4 服务商管理

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/service-provider` | ServiceProviderInfo | 服务商信息 | 服务商管理员 |
| `/service-providers` | ServiceProviders | 服务商列表 | 超管 |
| `/service-provider/:id` | ServiceProviderDetail | 服务商详情 | 服务商/超管 |

#### 7.2.5 达人管理

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/creator` | CreatorProfile | 达人资料 | 达人角色 |

#### 7.2.6 任务管理

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/task-hall` | TaskHall | 任务大厅 | 达人角色 |
| `/my-tasks` | MyTasks | 我的任务 | 达人角色 |

#### 7.2.7 活动管理

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/create-campaign` | CreateCampaign | 创建活动 | 商家/服务商 |
| `/campaign-approval` | CampaignApproval | 活动审核 | 服务商 |

#### 7.2.8 财务管理

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/recharge` | Recharge | 账户充值 | 商家/超管 |
| `/credit-transactions` | CreditTransactions | 交易记录 | 所有角色 |
| `/withdrawal` | Withdrawal | 提现申请 | 所有角色 |
| `/withdrawals` | Withdrawals | 提现记录 | 所有角色 |
| `/withdrawal-review` | WithdrawalReview | 提现审核 | 超管 |

#### 7.2.9 邀请管理

| 路由 | 页面组件 | 功能 | 权限 |
|------|---------|------|------|
| `/invitations` | InvitationCodes | 邀请码管理 | 管理员/员工 |

### 7.3 路由配置（App.tsx）

```typescript
// 受保护路由
<ProtectedRoute>
  <Route path="/dashboard" element={<Dashboard />} />
  <Route path="/invitations" element={<InvitationCodes />} />
  <Route path="/user-management" element={<UserManagement />} />
  <Route path="/merchant" element={<MerchantInfo />} />
  <Route path="/merchants" element={<Merchants />} />
  <Route path="/service-provider/:id" element={<ServiceProviderDetail />} />
  <Route path="/service-provider" element={<ServiceProviderInfo />} />
  <Route path="/service-providers" element={<ServiceProviders />} />
  <Route path="/creator" element={<CreatorProfile />} />
  <Route path="/task-hall" element={<TaskHall />} />
  <Route path="/my-tasks" element={<MyTasks />} />
  <Route path="/create-campaign" element={<CreateCampaign />} />
  <Route path="/campaign-approval" element={<CampaignApproval />} />
  <Route path="/recharge" element={<Recharge />} />
  <Route path="/credit-transactions" element={<CreditTransactions />} />
  <Route path="/withdrawal" element={<Withdrawal />} />
  <Route path="/withdrawals" element={<Withdrawals />} />
  <Route path="/withdrawal-review" element={<WithdrawalReview />} />
</ProtectedRoute>
```

---

## 8. API端点

### 8.1 API概览

系统共有**100+个API端点**，分为15个功能模块：

| 模块 | 端点数量 | 主要功能 |
|------|---------|----------|
| 认证模块 | 2 | 微信登录、用户信息 |
| 用户管理 | 2 | 用户列表、当前用户 |
| 邀请码管理 | 6 | 固定邀请码管理 |
| 商家管理 | 11 | 商家CRUD、员工管理 |
| 服务商管理 | 11 | 服务商CRUD、员工管理 |
| 达人管理 | 8 | 达人资料管理 |
| 活动管理 | 7 | 活动CRUD、审核 |
| 任务邀请 | 3 | 任务邀请码管理 |
| 任务管理 | 8 | 任务CRUD、审核 |
| 积分管理 | 4 | 积分账户、交易记录 |
| 充值订单 | 3 | 充值订单管理 |
| 提现管理 | 5 | 提现申请、审核 |
| 增强提现 | 5 | 带冻结机制的提现 |
| 现金账户 | 3 | 现金账户管理 |
| 系统账户 | 2 | 系统账户、财务汇总 |
| 财务审计 | 1 | 审计日志查询 |

### 8.2 详细API列表

#### 8.2.1 认证模块

| 方法 | 路径 | 功能 | 认证 |
|------|------|------|------|
| GET | `/api/v1/auth/wechat/login` | 发起微信登录 | ❌ |
| GET | `/api/v1/user/me` | 获取当前用户信息 | ✅ |

#### 8.2.2 用户管理

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | `/api/v1/users` | 获取用户列表 | 超级管理员 |

#### 8.2.3 邀请码管理

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/v1/invitations/fixed-codes` | 获取我的固定邀请码 |
| POST | `/api/v1/invitations/use` | 使用邀请码 |
| GET | `/api/v1/invitations` | 获取邀请码列表 |
| GET | `/api/v1/invitations/my` | 获取我的邀请列表 |
| GET | `/api/v1/invitations/:code` | 获取邀请码详情 |
| POST | `/api/v1/invitations/:code/disable` | 禁用邀请码 |

#### 8.2.4 商家管理

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| POST | `/api/v1/merchants` | 创建商家 | 商家权限 |
| GET | `/api/v1/merchants` | 获取商家列表 | 商家权限 |
| GET | `/api/v1/merchants/:id` | 获取商家详情 | 商家权限 |
| PUT | `/api/v1/merchants/:id` | 更新商家信息 | 商家权限 |
| DELETE | `/api/v1/merchants/:id` | 删除商家 | 商家权限 |
| GET | `/api/v1/merchant/me` | 获取我的商家信息 | 商家权限 |
| POST | `/api/v1/merchants/:id/staff` | 添加员工 | 商家管理员 |
| GET | `/api/v1/merchants/:id/staff` | 获取员工列表 | 商家管理员 |
| PUT | `/api/v1/merchants/:id/staff/:staff_id/permissions` | 更新员工权限 | 商家管理员 |
| DELETE | `/api/v1/merchants/:id/staff/:staff_id` | 删除员工 | 商家管理员 |
| GET | `/api/v1/merchants/permissions` | 获取权限配置 | 商家权限 |

#### 8.2.5 服务商管理

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| POST | `/api/v1/service-providers` | 创建服务商 | 服务商权限 |
| GET | `/api/v1/service-providers` | 获取服务商列表 | 服务商权限 |
| GET | `/api/v1/service-providers/:id` | 获取服务商详情 | 服务商权限 |
| PUT | `/api/v1/service-providers/:id` | 更新服务商信息 | 服务商权限 |
| DELETE | `/api/v1/service-providers/:id` | 删除服务商 | 服务商权限 |
| GET | `/api/v1/service-provider/me` | 获取我的服务商信息 | 服务商权限 |
| POST | `/api/v1/service-providers/:id/staff` | 添加员工 | 服务商管理员 |
| GET | `/api/v1/service-providers/:id/staff` | 获取员工列表 | 服务商管理员 |
| PUT | `/api/v1/service-providers/:id/staff/:staff_id/permissions` | 更新员工权限 | 服务商管理员 |
| DELETE | `/api/v1/service-providers/:id/staff/:staff_id` | 删除员工 | 服务商管理员 |
| GET | `/api/v1/service-providers/permissions` | 获取权限配置 | 服务商权限 |

#### 8.2.6 达人管理

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/v1/creators` | 获取达人列表 |
| GET | `/api/v1/creators/stats/level` | 获取达人等级统计 |
| GET | `/api/v1/creators/:id` | 获取达人详情 |
| PUT | `/api/v1/creators/:id` | 更新达人信息 |
| GET | `/api/v1/creators/:id/inviter` | 获取达人邀请关系 |
| POST | `/api/v1/creators/:id/break-relationship` | 断开邀请关系 |
| GET | `/api/v1/creator/me` | 获取我的达人资料 |
| PUT | `/api/v1/creator/me` | 更新我的达人资料 |

#### 8.2.7 活动管理

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| POST | `/api/v1/campaigns` | 创建活动 | 活动创建权限 |
| GET | `/api/v1/campaigns` | 获取活动列表 | 认证即可 |
| GET | `/api/v1/campaigns/:id` | 获取活动详情 | 认证即可 |
| POST | `/api/v1/campaigns/:id/approve` | 审核活动 | 活动审核权限 |
| PUT | `/api/v1/campaigns/:id` | 更新活动 | 活动编辑权限 |
| DELETE | `/api/v1/campaigns/:id` | 删除活动 | 活动删除权限 |
| GET | `/api/v1/campaigns/my` | 获取我的活动 | 认证即可 |

#### 8.2.8 任务邀请

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | `/api/v1/task-invitations/generate` | 生成任务邀请码 |
| POST | `/api/v1/task-invitations/use` | 使用任务邀请码 |
| GET | `/api/v1/task-invitations/validate/:code` | 验证任务邀请码 |

#### 8.2.9 任务管理

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | `/api/v1/tasks` | 获取任务列表 | 认证即可 |
| GET | `/api/v1/tasks/:id` | 获取任务详情 | 认证即可 |
| GET | `/api/v1/tasks/hall` | 任务大厅 | 达人角色 |
| GET | `/api/v1/tasks/my` | 我的任务 | 认证即可 |
| GET | `/api/v1/tasks/pending-review` | 待审核任务 | 审核权限 |
| POST | `/api/v1/tasks/:id/accept` | 接取任务 | 达人角色 |
| POST | `/api/v1/tasks/:id/submit` | 提交任务 | 达人角色 |
| POST | `/api/v1/tasks/:id/audit` | 审核任务 | REVIEW_TASK |

#### 8.2.10 积分管理

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/v1/credit/accounts` | 获取积分账户列表 |
| GET | `/api/v1/credit/balance` | 获取账户余额 |
| GET | `/api/v1/credit/transactions` | 获取交易记录 |
| POST | `/api/v1/credit/recharge` | 积分充值 |

#### 8.2.11 充值订单

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | `/api/v1/recharge-orders` | 创建充值订单 |
| GET | `/api/v1/recharge-orders` | 获取充值订单列表 |
| POST | `/api/v1/recharge-orders/:id/audit` | 审核充值订单 |

#### 8.2.12 提现管理

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | `/api/v1/withdrawals` | 创建提现申请 |
| GET | `/api/v1/withdrawals` | 获取提现申请列表 |
| GET | `/api/v1/withdrawals/:id` | 获取提现申请详情 |
| POST | `/api/v1/withdrawals/:id/audit` | 审核提现申请 |
| POST | `/api/v1/withdrawals/:id/process` | 处理提现 |

#### 8.2.13 增强提现（带冻结）

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | `/api/v1/withdrawals/enhanced` | 创建提现请求 |
| GET | `/api/v1/withdrawals/enhanced` | 获取提现请求列表 |
| GET | `/api/v1/withdrawals/enhanced/:id` | 获取提现请求详情 |
| POST | `/api/v1/withdrawals/enhanced/:id/approve` | 批准提现请求 |
| POST | `/api/v1/withdrawals/enhanced/:id/reject` | 拒绝提现请求 |

#### 8.2.14 现金账户

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/v1/cash-accounts` | 获取现金账户列表 |
| POST | `/api/v1/cash-accounts` | 创建现金账户 |
| POST | `/api/v1/cash-accounts/:id/balance` | 更新账户余额 |

#### 8.2.15 系统账户

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/v1/system-accounts` | 获取系统账户列表 |
| GET | `/api/v1/system-accounts/summary` | 获取财务汇总 |

#### 8.2.16 财务审计

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/v1/financial-audit-logs` | 获取财务审计日志 |

---

## 9. 数据模型

### 9.1 核心数据表

#### 9.1.1 用户相关

**users** - 用户表
```go
- ID                  uuid.UUID   (主键)
- AuthCenterUserID    string      (认证中心用户ID)
- UnionID             string      (微信UnionID)
- OpenID              string      (微信OpenID)
- Nickname            string      (昵称)
- AvatarURL           string      (头像)
- Roles               Roles       (角色数组)
- Status              string      (状态)
- CreatedAt           time.Time
- UpdatedAt           time.Time
```

**service_providers** - 服务商表
```go
- ID                  uuid.UUID   (主键)
- AdminID             uuid.UUID   (管理员ID)
- Name                string      (名称)
- Description         string      (描述)
- Status              string      (状态)
- CreatedAt           time.Time
- UpdatedAt           time.Time
```

**service_provider_staff** - 服务商员工表
```go
- ID                  uuid.UUID   (主键)
- ProviderID          uuid.UUID   (服务商ID)
- UserID              uuid.UUID   (用户ID)
- Status              string      (状态)
- CreatedAt           time.Time
- UpdatedAt           time.Time
```

**merchants** - 商家表
```go
- ID                  uuid.UUID   (主键)
- ProviderID          *uuid.UUID  (服务商ID)
- AdminID             uuid.UUID   (管理员ID)
- Name                string      (名称)
- Description         string      (描述)
- Status              string      (状态)
- CreatedAt           time.Time
- UpdatedAt           time.Time
```

**merchant_staff** - 商家员工表
```go
- ID                  uuid.UUID   (主键)
- MerchantID          uuid.UUID   (商家ID)
- UserID              uuid.UUID   (用户ID)
- Status              string      (状态)
- CreatedAt           time.Time
- UpdatedAt           time.Time
```

**creators** - 达人表
```go
- ID                  uuid.UUID   (主键)
- UserID              uuid.UUID   (用户ID)
- InviterID           *string     (邀请人ID)
- InviterType         *string     (邀请人类型)
- Level               string      (达人等级)
- FollowersCount      int         (粉丝数)
- WechatOpenID        string      (微信OpenID)
- WechatNickname      string      (微信昵称)
- WechatAvatar        string      (微信头像)
- Status              string      (状态)
- IsPrimary           bool        (是否主档案)
- CreatedAt           time.Time
- UpdatedAt           time.Time
```

#### 9.1.2 活动和任务

**campaigns** - 活动表
```go
- ID                    uuid.UUID      (主键)
- ProviderID            *uuid.UUID     (服务商ID)
- MerchantID            uuid.UUID      (商家ID)
- Title                 string         (标题)
- Requirements          string         (要求)
- Platforms             string         (平台)
- TaskAmount            int            (任务单价)
- Quota                 int            (名额)
- CreatorAmount         *int           (达人佣金)
- StaffReferralAmount   *int           (员工返佣)
- ProviderAmount        *int           (服务商分成)
- CampaignAmount        int            (活动总金额)
- TaskDeadline          time.Time      (任务截止)
- SubmissionDeadline    time.Time      (提交截止)
- Status                string         (状态)
- CreatedAt             time.Time
- UpdatedAt             time.Time
```

**tasks** - 任务表
```go
- ID              uuid.UUID       (主键)
- CampaignID      uuid.UUID       (活动ID)
- CreatorID       *uuid.UUID      (达人ID)
- InviterID       *string         (邀请人ID)
- InviterType     *string         (邀请人类型)
- Platform        string          (执行平台)
- PlatformURL     string          (作品链接)
- Screenshots     string          (截图)
- Notes           string          (备注)
- Status          string          (状态)
- AssignedAt      *time.Time      (接取时间)
- SubmittedAt     *time.Time      (提交时间)
- AuditedBy       *uuid.UUID      (审核人ID)
- AuditedAt       *time.Time      (审核时间)
- AuditNote       string          (审核备注)
- Version         int             (版本号)
- CreatedAt       time.Time
- UpdatedAt       time.Time
```

#### 9.1.3 积分和财务

**credit_accounts** - 积分账户表
```go
- ID            uuid.UUID  (主键)
- OwnerID       uuid.UUID  (所有者ID)
- OwnerType     string     (所有者类型)
- Balance       int        (可用余额)
- FrozenBalance int        (冻结余额)
- CreatedAt     time.Time
- UpdatedAt     time.Time
```

**credit_transactions** - 积分交易表
```go
- ID                  uuid.UUID   (主键)
- AccountID           uuid.UUID   (账户ID)
- TransactionTypeID   string      (交易类型)
- Amount              int         (金额)
- BalanceBefore       int         (交易前余额)
- BalanceAfter        int         (交易后余额)
- TransactionGroupID  *string     (交易组ID)
- GroupSequence       int         (组内序号)
- RelatedCampaignID   *uuid.UUID  (关联活动)
- RelatedTaskID       *uuid.UUID  (关联任务)
- Description         string      (描述)
- CreatedAt           time.Time
```

**recharge_orders** - 充值订单表
```go
- ID              uuid.UUID   (主键)
- MerchantID      uuid.UUID   (商家ID)
- Amount          int         (充值积分)
- YuanAmount      float64     (充值金额)
- PaymentMethod   string      (支付方式)
- PaymentProof    string      (支付凭证)
- Status          string      (状态)
- AuditNote       string      (审核备注)
- AuditedBy       *uuid.UUID  (审核人)
- AuditedAt       *time.Time  (审核时间)
- CreatedAt       time.Time
- UpdatedAt       time.Time
```

**withdrawal_requests** - 提现申请表
```go
- ID              uuid.UUID   (主键)
- AccountID       uuid.UUID   (账户ID)
- Amount          int         (提现积分)
- YuanAmount      float64     (提现金额)
- Status          string      (状态)
- CashAccountType string      (现金账户类型)
- Description     string      (描述)
- RejectReason    *string     (拒绝原因)
- ReviewedBy      string      (审核人)
- ReviewedAt      *time.Time  (审核时间)
- CompletedAt     *time.Time  (完成时间)
- CreatedAt       time.Time
```

**cash_accounts** - 现金账户表
```go
- ID          uuid.UUID   (主键)
- Name        string      (账户名称)
- AccountType string      (账户类型)
- Balance     int         (余额，单位：分)
- Status      string      (状态)
- CreatedAt   time.Time
- UpdatedAt   time.Time
```

**financial_audit_logs** - 财务审计日志表
```go
- ID          uuid.UUID   (主键)
- UserID      string      (操作用户ID)
- Action      string      (操作类型)
- ResourceType string     (资源类型)
- ResourceID  string      (资源ID)
- Details     jsonb       (详细信息)
- IPAddress   string      (IP地址)
- UserAgent   string      (User-Agent)
- CreatedAt   time.Time
```

#### 9.1.4 邀请系统

**invitation_codes** - 邀请码表
```go
- ID              uuid.UUID   (主键)
- Code            string      (邀请码)
- GeneratorID     uuid.UUID   (生成者ID)
- GeneratorType   string      (生成者类型)
- InvitedRole     string      (邀请角色)
- Status          string      (状态)
- UsageLimit      int         (使用次数限制)
- UsedCount       int         (已使用次数)
- ExpiresAt       *time.Time  (过期时间)
- CreatedAt       time.Time
- UpdatedAt       time.Time
```

**invitation_relationships** - 邀请关系表
```go
- ID          uuid.UUID   (主键)
- InviterID   uuid.UUID   (邀请人ID)
- InviteeID   uuid.UUID   (被邀请人ID)
- Role        string      (角色)
- Status      string      (状态)
- CreatedAt   time.Time
```

**task_invitations** - 任务邀请表
```go
- ID              uuid.UUID   (主键)
- Code            string      (邀请码)
- CreatorID       uuid.UUID   (创建者ID)
- TaskID          *uuid.UUID  (任务ID)
- CampaignID      *uuid.UUID  (活动ID)
- Status          string      (状态)
- ExpiresAt       *time.Time  (过期时间)
- CreatedAt       time.Time
- UpdatedAt       time.Time
```

#### 9.1.5 权限系统

**service_provider_staff_permissions** - 服务商员工权限表
```go
- ID              uuid.UUID   (主键)
- StaffID         uuid.UUID   (员工ID)
- PermissionCode  string      (权限代码)
- GrantedAt       time.Time  (授权时间)
- GrantedBy       uuid.UUID   (授权人)
```

**merchant_staff_permissions** - 商家员工权限表
```go
- ID              uuid.UUID   (主键)
- StaffID         uuid.UUID   (员工ID)
- PermissionCode  string      (权限代码)
- GrantedAt       time.Time  (授权时间)
- GrantedBy       uuid.UUID   (授权人)
```

---

## 10. 技术架构

### 10.1 技术栈

#### 后端技术栈

| 组件 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 语言 | Go | 1.25+ | 后端开发语言 |
| 框架 | Gin | - | Web框架 |
| ORM | GORM | - | 数据库ORM |
| 数据库 | PostgreSQL | 14+ | 关系型数据库 |
| 认证 | JWT | - | Token认证 |
| 认证中心 | auth-center | - | 统一认证 |

#### 前端技术栈

| 组件 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 语言 | TypeScript | 5+ | 前端开发语言 |
| 框架 | React | 18+ | UI框架 |
| 构建 | Vite | 5+ | 构建工具 |
| 路由 | React Router | 6+ | 路由管理 |
| 样式 | Tailwind CSS | 3+ | CSS框架 |
| 组件 | Shadcn/UI | - | UI组件库 |
| 状态 | Context API | - | 状态管理 |

### 10.2 后端架构

#### 分层架构

```
┌─────────────────────────────────────────┐
│           Routes Layer                  │  路由层
│   (routes/routes.go)                    │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────┴───────────────────────┐
│        Controllers Layer                │  控制器层
│   (controllers/*.go)                     │
│   - HTTP请求处理                         │
│   - 参数验证                             │
│   - 调用Service层                        │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────┴───────────────────────┐
│         Services Layer                  │  业务逻辑层
│   (services/*.go)                        │
│   - 核心业务逻辑                         │
│   - 事务管理                             │
│   - 跨模型协调                           │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────┴───────────────────────┐
│          Models Layer                   │  数据模型层
│   (models/*.go)                          │
│   - 数据库表结构                         │
│   - GORM模型定义                         │
│   - 数据验证                             │
└─────────────────────────────────────────┘
```

#### 目录结构

```
backend/
├── cmd/                # 可执行文件
├── config/             # 配置文件
├── constants/          # 常量定义
├── controllers/        # 控制器（15个）
├── middlewares/        # 中间件（4个）
├── migrations/         # 数据库迁移（27个）
├── models/             # 数据模型（14个）
├── routes/             # 路由定义
├── services/           # 业务逻辑（8个）
└── utils/              # 工具函数
```

### 10.3 前端架构

#### 组件结构

```
src/
├── components/         # 通用组件
│   ├── ui/             # UI组件（shadcn/ui）
│   ├── Sidebar.tsx     # 侧边栏
│   ├── MobileNav.tsx   # 移动端导航
│   └── ...
├── contexts/           # Context API
│   ├── AuthContext.tsx # 认证上下文
│   └── ToastContext.tsx # 提示上下文
├── lib/                # 工具函数
│   ├── roles.ts        # 角色工具
│   └── api.ts          # API客户端
├── pages/              # 页面组件（20个）
├── types/              # TypeScript类型
├── App.tsx             # 根组件
└── main.tsx            # 入口文件
```

#### 状态管理

使用React Context API进行状态管理：

- **AuthContext** - 用户认证状态
  - 当前用户信息
  - 登录/登出方法
  - 角色检查方法

- **ToastContext** - 提示消息
  - 成功/错误/信息提示
  - 自动消失机制

### 10.4 数据流

#### 请求流程

```
用户操作
    ↓
前端页面
    ↓
API客户端（api.ts）
    ↓
HTTP请求（带JWT Token）
    ↓
后端中间件（认证）
    ↓
路由分发
    ↓
Controller（处理请求）
    ↓
Service（业务逻辑）
    ↓
Model（数据库操作）
    ↓
数据库
    ↓
返回结果
    ↓
前端更新UI
```

#### 认证流程

```
1. 用户登录
    ↓
2. 微信OAuth2.0授权
    ↓
3. 获取用户信息
    ↓
4. 生成JWT Token
    ↓
5. 前端保存Token
    ↓
6. 后续请求携带Token
    ↓
7. 后端验证Token
    ↓
8. 获取用户信息
    ↓
9. 执行业务逻辑
```

### 10.5 数据库设计

#### 设计原则

1. **UUID主键** - 所有表使用UUID作为主键
2. **软删除** - 使用DeletedAt字段实现软删除
3. **时间戳** - CreatedAt和UpdatedAt记录创建和更新时间
4. **索引优化** - 为常用查询字段添加索引
5. **外键约束** - 保证数据完整性

#### 迁移文件

共有27个迁移文件，按序号顺序执行：

- `001_init_schema.sql` - 初始化数据库结构
- `002_seed_data.sql` - 基础数据填充
- `003-017` - 功能迭代迁移
- `018-022` - 财务系统迁移
- `023-026` - 邀请关系迁移
- `027+` - 后续功能迁移

### 10.6 安全设计

#### 认证安全

1. **JWT Token** - 使用Token进行无状态认证
2. **Token过期** - Token有过期时间
3. **刷新机制** - 支持Token刷新
4. **HTTPS** - 生产环境强制HTTPS

#### 权限安全

1. **角色验证** - 每个请求验证用户角色
2. **权限检查** - 细粒度权限控制
3. **资源隔离** - 用户只能访问有权限的资源
4. **审计日志** - 记录所有敏感操作

#### 数据安全

1. **事务保护** - 关键操作使用数据库事务
2. **行锁机制** - 防止并发问题
3. **数据验证** - 严格的输入验证
4. **SQL注入防护** - 使用参数化查询

### 10.7 性能优化

#### 后端优化

1. **数据库索引** - 为常用查询添加索引
2. **查询优化** - 使用Preload预加载关联数据
3. **分页查询** - 大数据量使用分页
4. **并发控制** - 使用行锁防止并发冲突

#### 前端优化

1. **代码分割** - React懒加载
2. **缓存策略** - API响应缓存
3. **虚拟滚动** - 大列表虚拟滚动
4. **防抖节流** - 减少不必要的请求

---

## 附录

### A. 术语表

| 术语 | 解释 |
|------|------|
| PR | Public Relations，公共关系 |
| 达人/创作者 | Creator，内容创作者 |
| 商家 | Merchant，品牌方 |
| 服务商 | Service Provider，平台运营方 |
| 活动 | Campaign，营销活动 |
| 任务 | Task，具体执行任务 |
| 积分 | Credit，虚拟货币单位 |
| 冻结积分 | FrozenBalance，被占用的积分 |
| 可用积分 | Balance，可自由使用的积分 |
| 邀请码 | Invitation Code，邀请用户注册的码 |
| 佣金 | Commission，任务报酬 |
| 返佣 | Referral Commission，推荐奖励 |

### B. 状态码

**活动状态（CampaignStatus）**:
- `DRAFT` - 草稿
- `PENDING_APPROVAL` - 待审核
- `OPEN` - 开放中
- `CLOSED` - 已关闭

**任务状态（TaskStatus）**:
- `OPEN` - 开放中
- `ASSIGNED` - 已分配
- `SUBMITTED` - 已提交
- `APPROVED` - 已通过
- `REJECTED` - 已拒绝

**提现状态（WithdrawalStatus）**:
- `pending` - 待审核
- `processing` - 处理中
- `completed` - 已完成
- `rejected` - 已拒绝

**充值订单状态**:
- `pending` - 待审核
- `approved` - 已批准
- `completed` - 已完成
- `rejected` - 已拒绝

### C. 环境变量

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=pr_business
DB_USER=postgres
DB_PASSWORD=your_password

# JWT配置
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION=24h

# 认证中心配置
AUTH_CENTER_URL=https://auth.example.com
AUTH_CLIENT_ID=your_client_id
AUTH_CLIENT_SECRET=your_client_secret

# 微信配置
WECHAT_APP_ID=your_app_id
WECHAT_APP_SECRET=your_app_secret

# 服务配置
SERVER_PORT=8080
SERVER_MODE=release
```

### D. 部署架构

```
┌─────────────┐
│   Nginx     │  反向代理 + 静态文件
└──────┬──────┘
       │
┌──────┴──────┐
│  Go Backend │  API服务
│  (pr-api)   │
└──────┬──────┘
       │
┌──────┴──────┐
│ PostgreSQL  │  数据库
└─────────────┘
```

---

**文档版本**: 1.0
**最后更新**: 2026-02-17
**维护者**: Development Team
