-- ============================================
-- 扩展 CreditAccount 表添加冻结余额字段
-- 用途：支持积分冻结机制
-- 创建时间:2026-02-12
-- ============================================

-- 添加冻结余额字段（用户余额冻结模式）
ALTER TABLE credit_accounts
  ADD COLUMN IF NOT EXISTS frozen_balance INT DEFAULT 0;

-- 添加字段注释
COMMENT ON COLUMN credit_accounts.frozen_balance IS '冻结的积分余额（单位：积分），用于提现申请、任务发布等场景';

-- ============================================
-- 设计说明
-- ============================================
-- 字段含义：
-- - balance: 可用积分余额（用户可以立即使用的积分）
-- - frozen_balance: 冻结积分余额（暂时冻结的积分，如提现申请中、任务发布中等）
--
-- 数据流：
-- 1. 用户申请提现 → balance -= amount, frozen_balance += amount
-- 2. 提现审核通过 → frozen_balance -= amount（积分被消耗）
-- 3. 提现审核拒绝 → frozen_balance -= amount, balance += amount（退回可用余额）
--
-- 总积分资产 = balance + frozen_balance
--
-- ============================================
-- 注意：已移除 cash_balance_diamond 和 cash_balance_gold 字段
-- 原因：统一使用"积分"作为唯一术语，不再区分钻石和金币
-- ============================================
