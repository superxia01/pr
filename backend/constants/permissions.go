package constants

// 权限码常量定义
// 用于服务商员工和商家员工的细粒度权限控制

const (
	// 达人管理权限组（3个）
	PermissionManageCreators    = "MANAGE_CREATORS"    // 管理达人（已有）
	PermissionEditCreatorInfo    = "EDIT_CREATOR_INFO"    // 编辑达人资料
	PermissionDeleteCreator     = "DELETE_CREATOR"     // 删除达人

	// 活动管理权限组（4个）
	PermissionPublishCampaign    = "PUBLISH_CAMPAIGN"    // 发布活动
	PermissionEditCommission     = "EDIT_CAMPAIGN_COMMISSION" // 编辑佣金分配
	PermissionEditCampaignInfo   = "EDIT_CAMPAIGN_INFO"   // 编辑活动信息
	PermissionDeleteCampaign     = "DELETE_CAMPAIGN"     // 删除活动

	// 任务管理权限组（2个）
	PermissionReviewTask         = "REVIEW_TASK"         // 审核任务
	PermissionViewAllTasks       = "VIEW_ALL_TASKS"       // 查看所有任务

	// 财务管理权限组（3个）
	PermissionRecharge            = "RECHARGE"            // 充值
	PermissionApproveWithdrawal  = "APPROVE_WITHDRAWAL"  // 审核提现
	PermissionViewFinancialReports = "VIEW_FINANCIAL_REPORTS" // 查看财务报表

	// 组织管理权限组（3个）
	PermissionManageMerchant     = "MANAGE_MERCHANT"     // 管理商家
	PermissionManageStaff        = "MANAGE_STAFF"        // 管理员工
	PermissionCreateInvitationCode = "CREATE_INVITATION_CODE" // 创建邀请码
)

// PermissionDefinition 权限定义（用于前端显示）
type PermissionDefinition struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Group       string   `json:"group"`
	ForProvider []string `json:"forProvider"` // 适用对象：service_provider, merchant
	ForMerchant  []string `json:"forMerchant"`  // 适用对象：merchant
}

// GetAllPermissionDefinitions 获取所有权限定义
func GetAllPermissionDefinitions() []PermissionDefinition {
	return []PermissionDefinition{
		// 达人管理组（仅服务商员工）
		{
			Code:        PermissionManageCreators,
			Name:        "管理达人",
			Description: "查看和管理本组织所有达人（默认只能查看自己邀请的）",
			Group:       "达人管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{},
		},
		{
			Code:        PermissionEditCreatorInfo,
			Name:        "编辑达人资料",
			Description: "编辑达人昵称、头像等基本信息",
			Group:       "达人管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{},
		},
		{
			Code:        PermissionDeleteCreator,
			Name:        "删除达人",
			Description: "从本组织移除达人",
			Group:       "达人管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{},
		},

		// 活动管理组
		{
			Code:        PermissionPublishCampaign,
			Name:        "发布活动",
			Description: "将活动状态改为 OPEN（公开接受任务）",
			Group:       "活动管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},
		{
			Code:        PermissionEditCommission,
			Name:        "编辑佣金分配",
			Description: "设置和修改活动的佣金分配方案",
			Group:       "活动管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},
		{
			Code:        PermissionEditCampaignInfo,
			Name:        "编辑活动信息",
			Description: "编辑活动的标题、要求、平台、截止时间",
			Group:       "活动管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},
		{
			Code:        PermissionDeleteCampaign,
			Name:        "删除活动",
			Description: "删除活动",
			Group:       "活动管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},

		// 任务管理组
		{
			Code:        PermissionReviewTask,
			Name:        "审核任务",
			Description: "审核达人提交的任务作品",
			Group:       "任务管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},
		{
			Code:        PermissionViewAllTasks,
			Name:        "查看所有任务",
			Description: "查看本组织所有活动的任务情况",
			Group:       "任务管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},

		// 财务管理组
		{
			Code:        PermissionRecharge,
			Name:        "充值",
			Description: "为积分账户充值",
			Group:       "财务管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},
		{
			Code:        PermissionApproveWithdrawal,
			Name:        "审核提现",
			Description: "审批用户的提现申请（仅超级管理员）",
			Group:       "财务管理",
			ForProvider: []string{}, // 仅超级管理员
			ForMerchant:  []string{},
		},
		{
			Code:        PermissionViewFinancialReports,
			Name:        "查看财务报表",
			Description: "查看财务统计和报表",
			Group:       "财务管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},

		// 组织管理组
		{
			Code:        PermissionManageMerchant,
			Name:        "管理商家",
			Description: "创建、编辑、删除商家",
			Group:       "组织管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{}, // 商家员工不能管理商家
		},
		{
			Code:        PermissionManageStaff,
			Name:        "管理员工",
			Description: "添加、删除员工，管理员工权限",
			Group:       "组织管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},
		{
			Code:        PermissionCreateInvitationCode,
			Name:        "创建邀请码",
			Description: "生成各类角色的邀请码",
			Group:       "组织管理",
			ForProvider: []string{"SERVICE_PROVIDER_STAFF"},
			ForMerchant:  []string{"MERCHANT_STAFF"},
		},
	}
}

// GetPermissionDefinitionsForRole 根据角色获取可用的权限定义
func GetPermissionDefinitionsForRole(role string) []PermissionDefinition {
	allPerms := GetAllPermissionDefinitions()
	result := make([]PermissionDefinition, 0)

	for _, perm := range allPerms {
		if role == "SERVICE_PROVIDER_STAFF" {
			// 检查是否适用于服务商员工
			for _, p := range perm.ForProvider {
				if p == "SERVICE_PROVIDER_STAFF" {
					result = append(result, perm)
					break
				}
			}
		} else if role == "MERCHANT_STAFF" {
			// 检查是否适用于商家员工
			for _, p := range perm.ForMerchant {
				if p == "MERCHANT_STAFF" {
					result = append(result, perm)
					break
				}
			}
		}
	}

	return result
}
