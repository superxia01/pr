-- 修改 invitation_codes 表的 generator_id 字段类型
-- 从 uuid 改为 varchar(255) 以匹配 users.id 的格式（usr_xxx）

-- 1. 删除相关约束
ALTER TABLE invitation_codes DROP CONSTRAINT IF EXISTS check_generator_consistency;
DROP INDEX IF EXISTS idx_invitation_codes_generator_id;

-- 2. 修改字段类型
ALTER TABLE invitation_codes
ALTER COLUMN generator_id TYPE varchar(255);

-- 3. 重新创建索引
CREATE INDEX idx_invitation_codes_generator_id
ON invitation_codes(generator_id);

-- 4. 重新创建约束（调整为 varchar 兼容）
ALTER TABLE invitation_codes
ADD CONSTRAINT check_generator_consistency
CHECK (
    (generator_type::text = 'system'::text AND generator_id IS NULL) OR
    (generator_type::text = ANY (ARRAY['super_admin'::text, 'merchant_admin'::text, 'sp_admin'::text, 'sp_staff'::text]) AND generator_id IS NOT NULL AND generator_id <> '')
);
