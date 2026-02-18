package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"pr-business/models"
	"pr-business/services"
	"pr-business/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CampaignController struct {
	db                *gorm.DB
	settlementService  *services.SettlementService
}

func NewCampaignController(db *gorm.DB) *CampaignController {
	// 初始化依赖服务
	settlementService := services.NewSettlementService(
		db,
		nil, // permissionService - 稍后实现
		nil, // validatorService
		nil, // cashAccountService
	)

	return &CampaignController{
		db:                db,
		settlementService: settlementService,
	}
}

// CreateCampaignRequest 创建营销活动请求
type CreateCampaignRequest struct {
	MerchantID         *string   `json:"merchantId" binding:"omitempty"` // 服务商创建时必填
	Title              string    `json:"title" binding:"required,min=1,max=100"`
	Requirements       string    `json:"requirements" binding:"required"`
	Platforms          string    `json:"platforms" binding:"required"` // JSONB string
	TaskAmount         int       `json:"taskAmount" binding:"required,min=0"`
	CreatorAmount      *int      `json:"creatorAmount" binding:"omitempty,min=0"`
	StaffReferralAmount *int     `json:"staffReferralAmount" binding:"omitempty,min=0"`
	ProviderAmount     *int      `json:"providerAmount" binding:"omitempty,min=0"`
	Quota              int       `json:"quota" binding:"required,min=1"`
	TaskDeadline       time.Time `json:"taskDeadline" binding:"required"`
	SubmissionDeadline time.Time `json:"submissionDeadline" binding:"required"`
}

// CreateCampaign 创建营销活动
// @Summary 创建营销活动
// @Description 商家管理员或服务商管理员创建营销活动
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param request body CreateCampaignRequest true "创建营销活动请求"
// @Success 200 {object} models.Campaign
// @Router /api/v1/campaigns [post]
func (ctrl *CampaignController) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
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

	var merchantID uuid.UUID
	var providerID *uuid.UUID
	var creatorType string
	var status models.CampaignStatus

	// 判断创建者类型并验证权限
	if utils.IsMerchantAdmin(user) {
		// 商家管理员创建活动
		creatorType = "MERCHANT_ADMIN"

		// 获取商家信息
		var merchant models.Merchant
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&merchant).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "您不是商家管理员"})
			return
		}
		merchantID = merchant.ID
		providerID = &merchant.ProviderID

		// 商家创建活动：状态为待审核，不设置佣金分配
		status = models.CampaignStatusPendingApproval

		// 商家不应该填写佣金分配（清空）
		req.CreatorAmount = nil
		req.StaffReferralAmount = nil
		req.ProviderAmount = nil

	} else if utils.IsServiceProviderAdmin(user) {
		// 服务商管理员创建活动
		creatorType = "SERVICE_PROVIDER_ADMIN"

		// 必须指定商家
		if req.MerchantID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "服务商创建活动必须指定商家"})
			return
		}

		// 获取服务商信息
		var provider models.ServiceProvider
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&provider).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "您不是服务商管理员"})
			return
		}
		providerID = &provider.ID

		// 解析商家ID
		parsedMerchantID, err := uuid.Parse(*req.MerchantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "商家ID格式错误"})
			return
		}
		merchantID = parsedMerchantID

		// 验证商家属于该服务商
		var merchant models.Merchant
		if err := ctrl.db.Where("id = ?", merchantID).First(&merchant).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "商家不存在"})
			return
		}

		if merchant.ProviderID != provider.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "该商家不属于您的服务商"})
			return
		}

		// 服务商创建活动：直接发布
		status = models.CampaignStatusOpen

	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有商家管理员或服务商管理员可以创建营销活动"})
		return
	}

	// 计算活动总金额
	campaignAmount := req.TaskAmount * req.Quota

	// 开始事务创建活动（如果直接发布，需要冻结积分）
	err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		// 创建营销活动
		campaign := models.Campaign{
			MerchantID:          merchantID,
			ProviderID:          providerID,
			CreatedBy:           user.AuthCenterUserID,
			CreatorType:         creatorType,
			Title:               req.Title,
			Requirements:        req.Requirements,
			Platforms:           req.Platforms,
			TaskAmount:          req.TaskAmount,
			CampaignAmount:      campaignAmount,
			CreatorAmount:       req.CreatorAmount,
			StaffReferralAmount: req.StaffReferralAmount,
			ProviderAmount:      req.ProviderAmount,
			Quota:               req.Quota,
			TaskDeadline:        req.TaskDeadline,
			SubmissionDeadline:  req.SubmissionDeadline,
			Status:              status,
		}

		if err := tx.Create(&campaign).Error; err != nil {
			return fmt.Errorf("创建营销活动失败: %w", err)
		}

		// 如果服务商直接创建活动（状态为 OPEN），需要冻结商家积分
		if status == models.CampaignStatusOpen {
			// 验证佣金分配总和
			totalCommission := 0
			if campaign.CreatorAmount != nil {
				totalCommission += *campaign.CreatorAmount
			}
			if campaign.StaffReferralAmount != nil {
				totalCommission += *campaign.StaffReferralAmount
			}
			if campaign.ProviderAmount != nil {
				totalCommission += *campaign.ProviderAmount
			}

			if totalCommission != campaign.TaskAmount {
				return errors.New("佣金分配总和必须等于任务总金额")
			}

			// 获取商家积分账户
			var creditAccount models.CreditAccount
			if err := tx.Where("owner_id = ? AND owner_type = ?", campaign.MerchantID, models.OwnerTypeOrgMerchant).First(&creditAccount).Error; err != nil {
				return fmt.Errorf("获取商家积分账户失败: %w", err)
			}

			// 验证可用余额是否足够
			if creditAccount.Balance < campaignAmount {
				return errors.New("商家可用积分不足")
			}

			// 冻结积分：从可用余额转移到冻结余额
			creditAccount.Balance -= campaignAmount
			creditAccount.FrozenBalance += campaignAmount

			// 记录冻结流水
			freezeTransaction := models.CreditTransaction{
				AccountID:      creditAccount.ID,
				Type:           "CAMPAIGN_FREEZE",
				Amount:         campaignAmount,
				BalanceBefore:  creditAccount.Balance + campaignAmount,
				BalanceAfter:   creditAccount.Balance,
				RelatedCampaignID: &campaign.ID,
				Description:    fmt.Sprintf("发布活动冻结积分：%s", campaign.Title),
			}
			if err := tx.Create(&freezeTransaction).Error; err != nil {
				return fmt.Errorf("记录冻结流水失败: %w", err)
			}

			// 保存积分账户更新
			if err := tx.Save(&creditAccount).Error; err != nil {
				return fmt.Errorf("更新积分账户失败: %w", err)
			}
		}

		// 预创建任务名额
		for i := 1; i <= req.Quota; i++ {
			task := models.Task{
				CampaignID:     campaign.ID,
				TaskSlotNumber: i,
				Status:         models.TaskStatusOpen,
			}
			if err := tx.Create(&task).Error; err != nil {
				return fmt.Errorf("创建任务名额失败: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新加载活动数据
	var campaign models.Campaign
	if err := ctrl.db.Preload("Merchant").Preload("Provider").Where("merchant_id = ? AND created_by = ?", merchantID, user.AuthCenterUserID).Order("created_at DESC").First(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载活动数据失败"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// UpdateCampaignRequest 更新营销活动请求
type UpdateCampaignRequest struct {
	Title              *string   `json:"title" binding:"omitempty,min=1,max=100"`
	Requirements       *string   `json:"requirements" binding:"omitempty"`
	Platforms          *string   `json:"platforms" binding:"omitempty"`
	TaskDeadline       *time.Time `json:"taskDeadline" binding:"omitempty"`
	SubmissionDeadline *time.Time `json:"submissionDeadline" binding:"omitempty"`
	Status             *string   `json:"status" binding:"omitempty,oneof=DRAFT PENDING_APPROVAL OPEN CLOSED"`
}

// UpdateCampaign 更新营销活动
// @Summary 更新营销活动
// @Description 更新活动信息或关闭活动。只有DRAFT状态的活动可以修改基本信息。关闭活动时会解冻未完成任务的积分。
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param id path string true "营销活动ID"
// @Param request body UpdateCampaignRequest true "更新营销活动请求"
// @Success 200 {object} models.Campaign
// @Router /api/v1/campaigns/{id} [put]
func (ctrl *CampaignController) UpdateCampaign(c *gin.Context) {
	id := c.Param("id")

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	var req UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找活动
	var campaign models.Campaign
	if err := ctrl.db.Where("id = ?", id).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
		return
	}

	// 权限检查：只有创建者、服务商管理员、超级管理员可以更新
	if !utils.IsServiceProviderAdmin(user) && !utils.IsSuperAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新活动"})
		return
	}

	// 如果是服务商管理员，验证活动属于该服务商
	if utils.IsServiceProviderAdmin(user) {
		var provider models.ServiceProvider
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&provider).Error; err == nil {
			if campaign.ProviderID == nil || *campaign.ProviderID != provider.ID {
				c.JSON(http.StatusForbidden, gin.H{"error": "该活动不属于您的服务商"})
				return
			}
		}
	}

	// 处理状态变更
	if req.Status != nil {
		newStatus := models.CampaignStatus(*req.Status)

		// 验证状态转换的合法性
		if campaign.Status == models.CampaignStatusClosed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "活动已关闭，无法再次操作"})
			return
		}

		// 关闭活动
		if newStatus == models.CampaignStatusClosed {
			if campaign.Status != models.CampaignStatusOpen {
				c.JSON(http.StatusBadRequest, gin.H{"error": "只能关闭开放中的活动"})
				return
			}

			// 调用积分解冻服务
			if err := ctrl.settlementService.SettleCampaignAfterClose(&campaign); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			campaign.Status = models.CampaignStatusClosed
		}

		// 其他状态转换
		if newStatus == models.CampaignStatusOpen {
			if campaign.Status != models.CampaignStatusPendingApproval {
				c.JSON(http.StatusBadRequest, gin.H{"error": "只能审核待审核状态的活动"})
				return
			}
			// 应该使用 ApproveCampaign 接口，这里不允许直接改为 OPEN
			c.JSON(http.StatusBadRequest, gin.H{"error": "请使用审核接口来发布活动"})
			return
		}
	}

	// 更新基本信息（仅允许 DRAFT 或 PENDING_APPROVAL 状态）
	if campaign.Status == models.CampaignStatusDraft || campaign.Status == models.CampaignStatusPendingApproval {
		if req.Title != nil {
			campaign.Title = *req.Title
		}
		if req.Requirements != nil {
			campaign.Requirements = *req.Requirements
		}
		if req.Platforms != nil {
			campaign.Platforms = *req.Platforms
		}
		if req.TaskDeadline != nil {
			campaign.TaskDeadline = *req.TaskDeadline
		}
		if req.SubmissionDeadline != nil {
			campaign.SubmissionDeadline = *req.SubmissionDeadline
		}
	} else if req.Status == nil {
		// 活动已发布，不允许修改基本信息
		c.JSON(http.StatusBadRequest, gin.H{"error": "活动已发布，不允许修改基本信息"})
		return
	}

	if err := ctrl.db.Save(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新活动失败"})
		return
	}

	// 重新加载活动数据
	ctrl.db.Preload("Merchant").Preload("Provider").First(&campaign, campaign.ID)

	c.JSON(http.StatusOK, campaign)
}

// ApproveCampaignRequest 审核营销活动请求
type ApproveCampaignRequest struct {
	CreatorAmount       *int `json:"creatorAmount" binding:"required,min=0"`        // 达人佣金
	StaffReferralAmount *int `json:"staffReferralAmount" binding:"required,min=0"` // 员工推荐佣金
	ProviderAmount      *int `json:"providerAmount" binding:"required,min=0"`      // 服务商佣金
}

// ApproveCampaign 审核并发布营销活动
// @Summary 审核并发布营销活动
// @Description 服务商审核商家提交的活动，填写佣金分配并发布
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param id path string true "营销活动ID"
// @Param request body ApproveCampaignRequest true "审核营销活动请求"
// @Success 200 {object} models.Campaign
// @Router /api/v1/campaigns/{id}/approve [post]
func (ctrl *CampaignController) ApproveCampaign(c *gin.Context) {
	campaignID := c.Param("id")

	var req ApproveCampaignRequest
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

	// 只有服务商管理员可以审核
	if !utils.IsServiceProviderAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有服务商管理员可以审核活动"})
		return
	}

	// 获取服务商信息
	var provider models.ServiceProvider
	if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&provider).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是服务商管理员"})
		return
	}

	// 查找活动
	var campaign models.Campaign
	if err := ctrl.db.Where("id = ?", campaignID).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
		return
	}

	// 验证活动属于该服务商
	if campaign.ProviderID == nil || *campaign.ProviderID != provider.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "该活动不属于您的服务商"})
		return
	}

	// 验证活动状态
	if campaign.Status != models.CampaignStatusPendingApproval {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能审核待审核状态的活动"})
		return
	}

	// 验证佣金分配总和等于任务总金额
	totalCommission := 0
	if req.CreatorAmount != nil {
		totalCommission += *req.CreatorAmount
	}
	if req.StaffReferralAmount != nil {
		totalCommission += *req.StaffReferralAmount
	}
	if req.ProviderAmount != nil {
		totalCommission += *req.ProviderAmount
	}

	if totalCommission != campaign.TaskAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "佣金分配总和必须等于任务总金额"})
		return
	}

	// 冻结商家积分（从可用余额转移到冻结余额）
	campaignAmount := campaign.TaskAmount * campaign.Quota

	// 开始冻结积分事务
	if err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		// 获取商家积分账户
		var creditAccount models.CreditAccount
		if err := tx.Where("owner_id = ? AND owner_type = ?", campaign.MerchantID, models.OwnerTypeOrgMerchant).First(&creditAccount).Error; err != nil {
			return fmt.Errorf("获取商家积分账户失败: %w", err)
		}

		// 验证可用余额是否足够
		if creditAccount.Balance < campaignAmount {
			return errors.New("商家可用积分不足")
		}

		// 冻结积分：从可用余额转移到冻结余额
		creditAccount.Balance -= campaignAmount
		creditAccount.FrozenBalance += campaignAmount

		// 记录冻结流水
		freezeTransaction := models.CreditTransaction{
			AccountID:      creditAccount.ID,
			Type:           "CAMPAIGN_FREEZE",
			Amount:         campaignAmount,
			BalanceBefore:  creditAccount.Balance + campaignAmount,
			BalanceAfter:   creditAccount.Balance,
			Description:    fmt.Sprintf("发布活动冻结积分：%s", campaign.Title),
		}
		if err := tx.Create(&freezeTransaction).Error; err != nil {
			return fmt.Errorf("记录冻结流水失败: %w", err)
		}

		// 更新活动状态
	updates := map[string]interface{}{
			"creator_amount":        req.CreatorAmount,
			"staff_referral_amount": req.StaffReferralAmount,
			"provider_amount":       req.ProviderAmount,
			"status":                models.CampaignStatusOpen,
		}

		if err := tx.Model(&campaign).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新活动状态失败: %w", err)
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新加载活动数据
	ctrl.db.Preload("Merchant").Preload("Provider").First(&campaign, campaign.ID)

	c.JSON(http.StatusOK, campaign)
}

// DeleteCampaign 删除营销活动
// @Summary 删除营销活动
// @Description 只能删除DRAFT或PENDING_APPROVAL状态的活动。已开放的活动必须先关闭才能删除。
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param id path string true "营销活动ID"
// @Success 200 {string} success message
// @Router /api/v1/campaigns/{id} [delete]
func (ctrl *CampaignController) DeleteCampaign(c *gin.Context) {
	id := c.Param("id")

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 查找活动
	var campaign models.Campaign
	if err := ctrl.db.Where("id = ?", id).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
		return
	}

	// 权限检查：只有服务商管理员、超级管理员可以删除
	if !utils.IsServiceProviderAdmin(user) && !utils.IsSuperAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除活动"})
		return
	}

	// 如果是服务商管理员，验证活动属于该服务商
	if utils.IsServiceProviderAdmin(user) {
		var provider models.ServiceProvider
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&provider).Error; err == nil {
			if campaign.ProviderID == nil || *campaign.ProviderID != provider.ID {
				c.JSON(http.StatusForbidden, gin.H{"error": "该活动不属于您的服务商"})
				return
			}
		}
	}

	// 验证活动状态：只能删除 DRAFT 或 PENDING_APPROVAL 状态的活动
	if campaign.Status == models.CampaignStatusOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "活动已开放，请先关闭活动再删除"})
		return
	}
	if campaign.Status == models.CampaignStatusClosed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "活动已关闭，不允许删除"})
		return
	}

	// 开始事务删除活动
	err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的任务名额
		if err := tx.Where("campaign_id = ?", campaign.ID).Delete(&models.Task{}).Error; err != nil {
			return fmt.Errorf("删除任务名额失败: %w", err)
		}

		// 删除关联的活动邀请（如果有）
		if err := tx.Where("campaign_id = ?", campaign.ID).Delete(&models.CampaignInvitation{}).Error; err != nil {
			// 忽略错误，因为可能没有邀请码
		}

		// 软删除活动
		if err := tx.Delete(&campaign).Error; err != nil {
			return fmt.Errorf("删除活动失败: %w", err)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetCampaigns 获取营销活动列表
// @Summary 获取营销活动列表
// @Description 获取营销活动列表，支持过滤
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param status query string false "状态过滤"
// @Param merchant_id query string false "商家ID过滤"
// @Success 200 {array} models.Campaign
// @Router /api/v1/campaigns [get]
func (ctrl *CampaignController) GetCampaigns(c *gin.Context) {
	var campaigns []models.Campaign
	query := ctrl.db.Model(&models.Campaign{})

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限过滤
	if utils.IsMerchantAdmin(user) && !utils.IsSuperAdmin(user) {
		// 商家管理员只能看到自己的营销活动
		var merchant models.Merchant
		if err := ctrl.db.Where("admin_id = ?", user.ID).First(&merchant).Error; err == nil {
			query = query.Where("merchant_id = ?", merchant.ID)
		}
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 商家过滤
	if merchantID := c.Query("merchant_id"); merchantID != "" {
		query = query.Where("merchant_id = ?", merchantID)
	}

	if err := query.Preload("Merchant").Preload("Provider").Order("created_at DESC").Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取营销活动列表失败"})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

// GetCampaign 获取营销活动详情
// @Summary 获取营销活动详情
// @Description 根据ID获取营销活动详情
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param id path string true "营销活动ID"
// @Success 200 {object} models.Campaign
// @Router /api/v1/campaigns/{id} [get]
func (ctrl *CampaignController) GetCampaign(c *gin.Context) {
	id := c.Param("id")

	var campaign models.Campaign
	if err := ctrl.db.Where("id = ?", id).Preload("Merchant").Preload("Provider").Preload("Tasks").First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "营销活动不存在"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// GetMyCampaigns 获取我的营销活动列表
// @Summary 获取我的营销活动列表
// @Description 根据用户角色返回相关的营销活动列表
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param status query string false "状态过滤"
// @Success 200 {array} models.Campaign
// @Router /api/v1/campaigns/my [get]
func (ctrl *CampaignController) GetMyCampaigns(c *gin.Context) {
	var campaigns []models.Campaign
	query := ctrl.db.Model(&models.Campaign{})

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 根据角色返回不同的活动列表
	if utils.IsMerchantAdmin(user) && !utils.IsSuperAdmin(user) {
		// 商家管理员：返回自己创建的所有活动
		var merchant models.Merchant
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&merchant).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "您不是商家管理员"})
			return
		}
		query = query.Where("merchant_id = ?", merchant.ID)
	} else if utils.IsServiceProviderAdmin(user) && !utils.IsSuperAdmin(user) {
		// 服务商管理员：返回管理的所有活动
		var provider models.ServiceProvider
		if err := ctrl.db.Where("admin_id::text = ?", user.AuthCenterUserID).First(&provider).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "您不是服务商管理员"})
			return
		}
		query = query.Where("provider_id = ?", provider.ID)
	} else if utils.IsSuperAdmin(user) {
		// 超级管理员：返回所有活动
		// 不过滤
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看活动"})
		return
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Preload("Merchant").Preload("Provider").Order("created_at DESC").Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取营销活动列表失败"})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}
