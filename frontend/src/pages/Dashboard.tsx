import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { getRoleName } from '../lib/roles'
import { useState } from 'react'
import {
  Gift,
  Users,
  Building2,
  Briefcase,
  Sparkles,
  Coins,
  Receipt,
  ShieldCheck,
  ArrowRight,
  User,
  Badge,
} from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/card'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Label } from '../components/ui/label'
import { invitationApi, authApi } from '../services/api'

export default function Dashboard() {
  const { user, updateUser } = useAuth()
  const [invitationCode, setInvitationCode] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  // 多角色：按「拥有的角色」显示所有对应快捷操作（与侧边栏一致）
  const roles = user?.roles ?? []
  const isSuperAdmin = () => roles.some((r) => r.toUpperCase() === 'SUPER_ADMIN')
  const isMerchantAdmin = () => roles.some((r) => r.toUpperCase() === 'MERCHANT_ADMIN')
  const isServiceProviderAdmin = () => roles.some((r) => r.toUpperCase() === 'SERVICE_PROVIDER_ADMIN')
  const isCreator = () => roles.some((r) => r.toUpperCase() === 'CREATOR')
  const isMerchantStaff = () => roles.some((r) => r.toUpperCase() === 'MERCHANT_STAFF')
  const isBasicUser = () => roles.some((r) => r.toUpperCase() === 'BASIC_USER')
  const canManageInvitations = () => isSuperAdmin() || isServiceProviderAdmin() || isMerchantAdmin() || isBasicUser()

  const quickActions = [
    // 达人相关
    ...(isCreator() ? [{
      title: '任务大厅',
      description: '浏览和接受任务',
      href: '/task-hall',
      icon: Briefcase,
      color: 'bg-blue-500',
    }] : []),
    // 商家/服务商相关
    ...(isMerchantAdmin() || isMerchantStaff() ? [{
      title: '我的任务',
      description: '管理已接任务',
      href: '/my-tasks',
      icon: Sparkles,
      color: 'bg-purple-500',
    }] : []),
    ...(isMerchantAdmin() || isMerchantStaff() ? [{
      title: '邀请码管理',
      description: '管理邀请码',
      href: '/invitations',
      icon: Gift,
      color: 'bg-orange-500',
    }] : []),
    // 财务相关
    ...(isMerchantAdmin() || isMerchantStaff() ? [{
      title: '充值',
      description: '账户充值',
      href: '/recharge',
      icon: Coins,
      color: 'bg-green-500',
    }, {
      title: '提现记录',
      description: '查看提现记录',
      href: '/withdrawals',
      icon: Receipt,
      color: 'bg-emerald-500',
    }] : []),
    ...(isSuperAdmin() ? [{
      title: '提现审核',
      description: '审核提现申请',
      href: '/withdrawal-review',
      icon: ShieldCheck,
      color: 'bg-amber-500',
    }] : []),
    // 管理功能 - 商家管理
    ...(isMerchantAdmin() || isMerchantStaff() ? [{
      title: '商家管理',
      description: '管理商家资料',
      href: '/merchant',
      icon: Building2,
      color: 'bg-cyan-500',
    }] : []),
    // 管理功能 - 服务商管理
    ...(isSuperAdmin() ? [{
      title: '服务商管理',
      description: '管理服务商组织',
      href: '/service-providers',
      icon: Building2,
      color: 'bg-teal-500',
    }] : []),
    // 管理功能 - 达人管理
    ...(isSuperAdmin() ? [{
      title: '达人管理',
      description: '管理本组织达人',
      href: '/user-management',
      icon: Users,
      color: 'bg-orange-500',
    }] : []),
  ]

  // 处理邀请码提交
  const handleInvitationCodeSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    setSuccess(false)

    const code = invitationCode.trim()

    if (!code) {
      setError('请输入邀请码')
      setLoading(false)
      return
    }

    try {
      // 调用后端API使用邀请码
      const response = await invitationApi.useInvitationCode({
        code: invitationCode,
        userId: user!.id,
      })

      // 检查响应是否成功（后端返回 200 且有 message 字段）
      if (response.error) {
        setError(response.error || '使用邀请码失败')
        setLoading(false)
        return
      }

      if (response.message) {
        // 成功！显示成功提示
        setSuccess(true)
        setLoading(false)

        // 重新获取用户信息
        try {
          const updatedUser = await authApi.getCurrentUser()
          // 更新 AuthContext 中的用户信息
          updateUser(updatedUser)
        } catch (err) {
          console.error('获取用户信息失败:', err)
        }

        // 3秒后跳转到 dashboard（此时用户信息已更新）
        setTimeout(() => {
          window.location.href = '/dashboard'
        }, 3000)
      } else {
        setError('使用邀请码失败，请重试')
        setLoading(false)
      }
    } catch (err: any) {
      console.error('使用邀请码失败:', err)
      setError('使用邀请码失败，请重试')
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* 顶部导航 */}
      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">PR Business</h1>
              <p className="text-gray-600 text-sm">营销任务管理平台</p>
            </div>
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-600">欢迎回来，{user?.nickname || '用户'}</span>
              <Link to="/dashboard" className="text-sm text-blue-600 hover:text-blue-800">
                刷新
              </Link>
            </div>
          </div>
        </div>
      </nav>

      {/* 主要内容 */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {!user ? (
          <div className="flex items-center justify-center min-h-[400px]">
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
              <p className="text-center">未登录</p>
            </div>
          </div>
        ) : (
          <div className="space-y-6">
            {/* 用户信息卡片 */}
            <Card className="hover:shadow-lg transition-all">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <User className="h-5 w-5" />
                  个人信息
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-start gap-6">
                  {/* 头像 */}
                  <div className="flex-shrink-0">
                    <img
                      src={user.avatarUrl || '/default-avatar.png'}
                      alt={user.nickname || '用户头像'}
                      className="w-20 h-20 rounded-full object-cover border-2 border-gray-200"
                    />
                  </div>

                  {/* 用户详情 */}
                  <div className="flex-1 space-y-3">
                    <div>
                      <div className="text-sm text-gray-500">昵称</div>
                      <div className="text-lg font-semibold">{user.nickname || '未设置'}</div>
                    </div>

                    <div>
                      <div className="text-sm text-gray-500">用户ID</div>
                      <div className="text-sm font-mono text-gray-700">{user.id}</div>
                    </div>

                    <div>
                      <div className="text-sm text-gray-500 mb-2">
                        拥有角色 ({user.roles?.length || 0} 个)
                      </div>
                      <div className="flex flex-wrap gap-2">
                        {user.roles && user.roles.length > 0 ? (
                          user.roles.map((role) => (
                            <span
                              key={role}
                              className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800"
                            >
                              <Badge className="h-3 w-3" />
                              {getRoleName(role)}
                            </span>
                          ))
                        ) : (
                          <span className="text-sm text-gray-400">暂无角色</span>
                        )}
                      </div>
                    </div>

                    <div>
                      <div className="text-sm text-gray-500">账号状态</div>
                      <div className={`text-sm font-medium ${
                        user.status === 'active' ? 'text-green-600' : 'text-gray-400'
                      }`}>
                        {user.status === 'active' ? '✓ 正常' : user.status || '未知'}
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* 邀请码管理卡片 */}
            {canManageInvitations() && (
              <Card className="hover:shadow-lg transition-all">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Gift className="h-5 w-5" />
                    邀请码管理
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <p className="text-sm text-gray-600">
                      根据您的角色和管理组织，生成邀请码来邀请用户加入您的组织。
                    </p>

                    <div className="space-y-2">
                      <Label htmlFor="invitation-code">填写邀请码</Label>
                      <Input
                        id="invitation-code"
                        value={invitationCode}
                        onChange={(e) => setInvitationCode(e.target.value)}
                        placeholder="请输入邀请码"
                        disabled={loading}
                      />
                    </div>

                    {error && (
                      <div className="p-3 bg-red-50 border border-red-200 text-red-700 rounded text-sm">
                        {error}
                      </div>
                    )}

                    {success && (
                      <div className="p-3 bg-green-50 border border-green-200 text-green-700 rounded text-sm">
                        ✓ 邀请码使用成功！
                      </div>
                    )}

                    <form onSubmit={handleInvitationCodeSubmit}>
                      <Button
                        type="submit"
                        size="lg"
                        disabled={loading || !invitationCode}
                        className="w-full"
                      >
                        {loading ? '提交中...' : '使用邀请码'}
                      </Button>
                    </form>
                  </div>
                </CardContent>
              </Card>
            )}

            {!canManageInvitations() && (
              <Card>
                <CardContent className="py-12">
                  <p className="text-center text-gray-500">您当前的角色无法邀请他人。</p>
                </CardContent>
              </Card>
            )}

            {/* 快捷操作 */}
            {quickActions.length > 0 && (
              <div>
                <h2 className="text-xl font-semibold mb-4">快捷操作</h2>
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                  {quickActions.map((action) => (
                    <Link
                      key={action.title}
                      to={action.href}
                      className="group flex items-center gap-4 p-4 rounded-lg border bg-white hover:bg-gray-50 transition-colors"
                    >
                      <div className={`p-3 rounded-lg ${action.color} text-white`}>
                        <action.icon className="h-5 w-5" />
                      </div>
                      <div className="flex-1">
                        <div className="font-medium">{action.title}</div>
                        <div className="text-sm text-gray-500">{action.description}</div>
                      </div>
                      <ArrowRight className="h-5 w-5 text-gray-400 group-hover:text-gray-600" />
                    </Link>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </main>
    </div>
  )
}
