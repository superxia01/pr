-- ============================================
-- PR Business 账户权限系统
-- 版本: v1.1
-- 日期: 2026-02-04
-- 说明: 添加 account_permissions 关联表，支持多对多用户-账户关系
-- ============================================

-- ============================================
-- 1. 修改 credit_accounts 表，移除 user_id 的 NOT NULL 和 UNIQUE 约束
-- ============================================

-- 删除旧的 UNIQUE 约束（如果存在）
ALTER TABLE credit_accounts DROP CONSTRAINT IF EXISTS credit_accounts_user_id_key;

-- 允许 user_id 为空
ALTER TABLE credit_accounts ALTER COLUMN user_id DROP NOT NULL;

-- ============================================
-- 2. 创建账户权限关联表
-- ============================================

CREATE TABLE IF NOT EXISTS account_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    user_id UUID NOT NULL,
    account_type VARCHAR(20) NOT NULL,  -- 'PERSONAL', 'ORG_MERCHANT', 'ORG_PROVIDER'
    can_view BOOLEAN NOT NULL DEFAULT TRUE,
    can_operate BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(account_id, user_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_account_permissions_user_id ON account_permissions(user_id);
CREATE INDEX IF NOT EXISTS idx_account_permissions_account_id ON account_permissions(account_id);

-- 外键约束
ALTER TABLE account_permissions DROP CONSTRAINT IF EXISTS account_permissions_account_id_fkey;
ALTER TABLE account_permissions ADD CONSTRAINT account_permissions_account_id_fkey
    FOREIGN KEY (account_id) REFERENCES credit_accounts(id) ON DELETE CASCADE;

-- 添加注释
COMMENT ON TABLE account_permissions IS '账户权限关联表：定义哪些用户可以访问哪些账户';
COMMENT ON COLUMN account_permissions.account_id IS '关联的账户ID';
COMMENT ON COLUMN account_permissions.user_id IS '用户ID';
COMMENT ON COLUMN account_permissions.can_view IS '是否可以查看账户';
COMMENT ON COLUMN account_permissions.can_operate IS '是否可以操作账户（提现、转账等）';

-- ============================================
-- 3. 为现有用户自动创建账户权限（数据迁移）
-- ============================================

-- 3.1 为所有用户创建个人账户权限
INSERT INTO account_permissions (account_id, user_id, account_type, can_view, can_operate)
SELECT
    ca.id,
    u.auth_center_user_id::uuid,
    'PERSONAL',
    true,
    true
FROM users u
JOIN credit_accounts ca ON ca.owner_id::text = u.auth_center_user_id AND ca.owner_type = 'USER_PERSONAL'
WHERE NOT EXISTS (
    SELECT 1 FROM account_permissions ap
    WHERE ap.account_id = ca.id AND ap.user_id::text = u.auth_center_user_id
);

-- 3.2 为商家管理员创建商家账户权限
INSERT INTO account_permissions (account_id, user_id, account_type, can_view, can_operate)
SELECT
    ca.id,
    u.auth_center_user_id::uuid,
    'ORG_MERCHANT',
    true,
    true
FROM users u
JOIN merchants m ON m.admin_id::text = u.auth_center_user_id
JOIN credit_accounts ca ON ca.owner_id = m.id AND ca.owner_type = 'ORG_MERCHANT'
WHERE NOT EXISTS (
    SELECT 1 FROM account_permissions ap
    WHERE ap.account_id = ca.id AND ap.user_id::text = u.auth_center_user_id
);

-- 3.3 为商家员工创建商家账户权限（有财务权限的）
INSERT INTO account_permissions (account_id, user_id, account_type, can_view, can_operate)
SELECT
    ca.id,
    u.auth_center_user_id::uuid,
    'ORG_MERCHANT',
    true,
    false  -- 员工默认不能操作
FROM users u
JOIN merchant_staff ms ON ms.user_id::text = u.auth_center_user_id AND ms.status = 'active'
JOIN credit_accounts ca ON ca.owner_id = ms.merchant_id AND ca.owner_type = 'ORG_MERCHANT'
WHERE NOT EXISTS (
    SELECT 1 FROM account_permissions ap
    WHERE ap.account_id = ca.id AND ap.user_id::text = u.auth_center_user_id
);
-- 注意：这里应该检查财务权限，但由于当前没有权限系统，暂时给所有员工查看权限

-- 3.4 为服务商管理员创建服务商账户权限
INSERT INTO account_permissions (account_id, user_id, account_type, can_view, can_operate)
SELECT
    ca.id,
    u.auth_center_user_id::uuid,
    'ORG_PROVIDER',
    true,
    true
FROM users u
JOIN service_providers sp ON sp.admin_id::text = u.auth_center_user_id
JOIN credit_accounts ca ON ca.owner_id = sp.id AND ca.owner_type = 'ORG_PROVIDER'
WHERE NOT EXISTS (
    SELECT 1 FROM account_permissions ap
    WHERE ap.account_id = ca.id AND ap.user_id::text = u.auth_center_user_id
);

-- 3.5 为服务商员工创建服务商账户权限
INSERT INTO account_permissions (account_id, user_id, account_type, can_view, can_operate)
SELECT
    ca.id,
    u.auth_center_user_id::uuid,
    'ORG_PROVIDER',
    true,
    false  -- 员工默认不能操作
FROM users u
JOIN service_provider_staff sps ON sps.user_id::text = u.auth_center_user_id AND sps.status = 'active'
JOIN credit_accounts ca ON ca.owner_id = sps.provider_id AND ca.owner_type = 'ORG_PROVIDER'
WHERE NOT EXISTS (
    SELECT 1 FROM account_permissions ap
    WHERE ap.account_id = ca.id AND ap.user_id::text = u.auth_center_user_id
);

-- ============================================
-- 完成提示
-- ============================================

DO $$
BEGIN
    RAISE NOTICE '账户权限系统创建完成！';
    RAISE NOTICE '已为现有用户自动创建账户权限关联';
    RAISE NOTICE '下一步：更新后端代码以支持新的权限系统';
END $$;
