-- 修复 merchants 表的 user_id 和 admin_id 字段类型
-- 从 UUID 改为 varchar(255) 以支持自定义格式的用户ID（usr_ 前缀）

-- 1. 修改 admin_id 字段类型
ALTER TABLE merchants ALTER COLUMN admin_id TYPE VARCHAR(255);

-- 2. 修改 user_id 字段类型
ALTER TABLE merchants ALTER COLUMN user_id TYPE VARCHAR(255);

-- 3. 删除 logo_url 字段（不再需要）
ALTER TABLE merchants DROP COLUMN IF EXISTS logo_url;

-- 4. 修改 merchant_staff 表的 user_id 字段类型
ALTER TABLE merchant_staff ALTER COLUMN user_id TYPE VARCHAR(255);
