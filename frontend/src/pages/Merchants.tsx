import { useState, useEffect } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { merchantApi, serviceProviderApi } from '../services/api'
import type { Merchant, ServiceProvider, User } from '../types'
import { DataTable } from '../components/ui/data-table'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Label } from '../components/ui/label'
import { Badge } from '../components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '../components/ui/dialog'
import { Textarea } from '../components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../components/ui/select'

// 辅助函数：创建完整 User 对象
const createMockUser = (id: string, nickname: string): User => ({
  id,
  nickname,
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
})

export default function Merchants() {
  const { user } = useAuth()
  const [merchants, setMerchants] = useState<Merchant[]>([])
  const [serviceProviders, setServiceProviders] = useState<ServiceProvider[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    // 🔴 开发模式：使用模拟数据
    if (import.meta.env.DEV) {
      setTimeout(() => {
        const mockMerchants: any[] = [
          {
            id: 'mch_001',
            name: '测试商家A',
            description: '这是第一个测试商家',
            providerId: 'sp_001',
            adminId: 'usr_001',
            status: 'active',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            provider: {
              id: 'sp_001',
              name: '测试服务商A',
              adminId: 'usr_001',
              userId: 'usr_001',
              description: '测试服务商',
              status: 'active',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
            admin: createMockUser('usr_001', '管理员A'),
          },
          {
            id: 'mch_002',
            name: '测试商家B',
            description: '这是第二个测试商家',
            providerId: 'sp_001',
            adminId: 'usr_002',
            status: 'active',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            provider: {
              id: 'sp_001',
              name: '测试服务商A',
              adminId: 'usr_001',
              userId: 'usr_001',
              description: '测试服务商',
              status: 'active',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
            admin: createMockUser('usr_002', '管理员B'),
          },
        ]

        const mockProviders: any[] = [
          {
            id: 'sp_001',
            name: '测试服务商A',
            adminId: 'usr_001',
            userId: 'usr_001',
            description: '测试服务商',
            status: 'active',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            admin: createMockUser('usr_001', '管理员A'),
          },
        ]

        setMerchants(mockMerchants)
        if (user?.currentRole === 'super_admin') {
          setServiceProviders(mockProviders)
        }
        setLoading(false)
        console.log('🔴 开发模式：使用模拟商家数据')
      }, 500)
      return
    }

    // 生产模式：调用真实 API
    try {
      setLoading(true)
      // 根据角色决定加载哪些商家
      let merchantData: Merchant[]
      if (user?.currentRole === 'super_admin') {
        // 超管看所有商家
        merchantData = await merchantApi.getMerchants()
      } else if (user?.currentRole === 'provider_admin') {
        // 服务商管理员只看自己服务商下的商家
        // 先获取服务商信息
        try {
          const provider = await serviceProviderApi.getMyServiceProvider()
          merchantData = await merchantApi.getMerchants({ provider_id: provider.id })
        } catch {
          // 如果没有绑定服务商，返回空列表
          merchantData = []
        }
      } else {
        merchantData = []
      }
      setMerchants(merchantData)

      // 如果是超管，还需要加载服务商列表用于创建时选择
      if (user?.currentRole === 'super_admin') {
        const providerData = await serviceProviderApi.getServiceProviders()
        setServiceProviders(providerData)
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '加载商家列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateMerchant = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')

    const formData = new FormData(e.currentTarget)
    const data: any = {
      name: formData.get('name') as string,
      description: formData.get('description') as string,
      userId: formData.get('userId') as string || user?.id, // 默认使用当前用户
    }

    // 如果是服务商管理员，自动使用自己绑定的服务商
    if (user?.currentRole === 'provider_admin') {
      try {
        const provider = await serviceProviderApi.getMyServiceProvider()
        data.providerId = provider.id
      } catch {
        setError('请先绑定服务商')
        return
      }
    } else {
      // 超管需要选择服务商
      const providerId = formData.get('providerId') as string
      if (!providerId) {
        setError('请选择服务商')
        return
      }
      data.providerId = providerId
    }

    try {
      await merchantApi.createMerchant(data)
      setShowCreateModal(false)
      loadData()
      e.currentTarget.reset()
    } catch (err: any) {
      setError(err.response?.data?.error || '创建商家失败')
    }
  }

  const handleDeleteMerchant = async (id: string) => {
    if (!confirm('确定要删除此商家吗？')) return

    try {
      await merchantApi.deleteMerchant(id)
      loadData()
    } catch (err: any) {
      setError(err.response?.data?.error || '删除失败')
    }
  }

  const getStatusBadge = (status: string) => {
    const variant = status === 'active' ? 'default' :
                    status === 'suspended' ? 'destructive' :
                    'secondary'
    const label = status === 'active' ? '正常' :
                 status === 'suspended' ? '暂停' :
                 status === 'inactive' ? '未激活' : status
    return <Badge variant={variant}>{label}</Badge>
  }

  // DataTable列定义
  const columns = [
    {
      id: 'name',
      header: '商家名称',
      accessorKey: 'name' as const,
      sortable: true,
      cell: ({ value }: { value: string }) => (
        <div className="flex items-center gap-3">
          <span className="font-medium">{value}</span>
        </div>
      ),
    },
    ...(user?.currentRole === 'super_admin' ? [{
      id: 'provider' as const,
      header: '所属服务商',
      accessorKey: 'provider' as const,
      cell: ({ value }: { value: any }) => (
        <span className="text-sm">
          {value?.name || '-'}
        </span>
      ),
    }] : []),
    {
      id: 'admin',
      header: '管理员',
      accessorKey: 'admin' as const,
      cell: ({ value }: { value: any }) => (
        <span className="text-sm">
          {value?.nickname || value?.id || '-'}
        </span>
      ),
    },
    {
      id: 'status',
      header: '状态',
      accessorKey: 'status' as const,
      sortable: true,
      filterable: true,
      cell: ({ value }: { value: string }) => getStatusBadge(value),
    },
    {
      id: 'createdAt',
      header: '创建时间',
      accessorKey: 'createdAt' as const,
      sortable: true,
      cell: ({ value }: { value: string }) => (
        <span className="text-sm text-muted-foreground">
          {new Date(value).toLocaleDateString('zh-CN')}
        </span>
      ),
    },
    {
      id: 'actions',
      header: '操作',
      cell: ({ row }: { row: Merchant }) => (
        <div className="flex gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => window.location.href = `/merchant/${row.id}`}
          >
            查看详情
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => handleDeleteMerchant(row.id)}
            className="text-destructive hover:text-destructive"
          >
            删除
          </Button>
        </div>
      ),
    },
  ]

  // 权限检查
  const canAccess = user?.currentRole === 'super_admin' ||
                    user?.currentRole === 'provider_admin'

  if (!canAccess) {
    return (
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-7xl mx-auto">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            您没有权限访问此页面
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* 顶部导航 */}
      <nav className="bg-card shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold">
                PR Business - {user?.currentRole === 'super_admin' ? '商家管理' : '我的商家'}
              </h1>
            </div>
            <div className="flex items-center gap-4">
              <a href="/dashboard" className="text-sm text-muted-foreground hover:text-foreground">
                返回工作台
              </a>
            </div>
          </div>
        </div>
      </nav>

      {/* 主要内容 */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6 flex justify-between items-center">
          <div>
            <h2 className="text-2xl font-bold">
              {user?.currentRole === 'super_admin' ? '商家管理' : '我的商家'}
            </h2>
            <p className="text-muted-foreground mt-1">
              {user?.currentRole === 'super_admin' ? '管理所有商家' : '管理您服务商下的商家'}
            </p>
          </div>
          <Button onClick={() => setShowCreateModal(true)}>
            创建商家
          </Button>
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="mb-4">
            <Badge variant="destructive" className="w-full py-2 px-3 justify-center">
              {error}
            </Badge>
          </div>
        )}

        {/* 商家列表 */}
        {loading ? (
          <div className="text-center py-12">
            <div className="text-muted-foreground">加载中...</div>
          </div>
        ) : (
          <DataTable
            data={merchants}
            columns={columns}
            searchable
            pageSize={10}
            emptyMessage="暂无商家"
          />
        )}
      </main>

      {/* 创建商家Dialog */}
      <Dialog open={showCreateModal} onOpenChange={setShowCreateModal}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>创建商家</DialogTitle>
          </DialogHeader>

          <form onSubmit={handleCreateMerchant} className="space-y-4">
            {/* 隐藏字段：商家管理员用户ID，默认当前用户 */}
            <input
              type="hidden"
              name="userId"
              value={user?.id || ''}
            />

            <div className="space-y-2">
              <Label htmlFor="name">商家名称 *</Label>
              <Input
                id="name"
                name="name"
                placeholder="请输入商家名称"
                required
              />
            </div>

            {/* 只有超管需要选择服务商 */}
            {user?.currentRole === 'super_admin' && (
              <div className="space-y-2">
                <Label htmlFor="providerId">所属服务商 *</Label>
                <Select name="providerId" required>
                  <SelectTrigger className="bg-background">
                    <SelectValue placeholder="请选择服务商" />
                  </SelectTrigger>
                  <SelectContent className="bg-popover border-border">
                    {serviceProviders.map((provider) => (
                      <SelectItem key={provider.id} value={provider.id}>
                        {provider.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="description">商家描述</Label>
              <Textarea
                id="description"
                name="description"
                placeholder="请输入商家描述"
                rows={3}
              />
            </div>

            <DialogFooter>
              <Button type="submit">创建</Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => setShowCreateModal(false)}
              >
                取消
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
