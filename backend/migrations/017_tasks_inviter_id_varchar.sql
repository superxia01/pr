-- 任务表 inviter_id 与 users.id 类型一致（VARCHAR）

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_inviter_id;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'tasks' AND column_name = 'inviter_id'
    AND data_type = 'uuid'
  ) THEN
    ALTER TABLE tasks ALTER COLUMN inviter_id TYPE VARCHAR(255) USING inviter_id::text;
  END IF;
END $$;

ALTER TABLE tasks ADD CONSTRAINT fk_tasks_inviter_id FOREIGN KEY (inviter_id) REFERENCES users(id);

COMMENT ON COLUMN tasks.inviter_id IS '邀请人用户ID，与 users.id 一致（VARCHAR）';
