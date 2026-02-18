-- ============================================
-- 删除 last_used_role 遗留字段
-- 版本: 1.0
-- 日期: 2026-02-15
-- 说明: 系统已改为多角色设计，不再需要"上次使用的角色"记录
-- ============================================

-- 删除 last_used_role 列
ALTER TABLE users DROP COLUMN IF EXISTS last_used_role;

-- 验证
\d+ users

echo '✅ Migration 030 completed: removed last_used_role column'
