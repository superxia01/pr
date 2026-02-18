-- 统一角色值为大写下划线格式
-- 迁移 012: normalize_role_values.sql (修复版)

-- 1. 更新 users 表的 active_role
UPDATE users SET active_role = 'SUPER_ADMIN' WHERE active_role = 'super_admin';
UPDATE users SET active_role = 'SERVICE_PROVIDER_ADMIN' WHERE active_role = 'SERVICE_PROVIDER_ADMIN' OR active_role = 'provider_admin' OR active_role = 'SP_ADMIN';
UPDATE users SET active_role = 'SERVICE_PROVIDER_STAFF' WHERE active_role = 'SERVICE_PROVIDER_STAFF' OR active_role = 'provider_staff' OR active_role = 'SP_STAFF';
UPDATE users SET active_role = 'MERCHANT_ADMIN' WHERE active_role = 'merchant_admin' OR active_role = 'MERCHANT_ADMIN';
UPDATE users SET active_role = 'MERCHANT_STAFF' WHERE active_role = 'merchant_staff' OR active_role = 'MERCHANT_STAFF';
UPDATE users SET active_role = 'CREATOR' WHERE active_role = 'creator' OR active_role = 'CREATOR';

-- 2. 更新 users 表的 roles 数组 (JSONB)
CREATE OR REPLACE FUNCTION normalize_user_roles() RETURNS VOID AS $$
DECLARE
    user_record RECORD;
    new_roles JSONB;
BEGIN
    FOR user_record IN SELECT id, roles FROM users WHERE roles IS NOT NULL AND roles != '[]'::jsonb LOOP
        -- 将 JSONB 数组转换为文本数组进行处理
        new_roles := (
            SELECT jsonb_agg(
                CASE
                    WHEN role_text = 'super_admin' THEN 'SUPER_ADMIN'::text
                    WHEN role_text = 'SERVICE_PROVIDER_ADMIN' OR role_text = 'provider_admin' OR role_text = 'SP_ADMIN' THEN 'SERVICE_PROVIDER_ADMIN'::text
                    WHEN role_text = 'SERVICE_PROVIDER_STAFF' OR role_text = 'provider_staff' OR role_text = 'SP_STAFF' THEN 'SERVICE_PROVIDER_STAFF'::text
                    WHEN role_text = 'merchant_admin' OR role_text = 'MERCHANT_ADMIN' THEN 'MERCHANT_ADMIN'::text
                    WHEN role_text = 'merchant_staff' OR role_text = 'MERCHANT_STAFF' THEN 'MERCHANT_STAFF'::text
                    WHEN role_text = 'creator' OR role_text = 'CREATOR' THEN 'CREATOR'::text
                    ELSE UPPER(role_text)
                END
            )
            FROM jsonb_array_elements_text(user_record.roles) AS role_text
        );

        UPDATE users
        SET roles = new_roles
        WHERE id = user_record.id;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- 执行函数
SELECT normalize_user_roles();

-- 删除临时函数
DROP FUNCTION normalize_user_roles();

-- 3. 添加 CHECK 约束确保角色值正确（如果约束已存在则忽略）
DO $$
BEGIN
    ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS users_active_role_check
        CHECK (active_role IN ('SUPER_ADMIN', 'SERVICE_PROVIDER_ADMIN', 'SERVICE_PROVIDER_STAFF', 'MERCHANT_ADMIN', 'MERCHANT_STAFF', 'CREATOR') OR active_role IS NULL);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END;
$$;

-- 4. 验证结果
SELECT id, active_role, roles
FROM users
WHERE active_role IS NOT NULL
LIMIT 10;
