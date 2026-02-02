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

// 邀请码类型
export interface InvitationCode {
  id: string
  code: string
  codeType: 'ADMIN_MASTER' | 'SP_ADMIN' | 'MERCHANT' | 'CREATOR' | 'STAFF'
  ownerId: string
  ownerType: 'super_admin' | 'service_provider' | 'merchant'
  status: 'active' | 'disabled' | 'expired' | 'used'
  maxUses: number
  useCount: number
  expiresAt: string | null
  usedBy: string[]
  metadata: Record<string, any>
  createdAt: string
  updatedAt: string
  deletedAt: string | null
}

// 创建邀请码请求
export interface CreateInvitationCodeRequest {
  codeType: string
  ownerId: string
  ownerType: string
  maxUses?: number
  expiresAt?: string
}

// 使用邀请码请求
export interface UseInvitationCodeRequest {
  code: string
  userId: string
}

// 使用邀请码响应
export interface UseInvitationCodeResponse {
  message: string
  code: string
  codeType: string
  ownerId: string
}
