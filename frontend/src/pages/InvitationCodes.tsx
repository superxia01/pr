import { useState, useEffect } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { invitationApi } from '../services/api'
import type { InvitationCode } from '../types'
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
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../components/ui/select'

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

  // DataTable列定义
  const columns = [
    {
      id: 'code',
      header: '邀请码',
      accessorKey: 'code' as const,
      sortable: true,
      cell: ({ value }: { value: string }) => (
        <span className="font-mono text-sm font-medium">{value}</span>
      ),
    },
    {
      id: 'codeType',
      header: '类型',
      accessorKey: 'codeType' as const,
      sortable: true,
      filterable: true,
      cell: ({ value }: { value: string }) => (
        <Badge variant="outline">{getCodeTypeLabel(value)}</Badge>
      ),
    },
    {
      id: 'status',
      header: '状态',
      accessorKey: 'status' as const,
      sortable: true,
      filterable: true,
      cell: ({ value }: { value: string }) => {
        const variant = value === 'active' ? 'default' :
                        value === 'disabled' ? 'secondary' :
                        value === 'expired' ? 'destructive' : 'outline'
        const label = value === 'active' ? '有效' :
                     value === 'disabled' ? '已禁用' :
                     value === 'expired' ? '已过期' : '已用完'
        return <Badge variant={variant}>{label}</Badge>
      },
    },
    {
      id: 'useCount',
      header: '使用情况',
      accessorKey: 'useCount' as const,
      cell: ({ row }: { row: InvitationCode }) => (
        <span className="text-sm text-muted-foreground">
          {row.useCount}/{row.maxUses === 0 ? '无限' : row.maxUses}
        </span>
      ),
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
      cell: ({ row }: { row: InvitationCode }) => (
        row.status === 'active' && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => handleDisableCode(row.code)}
            className="text-destructive hover:text-destructive"
          >
            禁用
          </Button>
        )
      ),
    },
  ]

  return (
    <div className="min-h-screen bg-background">
      {/* 顶部导航 */}
      <nav className="bg-card shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold">PR Business</h1>
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
        <div className="mb-6">
          <h2 className="text-2xl font-bold">邀请码管理</h2>
          <p className="text-muted-foreground mt-1">查看和管理您的邀请码</p>
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="mb-4">
            <Badge variant="destructive" className="w-full py-2 px-3 justify-center">
              {error}
            </Badge>
          </div>
        )}

        {/* 操作栏 */}
        <div className="mb-6 flex justify-between items-center">
          <div className="text-sm text-muted-foreground">
            共 <span className="font-medium">{codes.length}</span> 个邀请码
          </div>
          {canCreateCode() && (
            <Button onClick={() => setShowCreateModal(true)}>
              创建邀请码
            </Button>
          )}
        </div>

        {/* 邀请码列表 */}
        {loading ? (
          <div className="text-center py-12">
            <div className="text-muted-foreground">加载中...</div>
          </div>
        ) : (
          <DataTable
            data={codes}
            columns={columns}
            searchable
            pageSize={10}
            emptyMessage="暂无邀请码"
          />
        )}
      </main>

      {/* 创建邀请码Dialog */}
      <Dialog open={showCreateModal} onOpenChange={setShowCreateModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>创建邀请码</DialogTitle>
          </DialogHeader>

          <form onSubmit={handleCreateCode} className="space-y-4">
            <div className="space-y-2">
              <Label>邀请码类型</Label>
              <Select
                value={createForm.codeType}
                onValueChange={(value) => setCreateForm({ ...createForm, codeType: value })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="MERCHANT">商家</SelectItem>
                  <SelectItem value="CREATOR">达人</SelectItem>
                  <SelectItem value="STAFF">员工</SelectItem>
                  {user?.currentRole === 'super_admin' && (
                    <SelectItem value="SP_ADMIN">服务商管理员</SelectItem>
                  )}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>最大使用次数</Label>
              <Input
                type="number"
                min="1"
                value={createForm.maxUses}
                onChange={(e) => setCreateForm({ ...createForm, maxUses: parseInt(e.target.value) })}
                required
              />
            </div>

            <div className="space-y-2">
              <Label>过期时间（可选）</Label>
              <Input
                type="datetime-local"
                value={createForm.expiresAt}
                onChange={(e) => setCreateForm({ ...createForm, expiresAt: e.target.value })}
              />
            </div>

            <DialogFooter>
              <Button type="submit">创建</Button>
              <Button type="button" variant="outline" onClick={() => setShowCreateModal(false)}>
                取消
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
