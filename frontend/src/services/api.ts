import axios, { AxiosError } from 'axios'
import type {
  LoginRequest,
  PasswordLoginRequest,
  LoginResponse,
  RefreshTokenRequest,
  SwitchRoleRequest,
  User,
  ApiError
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

export default api
