package controllers

import (
	"net/http"
	"pr-business/models"
	"pr-business/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CampaignController struct {
	db *gorm.DB
}

func NewCampaignController(db *gorm.DB) *CampaignController {
	return &CampaignController{db: db}
}

// CreateCampaignRequest 创建营销活动请求
type CreateCampaignRequest struct {
	Title              string    `json:"title" binding:"required,min=1,max=100"`
	Requirements       string    `json:"requirements" binding:"required"`
	Platforms          string    `json:"platforms" binding:"required"` // JSONB string
	TaskAmount         int       `json:"taskAmount" binding:"required,min=0"`
	CampaignAmount     int       `json:"campaignAmount" binding:"required,min=0"`
	CreatorAmount      *int      `json:"creatorAmount" binding:"omitempty,min=0"`
	StaffReferralAmount *int     `json:"staffReferralAmount" binding:"omitempty,min=0"`
	ProviderAmount     *int      `json:"providerAmount" binding:"omitempty,min=0"`
	Quota              int       `json:"quota" binding:"required,min=1"`
	TaskDeadline       time.Time `json:"taskDeadline" binding:"required"`
	SubmissionDeadline time.Time `json:"submissionDeadline" binding:"required"`
}

// UpdateCampaignRequest 更新营销活动请求
type UpdateCampaignRequest struct {
	Title              *string    `json:"title" binding:"omitempty,min=1,max=100"`
	Requirements       *string    `json:"requirements"`
	Platforms          *string    `json:"platforms"`
	TaskAmount         *int       `json:"taskAmount" binding:"omitempty,min=0"`
	CampaignAmount     *int       `json:"campaignAmount" binding:"omitempty,min=0"`
	CreatorAmount      *int       `json:"creatorAmount" binding:"omitempty,min=0"`
	StaffReferralAmount *int      `json:"staffReferralAmount" binding:"omitempty,min=0"`
	ProviderAmount     *int       `json:"providerAmount" binding:"omitempty,min=0"`
	Quota              *int       `json:"quota" binding:"omitempty,min=1"`
	TaskDeadline       *time.Time `json:"taskDeadline"`
	SubmissionDeadline *time.Time `json:"submissionDeadline"`
	Status             *models.CampaignStatus `json:"status" binding:"omitempty,oneof=DRAFT OPEN CLOSED"`
}

// CreateCampaign 创建营销活动
// @Summary 创建营销活动
// @Description 商家管理员创建营销活动
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

	// 权限检查：只有商家管理员可以创建营销活动
	if !utils.HasRole(user, "MERCHANT_ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有商家管理员可以创建营销活动"})
		return
	}

	// 获取商家信息
	var merchant models.Merchant
	if err := ctrl.db.Where("admin_id = ?", user.ID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是商家管理员"})
		return
	}

	// 创建营销活动
	campaign := models.Campaign{
		MerchantID:          merchant.ID,
		ProviderID:          &merchant.ProviderID,
		Title:               req.Title,
		Requirements:        req.Requirements,
		Platforms:           req.Platforms,
		TaskAmount:          req.TaskAmount,
		CampaignAmount:      req.CampaignAmount,
		CreatorAmount:       req.CreatorAmount,
		StaffReferralAmount: req.StaffReferralAmount,
		ProviderAmount:      req.ProviderAmount,
		Quota:               req.Quota,
		TaskDeadline:        req.TaskDeadline,
		SubmissionDeadline:  req.SubmissionDeadline,
		Status:              models.CampaignStatusDraft,
	}

	if err := ctrl.db.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建营销活动失败"})
		return
	}

	// 预创建任务名额
	for i := 1; i <= req.Quota; i++ {
		task := models.Task{
			CampaignID:     campaign.ID,
			TaskSlotNumber: i,
			Status:         models.TaskStatusOpen,
		}
		ctrl.db.Create(&task)
	}

	c.JSON(http.StatusOK, campaign)
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
	if utils.HasRole(user, "MERCHANT_ADMIN") && !utils.HasRole(user, "SUPER_ADMIN") {
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

// UpdateCampaign 更新营销活动
// @Summary 更新营销活动
// @Description 只有草稿状态的营销活动可以更新
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param id path string true "营销活动ID"
// @Param request body UpdateCampaignRequest true "更新营销活动请求"
// @Success 200 {object} models.Campaign
// @Router /api/v1/campaigns/{id} [put]
func (ctrl *CampaignController) UpdateCampaign(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var campaign models.Campaign
	if err := ctrl.db.Where("id = ?", id).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "营销活动不存在"})
		return
	}

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查：只有营销活动所属商家的管理员可以更新
	hasPermission := utils.HasRole(user, "SUPER_ADMIN") ||
		(utils.HasRole(user, "MERCHANT_ADMIN") && user.ID == campaign.MerchantID.String())

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新此营销活动"})
		return
	}

	// 只有草稿状态可以修改基本信息
	if campaign.Status != models.CampaignStatusDraft && req.Status == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有草稿状态的营销活动可以修改"})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Requirements != nil {
		updates["requirements"] = *req.Requirements
	}
	if req.Platforms != nil {
		updates["platforms"] = *req.Platforms
	}
	if req.TaskAmount != nil {
		updates["task_amount"] = *req.TaskAmount
	}
	if req.CampaignAmount != nil {
		updates["campaign_amount"] = *req.CampaignAmount
	}
	if req.CreatorAmount != nil {
		updates["creator_amount"] = *req.CreatorAmount
	}
	if req.StaffReferralAmount != nil {
		updates["staff_referral_amount"] = *req.StaffReferralAmount
	}
	if req.ProviderAmount != nil {
		updates["provider_amount"] = *req.ProviderAmount
	}
	if req.Quota != nil {
		updates["quota"] = *req.Quota
	}
	if req.TaskDeadline != nil {
		updates["task_deadline"] = *req.TaskDeadline
	}
	if req.SubmissionDeadline != nil {
		updates["submission_deadline"] = *req.SubmissionDeadline
	}
	if req.Status != nil {
		// 状态流转检查
		if campaign.Status == models.CampaignStatusDraft && *req.Status == models.CampaignStatusOpen {
			// 草稿 -> 开放：允许
			updates["status"] = *req.Status
		} else if campaign.Status == models.CampaignStatusOpen && *req.Status == models.CampaignStatusClosed {
			// 开放 -> 关闭：允许
			updates["status"] = *req.Status
		} else if campaign.Status == *req.Status {
			// 状态不变：允许
			updates["status"] = *req.Status
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态流转"})
			return
		}
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有可更新的字段"})
		return
	}

	if err := ctrl.db.Model(&campaign).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新营销活动失败"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// DeleteCampaign 删除营销活动（软删除）
// @Summary 删除营销活动
// @Description 只有草稿状态的营销活动可以删除
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Param id path string true "营销活动ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/campaigns/{id} [delete]
func (ctrl *CampaignController) DeleteCampaign(c *gin.Context) {
	id := c.Param("id")
	var campaign models.Campaign

	if err := ctrl.db.Where("id = ?", id).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "营销活动不存在"})
		return
	}

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查
	hasPermission := utils.HasRole(user, "SUPER_ADMIN") ||
		(utils.HasRole(user, "MERCHANT_ADMIN") && user.ID == campaign.MerchantID.String())

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此营销活动"})
		return
	}

	// 只有草稿状态可以删除
	if campaign.Status != models.CampaignStatusDraft {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有草稿状态的营销活动可以删除"})
		return
	}

	// 删除关联的任务名额
	ctrl.db.Where("campaign_id = ?", id).Delete(&models.Task{})

	// 删除营销活动
	if err := ctrl.db.Delete(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除营销活动失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetMyCampaigns 获取当前用户商家的营销活动
// @Summary 获取当前用户商家的营销活动
// @Description 商家管理员获取自己商家的营销活动列表
// @Tags 营销活动管理
// @Accept json
// @Produce json
// @Success 200 {array} models.Campaign
// @Router /api/v1/campaigns/my [get]
func (ctrl *CampaignController) GetMyCampaigns(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 获取商家信息
	var merchant models.Merchant
	if err := ctrl.db.Where("admin_id = ?", user.ID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "您不是商家管理员"})
		return
	}

	var campaigns []models.Campaign
	if err := ctrl.db.Where("merchant_id = ?", merchant.ID).Preload("Provider").Order("created_at DESC").Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取营销活动列表失败"})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}
