import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { serviceProviderApi } from '../services/api'
import { PermissionDialog } from '../components/PermissionDialog'
import type { ServiceProvider, ServiceProviderStaff } from '../types'

export default function ServiceProviderInfo() {
  const [provider, setProvider] = useState<ServiceProvider | null>(null)
  const [staff, setStaff] = useState<ServiceProviderStaff[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editMode, setEditMode] = useState(false)
  const [permissionDialogOpen, setPermissionDialogOpen] = useState(false)
  const [selectedStaff, setSelectedStaff] = useState<ServiceProviderStaff | null>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    setError(null)
    try {
      const providerData = await serviceProviderApi.getMyServiceProvider()
      setProvider(providerData)
      if (providerData.id) {
        const staffData = await serviceProviderApi.getServiceProviderStaff(providerData.id)
        setStaff(staffData)
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '加载数据失败')
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateProvider = async (data: any) => {
    if (!provider) return
    try {
      await serviceProviderApi.updateServiceProvider(provider.id, data)
      loadData()
      setEditMode(false)
    } catch (err: any) {
      alert(err.response?.data?.error || '更新失败')
    }
  }

  const handleDeleteStaff = async (staffId: string) => {
    if (!provider || !confirm('确定要删除此员工吗？')) return
    try {
      await serviceProviderApi.deleteServiceProviderStaff(provider.id, staffId)
      loadData()
    } catch (err: any) {
      alert(err.response?.data?.error || '删除失败')
    }
  }

  const handleOpenPermissionDialog = (staffMember: ServiceProviderStaff) => {
    setSelectedStaff(staffMember)
    setPermissionDialogOpen(true)
  }

  const handleManagePermissions = async (permissions: string[]) => {
    if (!provider || !selectedStaff) return

    try {
      // Get current permissions
      const currentPermissions = selectedStaff.permissions?.map(p => p.permissionCode) || []

      // Grant new permissions
      for (const permCode of permissions) {
        if (!currentPermissions.includes(permCode)) {
          await serviceProviderApi.updateServiceProviderStaffPermission(
            provider.id,
            selectedStaff.id,
            { permissionCode: permCode, action: 'grant' }
          )
        }
      }

      // Revoke removed permissions
      for (const permCode of currentPermissions) {
        if (!permissions.includes(permCode)) {
          await serviceProviderApi.updateServiceProviderStaffPermission(
            provider.id,
            selectedStaff.id,
            { permissionCode: permCode, action: 'revoke' }
          )
        }
      }

      loadData()
    } catch (err: any) {
      throw err
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

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        </div>
      </div>
    )
  }

  if (!provider) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto text-center py-12">
          <p className="text-gray-500">您不是服务商管理员</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">服务商信息</h1>
              <p className="mt-1 text-sm text-gray-500">管理您的服务商信息和员工</p>
            </div>
            {!editMode && (
              <button
                onClick={() => setEditMode(true)}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                编辑信息
              </button>
            )}
          </div>

          <div className="p-6">
            {/* 服务商基本信息 */}
            <div className="bg-white rounded-lg shadow mb-8 p-6">
              <h2 className="text-lg font-medium text-gray-900 mb-4">基本信息</h2>
              {editMode ? (
                <EditProviderForm
                  provider={provider}
                  onSave={handleUpdateProvider}
                  onCancel={() => setEditMode(false)}
                />
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700">服务商名称</label>
                    <p className="mt-1 text-sm text-gray-900">{provider.name}</p>
                  </div>
                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700">服务商描述</label>
                    <p className="mt-1 text-sm text-gray-900">{provider.description || '-'}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">商家数量</label>
                    <p className="mt-1 text-sm text-gray-900">{provider.merchants?.length || 0}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">状态</label>
                    <div className="mt-1">
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
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">创建时间</label>
                    <p className="mt-1 text-sm text-gray-500">{new Date(provider.createdAt).toLocaleString()}</p>
                  </div>
                </div>
              )}
            </div>

            {/* 员工管理 */}
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-lg font-medium text-gray-900">员工管理</h2>
                <Link
                  to="/invitations"
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 inline-block"
                >
                  邀请员工
                </Link>
              </div>
              {staff.length === 0 ? (
                <div className="text-center py-12 bg-gray-50 rounded-lg">
                  <p className="text-gray-500">暂无员工</p>
                </div>
              ) : (
                <div className="overflow-hidden">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          姓名
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          职位
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          权限数量
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
                      {staff.map((s) => (
                        <tr key={s.id}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {s.user?.nickname || s.userId}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {s.title || '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {s.permissions?.length || 0}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                              s.status === 'active' ? 'bg-blue-100 text-blue-800' : 'bg-gray-100 text-gray-800'
                            }`}>
                              {s.status === 'active' ? '正常' : '未激活'}
                            </span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                            <button
                              onClick={() => handleOpenPermissionDialog(s)}
                              className="text-blue-600 hover:text-blue-900 mr-3"
                            >
                              管理权限
                            </button>
                            <button
                              onClick={() => handleDeleteStaff(s.id)}
                              className="text-red-600 hover:text-red-900"
                            >
                              删除
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* 商家列表 */}
        {provider.merchants && provider.merchants.length > 0 && (
          <div className="bg-white rounded-lg shadow mt-6">
            <div className="px-6 py-4 border-b border-gray-200">
              <h2 className="text-lg font-medium text-gray-900">关联商家</h2>
            </div>
            <div className="p-6">
              <div className="overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        商家名称
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        管理员
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        状态
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {provider.merchants.map((merchant) => (
                      <tr key={merchant.id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {merchant.name}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {merchant.admin?.nickname || merchant.adminId}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                            merchant.status === 'active' ? 'bg-blue-100 text-blue-800' :
                            merchant.status === 'suspended' ? 'bg-red-100 text-red-800' :
                            'bg-gray-100 text-gray-800'
                          }`}>
                            {merchant.status === 'active' ? '正常' :
                             merchant.status === 'suspended' ? '暂停' :
                             merchant.status === 'inactive' ? '未激活' : merchant.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}

        {/* 权限管理对话框 */}
        {selectedStaff && (
          <PermissionDialog
            open={permissionDialogOpen}
            onOpenChange={setPermissionDialogOpen}
            staffName={selectedStaff.user?.nickname || selectedStaff.userId}
            staffRole="SERVICE_PROVIDER_STAFF"
            currentPermissions={selectedStaff.permissions?.map(p => p.permissionCode) || []}
            onSave={handleManagePermissions}
          />
        )}
      </div>
    </div>
  )
}

// 编辑服务商信息表单
function EditProviderForm({
  provider,
  onSave,
  onCancel
}: {
  provider: ServiceProvider
  onSave: (data: any) => void
  onCancel: () => void
}) {
  const [formData, setFormData] = useState({
    name: provider.name,
    description: provider.description || '',
    logoUrl: provider.logoUrl || '',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(formData)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-gray-700">服务商名称</label>
        <input
          id="name"
          type="text"
          value={formData.name}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, name: e.target.value })}
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
          required
        />
      </div>
      <div>
        <label htmlFor="logoUrl" className="block text-sm font-medium text-gray-700">Logo URL</label>
        <input
          id="logoUrl"
          type="url"
          value={formData.logoUrl}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData({ ...formData, logoUrl: e.target.value })}
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
        />
      </div>
      <div>
        <label htmlFor="description" className="block text-sm font-medium text-gray-700">服务商描述</label>
        <textarea
          id="description"
          value={formData.description}
          onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setFormData({ ...formData, description: e.target.value })}
          rows={3}
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
        />
      </div>
      <div className="flex justify-end space-x-3">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
        >
          取消
        </button>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          保存
        </button>
      </div>
    </form>
  )
}
