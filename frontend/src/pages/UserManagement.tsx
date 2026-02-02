import { useEffect, useState } from 'react'
import { merchantApi, serviceProviderApi, creatorApi } from '../services/api'
import type { Merchant, ServiceProvider, Creator } from '../types'

export default function UserManagement() {
  const [activeTab, setActiveTab] = useState<'merchants' | 'providers' | 'creators'>('merchants')
  const [merchants, setMerchants] = useState<Merchant[]>([])
  const [providers, setProviders] = useState<ServiceProvider[]>([])
  const [creators, setCreators] = useState<Creator[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadData()
  }, [activeTab])

  const loadData = async () => {
    setLoading(true)
    setError(null)
    try {
      switch (activeTab) {
        case 'merchants':
          const merchantsData = await merchantApi.getMerchants()
          setMerchants(merchantsData)
          break
        case 'providers':
          const providersData = await serviceProviderApi.getServiceProviders()
          setProviders(providersData)
          break
        case 'creators':
          const creatorsData = await creatorApi.getCreators({ page: 1, page_size: 100 })
          setCreators(creatorsData.data)
          break
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '加载数据失败')
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteMerchant = async (id: string) => {
    if (!confirm('确定要删除此商家吗？')) return
    try {
      await merchantApi.deleteMerchant(id)
      loadData()
    } catch (err: any) {
      alert(err.response?.data?.error || '删除失败')
    }
  }

  const handleDeleteProvider = async (id: string) => {
    if (!confirm('确定要删除此服务商吗？')) return
    try {
      await serviceProviderApi.deleteServiceProvider(id)
      loadData()
    } catch (err: any) {
      alert(err.response?.data?.error || '删除失败')
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">用户管理</h1>
            <p className="mt-1 text-sm text-gray-500">管理系统中的所有用户和组织</p>
          </div>

          {/* 标签页 */}
          <div className="border-b border-gray-200">
            <nav className="flex -mb-px">
              <button
                onClick={() => setActiveTab('merchants')}
                className={`px-6 py-4 text-sm font-medium border-b-2 ${
                  activeTab === 'merchants'
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                商家
              </button>
              <button
                onClick={() => setActiveTab('providers')}
                className={`px-6 py-4 text-sm font-medium border-b-2 ${
                  activeTab === 'providers'
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                服务商
              </button>
              <button
                onClick={() => setActiveTab('creators')}
                className={`px-6 py-4 text-sm font-medium border-b-2 ${
                  activeTab === 'creators'
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                达人
              </button>
            </nav>
          </div>

          {/* 内容 */}
          <div className="p-6">
            {loading && (
              <div className="text-center py-12">
                <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                <p className="mt-2 text-gray-500">加载中...</p>
              </div>
            )}

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                {error}
              </div>
            )}

            {/* 商家列表 */}
            {!loading && activeTab === 'merchants' && (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        名称
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        管理员
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        服务商
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        状态
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        操作
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {merchants.map((merchant) => (
                      <tr key={merchant.id}>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-medium text-gray-900">{merchant.name}</div>
                          {merchant.description && (
                            <div className="text-sm text-gray-500 truncate max-w-xs">{merchant.description}</div>
                          )}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {merchant.admin?.nickname || merchant.adminId}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {merchant.provider?.name || merchant.providerId}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                            merchant.status === 'active' ? 'bg-green-100 text-green-800' :
                            merchant.status === 'suspended' ? 'bg-red-100 text-red-800' :
                            'bg-gray-100 text-gray-800'
                          }`}>
                            {merchant.status === 'active' ? '正常' :
                             merchant.status === 'suspended' ? '暂停' :
                             merchant.status === 'inactive' ? '未激活' : merchant.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                          <button className="text-blue-600 hover:text-blue-900 mr-3">编辑</button>
                          <button
                            onClick={() => handleDeleteMerchant(merchant.id)}
                            className="text-red-600 hover:text-red-900"
                          >
                            删除
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {merchants.length === 0 && (
                  <div className="text-center py-12 text-gray-500">暂无商家数据</div>
                )}
              </div>
            )}

            {/* 服务商列表 */}
            {!loading && activeTab === 'providers' && (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        名称
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        管理员
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        商家数量
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        状态
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        操作
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {providers.map((provider) => (
                      <tr key={provider.id}>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-medium text-gray-900">{provider.name}</div>
                          {provider.description && (
                            <div className="text-sm text-gray-500 truncate max-w-xs">{provider.description}</div>
                          )}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {provider.admin?.nickname || provider.adminId}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {provider.merchants?.length || 0}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                            provider.status === 'active' ? 'bg-green-100 text-green-800' :
                            provider.status === 'suspended' ? 'bg-red-100 text-red-800' :
                            'bg-gray-100 text-gray-800'
                          }`}>
                            {provider.status === 'active' ? '正常' :
                             provider.status === 'suspended' ? '暂停' :
                             provider.status === 'inactive' ? '未激活' : provider.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                          <button className="text-blue-600 hover:text-blue-900 mr-3">编辑</button>
                          <button
                            onClick={() => handleDeleteProvider(provider.id)}
                            className="text-red-600 hover:text-red-900"
                          >
                            删除
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {providers.length === 0 && (
                  <div className="text-center py-12 text-gray-500">暂无服务商数据</div>
                )}
              </div>
            )}

            {/* 达人列表 */}
            {!loading && activeTab === 'creators' && (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        昵称
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        等级
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        粉丝数
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        邀请人
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        状态
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        操作
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {creators.map((creator) => (
                      <tr key={creator.id}>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-medium text-gray-900">
                            {creator.user?.nickname || creator.wechatNickname}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                            creator.level === 'KOL' ? 'bg-purple-100 text-purple-800' :
                            creator.level === 'INF' ? 'bg-blue-100 text-blue-800' :
                            creator.level === 'KOC' ? 'bg-green-100 text-green-800' :
                            'bg-gray-100 text-gray-800'
                          }`}>
                            {creator.level}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {creator.followersCount.toLocaleString()}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {creator.inviter?.nickname || creator.inviterId || '-'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                            creator.status === 'active' ? 'bg-green-100 text-green-800' :
                            creator.status === 'banned' ? 'bg-red-100 text-red-800' :
                            'bg-gray-100 text-gray-800'
                          }`}>
                            {creator.status === 'active' ? '正常' :
                             creator.status === 'banned' ? '封禁' :
                             creator.status === 'inactive' ? '未激活' : creator.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                          <button className="text-blue-600 hover:text-blue-900">查看</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {creators.length === 0 && (
                  <div className="text-center py-12 text-gray-500">暂无达人数据</div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
