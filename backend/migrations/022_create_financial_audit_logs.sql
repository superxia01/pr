-- ============================================
-- 财务审计日志表
-- 用途：记录所有财务关键操作，确保可追溯
-- 创建时间:2026-02-12
-- ============================================

-- 创建审计日志表
CREATE TABLE financial_audit_logs (
    -- 主键
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- 用户ID（操作人）
    user_id VARCHAR(255) NOT NULL,

    -- 操作类型（区分不同财务操作）
    action VARCHAR(100) NOT NULL,
    -- 可能值：
    -- 提现相关：
    --   'WITHDRAWAL_APPROVE' - 审核通过
    --   'WITHDRAWAL_REJECT' - 审核拒绝
    -- 充值相关：
    --   'CREDIT_RECHARGE' - 用户充值
    -- 任务相关：
    --   'CAMPAIGN_PUBLISH' - 发布营销活动（冻结积分）
    --   'CAMPAIGN_TASK_CREATE' - 创建任务名额
    -- 系统相关：
    --   'SYSTEM_ADJUST' - 系统手动调整

    -- 资源类型（区分不同的资源）
    resource_type VARCHAR(50) NOT NULL,
    -- 可能值：
    -- 'WITHDRAWAL_REQUEST' - 提现申请
    -- 'PAYMENT_ORDER' - 支付订单
    -- 'CREDIT_ACCOUNT' - 积分账户
    -- 'CASH_ACCOUNT' - 现金账户
    -- 'SYSTEM_ACCOUNT' - 系统账户
    -- 'CAMPAIGN' - 营销活动

    -- 资源ID（对应资源的唯一标识）
    resource_id VARCHAR(255) NOT NULL,

    -- 变更记录（JSON格式，记录 before/after 状态）
    changes JSONB DEFAULT '{}',

    -- IP 地址（用于安全追踪）
    ip_address INET,

    -- User Agent（浏览器信息）
    user_agent TEXT,

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建索引（优化查询）
CREATE INDEX idx_audit_logs_user_action ON financial_audit_logs(user_id, action);
COMMENT ON idx_audit_logs_user_action IS '按用户和操作类型查询审计日志';

CREATE INDEX idx_audit_logs_resource ON financial_audit_logs(resource_type, resource_id);
COMMENT ON idx_audit_logs_resource IS '按资源类型和资源ID查询审计日志';

CREATE INDEX idx_audit_logs_created_at ON financial_audit_logs(created_at DESC);
COMMENT ON idx_audit_logs_created_at IS '按时间倒序查询审计日志';

-- ============================================
-- 操作类型枚举说明
-- ============================================
-- 提现审核：
-- WITHDRAWAL_APPROVE - 提现申请审核通过
--   changes: {"status": {"from": "PENDING", "to": "COMPLETED", "amount": 10000}

-- 充值记录：
-- CREDIT_RECHARGE - changes: {"before": 0, "after": 10000}

-- 任务发布：
-- CAMPAIGN_PUBLISH - changes: {"frozen_amount": 100000}

-- 系统调整：
-- SYSTEM_ADJUST - changes: {"adjustment": "冻结积分转现金"}
-- ============================================