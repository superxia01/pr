-- ============================================
-- PR Business 财务系统数据库迁移脚本
-- 统一积分术语，移除 DIAMOND/GOLD 区分
-- 生成时间: 2026-02-12
-- ============================================

-- 使用说明：
-- 1. 确保数据库连接配置正确 (.env 文件)
-- 2. 执行顺序：按照文件编号顺序执行
-- 3. 验证：每步执行后检查是否成功

-- ============================================
-- 迁移 018: 创建现金账户表
-- ============================================

\echo '执行迁移 018: 创建现金账户表...'
CREATE TABLE IF NOT EXISTS cash_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_type VARCHAR(50) NOT NULL,
    balance INT NOT NULL DEFAULT 0,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cash_accounts_type_active ON cash_accounts(account_type, is_active);

COMMENT ON COLUMN cash_accounts.account_type IS '账户类型：WECHAT, ALIPAY, BANK_TRANSFER, MARKETING, OPERATIONS';
COMMENT ON COLUMN cash_accounts.balance IS '余额（单位：分）';

-- ============================================
-- 迁移 019: 创建系统账户表
-- ============================================

\echo '执行迁移 019: 创建系统账户表...'
CREATE TABLE IF NOT EXISTS system_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_type VARCHAR(50) NOT NULL,
    balance INT NOT NULL DEFAULT 0,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_system_accounts_type ON system_accounts(account_type) WHERE is_active = true;

COMMENT ON COLUMN system_accounts.account_type IS '账户类型：TICKET_ESCROW, TASK_ESCROW, PLATFORM_REVENUE';
COMMENT ON COLUMN system_accounts.balance IS '余额（单位：积分）';

-- ============================================
-- 迁移 020: 扩展积分账户表添加冻结余额字段
-- ============================================

\echo '执行迁移 020: 添加冻结余额字段...'
DO $$
BEGIN
    -- 检查列是否已存在
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'credit_accounts'
        AND column_name = 'frozen_balance'
    ) THEN
        ALTER TABLE credit_accounts ADD COLUMN frozen_balance INT DEFAULT 0;
        COMMENT ON COLUMN credit_accounts.frozen_balance IS '冻结的积分余额（单位：积分），用于提现申请、任务发布等场景';
    END IF;
END $$;

-- ============================================
-- 迁移 021: 创建提现申请表
-- ============================================

\echo '执行迁移 021: 创建提现申请表...'
CREATE TABLE IF NOT EXISTS withdrawal_requests_enhanced (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES credit_accounts(id) ON DELETE SET NULL,
    amount INT NOT NULL,
    yuan_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    cash_account_type VARCHAR(50),
    description TEXT,
    reject_reason TEXT,
    reviewed_by VARCHAR(255),
    reviewed_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_withdrawal_requests_account_status ON withdrawal_requests_enhanced(account_id, status);

COMMENT ON TABLE withdrawal_requests_enhanced IS '提现申请表（增强版-带冻结机制）';
COMMENT ON COLUMN withdrawal_requests_enhanced.amount IS '提现金额（单位：积分）';
COMMENT ON COLUMN withdrawal_requests_enhanced.yuan_amount IS '提现金额（单位：元）';

-- ============================================
-- 迁移 022: 创建财务审计日志表
-- ============================================

\echo '执行迁移 022: 创建财务审计日志表...'
CREATE TABLE IF NOT EXISTS financial_audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id VARCHAR(255) NOT NULL,
    changes JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user_action ON financial_audit_logs(user_id, action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON financial_audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON financial_audit_logs(created_at DESC);

COMMENT ON TABLE financial_audit_logs IS '财务审计日志表（简化版）';

-- ============================================
-- 验证迁移结果
-- ============================================

\echo ''
\echo '============================================'
\echo '迁移验证'
\echo '============================================'
\echo ''

-- 检查新创建的表
DO $$
BEGIN
    DECLARE
        table_count INT;
    BEGIN
        -- 检查 cash_accounts 表
        SELECT COUNT(*) INTO table_count FROM information_schema.tables
        WHERE table_name = 'cash_accounts';
        IF table_count > 0 THEN
            RAISE NOTICE '✅ cash_accounts 表已创建';
        ELSE
            RAISE NOTICE '❌ cash_accounts 表创建失败';
        END IF;

        -- 检查 system_accounts 表
        SELECT COUNT(*) INTO table_count FROM information_schema.tables
        WHERE table_name = 'system_accounts';
        IF table_count > 0 THEN
            RAISE NOTICE '✅ system_accounts 表已创建';
        ELSE
            RAISE NOTICE '❌ system_accounts 表创建失败';
        END IF;

        -- 检查 withdrawal_requests_enhanced 表
        SELECT COUNT(*) INTO table_count FROM information_schema.tables
        WHERE table_name = 'withdrawal_requests_enhanced';
        IF table_count > 0 THEN
            RAISE NOTICE '✅ withdrawal_requests_enhanced 表已创建';
        ELSE
            RAISE NOTICE '❌ withdrawal_requests_enhanced 表创建失败';
        END IF;

        -- 检查 financial_audit_logs 表
        SELECT COUNT(*) INTO table_count FROM information_schema.tables
        WHERE table_name = 'financial_audit_logs';
        IF table_count > 0 THEN
            RAISE NOTICE '✅ financial_audit_logs 表已创建';
        ELSE
            RAISE NOTICE '❌ financial_audit_logs 表创建失败';
        END IF;

        -- 检查 frozen_balance 字段
        SELECT COUNT(*) INTO table_count FROM information_schema.columns
        WHERE table_name = 'credit_accounts'
        AND column_name = 'frozen_balance';
        IF table_count > 0 THEN
            RAISE NOTICE '✅ frozen_balance 字段已添加';
        ELSE
            RAISE NOTICE '❌ frozen_balance 字段添加失败';
        END IF;

    END;
END $$;

-- 显示当前积分账户表结构
\echo ''
\echo 'credit_accounts 表结构：'
\SELECT
    column_name AS "字段名",
    data_type AS "数据类型",
    is_nullable AS "可空",
    column_default AS "默认值"
FROM information_schema.columns
WHERE table_name = 'credit_accounts'
ORDER BY ordinal_position;

\echo ''
\echo '============================================'
\echo '迁移完成！'
\echo '============================================'
