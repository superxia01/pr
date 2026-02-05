-- 修改 merchants 表的 admin_id 字段为可空
-- 这样创建商家时不需要立即指定管理员，可以通过邀请码后续绑定

-- 1. 修改 admin_id 为可空
ALTER TABLE merchants ALTER COLUMN admin_id DROP NOT NULL;

-- 2. 删除现有的 admin_id 唯一索引（因为可以为空了）
-- ALTER TABLE merchants DROP INDEX IF EXISTS merchants_admin_id_key;

-- 3. 创建新的唯一索引，只对非空值建立唯一约束
CREATE UNIQUE INDEX merchants_admin_id_key ON merchants(admin_id) WHERE admin_id IS NOT NULL AND admin_id != '';
