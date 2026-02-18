-- 达人表 inviter_id 与 users.id 类型一致，支持邀请人ID为 usr_xxx 等字符串格式

ALTER TABLE creators DROP CONSTRAINT IF EXISTS fk_creators_inviter_id;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'creators' AND column_name = 'inviter_id'
    AND data_type = 'uuid'
  ) THEN
    ALTER TABLE creators ALTER COLUMN inviter_id TYPE VARCHAR(255) USING inviter_id::text;
  END IF;
END $$;

ALTER TABLE creators ADD CONSTRAINT fk_creators_inviter_id FOREIGN KEY (inviter_id) REFERENCES users(id);

COMMENT ON COLUMN creators.inviter_id IS '邀请人用户ID，与 users.id 一致（VARCHAR）';
