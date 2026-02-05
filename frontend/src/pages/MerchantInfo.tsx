import { useEffect, useState } from 'react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Label } from '../components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'
import { DataTable } from '../components/ui/data-table'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '../components/ui/dialog'
import { Badge } from '../components/ui/badge'
import { merchantApi } from '../services/api'
import type { Merchant, MerchantStaff } from '../types'

export default function MerchantInfo() {
  const [merchant, setMerchant] = useState<Merchant | null>(null)
  const [staff, setStaff] = useState<MerchantStaff[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editMode, setEditMode] = useState(false)
  const [showAddStaff, setShowAddStaff] = useState(false)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    setError(null)

    // 🔴 开发模式：使用模拟数据
    if (import.meta.env.DEV) {
      setTimeout(() => {
        const mockMerchant: any = {
          id: 'mch_001',
          name: '测试商家',
          description: '这是一个测试商家',
          logoUrl: '',
          providerId: 'sp_001',
          adminId: 'usr_001',
          status: 'active',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          provider: {
            id: 'sp_001',
            name: '测试服务商',
            adminId: 'usr_001',
            userId: 'usr_001',
            description: '测试服务商',
            logoUrl: '',
            status: 'active',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
          admin: {
            id: 'usr_001',
            nickname: '测试管理员',
            authCenterUserId: '',
            avatarUrl: '',
            profile: {},
            roles: [],
            currentRole: '',
            lastUsedRole: '',
            status: 'active',
            lastLoginAt: null,
            lastLoginIp: '',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
        }

        const mockStaff: any[] = [
          {
            id: 'staff_001',
            merchantId: 'mch_001',
            userId: 'usr_002',
            role: 'STAFF',
            joinedAt: new Date().toISOString(),
            user: {
              id: 'usr_002',
              nickname: '员工A',
              authCenterUserId: '',
              avatarUrl: '',
              profile: {},
              roles: [],
              currentRole: '',
              lastUsedRole: '',
              status: 'active',
              lastLoginAt: null,
              lastLoginIp: '',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          },
        ]

        setMerchant(mockMerchant)
        setStaff(mockStaff)
        setLoading(false)
        console.log('🔴 开发模式：使用模拟商家数据')
      }, 500)
      return
    }

    // 生产模式：调用真实 API
    try {
      const merchantData = await merchantApi.getMyMerchant()
      setMerchant(merchantData)
      if (merchantData.id) {
        const staffData = await merchantApi.getMerchantStaff(merchantData.id)
        setStaff(staffData)
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '加载数据失败')
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateMerchant = async (data: any) => {
    if (!merchant) return
    try {
      await merchantApi.updateMerchant(merchant.id, data)
      loadData()
      setEditMode(false)
    } catch (err: any) {
      alert(err.response?.data?.error || '更新失败')
    }
  }

  const handleAddStaff = async (staffData: any) => {
    if (!merchant) return
    try {
      await merchantApi.addMerchantStaff(merchant.id, staffData)
      setShowAddStaff(false)
      loadData()
    } catch (err: any) {
      alert(err.response?.data?.error || '添加失败')
    }
  }

  const handleDeleteStaff = async (staffId: string) => {
    if (!merchant || !confirm('确定要删除此员工吗？')) return
    try {
      await merchantApi.deleteMerchantStaff(merchant.id, staffId)
      loadData()
    } catch (err: any) {
      alert(err.response?.data?.error || '删除失败')
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

  if (!merchant) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto text-center py-12">
          <p className="text-gray-500">您不是商家管理员</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <Card>
          {/* 标题 */}
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-2xl">商家信息</CardTitle>
              <CardDescription className="mt-1">管理您的商家信息和员工</CardDescription>
            </div>
            {!editMode && (
              <Button onClick={() => setEditMode(true)}>
                编辑信息
              </Button>
            )}
          </CardHeader>

          <CardContent>
            {/* 商家基本信息 */}
            <Card className="mb-8">
              <CardHeader>
                <CardTitle className="text-lg">基本信息</CardTitle>
              </CardHeader>
              <CardContent>
                {editMode ? (
                  <EditMerchantForm
                    merchant={merchant}
                    onSave={handleUpdateMerchant}
                    onCancel={() => setEditMode(false)}
                  />
                ) : (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <Label className="text-sm font-medium">商家名称</Label>
                      <p className="mt-1 text-sm text-gray-900">{merchant.name}</p>
                    </div>
                    <div>
                      <Label className="text-sm font-medium">所属服务商</Label>
                      <p className="mt-1 text-sm text-gray-900">{merchant.provider?.name || merchant.providerId}</p>
                    </div>
                    <div className="md:col-span-2">
                      <Label className="text-sm font-medium">商家描述</Label>
                      <p className="mt-1 text-sm text-gray-900">{merchant.description || '-'}</p>
                    </div>
                    {merchant.industry && (
                      <div>
                        <Label className="text-sm font-medium">行业</Label>
                        <p className="mt-1 text-sm text-gray-900">{merchant.industry}</p>
                      </div>
                    )}
                    <div>
                      <Label className="text-sm font-medium">状态</Label>
                      <div className="mt-1">
                        <Badge variant={
                          merchant.status === 'active' ? 'default' :
                          merchant.status === 'suspended' ? 'destructive' : 'secondary'
                        }>
                          {merchant.status === 'active' ? '正常' :
                           merchant.status === 'suspended' ? '暂停' :
                           merchant.status === 'inactive' ? '未激活' : merchant.status}
                        </Badge>
                      </div>
                    </div>
                    <div>
                      <Label className="text-sm font-medium">创建时间</Label>
                      <p className="mt-1 text-sm text-gray-500">{new Date(merchant.createdAt).toLocaleString()}</p>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* 员工管理 */}
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="text-lg">员工管理</CardTitle>
                <Button onClick={() => setShowAddStaff(true)}>
                  添加员工
                </Button>
              </CardHeader>
              <CardContent>
                {staff.length === 0 ? (
                  <div className="text-center py-12 bg-gray-50 rounded-lg">
                    <p className="text-gray-500">暂无员工</p>
                  </div>
                ) : (
                  <DataTable
                    data={staff}
                    columns={[
                      {
                        id: 'name',
                        header: '姓名',
                        accessorKey: 'user',
                        cell: ({ row }) => <span>{row.user?.nickname || row.userId}</span>
                      },
                      {
                        id: 'title',
                        header: '职位',
                        accessorKey: 'title',
                        cell: ({ value }) => <span>{value || '-'}</span>
                      },
                      {
                        id: 'permissions',
                        header: '权限数量',
                        accessorKey: 'permissions',
                        cell: ({ row }) => <span>{row.permissions?.length || 0}</span>
                      },
                      {
                        id: 'status',
                        header: '状态',
                        accessorKey: 'status',
                        cell: ({ row }) => (
                          <Badge variant={row.status === 'active' ? 'default' : 'secondary'}>
                            {row.status === 'active' ? '正常' : '未激活'}
                          </Badge>
                        )
                      },
                      {
                        id: 'actions',
                        header: '操作',
                        cell: ({ row }) => (
                          <div className="space-x-2">
                            <Button variant="outline" size="sm" className="text-blue-600 hover:text-blue-900">
                              编辑权限
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              className="text-red-600 hover:text-red-900"
                              onClick={() => handleDeleteStaff(row.id)}
                            >
                              删除
                            </Button>
                          </div>
                        )
                      }
                    ]}
                    emptyMessage="暂无员工"
                  />
                )}
              </CardContent>
            </Card>

            {/* 添加员工对话框 */}
            <Dialog open={showAddStaff} onOpenChange={setShowAddStaff}>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>添加员工</DialogTitle>
                </DialogHeader>
                <AddStaffForm
                  onSave={handleAddStaff}
                />
              </DialogContent>
            </Dialog>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

// 编辑商家信息表单
function EditMerchantForm({
  merchant,
  onSave,
  onCancel
}: {
  merchant: Merchant
  onSave: (data: any) => void
  onCancel: () => void
}) {
  const [formData, setFormData] = useState({
    name: merchant.name,
    description: merchant.description || '',
    logoUrl: merchant.logoUrl || '',
    industry: merchant.industry || '',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(formData)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <Label htmlFor="name">商家名称</Label>
          <Input
            id="name"
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
          />
        </div>
        <div>
          <Label htmlFor="industry">行业</Label>
          <Input
            id="industry"
            type="text"
            value={formData.industry}
            onChange={(e) => setFormData({ ...formData, industry: e.target.value })}
          />
        </div>
        <div className="md:col-span-2">
          <Label htmlFor="logoUrl">Logo URL</Label>
          <Input
            id="logoUrl"
            type="url"
            value={formData.logoUrl}
            onChange={(e) => setFormData({ ...formData, logoUrl: e.target.value })}
          />
        </div>
        <div className="md:col-span-2">
          <Label htmlFor="description">商家描述</Label>
          <textarea
            id="description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            rows={3}
            className="flex min-h-[60px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          />
        </div>
      </div>
      <div className="flex justify-end space-x-3">
        <Button variant="outline" type="button" onClick={onCancel}>
          取消
        </Button>
        <Button type="submit">
          保存
        </Button>
      </div>
    </form>
  )
}

// 添加员工表单
function AddStaffForm({
  onSave
}: {
  onSave: (data: any) => void
}) {
  const [formData, setFormData] = useState({
    userId: '',
    title: '',
    permissions: [] as string[],
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(formData)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <Label htmlFor="userId">用户ID</Label>
        <Input
          id="userId"
          type="text"
          value={formData.userId}
          onChange={(e) => setFormData({ ...formData, userId: e.target.value })}
          required
        />
      </div>
      <div>
        <Label htmlFor="title">职位</Label>
        <Input
          id="title"
          type="text"
          value={formData.title}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
        />
      </div>
      <DialogFooter>
        <Button variant="outline" type="button">
          取消
        </Button>
        <Button type="submit">
          添加
        </Button>
      </DialogFooter>
    </form>
  )
}
