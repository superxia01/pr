import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { Button } from '../components/ui/button'
import { Card, CardContent, CardHeader } from '../components/ui/card'
import { Badge } from '../components/ui/badge'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'https://pr.crazyaigc.com'

export default function Login() {
  const [loginType, setLoginType] = useState<'wechat' | 'password'>('wechat')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login } = useAuth()
  const [searchParams] = useSearchParams()

  // 处理 auth-center 回调（V3.1 统一 Token 模式）
  useEffect(() => {
    const token = searchParams.get('token')

    if (!token) {
      // 没有 token，不处理回调
      return
    }

    setLoading(true)
    // 直接使用 auth-center token 调用 /api/v1/user/me
    // AuthCenterMiddleware 会自动验证 token 并获取用户信息
    fetch(`${API_BASE_URL}/api/v1/user/me`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
      .then(res => res.json())
      .then(data => {
        if (!data.success) {
          setError(data.error || '登录失败')
          return
        }

        // 保存 auth-center token
        login(
          token,
          '', // 不需要 refresh token
          {
            id: data.data.id,
            authCenterUserId: data.data.authCenterUserId,
            unionId: data.data.unionId || null,
            nickname: data.data.nickname,
            avatarUrl: data.data.avatarUrl,
            profile: data.data.profile || {},
            roles: data.data.roles || [],
            status: data.data.status,
            lastLoginAt: data.data.lastLoginAt,
            lastLoginIp: data.data.lastLoginIp,
            fixedInvitationCode: '',
            createdAt: data.data.createdAt,
            updatedAt: data.data.updatedAt,
          }
        )

        // 清除 URL 中的 token 参数
        navigate('/dashboard', { replace: true })
      })
      .catch((err) => {
        console.error('登录失败:', err)
        setError('登录失败，请重试')
      })
      .finally(() => {
        setLoading(false)
      })
  }, [searchParams, login, navigate])

  const handleWeChatLogin = () => {
    setError('')
    setLoading(true)

    // 重定向到后端的微信登录接口
    // 后端会重定向到 auth-center 进行微信授权
    // 授权成功后 auth-center 会回调到前端 /login?token=xxx
    window.location.href = `${API_BASE_URL}/api/v1/auth/wechat/login`
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4">
      <div className="max-w-md w-full">
        {/* Logo和标题 */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">PR Business</h1>
          <p className="text-gray-600">营销任务管理平台</p>
        </div>

        {/* 登录卡片 */}
        <Card className="shadow-xl">
          <CardHeader className="space-y-4">
            {/* 登录方式切换 */}
            <div className="flex bg-muted rounded-lg p-1">
              <Button
                onClick={() => setLoginType('wechat')}
                variant={loginType === 'wechat' ? 'default' : 'ghost'}
                className="flex-1"
              >
                微信登录
              </Button>
            </div>

            {/* 错误提示 */}
            {error && (
              <Badge variant="destructive" className="w-full py-2 px-3 justify-center">
                {error}
              </Badge>
            )}
          </CardHeader>

          <CardContent>
            {/* 微信登录 */}
            <div className="space-y-4">
              <div className="text-center text-muted-foreground mb-6">
                <p>点击下方按钮使用微信账号登录</p>
              </div>
              <Button
                onClick={handleWeChatLogin}
                disabled={loading}
                className="w-full bg-green-500 hover:bg-green-600"
                size="lg"
              >
                <svg className="w-6 h-6 mr-2" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.27-.027-.407-.03zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.969-.982z" />
                </svg>
                {loading ? '登录中...' : '微信登录'}
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* 底部提示 */}
        <p className="text-center text-muted-foreground text-sm mt-6">
          登录即表示同意《用户协议》和《隐私政策》
        </p>
      </div>
    </div>
  )
}
