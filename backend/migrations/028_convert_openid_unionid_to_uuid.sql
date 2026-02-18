-- ============================================
-- 转换 open_id/union_id 列为 UUID 类型
-- 版本: 1.0
-- 日期: 2026-02-15
-- 说明: 将 users 表的 open_id 和 union_id 从 varchar(255) 转换为 uuid 类型
-- ============================================

-- 删除相关索引（如果存在）
DROP INDEX IF EXISTS idx_users_open_id;
DROP INDEX IF EXISTS idx_users_union_id;

-- 转换 open_id 列为 uuid
ALTER TABLE users ALTER COLUMN open_id TYPE uuid USING (
  CASE
    WHEN open_id ~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$' THEN open_id::uuid
    ELSE NULL::uuid
  END
);

-- 转换 union_id 列为 uuid
ALTER TABLE users ALTER COLUMN union_id TYPE uuid USING (
  CASE
    WHEN union_id ~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$' THEN union_id::uuid
    ELSE NULL::uuid
  END
);

-- 重新创建索引
CREATE INDEX idx_users_open_id ON users(open_id);
CREATE INDEX idx_users_union_id ON users(union_id);

-- 验证列类型已更改
\d+ users

-- 验证索引已创建
\di idx_users_open_id
\di idx_users_union_id

echo '✅ Migration 028 completed: open_id/union_id converted to uuid'
