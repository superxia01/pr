-- 添加 auth_center_token 字段，用于存储 auth-center 的 token
-- 这个 token 用来调用 auth-center 的用户信息接口

-- 1. 添加字段
ALTER TABLE users ADD COLUMN auth_center_token VARCHAR(500);

-- 2. 添加注释
COMMENT ON COLUMN users.auth_center_token IS 'auth-center 的 access token，用于调用 auth-center API';
