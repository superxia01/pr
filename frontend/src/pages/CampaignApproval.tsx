import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { campaignApi } from '../services/api'
import type { Campaign } from '../types'
import { useToast } from '../contexts/ToastContext'

export default function CampaignApproval() {
  const navigate = useNavigate()
  const toast = useToast()
  const [campaigns, setCampaigns] = useState<Campaign[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [submittingId, setSubmittingId] = useState<string | null>(null)

  // 审核表单数据
  const [approvalForms, setApprovalForms] = useState<Record<string, {
    creatorAmount: number
    staffReferralAmount: number
    providerAmount: number
  }>>({})

  useEffect(() => {
    loadPendingCampaigns()
  }, [])

  const loadPendingCampaigns = async () => {
    setLoading(true)
    try {
      // 获取所有活动，然后过滤出待审核的
      const allCampaigns = await campaignApi.getCampaigns()
      const pendingCampaigns = allCampaigns.filter(c => c.status === 'PENDING_APPROVAL')
      setCampaigns(pendingCampaigns)

      // 初始化表单数据
      const forms: Record<string, any> = {}
      pendingCampaigns.forEach(campaign => {
        forms[campaign.id] = {
          creatorAmount: 0,
          staffReferralAmount: 0,
          providerAmount: 0,
        }
      })
      setApprovalForms(forms)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载活动列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (campaignId: string) => {
    setSubmittingId(campaignId)
    setError(null)

    try {
      const formData = approvalForms[campaignId]

      // 验证佣金分配总和
      const total = formData.creatorAmount + formData.staffReferralAmount + formData.providerAmount
      const campaign = campaigns.find(c => c.id === campaignId)

      if (total !== campaign?.taskAmount) {
        setError('佣金分配总和必须等于任务总金额')
        setSubmittingId(null)
        return
      }

      await campaignApi.approveCampaign(campaignId, formData)

      // 从列表中移除已审核的活动
      setCampaigns(prev => prev.filter(c => c.id !== campaignId))
      toast.showSuccess('审核通过，活动已发布！')
    } catch (err: any) {
      setError(err.response?.data?.error || '审核失败')
      setSubmittingId(null)
    }
  }

  const updateForm = (campaignId: string, field: string, value: number) => {
    setApprovalForms(prev => ({
      ...prev,
      [campaignId]: {
        ...prev[campaignId],
        [field]: value,
      },
    }))
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
      <div className="max-w-7xl mx-auto">
        {/* 标题 */}
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900">活动审核</h1>
          <p className="mt-1 text-sm text-gray-500">
            审核商家提交的活动申请，设置佣金分配并发布
          </p>
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        )}

        {/* 活动列表 */}
        {campaigns.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <p className="text-gray-500">暂无待审核的活动</p>
            <button
              onClick={() => navigate('/dashboard')}
              className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              返回工作台
            </button>
          </div>
        ) : (
          <div className="space-y-6">
            {campaigns.map((campaign) => {
              const formData = approvalForms[campaign.id] || {
                creatorAmount: 0,
                staffReferralAmount: 0,
                providerAmount: 0,
              }
              const total = formData.creatorAmount + formData.staffReferralAmount + formData.providerAmount
              const isValid = total === campaign.taskAmount

              return (
                <div key={campaign.id} className="bg-white rounded-lg shadow">
                  <div className="p-6">
                    {/* 活动基本信息 */}
                    <div className="mb-6">
                      <h3 className="text-lg font-medium text-gray-900">{campaign.title}</h3>
                      <p className="mt-2 text-sm text-gray-600">{campaign.requirements}</p>
                      <div className="mt-4 grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <span className="font-medium text-gray-700">商家：</span>
                          <span className="text-gray-600">{campaign.merchant?.name || '-'}</span>
                        </div>
                        <div>
                          <span className="font-medium text-gray-700">每任务积分：</span>
                          <span className="text-gray-600">{campaign.taskAmount}</span>
                        </div>
                        <div>
                          <span className="font-medium text-gray-700">任务名额：</span>
                          <span className="text-gray-600">{campaign.quota}</span>
                        </div>
                        <div>
                          <span className="font-medium text-gray-700">活动总积分：</span>
                          <span className="text-gray-600">{campaign.campaignAmount}</span>
                        </div>
                      </div>
                    </div>

                    {/* 佣金分配表单 */}
                    <div className="border-t pt-6">
                      <h4 className="text-sm font-medium text-gray-900 mb-4">设置佣金分配</h4>
                      <div className="grid grid-cols-3 gap-4 mb-4">
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">
                            达人佣金
                          </label>
                          <input
                            type="number"
                            value={formData.creatorAmount}
                            onChange={(e) => updateForm(campaign.id, 'creatorAmount', parseInt(e.target.value) || 0)}
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                            min="0"
                            required
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">
                            员工推荐佣金
                          </label>
                          <input
                            type="number"
                            value={formData.staffReferralAmount}
                            onChange={(e) => updateForm(campaign.id, 'staffReferralAmount', parseInt(e.target.value) || 0)}
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                            min="0"
                            required
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">
                            服务商佣金
                          </label>
                          <input
                            type="number"
                            value={formData.providerAmount}
                            onChange={(e) => updateForm(campaign.id, 'providerAmount', parseInt(e.target.value) || 0)}
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                            min="0"
                            required
                          />
                        </div>
                      </div>

                      {/* 佣金总和提示 */}
                      <div className={`mb-4 p-3 rounded-md text-sm ${
                        isValid ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'
                      }`}>
                        <div className="flex justify-between items-center">
                          <span>佣金总和：{total} / {campaign.taskAmount}</span>
                          <span className="font-medium">
                            {isValid ? '✓ 匹配' : '✗ 总和必须等于任务总金额'}
                          </span>
                        </div>
                      </div>

                      {/* 操作按钮 */}
                      <div className="flex justify-end space-x-3">
                        <button
                          onClick={() => loadPendingCampaigns()}
                          className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                        >
                          刷新
                        </button>
                        <button
                          onClick={() => handleApprove(campaign.id)}
                          disabled={!isValid || submittingId === campaign.id}
                          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          {submittingId === campaign.id ? '提交中...' : '审核通过并发布'}
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
