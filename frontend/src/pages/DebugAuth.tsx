import { useState, useEffect } from 'react'
import { authApi } from '../services/api'

// 简单的JWT解析函数（不验证签名）
function parseJWT(token: string) {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
    return JSON.parse(jsonPayload)
  } catch (e) {
    return null
  }
}

export default function DebugAuth() {
  const [debugInfo, setDebugInfo] = useState<any>({})

  useEffect(() => {
    // 收集所有调试信息
    const collectDebugInfo = () => {
      const accessToken = localStorage.getItem('accessToken')
      const refreshToken = localStorage.getItem('refreshToken')
      const userStr = localStorage.getItem('user')

      let decodedToken = null
      if (accessToken) {
        decodedToken = parseJWT(accessToken)
      }

      let user = null
      if (userStr) {
        try {
          user = JSON.parse(userStr)
        } catch (e) {
          user = { error: 'Failed to parse user' }
        }
      }

      setDebugInfo({
        localStorage: {
          accessToken: accessToken ? `${accessToken.substring(0, 30)}... (length: ${accessToken.length})` : null,
          refreshToken: refreshToken ? `${refreshToken.substring(0, 30)}... (length: ${refreshToken.length})` : null,
          user: user ? `Present (id: ${user.id}, nickname: ${user.nickname})` : null,
        },
        decodedToken: {
          header: decodedToken ? 'See below' : null,
          payload: decodedToken,
        },
        currentUser: user,
        timestamp: new Date().toISOString(),
      })
    }

    collectDebugInfo()
  }, [])

  const testAPI = async () => {
    console.log('🧪 测试 /api/v1/user/me API...')
    try {
      const response = await authApi.getCurrentUser()
      console.log('✅ API调用成功:', response)
      alert('✅ API调用成功！\n' + JSON.stringify(response, null, 2))
    } catch (error: any) {
      console.error('❌ API调用失败:', error)
      alert('❌ API调用失败！\n' +
        'Status: ' + (error.response?.status) + '\n' +
        'Error: ' + (error.response?.data?.error || error.message))
    }
  }

  const testLoginAPI = async () => {
    console.log('🧪 测试密码登录API...')
    try {
      const phoneNumber = prompt('请输入手机号:')
      const password = prompt('请输入密码:')
      if (!phoneNumber || !password) return

      const response = await authApi.passwordLogin({ phoneNumber, password })
      console.log('✅ 登录成功:', response)
      alert('✅ 登录成功！\n' +
        'AccessToken: ' + response.accessToken.substring(0, 30) + '...\n' +
        'RefreshToken: ' + response.refreshToken.substring(0, 30) + '...\n' +
        'UserID: ' + response.userId + '\n' +
        'Roles: ' + response.roles.join(', ') + '\n' +
        'CurrentRole: ' + response.currentRole)

      // 刷新页面
      setTimeout(() => window.location.reload(), 2000)
    } catch (error: any) {
      console.error('❌ 登录失败:', error)
      alert('❌ 登录失败！\n' + (error.response?.data?.error || error.message))
    }
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">🔍 认证调试工具</h1>

        {/* 操作按钮 */}
        <div className="bg-white rounded-lg shadow p-6 mb-6 space-x-4">
          <button
            onClick={testAPI}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            测试 /api/v1/user/me
          </button>
          <button
            onClick={testLoginAPI}
            className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
          >
            测试登录
          </button>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            刷新页面
          </button>
          <button
            onClick={() => {
              console.log('完整调试信息:', debugInfo)
              alert('调试信息已输出到控制台（F12）')
            }}
            className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600"
          >
            输出到控制台
          </button>
        </div>

        {/* localStorage状态 */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-xl font-bold mb-4">📦 localStorage 状态</h2>
          <pre className="bg-gray-100 p-4 rounded overflow-x-auto">
            {JSON.stringify(debugInfo.localStorage, null, 2)}
          </pre>
        </div>

        {/* 解码后的Token */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-xl font-bold mb-4">🔑 解码后的 JWT Token</h2>
          {debugInfo.decodedToken?.payload ? (
            <pre className="bg-gray-100 p-4 rounded overflow-x-auto">
              {JSON.stringify(debugInfo.decodedToken.payload, null, 2)}
            </pre>
          ) : (
            <p className="text-red-500">❌ 无法解析Token（可能没有token或格式错误）</p>
          )}
        </div>

        {/* 当前用户信息 */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-xl font-bold mb-4">👤 当前用户信息</h2>
          <pre className="bg-gray-100 p-4 rounded overflow-x-auto">
            {JSON.stringify(debugInfo.currentUser, null, 2)}
          </pre>
        </div>

        {/* 检查清单 */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold mb-4">✅ 问题检查清单</h2>
          <ul className="space-y-2">
            <li className={debugInfo.localStorage?.accessToken ? 'text-green-600' : 'text-red-600'}>
              {debugInfo.localStorage?.accessToken ? '✅' : '❌'} accessToken 存在于 localStorage
            </li>
            <li className={debugInfo.localStorage?.refreshToken ? 'text-green-600' : 'text-red-600'}>
              {debugInfo.localStorage?.refreshToken ? '✅' : '❌'} refreshToken 存在于 localStorage
            </li>
            <li className={debugInfo.localStorage?.user ? 'text-green-600' : 'text-red-600'}>
              {debugInfo.localStorage?.user ? '✅' : '❌'} user 存在于 localStorage
            </li>
            <li className={debugInfo.decodedToken?.payload ? 'text-green-600' : 'text-red-600'}>
              {debugInfo.decodedToken?.payload ? '✅' : '❌'} JWT Token 可以成功解析
            </li>
            <li className={debugInfo.decodedToken?.payload?.userId ? 'text-green-600' : 'text-red-600'}>
              {debugInfo.decodedToken?.payload?.userId ? '✅' : '❌'} Token 包含 userId
            </li>
            <li className={debugInfo.decodedToken?.payload?.role ? 'text-green-600' : 'text-yellow-600'}>
              {debugInfo.decodedToken?.payload?.role ? '✅' : '⚠️'} Token 包含 role: {debugInfo.decodedToken?.payload?.role || '(空)'}
            </li>
            <li className={debugInfo.decodedToken?.payload?.roles?.length > 0 ? 'text-green-600' : 'text-yellow-600'}>
              {debugInfo.decodedToken?.payload?.roles?.length > 0 ? '✅' : '⚠️'} Token 包含 roles数组: {JSON.stringify(debugInfo.decodedToken?.payload?.roles || [])}
            </li>
            <li className={debugInfo.currentUser?.id ? 'text-green-600' : 'text-red-600'}>
              {debugInfo.currentUser?.id ? '✅' : '❌'} 用户对象包含 id
            </li>
            <li className={debugInfo.currentUser?.currentRole ? 'text-green-600' : 'text-yellow-600'}>
              {debugInfo.currentUser?.currentRole ? '✅' : '⚠️'} 用户对象包含 currentRole: {debugInfo.currentUser?.currentRole || '(空)'}
            </li>
            <li className={debugInfo.currentUser?.roles?.length > 0 ? 'text-green-600' : 'text-yellow-600'}>
              {debugInfo.currentUser?.roles?.length > 0 ? '✅' : '⚠️'} 用户对象包含 roles数组: {JSON.stringify(debugInfo.currentUser?.roles || [])}
            </li>
          </ul>
        </div>

        {/* 常见问题 */}
        <div className="bg-yellow-50 rounded-lg shadow p-6 mt-6">
          <h2 className="text-xl font-bold mb-4">💡 常见401错误原因</h2>
          <ol className="list-decimal list-inside space-y-2 text-gray-700">
            <li><strong>Token未保存</strong>: 登录成功但localStorage中没有token</li>
            <li><strong>Token格式错误</strong>: Bearer前缀缺失或格式不对</li>
            <li><strong>Token已过期</strong>: 检查Token的exp时间</li>
            <li><strong>后端JWT密钥不一致</strong>: 生成和验证使用了不同的密钥</li>
            <li><strong>用户角色为空</strong>: 新用户可能没有角色，导致某些API返回401</li>
            <li><strong>CORS问题</strong>: 跨域请求时token没有正确发送</li>
          </ol>
        </div>
      </div>
    </div>
  )
}
