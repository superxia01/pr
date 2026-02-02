import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { creditApi } from '../services/api'
import type { CreditAccount } from '../types'

export default function Recharge() {
  const navigate = useNavigate()
  const [account, setAccount] = useState<CreditAccount | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [amount, setAmount] = useState('')

  // 预设金额选项
  const presetAmounts = [100, 500, 1000, 5000, 10000]

  // 充值方式选项
  const paymentMethods = [
    { id: 'alipay', name: '支付宝', icon: '💙' },
    { id: 'wechat', name: '微信支付', icon: '💚' },
    { id: 'bank', name: '银行转账', icon: '🏦' },
  ]

  const [selectedMethod, setSelectedMethod] = useState(paymentMethods[0].id)

  useEffect(() => {
    loadAccount()
  }, [])

  const loadAccount = async () => {
    setLoading(true)
    try {
      const accountData = await creditApi.getBalance()
      setAccount(accountData)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载账户信息失败')
    } finally {
      setLoading(false)
    }
  }

  const handlePresetAmount = (amt: number) => {
    setAmount(amt.toString())
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    setError(null)

    const rechargeAmount = parseInt(amount)
    if (!rechargeAmount || rechargeAmount <= 0) {
      setError('请输入有效的充值金额')
      setSubmitting(false)
      return
    }

    try {
      const response = await creditApi.recharge({ amount: rechargeAmount })
      alert(`充值成功！当前余额：${response.account.balance} 积分`)
      setAmount('')
      loadAccount()
    } catch (err: any) {
      setError(err.response?.data?.error || '充值失败')
      setSubmitting(false)
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

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-2xl mx-auto">
        {/* 返回按钮 */}
        <button
          onClick={() => navigate(-1)}
          className="mb-4 text-blue-600 hover:text-blue-700 flex items-center gap-2"
        >
          <span>←</span> 返回
        </button>

        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">积分充值</h1>
            <p className="mt-1 text-sm text-gray-500">为您的账户充值积分</p>
          </div>

          {/* 内容 */}
          <div className="p-6">
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
                {error}
              </div>
            )}

            {/* 当前余额 */}
            {account && (
              <div className="bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg p-6 mb-6 text-white">
                <p className="text-sm opacity-90">当前余额</p>
                <p className="text-4xl font-bold mt-2">{account.balance.toLocaleString()}</p>
                <p className="text-sm opacity-90 mt-1">积分</p>
                {account.frozenBalance > 0 && (
                  <p className="text-sm opacity-90 mt-2">
                    冻结余额：{account.frozenBalance.toLocaleString()} 积分
                  </p>
                )}
              </div>
            )}

            {/* 充值表单 */}
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* 预设金额 */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-3">
                  快捷充值
                </label>
                <div className="grid grid-cols-5 gap-2">
                  {presetAmounts.map((amt) => (
                    <button
                      key={amt}
                      type="button"
                      onClick={() => handlePresetAmount(amt)}
                      className={`py-3 px-4 border rounded-md text-sm font-medium transition-colors ${
                        amount === amt.toString()
                          ? 'bg-blue-500 text-white border-blue-500'
                          : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                      }`}
                    >
                      {amt}
                    </button>
                  ))}
                </div>
              </div>

              {/* 自定义金额 */}
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  自定义金额
                </label>
                <div className="mt-1 relative rounded-md shadow-sm">
                  <input
                    type="number"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    className="block w-full rounded-md border-gray-300 pr-12 sm:text-sm border p-2"
                    placeholder="请输入充值金额"
                    min="1"
                  />
                  <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                    <span className="text-gray-500 sm:text-sm">积分</span>
                  </div>
                </div>
              </div>

              {/* 充值方式 */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-3">
                  充值方式 <span className="text-red-500">*</span>
                </label>
                <div className="grid grid-cols-3 gap-4">
                  {paymentMethods.map((method) => (
                    <button
                      key={method.id}
                      type="button"
                      onClick={() => setSelectedMethod(method.id)}
                      className={`p-4 border-2 rounded-lg text-center transition-colors ${
                        selectedMethod === method.id
                          ? 'border-blue-500 bg-blue-50'
                          : 'border-gray-300 hover:border-gray-400'
                      }`}
                    >
                      <div className="text-2xl mb-2">{method.icon}</div>
                      <div className="text-sm font-medium">{method.name}</div>
                    </button>
                  ))}
                </div>
              </div>

              {/* 充值说明 */}
              <div className="bg-gray-50 rounded-lg p-4">
                <h4 className="text-sm font-medium text-gray-900 mb-2">充值说明</h4>
                <ul className="text-xs text-gray-600 space-y-1">
                  <li>• 单笔充值金额最低 100 积分</li>
                  <li>• 充值到账时间为实时</li>
                  <li>• 如有疑问请联系客服</li>
                  <li>• 本服务为模拟充值，实际环境需对接第三方支付</li>
                </ul>
              </div>

              {/* 按钮 */}
              <div className="flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={() => navigate(-1)}
                  className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={submitting || !amount}
                  className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {submitting ? '充值中...' : '确认充值'}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
