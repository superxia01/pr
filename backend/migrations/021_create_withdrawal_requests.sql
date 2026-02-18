-- ============================================
-- 提现申请表
-- 用途：记录用户的提现申请流程
-- 创建时间:2026-02-12
-- ============================================

-- 创建提现申请表
CREATE TABLE withdrawal_requests (
    -- 主键
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- 关联积分账户（引用 credit_accounts.id）
    account_id UUID NOT NULL REFERENCES credit_accounts(id),

    -- 提现金额（积分）
    amount INT NOT NULL,

    -- 提现金额（元，保留两位小数）
    yuan_amount DECIMAL(10, 2) NOT NULL,

    -- 状态：PENDING, PROCESSING, COMPLETED, REJECTED
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',

    -- 现金账户类型（用于打款）
    cash_account_type VARCHAR(50),
    -- 可能值：'WECHAT', 'ALIPAY', 'BANK_TRANSFER',

    -- 描述和拒绝原因
    description TEXT,
    reject_reason TEXT,

    -- 审核信息
    reviewed_by VARCHAR(255),
    reviewed_at TIMESTAMP,
    completed_at TIMESTAMP,

    -- 申请时间
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_withdrawal_requests_account_status ON withdrawal_requests(account_id, status);
COMMENT ON idx_withdrawal_requests_account_status IS '查询用户的提现申请状态';

-- ============================================
-- 外键说明（积分账户已存在）
-- ============================================
-- 积分账户表结构：
-- credit_accounts (用户积分)
--   - id: UUID
--   - user_id: VARCHAR(255)
--   - balance: INT (可用余额)
--   - frozen_balance: INT (冻结余额) - 本次优化新增
--
-- withdrawal_requests.account_id 引用的是 credit_accounts.id
-- 外键约束在 ON DELETE SET NULL，允许级联删除（积分账户被删除时，提现申请记录保留或级联删除）
--
-- ============================================
-- 状态说明
-- ============================================
-- PENDING: 待审核（冻结状态，积分已冻结到 frozen_balance）
-- PROCESSING: 处理中（财务打款中）
-- COMPLETED: 已完成（积分已扣除，可打款）
-- REJECTED: 已拒绝（退还积分到可用余额）
--
-- 审核流程：
-- 1. 用户申请 → status=PENDING, frozen_balance+=amount
-- 2. 管理员审核通过 → status=COMPLETED, frozen_balance-=amount
-- 3. 管理员审核拒绝 → status=REJECTED, frozen_balance-=amount, balance+=amount
--
-- 打款流程：
-- 1. COMPLETED → 管理员从现金账户打款给用户
-- 2. 扣除 cash_account.balance (对应 yuan_amount)
-- 3. withdrawal_requests.completed_at = NOW()
-- ============================================
