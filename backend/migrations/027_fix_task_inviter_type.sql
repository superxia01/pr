-- ============================================
-- 修复 Task 模型 inviter_type 约束值
-- 版本: 1.0
-- 日期: 2026-02-15
-- 说明: 统一 tasks 表的 inviter_type 约束值与 creators 表一致
-- ============================================

-- 删除旧的约束
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_inviter_type_check;

-- 添加新的约束（与 creators 表一致）
ALTER TABLE tasks ADD CONSTRAINT tasks_inviter_type_check 
  CHECK (inviter_type IS NULL OR inviter_type::text = ANY(ARRAY[
    'SERVICE_PROVIDER_STAFF'::character varying::text,
    'SERVICE_PROVIDER_ADMIN'::character varying::text,
    'OTHER'::character varying::text
  ]));

-- 验证约束已添加
\echo '✅ Migration 027 completed: tasks.inviter_type constraint fixed'
