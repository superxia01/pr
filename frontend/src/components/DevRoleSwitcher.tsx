import { useAuth } from '../contexts/AuthContext'
import { getRoleName } from '../lib/roles'

// 开发模式：所有可用的角色
const ALL_ROLES = [
  'SUPER_ADMIN',
  'ADMIN',
  'MERCHANT_ADMIN',
  'MERCHANT_STAFF',
  'SP_ADMIN',
  'CREATOR',
]

export default function DevRoleSwitcher() {
  const { user, updateUser } = useAuth()

  if (!user || !import.meta.env.DEV) {
    return null
  }

  const handleSwitchRole = (newRole: string) => {
    // 转换为小写下划线格式（当前角色字段格式）
    const roleMapping: Record<string, string> = {
      'SUPER_ADMIN': 'super_admin',
      'ADMIN': 'admin',
      'MERCHANT_ADMIN': 'merchant_admin',
      'MERCHANT_STAFF': 'merchant_staff',
      'SP_ADMIN': 'service_provider_admin',
      'CREATOR': 'creator',
    }

    const newActiveRole = roleMapping[newRole] || newRole.toLowerCase()

    const updatedUser = {
      ...user,
      currentRole: newActiveRole,
      roles: ALL_ROLES,
    }

    // 更新用户上下文
    updateUser(updatedUser)

    // 直接保存到 localStorage，确保刷新后保持
    localStorage.setItem('user', JSON.stringify(updatedUser))

    console.log('🔴 切换角色:', newRole, '->', newActiveRole)

    // 重新加载页面以刷新菜单和权限
    setTimeout(() => {
      window.location.reload()
    }, 100)
  }

  return (
    <div className="fixed bottom-4 right-4 z-50">
      <div className="bg-gray-900 text-white rounded-lg shadow-xl p-4 min-w-[200px]">
        <div className="text-xs font-semibold text-gray-400 mb-2">
          🔴 开发模式 - 切换角色
        </div>
        <div className="space-y-1">
          {ALL_ROLES.map((role) => {
            const isActive = user.currentRole === role.toLowerCase() ||
                           user.currentRole === 'service_provider_admin' && role === 'SP_ADMIN' ||
                           user.currentRole === 'merchant_admin' && role === 'MERCHANT_ADMIN'

            return (
              <button
                key={role}
                onClick={() => handleSwitchRole(role)}
                className={`w-full text-left px-3 py-2 rounded text-sm transition-all ${
                  isActive
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
                }`}
              >
                <span className="mr-2">{isActive ? '✓' : '○'}</span>
                {getRoleName(role)}
              </button>
            )
          })}
        </div>
      </div>
    </div>
  )
}
