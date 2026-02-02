import { useEffect, useState } from 'react'
import { creatorApi } from '../services/api'
import type { Creator } from '../types'

export default function CreatorProfile() {
  const [creator, setCreator] = useState<Creator | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editMode, setEditMode] = useState(false)
  const [inviterInfo, setInviterInfo] = useState<any>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    setError(null)
    try {
      const creatorData = await creatorApi.getMyCreatorProfile()
      setCreator(creatorData)
      // 获取邀请关系
      const inviterData = await creatorApi.getCreatorInviterRelationship(creatorData.id)
      setInviterInfo(inviterData)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载数据失败')
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateProfile = async (data: any) => {
    if (!creator) return
    try {
      await creatorApi.updateMyCreatorProfile(data)
      loadData()
      setEditMode(false)
    } catch (err: any) {
      alert(err.response?.data?.error || '更新失败')
    }
  }

  const handleBreakRelationship = async () => {
    if (!creator || !confirm('确定要解除与邀请人的关系吗？')) return
    try {
      await creatorApi.breakInviterRelationship(creator.id)
      loadData()
    } catch (err: any) {
      alert(err.response?.data?.error || '操作失败')
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

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        </div>
      </div>
    )
  }

  if (!creator) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto text-center py-12">
          <p className="text-gray-500">您不是达人</p>
        </div>
      </div>
    )
  }

  const levelColors = {
    UGC: 'bg-gray-100 text-gray-800',
    KOC: 'bg-green-100 text-green-800',
    INF: 'bg-blue-100 text-blue-800',
    KOL: 'bg-purple-100 text-purple-800',
  }

  const statusColors = {
    active: 'bg-green-100 text-green-800',
    banned: 'bg-red-100 text-red-800',
    inactive: 'bg-gray-100 text-gray-800',
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">达人个人中心</h1>
              <p className="mt-1 text-sm text-gray-500">管理您的达人资料和信息</p>
            </div>
            {!editMode && (
              <button
                onClick={() => setEditMode(true)}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                编辑资料
              </button>
            )}
          </div>

          <div className="p-6">
            {editMode ? (
              <EditProfileForm
                creator={creator}
                onSave={handleUpdateProfile}
                onCancel={() => setEditMode(false)}
              />
            ) : (
              <div className="space-y-8">
                {/* 基本信息 */}
                <div>
                  <h2 className="text-lg font-medium text-gray-900 mb-4">基本信息</h2>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div>
                      <label className="block text-sm font-medium text-gray-700">昵称</label>
                      <p className="mt-1 text-sm text-gray-900">{creator.user?.nickname || '-'}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">达人等级</label>
                      <span className={`mt-1 px-3 inline-flex text-sm leading-5 font-semibold rounded-full ${levelColors[creator.level]}`}>
                        {creator.level}
                      </span>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">粉丝数</label>
                      <p className="mt-1 text-sm text-gray-900">{creator.followersCount.toLocaleString()}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">微信昵称</label>
                      <p className="mt-1 text-sm text-gray-900">{creator.wechatNickname || '-'}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">状态</label>
                      <span className={`mt-1 px-3 inline-flex text-sm leading-5 font-semibold rounded-full ${statusColors[creator.status]}`}>
                        {creator.status === 'active' ? '正常' :
                         creator.status === 'banned' ? '封禁' :
                         creator.status === 'inactive' ? '未激活' : creator.status}
                      </span>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">注册时间</label>
                      <p className="mt-1 text-sm text-gray-500">{new Date(creator.createdAt).toLocaleString()}</p>
                    </div>
                  </div>
                </div>

                {/* 邀请关系 */}
                <div>
                  <h2 className="text-lg font-medium text-gray-900 mb-4">邀请关系</h2>
                  {inviterInfo?.inviter ? (
                    <div className="bg-gray-50 rounded-lg p-4">
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div>
                          <label className="block text-sm font-medium text-gray-700">邀请人</label>
                          <p className="mt-1 text-sm text-gray-900">{inviterInfo.inviter.nickname || '-'}</p>
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700">邀请人类型</label>
                          <p className="mt-1 text-sm text-gray-900">
                            {inviterInfo.inviterType === 'PROVIDER_STAFF' ? '服务商员工' :
                             inviterInfo.inviterType === 'PROVIDER_ADMIN' ? '服务商管理员' :
                             inviterInfo.inviterType === 'OTHER' ? '其他' : '-'}
                          </p>
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700">关系状态</label>
                          <div className="mt-1 flex items-center gap-2">
                            <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                              inviterInfo.relationship_broken ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'
                            }`}>
                              {inviterInfo.relationship_broken ? '已解除' : '正常'}
                            </span>
                            {!inviterInfo.relationship_broken && (
                              <button
                                onClick={handleBreakRelationship}
                                className="text-sm text-red-600 hover:text-red-900"
                              >
                                解除关系
                              </button>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="text-center py-8 bg-gray-50 rounded-lg">
                      <p className="text-gray-500">无邀请关系</p>
                    </div>
                  )}
                </div>

                {/* 等级说明 */}
                <div>
                  <h2 className="text-lg font-medium text-gray-900 mb-4">达人等级说明</h2>
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <div className="bg-gray-50 rounded-lg p-4 text-center">
                      <span className={`px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full ${levelColors.UGC}`}>
                        UGC
                      </span>
                      <p className="mt-2 text-sm text-gray-600">普通用户</p>
                    </div>
                    <div className="bg-gray-50 rounded-lg p-4 text-center">
                      <span className={`px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full ${levelColors.KOC}`}>
                        KOC
                      </span>
                      <p className="mt-2 text-sm text-gray-600">关键意见消费者</p>
                    </div>
                    <div className="bg-gray-50 rounded-lg p-4 text-center">
                      <span className={`px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full ${levelColors.INF}`}>
                        INF
                      </span>
                      <p className="mt-2 text-sm text-gray-600">达人</p>
                    </div>
                    <div className="bg-gray-50 rounded-lg p-4 text-center">
                      <span className={`px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full ${levelColors.KOL}`}>
                        KOL
                      </span>
                      <p className="mt-2 text-sm text-gray-600">关键意见领袖</p>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

// 编辑资料表单
function EditProfileForm({
  creator,
  onSave,
  onCancel
}: {
  creator: Creator
  onSave: (data: any) => void
  onCancel: () => void
}) {
  const [formData, setFormData] = useState({
    wechatOpenId: creator.wechatOpenId || '',
    wechatNickname: creator.wechatNickname || '',
    wechatAvatar: creator.wechatAvatar || '',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(formData)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-w-2xl">
      <div>
        <label className="block text-sm font-medium text-gray-700">微信 OpenID</label>
        <input
          type="text"
          value={formData.wechatOpenId}
          onChange={(e) => setFormData({ ...formData, wechatOpenId: e.target.value })}
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700">微信昵称</label>
        <input
          type="text"
          value={formData.wechatNickname}
          onChange={(e) => setFormData({ ...formData, wechatNickname: e.target.value })}
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700">微信头像 URL</label>
        <input
          type="url"
          value={formData.wechatAvatar}
          onChange={(e) => setFormData({ ...formData, wechatAvatar: e.target.value })}
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
        />
      </div>
      <div className="flex justify-end space-x-3">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
        >
          取消
        </button>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          保存
        </button>
      </div>
    </form>
  )
}
