// 用户类型
export interface User {
  id: string
  authCenterUserId: string
  nickname: string
  avatarUrl: string
  profile: Record<string, any>
  roles: string[]
  currentRole: string
  lastUsedRole: string
  status: string
  lastLoginAt: string | null
  lastLoginIp: string
  createdAt: string
  updatedAt: string
}

// 登录请求
export interface LoginRequest {
  authCode: string
}

export interface PasswordLoginRequest {
  username: string
  password: string
}

// 登录响应
export interface LoginResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  userId: string
  roles: string[]
  currentRole: string
}

// 刷新令牌请求
export interface RefreshTokenRequest {
  refreshToken: string
}

// 切换角色请求
export interface SwitchRoleRequest {
  newRole: string
}

// API错误响应
export interface ApiError {
  error: string
}
