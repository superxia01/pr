import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { withdrawalApi } from '../services/api'
import type { Withdrawal } from '../types'

export default function Withdrawals() {
  const navigate = useNavigate()
  const [withdrawals, setWithdrawals] = useState<Withdrawal[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(1)
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
        page,
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
            <button
              onClick={() => navigate('/withdrawal')}
              className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 text-sm font-medium"
            >
              申请提现
            </button>
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

            {!loading && withdrawals.length === 0 && (
              <div className="text-center py-12 bg-gray-50 rounded-lg">
                <p className="text-gray-500 mb-4">暂无提现记录</p>
                <button
                  onClick={() => navigate('/withdrawal')}
                  className="px-6 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
                >
                  立即申请提现
                </button>
              </div>
            )}

            {!loading && withdrawals.length > 0 && (
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
                    <option value="">全部状态</option>
                    <option value="pending">待审核</option>
                    <option value="approved">审核通过</option>
                    <option value="rejected">审核拒绝</option>
                    <option value="completed">已完成</option>
                  </select>
                </div>

                {/* 提现列表 */}
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          申请时间
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          提现方式
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          申请金额
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          手续费
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          实际到账
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          状态
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                          备注
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {withdrawals.map((withdrawal) => (
                        <tr key={withdrawal.id}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {new Date(withdrawal.createdAt).toLocaleString()}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {getMethodText(withdrawal.method)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                            {withdrawal.amount.toLocaleString()}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {withdrawal.fee > 0 ? withdrawal.fee.toLocaleString() : '-'}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-green-600">
                            {withdrawal.actualAmount.toLocaleString()}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm">
                            <span
                              className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(
                                withdrawal.status
                              )}`}
                            >
                              {getStatusText(withdrawal.status)}
                            </span>
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-500">
                            {withdrawal.auditNote || withdrawal.completedAt
                              ? new Date(withdrawal.completedAt || '').toLocaleDateString()
                              : '-'}
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
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                      disabled={page === 1}
                      className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      上一页
                    </button>
                    <span className="text-sm text-gray-600">
                      第 {page} 页，共 {Math.ceil(total / pageSize)} 页
                    </span>
                    <button
                      onClick={() => setPage((p) => p + 1)}
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
