import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { cn } from '../lib/utils'
import { getRoleName } from '../lib/roles'
import {
  LayoutDashboard,
  Gift,
  Users,
  Building2,
  Briefcase,
  Sparkles,
  Coins,
  Receipt,
  ShieldCheck,
  X,
  ChevronRight,
  LogOut,
  CreditCard,
} from 'lucide-react'

interface SidebarProps {
  isOpen: boolean
  onClose: () => void
}

export function Sidebar({ isOpen, onClose }: SidebarProps) {
  const location = useLocation()
  const { user, logout } = useAuth()

  // 多角色：按「拥有的角色」显示所有对应菜单（不再使用角色切换器）
  const roles = user?.roles ?? []
  const isSuperAdmin = () => roles.some((r) => r.toUpperCase() === 'SUPER_ADMIN')
  const isMerchantAdmin = () => roles.some((r) => r.toUpperCase() === 'MERCHANT_ADMIN')
  const isServiceProviderAdmin = () => roles.some((r) => r.toUpperCase() === 'SERVICE_PROVIDER_ADMIN')
  const isCreator = () => roles.some((r) => r.toUpperCase() === 'CREATOR')
  const isMerchantStaff = () => roles.some((r) => r.toUpperCase() === 'MERCHANT_STAFF')
  const isBasicUser = () => roles.some((r) => r.toUpperCase() === 'BASIC_USER')

  const canManageInvitations = () => {
    return isSuperAdmin() || isServiceProviderAdmin() || isMerchantAdmin() || isBasicUser()
  }

  const navItems = [
    // 工作台 - 始终显示
    {
      title: '工作台',
      href: '/',
      icon: LayoutDashboard,
    },
    // 达人相关（拥有达人角色即显示）
    ...(isCreator() ? [{
      title: '任务大厅',
      href: '/task-hall',
      icon: Briefcase,
    }, {
      title: '我的任务',
      href: '/my-tasks',
      icon: Sparkles,
    }, {
      title: '达人中心',
      href: '/creator',
      icon: Sparkles,
    }] : []),
    // 活动创建（商家管理员、商家员工、服务商管理员）
    ...(isMerchantAdmin() || isMerchantStaff() || isServiceProviderAdmin() ? [{
      title: '创建活动',
      href: '/create-campaign',
      icon: Sparkles,
    }] : []),
    // 活动审核（仅服务商管理员）
    ...(isServiceProviderAdmin() ? [{
      title: '活动审核',
      href: '/campaign-approval',
      icon: ShieldCheck,
    }] : []),
    ...(isMerchantAdmin() ? [{
      title: '商家信息',
      href: '/merchant',
      icon: Building2,
    }] : []),
    // 服务商相关
    ...(isServiceProviderAdmin() ? [{
      title: '服务商信息',
      href: '/service-provider',
      icon: Building2,
    }] : []),
    // 超级管理员 / 服务商管理员
    ...(isSuperAdmin() || isServiceProviderAdmin() ? [{
      title: '商家管理',
      href: '/merchants',
      icon: Briefcase,
    }] : []),
    // 超级管理员专有
    ...(isSuperAdmin() ? [{
      title: '服务商管理',
      href: '/service-providers',
      icon: Building2,
    }, {
      title: '用户管理',
      href: '/user-management',
      icon: Users,
    }] : []),
    ...(canManageInvitations() ? [{
      title: '邀请码管理',
      href: '/invitations',
      icon: Gift,
    }] : []),
    // 财务相关 - 始终显示
    {
      title: '充值',
      href: '/recharge',
      icon: CreditCard,
    },
    {
      title: '积分明细',
      href: '/credit-transactions',
      icon: Receipt,
    },
    {
      title: '提现记录',
      href: '/withdrawals',
      icon: Coins,
    },
    ...(isSuperAdmin() ? [{
      title: '提现审核',
      href: '/withdrawal-review',
      icon: ShieldCheck,
    }] : []),
  ]

  return (
    <>
      {/* 移动端遮罩 */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={onClose}
        />
      )}

      {/* 侧边栏 */}
      <aside
        className={cn(
          "fixed top-0 left-0 z-50 h-screen w-64 bg-gray-900 text-gray-100 transition-transform duration-200 ease-in-out lg:translate-x-0",
          isOpen ? "translate-x-0" : "-translate-x-full"
        )}
      >
        <div className="flex h-full flex-col">
          {/* Logo */}
          <div className="flex h-16 items-center justify-between px-6 border-b border-gray-800">
            <Link to="/" className="flex items-center gap-2 font-semibold text-white">
              <div className="h-8 w-8 rounded-lg bg-blue-600 flex items-center justify-center">
                <span className="text-white font-bold">PR</span>
              </div>
              <span className="text-lg">Business</span>
            </Link>
            <button
              onClick={onClose}
              className="lg:hidden text-gray-400 hover:text-white"
            >
              <X className="h-6 w-6" />
            </button>
          </div>

          {/* 导航菜单 */}
          <nav className="flex-1 overflow-y-auto px-3 py-4">
            <div className="space-y-1">
              {navItems.map((item) => {
                const isActive = location.pathname === item.href
                const Icon = item.icon

                return (
                  <Link
                    key={item.href}
                    to={item.href}
                    onClick={() => {
                      // 移动端点击后关闭侧边栏
                      if (window.innerWidth < 1024) {
                        onClose()
                      }
                    }}
                    className={cn(
                      "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
                      isActive
                        ? "bg-blue-600 text-white"
                        : "text-gray-300 hover:bg-gray-800 hover:text-white"
                    )}
                  >
                    <Icon className="h-5 w-5 flex-shrink-0" />
                    <span>{item.title}</span>
                    {isActive && <ChevronRight className="ml-auto h-4 w-4" />}
                  </Link>
                )
              })}
            </div>
          </nav>

          {/* 底部用户信息 */}
          <div className="border-t border-gray-800 p-4">
            <div className="flex items-center gap-3 mb-3">
              {/* 微信头像或首字母占位图 */}
              {user?.avatarUrl ? (
                <img
                  src={user.avatarUrl}
                  alt={user.nickname || '用户'}
                  className="h-10 w-10 rounded-full object-cover border-2 border-gray-700"
                />
              ) : (
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                  {(user?.nickname || '用户').charAt(0).toUpperCase()}
                </div>
              )}
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-white truncate">
                  {user?.nickname || '用户'}
                </p>
                <p className="text-xs text-gray-400 truncate">
                  {roles.length > 0
                    ? roles.length === 1
                      ? getRoleName(roles[0].toUpperCase())
                      : `多角色 (${roles.length})`
                    : ''}
                </p>
              </div>
            </div>

            <button
              onClick={logout}
              className="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-300 hover:bg-gray-800 hover:text-white rounded-lg transition-colors"
            >
              <LogOut className="h-4 w-4" />
              <span>退出登录</span>
            </button>
          </div>
        </div>
      </aside>
    </>
  )
}
