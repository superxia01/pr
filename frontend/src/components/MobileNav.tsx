import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { getRoleName } from '../lib/roles'

interface MobileNavProps {
  onLogout: () => void
}

export default function MobileNav({ onLogout }: MobileNavProps) {
  const { user } = useAuth()
  const [isOpen, setIsOpen] = useState(false)

  const canManageInvitations = () => {
    const role = user?.currentRole
    return role === 'super_admin' || role === 'service_provider_admin' || role === 'merchant_admin'
  }

  const isSuperAdmin = () => {
    return user?.roles.includes('SUPER_ADMIN')
  }

  const isMerchantAdmin = () => {
    return user?.roles.includes('MERCHANT_ADMIN')
  }

  const isServiceProviderAdmin = () => {
    return user?.roles.includes('SP_ADMIN')
  }

  const isCreator = () => {
    return user?.roles.includes('CREATOR')
  }

  const isMerchantStaff = () => {
    return user?.roles.includes('MERCHANT_STAFF')
  }

  const navItems = [
    ...(canManageInvitations()
      ? [{ href: '/invitations', label: '邀请码管理', icon: '🎫' }]
      : []),
    ...(isSuperAdmin() ? [{ href: '/user-management', label: '用户管理', icon: '👥' }] : []),
    ...(isMerchantAdmin() ? [{ href: '/merchant', label: '商家信息', icon: '🏢' }] : []),
    ...(isServiceProviderAdmin()
      ? [{ href: '/service-provider', label: '服务商信息', icon: '🏪' }]
      : []),
    ...(isCreator() ? [{ href: '/creator', label: '达人中心', icon: '⭐' }] : []),
    ...(isMerchantAdmin() || isMerchantStaff()
      ? [{ href: '/create-campaign', label: '创建活动', icon: '📢' }]
      : []),
    ...(isCreator()
      ? [
          { href: '/task-hall', label: '任务大厅', icon: '📋' },
          { href: '/my-tasks', label: '我的任务', icon: '✅' },
        ]
      : []),
    { href: '/recharge', label: '充值', icon: '💰' },
    { href: '/credit-transactions', label: '积分明细', icon: '📊' },
    { href: '/withdrawals', label: '提现记录', icon: '💸' },
    ...(isSuperAdmin() ? [{ href: '/withdrawal-review', label: '提现审核', icon: '✍️' }] : []),
  ]

  return (
    <>
      {/* 汉堡菜单按钮 */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="md:hidden p-2 rounded-md text-gray-600 hover:text-gray-900 hover:bg-gray-100"
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          {isOpen ? (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          ) : (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          )}
        </svg>
      </button>

      {/* 移动端导航菜单 */}
      {isOpen && (
        <>
          {/* 遮罩层 */}
          <div
            className="fixed inset-0 bg-black bg-opacity-50 z-40 md:hidden"
            onClick={() => setIsOpen(false)}
          ></div>

          {/* 侧边栏 */}
          <div className="fixed inset-y-0 left-0 w-64 bg-white shadow-xl z-50 md:hidden overflow-y-auto">
            <div className="p-4">
              {/* 标题 */}
              <div className="flex items-center justify-between mb-6">
                <h1 className="text-xl font-bold text-gray-900">PR Business</h1>
                <button
                  onClick={() => setIsOpen(false)}
                  className="p-2 rounded-md text-gray-600 hover:text-gray-900 hover:bg-gray-100"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              {/* 当前角色 */}
              {user?.currentRole && (
                <div className="mb-4 p-3 bg-blue-50 rounded-lg">
                  <p className="text-xs text-blue-600">当前角色</p>
                  <p className="text-sm font-semibold text-blue-900">{getRoleName(user.currentRole)}</p>
                </div>
              )}

              {/* 导航链接 */}
              <nav className="space-y-1">
                {navItems.map((item) => (
                  <a
                    key={item.href}
                    href={item.href}
                    onClick={() => setIsOpen(false)}
                    className="flex items-center gap-3 px-4 py-3 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
                  >
                    <span className="text-xl">{item.icon}</span>
                    <span className="text-sm font-medium">{item.label}</span>
                  </a>
                ))}
              </nav>

              {/* 用户信息和退出 */}
              <div className="mt-6 pt-6 border-t border-gray-200">
                <div className="px-4 mb-4">
                  <p className="text-sm text-gray-900 font-medium">{user?.nickname || '用户'}</p>
                  <p className="text-xs text-gray-500">{user?.id}</p>
                </div>
                <button
                  onClick={onLogout}
                  className="w-full flex items-center gap-3 px-4 py-3 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                    />
                  </svg>
                  <span className="text-sm font-medium">退出登录</span>
                </button>
              </div>
            </div>
          </div>
        </>
      )}
    </>
  )
}
