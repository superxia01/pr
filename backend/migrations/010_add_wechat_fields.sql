-- 添加微信相关独立字段
-- 版本: v1.0
-- 日期: 2026-02-06
-- 说明: 将 union_id 和 open_id 从 Profile JSONB 提取为独立字段，提高查询效率

-- 添加 open_id 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS open_id VARCHAR(255);

-- 添加 union_id 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS union_id VARCHAR(255);

-- 添加 login_type 字段（区分登录方式）
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_type VARCHAR(20) DEFAULT 'local' CHECK (login_type IN ('local', 'wechat'));

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_open_id ON users(open_id);
CREATE INDEX IF NOT EXISTS idx_users_union_id ON users(union_id);

-- 添加注释
COMMENT ON COLUMN users.open_id IS '微信 OpenID（应用唯一）';
COMMENT ON COLUMN users.union_id IS '微信 UnionID（跨应用唯一）';
COMMENT ON COLUMN users.login_type IS '登录类型：local（本地密码）或 wechat（微信）';

-- 迁移现有数据：从 Profile 中提取 union_id 和 open_id
UPDATE users
SET
    open_id = (profile->>'open_id')::VARCHAR(255),
    union_id = (profile->>'union_id')::VARCHAR(255),
    login_type = CASE WHEN (profile->>'union_id') IS NOT NULL THEN 'wechat' ELSE 'local' END
WHERE profile IS NOT NULL;

-- 迁移完成后，可以将 Profile 中的旧数据设为 NULL（可选）
-- UPDATE users SET profile = profile - 'open_id' - 'union_id' WHERE profile ? 'open_id' OR profile ? 'union_id';
