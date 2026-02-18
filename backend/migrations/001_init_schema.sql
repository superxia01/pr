-- PR Business 数据库初始化脚本
-- 版本: v1.0
-- 日期: 2026-02-02
-- 说明: 创建所有核心表、索引、外键、触发器

-- ============================================
-- 扩展模块
-- ============================================

-- 启用UUID扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 启用JSONB索引优化
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================
-- 1. 用户表 (users)
-- ============================================

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY DEFAULT (encode(decode(random_uuid()::text, 'base64'), 'hex')),
    auth_center_user_id UUID NOT NULL,
    nickname VARCHAR(50),
    avatar_url VARCHAR(500),
    profile JSONB DEFAULT '{}'::jsonb,
    roles JSONB NOT NULL DEFAULT '[]'::jsonb,
    current_role VARCHAR(50),
    last_used_role VARCHAR(50),
    invited_by VARCHAR(255),
    invitation_code_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'banned', 'inactive')),
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引
CREATE UNIQUE INDEX idx_users_auth_center_user_id ON users(auth_center_user_id);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_current_role ON users(current_role);
CREATE INDEX idx_users_invited_by ON users(invited_by);
CREATE INDEX idx_users_invitation_code_id ON users(invitation_code_id);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_roles ON users USING GIN(roles);
CREATE INDEX idx_users_profile ON users USING GIN(profile);

-- 注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '用户唯一标识（本地主键，CUID格式）';
COMMENT ON COLUMN users.auth_center_user_id IS '关联 auth_center_db.users.user_id（核心字段）';
COMMENT ON COLUMN users.roles IS '用户拥有的所有角色，例如：["MERCHANT_ADMIN", "CREATOR"]';
COMMENT ON COLUMN users.current_role IS '用户当前激活的角色，用于工作台切换';

-- ============================================
-- 2. 邀请码表 (invitation_codes)
-- ============================================

CREATE TABLE invitation_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(30) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    target_role VARCHAR(50) NOT NULL,
    generator_id UUID NOT NULL,
    generator_type VARCHAR(50) NOT NULL CHECK (generator_type IN ('system', 'super_admin', 'merchant_admin', 'sp_admin', 'sp_staff')),
    organization_id UUID,
    organization_type VARCHAR(50) CHECK (organization_type IS NULL OR organization_type IN ('merchant', 'service_provider')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_one_time BOOLEAN NOT NULL DEFAULT false,
    max_uses INT,
    use_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_invitation_codes_type ON invitation_codes(type);
CREATE INDEX idx_invitation_codes_generator_id ON invitation_codes(generator_id);
CREATE INDEX idx_invitation_codes_organization_id ON invitation_codes(organization_id);
CREATE INDEX idx_invitation_codes_is_active ON invitation_codes(is_active);

-- 外键
ALTER TABLE invitation_codes ADD CONSTRAINT fk_invitation_codes_generator_id FOREIGN KEY (generator_id) REFERENCES users(id);

-- 数据完整性约束
ALTER TABLE invitation_codes ADD CONSTRAINT check_generator_consistency
CHECK (
    (generator_type = 'system' AND generator_id IS NULL) OR
    (generator_type IN ('super_admin', 'merchant_admin', 'sp_admin', 'sp_staff') AND generator_id IS NOT NULL)
);

COMMENT ON TABLE invitation_codes IS '邀请码表';

-- ============================================
-- 3. 商家表 (merchants)
-- ============================================

CREATE TABLE merchants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL,
    provider_id UUID NOT NULL,
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    logo_url VARCHAR(500),
    industry VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引
CREATE UNIQUE INDEX idx_merchants_admin_id ON merchants(admin_id);
CREATE INDEX idx_merchants_provider_id ON merchants(provider_id);
CREATE INDEX idx_merchants_user_id ON merchants(user_id);
CREATE INDEX idx_merchants_status ON merchants(status);
CREATE INDEX idx_merchants_name ON merchants(name);

-- 外键
ALTER TABLE merchants ADD CONSTRAINT fk_merchants_admin_id FOREIGN KEY (admin_id) REFERENCES users(id);
ALTER TABLE merchants ADD CONSTRAINT fk_merchants_provider_id FOREIGN KEY (provider_id) REFERENCES service_providers(id);
ALTER TABLE merchants ADD CONSTRAINT fk_merchants_user_id FOREIGN KEY (user_id) REFERENCES users(id);

COMMENT ON TABLE merchants IS '商家表';
COMMENT ON COLUMN merchants.provider_id IS '关联 service_providers.id（创建该商家的服务商）';

-- ============================================
-- 4. 服务商表 (service_providers)
-- ============================================

CREATE TABLE service_providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL,
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    logo_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引
CREATE UNIQUE INDEX idx_service_providers_admin_id ON service_providers(admin_id);
CREATE INDEX idx_service_providers_user_id ON service_providers(user_id);
CREATE INDEX idx_service_providers_status ON service_providers(status);
CREATE INDEX idx_service_providers_name ON service_providers(name);

-- 外键
ALTER TABLE service_providers ADD CONSTRAINT fk_service_providers_admin_id FOREIGN KEY (admin_id) REFERENCES users(id);
ALTER TABLE service_providers ADD CONSTRAINT fk_service_providers_user_id FOREIGN KEY (user_id) REFERENCES users(id);

COMMENT ON TABLE service_providers IS '服务商表';

-- ============================================
-- 5. 达人表 (creators)
-- ============================================

CREATE TABLE creators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT TRUE,
    level VARCHAR(20) NOT NULL DEFAULT 'UGC' CHECK (level IN ('UGC', 'KOC', 'INF', 'KOL')),
    followers_count INT NOT NULL DEFAULT 0,
    wechat_openid VARCHAR(100),
    wechat_nickname VARCHAR(100),
    wechat_avatar VARCHAR(500),
    inviter_id UUID,
    inviter_type VARCHAR(50) CHECK (inviter_type IS NULL OR inviter_type IN ('PROVIDER_STAFF', 'PROVIDER_ADMIN', 'OTHER')),
    inviter_relationship_broken BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'banned', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_creators_user_id ON creators(user_id);
CREATE INDEX idx_creators_level ON creators(level);
CREATE INDEX idx_creators_inviter_id ON creators(inviter_id);
CREATE INDEX idx_creators_status ON creators(status);
CREATE INDEX idx_creators_is_primary ON creators(is_primary);
CREATE UNIQUE INDEX idx_creators_user_primary ON creators(user_id, is_primary) WHERE is_primary = true;

-- 外键
ALTER TABLE creators ADD CONSTRAINT fk_creators_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE creators ADD CONSTRAINT fk_creators_inviter_id FOREIGN KEY (inviter_id) REFERENCES users(id);

COMMENT ON TABLE creators IS '达人表';

-- ============================================
-- 6. 营销活动表 (campaigns)
-- ============================================

CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_id UUID NOT NULL,
    provider_id UUID,
    title VARCHAR(100) NOT NULL,
    requirements TEXT NOT NULL,
    platforms JSONB NOT NULL,
    task_amount INT NOT NULL,
    campaign_amount INT NOT NULL,
    creator_amount INT,
    staff_referral_amount INT,
    provider_amount INT,
    quota INT NOT NULL,
    task_deadline TIMESTAMP NOT NULL,
    submission_deadline TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'OPEN', 'CLOSED')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_campaigns_merchant_id ON campaigns(merchant_id);
CREATE INDEX idx_campaigns_provider_id ON campaigns(provider_id);
CREATE INDEX idx_campaigns_status ON campaigns(status);
CREATE INDEX idx_campaigns_created_at ON campaigns(created_at);
CREATE INDEX idx_campaigns_platforms ON campaigns USING GIN(platforms);

-- 外键
ALTER TABLE campaigns ADD CONSTRAINT fk_campaigns_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(id);
ALTER TABLE campaigns ADD CONSTRAINT fk_campaigns_provider_id FOREIGN KEY (provider_id) REFERENCES service_providers(id);

COMMENT ON TABLE campaigns IS '营销活动表';

-- ============================================
-- 7. 任务名额表 (tasks)
-- ============================================

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL,
    task_slot_number INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'ASSIGNED', 'SUBMITTED', 'APPROVED', 'REJECTED')),
    creator_id UUID,
    assigned_at TIMESTAMP,
    platform VARCHAR(50),
    platform_url VARCHAR(500),
    screenshots JSONB,
    submitted_at TIMESTAMP,
    notes TEXT,
    audited_by UUID,
    audited_at TIMESTAMP,
    audit_note TEXT,
    inviter_id UUID,
    inviter_type VARCHAR(50) CHECK (inviter_type IS NULL OR inviter_type IN ('PROVIDER_STAFF', 'MERCHANT_STAFF', 'OTHER')),
    priority VARCHAR(10) NOT NULL DEFAULT 'MEDIUM' CHECK (priority IN ('HIGH', 'MEDIUM', 'LOW')),
    tags TEXT[],
    version INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_tasks_campaign_id ON tasks(campaign_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_creator_id ON tasks(creator_id);
CREATE INDEX idx_tasks_inviter_id ON tasks(inviter_id);
CREATE INDEX idx_tasks_audited_by ON tasks(audited_by);
CREATE UNIQUE INDEX idx_tasks_campaign_slot ON tasks(campaign_id, task_slot_number);
CREATE UNIQUE INDEX idx_tasks_campaign_creator ON tasks(campaign_id, creator_id) WHERE creator_id IS NOT NULL;

-- 外键
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_campaign_id FOREIGN KEY (campaign_id) REFERENCES campaigns(id);
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_creator_id FOREIGN KEY (creator_id) REFERENCES creators(id);
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_audited_by FOREIGN KEY (audited_by) REFERENCES users(id);
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_inviter_id FOREIGN KEY (inviter_id) REFERENCES users(id);

COMMENT ON TABLE tasks IS '任务名额表';

-- ============================================
-- 8. 商家员工表 (merchant_staff)
-- ============================================

CREATE TABLE merchant_staff (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    merchant_id UUID NOT NULL,
    title VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_merchant_staff_merchant_id ON merchant_staff(merchant_id);
CREATE INDEX idx_merchant_staff_status ON merchant_staff(status);

-- 外键
ALTER TABLE merchant_staff ADD CONSTRAINT fk_merchant_staff_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE merchant_staff ADD CONSTRAINT fk_merchant_staff_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(id);

COMMENT ON TABLE merchant_staff IS '商家员工表';

-- ============================================
-- 9. 商家员工权限表 (merchant_staff_permissions)
-- ============================================

CREATE TABLE merchant_staff_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    staff_id UUID NOT NULL,
    permission_code VARCHAR(50) NOT NULL,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    granted_by UUID NOT NULL
);

-- 索引
CREATE UNIQUE INDEX idx_merchant_staff_permissions_staff_permission ON merchant_staff_permissions(staff_id, permission_code);
CREATE INDEX idx_merchant_staff_permissions_staff_id ON merchant_staff_permissions(staff_id);
CREATE INDEX idx_merchant_staff_permissions_permission_code ON merchant_staff_permissions(permission_code);

-- 外键
ALTER TABLE merchant_staff_permissions ADD CONSTRAINT fk_merchant_staff_permissions_staff_id FOREIGN KEY (staff_id) REFERENCES merchant_staff(id);
ALTER TABLE merchant_staff_permissions ADD CONSTRAINT fk_merchant_staff_permissions_granted_by FOREIGN KEY (granted_by) REFERENCES merchant_staff(id);

COMMENT ON TABLE merchant_staff_permissions IS '商家员工权限表';

-- ============================================
-- 10. 服务商员工表 (service_provider_staff)
-- ============================================

CREATE TABLE service_provider_staff (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    provider_id UUID NOT NULL,
    title VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_service_provider_staff_provider_id ON service_provider_staff(provider_id);
CREATE INDEX idx_service_provider_staff_status ON service_provider_staff(status);

-- 外键
ALTER TABLE service_provider_staff ADD CONSTRAINT fk_service_provider_staff_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE service_provider_staff ADD CONSTRAINT fk_service_provider_staff_provider_id FOREIGN KEY (provider_id) REFERENCES service_providers(id);

COMMENT ON TABLE service_provider_staff IS '服务商员工表';

-- ============================================
-- 11. 服务商员工权限表 (provider_staff_permissions)
-- ============================================

CREATE TABLE provider_staff_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    staff_id UUID NOT NULL,
    permission_code VARCHAR(50) NOT NULL,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    granted_by UUID NOT NULL
);

-- 索引
CREATE UNIQUE INDEX idx_provider_staff_permissions_staff_permission ON provider_staff_permissions(staff_id, permission_code);
CREATE INDEX idx_provider_staff_permissions_staff_id ON provider_staff_permissions(staff_id);
CREATE INDEX idx_provider_staff_permissions_permission_code ON provider_staff_permissions(permission_code);

-- 外键
ALTER TABLE provider_staff_permissions ADD CONSTRAINT fk_provider_staff_permissions_staff_id FOREIGN KEY (staff_id) REFERENCES service_provider_staff(id);
ALTER TABLE provider_staff_permissions ADD CONSTRAINT fk_provider_staff_permissions_granted_by FOREIGN KEY (granted_by) REFERENCES service_provider_staff(id);

COMMENT ON TABLE provider_staff_permissions IS '服务商员工权限表';

-- ============================================
-- 12. 积分账户表 (credit_accounts)
-- ============================================

CREATE TABLE credit_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL,
    owner_type VARCHAR(20) NOT NULL CHECK (owner_type IN ('ORG_MERCHANT', 'ORG_PROVIDER', 'USER_PERSONAL')),
    balance INT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    frozen_balance INT NOT NULL DEFAULT 0 CHECK (frozen_balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE UNIQUE INDEX idx_credit_accounts_owner_type ON credit_accounts(owner_id, owner_type);
CREATE INDEX idx_credit_accounts_owner_type_only ON credit_accounts(owner_type);

COMMENT ON TABLE credit_accounts IS '积分账户表';
COMMENT ON COLUMN credit_accounts.frozen_balance IS '冻结余额：商家/服务商为预扣费用，达人/员工为待确认收入+待打款金额';

-- ============================================
-- 13. 积分交易类型表 (transaction_types)
-- ============================================

CREATE TABLE transaction_types (
    code VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    account_types VARCHAR(20)[] NOT NULL,
    amount_direction VARCHAR(10) NOT NULL CHECK (amount_direction IN ('positive', 'negative')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deprecated_at TIMESTAMP
);

-- 初始化数据
INSERT INTO transaction_types (code, name, description, account_types, amount_direction) VALUES
('RECHARGE', '商家充值', '商家充值积分', ARRAY['ORG_MERCHANT'], 'positive'),
('TASK_INCOME', '任务收入', '达人任务审核通过，frozen_balance转入balance', ARRAY['USER_PERSONAL'], 'positive'),
('TASK_SUBMIT', '任务提交', '达人提交任务，预计收入冻结到frozen_balance', ARRAY['USER_PERSONAL'], 'positive'),
('STAFF_REFERRAL', '员工返佣', '员工邀请的达人完成任务时的返佣', ARRAY['USER_PERSONAL'], 'positive'),
('PROVIDER_INCOME', '服务商收入', '服务商完成任务审核通过获得的收入', ARRAY['ORG_PROVIDER'], 'positive'),
('TASK_PUBLISH', '发布任务', '商家发布任务扣除的积分（balance→frozen_balance）', ARRAY['ORG_MERCHANT'], 'negative'),
('TASK_ACCEPT', '接任务扣除', '达人接任务时从商家frozen_balance扣除', ARRAY['ORG_MERCHANT'], 'negative'),
('TASK_REJECT', '审核拒绝', '达人任务被拒绝，释放frozen_balance', ARRAY['USER_PERSONAL'], 'negative'),
('TASK_ESCALATE', '超时拒绝', '提交截止时间到，自动拒绝', ARRAY['ORG_MERCHANT'], 'positive'),
('TASK_REFUND', '任务退款', '任务关闭时，未完成名额费用退回balance', ARRAY['ORG_MERCHANT'], 'positive'),
('WITHDRAW', '提现', '用户申请提现（扣除balance）', ARRAY['USER_PERSONAL', 'ORG_MERCHANT', 'ORG_PROVIDER'], 'negative'),
('WITHDRAW_FREEZE', '提现冻结', '提现申请通过，从balance转入frozen_balance', ARRAY['USER_PERSONAL', 'ORG_MERCHANT', 'ORG_PROVIDER'], 'negative'),
('WITHDRAW_REFUND', '提现拒绝', '提现被拒绝，frozen_balance转回balance', ARRAY['USER_PERSONAL', 'ORG_MERCHANT', 'ORG_PROVIDER'], 'positive'),
('BONUS_GIFT', '系统赠送', '平台赠送的积分', ARRAY['USER_PERSONAL'], 'positive');

COMMENT ON TABLE transaction_types IS '积分交易类型表';

-- ============================================
-- 14. 积分流水表 (credit_transactions)
-- ============================================

CREATE TABLE credit_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    amount INT NOT NULL,
    balance_before INT NOT NULL,
    balance_after INT NOT NULL,
    transaction_group_id UUID,
    group_sequence INT,
    related_campaign_id UUID,
    related_task_id UUID,
    description VARCHAR(200),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_credit_transactions_account_id ON credit_transactions(account_id);
CREATE INDEX idx_credit_transactions_type ON credit_transactions(type);
CREATE INDEX idx_credit_transactions_group_id ON credit_transactions(transaction_group_id);
CREATE INDEX idx_credit_transactions_created_at ON credit_transactions(created_at);

-- 外键
ALTER TABLE credit_transactions ADD CONSTRAINT fk_credit_transactions_account_id FOREIGN KEY (account_id) REFERENCES credit_accounts(id);
ALTER TABLE credit_transactions ADD CONSTRAINT fk_credit_transactions_type FOREIGN KEY (type) REFERENCES transaction_types(code);
ALTER TABLE credit_transactions ADD CONSTRAINT fk_credit_transactions_campaign_id FOREIGN KEY (related_campaign_id) REFERENCES campaigns(id);
ALTER TABLE credit_transactions ADD CONSTRAINT fk_credit_transactions_task_id FOREIGN KEY (related_task_id) REFERENCES tasks(id);

COMMENT ON TABLE credit_transactions IS '积分流水表';

-- ============================================
-- 15. 提现表 (withdrawals)
-- ============================================

CREATE TABLE withdrawals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    amount INT NOT NULL,
    fee INT NOT NULL DEFAULT 0,
    actual_amount INT NOT NULL,
    method VARCHAR(20) NOT NULL CHECK (method IN ('ALIPAY', 'WECHAT', 'BANK')),
    account_info JSONB NOT NULL,
    account_info_hash VARCHAR(64),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'completed')),
    audit_note VARCHAR(200),
    audited_by UUID,
    audited_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_withdrawals_account_id ON withdrawals(account_id);
CREATE INDEX idx_withdrawals_status ON withdrawals(status);
CREATE INDEX idx_withdrawals_created_at ON withdrawals(created_at);

-- 外键
ALTER TABLE withdrawals ADD CONSTRAINT fk_withdrawals_account_id FOREIGN KEY (account_id) REFERENCES credit_accounts(id);
ALTER TABLE withdrawals ADD CONSTRAINT fk_withdrawals_audited_by FOREIGN KEY (audited_by) REFERENCES users(id);

COMMENT ON TABLE withdrawals IS '提现表';

-- ============================================
-- 16. 商家-服务商绑定表 (merchant_provider_bindings)
-- ============================================

CREATE TABLE merchant_provider_bindings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_id UUID NOT NULL,
    provider_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE UNIQUE INDEX idx_merchant_provider_bindings_merchant_provider ON merchant_provider_bindings(merchant_id, provider_id);
CREATE INDEX idx_merchant_provider_bindings_merchant_id ON merchant_provider_bindings(merchant_id);
CREATE INDEX idx_merchant_provider_bindings_provider_id ON merchant_provider_bindings(provider_id);
CREATE INDEX idx_merchant_provider_bindings_status ON merchant_provider_bindings(status);

-- 外键
ALTER TABLE merchant_provider_bindings ADD CONSTRAINT fk_mpb_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(id);
ALTER TABLE merchant_provider_bindings ADD CONSTRAINT fk_mpb_provider_id FOREIGN KEY (provider_id) REFERENCES service_providers(id);

COMMENT ON TABLE merchant_provider_bindings IS '商家-服务商绑定表';

-- ============================================
-- 函数：自动更新 updated_at
-- ============================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为所有表添加updated_at触发器
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_invitation_codes_updated_at BEFORE UPDATE ON invitation_codes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_merchants_updated_at BEFORE UPDATE ON merchants FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_service_providers_updated_at BEFORE UPDATE ON service_providers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_creators_updated_at BEFORE UPDATE ON creators FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_campaigns_updated_at BEFORE UPDATE ON campaigns FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_merchant_staff_updated_at BEFORE UPDATE ON merchant_staff FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_service_provider_staff_updated_at BEFORE UPDATE ON service_provider_staff FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_credit_accounts_updated_at BEFORE UPDATE ON credit_accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_transaction_types_updated_at BEFORE UPDATE ON transaction_types FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_credit_transactions_updated_at BEFORE UPDATE ON credit_transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_withdrawals_updated_at BEFORE UPDATE ON withdrawals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_merchant_provider_bindings_updated_at BEFORE UPDATE ON merchant_provider_bindings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_merchant_staff_permissions_updated_at BEFORE UPDATE ON merchant_staff_permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_provider_staff_permissions_updated_at BEFORE UPDATE ON provider_staff_permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 完成提示
-- ============================================

DO $$
BEGIN
    RAISE NOTICE '数据库初始化完成！';
    RAISE NOTICE '已创建 % 个表', (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public');
    RAISE NOTICE '下一步：运行 002_seed_data.sql 插入测试数据';
END $$;
