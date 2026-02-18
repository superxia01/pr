-- 邀请关系表
-- 用于记录用户之间的邀请关系（谁邀请了谁、什么角色、绑定哪个组织）

CREATE TABLE IF NOT EXISTS invitation_relationships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inviter_id VARCHAR(255) NOT NULL,
    inviter_role VARCHAR(50) NOT NULL,
    invitee_id VARCHAR(255) NOT NULL,
    invitee_role VARCHAR(50) NOT NULL,
    organization_id UUID,
    organization_type VARCHAR(50),
    invitation_code VARCHAR(30) NOT NULL,
    invited_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'cancelled')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_invitation_relationships_inviter ON invitation_relationships(inviter_id);
CREATE INDEX IF NOT EXISTS idx_invitation_relationships_invitee ON invitation_relationships(invitee_id);
CREATE INDEX IF NOT EXISTS idx_invitation_relationships_code ON invitation_relationships(invitation_code);
CREATE INDEX IF NOT EXISTS idx_invitation_relationships_org ON invitation_relationships(organization_id);

-- 外键约束（可选，depending on whether users table uses these IDs）
-- ALTER TABLE invitation_relationships
--     ADD CONSTRAINT fk_invitation_relationships_inviter
--     FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE;
-- ALTER TABLE invitation_relationships
--     ADD CONSTRAINT fk_invitation_relationships_invitee
--     FOREIGN KEY (invitee_id) REFERENCES users(id) ON DELETE CASCADE;

-- 注释
COMMENT ON TABLE invitation_relationships IS '邀请关系记录表';
COMMENT ON COLUMN invitation_relationships.inviter_id IS '邀请人ID（users.id）';
COMMENT ON COLUMN invitation_relationships.inviter_role IS '邀请人角色（如：SERVICE_PROVIDER_ADMIN, MERCHANT_ADMIN, CREATOR）';
COMMENT ON COLUMN invitation_relationships.invitee_id IS '被邀请人ID（users.id）';
COMMENT ON COLUMN invitation_relationships.invitee_role IS '被邀请人角色';
COMMENT ON COLUMN invitation_relationships.organization_id IS '绑定的组织ID（可选）';
COMMENT ON COLUMN invitation_relationships.organization_type IS '绑定的组织类型（service_provider/merchant，可选）';
COMMENT ON COLUMN invitation_relationships.invitation_code IS '使用的邀请码';
COMMENT ON COLUMN invitation_relationships.invited_at IS '邀请成功时间（使用邀请码的时间）';
COMMENT ON COLUMN invitation_relationships.status IS '状态：active-有效, cancelled-已取消';
