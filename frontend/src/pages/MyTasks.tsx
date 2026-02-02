import { useEffect, useState } from 'react'
import { taskApi } from '../services/api'
import type { Task } from '../types'

export default function MyTasks() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [filter, setFilter] = useState<string>('')
  const [submitModalOpen, setSubmitModalOpen] = useState(false)
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [platformUrl, setPlatformUrl] = useState('')
  const [notes, setNotes] = useState('')

  useEffect(() => {
    loadTasks()
  }, [filter])

  const loadTasks = async () => {
    setLoading(true)
    setError(null)
    try {
      const params = filter ? { status: filter } : undefined
      const response = await taskApi.getMyTasks(params)
      setTasks(response)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载任务失败')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmitTask = async () => {
    if (!selectedTask) return

    try {
      await taskApi.submitTask(selectedTask.id, {
        platformUrl,
        notes,
      })
      alert('提交成功！')
      setSubmitModalOpen(false)
      setSelectedTask(null)
      setPlatformUrl('')
      setNotes('')
      loadTasks()
    } catch (err: any) {
      alert(err.response?.data?.error || '提交失败')
    }
  }

  const openSubmitModal = (task: Task) => {
    setSelectedTask(task)
    setPlatformUrl(task.platformUrl || '')
    setNotes(task.notes || '')
    setSubmitModalOpen(true)
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'OPEN': return 'bg-green-100 text-green-800'
      case 'ASSIGNED': return 'bg-blue-100 text-blue-800'
      case 'SUBMITTED': return 'bg-yellow-100 text-yellow-800'
      case 'APPROVED': return 'bg-purple-100 text-purple-800'
      case 'REJECTED': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const getStatusText = (status: string) => {
    switch (status) {
      case 'OPEN': return '开放中'
      case 'ASSIGNED': return '已接取'
      case 'SUBMITTED': return '已提交'
      case 'APPROVED': return '已通过'
      case 'REJECTED': return '已拒绝'
      default: return status
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题和筛选 */}
          <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">我的任务</h1>
              <p className="mt-1 text-sm text-gray-500">管理您接取的任务</p>
            </div>
            <select
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              className="px-4 py-2 border border-gray-300 rounded-md text-sm"
            >
              <option value="">全部状态</option>
              <option value="ASSIGNED">已接取</option>
              <option value="SUBMITTED">已提交</option>
              <option value="APPROVED">已通过</option>
              <option value="REJECTED">已拒绝</option>
            </select>
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

            {!loading && tasks.length === 0 && (
              <div className="text-center py-12 bg-gray-50 rounded-lg">
                <p className="text-gray-500">暂无任务</p>
              </div>
            )}

            {!loading && tasks.length > 0 && (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        营销活动
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        积分奖励
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        状态
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        截止时间
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        操作
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {tasks.map((task) => (
                      <tr key={task.id}>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-medium text-gray-900">
                            {task.campaign?.title || '未知活动'}
                          </div>
                          <div className="text-sm text-gray-500">
                            名额 #{task.taskSlotNumber}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className="text-sm font-bold text-blue-600">
                            {task.campaign?.taskAmount || 0}
                          </span>
                          <span className="text-xs text-gray-500 ml-1">积分</span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(task.status)}`}>
                            {getStatusText(task.status)}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {task.campaign?.submissionDeadline ? new Date(task.campaign.submissionDeadline).toLocaleString() : '-'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                          {task.status === 'ASSIGNED' && (
                            <button
                              onClick={() => openSubmitModal(task)}
                              className="text-blue-600 hover:text-blue-900"
                            >
                              提交任务
                            </button>
                          )}
                          {task.status === 'SUBMITTED' && (
                            <span className="text-gray-400">等待审核</span>
                          )}
                          {task.status === 'APPROVED' && (
                            <span className="text-green-600">已完成</span>
                          )}
                          {task.status === 'REJECTED' && (
                            <button
                              onClick={() => openSubmitModal(task)}
                              className="text-red-600 hover:text-red-900"
                            >
                              重新提交
                            </button>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>

        {/* 提交任务对话框 */}
        {submitModalOpen && selectedTask && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
              <h3 className="text-lg font-medium text-gray-900 mb-4">提交任务</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">营销活动</label>
                  <p className="mt-1 text-sm text-gray-900">{selectedTask.campaign?.title}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">平台链接</label>
                  <input
                    type="url"
                    value={platformUrl}
                    onChange={(e) => setPlatformUrl(e.target.value)}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                    placeholder="https://..."
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">备注说明</label>
                  <textarea
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    rows={3}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm border p-2"
                    placeholder="请填写任务说明..."
                  />
                </div>
                <div className="flex justify-end space-x-3">
                  <button
                    onClick={() => setSubmitModalOpen(false)}
                    className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                  >
                    取消
                  </button>
                  <button
                    onClick={handleSubmitTask}
                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                  >
                    提交
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
