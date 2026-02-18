-- ============================================
-- 系统账户表
-- 用途：托管平台资产（任务托管、平台收益）
-- 创建时间:2026-02-12
-- ============================================

-- 创建系统账户表
CREATE TABLE system_accounts (
    -- 主键
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- 账户类型（区分不同托管）
    account_type VARCHAR(50) NOT NULL,
    -- 可能的值：
    -- 'TICKET_ESCROW'       - 票务托管（积分）
    -- 'TASK_ESCROW'         - 任务托管（积分）
    -- 'PLATFORM_REVENUE'    - 平台收益（积分）

    -- 余额（单位：积分）
    balance INT NOT NULL DEFAULT 0,

    -- 账户描述
    description TEXT,

    -- 是否激活（用于软删除）
    is_active BOOLEAN NOT NULL DEFAULT true,

    -- 元数据（记录业务相关配置）
    metadata JSONB DEFAULT '{}',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建索引
CREATE UNIQUE INDEX idx_system_accounts_type ON system_accounts(account_type, is_active)
WHERE account_type IS NOT NULL AND is_active = true
COMMENT ON idx_system_accounts_type IS '同类型激活账户的唯一索引';

-- ============================================
-- 外键说明（将来需要关联）
-- ============================================
-- task_escrow 表会引用 system_accounts(id)
-- platform_revenue 表会引用 system_accounts(id)
--
-- 外键约束示例（将来实施时添加）
-- ALTER TABLE tasks ADD CONSTRAINT fk_task_escrow_system_account
--     FOREIGN KEY (task_escrow) REFERENCES system_accounts(id) ON DELETE SET NULL;
--
-- ALTER TABLE platform_revenue ADD CONSTRAINT fk_platform_revenue_system_account
--     FOREIGN KEY (platform_revenue) REFERENCES system_accounts(id) ON DELETE SET NULL;
-- ============================================
