# PR Business 微信登录流程文档

## 概述

本文档描述 PR Business 系统如何通过 auth-center 实现微信登录，适用于其他业务系统参考实现。

---

## 架构设计

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   前端 Vue  │◄───────►│  后端 Gin    │◄───────►│ auth-center │
│  pr.crazyaigc.com  │         │ pr-backend   │         │os.crazyaigc.com│
└─────────────┘         └──────────────┘         └─────────────┘
```

**关键原则**：
1. ✅ 前端只能与自己的后端通信，**不能直接调用 auth-center**
2. ✅ auth-center 的 token **只用于**调用 auth-center API，不保存到数据库
3. ✅ 后端生成自己的 JWT token 用于业务系统认证
4. ✅ 用户资料（昵称、头像）在登录响应中直接返回，无需额外调用

---

## 完整登录流程

### 步骤1：用户发起微信登录

**前端操作**：
```typescript
// 用户点击"微信登录"按钮
window.location.href = `${API_BASE_URL}/api/v1/auth/wechat/login`
```

**后端接口**：`GET /api/v1/auth/wechat/login`

```go
// 后端代码
func (ctrl *AuthController) WeChatLoginRedirect(c *gin.Context) {
    // 重定向到 auth-center 的微信登录页面
    authCenterURL := fmt.Sprintf(
        "%s/api/auth/wechat/login?callbackUrl=%s",
        ctrl.cfg.AuthCenterURL,
        url.QueryEscape(ctrl.cfg.FrontendURL+"/login"), // ✅ 回调到前端！
    )
    c.Redirect(http.StatusFound, authCenterURL)
}
```

**请求示例**：
```
GET https://os.crazyaigc.com/api/auth/wechat/login?callbackUrl=https%3A%2F%2Fpr.crazyaigc.com%2Flogin
```

---

### 步骤2：auth-center 处理微信授权

1. auth-center 引导用户完成微信授权
2. 授权成功后，auth-center **重定向回前端登录页**（callbackUrl）

---

### 步骤3：前端接收回调

前端登录页 `https://pr.crazyaigc.com/login` 接收 URL 参数：

**情况A：微信内登录**（直接有 token）
```
URL: https://pr.crazyaigc.com/login?token=xxx&userId=yyy
```

**情况B：PC扫码登录**（需要用 code 换 token）
```
URL: https://pr.crazyaigc.com/login?code=xxx&type=open
```

---

### 步骤4A：微信内登录流程（推荐）

**前端代码**：
```typescript
const token = searchParams.get('token')
const userId = searchParams.get('userId')

if (token && userId) {
  // ✅ 使用 auth-center token 调用后端获取完整用户信息
  fetch(`${API_BASE_URL}/api/v1/user/me`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
    .then(res => res.json())
    .then(data => {
      // 后端返回的 LoginResponse 包含：
      // - accessToken: 后端生成的 JWT token
      // - userId: 用户ID
      // - nickname: 微信昵称
      // - avatarUrl: 微信头像
      // - roles: 角色列表
      // - currentRole: 当前角色

      login(data.accessToken, data.refreshToken, data)
      navigate('/dashboard')
    })
}
```

**后端接口**：`GET /api/v1/user/me`

```go
func (ctrl *AuthController) GetCurrentUser(c *gin.Context) {
    userID := c.GetString("userId")
    var user models.User
    ctrl.db.Where("id = ?", userID).First(&user)

    // ✅ 从请求头获取 auth-center token
    authHeader := c.GetHeader("Authorization")
    if strings.HasPrefix(authHeader, "Bearer ") {
        authCenterToken := strings.TrimPrefix(authHeader, "Bearer ")

        // ✅ 调用 auth-center 获取最新用户资料（昵称、头像）
        userInfo, _ := ctrl.callAuthCenterGetUserInfo(authCenterToken)
        if userInfo.Success {
            // 更新本地用户资料
            if userInfo.Data.Profile.Nickname != "" {
                user.Nickname = userInfo.Data.Profile.Nickname
            }
            if userInfo.Data.Profile.AvatarURL != "" {
                user.AvatarURL = userInfo.Data.Profile.AvatarURL
            }
            ctrl.db.Save(&user)
        }
    }

    // ✅ 返回完整用户信息
    c.JSON(http.StatusOK, gin.H{
        "id":                 user.ID,
        "nickname":           user.Nickname,      // ✅ 微信昵称
        "avatarUrl":          user.AvatarURL,     // ✅ 微信头像
        "roles":              convertRolesToUpperCase(user.Roles),
        "currentRole":        user.ActiveRole,
        // ...
    })
}
```

**auth-center API 调用**：
```go
func (ctrl *AuthController) callAuthCenterGetUserInfo(token string) (*AuthCenterUserInfo, error) {
    req, _ := http.NewRequest("GET", ctrl.cfg.AuthCenterURL+"/api/auth/user-info", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, _ := client.Do(req)
    defer resp.Body.Close()

    var userInfoResp AuthCenterUserInfo
    json.NewDecoder(resp.Body).Decode(&userInfoResp)
    return &userInfoResp, nil
}
```

---

### 步骤4B：PC扫码登录流程

**前端代码**：
```typescript
const code = searchParams.get('code')

if (code) {
  // ✅ 用 code 换取 token
  fetch(`${API_BASE_URL}/api/v1/auth/wechat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ authCode: code }),
  })
    .then(res => res.json())
    .then(data => {
      // 后端返回的 LoginResponse 包含：
      // - accessToken: 后端生成的 JWT token
      // - nickname: 微信昵称（从 auth-center 登录响应中获取）
      // - avatarUrl: 微信头像（从 auth-center 登录响应中获取）

      login(data.accessToken, data.refreshToken, data)
      navigate('/dashboard')
    })
}
```

**后端接口**：`POST /api/v1/auth/wechat`

```go
func (ctrl *AuthController) WeChatLogin(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)

    // 1. 调用 auth-center 用 code 换取用户信息
    authCenterResp, err := ctrl.callAuthCenterWechatLogin(req.AuthCode, "open")
    if !authCenterResp.Success {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "微信授权失败"})
        return
    }

    // 2. 创建或获取本地用户
    user, err := ctrl.findOrCreateUser(authCenterResp.Data.UserID)

    // 3. ✅ 直接从登录响应获取用户资料（无需再调用 user-info）
    if authCenterResp.Data.Profile.Nickname != "" {
        user.Nickname = authCenterResp.Data.Profile.Nickname
    }
    if authCenterResp.Data.Profile.AvatarURL != "" {
        user.AvatarURL = authCenterResp.Data.Profile.AvatarURL
    }
    ctrl.db.Save(&user)

    // 4. 生成后端自己的 JWT token
    accessToken, refreshToken := ctrl.generateTokens(user)

    // 5. ✅ 返回包含昵称和头像的登录响应
    c.JSON(http.StatusOK, LoginResponse{
        AccessToken:  accessToken,          // 后端 JWT token
        RefreshToken: refreshToken,
        ExpiresIn:    int64(ctrl.cfg.JWTAccessTokenExpire.Seconds()),
        UserID:       user.ID,
        Nickname:     user.Nickname,        // ✅ 微信昵称
        AvatarURL:    user.AvatarURL,       // ✅ 微信头像
        Roles:        convertRolesToUpperCase(user.Roles),
        CurrentRole:  user.ActiveRole,
    })
}
```

**auth-center 登录 API 调用**：
```go
func (ctrl *AuthController) callAuthCenterWechatLogin(code, loginType string) (*AuthCenterResponse, error) {
    reqBody := map[string]string{
        "code": code,
        "type": loginType,
    }
    jsonBody, _ := json.Marshal(reqBody)

    resp, err := http.Post(
        fmt.Sprintf("%s/api/auth/wechat/login", ctrl.cfg.AuthCenterURL),
        "application/json",
        bytes.NewBuffer(jsonBody),
    )
    defer resp.Body.Close()

    var result AuthCenterResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

**auth-center 响应结构**：
```json
{
  "success": true,
  "data": {
    "userId": "300d0851-7a28-4ad0-98dc-98ac29811945",
    "token": "",
    "unionId": "oZh_a67J99sgfrHFX5pRPcXr0uQA",
    "phoneNumber": "13777076463",
    "profile": {
      "nickname": "微信昵称",
      "avatarUrl": "https://xxx.jpg"
    }
  }
}
```

**⚠️ 注意**：auth-center 的 `/api/auth/wechat/login` 响应中 `token` 字段为空，但 `data.profile` 包含用户资料。

---

## 数据结构

### 后端结构体定义

```go
// auth-center 登录响应
type AuthCenterResponse struct {
    Success bool `json:"success"`
    Data    struct {
        UserID      string `json:"userId"`
        Token       string `json:"token"`        // ⚠️ 可能为空
        UnionID     string `json:"unionId"`
        PhoneNumber string `json:"phoneNumber"`
        Email       string `json:"email"`
        Profile     struct {
            Nickname  string `json:"nickname"`
            AvatarURL string `json:"avatarUrl"`
        } `json:"profile"`  // ✅ 关键：用户资料在这里
    } `json:"data"`
    Error string `json:"error,omitempty"`
}

// auth-center 用户信息响应
type AuthCenterUserInfo struct {
    Success bool `json:"success"`
    Data    struct {
        UserID  string `json:"userId"`
        Profile struct {
            Nickname  string `json:"nickname"`
            AvatarURL string `json:"avatarUrl"`
        } `json:"profile"`
    } `json:"data"`
}

// 返回给前端的登录响应
type LoginResponse struct {
    AccessToken  string   `json:"accessToken"`   // 后端 JWT token
    RefreshToken string   `json:"refreshToken"`
    ExpiresIn    int64    `json:"expiresIn"`
    UserID       string   `json:"userId"`
    Nickname     string   `json:"nickname"`      // ✅ 微信昵称
    AvatarURL    string   `json:"avatarUrl"`     // ✅ 微信头像
    Roles        []string `json:"roles"`         // 转大写：SUPER_ADMIN
    CurrentRole  string   `json:"currentRole"`   // 小写：super_admin
}
```

### 前端类型定义

```typescript
// 登录响应
interface LoginResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  userId: string
  nickname?: string       // ✅ 微信昵称
  avatarUrl?: string      // ✅ 微信头像
  roles: string[]
  currentRole: string
}

// 用户信息
interface User {
  id: string
  nickname: string
  avatarUrl: string
  profile: Record<string, any>
  roles: string[]
  currentRole: string
  // ...
}
```

---

## 关键配置

### 后端配置 (.env)

```bash
# auth-center 配置
AUTH_CENTER_URL=https://os.crazyaigc.com
AUTH_CENTER_REDIRECT_URI=https://pr.crazyaigc.com/api/v1/auth/callback  # ⚠️ 已废弃

# 前端配置
FRONTEND_URL=https://pr.crazyaigc.com

# JWT 配置（业务系统自己的 token）
JWT_SECRET=your-secret-key
JWT_ACCESS_TOKEN_EXPIRE=24h
JWT_REFRESH_TOKEN_EXPIRE=168h
```

---

## 最佳实践总结

### ✅ DO（推荐做法）

1. **auth-center 回调到前端**，不是后端
   ```go
   callbackUrl=https://pr.crazyaigc.com/login  // ✅ 前端
   ```

2. **直接从登录响应获取用户资料**，不要依赖 token 字段
   ```go
   // ✅ 从 data.profile 获取
   user.Nickname = authCenterResp.Data.Profile.Nickname
   user.AvatarURL = authCenterResp.Data.Profile.AvatarURL
   ```

3. **前端只与自己的后端通信**
   ```typescript
   // ✅ 正确：调用后端
   fetch('/api/v1/user/me', { headers: { 'Authorization': `Bearer ${authCenterToken}` }})

   // ❌ 错误：直接调用 auth-center
   fetch('https://os.crazyaigc.com/api/auth/user-info', ...)
   ```

4. **auth-center token 只用于 API 调用**，不保存到数据库
   ```go
   // ✅ 正确：用 token 调用 auth-center API
   userInfo := callAuthCenterGetUserInfo(authCenterToken)

   // ❌ 错误：保存 token 到数据库
   user.AuthCenterToken = token  // 不要这样做
   ```

5. **登录响应中包含完整的用户资料**
   ```go
   c.JSON(http.StatusOK, LoginResponse{
       Nickname:  user.Nickname,   // ✅ 包含昵称
       AvatarURL: user.AvatarURL,  // ✅ 包含头像
   })
   ```

### ❌ DON'T（避免的做法）

1. ❌ **不要**让 auth-center 回调到后端
2. ❌ **不要**保存 auth-center 的 token 到数据库
3. ❌ **不要**期望 auth-center 登录响应的 token 字段有值（它是空的）
4. ❌ **不要**前端直接调用 auth-center API
5. ❌ **不要**在数据库字段中存储 auth_center_token

---

## 流程图

```
用户点击微信登录
    │
    ▼
前端: GET /api/v1/auth/wechat/login
    │
    ▼
后端: 302 重定向到 auth-center
    │
    ▼
auth-center: 微信授权
    │
    ▼
auth-center: 302 重定向回前端（带 token 或 code）
    │
    ├──────────────────────────────────────┐
    ▼                                      ▼
微信内登录（有 token）              PC扫码登录（有 code）
    │                                      │
    ▼                                      ▼
前端: GET /api/v1/user/me          前端: POST /api/v1/auth/wechat
(带 auth-center token)            (带 code)
    │                                      │
    ▼                                      ▼
后端: 用 token 调用 auth-center      后端: 用 code 调用 auth-center
      /api/auth/user-info                  /api/auth/wechat/login
    │                                      │
    ▼                                      ▼
返回: JWT token + 用户资料        返回: JWT token + 用户资料
    │                                      │
    └──────────────────┬──────────────────┘
                       ▼
                  前端保存 JWT token
                       │
                       ▼
                  后续请求使用 JWT token
```

---

## 常见问题

### Q1: 为什么 auth-center 登录响应的 token 是空的？
**A**: 这是正常设计。token 字段确实为空，但 `data.profile` 包含用户资料。

### Q2: 前端能否直接调用 auth-center API？
**A**: ❌ 不能。所有与 auth-center 的通信必须通过后端。

### Q3: 是否需要保存 auth-center 的 token？
**A**: ❌ 不需要。auth-center token 只用于调用 auth-center API，用完即丢。

### Q4: 微信昵称和头像从哪里获取？
**A**:
- **PC扫码登录**：从 auth-center `/api/auth/wechat/login` 响应的 `data.profile` 获取
- **微信内登录**：从 auth-center `/api/auth/user-info` 响应的 `data.profile` 获取

### Q5: 如何刷新用户资料？
**A**: 前端调用 `/api/v1/user/me` 时，后端会用 auth-center token 同步最新资料。

---

## 文件位置

### 后端
- 控制器：`backend/controllers/auth.go`
- 路由：`backend/routes/routes.go`
- 模型：`backend/models/user.go`
- 配置：`backend/.env.production`

### 前端
- 登录页：`frontend/src/pages/Login.tsx`
- AuthContext：`frontend/src/contexts/AuthContext.tsx`
- 类型定义：`frontend/src/types/index.ts`

---

**文档版本**：v1.0
**最后更新**：2026-02-05
**维护者**：PR Business Team
