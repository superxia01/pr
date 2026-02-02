import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { authApi } from '../services/api'
import { getRoleName, getRoleDescription, getRoleColorClass } from '../lib/roles'

export default function RoleSwitcher() {
  const { user, updateUser } = useAuth()
  const [switching, setSwitching] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!user || user.roles.length <= 1) {
    return null
  }

  const handleSwitchRole = async (newRole: string) => {
    if (newRole === user.currentRole || switching) return

    // 确认提示
    const currentRoleName = getRoleName(user.currentRole)
    const newRoleName = getRoleName(newRole)
    const confirmed = confirm(
      `确认切换角色吗？\n\n当前角色：${currentRoleName}\n新角色：${newRoleName}\n\n切换后将重新加载页面。`
    )

    if (!confirmed) return

    setSwitching(true)
    setError(null)

    try {
      const response = await authApi.switchRole({ newRole })
      localStorage.setItem('accessToken', response.accessToken)

      // 更新用户上下文
      updateUser({
        ...user,
        currentRole: response.currentRole,
        lastUsedRole: response.lastUsedRole,
      })

      // 成功提示
      alert(`角色已切换为：${getRoleName(newRole)}`)

      // 重新加载页面以刷新所有状态
      setTimeout(() => {
        window.location.reload()
      }, 300)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || '切换角色失败，请重试'
      setError(errorMsg)
      console.error('切换角色失败:', err)
      alert(errorMsg)
    } finally {
      setSwitching(false)
    }
  }

  return (
    <div className="fixed top-4 right-4 z-50">
      <div className="bg-white rounded-lg shadow-xl border border-gray-200 p-4 min-w-[280px]">
        {/* 标题 */}
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-sm font-semibold text-gray-900">切换角色</h3>
          {switching && (
            <div className="flex items-center gap-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
              <span className="text-xs text-gray-500">切换中...</span>
            </div>
          )}
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="mb-3 p-2 bg-red-50 border border-red-200 rounded-md">
            <p className="text-xs text-red-600">{error}</p>
          </div>
        )}

        {/* 当前角色提示 */}
        <div className="mb-3 p-2 bg-blue-50 border border-blue-200 rounded-md">
          <p className="text-xs text-blue-700">
            当前角色：<span className="font-medium">{getRoleName(user.currentRole)}</span>
          </p>
        </div>

        {/* 角色列表 */}
        <div className="space-y-2">
          {user.roles.map((role) => {
            const isActive = role === user.currentRole
            return (
              <button
                key={role}
                onClick={() => handleSwitchRole(role)}
                disabled={switching || isActive}
                className={`w-full text-left px-3 py-2.5 rounded-lg border transition-all ${
                  isActive
                    ? 'border-blue-300 bg-blue-50 shadow-sm'
                    : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
                } ${switching ? 'opacity-50 cursor-not-allowed' : ''}`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      {isActive && (
                        <span className="text-blue-600">
                          <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                            <path
                              fillRule="evenodd"
                              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                              clipRule="evenodd"
                            />
                          </svg>
                        </span>
                      )}
                      <span className={`text-sm font-medium ${isActive ? 'text-blue-700' : 'text-gray-900'}`}>
                        {getRoleName(role)}
                      </span>
                    </div>
                    <p className={`text-xs mt-0.5 ${isActive ? 'text-blue-600' : 'text-gray-500'}`}>
                      {getRoleDescription(role)}
                    </p>
                  </div>
                  <div className="ml-2">
                    <span
                      className={`px-2 py-0.5 rounded text-xs font-medium ${getRoleColorClass(role)}`}
                    >
                      {role}
                    </span>
                  </div>
                </div>
              </button>
            )
          })}
        </div>

        {/* 底部提示 */}
        <div className="mt-3 pt-3 border-t border-gray-200">
          <p className="text-xs text-gray-500">
            💡 提示：切换角色后，工作台菜单和权限会相应变化
          </p>
        </div>
      </div>
    </div>
  )
}
