// 角色配置
export const ROLE_CONFIG: Record<string, { name: string; description: string; color: string }> = {
  SUPER_ADMIN: {
    name: '超级管理员',
    description: '系统最高权限',
    color: 'purple',
  },
  SP_ADMIN: {
    name: '服务商管理员',
    description: '管理服务商和员工',
    color: 'blue',
  },
  SP_STAFF: {
    name: '服务商员工',
    description: '服务商工作人员',
    color: 'blue',
  },
  MERCHANT_ADMIN: {
    name: '商家管理员',
    description: '管理商家和员工',
    color: 'green',
  },
  MERCHANT_STAFF: {
    name: '商家员工',
    description: '商家工作人员',
    color: 'green',
  },
  CREATOR: {
    name: '达人',
    description: '内容创作者',
    color: 'orange',
  },
}

// 获取角色显示名称
export function getRoleName(role: string): string {
  return ROLE_CONFIG[role]?.name || role
}

// 获取角色描述
export function getRoleDescription(role: string): string {
  return ROLE_CONFIG[role]?.description || ''
}

// 获取角色颜色
export function getRoleColor(role: string): string {
  return ROLE_CONFIG[role]?.color || 'gray'
}

// 获取角色颜色类名（用于Tailwind）
export function getRoleColorClass(role: string): string {
  const color = getRoleColor(role)
  const colorMap: Record<string, string> = {
    purple: 'bg-purple-100 text-purple-700 border-purple-200',
    blue: 'bg-blue-100 text-blue-700 border-blue-200',
    green: 'bg-green-100 text-green-700 border-green-200',
    orange: 'bg-orange-100 text-orange-700 border-orange-200',
    gray: 'bg-gray-100 text-gray-700 border-gray-200',
  }
  return colorMap[color] || colorMap.gray
}
