package routes

import (
	"pr-business/config"
	"pr-business/controllers"
	"pr-business/middlewares"
	"pr-business/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, redisClient interface{}, cfg *config.Config) {
	// 将 db 设置到全局上下文，供中间件使用
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// 初始化服务层
	validatorService := services.NewValidatorService(db)
	cashAccountService := services.NewCashAccountService(db, validatorService)
	systemAccountService := services.NewSystemAccountService(db, validatorService)
	auditService := services.NewAuditService(db)
	withdrawalEnhancedService := services.NewWithdrawalEnhancedService(
		db,
		validatorService,
		cashAccountService,
		systemAccountService,
	)

	// 初始化controllers
	authController := controllers.NewAuthController(cfg, db)
	invitationController := controllers.NewInvitationController(db, cfg)
	merchantController := controllers.NewMerchantController(db)
	serviceProviderController := controllers.NewServiceProviderController(db)
	creatorController := controllers.NewCreatorController(db)
	campaignController := controllers.NewCampaignController(db)
	taskController := controllers.NewTaskController(db)
	creditController := controllers.NewCreditController(db)
	withdrawalController := controllers.NewWithdrawalController(db)
	taskInvitationController := controllers.NewTaskInvitationController(db)
	rechargeOrderController := controllers.NewRechargeOrderController(db)

	// 新增：财务相关控制器
	withdrawalEnhancedController := controllers.NewWithdrawalEnhancedController(db, withdrawalEnhancedService, auditService)
	cashAccountController := controllers.NewCashAccountController(cashAccountService, auditService)
	systemAccountController := controllers.NewSystemAccountController(systemAccountService, auditService)
	financialAuditController := controllers.NewFinancialAuditController(auditService)

	// API路由组
	v1 := r.Group("/api/v1")
	{
		// 认证路由（无需认证）
		auth := v1.Group("/auth")
		{
			// 发起微信登录
			auth.GET("/wechat/login", authController.WeChatLoginRedirect)
		}

		// 用户路由（需要认证）
		user := v1.Group("/user")
		user.Use(middlewares.AuthCenterMiddleware(cfg, db))
		{
			// 获取当前用户信息
			user.GET("/me", authController.GetCurrentUser)
		}

		// 用户管理路由（需要认证+超级管理员权限）
		users := v1.Group("/users")
		users.Use(middlewares.AuthCenterMiddleware(cfg, db))
		{
			// 获取用户列表（仅超级管理员）
			users.GET("", authController.GetUsers)
		}

		// 邀请码路由（需要认证）
		invitations := v1.Group("/invitations")
		invitations.Use(middlewares.AuthCenterMiddleware(cfg, db))
		{
			// 获取我的固定邀请码列表（人邀请人）
			invitations.GET("/fixed-codes", invitationController.GetMyFixedInvitationCodes)
			// 使用邀请码（必须写在 /:code 之前，否则 "use" 会被当作 code）
			invitations.POST("/use", invitationController.UseInvitationCode)
			// 获取邀请码列表（旧版兼容）
			invitations.GET("", invitationController.ListInvitationCodes)
			// 获取我的邀请列表
			invitations.GET("/my", invitationController.GetMyInvitations)
			// 获取邀请码详情
			invitations.GET("/:code", invitationController.GetInvitationCode)
			// 禁用邀请码
			invitations.POST("/:code/disable", invitationController.DisableInvitationCode)
		}

		// 需要认证的路由
		protected := v1.Group("")
		protected.Use(middlewares.AuthCenterMiddleware(cfg, db))
		{


			// 商家管理
			protected.POST("/merchants", merchantController.CreateMerchant)
			protected.GET("/merchants", merchantController.GetMerchants)
			protected.GET("/merchants/:id", merchantController.GetMerchant)
			protected.PUT("/merchants/:id", merchantController.UpdateMerchant)
			protected.DELETE("/merchants/:id", merchantController.DeleteMerchant)
			protected.GET("/merchant/me", merchantController.GetMyMerchant)
			protected.POST("/merchants/:id/staff", merchantController.AddMerchantStaff)
			protected.GET("/merchants/:id/staff", merchantController.GetMerchantStaff)
			protected.PUT("/merchants/:id/staff/:staff_id/permissions", merchantController.UpdateMerchantStaffPermission)
			protected.DELETE("/merchants/:id/staff/:staff_id", merchantController.DeleteMerchantStaff)
			protected.GET("/merchants/permissions", merchantController.GetPermissions)

			// 服务商管理
			protected.POST("/service-providers", serviceProviderController.CreateServiceProvider)
			protected.GET("/service-providers", serviceProviderController.GetServiceProviders)
			protected.GET("/service-providers/:id", serviceProviderController.GetServiceProvider)
			protected.PUT("/service-providers/:id", serviceProviderController.UpdateServiceProvider)
			protected.DELETE("/service-providers/:id", serviceProviderController.DeleteServiceProvider)
			protected.GET("/service-provider/me", serviceProviderController.GetMyServiceProvider)
			protected.POST("/service-providers/:id/staff", serviceProviderController.AddServiceProviderStaff)
			protected.GET("/service-providers/:id/staff", serviceProviderController.GetServiceProviderStaff)
			protected.PUT("/service-providers/:id/staff/:staff_id/permissions", serviceProviderController.UpdateServiceProviderStaffPermission)
			protected.DELETE("/service-providers/:id/staff/:staff_id", serviceProviderController.DeleteServiceProviderStaff)
			protected.GET("/service-providers/permissions", serviceProviderController.GetPermissions)

			// 达人管理
			protected.GET("/creators", creatorController.GetCreators)
			protected.GET("/creators/stats/level", creatorController.GetCreatorLevelStats)
			protected.GET("/creators/:id", creatorController.GetCreator)
			protected.PUT("/creators/:id", creatorController.UpdateCreator)
			protected.GET("/creators/:id/inviter", creatorController.GetCreatorInviterRelationship)
			protected.POST("/creators/:id/break-relationship", creatorController.BreakInviterRelationship)
			protected.GET("/creator/me", creatorController.GetMyCreatorProfile)
			protected.PUT("/creator/me", creatorController.UpdateMyCreatorProfile)

			// 营销活动管理
			protected.POST("/campaigns", campaignController.CreateCampaign)
			protected.GET("/campaigns", campaignController.GetCampaigns)
			protected.GET("/campaigns/:id", campaignController.GetCampaign)
			protected.POST("/campaigns/:id/approve", campaignController.ApproveCampaign)
			protected.PUT("/campaigns/:id", campaignController.UpdateCampaign)
			protected.DELETE("/campaigns/:id", campaignController.DeleteCampaign)
			protected.GET("/campaigns/my", campaignController.GetMyCampaigns)

			// 任务邀请管理
			protected.POST("/task-invitations/generate", taskInvitationController.GenerateInvitationCode)
			protected.POST("/task-invitations/use", taskInvitationController.UseTaskInvitationCode)
			protected.GET("/task-invitations/validate/:code", taskInvitationController.ValidateInvitationCode)

			// 任务管理
			protected.GET("/tasks", taskController.GetTasks)
			protected.GET("/tasks/hall", taskController.GetTaskHall)
			protected.GET("/tasks/my", taskController.GetMyTasks)
			protected.GET("/tasks/pending-review", taskController.GetTasksForReview)
			protected.GET("/tasks/:id", taskController.GetTask)
			protected.POST("/tasks/:id/accept", taskController.AcceptTask)
			protected.POST("/tasks/:id/submit", taskController.SubmitTask)
			protected.POST("/tasks/:id/audit", taskController.AuditTask)

			// 积分管理
			protected.GET("/credit/accounts", creditController.GetUserAccounts)
			protected.GET("/credit/balance", creditController.GetAccountBalance)
			protected.GET("/credit/transactions", creditController.GetTransactions)
			protected.POST("/credit/recharge", creditController.Recharge)

			// 充值订单管理（线下充值流程）
			protected.POST("/recharge-orders", rechargeOrderController.CreateRechargeOrder)
			protected.GET("/recharge-orders", rechargeOrderController.GetRechargeOrders)
			protected.POST("/recharge-orders/:id/audit", rechargeOrderController.AuditRechargeOrder)

			// 提现管理
			protected.POST("/withdrawals", withdrawalController.CreateWithdrawal)
			protected.GET("/withdrawals", withdrawalController.GetWithdrawals)
			protected.GET("/withdrawals/:id", withdrawalController.GetWithdrawal)
			protected.POST("/withdrawals/:id/audit", withdrawalController.AuditWithdrawal)
			protected.POST("/withdrawals/:id/process", withdrawalController.ProcessWithdrawal)

			// 新增：增强提现管理（带冻结机制）
			protected.POST("/withdrawals/enhanced", withdrawalEnhancedController.CreateWithdrawalRequest)
			protected.GET("/withdrawals/enhanced", withdrawalEnhancedController.GetWithdrawalRequests)
			protected.GET("/withdrawals/enhanced/:id", withdrawalEnhancedController.GetWithdrawalRequest)
			protected.POST("/withdrawals/enhanced/:id/approve", withdrawalEnhancedController.ApproveWithdrawalRequest)
			protected.POST("/withdrawals/enhanced/:id/reject", withdrawalEnhancedController.RejectWithdrawalRequest)

			// 新增：现金账户管理
			protected.GET("/cash-accounts", cashAccountController.GetCashAccounts)
			protected.POST("/cash-accounts", cashAccountController.CreateCashAccount)
			protected.POST("/cash-accounts/:id/balance", cashAccountController.UpdateCashAccountBalance)

			// 新增：系统账户管理
			protected.GET("/system-accounts", systemAccountController.GetSystemAccounts)
			protected.GET("/system-accounts/summary", systemAccountController.GetFinancialSummary)

			// 新增：财务审计日志
			protected.GET("/financial-audit-logs", financialAuditController.GetAuditLogs)
		}
	}
}
