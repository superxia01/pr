// 角色常量（与后端保持一致）
export const ROLES = {
  SUPER_ADMIN: 'SUPER_ADMIN',
  SERVICE_PROVIDER_ADMIN: 'SERVICE_PROVIDER_ADMIN',
  SERVICE_PROVIDER_STAFF: 'SERVICE_PROVIDER_STAFF',
  MERCHANT_ADMIN: 'MERCHANT_ADMIN',
  MERCHANT_STAFF: 'MERCHANT_STAFF',
  CREATOR: 'CREATOR',
} as const

export type Role = typeof ROLES[keyof typeof ROLES]

// 角色配置
export const ROLE_CONFIG: Record<string, { name: string; shortName: string; description: string; color: string }> = {
  [ROLES.SUPER_ADMIN]: {
    name: '超级管理员',
    shortName: '超管',
    description: '系统最高权限',
    color: 'purple',
  },
  [ROLES.SERVICE_PROVIDER_ADMIN]: {
    name: '服务商管理员',
    shortName: '服务商管理员',
    description: '管理服务商和员工',
    color: 'blue',
  },
  [ROLES.SERVICE_PROVIDER_STAFF]: {
    name: '服务商员工',
    shortName: '服务商员工',
    description: '服务商工作人员',
    color: 'blue',
  },
  [ROLES.MERCHANT_ADMIN]: {
    name: '商家管理员',
    shortName: '商家管理员',
    description: '管理商家和员工',
    color: 'green',
  },
  [ROLES.MERCHANT_STAFF]: {
    name: '商家员工',
    shortName: '商家员工',
    description: '商家工作人员',
    color: 'green',
  },
  [ROLES.CREATOR]: {
    name: '达人',
    shortName: '达人',
    description: '内容创作者',
    color: 'orange',
  },
}

// 工具函数
export function getRoleDisplayName(role: string): string {
  return ROLE_CONFIG[role]?.name || role
}

export function getRoleShortName(role: string): string {
  return ROLE_CONFIG[role]?.shortName || role
}

export function getRoleDescription(role: string): string {
  return ROLE_CONFIG[role]?.description || ''
}

export function getRoleColor(role: string): string {
  return ROLE_CONFIG[role]?.color || 'gray'
}

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

export function isValidRole(role: string): role is Role {
  return Object.values(ROLES).includes(role as Role)
}

export function getAllRoles(): Role[] {
  return Object.values(ROLES)
}

// 向后兼容的旧常量（废弃，使用 ROLES 替代）
/** @deprecated 使用 ROLES.SERVICE_PROVIDER_ADMIN 替代 */
export const SP_ADMIN = ROLES.SERVICE_PROVIDER_ADMIN

/** @deprecated 使用 ROLES.SERVICE_PROVIDER_STAFF 替代 */
export const SP_STAFF = ROLES.SERVICE_PROVIDER_STAFF

// 旧函数名（向后兼容）
/** @deprecated 使用 getRoleDisplayName 替代 */
export const getRoleName = getRoleDisplayName
