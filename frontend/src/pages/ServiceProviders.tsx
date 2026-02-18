import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { serviceProviderApi } from '../services/api'
import type { ServiceProvider } from '../types'
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

export default function ServiceProviders() {
  const { user } = useAuth()
  const [providers, setProviders] = useState<ServiceProvider[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    loadProviders()
  }, [])

  const loadProviders = async () => {
    try {
      setLoading(true)
      const data = await serviceProviderApi.getServiceProviders()
      setProviders(data)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载服务商列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateProvider = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')

    const formData = new FormData(e.currentTarget)
    const data = {
      name: formData.get('name') as string,
      description: formData.get('description') as string,
    }

    try {
      await serviceProviderApi.createServiceProvider(data)
      setShowCreateModal(false)
      loadProviders()
      // 重置表单
      e.currentTarget.reset()
    } catch (err: any) {
      setError(err.response?.data?.error || '创建服务商失败')
    }
  }

  const handleDeleteProvider = async (id: string) => {
    if (!confirm('确定要删除此服务商吗？')) return

    try {
      await serviceProviderApi.deleteServiceProvider(id)
      loadProviders()
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
      header: '服务商名称',
      accessorKey: 'name' as const,
      sortable: true,
      cell: ({ value }: { value: string }) => (
        <div className="flex items-center gap-3">
          <span className="font-medium">{value}</span>
        </div>
      ),
    },
    {
      id: 'description',
      header: '描述',
      accessorKey: 'description' as const,
      cell: ({ value }: { value: string }) => (
        <span className="text-sm text-muted-foreground">
          {value || '-'}
        </span>
      ),
    },
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
      id: 'merchants',
      header: '商家数量',
      accessorKey: 'merchants' as const,
      cell: ({ value }: { value: any[] | undefined }) => (
        <span className="text-sm font-medium">
          {value?.length || 0}
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
      cell: ({ row }: { row: ServiceProvider }) => (
        <div className="flex gap-2">
          <Button variant="ghost" size="sm" asChild>
            <Link to={`/service-provider/${row.id}`}>查看详情</Link>
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => handleDeleteProvider(row.id)}
            className="text-destructive hover:text-destructive"
          >
            删除
          </Button>
        </div>
      ),
    },
  ]

  // 超级管理员和服务商管理员可以访问
  const canAccess = user?.roles.some(r =>
    r.toUpperCase() === 'SUPER_ADMIN' || r.toUpperCase() === 'SERVICE_PROVIDER_ADMIN'
  )

  if (!canAccess) {
    return (
      <div className="min-h-screen bg-background p-8">
        <div className="max-w-7xl mx-auto">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            您不是服务商管理员
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
              <h1 className="text-xl font-bold">PR Business - 服务商管理</h1>
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
            <h2 className="text-2xl font-bold">服务商管理</h2>
            <p className="text-muted-foreground mt-1">管理所有服务商组织</p>
          </div>
          <Button onClick={() => setShowCreateModal(true)}>
            创建服务商
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

        {/* 统计信息 */}
        <div className="mb-6 grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-card rounded-lg border p-4">
            <div className="text-sm text-muted-foreground">总服务商数</div>
            <div className="text-2xl font-bold mt-1">{providers.length}</div>
          </div>
          <div className="bg-card rounded-lg border p-4">
            <div className="text-sm text-muted-foreground">活跃服务商</div>
            <div className="text-2xl font-bold mt-1">
              {providers.filter(p => p.status === 'active').length}
            </div>
          </div>
          <div className="bg-card rounded-lg border p-4">
            <div className="text-sm text-muted-foreground">关联商家总数</div>
            <div className="text-2xl font-bold mt-1">
              {providers.reduce((sum, p) => sum + (p.merchants?.length || 0), 0)}
            </div>
          </div>
        </div>

        {/* 服务商列表 */}
        {loading ? (
          <div className="text-center py-12">
            <div className="text-muted-foreground">加载中...</div>
          </div>
        ) : (
          <DataTable
            data={providers}
            columns={columns}
            searchable
            pageSize={10}
            emptyMessage="暂无服务商"
          />
        )}
      </main>

      {/* 创建服务商Dialog */}
      <Dialog open={showCreateModal} onOpenChange={setShowCreateModal}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>创建服务商</DialogTitle>
          </DialogHeader>

          <form onSubmit={handleCreateProvider} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">服务商名称 *</Label>
              <Input
                id="name"
                name="name"
                placeholder="请输入服务商名称"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">服务商描述</Label>
              <Textarea
                id="description"
                name="description"
                placeholder="请输入服务商描述"
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
