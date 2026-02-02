import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { creditApi, withdrawalApi } from '../services/api'
import type { CreditAccount } from '../types'

export default function Withdrawal() {
  const navigate = useNavigate()
  const [account, setAccount] = useState<CreditAccount | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [amount, setAmount] = useState('')
  const [method, setMethod] = useState<'ALIPAY' | 'WECHAT' | 'BANK'>('ALIPAY')
  const [accountInfo, setAccountInfo] = useState<Record<string, any>>({})

  // 预设金额选项
  const presetAmounts = [100, 500, 1000, 5000, 10000]

  // 提现方式配置
  const withdrawalMethods = [
    {
      id: 'ALIPAY' as const,
      name: '支付宝',
      icon: '💙',
      fields: [
        { name: 'account', label: '支付宝账号', placeholder: '请输入支付宝账号或手机号', required: true },
        { name: 'name', label: '真实姓名', placeholder: '请输入真实姓名', required: true },
      ],
    },
    {
      id: 'WECHAT' as const,
      name: '微信支付',
      icon: '💚',
      fields: [
        { name: 'wechatId', label: '微信号', placeholder: '请输入微信号', required: true },
        { name: 'name', label: '真实姓名', placeholder: '请输入真实姓名', required: true },
      ],
    },
    {
      id: 'BANK' as const,
      name: '银行转账',
      icon: '🏦',
      fields: [
        { name: 'bankName', label: '开户银行', placeholder: '如：中国工商银行', required: true },
        { name: 'account', label: '银行卡号', placeholder: '请输入银行卡号', required: true },
        { name: 'name', label: '持卡人姓名', placeholder: '请输入持卡人姓名', required: true },
        { name: 'branch', label: '开户行支行', placeholder: '如：北京分行朝阳支行', required: false },
      ],
    },
  ]

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

  const handleAccountInfoChange = (field: string, value: string) => {
    setAccountInfo((prev) => ({ ...prev, [field]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    setError(null)

    const withdrawalAmount = parseInt(amount)
    if (!withdrawalAmount || withdrawalAmount <= 0) {
      setError('请输入有效的提现金额')
      setSubmitting(false)
      return
    }

    // 验证账户信息
    const currentMethod = withdrawalMethods.find((m) => m.id === method)
    if (!currentMethod) {
      setError('无效的提现方式')
      setSubmitting(false)
      return
    }

    for (const field of currentMethod.fields) {
      if (field.required && !accountInfo[field.name]) {
        setError(`请填写${field.label}`)
        setSubmitting(false)
        return
      }
    }

    try {
      await withdrawalApi.createWithdrawal({
        amount: withdrawalAmount,
        method,
        accountInfo,
      })
      alert('提现申请已提交，等待审核')
      navigate('/withdrawals')
    } catch (err: any) {
      setError(err.response?.data?.error || '提现申请失败')
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

  const currentMethod = withdrawalMethods.find((m) => m.id === method)
  const availableBalance = account?.balance || 0

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
            <h1 className="text-2xl font-bold text-gray-900">申请提现</h1>
            <p className="mt-1 text-sm text-gray-500">提取您的积分到账户</p>
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
              <div className="bg-gradient-to-r from-green-500 to-green-600 rounded-lg p-6 mb-6 text-white">
                <p className="text-sm opacity-90">可提现余额</p>
                <p className="text-4xl font-bold mt-2">{availableBalance.toLocaleString()}</p>
                <p className="text-sm opacity-90 mt-1">积分</p>
                {account.frozenBalance > 0 && (
                  <p className="text-sm opacity-90 mt-2">
                    冻结余额：{account.frozenBalance.toLocaleString()} 积分（不可提现）
                  </p>
                )}
              </div>
            )}

            {/* 提现表单 */}
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* 预设金额 */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-3">快捷提现</label>
                <div className="grid grid-cols-5 gap-2">
                  {presetAmounts.map((amt) => (
                    <button
                      key={amt}
                      type="button"
                      onClick={() => handlePresetAmount(amt)}
                      disabled={amt > availableBalance}
                      className={`py-3 px-4 border rounded-md text-sm font-medium transition-colors ${
                        amount === amt.toString()
                          ? 'bg-green-500 text-white border-green-500'
                          : amt > availableBalance
                          ? 'bg-gray-100 text-gray-400 border-gray-300 cursor-not-allowed'
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
                <label className="block text-sm font-medium text-gray-700">提现金额</label>
                <div className="mt-1 relative rounded-md shadow-sm">
                  <input
                    type="number"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    className="block w-full rounded-md border-gray-300 pr-12 sm:text-sm border p-2"
                    placeholder="请输入提现金额"
                    min="1"
                    max={availableBalance}
                  />
                  <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                    <span className="text-gray-500 sm:text-sm">积分</span>
                  </div>
                </div>
                <p className="mt-1 text-xs text-gray-500">可提现金额：{availableBalance.toLocaleString()} 积分</p>
              </div>

              {/* 提现方式 */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-3">
                  提现方式 <span className="text-red-500">*</span>
                </label>
                <div className="grid grid-cols-3 gap-4">
                  {withdrawalMethods.map((methodInfo) => (
                    <button
                      key={methodInfo.id}
                      type="button"
                      onClick={() => {
                        setMethod(methodInfo.id)
                        setAccountInfo({})
                      }}
                      className={`p-4 border-2 rounded-lg text-center transition-colors ${
                        method === methodInfo.id
                          ? 'border-green-500 bg-green-50'
                          : 'border-gray-300 hover:border-gray-400'
                      }`}
                    >
                      <div className="text-2xl mb-2">{methodInfo.icon}</div>
                      <div className="text-sm font-medium">{methodInfo.name}</div>
                    </button>
                  ))}
                </div>
              </div>

              {/* 账户信息 */}
              {currentMethod && (
                <div className="space-y-4">
                  <label className="block text-sm font-medium text-gray-700">
                    账户信息 <span className="text-red-500">*</span>
                  </label>
                  {currentMethod.fields.map((field) => (
                    <div key={field.name}>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        {field.label}
                        {field.required && <span className="text-red-500">*</span>}
                      </label>
                      <input
                        type="text"
                        value={accountInfo[field.name] || ''}
                        onChange={(e) => handleAccountInfoChange(field.name, e.target.value)}
                        className="block w-full rounded-md border-gray-300 sm:text-sm border p-2"
                        placeholder={field.placeholder}
                        required={field.required}
                      />
                    </div>
                  ))}
                </div>
              )}

              {/* 提现说明 */}
              <div className="bg-gray-50 rounded-lg p-4">
                <h4 className="text-sm font-medium text-gray-900 mb-2">提现说明</h4>
                <ul className="text-xs text-gray-600 space-y-1">
                  <li>• 单笔提现金额最低 100 积分</li>
                  <li>• 提现申请提交后，积分将被冻结，等待审核</li>
                  <li>• 审核通过后1-3个工作日内打款</li>
                  <li>• 如有疑问请联系客服</li>
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
                  disabled={submitting || !amount || parseInt(amount) > availableBalance}
                  className="px-6 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {submitting ? '提交中...' : '提交申请'}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
