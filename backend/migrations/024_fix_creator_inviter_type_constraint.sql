-- 修复 creators 表的 inviter_type 约束拼写错误
-- 原约束：PROVIDER_STAFF, PROVIDER_ADMIN
-- 修正为：SERVICE_PROVIDER_STAFF, SERVICE_PROVIDER_ADMIN

-- 1. 删除旧约束
ALTER TABLE creators DROP CONSTRAINT IF EXISTS creators_inviter_type_check;

-- 2. 添加正确的新约束
ALTER TABLE creators
    ADD CONSTRAINT creators_inviter_type_check
    CHECK (inviter_type IS NULL OR inviter_type IN ('SERVICE_PROVIDER_STAFF', 'SERVICE_PROVIDER_ADMIN', 'OTHER'));

-- 3. 修正现有数据中的错误值（如果存在）
UPDATE creators
SET inviter_type = CASE
    WHEN inviter_type = 'PROVIDER_STAFF' THEN 'SERVICE_PROVIDER_STAFF'
    WHEN inviter_type = 'PROVIDER_ADMIN' THEN 'SERVICE_PROVIDER_ADMIN'
    ELSE inviter_type
END
WHERE inviter_type IN ('PROVIDER_STAFF', 'PROVIDER_ADMIN');

-- 注释
COMMENT ON COLUMN creators.inviter_type IS '邀请人类型：SERVICE_PROVIDER_STAFF-服务商员工, SERVICE_PROVIDER_ADMIN-服务商管理员, OTHER-其他';
