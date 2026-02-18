-- 邀请码表 generator_id 与 users.id 类型一致
-- users.id 为 VARCHAR(255)，generator_id 原为 UUID，现改为 VARCHAR(255) 以支持 usr_xxx 等格式

-- 1. 删除外键（若存在；001 中 UUID 与 users.id 类型不一致可能导致外键未创建成功）
ALTER TABLE invitation_codes DROP CONSTRAINT IF EXISTS fk_invitation_codes_generator_id;

-- 2. 仅当列为 UUID 时改为 VARCHAR(255)（若已是 varchar 则跳过）
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'invitation_codes' AND column_name = 'generator_id'
    AND data_type = 'uuid'
  ) THEN
    ALTER TABLE invitation_codes ALTER COLUMN generator_id TYPE VARCHAR(255) USING generator_id::text;
  END IF;
END $$;

-- 3. 允许 system 类型邀请码的 generator_id 为空（与 check_generator_consistency 一致）
ALTER TABLE invitation_codes ALTER COLUMN generator_id DROP NOT NULL;

-- 4. 重新添加外键
ALTER TABLE invitation_codes ADD CONSTRAINT fk_invitation_codes_generator_id
  FOREIGN KEY (generator_id) REFERENCES users(id);

COMMENT ON COLUMN invitation_codes.generator_id IS '生成者用户ID，与 users.id 一致（VARCHAR）';
