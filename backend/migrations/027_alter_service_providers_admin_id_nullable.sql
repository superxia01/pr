-- 修改 service_providers 表的 admin_id 字段为可空
-- 允许多个服务商暂时没有管理员（NULL 值不违反唯一约束）

-- 1. 删除现有的唯一约束（如果存在）
DROP INDEX IF EXISTS idx_service_providers_admin_id;

-- 2. 修改 admin_id 字段为可空
ALTER TABLE service_providers
ALTER COLUMN admin_id DROP NOT NULL;

-- 3. 将现有的空字符串更新为 NULL
UPDATE service_providers
SET admin_id = NULL
WHERE admin_id = '';

-- 4. 重新创建唯一约束（允许 NULL 值）
CREATE UNIQUE INDEX idx_service_providers_admin_id
ON service_providers(admin_id)
WHERE admin_id IS NOT NULL;

-- 注：PostgreSQL 的部分索引（partial index）允许在 WHERE 条件下创建唯一索引
-- 这样多个 NULL 值不会违反唯一约束
