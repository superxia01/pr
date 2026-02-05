import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { authApi } from '../services/api'
import { getRoleName } from '../lib/roles'
import { Button } from './ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from './ui/dialog'

export default function RoleSwitcher() {
  const { user, updateUser } = useAuth()
  const [switching, setSwitching] = useState(false)
  const [open, setOpen] = useState(false)

  if (!user || user.roles.length <= 1) {
    return null
  }

  const handleSwitchRole = async (newRole: string) => {
    if (newRole === user.currentRole || switching) return

    setSwitching(true)

    try {
      // 转换为小写格式（后端期望小写）
      const roleMapping: Record<string, string> = {
        'SUPER_ADMIN': 'super_admin',
        'MERCHANT_ADMIN': 'merchant_admin',
        'MERCHANT_STAFF': 'merchant_staff',
        'SP_ADMIN': 'provider_admin',
        'SP_STAFF': 'provider_staff',
        'CREATOR': 'creator',
      }

      const lowerCaseRole = roleMapping[newRole] || newRole.toLowerCase()

      const response = await authApi.switchRole({ newRole: lowerCaseRole })
      localStorage.setItem('accessToken', response.accessToken)

      // 更新用户上下文
      updateUser({
        ...user,
        currentRole: response.currentRole,
        lastUsedRole: response.lastUsedRole,
      })

      // 关闭对话框
      setOpen(false)

      // 重新加载页面
      setTimeout(() => {
        window.location.reload()
      }, 300)
    } catch (err: any) {
      alert(err.response?.data?.error || '切换角色失败，请重试')
    } finally {
      setSwitching(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost" className="w-full justify-start text-gray-300 hover:text-white hover:bg-gray-800">
          切换角色
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md bg-gray-800 text-gray-100 border-gray-700">
        <DialogHeader>
          <DialogTitle className="text-white">切换角色</DialogTitle>
          <DialogDescription className="text-gray-400">
            选择要使用的角色，切换后页面菜单会相应变化
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-2 py-4">
          {user.roles.map((role) => {
            const isActive = role === 'SUPER_ADMIN' && user.currentRole === 'super_admin' ||
                           role === 'MERCHANT_ADMIN' && user.currentRole === 'merchant_admin' ||
                           role === 'SP_ADMIN' && user.currentRole === 'provider_admin' ||
                           role === 'CREATOR' && user.currentRole === 'creator'

            return (
              <button
                key={role}
                onClick={() => handleSwitchRole(role)}
                disabled={switching || isActive}
                className={`w-full text-left px-4 py-3 rounded-lg border transition-all ${
                  isActive
                    ? 'border-blue-500 bg-blue-600 text-white'
                    : 'border-gray-600 hover:border-gray-500 hover:bg-gray-700 text-gray-200'
                } ${switching ? 'opacity-50 cursor-not-allowed' : ''}`}
              >
                <div className="flex items-center justify-between">
                  <div>
                    <div className="font-medium">{getRoleName(role)}</div>
                    <div className={`text-xs mt-0.5 ${isActive ? 'text-blue-100' : 'text-gray-400'}`}>
                      {isActive ? '当前角色' : '点击切换'}
                    </div>
                  </div>
                  {isActive && (
                    <span className="text-white">
                      <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                        <path
                          fillRule="evenodd"
                          d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                          clipRule="evenodd"
                        />
                      </svg>
                    </span>
                  )}
                </div>
              </button>
            )
          })}
        </div>

        {switching && (
          <div className="text-center text-sm text-gray-400">
            切换中...
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
