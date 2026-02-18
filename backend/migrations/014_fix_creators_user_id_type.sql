-- 修复 creators 表的 user_id 字段类型
-- 从 UUID 改为 VARCHAR(255) 以与 users(id) 一致，支持懒创建达人记录（usr_ 前缀）

-- 1. 删除外键（若存在）
ALTER TABLE creators DROP CONSTRAINT IF EXISTS fk_creators_user_id;

-- 2. 修改 user_id 类型（若有数据则用 ::text 转换）
ALTER TABLE creators ALTER COLUMN user_id TYPE VARCHAR(255) USING user_id::text;

-- 3. 重新添加外键
ALTER TABLE creators ADD CONSTRAINT fk_creators_user_id FOREIGN KEY (user_id) REFERENCES users(id);

COMMENT ON COLUMN creators.user_id IS '关联 users.id（与 users 主键类型一致）';
