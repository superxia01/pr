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
		}
	}
}
