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
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)

  // 从localStorage恢复登录状态
  useEffect(() => {
    const storedToken = localStorage.getItem('accessToken')
    const storedUser = localStorage.getItem('user')
    if (storedToken && storedUser) {
      setToken(storedToken)
      setUser(JSON.parse(storedUser))
    }
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
