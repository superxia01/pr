import { useState, useEffect } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { invitationApi } from '../services/api'
import { Button } from '../components/ui/button'
import { Badge } from '../components/ui/badge'

interface InvitationRelationship {
	id: string
	inviterId: string
	inviterRole: string
	inviteeId: string
	inviteeRole: string
	organizationId: string | null
	organizationType: string | null
	invitationCode: string
	invitedAt: string
	status: string
	invitee?: {
		id: string
		nickname: string
		avatarUrl: string
	}
}

interface InvitationCode {
	code: string
	role: string
	roleLabel: string
	orgId?: string
	orgName?: string
	orgType?: string
}

export default function InvitationCodes() {
	const { user } = useAuth()
	const [invitationCodes, setInvitationCodes] = useState<InvitationCode[]>([])
	const [invitations, setInvitations] = useState<InvitationRelationship[]>([])
	const [loading, setLoading] = useState(true)
	const [error, setError] = useState('')
	const [copiedCode, setCopiedCode] = useState<string | null>(null)
	const [currentPage, setCurrentPage] = useState(1)
	const [total, setTotal] = useState(0)
	const pageSize = 20

	useEffect(() => {
		loadMyFixedInvitationCodes()
		loadMyInvitations()
	}, [])

	const loadMyFixedInvitationCodes = async () => {
		try {
			setLoading(true)
			const response = await invitationApi.getMyFixedInvitationCodes()
			setInvitationCodes(response.invitationCodes || [])
		} catch (err: any) {
			setError(err.response?.data?.error || '加载邀请码失败')
		} finally {
			setLoading(false)
		}
	}

	const loadMyInvitations = async (page: number = 1) => {
		try {
			setLoading(true)
			const response = await invitationApi.getMyInvitations({
				page,
				page_size: pageSize,
			})
			setInvitations(response.invitations || [])
			setTotal(response.total || 0)
			setCurrentPage(page)
		} catch (err: any) {
			setError(err.response?.data?.error || '加载邀请列表失败')
		} finally {
			setLoading(false)
		}
	}

	const handleCopyCode = async (code: string) => {
		try {
			await navigator.clipboard.writeText(code)
			setCopiedCode(code)
			setTimeout(() => setCopiedCode(null), 2000)
		} catch (err) {
			// Fallback for older browsers
			const textArea = document.createElement('textarea')
			textArea.value = code
			document.body.appendChild(textArea)
			textArea.select()
			document.execCommand('copy')
			document.body.removeChild(textArea)
			setCopiedCode(code)
			setTimeout(() => setCopiedCode(null), 2000)
		}
	}

	const getRoleLabel = (role: string) => {
		const labels: Record<string, string> = {
			SUPER_ADMIN: '超级管理员',
			'SERVICE_PROVIDER_ADMIN': '服务商管理员',
			'SERVICE_PROVIDER_STAFF': '服务商员工',
			'MERCHANT_ADMIN': '商家管理员',
			'MERCHANT_STAFF': '商家员工',
			'CREATOR': '达人',
		}
		return labels[role] || role
	}

	const getOrganizationLabel = (orgType: string | null, orgId: string | null) => {
		if (!orgType || !orgId) return '-'
		if (orgType === 'service_provider') return '服务商'
		if (orgType === 'merchant') return '商家'
		if (orgType === 'CAMPAIGN') return '营销活动'
		return '-'
	}

	const handleNextPage = () => {
		if (currentPage * pageSize < total) {
			loadMyInvitations(currentPage + 1)
		}
	}

	const handlePrevPage = () => {
		if (currentPage > 1) {
			loadMyInvitations(currentPage - 1)
		}
	}

	if (!user) {
		return (
			<div className="min-h-screen bg-background p-8">
				<div className="max-w-7xl mx-auto">
					<div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
						未登录
					</div>
				</div>
			</div>
		)
	}

	return (
		<div className="min-h-screen bg-background">
			{/* 顶部导航 */}
			<nav className="bg-card shadow-sm border-b">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
					<div className="flex justify-between h-16">
						<div className="flex items-center">
							<h1 className="text-xl font-bold">邀请码管理</h1>
						</div>
						<div className="flex items-center gap-4">
							<a href="/dashboard" className="text-sm text-muted-foreground hover:text-foreground">
								返回工作台
							</a>
						</div>
					</div>
				</div>
			</nav>

			{/* 主要内容 */}
			<main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				{/* 错误提示 */}
				{error && (
					<div className="mb-4">
						<Badge variant="destructive" className="w-full py-2 px-3 justify-center">
							{error}
						</Badge>
					</div>
				)}

				{/* 我的固定邀请码 */}
				<div className="mb-8 bg-card rounded-lg border p-6">
					<h2 className="text-xl font-bold mb-4">我的邀请码</h2>
					<p className="text-sm text-muted-foreground mb-4">
						根据您的角色和管理组织，系统已为您生成以下邀请码。
					</p>

					{invitationCodes.length === 0 ? (
						<div className="text-center py-8 text-sm text-muted-foreground">
							您当前的角色无法邀请他人
						</div>
					) : (
						<div className="space-y-3">
							{invitationCodes.map((ic) => (
								<div key={ic.code} className="border rounded-lg p-4 bg-muted/30">
									<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
										{/* 角色信息 */}
										<div>
											<div className="text-xs text-muted-foreground mb-1">邀请角色</div>
											<div className="text-lg font-bold">{ic.roleLabel}</div>
										</div>

										{/* 组织信息（如果有） */}
										{ic.orgName && (
											<div>
												<div className="text-xs text-muted-foreground mb-1">绑定组织</div>
												<div className="text-sm font-medium">{ic.orgName}</div>
											</div>
										)}

										{/* 邀请码 */}
										<div className="md:col-span-2">
											<div className="text-xs text-muted-foreground mb-1">邀请码</div>
											<div className="flex items-center gap-2">
												<code className="flex-1 px-3 py-2 bg-background rounded text-sm font-mono text-center break-all">
													{ic.code}
												</code>
												<Button
													variant="ghost"
													size="sm"
													onClick={() => handleCopyCode(ic.code)}
												>
													{copiedCode === ic.code ? (
														<span className="text-green-600">✓ 已复制</span>
													) : (
														<>
															<svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
																<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 12-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8" />
															</svg>
															复制
														</>
													)}
												</Button>
											</div>
										</div>
									</div>
								</div>
							))}
						</div>
					)}
				</div>

				{/* 邀请列表 */}
				<div className="bg-card rounded-lg border">
					<div className="p-6 border-b">
						<h2 className="text-xl font-bold">我邀请的人</h2>
						<p className="text-sm text-muted-foreground mt-1">
							查看所有通过你的邀请码注册的用户
						</p>
					</div>

					{loading ? (
						<div className="text-center py-12">
							<div className="text-muted-foreground">加载中...</div>
						</div>
					) : invitations.length === 0 ? (
						<div className="text-center py-12">
							<div className="text-muted-foreground">
								暂无邀请记录
							</div>
						</div>
					) : (
						<>
							<div className="overflow-x-auto">
								<table className="w-full">
									<thead className="bg-muted">
										<tr>
											<th className="px-4 py-3 text-left text-sm font-medium">被邀请人</th>
											<th className="px-4 py-3 text-left text-sm font-medium">角色</th>
											<th className="px-4 py-3 text-left text-sm font-medium">绑定组织</th>
											<th className="px-4 py-3 text-left text-sm font-medium">邀请码</th>
											<th className="px-4 py-3 text-left text-sm font-medium">邀请时间</th>
											<th className="px-4 py-3 text-left text-sm font-medium">状态</th>
										</tr>
									</thead>
									<tbody>
										{invitations.map((inv) => (
											<tr key={inv.id} className="border-b hover:bg-muted/50">
												<td className="px-4 py-3">
													<div className="flex items-center gap-3">
														<img
															src={inv.invitee?.avatarUrl || '/default-avatar.png'}
															alt={inv.invitee?.nickname || '用户'}
															className="w-8 h-8 rounded-full"
														/>
														<span className="text-sm font-medium">{inv.invitee?.nickname || '未知用户'}</span>
													</div>
												</td>
												<td className="px-4 py-3">
													<Badge>{getRoleLabel(inv.inviteeRole)}</Badge>
												</td>
												<td className="px-4 py-3 text-sm">
													{getOrganizationLabel(inv.organizationType, inv.organizationId)}
												</td>
												<td className="px-4 py-3">
													<code className="text-xs bg-muted px-2 py-1 rounded">{inv.invitationCode}</code>
												</td>
												<td className="px-4 py-3 text-sm text-muted-foreground">
													{new Date(inv.invitedAt).toLocaleString('zh-CN')}
												</td>
												<td className="px-4 py-3">
													<Badge variant={inv.status === 'active' ? 'default' : 'secondary'}>
														{inv.status === 'active' ? '激活' : '已取消'}
													</Badge>
												</td>
											</tr>
										))}
									</tbody>
								</table>
							</div>

							{/* 分页 */}
							{total > pageSize && (
								<div className="flex justify-center items-center gap-4 p-4 border-t">
									<Button
										variant="outline"
										size="sm"
										onClick={handlePrevPage}
										disabled={currentPage <= 1}
									>
										上一页
									</Button>
									<span className="text-sm text-muted-foreground">
										第 {currentPage} 页，共 {Math.ceil(total / pageSize)} 页
									</span>
									<Button
										variant="outline"
										size="sm"
										onClick={handleNextPage}
										disabled={currentPage * pageSize >= total}
									>
										下一页
									</Button>
								</div>
							)}
						</>
					)}
				</div>
			</main>
		</div>
	)
}
