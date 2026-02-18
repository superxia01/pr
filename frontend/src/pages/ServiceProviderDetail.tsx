import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { serviceProviderApi } from '../services/api'
import type { ServiceProvider, ServiceProviderStaff } from '../types'

export default function ServiceProviderDetail() {
  const { id } = useParams<{ id: string }>()
  const [provider, setProvider] = useState<ServiceProvider | null>(null)
  const [staff, setStaff] = useState<ServiceProviderStaff[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (id) loadData()
  }, [id])

  const loadData = async () => {
    if (!id) return
    setLoading(true)
    setError(null)
    try {
      const providerData = await serviceProviderApi.getServiceProvider(id)
      setProvider(providerData)
      const staffData = await serviceProviderApi.getServiceProviderStaff(id)
      setStaff(staffData)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p className="mt-2 text-gray-500">加载中...</p>
        </div>
      </div>
    )
  }

  if (error || !provider) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error || '服务商不存在'}
          </div>
          <Link to="/service-providers" className="text-blue-600 hover:underline">返回服务商列表</Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="mb-4">
          <Link to="/service-providers" className="text-blue-600 hover:underline text-sm">← 返回服务商列表</Link>
        </div>
        <div className="bg-white rounded-lg shadow">
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">服务商详情</h1>
            <p className="mt-1 text-sm text-gray-500">{provider.name}</p>
          </div>
          <div className="p-6 space-y-6">
            <div>
              <h2 className="text-lg font-medium text-gray-900 mb-4">基本信息</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700">服务商名称</label>
                  <p className="mt-1 text-gray-900">{provider.name}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">描述</label>
                  <p className="mt-1 text-gray-900">{provider.description || '-'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">管理员</label>
                  <p className="mt-1 text-gray-900">{provider.admin?.nickname || provider.adminId || '-'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">状态</label>
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    provider.status === 'active' ? 'bg-blue-100 text-blue-800' :
                    provider.status === 'suspended' ? 'bg-red-100 text-red-800' :
                    'bg-gray-100 text-gray-800'
                  }`}>
                    {provider.status === 'active' ? '正常' :
                     provider.status === 'suspended' ? '暂停' :
                     provider.status === 'inactive' ? '未激活' : provider.status}
                  </span>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">创建时间</label>
                  <p className="mt-1 text-gray-500">{new Date(provider.createdAt).toLocaleString()}</p>
                </div>
              </div>
            </div>

            <div>
              <h2 className="text-lg font-medium text-gray-900 mb-4">员工 ({staff.length})</h2>
              {staff.length === 0 ? (
                <p className="text-gray-500 text-sm">暂无员工</p>
              ) : (
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">姓名/用户</th>
                        <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">职位</th>
                        <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">状态</th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {staff.map((s) => (
                        <tr key={s.id}>
                          <td className="px-4 py-3 text-sm text-gray-900">{s.user?.nickname || s.userId}</td>
                          <td className="px-4 py-3 text-sm text-gray-500">{s.title || '-'}</td>
                          <td className="px-4 py-3">
                            <span className={`inline-flex px-2.5 py-0.5 rounded-full text-xs font-medium ${
                              s.status === 'active' ? 'bg-blue-100 text-blue-800' : 'bg-gray-100 text-gray-800'
                            }`}>
                              {s.status === 'active' ? '正常' : '未激活'}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>

            {provider.merchants && provider.merchants.length > 0 && (
              <div>
                <h2 className="text-lg font-medium text-gray-900 mb-4">关联商家 ({provider.merchants.length})</h2>
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">商家名称</th>
                        <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">管理员</th>
                        <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">状态</th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {provider.merchants.map((m) => (
                        <tr key={m.id}>
                          <td className="px-4 py-3 text-sm text-gray-900">{m.name}</td>
                          <td className="px-4 py-3 text-sm text-gray-500">{m.admin?.nickname || m.adminId || '-'}</td>
                          <td className="px-4 py-3">
                            <span className={`inline-flex px-2.5 py-0.5 rounded-full text-xs font-medium ${
                              m.status === 'active' ? 'bg-blue-100 text-blue-800' :
                              m.status === 'suspended' ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-gray-800'
                            }`}>
                              {m.status === 'active' ? '正常' : m.status === 'suspended' ? '暂停' : m.status}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
