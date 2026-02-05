import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import type { User } from '../types'

interface AuthContextType {
  user: User | null
  token: string | null
  login: (token: string, refreshToken: string, user: User) => void
  logout: () => void
  updateUser: (user: User) => void
  isAuthenticated: boolean
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // 🔴 开发模式：提供模拟的 super_admin 用户
  useEffect(() => {
    // 检查是否已经有用户数据（真实登录或开发模式）
    const storedUser = localStorage.getItem('user')

    if (storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser)
        setToken(localStorage.getItem('accessToken'))
        setUser(parsedUser)
        console.log('从 localStorage 恢复用户:', parsedUser.currentRole)
      } catch (e) {
        console.error('解析用户数据失败', e)
      }
    } else if (import.meta.env.DEV) {
      // 开发模式且未登录：使用模拟用户
      const mockUser: User = {
        id: 'usr_78a4785a-a30f-49ca-bfe9-680d59da13c7',
        authCenterUserId: '300d0851-7a28-4ad0-98dc-98ac29811945',
        nickname: '测试管理员',
        avatarUrl: '',
        profile: {},
        roles: ['SUPER_ADMIN', 'ADMIN', 'MERCHANT_ADMIN', 'SP_ADMIN', 'CREATOR'],
        currentRole: 'super_admin',
        lastUsedRole: '',
        status: 'active',
        lastLoginAt: new Date().toISOString(),
        lastLoginIp: '127.0.0.1',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      }
      setUser(mockUser)
      // 保存到 localStorage，这样刷新后不会丢失
      localStorage.setItem('user', JSON.stringify(mockUser))
      console.log('🔴 开发模式：使用模拟用户', mockUser)
    }
    setIsLoading(false)
  }, [])

  const login = (newToken: string, refreshToken: string, newUser: User) => {
    localStorage.setItem('accessToken', newToken)
    localStorage.setItem('refreshToken', refreshToken)
    localStorage.setItem('user', JSON.stringify(newUser))
    setToken(newToken)
    setUser(newUser)
  }

  const logout = () => {
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
    setToken(null)
    setUser(null)
  }

  const updateUser = (updatedUser: User) => {
    localStorage.setItem('user', JSON.stringify(updatedUser))
    setUser(updatedUser)
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        login,
        logout,
        updateUser,
        isAuthenticated: !!user,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
