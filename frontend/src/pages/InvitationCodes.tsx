import { useState, useEffect } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { invitationApi } from '../services/api'
import type { InvitationCode } from '../types'

export default function InvitationCodes() {
  const { user } = useAuth()
  const [codes, setCodes] = useState<InvitationCode[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [error, setError] = useState('')

  // 创建邀请码表单
  const [createForm, setCreateForm] = useState({
    codeType: 'MERCHANT',
    ownerId: user?.id || '',
    ownerType: 'merchant',
    maxUses: 1,
    expiresAt: '',
  })

  useEffect(() => {
    loadInvitationCodes()
  }, [])

  const loadInvitationCodes = async () => {
    try {
      setLoading(true)
      const response = await invitationApi.listInvitationCodes()
      setCodes(response.codes)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载邀请码失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateCode = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    try {
      await invitationApi.createInvitationCode(createForm)
      setShowCreateModal(false)
      loadInvitationCodes()
      // 重置表单
      setCreateForm({
        codeType: 'MERCHANT',
        ownerId: user?.id || '',
        ownerType: 'merchant',
        maxUses: 1,
        expiresAt: '',
      })
    } catch (err: any) {
      setError(err.response?.data?.error || '创建邀请码失败')
    }
  }

  const handleDisableCode = async (code: string) => {
    if (!confirm('确定要禁用此邀请码吗？')) return

    try {
      await invitationApi.disableInvitationCode(code)
      loadInvitationCodes()
    } catch (err: any) {
      setError(err.response?.data?.error || '禁用邀请码失败')
    }
  }

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800'
      case 'disabled':
        return 'bg-gray-100 text-gray-800'
      case 'expired':
        return 'bg-yellow-100 text-yellow-800'
      case 'used':
        return 'bg-blue-100 text-blue-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getCodeTypeLabel = (codeType: string) => {
    const labels: Record<string, string> = {
      'ADMIN_MASTER': '超级管理员',
      'SP_ADMIN': '服务商管理员',
      'MERCHANT': '商家',
      'CREATOR': '达人',
      'STAFF': '员工',
    }
    return labels[codeType] || codeType
  }

  const canCreateCode = () => {
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
              <a href="/dashboard" className="text-sm text-gray-700 hover:text-gray-900">
                返回工作台
              </a>
            </div>
          </div>
        </div>
      </nav>

      {/* 主要内容 */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-gray-900">邀请码管理</h2>
          <p className="text-gray-600 mt-1">查看和管理您的邀请码</p>
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-600">
            {error}
          </div>
        )}

        {/* 操作栏 */}
        <div className="mb-6 flex justify-between items-center">
          <div className="text-sm text-gray-700">
            共 <span className="font-medium">{codes.length}</span> 个邀请码
          </div>
          {canCreateCode() && (
            <button
              onClick={() => setShowCreateModal(true)}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              创建邀请码
            </button>
          )}
        </div>

        {/* 邀请码列表 */}
        {loading ? (
          <div className="text-center py-12">
            <div className="text-gray-500">加载中...</div>
          </div>
        ) : codes.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-lg">
            <p className="text-gray-500">暂无邀请码</p>
          </div>
        ) : (
          <div className="bg-white rounded-lg shadow overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    邀请码
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    类型
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    状态
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    使用情况
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    创建时间
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {codes.map((code) => (
                  <tr key={code.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900 font-mono">
                        {code.code}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {getCodeTypeLabel(code.codeType)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusBadgeClass(code.status)}`}>
                        {code.status === 'active' ? '有效' :
                         code.status === 'disabled' ? '已禁用' :
                         code.status === 'expired' ? '已过期' : '已用完'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {code.useCount}/{code.maxUses === 0 ? '无限' : code.maxUses}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(code.createdAt).toLocaleDateString('zh-CN')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      {code.status === 'active' && (
                        <button
                          onClick={() => handleDisableCode(code.code)}
                          className="text-red-600 hover:text-red-900"
                        >
                          禁用
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </main>

      {/* 创建邀请码模态框 */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-lg font-semibold mb-4">创建邀请码</h3>

            <form onSubmit={handleCreateCode} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  邀请码类型
                </label>
                <select
                  value={createForm.codeType}
                  onChange={(e) => setCreateForm({ ...createForm, codeType: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                >
                  <option value="MERCHANT">商家</option>
                  <option value="CREATOR">达人</option>
                  <option value="STAFF">员工</option>
                  {user?.currentRole === 'super_admin' && (
                    <>
                      <option value="SP_ADMIN">服务商管理员</option>
                    </>
                  )}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  最大使用次数
                </label>
                <input
                  type="number"
                  min="1"
                  value={createForm.maxUses}
                  onChange={(e) => setCreateForm({ ...createForm, maxUses: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  过期时间（可选）
                </label>
                <input
                  type="datetime-local"
                  value={createForm.expiresAt}
                  onChange={(e) => setCreateForm({ ...createForm, expiresAt: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                >
                  创建
                </button>
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  className="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
                >
                  取消
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
