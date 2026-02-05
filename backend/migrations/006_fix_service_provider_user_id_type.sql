-- 修复 service_providers 表的 user_id 和 admin_id 字段类型
-- 从 UUID 改为 varchar(255) 以支持自定义格式的用户ID（usr_ 前缀）

-- 1. 修改 admin_id 字段类型
ALTER TABLE service_providers ALTER COLUMN admin_id TYPE VARCHAR(255);

-- 2. 修改 user_id 字段类型
ALTER TABLE service_providers ALTER COLUMN user_id TYPE VARCHAR(255);

-- 3. 删除 logo_url 字段（不再需要）
ALTER TABLE service_providers DROP COLUMN IF EXISTS logo_url;

-- 4. 修改 service_provider_staff 表的 user_id 字段类型
ALTER TABLE service_provider_staff ALTER COLUMN user_id TYPE VARCHAR(255);
