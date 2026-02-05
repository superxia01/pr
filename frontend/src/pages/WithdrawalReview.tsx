import { useEffect, useState } from 'react'
import { withdrawalApi } from '../services/api'
import type { Withdrawal } from '../types'
import { DataTable } from '../components/ui/data-table'
import type { ColumnDef } from '../components/ui/data-table'
import { Button } from '../components/ui/button'

export default function WithdrawalReview() {
  const [withdrawals, setWithdrawals] = useState<Withdrawal[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(0) // DataTable使用0-based索引
  const [total, setTotal] = useState(0)
  const [filter, setFilter] = useState('')
  const [processingId, setProcessingId] = useState<string | null>(null)
  const pageSize = 20

  useEffect(() => {
    loadWithdrawals()
  }, [page, filter])

  const loadWithdrawals = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await withdrawalApi.getWithdrawals({
        page: page + 1, // API使用1-based索引
        page_size: pageSize,
        status: filter || undefined,
      })
      setWithdrawals(response.withdrawals)
      setTotal(response.total)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载提现记录失败')
    } finally {
      setLoading(false)
    }
  }

  const handleAudit = async (id: string, approved: boolean, auditNote?: string) => {
    setProcessingId(id)
    try {
      await withdrawalApi.auditWithdrawal(id, { approved, auditNote })
      alert(approved ? '已通过审核' : '已拒绝申请')
      loadWithdrawals()
    } catch (err: any) {
      alert(err.response?.data?.error || '操作失败')
    } finally {
      setProcessingId(null)
    }
  }

  const handleProcess = async (id: string) => {
    if (!confirm('确认已完成打款？')) return

    setProcessingId(id)
    try {
      await withdrawalApi.processWithdrawal(id)
      alert('打款成功')
      loadWithdrawals()
    } catch (err: any) {
      alert(err.response?.data?.error || '操作失败')
    } finally {
      setProcessingId(null)
    }
  }

  const getStatusText = (status: string) => {
    const statusMap: Record<string, string> = {
      pending: '待审核',
      approved: '已通过',
      rejected: '已拒绝',
      completed: '已完成',
    }
    return statusMap[status] || status
  }

  const getStatusColor = (status: string) => {
    const colorMap: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800',
      approved: 'bg-blue-100 text-blue-800',
      rejected: 'bg-red-100 text-red-800',
      completed: 'bg-green-100 text-green-800',
    }
    return colorMap[status] || 'bg-gray-100 text-gray-800'
  }

  const getMethodText = (method: string) => {
    const methodMap: Record<string, string> = {
      ALIPAY: '支付宝',
      WECHAT: '微信',
      BANK: '银行转账',
    }
    return methodMap[method] || method
  }

  const getMethodIcon = (method: string) => {
    const iconMap: Record<string, string> = {
      ALIPAY: '💙',
      WECHAT: '💚',
      BANK: '🏦',
    }
    return iconMap[method] || '💳'
  }

  // 定义表格列
  const columns: ColumnDef<Withdrawal>[] = [
    {
      id: 'info',
      header: '提现信息',
      accessorKey: 'method',
      cell: ({ row }) => {
        const accountInfo =
          typeof row.accountInfo === 'string'
            ? JSON.parse(row.accountInfo)
            : row.accountInfo

        return (
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <span className="text-xl">{getMethodIcon(row.method)}</span>
              <div>
                <div className="font-medium text-gray-900">
                  {getMethodText(row.method)}提现
                </div>
                <div className="text-xs text-gray-500">
                  {new Date(row.createdAt).toLocaleString()}
                </div>
              </div>
            </div>
            <div className="text-xs bg-gray-50 rounded p-2">
              <div className="font-medium text-gray-700 mb-1">账户信息</div>
              {Object.entries(accountInfo).map(([key, value]) => (
                <div key={key} className="text-gray-600">
                  <span className="text-gray-400">{key}:</span> {String(value)}
                </div>
              ))}
            </div>
          </div>
        )
      }
    },
    {
      id: 'amount',
      header: '金额信息',
      accessorKey: 'amount',
      cell: ({ row }) => (
        <div className="text-sm space-y-1">
          <div>
            <span className="text-gray-500">申请:</span>{' '}
            <span className="font-semibold">{row.amount.toLocaleString()}</span>
          </div>
          <div>
            <span className="text-gray-500">手续费:</span>{' '}
            <span className="font-semibold">{row.fee.toLocaleString()}</span>
          </div>
          <div>
            <span className="text-gray-500">到账:</span>{' '}
            <span className="font-semibold text-green-600">{row.actualAmount.toLocaleString()}</span>
          </div>
        </div>
      )
    },
    {
      id: 'status',
      header: '状态',
      accessorKey: 'status',
      cell: ({ row }) => (
        <div className="space-y-1">
          <span
            className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(row.status)}`}
          >
            {getStatusText(row.status)}
          </span>
          {row.auditedAt && (
            <div className="text-xs text-gray-500">
              审核: {new Date(row.auditedAt).toLocaleDateString()}
            </div>
          )}
          {row.completedAt && (
            <div className="text-xs text-green-600">
              完成: {new Date(row.completedAt).toLocaleDateString()}
            </div>
          )}
        </div>
      )
    },
    {
      id: 'actions',
      header: '操作',
      accessorKey: 'status',
      cell: ({ row }) => (
        <div className="flex flex-col gap-2">
          {row.status === 'pending' && (
            <>
              <Button
                size="sm"
                onClick={() => handleAudit(row.id, true)}
                disabled={processingId === row.id}
                className="bg-green-600 hover:bg-green-700"
              >
                {processingId === row.id ? '处理中...' : '通过'}
              </Button>
              <Button
                size="sm"
                variant="destructive"
                onClick={() => {
                  const note = prompt('请输入拒绝原因（可选）')
                  if (note !== null) {
                    handleAudit(row.id, false, note)
                  }
                }}
                disabled={processingId === row.id}
              >
                拒绝
              </Button>
            </>
          )}
          {row.status === 'approved' && (
            <Button
              size="sm"
              onClick={() => handleProcess(row.id)}
              disabled={processingId === row.id}
              className="bg-blue-600 hover:bg-blue-700"
            >
              {processingId === row.id ? '处理中...' : '确认打款'}
            </Button>
          )}
          {row.auditNote && (
            <div className="text-xs text-gray-500 mt-1">
              备注: {row.auditNote}
            </div>
          )}
        </div>
      )
    },
  ]

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">提现审核管理</h1>
            <p className="mt-1 text-sm text-gray-500">审核和处理用户的提现申请</p>
          </div>

          {/* 内容 */}
          <div className="p-6">
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
                {error}
              </div>
            )}

            {/* 筛选器 */}
            <div className="mb-4">
              <select
                value={filter}
                onChange={(e) => {
                  setFilter(e.target.value)
                  setPage(0) // 重置到第一页
                }}
                className="px-4 py-2 border border-gray-300 rounded-md text-sm"
              >
                <option value="">全部状态</option>
                <option value="pending">待审核</option>
                <option value="approved">已通过</option>
                <option value="rejected">已拒绝</option>
                <option value="completed">已完成</option>
              </select>
            </div>

            {/* 表格 */}
            <DataTable
              columns={columns}
              data={withdrawals}
              serverSide={true}
              total={total}
              pageSize={pageSize}
              onPageChange={(pageIndex) => setPage(pageIndex)}
              loading={loading}
              emptyMessage="暂无提现记录"
              searchable={false}
            />
          </div>
        </div>
      </div>
    </div>
  )
}
