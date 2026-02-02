import axios, { AxiosError } from 'axios'
import type {
  LoginRequest,
  PasswordLoginRequest,
  LoginResponse,
  RefreshTokenRequest,
  SwitchRoleRequest,
  User,
  ApiError,
  InvitationCode,
  CreateInvitationCodeRequest,
  UseInvitationCodeRequest,
  UseInvitationCodeResponse,
  Merchant,
  CreateMerchantRequest,
  UpdateMerchantRequest,
  MerchantStaff,
  AddMerchantStaffRequest,
  ServiceProvider,
  CreateServiceProviderRequest,
  UpdateServiceProviderRequest,
  ServiceProviderStaff,
  AddServiceProviderStaffRequest,
  Creator,
  UpdateCreatorRequest,
  CreatorsResponse,
  CreatorLevelStats,
} from '../types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// 创建axios实例
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 添加token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器 - 处理token过期
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiError>) => {
    if (error.response?.status === 401) {
      // Token过期，尝试刷新
      const refreshToken = localStorage.getItem('refreshToken')
      if (refreshToken) {
        try {
          const response = await axios.post<LoginResponse>(
            `${API_BASE_URL}/api/v1/auth/refresh`,
            { refreshToken }
          )
          const { accessToken } = response.data
          localStorage.setItem('accessToken', accessToken)
          // 重试原请求
          if (error.config) {
            error.config.headers.Authorization = `Bearer ${accessToken}`
            return api.request(error.config)
          }
        } catch (refreshError) {
          // 刷新失败，清除token并跳转登录
          localStorage.removeItem('accessToken')
          localStorage.removeItem('refreshToken')
          localStorage.removeItem('user')
          window.location.href = '/login'
        }
      } else {
        // 没有refreshToken，跳转登录
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

// 认证API
export const authApi = {
  // 微信登录
  wechatLogin: async (data: LoginRequest) => {
    const response = await api.post<LoginResponse>('/api/v1/auth/wechat', data)
    return response.data
  },

  // 密码登录
  passwordLogin: async (data: PasswordLoginRequest) => {
    const response = await api.post<LoginResponse>('/api/v1/auth/password', data)
    return response.data
  },

  // 刷新令牌
  refreshToken: async (data: RefreshTokenRequest) => {
    const response = await api.post<{ accessToken: string; expiresIn: number }>(
      '/api/v1/auth/refresh',
      data
    )
    return response.data
  },

  // 切换角色
  switchRole: async (data: SwitchRoleRequest) => {
    const response = await api.post<{ accessToken: string; currentRole: string; lastUsedRole: string }>(
      '/api/v1/user/switch-role',
      data
    )
    return response.data
  },

  // 获取当前用户
  getCurrentUser: async () => {
    const response = await api.get<User>('/api/v1/user/me')
    return response.data
  },
}

// 邀请码API
export const invitationApi = {
  // 创建邀请码
  createInvitationCode: async (data: CreateInvitationCodeRequest) => {
    const response = await api.post<InvitationCode>('/api/v1/invitations', data)
    return response.data
  },

  // 获取邀请码列表
  listInvitationCodes: async (params?: {
    codeType?: string
    ownerId?: string
    status?: string
  }) => {
    const response = await api.get<{ codes: InvitationCode[]; total: number }>(
      '/api/v1/invitations',
      { params }
    )
    return response.data
  },

  // 获取邀请码详情
  getInvitationCode: async (code: string) => {
    const response = await api.get<InvitationCode>(`/api/v1/invitations/${code}`)
    return response.data
  },

  // 使用邀请码
  useInvitationCode: async (data: UseInvitationCodeRequest) => {
    const response = await api.post<UseInvitationCodeResponse>(
      '/api/v1/invitations/use',
      data
    )
    return response.data
  },

  // 禁用邀请码
  disableInvitationCode: async (code: string) => {
    const response = await api.post<{ message: string }>(
      `/api/v1/invitations/${code}/disable`
    )
    return response.data
  },
}

// 商家API
export const merchantApi = {
  // 创建商家
  createMerchant: async (data: CreateMerchantRequest) => {
    const response = await api.post<Merchant>('/api/v1/merchants', data)
    return response.data
  },

  // 获取商家列表
  getMerchants: async (params?: { status?: string; provider_id?: string }) => {
    const response = await api.get<Merchant[]>('/api/v1/merchants', { params })
    return response.data
  },

  // 获取商家详情
  getMerchant: async (id: string) => {
    const response = await api.get<Merchant>(`/api/v1/merchants/${id}`)
    return response.data
  },

  // 更新商家
  updateMerchant: async (id: string, data: UpdateMerchantRequest) => {
    const response = await api.put<Merchant>(`/api/v1/merchants/${id}`, data)
    return response.data
  },

  // 删除商家
  deleteMerchant: async (id: string) => {
    const response = await api.delete<{ message: string }>(`/api/v1/merchants/${id}`)
    return response.data
  },

  // 获取当前用户管理的商家
  getMyMerchant: async () => {
    const response = await api.get<Merchant>('/api/v1/merchant/me')
    return response.data
  },

  // 添加商家员工
  addMerchantStaff: async (merchantId: string, data: AddMerchantStaffRequest) => {
    const response = await api.post<MerchantStaff>(
      `/api/v1/merchants/${merchantId}/staff`,
      data
    )
    return response.data
  },

  // 获取商家员工列表
  getMerchantStaff: async (merchantId: string) => {
    const response = await api.get<MerchantStaff[]>(
      `/api/v1/merchants/${merchantId}/staff`
    )
    return response.data
  },

  // 更新员工权限
  updateMerchantStaffPermission: async (
    merchantId: string,
    staffId: string,
    data: { permissionCode: string; action: 'grant' | 'revoke' }
  ) => {
    const response = await api.put<{ message: string }>(
      `/api/v1/merchants/${merchantId}/staff/${staffId}/permissions`,
      data
    )
    return response.data
  },

  // 删除商家员工
  deleteMerchantStaff: async (merchantId: string, staffId: string) => {
    const response = await api.delete<{ message: string }>(
      `/api/v1/merchants/${merchantId}/staff/${staffId}`
    )
    return response.data
  },
}

// 服务商API
export const serviceProviderApi = {
  // 创建服务商
  createServiceProvider: async (data: CreateServiceProviderRequest) => {
    const response = await api.post<ServiceProvider>('/api/v1/service-providers', data)
    return response.data
  },

  // 获取服务商列表
  getServiceProviders: async (params?: { status?: string }) => {
    const response = await api.get<ServiceProvider[]>('/api/v1/service-providers', { params })
    return response.data
  },

  // 获取服务商详情
  getServiceProvider: async (id: string) => {
    const response = await api.get<ServiceProvider>(`/api/v1/service-providers/${id}`)
    return response.data
  },

  // 更新服务商
  updateServiceProvider: async (id: string, data: UpdateServiceProviderRequest) => {
    const response = await api.put<ServiceProvider>(`/api/v1/service-providers/${id}`, data)
    return response.data
  },

  // 删除服务商
  deleteServiceProvider: async (id: string) => {
    const response = await api.delete<{ message: string }>(`/api/v1/service-providers/${id}`)
    return response.data
  },

  // 获取当前用户管理的服务商
  getMyServiceProvider: async () => {
    const response = await api.get<ServiceProvider>('/api/v1/service-provider/me')
    return response.data
  },

  // 添加服务商员工
  addServiceProviderStaff: async (providerId: string, data: AddServiceProviderStaffRequest) => {
    const response = await api.post<ServiceProviderStaff>(
      `/api/v1/service-providers/${providerId}/staff`,
      data
    )
    return response.data
  },

  // 获取服务商员工列表
  getServiceProviderStaff: async (providerId: string) => {
    const response = await api.get<ServiceProviderStaff[]>(
      `/api/v1/service-providers/${providerId}/staff`
    )
    return response.data
  },

  // 更新员工权限
  updateServiceProviderStaffPermission: async (
    providerId: string,
    staffId: string,
    data: { permissionCode: string; action: 'grant' | 'revoke' }
  ) => {
    const response = await api.put<{ message: string }>(
      `/api/v1/service-providers/${providerId}/staff/${staffId}/permissions`,
      data
    )
    return response.data
  },

  // 删除服务商员工
  deleteServiceProviderStaff: async (providerId: string, staffId: string) => {
    const response = await api.delete<{ message: string }>(
      `/api/v1/service-providers/${providerId}/staff/${staffId}`
    )
    return response.data
  },
}

// 达人API
export const creatorApi = {
  // 获取达人列表
  getCreators: async (params?: {
    level?: string
    status?: string
    inviter_id?: string
    page?: number
    page_size?: number
  }) => {
    const response = await api.get<CreatorsResponse>('/api/v1/creators', { params })
    return response.data
  },

  // 获取达人详情
  getCreator: async (id: string) => {
    const response = await api.get<Creator>(`/api/v1/creators/${id}`)
    return response.data
  },

  // 更新达人
  updateCreator: async (id: string, data: UpdateCreatorRequest) => {
    const response = await api.put<Creator>(`/api/v1/creators/${id}`, data)
    return response.data
  },

  // 获取达人等级统计
  getCreatorLevelStats: async () => {
    const response = await api.get<{ stats: CreatorLevelStats[] }>(
      '/api/v1/creators/stats/level'
    )
    return response.data
  },

  // 获取达人邀请关系
  getCreatorInviterRelationship: async (id: string) => {
    const response = await api.get<any>(`/api/v1/creators/${id}/inviter`)
    return response.data
  },

  // 解除邀请关系
  breakInviterRelationship: async (id: string) => {
    const response = await api.post<{ message: string }>(
      `/api/v1/creators/${id}/break-relationship`
    )
    return response.data
  },

  // 获取当前用户的达人资料
  getMyCreatorProfile: async () => {
    const response = await api.get<Creator>('/api/v1/creator/me')
    return response.data
  },

  // 更新当前用户的达人资料
  updateMyCreatorProfile: async (data: UpdateCreatorRequest) => {
    const response = await api.put<Creator>('/api/v1/creator/me', data)
    return response.data
  },
}

export default api
