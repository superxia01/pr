import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useToast } from '../contexts/ToastContext'
import { merchantApi, campaignApi } from '../services/api'
import type { Merchant } from '../types'

export default function CreateCampaign() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const toast = useToast()
  const [merchants, setMerchants] = useState<Merchant[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  // 判断角色
  const isMerchantAdmin = user?.roles.includes('MERCHANT_ADMIN')
  const isServiceProviderAdmin = user?.roles.includes('SERVICE_PROVIDER_ADMIN')

  // 表单数据
  const [formData, setFormData] = useState({
    merchantId: '',
    title: '',
    requirements: '',
    platforms: [] as string[],
    taskAmount: 100,
    creatorAmount: 0,
    staffReferralAmount: 0,
    providerAmount: 0,
    quota: 10,
    taskDeadline: '',
    submissionDeadline: '',
  })

  useEffect(() => {
    if (isServiceProviderAdmin) {
      loadMerchants()
    } else if (isMerchantAdmin) {
      loadMyMerchant()
    }
  }, [])

  const loadMyMerchant = async () => {
    setLoading(true)
    try {
      const merchant = await merchantApi.getMyMerchant()
      setFormData(prev => ({ ...prev, merchantId: merchant.id }))
      setMerchants([merchant])
    } catch (err: any) {
      setError(err.response?.data?.error || '加载商家信息失败')
    } finally {
      setLoading(false)
    }
  }

  const loadMerchants = async () => {
    setLoading(true)
    try {
      const data = await merchantApi.getMerchants()
      setMerchants(data)
      if (data.length > 0) {
        setFormData(prev => ({ ...prev, merchantId: data[0].id }))
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '加载商家信息失败')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    setError(null)

    try {
      const payload: any = {
        title: formData.title,
        requirements: formData.requirements,
        platforms: JSON.stringify(formData.platforms),
        taskAmount: formData.taskAmount,
        quota: formData.quota,
        taskDeadline: formData.taskDeadline,
        submissionDeadline: formData.submissionDeadline,
      }

      // 服务商创建时需要指定商家
      if (isServiceProviderAdmin) {
        payload.merchantId = formData.merchantId
        // 服务商设置佣金分配
        payload.creatorAmount = formData.creatorAmount || null
        payload.staffReferralAmount = formData.staffReferralAmount || null
        payload.providerAmount = formData.providerAmount || null
      }

      await campaignApi.createCampaign(payload)

      toast.showSuccess('创建成功！')
      navigate('/dashboard')
    } catch (err: any) {
      setError(err.response?.data?.error || '创建失败')
      setSubmitting(false)
    }
  }

  const platformOptions = [
    '小红书',
    '抖音',
    '快手',
    'B站',
    '微博',
    '微信视频号',
  ]

  const togglePlatform = (platform: string) => {
    const platforms = formData.platforms || []
    if (platforms.includes(platform)) {
      setFormData({
        ...formData,
        platforms: platforms.filter(p => p !== platform),
      })
    } else {
      setFormData({
        ...formData,
        platforms: [...platforms],
      })
    }
  }

  // 计算活动总积分
  const totalAmount = formData.taskAmount * formData.quota

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-3xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">
              {isServiceProviderAdmin ? '创建营销活动' : '提交活动申请'}
            </h1>
            <p className="mt-1 text-sm text-gray-500">
              {isServiceProviderAdmin
                ? '创建新的营销活动并设置任务名额和佣金分配，直接发布'
                : '提交活动申请，由服务商审核并设置佣金分配后发布'}
            </p>
          </div>

          {/* 内容 */}
          <div className="p-6">
            {loading && (
              <div className="text-center py-12">
                <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                <p className="mt-2 text-gray-500">加载中...</p>
              </div>
            )}

            {!loading && (
              <form onSubmit={handleSubmit} className="space-y-6">
                {/* 错误提示 */}
                {error && (
                  <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                    {error}
                  </div>
                )}

                {/* 服务商选择商家 */}
                {isServiceProviderAdmin && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      选择商家 <span className="text-red-500">*</span>
                    </label>
                    <select
                      value={formData.merchantId}
                      onChange={(e) => setFormData({ ...formData, merchantId: e.target.value })}
                      className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                      required
                    >
                      <option value="">请选择商家</option>
                      {merchants.map((merchant) => (
                        <option key={merchant.id} value={merchant.id}>
                          {merchant.name}
                        </option>
                      ))}
                    </select>
                  </div>
                )}

                {/* 活动标题 */}
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    活动标题 <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    value={formData.title}
                    onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                    placeholder="请输入活动标题"
                    required
                  />
                </div>

                {/* 活动要求 */}
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    活动要求 <span className="text-red-500">*</span>
                  </label>
                  <textarea
                    value={formData.requirements}
                    onChange={(e) => setFormData({ ...formData, requirements: e.target.value })}
                    rows={3}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                    placeholder="请详细描述活动要求、内容规范等"
                    required
                  />
                </div>

                {/* 平台选择 */}
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    平台选择 <span className="text-red-500">*</span>
                  </label>
                  <div className="mt-2 grid grid-cols-3 gap-2">
                    {platformOptions.map((platform) => {
                      const isSelected = (formData.platforms || []).includes(platform)
                      return (
                        <button
                          key={platform}
                          type="button"
                          onClick={() => togglePlatform(platform)}
                          className={`px-4 py-2 border rounded-md text-sm transition-colors ${
                            isSelected
                              ? 'bg-blue-500 text-white border-blue-500'
                              : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                          }`}
                        >
                          {platform}
                        </button>
                      )
                    })}
                  </div>
                </div>

                {/* 积分设置 */}
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      每任务积分 <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="number"
                      value={formData.taskAmount}
                      onChange={(e) => setFormData({ ...formData, taskAmount: parseInt(e.target.value) || 0 })}
                      className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                      min="0"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      活动总积分
                    </label>
                    <input
                      type="number"
                      value={totalAmount}
                      className="mt-1 block w-full rounded-md border-gray-300 bg-gray-50 sm:text-sm border p-2"
                      readOnly
                    />
                    <p className="mt-1 text-xs text-gray-500">自动计算：每任务积分 × 任务名额</p>
                  </div>
                </div>

                {/* 服务商设置佣金分配 */}
                {isServiceProviderAdmin && (
                  <div className="border rounded-lg p-4 bg-gray-50">
                    <h3 className="text-sm font-medium text-gray-700 mb-3">佣金分配（可选）</h3>
                    <div className="grid grid-cols-3 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700">
                          达人佣金
                        </label>
                        <input
                          type="number"
                          value={formData.creatorAmount || ''}
                          onChange={(e) => setFormData({ ...formData, creatorAmount: parseInt(e.target.value) || 0 })}
                          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                          min="0"
                          placeholder="0"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700">
                          员工推荐佣金
                        </label>
                        <input
                          type="number"
                          value={formData.staffReferralAmount || ''}
                          onChange={(e) => setFormData({ ...formData, staffReferralAmount: parseInt(e.target.value) || 0 })}
                          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                          min="0"
                          placeholder="0"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700">
                          服务商佣金
                        </label>
                        <input
                          type="number"
                          value={formData.providerAmount || ''}
                          onChange={(e) => setFormData({ ...formData, providerAmount: parseInt(e.target.value) || 0 })}
                          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                          min="0"
                          placeholder="0"
                        />
                      </div>
                    </div>
                    <p className="mt-2 text-xs text-gray-500">留空则稍后设置</p>
                  </div>
                )}

                {/* 任务名额 */}
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    任务名额数量 <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    value={formData.quota}
                    onChange={(e) => setFormData({ ...formData, quota: parseInt(e.target.value) || 0 })}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                    min="1"
                    required
                  />
                  <p className="mt-1 text-xs text-gray-500">将创建 {formData.quota} 个任务名额供达人接取</p>
                </div>

                {/* 时间设置 */}
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      任务截止时间 <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="datetime-local"
                      value={formData.taskDeadline}
                      onChange={(e) => setFormData({ ...formData, taskDeadline: e.target.value })}
                      className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      提交截止时间 <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="datetime-local"
                      value={formData.submissionDeadline}
                      onChange={(e) => setFormData({ ...formData, submissionDeadline: e.target.value })}
                      className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                      required
                    />
                  </div>
                </div>

                {/* 商家提示 */}
                {isMerchantAdmin && (
                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <h4 className="text-sm font-medium text-blue-800 mb-2">提交流程说明</h4>
                    <ul className="text-xs text-blue-700 space-y-1 list-disc list-inside">
                      <li>提交后活动状态为"待审核"</li>
                      <li>服务商将审核并设置佣金分配方案</li>
                      <li>审核通过后活动将自动发布为OPEN状态</li>
                      <li>请确保每任务积分和任务名额填写正确</li>
                    </ul>
                  </div>
                )}

                {/* 按钮 */}
                <div className="flex justify-end space-x-3">
                  <button
                    type="button"
                    onClick={() => navigate('/dashboard')}
                    className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    disabled={submitting}
                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {submitting ? '提交中...' : isServiceProviderAdmin ? '创建活动' : '提交申请'}
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
