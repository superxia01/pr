-- ============================================
-- Migration 004: 重命名 current_role 为 active_role
-- ============================================
-- 原因: current_role 是 PostgreSQL 保留字，会导致 SQL 查询冲突
-- 解决: 重命名为 active_role 避免冲突
-- ============================================

-- 1. 重命名列（必须加引号，因为 current_role 是 PostgreSQL 保留字）
ALTER TABLE users RENAME COLUMN "current_role" TO active_role;

-- 2. 重命名索引
DROP INDEX IF EXISTS idx_users_current_role;
CREATE INDEX idx_users_active_role ON users(active_role);

-- 3. 更新注释
COMMENT ON COLUMN users.active_role IS '用户当前激活的角色，用于工作台切换';

-- 4. 验证迁移
SELECT
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'users' AND column_name = 'active_role';
