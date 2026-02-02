import { useEffect, useState } from 'react'
import { withdrawalApi } from '../services/api'
import type { Withdrawal } from '../types'

export default function WithdrawalReview() {
  const [withdrawals, setWithdrawals] = useState<Withdrawal[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(1)
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
                <p className="text-gray-500">暂无提现记录</p>
              </div>
            )}

            {!loading && withdrawals.length > 0 && (
              <>
                {/* 筛选器 */}
                <div className="mb-4 flex gap-4">
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
                    <option value="approved">已通过</option>
                    <option value="rejected">已拒绝</option>
                    <option value="completed">已完成</option>
                  </select>
                </div>

                {/* 提现列表 */}
                <div className="space-y-4">
                  {withdrawals.map((withdrawal) => {
                    const accountInfo =
                      typeof withdrawal.accountInfo === 'string'
                        ? JSON.parse(withdrawal.accountInfo)
                        : withdrawal.accountInfo

                    return (
                      <div
                        key={withdrawal.id}
                        className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors"
                      >
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            {/* 头部信息 */}
                            <div className="flex items-center gap-4 mb-3">
                              <span className="text-2xl">{getMethodIcon(withdrawal.method)}</span>
                              <div>
                                <div className="font-medium text-gray-900">
                                  {getMethodText(withdrawal.method)}提现
                                </div>
                                <div className="text-sm text-gray-500">
                                  申请时间：{new Date(withdrawal.createdAt).toLocaleString()}
                                </div>
                              </div>
                              <span
                                className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(
                                  withdrawal.status
                                )}`}
                              >
                                {getStatusText(withdrawal.status)}
                              </span>
                            </div>

                            {/* 金额信息 */}
                            <div className="grid grid-cols-3 gap-4 mb-3">
                              <div>
                                <div className="text-xs text-gray-500">申请金额</div>
                                <div className="text-lg font-semibold text-gray-900">
                                  {withdrawal.amount.toLocaleString()} 积分
                                </div>
                              </div>
                              <div>
                                <div className="text-xs text-gray-500">手续费</div>
                                <div className="text-lg font-semibold text-gray-900">
                                  {withdrawal.fee.toLocaleString()} 积分
                                </div>
                              </div>
                              <div>
                                <div className="text-xs text-gray-500">实际到账</div>
                                <div className="text-lg font-semibold text-green-600">
                                  {withdrawal.actualAmount.toLocaleString()} 积分
                                </div>
                              </div>
                            </div>

                            {/* 账户信息 */}
                            <div className="bg-gray-50 rounded p-3 mb-3">
                              <div className="text-xs font-medium text-gray-700 mb-2">账户信息</div>
                              <div className="grid grid-cols-2 gap-2 text-xs">
                                {Object.entries(accountInfo).map(([key, value]) => (
                                  <div key={key}>
                                    <span className="text-gray-500">{key}: </span>
                                    <span className="text-gray-900">{String(value)}</span>
                                  </div>
                                ))}
                              </div>
                            </div>

                            {/* 审核信息 */}
                            {withdrawal.auditedAt && (
                              <div className="text-xs text-gray-500 mb-3">
                                审核时间：{new Date(withdrawal.auditedAt).toLocaleString()}
                                {withdrawal.auditNote && (
                                  <span className="ml-2">备注：{withdrawal.auditNote}</span>
                                )}
                              </div>
                            )}

                            {/* 完成时间 */}
                            {withdrawal.completedAt && (
                              <div className="text-xs text-green-600 mb-3">
                                完成时间：{new Date(withdrawal.completedAt).toLocaleString()}
                              </div>
                            )}
                          </div>

                          {/* 操作按钮 */}
                          <div className="flex flex-col gap-2 ml-4">
                            {withdrawal.status === 'pending' && (
                              <>
                                <button
                                  onClick={() => handleAudit(withdrawal.id, true)}
                                  disabled={processingId === withdrawal.id}
                                  className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50 text-sm"
                                >
                                  {processingId === withdrawal.id ? '处理中...' : '通过'}
                                </button>
                                <button
                                  onClick={() => {
                                    const note = prompt('请输入拒绝原因（可选）')
                                    if (note !== null) {
                                      handleAudit(withdrawal.id, false, note)
                                    }
                                  }}
                                  disabled={processingId === withdrawal.id}
                                  className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50 text-sm"
                                >
                                  拒绝
                                </button>
                              </>
                            )}
                            {withdrawal.status === 'approved' && (
                              <button
                                onClick={() => handleProcess(withdrawal.id)}
                                disabled={processingId === withdrawal.id}
                                className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 text-sm"
                              >
                                {processingId === withdrawal.id ? '处理中...' : '确认打款'}
                              </button>
                            )}
                          </div>
                        </div>
                      </div>
                    )
                  })}
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
