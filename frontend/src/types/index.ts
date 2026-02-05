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
  phoneNumber: string
  password: string
}

// 登录响应
export interface LoginResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  userId: string
  nickname?: string
  avatarUrl?: string
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
  name: string
  description?: string
}

export interface UpdateServiceProviderRequest {
  name?: string
  description?: string
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

// 营销活动类型
export interface Campaign {
  id: string
  merchantId: string
  providerId: string | null
  title: string
  requirements: string
  platforms: string // JSONB string
  taskAmount: number
  campaignAmount: number
  creatorAmount: number | null
  staffReferralAmount: number | null
  providerAmount: number | null
  quota: number
  taskDeadline: string
  submissionDeadline: string
  status: 'DRAFT' | 'OPEN' | 'CLOSED'
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  merchant?: Merchant
  provider?: ServiceProvider
  tasks?: Task[]
}

export interface CreateCampaignRequest {
  title: string
  requirements: string
  platforms: string
  taskAmount: number
  campaignAmount: number
  creatorAmount?: number
  staffReferralAmount?: number
  providerAmount?: number
  quota: number
  taskDeadline: string
  submissionDeadline: string
}

export interface UpdateCampaignRequest {
  title?: string
  requirements?: string
  platforms?: string
  taskAmount?: number
  campaignAmount?: number
  creatorAmount?: number
  staffReferralAmount?: number
  providerAmount?: number
  quota?: number
  taskDeadline?: string
  submissionDeadline?: string
  status?: 'DRAFT' | 'OPEN' | 'CLOSED'
}

// 任务类型
export interface Task {
  id: string
  campaignId: string
  taskSlotNumber: number
  status: 'OPEN' | 'ASSIGNED' | 'SUBMITTED' | 'APPROVED' | 'REJECTED'
  creatorId: string | null
  assignedAt: string | null
  platform: string
  platformUrl: string
  screenshots: string // JSONB string
  submittedAt: string | null
  notes: string
  auditedBy: string | null
  auditedAt: string | null
  auditNote: string
  inviterId: string | null
  inviterType: string
  priority: 'HIGH' | 'MEDIUM' | 'LOW'
  tags: string[]
  version: number
  createdAt: string
  updatedAt: string
  campaign?: Campaign
  creator?: Creator
  auditor?: User
  inviter?: User
}

export interface AcceptTaskRequest {
  platform: string
}

export interface SubmitTaskRequest {
  platformUrl: string
  screenshots?: string
  notes?: string
}

export interface AuditTaskRequest {
  action: 'approve' | 'reject'
  auditNote?: string
}

export interface TasksResponse {
  data: Task[]
  total: number
  page: number
  page_size: number
}

// 积分账户类型
export interface CreditAccount {
  id: string
  ownerId: string
  ownerType: 'ORG_MERCHANT' | 'ORG_PROVIDER' | 'USER_PERSONAL'
  balance: number
  frozenBalance: number
  createdAt: string
  updatedAt: string
}

// 积分流水类型
export interface CreditTransaction {
  id: string
  accountId: string
  type: string
  amount: number
  balanceBefore: number
  balanceAfter: number
  transactionGroupId: string | null
  groupSequence: number | null
  relatedCampaignId: string | null
  relatedTaskId: string | null
  description: string
  createdAt: string
}

export interface TransactionsResponse {
  data: CreditTransaction[]
  total: number
  page: number
  page_size: number
}

export interface RechargeRequest {
  amount: number
}

// 提现类型
export interface Withdrawal {
  id: string
  accountId: string
  amount: number
  fee: number
  actualAmount: number
  method: 'ALIPAY' | 'WECHAT' | 'BANK'
  accountInfo: string // JSONB string
  accountInfoHash: string
  status: 'pending' | 'approved' | 'rejected' | 'completed'
  auditNote: string
  auditedBy: string | null
  auditedAt: string | null
  completedAt: string | null
  createdAt: string
  updatedAt: string
  account?: CreditAccount
  auditor?: User
}

export interface CreateWithdrawalRequest {
  amount: number
  method: 'ALIPAY' | 'WECHAT' | 'BANK'
  accountInfo: Record<string, any>
}

export interface AuditWithdrawalRequest {
  approved: boolean
  auditNote?: string
}

export interface WithdrawalsResponse {
  withdrawals: Withdrawal[]
  total: number
}

