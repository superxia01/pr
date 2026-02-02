import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { merchantApi, campaignApi } from '../services/api'
import type { Merchant } from '../types'

export default function CreateCampaign() {
  const navigate = useNavigate()
  const [merchants, setMerchants] = useState<Merchant[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  // 表单数据
  const [formData, setFormData] = useState({
    title: '',
    requirements: '',
    platforms: [] as string[],
    taskAmount: 100,
    campaignAmount: 1000,
    quota: 10,
    taskDeadline: '',
    submissionDeadline: '',
  })

  useEffect(() => {
    loadMerchants()
  }, [])

  const loadMerchants = async () => {
    setLoading(true)
    try {
      // 获取当前用户管理的商家
      const merchant = await merchantApi.getMyMerchant()
      setMerchants([merchant])
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
      await campaignApi.createCampaign({
        title: formData.title,
        requirements: formData.requirements,
        platforms: JSON.stringify(formData.platforms),
        taskAmount: formData.taskAmount,
        campaignAmount: formData.campaignAmount,
        quota: formData.quota,
        taskDeadline: formData.taskDeadline,
        submissionDeadline: formData.submissionDeadline,
      })

      alert('创建成功！')
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
        platforms: [...platforms, platform],
      })
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-3xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">创建营销活动</h1>
            <p className="mt-1 text-sm text-gray-500">创建新的营销活动并设置任务名额</p>
          </div>

          {/* 内容 */}
          <div className="p-6">
            {loading && (
              <div className="text-center py-12">
                <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                <p className="mt-2 text-gray-500">加载中...</p>
              </div>
            )}

            {!loading && merchants.length === 0 && (
              <div className="text-center py-12 bg-red-50 rounded-lg">
                <p className="text-red-700">您不是商家管理员，无法创建营销活动</p>
              </div>
            )}

            {!loading && merchants.length > 0 && (
              <form onSubmit={handleSubmit} className="space-y-6">
                {/* 错误提示 */}
                {error && (
                  <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                    {error}
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
                      每任务积分（达人） <span className="text-red-500">*</span>
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
                      活动总积分 <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="number"
                      value={formData.campaignAmount}
                      onChange={(e) => setFormData({ ...formData, campaignAmount: parseInt(e.target.value) || 0 })}
                      className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                      min="0"
                      required
                    />
                  </div>
                </div>

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
                    {submitting ? '创建中...' : '创建活动'}
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
