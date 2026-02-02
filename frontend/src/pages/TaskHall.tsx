import { useEffect, useState } from 'react'
import { taskApi } from '../services/api'
import type { Task } from '../types'
import { DataTable } from '../components/ui/data-table'
import { Button } from '../components/ui/button'
import { Badge } from '../components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/card'

export default function TaskHall() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  // const [total, setTotal] = useState(0)
  // const [page, setPage] = useState(1)
  const pageSize = 20

  useEffect(() => {
    loadTasks()
  }, [])

  const loadTasks = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await taskApi.getTaskHall({ page: 1, page_size: pageSize })
      setTasks(response.data)
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

  // 定义表格列
  const columns = [
    {
      id: 'campaign',
      header: '活动信息',
      sortable: false,
      cell: ({ row }: { row: Task }) => (
        <div className="flex flex-col">
          <h4 className="font-medium text-gray-900">
            {row.campaign?.title || '未知活动'}
          </h4>
          <p className="text-sm text-gray-500 line-clamp-2 mt-1">
            {row.campaign?.requirements || '无要求'}
          </p>
        </div>
      )
    },
    {
      id: 'reward',
      header: '积分奖励',
      sortable: false,
      cell: ({ row }: { row: Task }) => (
        <div className="text-center">
          <span className="text-lg font-bold text-blue-600">
            {row.campaign?.taskAmount || 0}
          </span>
          <span className="text-sm text-gray-500 ml-1">积分</span>
        </div>
      )
    },
    {
      id: 'status',
      header: '状态',
      sortable: false,
      cell: ({ row }: { row: Task }) => (
        <div className="flex justify-center">
          <Badge
            variant={
              row.status === 'OPEN' ? 'default' :
              row.status === 'ASSIGNED' ? 'secondary' :
              row.status === 'SUBMITTED' ? 'outline' :
              row.status === 'APPROVED' ? 'secondary' :
              'destructive'
            }
          >
            {row.status === 'OPEN' && '可接取'}
            {row.status === 'ASSIGNED' && '已分配'}
            {row.status === 'SUBMITTED' && '已提交'}
            {row.status === 'APPROVED' && '已通过'}
            {row.status === 'REJECTED' && '已拒绝'}
          </Badge>
        </div>
      )
    },
    {
      id: 'priority',
      header: '优先级',
      sortable: false,
      cell: ({ row }: { row: Task }) => (
        <div className="flex justify-center">
          <Badge
            variant={
              row.priority === 'HIGH' ? 'destructive' :
              row.priority === 'MEDIUM' ? 'secondary' :
              'outline'
            }
          >
            {row.priority === 'HIGH' && '高'}
            {row.priority === 'MEDIUM' && '中'}
            {row.priority === 'LOW' && '低'}
          </Badge>
        </div>
      )
    },
    {
      id: 'platform',
      header: '平台',
      sortable: false,
      cell: ({ row }: { row: Task }) => (
        <div className="text-center">
          <span className="text-sm font-medium text-gray-700">
            {row.platform || '未指定'}
          </span>
        </div>
      )
    },
    {
      id: 'slot',
      header: '任务名额',
      sortable: false,
      cell: ({ row }: { row: Task }) => (
        <div className="text-center">
          <span className="text-sm font-medium text-gray-700">
            #{row.taskSlotNumber}
          </span>
        </div>
      )
    },
    {
      id: 'deadline',
      header: '截止时间',
      sortable: false,
      cell: ({ row }: { row: Task }) => (
        <div className="text-center">
          <span className="text-sm text-gray-700">
            {row.campaign?.taskDeadline ?
              new Date(row.campaign.taskDeadline).toLocaleDateString() :
              '未设置'
            }
          </span>
        </div>
      )
    },
    {
      id: 'actions',
      header: '操作',
      sortable: false,
      cell: ({ row }: { row: Task }) => (
        <div className="flex justify-center">
          <Button
            onClick={() => handleAcceptTask(row.id)}
            disabled={row.status !== 'OPEN'}
            className="w-full"
          >
            接任务
          </Button>
        </div>
      )
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <Card>
          <CardHeader>
            <CardTitle className="text-2xl font-bold text-gray-900">任务大厅</CardTitle>
            <CardDescription className="text-sm text-gray-500">
              查看并接取您感兴趣的任务
            </CardDescription>
          </CardHeader>
          <CardContent>
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
              <div className="space-y-6">
                <DataTable
                  data={tasks}
                  columns={columns}
                  searchable={true}
                  pageSize={pageSize}
                  onRowClick={(row) => {
                    if (row.status === 'OPEN') {
                      handleAcceptTask(row.id)
                    }
                  }}
                  emptyMessage="暂无任务数据"
                />
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
