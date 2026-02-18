import { useEffect, useState } from 'react'
import { creditApi } from '../services/api'
import type { CreditTransaction, CreditAccount } from '../types'
import { DataTable } from '../components/ui/data-table'
import type { ColumnDef } from '../components/ui/data-table'

export default function CreditTransactions() {
  const [transactions, setTransactions] = useState<CreditTransaction[]>([])
  const [balance, setBalance] = useState<CreditAccount | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(0) // DataTable使用0-based索引
  const [total, setTotal] = useState(0)
  const [filter, setFilter] = useState('')
  const pageSize = 20

  useEffect(() => {
    loadTransactions()
  }, [page, filter])

  useEffect(() => {
    creditApi.getBalance().then(setBalance).catch(() => setBalance(null))
  }, [])

  const loadTransactions = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await creditApi.getTransactions({
        page: page + 1, // API使用1-based索引
        page_size: pageSize,
        type: filter || undefined
      })
      setTransactions(response.data)
      setTotal(response.total || 0)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载流水失败')
    } finally {
      setLoading(false)
    }
  }

  const getTransactionName = (type: string) => {
    const names: Record<string, string> = {
      'RECHARGE': '充值',
      'TASK_INCOME': '任务收入',
      'TASK_SUBMIT': '任务提交',
      'STAFF_REFERRAL': '员工返佣',
      'PROVIDER_INCOME': '服务商收入',
      'TASK_PUBLISH': '发布任务',
      'TASK_ACCEPT': '接任务',
      'TASK_REJECT': '审核拒绝',
      'TASK_ESCALATE': '超时拒绝',
      'TASK_REFUND': '任务退款',
      'WITHDRAW': '提现',
      'WITHDRAW_FREEZE': '提现冻结',
      'WITHDRAW_REFUND': '提现拒绝',
      'BONUS_GIFT': '系统赠送',
    }
    return names[type] || type
  }

  const getAmountColor = (amount: number) => {
    return amount > 0 ? 'text-green-600' : 'text-red-600'
  }

  const getAmountPrefix = (amount: number) => {
    return amount > 0 ? '+' : ''
  }

  // 定义表格列
  const columns: ColumnDef<CreditTransaction>[] = [
    {
      id: 'createdAt',
      header: '时间',
      accessorKey: 'createdAt',
      cell: ({ row }) => (
        <span className="text-sm text-gray-500">
          {new Date(row.createdAt).toLocaleString()}
        </span>
      )
    },
    {
      id: 'type',
      header: '类型',
      accessorKey: 'type',
      cell: ({ row }) => (
        <span className="text-sm font-medium text-gray-900">
          {getTransactionName(row.type)}
        </span>
      )
    },
    {
      id: 'amount',
      header: '金额',
      accessorKey: 'amount',
      cell: ({ row }) => (
        <span className={`text-sm font-medium ${getAmountColor(row.amount)}`}>
          {getAmountPrefix(row.amount)}
          {row.amount.toLocaleString()}
        </span>
      )
    },
    {
      id: 'balanceAfter',
      header: '余额后',
      accessorKey: 'balanceAfter',
      cell: ({ row }) => (
        <span className="text-sm text-gray-900">
          {row.balanceAfter.toLocaleString()}
        </span>
      )
    },
    {
      id: 'description',
      header: '说明',
      accessorKey: 'description',
      cell: ({ row }) => (
        <span className="text-sm text-gray-500">
          {row.description || '-'}
        </span>
      )
    },
  ]

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">积分明细</h1>
            <p className="mt-1 text-sm text-gray-500">查看您的积分流水记录</p>
          </div>

          {/* 当前余额 */}
          {balance != null && (
            <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
              <div className="flex items-baseline gap-4">
                <span className="text-sm text-gray-600">当前余额</span>
                <span className="text-2xl font-bold text-gray-900">{balance.balance.toLocaleString()}</span>
                <span className="text-sm text-gray-500">积分</span>
                {balance.frozenBalance > 0 && (
                  <span className="text-sm text-gray-500">冻结：{balance.frozenBalance.toLocaleString()}</span>
                )}
              </div>
            </div>
          )}

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
                <option value="">全部类型</option>
                <option value="RECHARGE">充值</option>
                <option value="TASK_INCOME">任务收入</option>
                <option value="TASK_PUBLISH">发布任务</option>
                <option value="TASK_ACCEPT">接任务</option>
                <option value="WITHDRAW">提现</option>
              </select>
            </div>

            {/* 表格 */}
            <DataTable
              columns={columns}
              data={transactions}
              serverSide={true}
              total={total}
              pageSize={pageSize}
              onPageChange={(pageIndex) => setPage(pageIndex)}
              loading={loading}
              emptyMessage="暂无流水记录"
              searchable={false}
            />
          </div>
        </div>
      </div>
    </div>
  )
}
