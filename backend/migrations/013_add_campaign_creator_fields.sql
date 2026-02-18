-- 添加活动创建者字段
-- 用于区分商家和服务商创建的活动，并控制权限

-- 添加 created_by 字段（创建者ID）
ALTER TABLE campaigns ADD COLUMN created_by VARCHAR(255) NOT NULL DEFAULT '';

-- 添加 creator_type 字段（创建者类型：MERCHANT_ADMIN 或 SERVICE_PROVIDER_ADMIN）
ALTER TABLE campaigns ADD COLUMN creator_type VARCHAR(50) NOT NULL DEFAULT '';

-- 添加索引
CREATE INDEX idx_campaigns_creator_type ON campaigns(creator_type);

-- 添加注释
COMMENT ON COLUMN campaigns.created_by IS '活动创建者ID（auth_center_user_id）';
COMMENT ON COLUMN campaigns.creator_type IS '创建者类型：MERCHANT_ADMIN（商家管理员）或 SERVICE_PROVIDER_ADMIN（服务商管理员）';
