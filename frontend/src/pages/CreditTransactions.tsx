import { useEffect, useState } from 'react'
import { creditApi } from '../services/api'
import type { CreditTransaction } from '../types'

export default function CreditTransactions() {
  const [transactions, setTransactions] = useState<CreditTransaction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [filter, setFilter] = useState('')
  const pageSize = 20

  useEffect(() => {
    loadTransactions()
  }, [page, filter])

  const loadTransactions = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await creditApi.getTransactions({ page, page_size: pageSize })
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

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">积分明细</h1>
            <p className="mt-1 text-sm text-gray-500">查看您的积分流水记录</p>
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
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
                {error}
              </div>
            )}

            {!loading && transactions.length === 0 && (
              <div className="text-center py-12 bg-gray-50 rounded-lg">
                <p className="text-gray-500">暂无流水记录</p>
              </div>
            )}

            {!loading && transactions.length > 0 && (
              <>
                {/* 筛选器 */}
                <div className="mb-4">
                  <select
                    value={filter}
                    onChange={(e) => {
                      setFilter(e.target.value)
                      setPage(1)
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

                {/* 流水列表 */}
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          时间
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          类型
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          金额
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          余额后
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          说明
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {transactions.map((txn) => (
                        <tr key={txn.id}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {new Date(txn.createdAt).toLocaleString()}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {getTransactionName(txn.type)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                            <span className={getAmountColor(txn.amount)}>
                              {getAmountPrefix(txn.amount)}
                              {txn.amount.toLocaleString()}
                            </span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {txn.balanceAfter.toLocaleString()}
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-500">
                            {txn.description || '-'}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {/* 分页 */}
                {total > pageSize && (
                  <div className="mt-6 flex justify-center items-center gap-4">
                    <button
                      onClick={() => setPage(p => Math.max(1, p - 1))}
                      disabled={page === 1}
                      className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      上一页
                    </button>
                    <span className="text-sm text-gray-600">
                      第 {page} 页，共 {Math.ceil(total / pageSize)} 页
                    </span>
                    <button
                      onClick={() => setPage(p => p + 1)}
                      disabled={page >= Math.ceil(total / pageSize)}
                      className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      下一页
                    </button>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
