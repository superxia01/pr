-- PR Business 测试数据脚本
-- 版本: v1.0
-- 日期: 2026-02-02
-- 说明: 插入测试用户、邀请码、组织等数据

-- ============================================
-- 1. 创建超级管理员
-- ============================================

-- 创建超级管理员用户
INSERT INTO users (id, auth_center_user_id, nickname, roles, current_role, status)
VALUES (
    'usr_admin_super_001',
    '550e8400-e29b-41d4-a716-446655440000',
    '超级管理员',
    ARRAY['SUPER_ADMIN']::jsonb,
    'SUPER_ADMIN',
    'active'
);

-- 创建超级管理员固定邀请码
INSERT INTO invitation_codes (id, code, type, target_role, generator_id, generator_type, is_active, is_one_time, max_uses, use_count)
VALUES (
    uuid_generate_v4(),
    'ADMIN-MASTER',
    'ADMIN-MASTER',
    'SUPER_ADMIN',
    NULL,
    'system',
    true,
    false,
    NULL,
    0
);

-- 创建服务商管理员固定邀请码
INSERT INTO invitation_codes (id, code, type, target_role, generator_id, generator_type, is_active, is_one_time, max_uses, use_count)
VALUES (
    uuid_generate_v4(),
    'SP-ADMIN',
    'SP-ADMIN',
    'SERVICE_PROVIDER_ADMIN',
    NULL,
    'system',
    true,
    false,
    NULL,
    0
);

-- ============================================
-- 2. 创建测试服务商
-- ============================================

-- 创建服务商用户
INSERT INTO users (id, auth_center_user_id, nickname, roles, current_role, status)
VALUES (
    'usr_provider_admin_001',
    '660e8400-e29b-41d4-a716-446655440001',
    '测试服务商',
    ARRAY['SERVICE_PROVIDER_ADMIN']::jsonb,
    'SERVICE_PROVIDER_ADMIN',
    'active'
);

-- 创建服务商组织
INSERT INTO service_providers (id, admin_id, user_id, name, description, status)
VALUES (
    uuid_generate_v4(),
    'usr_provider_admin_001',
    'usr_provider_admin_001',
    '星云营销服务商',
    '专业的达人营销服务提供商',
    'active'
);

-- 获取服务商ID（用于后续生成邀请码）
DO $$
DECLARE
    v_provider_id UUID;
BEGIN
    SELECT id INTO v_provider_id FROM service_providers WHERE admin_id = 'usr_provider_admin_001' LIMIT 1;

    -- 生成服务商员工邀请码
    INSERT INTO invitation_codes (id, code, type, target_role, generator_id, generator_type, organization_id, organization_type, is_active, is_one_time, max_uses, use_count)
    VALUES (
        uuid_generate_v4(),
        'SPSTAFF-' || SUBSTRING(v_provider_id::text FROM LENGTH(v_provider_id::text) - 5 FOR 6),
        'SPSTAFF-{服务商ID后6位}',
        'SERVICE_PROVIDER_STAFF',
        'usr_provider_admin_001',
        'sp_admin',
        v_provider_id,
        'service_provider',
        true,
        false,
        NULL,
        0
    );

    -- 生成商家邀请码（核心功能）
    INSERT INTO invitation_codes (id, code, type, target_role, generator_id, generator_type, organization_id, organization_type, is_active, is_one_time, max_uses, use_count)
    VALUES (
        uuid_generate_v4(),
        'MERCHANT-' || SUBSTRING(v_provider_id::text FROM LENGTH(v_provider_id::text) - 5 FOR 6),
        'MERCHANT-{服务商ID后6位}',
        'MERCHANT_ADMIN',
        'usr_provider_admin_001',
        'sp_admin',
        v_provider_id,
        'service_provider',
        true,
        false,
        NULL,
        0
    );

    RAISE NOTICE '测试服务商创建完成，ID: %', v_provider_id;
END $$;

-- ============================================
-- 3. 创建测试商家
-- ============================================

-- 创建商家用户（使用服务商邀请码的场景）
INSERT INTO users (id, auth_center_user_id, nickname, roles, current_role, invitation_code_id, status)
VALUES (
    'usr_merchant_admin_001',
    '770e8400-e29b-41d4-a716-446655440002',
    '测试商家',
    ARRAY['MERCHANT_ADMIN']::jsonb,
    'MERCHANT_ADMIN',
    NULL,
    'active'
);

-- 获取服务商ID
DO $$
DECLARE
    v_provider_id UUID;
    v_merchant_code VARCHAR(30);
BEGIN
    SELECT id INTO v_provider_id FROM service_providers WHERE admin_id = 'usr_provider_admin_001' LIMIT 1;

    -- 生成商家邀请码
    v_merchant_code := 'MERCHANT-' || SUBSTRING(v_provider_id::text FROM LENGTH(v_provider_id::text) - 5 FOR 6);

    -- 更新用户的invitation_code_id（模拟使用邀请码）
    UPDATE invitation_codes
    SET use_count = use_count + 1
    WHERE code = v_merchant_code;

    -- 获取邀请码ID
    UPDATE users
    SET invitation_code_id = (SELECT id FROM invitation_codes WHERE code = v_merchant_code)
    WHERE id = 'usr_merchant_admin_001';

    -- 创建商家组织
    INSERT INTO merchants (id, admin_id, provider_id, user_id, name, description, industry, status)
    VALUES (
        uuid_generate_v4(),
        'usr_merchant_admin_001',
        v_provider_id,
        'usr_merchant_admin_001',
        '美妆品牌',
        '专注于美妆产品的达人营销推广',
        '美妆护肤',
        'active'
    );

    -- 生成商家员工邀请码
    INSERT INTO invitation_codes (id, code, type, target_role, generator_id, generator_type, is_active, is_one_time, max_uses, use_count)
    VALUES (
        uuid_generate_v4(),
        'MSTAFF-' || SUBSTRING(encode(gen_random_bytes(16), 'hex'), 1, 6),
        'MSTAFF-{商家ID后6位}',
        'MERCHANT_STAFF',
        'usr_merchant_admin_001',
        'merchant_admin',
        true,
        false,
        NULL,
        0
    );

    RAISE NOTICE '测试商家创建完成，已绑定到服务商';
END $$;

-- ============================================
-- 4. 创建测试达人
-- ============================================

-- 创建达人用户
INSERT INTO users (id, auth_center_user_id, nickname, roles, current_role, status)
VALUES (
    'usr_creator_001',
    '880e8400-e29b-41d4-a716-446655440003',
    '达人小李',
    ARRAY['CREATOR']::jsonb,
    'CREATOR',
    'active'
), (
    'usr_creator_002',
    '880e8400-e29b-41d4-a716-446655440004',
    '达人小张',
    ARRAY['CREATOR']::jsonb,
    'CREATOR',
    'active'
);

-- 创建达人记录
INSERT INTO creators (id, user_id, is_primary, level, followers_count, status)
VALUES
    (uuid_generate_v4(), 'usr_creator_001', true, 'KOC', 5000, 'active'),
    (uuid_generate_v4(), 'usr_creator_002', true, 'INF', 50000, 'active');

-- 创建员工邀请码（用于邀请达人）
INSERT INTO users (id, auth_center_user_id, nickname, roles, current_role, status)
VALUES (
    'usr_sp_staff_001',
    '990e8400-e29b-41d4-a716-446655440005',
    '服务商员工小王',
    ARRAY['SERVICE_PROVIDER_STAFF', 'CREATOR']::jsonb,
    'SERVICE_PROVIDER_STAFF',
    'active'
);

-- 生成员工邀请码
INSERT INTO invitation_codes (id, code, type, target_role, generator_id, generator_type, organization_id, organization_type, is_active, is_one_time, max_uses, use_count)
VALUES (
    uuid_generate_v4(),
    'CREATOR-' || SUBSTRING(encode(gen_random_bytes(16), 'hex'), 1, 6),
    'CREATOR-{员工ID后6位}',
    'CREATOR',
    'usr_sp_staff_001',
    'sp_staff',
    (SELECT id FROM service_providers WHERE admin_id = 'usr_provider_admin_001'),
    'service_provider',
    true,
    false,
    NULL,
    0
);

-- ============================================
-- 5. 创建积分账户
-- ============================================

-- 商家组织账户
DO $$
DECLARE
    v_merchant_id UUID;
BEGIN
    SELECT id INTO v_merchant_id FROM merchants WHERE admin_id = 'usr_merchant_admin_001' LIMIT 1;

    INSERT INTO credit_accounts (id, owner_id, owner_type, balance)
    VALUES (
        uuid_generate_v4(),
        v_merchant_id,
        'ORG_MERCHANT',
        0  -- 初始余额0
    );

    RAISE NOTICE '商家积分账户创建完成';
END $$;

-- 服务商组织账户
DO $$
DECLARE
    v_provider_id UUID;
BEGIN
    SELECT id INTO v_provider_id FROM service_providers WHERE admin_id = 'usr_provider_admin_001' LIMIT 1;

    INSERT INTO credit_accounts (id, owner_id, owner_type, balance)
    VALUES (
        uuid_generate_v4(),
        v_provider_id,
        'ORG_PROVIDER',
        0  -- 初始余额0
    );

    RAISE NOTICE '服务商积分账户创建完成';
END $$;

-- 达人个人账户
DO $$
DECLARE
    v_creator_id UUID;
BEGIN
    SELECT id INTO v_creator_id FROM creators WHERE user_id = 'usr_creator_001' LIMIT 1;

    INSERT INTO credit_accounts (id, owner_id, owner_type, balance)
    VALUES (
        uuid_generate_v4(),
        v_creator_id,
        'USER_PERSONAL',
        0  -- 初始余额0
    );

    RAISE NOTICE '达人积分账户创建完成';
END $$;

-- ============================================
-- 6. 创建测试营销活动
-- ============================================

DO $$
DECLARE
    v_merchant_id UUID;
    v_provider_id UUID;
BEGIN
    SELECT id INTO v_merchant_id FROM merchants WHERE admin_id = 'usr_merchant_admin_001' LIMIT 1;
    SELECT id INTO v_provider_id FROM service_providers WHERE admin_id = 'usr_provider_admin_001' LIMIT 1;

    -- 创建营销活动
    INSERT INTO campaigns (
        id,
        merchant_id,
        provider_id,
        title,
        requirements,
        platforms,
        task_amount,
        campaign_amount,
        creator_amount,
        staff_referral_amount,
        provider_amount,
        quota,
        task_deadline,
        submission_deadline,
        status
    ) VALUES (
        uuid_generate_v4(),
        v_merchant_id,
        v_provider_id,
        '美妆产品体验推广',
        '体验我们的新款护肤套装，分享使用心得到小红书。',
        ARRAY['xiaohongshu']::jsonb,
        100,  -- 任务金额：100元
        1000,  -- 活动总费用：1000元（10人×100元）
        80,   -- 达人收入：80元
        10,   -- 员工返佣：10元
        10,   -- 服务商收入：10元
        10,   -- 活动名额：10人
        NOW() + INTERVAL '30 days',   -- 接任务截止：30天后
        NOW() + INTERVAL '37 days',   -- 提交截止：37天后
        'OPEN'  -- 直接发布为开放状态
    );

    -- 为营销活动创建10个任务名额
    INSERT INTO tasks (id, campaign_id, task_slot_number, status, version)
    SELECT
        uuid_generate_v4(),
        c.id,
        generate_series(1, 10),
        'OPEN',
        0
    FROM campaigns c
    WHERE c.title = '美妆产品体验推广'
    ORDER BY c.id
    LIMIT 1;

    RAISE NOTICE '测试营销活动创建完成，包含10个任务名额';
END $$;

-- ============================================
-- 完成提示
-- ============================================

DO $$
BEGIN
    RAISE NOTICE '===========================================';
    RAISE NOTICE '测试数据插入完成！';
    RAISE NOTICE '===========================================';
    RAISE NOTICE '已创建：';
    RAISE NOTICE '  - 1个超级管理员';
    RAISE NOTICE '  - 1个服务商（含员工）';
    RAISE NOTICE '  - 1个商家（含员工）';
    RAISE NOTICE '  - 2个达人';
    RAISE NOTICE '  - 1个营销活动（10个任务名额）';
    RAISE NOTICE '  - 所有积分账户';
    RAISE NOTICE '固定邀请码：';
    RAISE NOTICE '  - ADMIN-MASTER（超级管理员）';
    RAISE NOTICE '  - SP-ADMIN（服务商管理员）';
    RAISE NOTICE '  - MERCHANT-{xxx}（商家邀请码）';
    RAISE NOTICE '  - MSTAFF-{xxx}（商家员工邀请码）';
    RAISE NOTICE '  - SPSTAFF-{xxx}（服务商员工邀请码）';
    RAISE NOTICE '  - CREATOR-{xxx}（达人邀请码）';
    RAISE NOTICE '';
    RAISE NOTICE '测试账号：';
    RAISE NOTICE '  超级管理员：usr_admin_super_001';
    RAISE NOTICE '  服务商管理员：usr_provider_admin_001';
    RAISE NOTICE '  商家管理员：usr_merchant_admin_001';
    RAISE NOTICE '  达人：usr_creator_001, usr_creator_002';
    RAISE NOTICE '===========================================';
END $$;
