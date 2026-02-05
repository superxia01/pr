-- ============================================
-- PR Business 数据修复脚本
-- 修复数据库中格式不正确的 profile 和 roles 字段
--
-- 问题描述：
--  - 表结构定义为 JSONB 类型
--  - 但已有数据以字符串格式存储（如 "{}" 而不是 {}）
--  - 导致 GORM Scan 时报错
--
-- 执行时机：在部署新版本前执行此脚本
-- ============================================

-- ============================================
-- 第一步：备份数据（安全起见）
-- ============================================

-- 备份 users 表
CREATE TABLE IF NOT EXISTS users_backup_20260204 AS
SELECT * FROM users;

-- 验证备份
SELECT COUNT(*) AS backup_count FROM users_backup_20260204;


-- ============================================
-- 第二步：修复 profile 字段
-- ============================================

-- 修复：将字符串格式转换为真正的JSONB
UPDATE users
SET profile = CASE
    WHEN profile IS NULL THEN '{}'::jsonb
    WHEN profile = '' THEN '{}'::jsonb
    WHEN profile::text = '{}' THEN '{}'::jsonb
    WHEN profile::text = '[]' THEN '[]'::jsonb
    WHEN substring(profile::text, 1, 1) = '{' THEN profile::jsonb  -- 已经是JSON格式
    ELSE '{}'::jsonb
END;

-- 验证修复
SELECT id, nickname, profile, jsonb_pretty(profile) as profile_formatted
FROM users
LIMIT 5;


-- ============================================
-- 第三步：修复 roles 字段
-- ============================================

-- 修复：将字符串格式转换为真正的JSONB
UPDATE users
SET roles = CASE
    WHEN roles IS NULL THEN '[]'::jsonb
    WHEN roles = '' THEN '[]'::jsonb
    WHEN roles::text = '{}' THEN '{}'::jsonb
    WHEN roles::text = '[]' THEN '[]'::jsonb
    WHEN substring(roles::text, 1, 1) = '[' THEN roles::jsonb  -- 已经是JSON格式
    ELSE '[]'::jsonb
END;

-- 验证修复
SELECT id, nickname, roles, jsonb_pretty(roles) as roles_formatted
FROM users
LIMIT 5;


-- ============================================
-- 第四步：验证修复结果
-- ============================================

-- 统计修复后的数据
SELECT
    COUNT(*) as total_users,
    COUNT(*) FILTER (profile IS NOT NULL) as users_with_profile,
    COUNT(*) FILTER (jsonb_typeof(profile) = 'object') as valid_profile_objects,
    COUNT(*) FILTER (array_length(roles::text[]::jsonb, 1) > 0) as valid_roles_arrays
FROM users;


-- ============================================
-- 第五步：回滚方案（如果需要）
-- ============================================

-- 如果修复有问题，可以从备份恢复：
-- DROP TABLE users;
-- CREATE TABLE users AS TABLE users_backup_20260204;
-- ALTER TABLE users OWNER TO nexus_user;
