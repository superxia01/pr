-- 更新超级管理员角色
-- 执行方式：psql -h localhost -U nexus_user -d pr_business_db -f update_super_admin.sql

-- 更新超级管理员
UPDATE users
SET
    roles = '["super_admin","admin","merchant_admin","provider_admin","creator"]'::jsonb,
    current_role = 'super_admin'
WHERE auth_center_user_id = 'oZh_a67J99sgfrHFX5pRPcXr0uQA';

-- 验证更新结果
SELECT
    id,
    auth_center_user_id,
    nickname,
    roles,
    current_role,
    status,
    created_at
FROM users
WHERE auth_center_user_id = 'oZh_a67J99sgfrHFX5pRPcXr0uQA';
