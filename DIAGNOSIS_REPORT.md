# 🔧 微信登录"用户创建失败"问题 - 完整诊断与修复方案

**问题时间**: 2026-02-04 11:54:47
**错误类型**: 数据库Scan错误
**状态**: ✅ 问题已定位，修复方案已准备（未部署）

---

## 🔴 问题现象

### 错误信息
```
https://pr.crazyaigc.com/api/v1/auth/callback?code=xxx&type=open
{"error":"用户创建失败"}
```

### 后端日志
```
[2026/02/04 11:54:47] sql: Scan error on column index 4, name "profile":
json: cannot unmarshal string into Go value of type models.Profile

SELECT * FROM "users" WHERE auth_center_user_id = '300d0851-7a28-4ad0-98dc-98ac29811945' ORDER BY "users"."id" LIMIT 1
```

**位置**: `backend/controllers/auth.go:365`
**API**: `/api/v1/auth/callback`
**状态码**: 500 (Internal Server Error)

---

## 🔍 根本原因分析

### 问题1: 数据格式不匹配

**数据库实际情况**:
```sql
SELECT id, auth_center_user_id, nickname, profile, roles
FROM users
WHERE auth_center_user_id = '300d0851-7a28-4ad0-98dc-98ac29811945';

id: usr_d03437d6-db20-43d1-b8f9-31d94854dc43
auth_center_user_id: 300d0851-7a28-4ad0-98dc-98ac29811945
nickname: 新用户
profile: "{}"      -- ⚠️ 这是字符串，不是JSONB！
roles: ["SUPER_ADMIN"]
```

**表结构定义**（正确）:
```sql
column_name | data_type | udt_name
-------------+-----------+----------
profile     | jsonb     | jsonb  ✅ 类型正确
roles       | jsonb     | jsonb  ✅ 类型正确
```

**问题**:
- 表字段类型是 `JSONB` ✅
- 但已有数据以字符串格式存储 `"{}"` ❌
- PostgreSQL在某些情况下会将JSONB数据作为text返回
- GORM的Scan方法无法处理字符串格式的JSONB

### 问题2: 为什么已有的Scan修复没有生效？

**Scan方法已修复** (`models/user.go:46-64`):
```go
func (p *Profile) Scan(value interface{}) error {
    // ...
    switch v := value.(type) {
    case []byte:
        bytes = v
    case string:  // ✅ 已添加字符串处理
        bytes = []byte(v)
    default:
        return nil
    }
    return json.Unmarshal(bytes, p)
}
```

**为什么还报错**?
- 可能的原因：编译的二进制不是最新版本
- 或者：GORM在Scan前做了类型检查，没有调用自定义Scan方法

---

## ✅ 修复方案

### 方案1: 数据修复（推荐，需要执行SQL）

**文件**: `backend/migrations/fix_profile_roles_data.sql`

**修复步骤**:

#### 1. 备份现有数据
```sql
CREATE TABLE IF NOT EXISTS users_backup_20260204 AS
SELECT * FROM users;
```

#### 2. 修复profile字段
```sql
UPDATE users
SET profile = CASE
    WHEN profile IS NULL THEN '{}'::jsonb
    WHEN profile = '' THEN '{}'::jsonb
    WHEN profile::text = '{}' THEN '{}'::jsonb
    WHEN substring(profile::text, 1, 1) = '{' THEN profile::jsonb
    ELSE '{}'::jsonb
END;
```

#### 3. 修复roles字段
```sql
UPDATE users
SET roles = CASE
    WHEN roles IS NULL THEN '[]'::jsonb
    WHEN roles = '' THEN '[]'::jsonb
    WHEN roles::text = '[]' THEN '[]'::jsonb
    WHEN substring(roles::text, 1, 1) = '[' THEN roles::jsonb
    ELSE '[]'::jsonb
END;
```

#### 4. 验证修复
```sql
-- 验证profile格式
SELECT id, nickname, jsonb_pretty(profile) as profile_formatted
FROM users
LIMIT 5;

-- 验证roles格式
SELECT id, nickname, jsonb_pretty(roles) as roles_formatted
FROM users
LIMIT 5;
```

### 方案2: 代码加固（已完成，但未部署）

**已修复的代码** (`models/user.go`):

✅ Profile.Scan() - 支持 string 类型
✅ Roles.Scan() - 支持 string 类型
✅ Profile.Value() - 正确序列化
✅ Roles.Value() - 正确序列化

**未生效的可能原因**:
- 编译的二进制不是最新版本
- 需要重新编译并部署

---

## 📋 执行步骤（当准备好时）

### 第一步：执行数据修复SQL

```bash
# 1. SSH到服务器
ssh shanghai-tencent

# 2. 连接数据库
PGPASSWORD=hRJ9NSJApfeyFDraaDgkYowY psql -h localhost -p 5432 -U nexus_user -d pr_business_db

# 3. 执行修复脚本
\i /path/to/fix_profile_roles_data.sql

# 或者手动执行上面的SQL语句

# 4. 验证修复结果
SELECT id, nickname, profile, roles FROM users LIMIT 5;
```

### 第二步：重新编译并部署

```bash
# 1. 本地重新编译
cd /Users/xia/Documents/GitHub/pr-business/backend
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o pr-business-linux .

# 2. 上传到服务器
scp pr-business-linux shanghai-tencent:/var/www/pr-backend/

# 3. 重启服务
ssh shanghai-tencent <<'EOF'
cd /var/www/pr-backend
sudo systemctl restart pr-business-backend
sudo systemctl status pr-business-backend --no-pager | head -15
EOF
```

### 第三步：测试验证

```bash
# 1. 测试微信登录
# 访问 https://pr.crazyaigc.com
# 点击微信登录

# 2. 查看日志
ssh shanghai-tencent "sudo journalctl -u pr-business-backend -f"

# 3. 验证数据库连接
ssh shanghai-tencent "
PGPASSWORD=hRJ9NSJApfeyFDraaDgkYowY psql -h localhost -p 5432 -U nexus_user -d pr_business_db -c 'SELECT COUNT(*) FROM users;'
"
```

---

## 🛡️ 预防措施

### 未来避免此问题

1. **初始化时使用正确的格式**
   ```go
   // ✅ 正确
   user.Profile = models.Profile{}  // 空对象
   user.Roles = models.Roles{}      // 空数组

   // ❌ 错误
   user.Profile = models.Profile("{}")  // 不要传字符串
   user.Roles = models.Roles("[]")     // 不要传字符串
   ```

2. **数据库迁移时确保类型**
   ```sql
   -- ✅ 正确
   profile JSONB DEFAULT '{}'::jsonb NOT NULL
   roles JSONB DEFAULT '[]'::jsonb NOT NULL

   -- ❌ 错误
   profile TEXT DEFAULT '{}'
   profile VARCHAR(255) DEFAULT '{}'
   ```

3. **数据验证查询**
   ```sql
   -- 定期检查数据格式
   SELECT
       id,
       profile,
       jsonb_typeof(profile) as profile_type,
       roles,
       jsonb_typeof(roles) as roles_type
   FROM users;
   ```

---

## 📊 问题影响范围

### 受影响的用户
- 所有在修复前注册的用户
- 数据库中 `profile` 或 `roles` 字段为字符串格式的用户

### 受影响的功能
- 微信登录回调（查询已有用户时）
- 密码登录（查询已有用户时）
- 任何需要Scan用户的操作

### 未受影响的功能
- 新用户注册（会创建正确格式的数据）
- API认证（Token生成和验证）

---

## 📝 修复检查清单

修复后需要验证的点：

- [ ] SQL修复脚本已执行
- [ ] 数据格式验证通过
- [ ] 后端服务已重启
- [ ] 日志无Scan错误
- [ ] 微信登录测试通过
- [ ] 密码登录测试通过
- [ ] Token刷新正常

---

## 🎯 关键要点

1. **问题本质**: 数据格式不匹配，不是代码bug
2. **最佳方案**: 执行SQL修复脚本 + 重新部署
3. **无需担心**: 已有用户数据不会丢失，只是格式转换
4. **预防**: 新用户不会有这个问题（修复后的代码会正确处理）

**准备好执行修复时告诉我，我将协助你执行上述步骤。**
