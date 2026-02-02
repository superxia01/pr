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

// 商家类型
export interface Merchant {
  id: string
  adminId: string
  providerId: string
  userId: string
  name: string
  description: string
  logoUrl: string
  industry: string
  status: 'active' | 'suspended' | 'inactive'
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  admin?: User
  provider?: ServiceProvider
  user?: User
  staff?: MerchantStaff[]
}

export interface CreateMerchantRequest {
  providerId: string
  userId: string
  name: string
  description?: string
  logoUrl?: string
  industry?: string
}

export interface UpdateMerchantRequest {
  name?: string
  description?: string
  logoUrl?: string
  industry?: string
  status?: 'active' | 'suspended' | 'inactive'
}

// 商家员工类型
export interface MerchantStaff {
  id: string
  userId: string
  merchantId: string
  title: string
  status: 'active' | 'inactive'
  createdAt: string
  updatedAt: string
  user?: User
  merchant?: Merchant
  permissions?: MerchantStaffPermission[]
}

export interface AddMerchantStaffRequest {
  userId: string
  title?: string
  permissions: string[]
}

export interface MerchantStaffPermission {
  id: string
  staffId: string
  permissionCode: string
  grantedAt: string
  grantedBy: string
}

// 服务商类型
export interface ServiceProvider {
  id: string
  adminId: string
  userId: string
  name: string
  description: string
  logoUrl: string
  status: 'active' | 'suspended' | 'inactive'
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  admin?: User
  user?: User
  staff?: ServiceProviderStaff[]
  merchants?: Merchant[]
}

export interface CreateServiceProviderRequest {
  userId: string
  name: string
  description?: string
  logoUrl?: string
}

export interface UpdateServiceProviderRequest {
  name?: string
  description?: string
  logoUrl?: string
  status?: 'active' | 'suspended' | 'inactive'
}

// 服务商员工类型
export interface ServiceProviderStaff {
  id: string
  userId: string
  providerId: string
  title: string
  status: 'active' | 'inactive'
  createdAt: string
  updatedAt: string
  user?: User
  provider?: ServiceProvider
  permissions?: ServiceProviderStaffPermission[]
}

export interface AddServiceProviderStaffRequest {
  userId: string
  title?: string
  permissions: string[]
}

export interface ServiceProviderStaffPermission {
  id: string
  staffId: string
  permissionCode: string
  grantedAt: string
  grantedBy: string
}

// 达人类型
export interface Creator {
  id: string
  userId: string
  isPrimary: boolean
  level: 'UGC' | 'KOC' | 'INF' | 'KOL'
  followersCount: number
  wechatOpenId: string
  wechatNickname: string
  wechatAvatar: string
  inviterId: string | null
  inviterType: string | null
  inviterRelationshipBroken: boolean
  status: 'active' | 'banned' | 'inactive'
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  user?: User
  inviter?: User
}

export interface UpdateCreatorRequest {
  level?: 'UGC' | 'KOC' | 'INF' | 'KOL'
  followersCount?: number
  wechatOpenId?: string
  wechatNickname?: string
  wechatAvatar?: string
  status?: 'active' | 'banned' | 'inactive'
}

export interface CreatorsResponse {
  data: Creator[]
  total: number
  page: number
  page_size: number
}

export interface CreatorLevelStats {
  level: string
  count: number
}
