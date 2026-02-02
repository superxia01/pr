import { useAuth } from '../contexts/AuthContext'

export default function Dashboard() {
  const { user, logout } = useAuth()

  const canManageInvitations = () => {
    const role = user?.currentRole
    return role === 'super_admin' || role === 'service_provider_admin' || role === 'merchant_admin'
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 顶部导航 */}
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900">PR Business</h1>
            </div>
            <div className="flex items-center gap-4">
              {canManageInvitations() && (
                <a
                  href="/invitations"
                  className="px-4 py-2 text-sm text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg transition-colors"
                >
                  邀请码管理
                </a>
              )}
              <span className="text-sm text-gray-700">
                欢迎，{user?.nickname || '用户'}
              </span>
              <button
                onClick={logout}
                className="px-4 py-2 text-sm text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
              >
                退出登录
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* 主要内容 */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">工作台</h2>

          <div className="space-y-4">
            <div className="p-4 bg-blue-50 rounded-lg">
              <h3 className="font-medium text-blue-900">当前角色</h3>
              <p className="text-blue-700">{user?.currentRole || '未设置'}</p>
            </div>

            <div className="p-4 bg-green-50 rounded-lg">
              <h3 className="font-medium text-green-900">拥有角色</h3>
              <ul className="text-green-700 list-disc list-inside">
                {user?.roles.map((role) => (
                  <li key={role}>{role}</li>
                ))}
              </ul>
            </div>

            <div className="p-4 bg-purple-50 rounded-lg">
              <h3 className="font-medium text-purple-900">用户ID</h3>
              <p className="text-purple-700 text-sm">{user?.id}</p>
            </div>

            {canManageInvitations() && (
              <div className="p-4 bg-yellow-50 rounded-lg">
                <h3 className="font-medium text-yellow-900">快捷操作</h3>
                <a
                  href="/invitations"
                  className="text-yellow-700 hover:text-yellow-800 underline"
                >
                  管理邀请码 →
                </a>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
