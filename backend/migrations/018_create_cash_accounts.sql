-- ============================================
-- 现金账户表
-- 用途：记录真实的资金流入流出（微信、支付宝、银行转账等）
-- 创建时间: 2026-02-12
-- ============================================

-- 创建现金账户表
CREATE TABLE cash_accounts (
    -- 主键
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- 账户类型（区分不同渠道）
    account_type VARCHAR(50) NOT NULL,
    -- 可能的值:
    -- 'WECHAT'      - 微信商户账户
    -- 'ALIPAY'      - 支付宝商户账户
    -- 'BANK_TRANSFER' - 银行转账账户
    -- 'MARKETING'     - 营销赠送账户
    -- 'OPERATIONS'   - 运营成本账户

    -- 余额（单位：分）
    balance INT NOT NULL DEFAULT 0,

    -- 账户描述
    description TEXT,

    -- 是否激活（用于软删除）
    is_active BOOLEAN NOT NULL DEFAULT true,

    -- 元数据（记录银行信息、账号等）
    metadata JSONB DEFAULT '{}',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_cash_accounts_type_active ON cash_accounts(account_type, is_active);
COMMENT ON idx_cash_accounts_type_active IS '按账户类型和激活状态索引，优化查询';

-- ============================================
-- 外键约束说明（将来需要）
-- ============================================
-- 提现表会引用 account_id，需要在外键建立后添加
-- ALTER TABLE withdrawal_requests ADD CONSTRAINT fk_withdrawal_account_id
--     FOREIGN KEY (account_id) REFERENCES cash_accounts(id) ON DELETE SET NULL;

COMMENT ON cash_accounts IS '现金账户：记录平台真实的资金流向';
