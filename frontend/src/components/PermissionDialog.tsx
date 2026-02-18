import { useState, useEffect } from 'react'
import { Button } from './ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from './ui/dialog'

interface Permission {
  code: string
  name: string
  description: string
  group: string
  forProvider?: string[]
  forMerchant?: string[]
}

interface PermissionDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  staffName: string
  staffRole: string // SERVICE_PROVIDER_STAFF or MERCHANT_STAFF
  currentPermissions: string[]
  onSave: (permissions: string[]) => Promise<void>
}

export function PermissionDialog({
  open,
  onOpenChange,
  staffName,
  staffRole,
  currentPermissions = [],
  onSave,
}: PermissionDialogProps) {
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [selectedCodes, setSelectedCodes] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)

  // 加载权限列表
  useEffect(() => {
    if (open) {
      loadPermissions()
    }
  }, [open])

  // 初始化已选权限
  useEffect(() => {
    setSelectedCodes([...currentPermissions])
  }, [currentPermissions, open])

  const loadPermissions = async () => {
    setLoading(true)
    try {
      // 根据员工角色选择API
      let allPermissions: Permission[] = []
      if (staffRole === 'SERVICE_PROVIDER_STAFF') {
        allPermissions = await import('../services/api').then(m => m.serviceProviderApi.getPermissions())
      } else if (staffRole === 'MERCHANT_STAFF') {
        allPermissions = await import('../services/api').then(m => m.merchantApi.getPermissions())
      }

      // 根据角色过滤权限
      let filteredPermissions = allPermissions
      if (staffRole === 'SERVICE_PROVIDER_STAFF') {
        filteredPermissions = allPermissions.filter((p: Permission) =>
          p.forProvider && p.forProvider.includes('SERVICE_PROVIDER_STAFF')
        )
      } else if (staffRole === 'MERCHANT_STAFF') {
        filteredPermissions = allPermissions.filter((p: Permission) =>
          p.forMerchant && p.forMerchant.includes('MERCHANT_STAFF')
        )
      }

      setPermissions(filteredPermissions)
    } catch (err: any) {
      console.error('加载权限列表失败:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      await onSave(selectedCodes)
      onOpenChange(false)
    } catch (err: any) {
      console.error('保存权限失败:', err)
      alert('保存权限失败：' + (err.response?.data?.error || err.message))
    } finally {
      setSaving(false)
    }
  }

  const handleTogglePermission = (code: string) => {
    setSelectedCodes(prev => {
      if (prev.includes(code)) {
        return prev.filter(c => c !== code)
      } else {
        return [...prev, code]
      }
    })
  }

  // 按分组显示权限
  const groupedPermissions = permissions.reduce((acc, permission) => {
    if (!acc[permission.group]) {
      acc[permission.group] = []
    }
    acc[permission.group].push(permission)
    return acc
  }, {} as Record<string, Permission[]>)

  if (loading) {
    return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="max-w-2xl max-h-[80vh]">
          <div className="flex items-center justify-center py-12">
            <div className="text-center">
              <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              <p className="mt-2 text-sm text-gray-500">加载权限列表...</p>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[80vh]">
        <DialogHeader>
          <DialogTitle>管理员工权限 - {staffName}</DialogTitle>
        </DialogHeader>

        <div className="max-h-[60vh] overflow-y-auto px-6 py-4">
          {Object.entries(groupedPermissions).map(([group, perms]) => (
            <div key={group} className="mb-6">
              <h3 className="text-sm font-semibold text-gray-900 mb-3 pb-2 border-b">
                {group}
              </h3>
              <div className="space-y-3">
                {perms.map((permission) => (
                  <div key={permission.code} className="flex items-start space-x-3">
                    <input
                      type="checkbox"
                      id={`perm-${permission.code}`}
                      checked={selectedCodes.includes(permission.code)}
                      onChange={() => handleTogglePermission(permission.code)}
                      className="mt-1 h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <div
                      className="flex-1 cursor-pointer"
                      onClick={() => handleTogglePermission(permission.code)}
                    >
                      <label
                        htmlFor={`perm-${permission.code}`}
                        className="text-sm font-medium text-gray-900 cursor-pointer"
                      >
                        {permission.name}
                      </label>
                      <p className="text-xs text-gray-500 mt-0.5">
                        {permission.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}

          {permissions.length === 0 && (
            <div className="text-center py-8 text-sm text-gray-500">
              暂无可授权的权限
            </div>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            取消
          </Button>
          <Button onClick={handleSave} disabled={saving}>
            {saving ? '保存中...' : '保存'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
