-- 任务邀请功能

-- 1. 任务邀请码表
CREATE TABLE IF NOT EXISTS task_invitation_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) NOT NULL UNIQUE,
    campaign_id UUID NOT NULL,
    generator_id VARCHAR(255) NOT NULL,
    generator_type VARCHAR(50) NOT NULL CHECK (generator_type IN ('merchant_admin', 'provider_admin', 'merchant_staff', 'provider_staff')),
    max_uses INT DEFAULT 100,
    use_count INT DEFAULT 0,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 2. 任务邀请记录表
CREATE TABLE IF NOT EXISTS task_invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invitation_code_id UUID NOT NULL REFERENCES task_invitation_codes(id) ON DELETE CASCADE,
    creator_id VARCHAR(255),
    task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'accepted' CHECK (status IN ('accepted', 'completed', 'expired')),
    accepted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_task_invitation_codes_campaign ON task_invitation_codes(campaign_id);
CREATE INDEX IF NOT EXISTS idx_task_invitation_codes_code ON task_invitation_codes(code);
CREATE INDEX IF NOT EXISTS idx_task_invitations_invitation ON task_invitations(invitation_code_id);
CREATE INDEX IF NOT EXISTS idx_task_invitations_creator ON task_invitations(creator_id);
CREATE INDEX IF NOT EXISTS idx_task_invitations_task ON task_invitations(task_id);

-- 注释
COMMENT ON TABLE task_invitation_codes IS '任务邀请码表';
COMMENT ON TABLE task_invitations IS '任务邀请记录表';
COMMENT ON COLUMN task_invitation_codes.code IS '邀请码（6位随机码）';
COMMENT ON COLUMN task_invitation_codes.generator_type IS '生成者类型';
COMMENT ON COLUMN task_invitations.status IS '状态: accepted-已接受, completed-已完成, expired-已过期';
