package routes

import (
	"pr-business/config"
	"pr-business/controllers"
	"pr-business/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, redisClient interface{}, cfg *config.Config) {
	// 初始化controllers
	authController := controllers.NewAuthController(db, cfg)
	invitationController := controllers.NewInvitationController(db, cfg)
	merchantController := controllers.NewMerchantController(db)
	serviceProviderController := controllers.NewServiceProviderController(db)
	creatorController := controllers.NewCreatorController(db)
	campaignController := controllers.NewCampaignController(db)
	taskController := controllers.NewTaskController(db)
	creditController := controllers.NewCreditController(db)

	// API路由组
	v1 := r.Group("/api/v1")
	{
		// 认证路由（无需JWT验证）
		auth := v1.Group("/auth")
		{
			auth.POST("/wechat", authController.WeChatLogin)
			auth.POST("/password", authController.PasswordLogin)
			auth.POST("/refresh", authController.RefreshToken)
		}

		// 邀请码路由（部分无需认证）
		invitations := v1.Group("/invitations")
		{
			// 使用邀请码（需要认证）
			invitations.POST("/use", middlewares.AuthMiddleware(cfg.JWTSecret), invitationController.UseInvitationCode)
		}

		// 需要认证的路由
		protected := v1.Group("")
		protected.Use(middlewares.AuthMiddleware(cfg.JWTSecret))
		{
			// 用户相关
			protected.GET("/user/me", authController.GetCurrentUser)
			protected.POST("/user/switch-role", authController.SwitchRole)

			// 邀请码管理
			protected.POST("/invitations", invitationController.CreateInvitationCode)
			protected.GET("/invitations", invitationController.ListInvitationCodes)
			protected.GET("/invitations/:code", invitationController.GetInvitationCode)
			protected.POST("/invitations/:code/disable", invitationController.DisableInvitationCode)

			// 商家管理
			protected.POST("/merchants", merchantController.CreateMerchant)
			protected.GET("/merchants", merchantController.GetMerchants)
			protected.GET("/merchants/:id", merchantController.GetMerchant)
			protected.PUT("/merchants/:id", merchantController.UpdateMerchant)
			protected.DELETE("/merchants/:id", merchantController.DeleteMerchant)
			protected.GET("/merchant/me", merchantController.GetMyMerchant)
			protected.POST("/merchants/:id/staff", merchantController.AddMerchantStaff)
			protected.GET("/merchants/:id/staff", merchantController.GetMerchantStaff)
			protected.PUT("/merchants/:merchant_id/staff/:staff_id/permissions", merchantController.UpdateMerchantStaffPermission)
			protected.DELETE("/merchants/:merchant_id/staff/:staff_id", merchantController.DeleteMerchantStaff)

			// 服务商管理
			protected.POST("/service-providers", serviceProviderController.CreateServiceProvider)
			protected.GET("/service-providers", serviceProviderController.GetServiceProviders)
			protected.GET("/service-providers/:id", serviceProviderController.GetServiceProvider)
			protected.PUT("/service-providers/:id", serviceProviderController.UpdateServiceProvider)
			protected.DELETE("/service-providers/:id", serviceProviderController.DeleteServiceProvider)
			protected.GET("/service-provider/me", serviceProviderController.GetMyServiceProvider)
			protected.POST("/service-providers/:id/staff", serviceProviderController.AddServiceProviderStaff)
			protected.GET("/service-providers/:id/staff", serviceProviderController.GetServiceProviderStaff)
			protected.PUT("/service-providers/:provider_id/staff/:staff_id/permissions", serviceProviderController.UpdateServiceProviderStaffPermission)
			protected.DELETE("/service-providers/:provider_id/staff/:staff_id", serviceProviderController.DeleteServiceProviderStaff)

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
			protected.GET("/campaigns/my", campaignController.GetMyCampaigns)
			protected.GET("/campaigns/:id", campaignController.GetCampaign)
			protected.PUT("/campaigns/:id", campaignController.UpdateCampaign)
			protected.DELETE("/campaigns/:id", campaignController.DeleteCampaign)

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
			protected.GET("/credit/balance", creditController.GetAccountBalance)
			protected.GET("/credit/transactions", creditController.GetTransactions)
			protected.POST("/credit/recharge", creditController.Recharge)
		}
	}
}
