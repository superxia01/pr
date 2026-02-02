import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { authApi } from '../services/api'

export default function RoleSwitcher() {
  const { user, updateUser } = useAuth()
  const [switching, setSwitching] = useState(false)

  if (!user || user.roles.length <= 1) {
    return null
  }

  const handleSwitchRole = async (newRole: string) => {
    if (newRole === user.currentRole || switching) return

    setSwitching(true)
    try {
      const response = await authApi.switchRole({ newRole })
      localStorage.setItem('accessToken', response.accessToken)
      updateUser({
        ...user,
        currentRole: response.currentRole,
        lastUsedRole: response.lastUsedRole,
      })
      window.location.reload()
    } catch (error) {
      console.error('切换角色失败:', error)
    } finally {
      setSwitching(false)
    }
  }

  return (
    <div className="fixed top-4 right-4 z-50">
      <div className="bg-white rounded-lg shadow-lg p-4 min-w-[200px]">
        <h3 className="text-sm font-medium text-gray-700 mb-2">切换角色</h3>
        <div className="space-y-1">
          {user.roles.map((role) => (
            <button
              key={role}
              onClick={() => handleSwitchRole(role)}
              disabled={switching || role === user.currentRole}
              className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors ${
                role === user.currentRole
                  ? 'bg-blue-50 text-blue-700 font-medium'
                  : 'text-gray-700 hover:bg-gray-50'
              } ${switching ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
              {role === user.currentRole && <span className="mr-2">✓</span>}
              {role}
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}
