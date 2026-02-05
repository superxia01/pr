import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { withdrawalApi } from '../services/api'
import type { Withdrawal } from '../types'
import { DataTable } from '../components/ui/data-table'
import type { ColumnDef } from '../components/ui/data-table'
import { Button } from '../components/ui/button'

export default function Withdrawals() {
  const navigate = useNavigate()
  const [withdrawals, setWithdrawals] = useState<Withdrawal[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(0) // DataTable使用0-based索引
  const [total, setTotal] = useState(0)
  const [filter, setFilter] = useState('')
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

  const getStatusText = (status: string) => {
    const statusMap: Record<string, string> = {
      pending: '待审核',
      approved: '审核通过',
      rejected: '审核拒绝',
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
      WECHAT: '微信支付',
      BANK: '银行转账',
    }
    return methodMap[method] || method
  }

  // 定义表格列
  const columns: ColumnDef<Withdrawal>[] = [
    {
      id: 'createdAt',
      header: '申请时间',
      accessorKey: 'createdAt',
      cell: ({ row }) => (
        <span className="text-sm text-gray-500">
          {new Date(row.createdAt).toLocaleString()}
        </span>
      )
    },
    {
      id: 'method',
      header: '提现方式',
      accessorKey: 'method',
      cell: ({ row }) => (
        <span className="text-sm text-gray-900">
          {getMethodText(row.method)}
        </span>
      )
    },
    {
      id: 'amount',
      header: '申请金额',
      accessorKey: 'amount',
      cell: ({ row }) => (
        <span className="text-sm text-gray-900">
          {row.amount.toLocaleString()}
        </span>
      )
    },
    {
      id: 'fee',
      header: '手续费',
      accessorKey: 'fee',
      cell: ({ row }) => (
        <span className="text-sm text-gray-500">
          {row.fee > 0 ? row.fee.toLocaleString() : '-'}
        </span>
      )
    },
    {
      id: 'actualAmount',
      header: '实际到账',
      accessorKey: 'actualAmount',
      cell: ({ row }) => (
        <span className="text-sm font-medium text-green-600">
          {row.actualAmount.toLocaleString()}
        </span>
      )
    },
    {
      id: 'status',
      header: '状态',
      accessorKey: 'status',
      cell: ({ row }) => (
        <span
          className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(row.status)}`}
        >
          {getStatusText(row.status)}
        </span>
      )
    },
    {
      id: 'auditNote',
      header: '备注',
      accessorKey: 'auditNote',
      cell: ({ row }) => (
        <span className="text-sm text-gray-500">
          {row.auditNote || row.completedAt
            ? new Date(row.completedAt || '').toLocaleDateString()
            : '-'}
        </span>
      )
    },
  ]

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">提现记录</h1>
              <p className="mt-1 text-sm text-gray-500">查看您的提现申请记录</p>
            </div>
            <Button
              onClick={() => navigate('/withdrawal')}
              className="bg-green-600 hover:bg-green-700"
            >
              申请提现
            </Button>
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
                <option value="approved">审核通过</option>
                <option value="rejected">审核拒绝</option>
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
