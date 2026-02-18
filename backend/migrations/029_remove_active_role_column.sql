-- ============================================
-- 删除 active_role 遗留字段
-- 版本: 1.0
-- 日期: 2026-02-15
-- 说明: 系统已改为多角色设计，不再使用"当前激活角色"概念
-- ============================================

-- 删除 active_role 列
ALTER TABLE users DROP COLUMN IF EXISTS active_role;

-- 删除相关索引
DROP INDEX IF EXISTS idx_users_active_role;

-- 删除相关约束
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_active_role_check;

-- 验证
\d+ users

echo '✅ Migration 029 completed: removed active_role column'
