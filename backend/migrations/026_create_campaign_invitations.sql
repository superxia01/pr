-- 活动邀请码表
-- 用于营销活动的邀请码管理

CREATE TABLE IF NOT EXISTS campaign_invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    code VARCHAR(30) NOT NULL UNIQUE,
    max_uses INT NOT NULL DEFAULT 1,
    use_count INT NOT NULL DEFAULT 0,
    expires_at TIMESTAMP NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'used')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_campaign_invitations_campaign_id ON campaign_invitations(campaign_id);
CREATE INDEX IF NOT EXISTS idx_campaign_invitations_code ON campaign_invitations(code);
CREATE INDEX IF NOT EXISTS idx_campaign_invitations_status ON campaign_invitations(status);

-- 注释
COMMENT ON TABLE campaign_invitations IS '活动邀请码表';
COMMENT ON COLUMN campaign_invitations.campaign_id IS '关联的营销活动';
COMMENT ON COLUMN campaign_invitations.code IS '邀请码（格式：CAMP-{活动ID后6位}-{随机4位}）';
COMMENT ON COLUMN campaign_invitations.max_uses IS '最大使用次数';
COMMENT ON COLUMN campaign_invitations.use_count IS '已使用次数';
COMMENT ON COLUMN campaign_invitations.expires_at IS '过期时间（NULL表示永不过期）';
COMMENT ON COLUMN campaign_invitations.status IS '状态：active-有效，inactive-失效，used-已使用';
