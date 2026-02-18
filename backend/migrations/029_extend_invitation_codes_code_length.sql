-- 扩展 invitation_codes 表的 code 字段长度
-- 从 varchar(30) 扩展到 varchar(100) 以支持包含完整组织ID的邀请码格式

-- 邀请码格式：INV-{用户ID后8位}-{组织ID (36字符)}-{角色代码}
-- 总长度：4 + 1 + 8 + 1 + 36 + 1 + 10 = 61字符
-- 留有余量，设置为100

ALTER TABLE invitation_codes
ALTER COLUMN code TYPE varchar(100);
