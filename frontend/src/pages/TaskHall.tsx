import { useEffect, useState } from 'react'
import { taskApi } from '../services/api'
import type { Task } from '../types'

export default function TaskHall() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const pageSize = 20

  useEffect(() => {
    loadTasks()
  }, [page])

  const loadTasks = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await taskApi.getTaskHall({ page, page_size: pageSize })
      setTasks(response.data)
      setTotal(response.total || 0)
    } catch (err: any) {
      setError(err.response?.data?.error || '加载任务失败')
    } finally {
      setLoading(false)
    }
  }

  const handleAcceptTask = async (taskId: string) => {
    const platform = prompt('请输入平台名称（如：小红书、抖音、快手等）：')
    if (!platform) return

    try {
      await taskApi.acceptTask(taskId, { platform })
      alert('接任务成功！')
      loadTasks()
    } catch (err: any) {
      alert(err.response?.data?.error || '接任务失败')
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-lg shadow">
          {/* 标题 */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">任务大厅</h1>
            <p className="mt-1 text-sm text-gray-500">查看并接取您感兴趣的任务</p>
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
                <p className="text-gray-500">暂无可接任务</p>
              </div>
            )}

            {!loading && tasks.length > 0 && (
              <>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {tasks.map((task) => (
                    <div key={task.id} className="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow">
                      {/* 营销活动标题 */}
                      <h3 className="text-lg font-semibold text-gray-900 mb-2">
                        {task.campaign?.title || '未知活动'}
                      </h3>

                      {/* 任务要求 */}
                      <p className="text-sm text-gray-600 mb-4 line-clamp-2">
                        {task.campaign?.requirements || '无要求'}
                      </p>

                      {/* 积分奖励 */}
                      <div className="mb-4">
                        <span className="text-2xl font-bold text-blue-600">
                          {task.campaign?.taskAmount || 0}
                        </span>
                        <span className="text-sm text-gray-500 ml-1">积分</span>
                      </div>

                      {/* 任务详情 */}
                      <div className="space-y-2 text-sm text-gray-600 mb-4">
                        <div className="flex justify-between">
                          <span>任务名额：</span>
                          <span className="font-medium">#{task.taskSlotNumber}</span>
                        </div>
                        {task.platform && (
                          <div className="flex justify-between">
                            <span>平台：</span>
                            <span className="font-medium">{task.platform}</span>
                          </div>
                        )}
                        {task.campaign?.taskDeadline && (
                          <div className="flex justify-between">
                            <span>截止时间：</span>
                            <span className="font-medium">
                              {new Date(task.campaign.taskDeadline).toLocaleDateString()}
                            </span>
                          </div>
                        )}
                      </div>

                      {/* 操作按钮 */}
                      <button
                        onClick={() => handleAcceptTask(task.id)}
                        className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                      >
                        接任务
                      </button>
                    </div>
                  ))}
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
