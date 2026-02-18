import { useState, useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { taskInvitationApi, taskApi } from '../services/api'
import type { ValidateTaskInvitationCodeResponse, Task } from '../types'
import { Button } from '../components/ui/button'
import { Badge } from '../components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/card'
import { useAuth } from '../contexts/AuthContext'

export default function TaskInvite() {
  const { code } = useParams<{ code: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [data, setData] = useState<ValidateTaskInvitationCodeResponse | null>(null)
  const usedInviteRef = useRef(false)

  useEffect(() => {
    if (code) {
      loadInvitation()
    }
  }, [code])

  // 已登录且邀请码校验通过后，自动「使用」邀请码：获得达人角色、绑定邀请人、绑定活动商家/服务商（仅调一次，重复调用后端返回已使用）
  useEffect(() => {
    if (!code || !data || !user || usedInviteRef.current) return
    usedInviteRef.current = true
    taskInvitationApi.useTaskInvitationCode(code).catch(() => { usedInviteRef.current = false })
  }, [code, data, user])

  const loadInvitation = async () => {
    try {
      setLoading(true)
      const response = await taskInvitationApi.validateInvitationCode(code!)
      setData(response)
    } catch (err: any) {
      setError(err.response?.data?.error || '邀请码无效或已过期')
    } finally {
      setLoading(false)
    }
  }

  const handleAcceptTask = async (task: Task) => {
    if (!user) {
      navigate('/login', { replace: true })
      return
    }

    try {
      await taskApi.acceptTask(task.id, { platform: '待选择' })
      navigate('/my-tasks')
    } catch (err: any) {
      setError(err.response?.data?.error || '接受任务失败')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-muted-foreground">加载中...</div>
      </div>
    )
  }

  if (error || !data) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center p-4">
        <Card className="max-w-md w-full">
          <CardContent className="pt-6">
            <div className="text-center">
              <div className="text-red-500 text-5xl mb-4">⚠️</div>
              <h2 className="text-xl font-bold mb-2">邀请码无效</h2>
              <p className="text-muted-foreground">{error || '邀请码不存在或已过期'}</p>
              <Button
                className="mt-4"
                onClick={() => navigate('/login')}
              >
                返回登录
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  const { invitationCode, campaign, tasks } = data

  return (
    <div className="min-h-screen bg-background">
      {/* 顶部导航 */}
      <nav className="bg-card shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold">PR Business - 任务邀请</h1>
            </div>
            <div className="flex items-center gap-4">
              {user ? (
                <>
                  <span className="text-sm text-muted-foreground">
                    {user.nickname}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => navigate('/dashboard')}
                  >
                    进入工作台
                  </Button>
                </>
              ) : (
                <Button size="sm" onClick={() => navigate('/login')}>
                  登录后接任务
                </Button>
              )}
            </div>
          </div>
        </div>
      </nav>

      {/* 主要内容 */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* 营销活动信息 */}
        <Card className="mb-6">
          <CardHeader>
            <div className="flex justify-between items-start">
              <div>
                <CardTitle className="text-2xl">{campaign.title}</CardTitle>
                <p className="text-muted-foreground mt-2">{campaign.requirements}</p>
              </div>
              <Badge variant="default">{campaign.status}</Badge>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <div className="text-sm text-muted-foreground">任务数量</div>
                <div className="text-2xl font-bold">{campaign.taskAmount}</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">任务金额</div>
                <div className="text-2xl font-bold">¥{campaign.creatorAmount || campaign.campaignAmount / campaign.taskAmount}</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">任务截止</div>
                <div className="text-sm font-medium">
                  {new Date(campaign.taskDeadline).toLocaleDateString('zh-CN')}
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">提交截止</div>
                <div className="text-sm font-medium">
                  {new Date(campaign.submissionDeadline).toLocaleDateString('zh-CN')}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* 邀请码信息 */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>邀请码信息</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <div className="text-sm text-muted-foreground">邀请码</div>
                <div className="text-lg font-mono font-bold">{invitationCode.code}</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">使用情况</div>
                <div className="text-sm font-medium">
                  {invitationCode.useCount} / {invitationCode.maxUses}
                </div>
              </div>
              {invitationCode.expiresAt && (
                <div>
                  <div className="text-sm text-muted-foreground">过期时间</div>
                  <div className="text-sm font-medium">
                    {new Date(invitationCode.expiresAt).toLocaleString('zh-CN')}
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* 可接任务列表 */}
        <div>
          <h2 className="text-xl font-bold mb-4">
            可接任务 ({data.totalTasks})
          </h2>

          {tasks.length === 0 ? (
            <Card>
              <CardContent className="pt-6">
                <div className="text-center py-8 text-muted-foreground">
                  暂无可接任务
                </div>
              </CardContent>
            </Card>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {tasks.map((task) => (
                <Card key={task.id}>
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <CardTitle className="text-lg">
                        任务 #{task.taskSlotNumber}
                      </CardTitle>
                      <Badge variant="default">开放中</Badge>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <div className="flex justify-between">
                        <span className="text-sm text-muted-foreground">金额</span>
                        <span className="font-medium">¥{campaign.creatorAmount || campaign.campaignAmount / campaign.taskAmount}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-sm text-muted-foreground">优先级</span>
                        <Badge variant={task.priority === 'HIGH' ? 'destructive' : 'secondary'}>
                          {task.priority === 'HIGH' ? '高' : task.priority === 'MEDIUM' ? '中' : '低'}
                        </Badge>
                      </div>
                      <Button
                        className="w-full mt-4"
                        onClick={() => handleAcceptTask(task)}
                        disabled={!user}
                      >
                        {user ? '立即接任务' : '登录后接任务'}
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
