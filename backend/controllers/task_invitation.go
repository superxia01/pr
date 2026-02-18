package controllers

import (
	"fmt"
	"net/http"
	"pr-business/constants"
	"pr-business/models"
	"pr-business/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskInvitationController struct {
	db *gorm.DB
}

func NewTaskInvitationController(db *gorm.DB) *TaskInvitationController {
	return &TaskInvitationController{db: db}
}

// GenerateInvitationCodeRequest 生成邀请码请求
type GenerateInvitationCodeRequest struct {
	CampaignID string `json:"campaignId" binding:"required"`
	MaxUses     int    `json:"maxUses"`
	ExpiresIn   int    `json:"expiresIn"` // 过期时间（秒）
}

// GenerateInvitationCodeResponse 生成邀请码响应
type GenerateInvitationCodeResponse struct {
	InvitationURL string `json:"invitationUrl"`
	Code           string `json:"code"`
	MaxUses        int    `json:"maxUses"`
	ExpiresAt      string `json:"expiresAt"`
}

// GenerateInvitationCode 生成任务邀请码
// POST /api/v1/task-invitations/generate
func (ctrl *TaskInvitationController) GenerateInvitationCode(c *gin.Context) {
	var req GenerateInvitationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查：只有商家和服务商管理员可以生成邀请码
	if !utils.HasRole(user, constants.RoleMerchantAdmin) && !utils.HasRole(user, constants.RoleServiceProviderAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有商家管理员和服务商管理员可以生成邀请码"})
		return
	}

	// 验证营销活动
	var campaign models.Campaign
	if err := ctrl.db.Where("id = ?", req.CampaignID).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "营销活动不存在"})
		return
	}

	// 检查活动状态
	if campaign.Status != "OPEN" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "营销活动未开放"})
		return
	}

	// 生成6位随机邀请码
	code := utils.GenerateInvitationCode()

	// 计算过期时间
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(req.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	// 确定生成者类型
	generatorType := "merchant_admin"
	if utils.HasRole(user, constants.RoleServiceProviderAdmin) {
		generatorType = "provider_admin"
	}

	// 创建邀请码记录
	invitationCode := models.TaskInvitationCode{
		Code:         code,
		CampaignID:  campaign.ID,
		GeneratorID:  user.ID,
		GeneratorType: generatorType,
		MaxUses:      req.MaxUses,
		ExpiresAt:   expiresAt,
	}

	if err := ctrl.db.Create(&invitationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建邀请码失败"})
		return
	}

	// 生成邀请链接
	invitationURL := fmt.Sprintf("https://pr.crazyaigc.com/invite/%s", code)

	c.JSON(http.StatusOK, GenerateInvitationCodeResponse{
		InvitationURL: invitationURL,
		Code:           code,
		MaxUses:        req.MaxUses,
		ExpiresAt:      expiresAt.Format(time.RFC3339),
	})
}

// ValidateInvitationCode 验证邀请码
// GET /api/v1/task-invitations/validate/:code
func (ctrl *TaskInvitationController) ValidateInvitationCode(c *gin.Context) {
	code := c.Param("code")

	var invitationCode models.TaskInvitationCode
	if err := ctrl.db.Where("code = ?", code).Preload("Campaign").First(&invitationCode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "邀请码不存在"})
		return
	}

	// 检查邀请码是否有效
	if invitationCode.ExpiresAt != nil && time.Now().After(*invitationCode.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "邀请码已过期"})
		return
	}

	if invitationCode.UseCount >= invitationCode.MaxUses {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邀请码使用次数已达上限"})
		return
	}

	// 检查活动状态
	if invitationCode.Campaign.Status != "OPEN" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "营销活动已结束"})
		return
	}

	// 获取活动的任务列表
	var tasks []models.Task
	if err := ctrl.db.Where("campaign_id = ? AND status = ?", invitationCode.CampaignID, models.TaskStatusOpen).
		Preload("Campaign.Merchant").
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务列表失败"})
		return
	}

	// 返回邀请码信息和任务列表
	c.JSON(http.StatusOK, gin.H{
		"invitationCode": invitationCode,
		"campaign":      invitationCode.Campaign,
		"tasks":         tasks,
		"totalTasks":    len(tasks),
	})
}

// UseTaskInvitationCodeRequest 使用营销活动邀请码请求
type UseTaskInvitationCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

// UseTaskInvitationCode 使用营销活动邀请码（登录后调用）
// 效果：1）被邀请人自动获得达人角色 2）与邀请人绑定（Creator.inviter） 3）与营销活动的商家/服务商产生绑定（通过 TaskInvitation + Campaign）
// POST /api/v1/task-invitations/use
func (ctrl *TaskInvitationController) UseTaskInvitationCode(c *gin.Context) {
	var req UseTaskInvitationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供邀请码"})
		return
	}

	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	var invitationCode models.TaskInvitationCode
	if err := ctrl.db.Where("code = ?", req.Code).Preload("Campaign").Preload("Campaign.Merchant").Preload("Campaign.Provider").First(&invitationCode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "邀请码不存在"})
		return
	}

	if invitationCode.ExpiresAt != nil && time.Now().After(*invitationCode.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "邀请码已过期"})
		return
	}
	if invitationCode.MaxUses > 0 && invitationCode.UseCount >= invitationCode.MaxUses {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邀请码使用次数已达上限"})
		return
	}
	if invitationCode.Campaign.Status != models.CampaignStatusOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "营销活动未开放"})
		return
	}

	// 已使用过该码则直接返回成功
	var existing models.TaskInvitation
	if ctrl.db.Where("invitation_code_id = ? AND creator_id = ?", invitationCode.ID, user.ID).First(&existing).Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":       "已使用过该邀请码",
			"campaign":      invitationCode.Campaign,
			"invitationCode": invitationCode,
		})
		return
	}

	err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		// 1. 若无达人角色则赋予
		if !utils.HasRole(user, constants.RoleCreator) {
			user.Roles = append(user.Roles, constants.RoleCreator)
			if err := tx.Save(user).Error; err != nil {
				return fmt.Errorf("更新用户角色失败: %w", err)
			}
		}

		// 2. 创建或更新达人档案，绑定邀请人（GeneratorID = 发码人）
		inviterType := taskInviterTypeFromGeneratorType(invitationCode.GeneratorType)
		var creator models.Creator
		if err := tx.Where("user_id = ? AND is_primary = ?", user.ID, true).First(&creator).Error; err != nil {
			creator = models.Creator{
				UserID:      user.ID,
				IsPrimary:   true,
				Level:       "UGC",
				InviterID:   &invitationCode.GeneratorID,
				InviterType: inviterType,
				Status:      "active",
			}
			if err := tx.Create(&creator).Error; err != nil {
				return fmt.Errorf("创建达人绑定失败: %w", err)
			}
		} else {
			updates := map[string]interface{}{
				"inviter_id":   invitationCode.GeneratorID,
				"inviter_type": inviterType,
			}
			if err := tx.Model(&creator).Updates(updates).Error; err != nil {
				return fmt.Errorf("更新达人邀请关系失败: %w", err)
			}
		}

		// 3. 记录与营销活动/商家/服务商的绑定（TaskInvitation 关联 invitation_code -> campaign -> merchant + provider）
		ti := models.TaskInvitation{
			InvitationCodeID: invitationCode.ID,
			CreatorID:        &user.ID,
			TaskID:           nil,
			Status:           "accepted",
		}
		if err := tx.Create(&ti).Error; err != nil {
			return fmt.Errorf("创建邀请绑定记录失败: %w", err)
		}

		// 4. 增加邀请码使用次数
		if err := tx.Model(&invitationCode).Update("use_count", gorm.Expr("use_count + 1")).Error; err != nil {
			return fmt.Errorf("更新使用次数失败: %w", err)
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "邀请码使用成功，已获得达人身份并绑定邀请人与活动",
		"campaign":       invitationCode.Campaign,
		"invitationCode": invitationCode,
	})
}

// taskInviterTypeFromGeneratorType 营销活动邀请码生成者类型 -> Creator.InviterType
func taskInviterTypeFromGeneratorType(gt string) string {
	switch strings.ToLower(gt) {
	case "provider_admin":
		return "PROVIDER_ADMIN"
	case "provider_staff", "merchant_staff":
		return "PROVIDER_STAFF"
	default:
		return "OTHER"
	}
}
