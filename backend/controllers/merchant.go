package controllers

import (
	"net/http"
	"pr-business/constants"
	"pr-business/models"
	"pr-business/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MerchantController struct {
	db *gorm.DB
}

func NewMerchantController(db *gorm.DB) *MerchantController {
	return &MerchantController{db: db}
}

// CreateMerchantRequest 创建商家请求
type CreateMerchantRequest struct {
	ProviderID  string `json:"providerId" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
	Industry    string `json:"industry"`
}

// UpdateMerchantRequest 更新商家请求
type UpdateMerchantRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description"`
	Industry    string `json:"industry"`
	Status      string `json:"status" binding:"omitempty,oneof=active suspended inactive"`
}

// AddMerchantStaffRequest 添加商家员工请求
type AddMerchantStaffRequest struct {
	UserID      string   `json:"userId" binding:"required"`
	Title       string   `json:"title" binding:"omitempty,max=50"`
	Permissions []string `json:"permissions" binding:"required"`
}

// UpdateMerchantStaffPermissionRequest 更新员工权限请求
type UpdateMerchantStaffPermissionRequest struct {
	PermissionCode string `json:"permissionCode" binding:"required"`
	Action         string `json:"action" binding:"required,oneof=grant revoke"`
}

// CreateMerchant 创建商家
// @Summary 创建商家
// @Description 只有超级管理员和服务商管理员可以创建商家
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param request body CreateMerchantRequest true "创建商家请求"
// @Success 200 {object} models.Merchant
// @Router /api/v1/merchants [post]
func (ctrl *MerchantController) CreateMerchant(c *gin.Context) {
	var req CreateMerchantRequest
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

	// 权限检查：只有超级管理员和服务商管理员可以创建商家
	if !utils.IsSuperAdmin(user) && !utils.IsServiceProviderAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限创建商家"})
		return
	}

	// 如果是服务商管理员，只能为自己所属的服务商创建商家
	if utils.IsServiceProviderAdmin(user) && !utils.IsSuperAdmin(user) {
		var serviceProvider models.ServiceProvider
		if err := ctrl.db.Where("admin_id = ?", user.ID).First(&serviceProvider).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "您不是服务商管理员"})
			return
		}
		if serviceProvider.ID.String() != req.ProviderID {
			c.JSON(http.StatusForbidden, gin.H{"error": "只能为自己所属的服务商创建商家"})
			return
		}
	}

	// 检查服务商是否存在
	var serviceProvider models.ServiceProvider
	if err := ctrl.db.Where("id = ?", req.ProviderID).First(&serviceProvider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务商不存在"})
		return
	}

	// 创建商家（admin_id 暂时为空，后续通过邀请码绑定管理员）
	merchant := models.Merchant{
		ProviderID:  serviceProvider.ID,
		UserID:      user.ID, // 记录创建者
		Name:        req.Name,
		Description: req.Description,
		Industry:    req.Industry,
		Status:      "active",
	}

	if err := ctrl.db.Create(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商家失败"})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

// GetMerchants 获取商家列表
// @Summary 获取商家列表
// @Description 获取商家列表，支持过滤
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param status query string false "状态过滤"
// @Param provider_id query string false "服务商ID过滤"
// @Success 200 {array} models.Merchant
// @Router /api/v1/merchants [get]
func (ctrl *MerchantController) GetMerchants(c *gin.Context) {
	var merchants []models.Merchant
	query := ctrl.db.Model(&models.Merchant{})

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限过滤
	if utils.IsServiceProviderAdmin(user) && !utils.IsSuperAdmin(user) {
		// 服务商管理员只能看到自己服务商下的商家
		var serviceProvider models.ServiceProvider
		if err := ctrl.db.Where("admin_id = ?", user.ID).First(&serviceProvider).Error; err == nil {
			query = query.Where("provider_id = ?", serviceProvider.ID)
		}
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 服务商过滤
	if providerID := c.Query("provider_id"); providerID != "" {
		query = query.Where("provider_id = ?", providerID)
	}

	if err := query.Preload("Admin").Preload("Provider").Preload("User").Find(&merchants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商家列表失败"})
		return
	}

	c.JSON(http.StatusOK, merchants)
}

// GetMerchant 获取商家详情
// @Summary 获取商家详情
// @Description 根据ID获取商家详情
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param id path string true "商家ID"
// @Success 200 {object} models.Merchant
// @Router /api/v1/merchants/{id} [get]
func (ctrl *MerchantController) GetMerchant(c *gin.Context) {
	id := c.Param("id")
	var merchant models.Merchant

	if err := ctrl.db.Where("id = ?", id).Preload("Admin").Preload("Provider").Preload("User").Preload("Staff").First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家不存在"})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

// UpdateMerchant 更新商家信息
// @Summary 更新商家信息
// @Description 只有超级管理员、服务商管理员和该商家的管理员可以更新
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param id path string true "商家ID"
// @Param request body UpdateMerchantRequest true "更新商家请求"
// @Success 200 {object} models.Merchant
// @Router /api/v1/merchants/{id} [put]
func (ctrl *MerchantController) UpdateMerchant(c *gin.Context) {
	id := c.Param("id")
	var req UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var merchant models.Merchant
	if err := ctrl.db.Where("id = ?", id).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家不存在"})
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
	hasPermission := utils.IsSuperAdmin(user) ||
		utils.IsServiceProviderAdmin(user) ||
		(utils.IsMerchantAdmin(user) && merchant.AdminID == user.ID)

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新商家信息"})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Industry != "" {
		updates["industry"] = req.Industry
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := ctrl.db.Model(&merchant).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新商家失败"})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

// DeleteMerchant 删除商家（软删除）
// @Summary 删除商家
// @Description 只有超级管理员可以删除商家
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param id path string true "商家ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{id} [delete]
func (ctrl *MerchantController) DeleteMerchant(c *gin.Context) {
	id := c.Param("id")
	var merchant models.Merchant

	if err := ctrl.db.Where("id = ?", id).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家不存在"})
		return
	}

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查：只有超级管理员可以删除
	if !utils.IsSuperAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除商家"})
		return
	}

	if err := ctrl.db.Delete(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除商家失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// AddMerchantStaff 添加商家员工
// @Summary 添加商家员工
// @Description 只有商家管理员可以添加员工
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param id path string true "商家ID"
// @Param request body AddMerchantStaffRequest true "添加员工请求"
// @Success 200 {object} models.MerchantStaff
// @Router /api/v1/merchants/{id}/staff [post]
func (ctrl *MerchantController) AddMerchantStaff(c *gin.Context) {
	merchantID := c.Param("id")
	var req AddMerchantStaffRequest
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

	// 获取商家
	var merchant models.Merchant
	if err := ctrl.db.Where("id = ?", merchantID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家不存在"})
		return
	}

	// 权限检查：只有商家管理员可以添加员工
	if !utils.IsSuperAdmin(user) && merchant.AdminID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限添加员工"})
		return
	}

	// 检查用户是否存在
	var targetUser models.User
	if err := ctrl.db.Where("id = ?", req.UserID).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查用户是否已经是该商家的员工
	var existingStaff models.MerchantStaff
	if err := ctrl.db.Where("user_id = ? AND merchant_id = ?", req.UserID, merchantID).First(&existingStaff).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该用户已经是该商家的员工"})
		return
	}

	// 创建员工
	staff := models.MerchantStaff{
		UserID:     req.UserID,
		MerchantID: merchant.ID,
		Title:      req.Title,
		Status:     "active",
	}

	if err := ctrl.db.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加员工失败"})
		return
	}

	// 添加权限
	for _, permCode := range req.Permissions {
		permission := models.MerchantStaffPermission{
			StaffID:        staff.ID,
			PermissionCode: permCode,
			GrantedBy:      staff.ID, // 暂时设置为自己，实际应该是管理员ID
		}
		ctrl.db.Create(&permission)
	}

	// 为用户添加 merchant_staff 角色
	if !utils.HasRole(&targetUser, "merchant_staff") {
		targetUser.Roles = append(targetUser.Roles, "merchant_staff")
		ctrl.db.Save(&targetUser)
	}

	c.JSON(http.StatusOK, staff)
}

// GetMerchantStaff 获取商家员工列表
// @Summary 获取商家员工列表
// @Description 获取指定商家的员工列表
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param id path string true "商家ID"
// @Success 200 {array} models.MerchantStaff
// @Router /api/v1/merchants/{id}/staff [get]
func (ctrl *MerchantController) GetMerchantStaff(c *gin.Context) {
	merchantID := c.Param("id")
	var staff []models.MerchantStaff

	if err := ctrl.db.Where("merchant_id = ?", merchantID).Preload("User").Preload("Permissions").Find(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取员工列表失败"})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// UpdateMerchantStaffPermission 更新员工权限
// @Summary 更新员工权限
// @Description 授予或撤销员工权限
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param merchant_id path string true "商家ID"
// @Param staff_id path string true "员工ID"
// @Param request body UpdateMerchantStaffPermissionRequest true "更新权限请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{merchant_id}/staff/{staff_id}/permissions [put]
func (ctrl *MerchantController) UpdateMerchantStaffPermission(c *gin.Context) {
	merchantID := c.Param("merchant_id")
	staffID := c.Param("staff_id")
	var req UpdateMerchantStaffPermissionRequest
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

	// 获取商家
	var merchant models.Merchant
	if err := ctrl.db.Where("id = ?", merchantID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家不存在"})
		return
	}

	// 权限检查：只有商家管理员可以更新员工权限
	if !utils.IsSuperAdmin(user) && merchant.AdminID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新员工权限"})
		return
	}

	// 获取员工
	var staff models.MerchantStaff
	if err := ctrl.db.Where("id = ? AND merchant_id = ?", staffID, merchantID).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "员工不存在"})
		return
	}

	if req.Action == "grant" {
		// 授予权限
		var existingPerm models.MerchantStaffPermission
		if err := ctrl.db.Where("staff_id = ? AND permission_code = ?", staffID, req.PermissionCode).First(&existingPerm).Error; err != nil {
			permission := models.MerchantStaffPermission{
				StaffID:        staff.ID,
				PermissionCode: req.PermissionCode,
				GrantedBy:      staff.ID,
			}
			ctrl.db.Create(&permission)
		}
	} else if req.Action == "revoke" {
		// 撤销权限
		ctrl.db.Where("staff_id = ? AND permission_code = ?", staffID, req.PermissionCode).Delete(&models.MerchantStaffPermission{})
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限更新成功"})
}

// DeleteMerchantStaff 删除商家员工
// @Summary 删除商家员工
// @Description 从商家移除员工
// @Tags 商家管理
// @Accept json
// @Produce json
// @Param merchant_id path string true "商家ID"
// @Param staff_id path string true "员工ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/merchants/{merchant_id}/staff/{staff_id} [delete]
func (ctrl *MerchantController) DeleteMerchantStaff(c *gin.Context) {
	merchantID := c.Param("merchant_id")
	staffID := c.Param("staff_id")

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 获取商家
	var merchant models.Merchant
	if err := ctrl.db.Where("id = ?", merchantID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家不存在"})
		return
	}

	// 权限检查：只有商家管理员可以删除员工
	if !utils.IsSuperAdmin(user) && merchant.AdminID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除员工"})
		return
	}

	// 删除员工权限
	ctrl.db.Where("staff_id = ?", staffID).Delete(&models.MerchantStaffPermission{})

	// 删除员工
	if err := ctrl.db.Where("id = ? AND merchant_id = ?", staffID, merchantID).Delete(&models.MerchantStaff{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除员工失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetMyMerchant 获取当前用户管理的商家
// @Summary 获取当前用户管理的商家
// @Description 商家管理员获取自己管理的商家信息
// @Tags 商家管理
// @Accept json
// @Produce json
// @Success 200 {object} models.Merchant
// @Router /api/v1/merchant/me [get]
func (ctrl *MerchantController) GetMyMerchant(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	var merchant models.Merchant
	if err := ctrl.db.Where("admin_id = ?", user.ID).Preload("Admin").Preload("Provider").Preload("User").Preload("Staff").First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "您不是商家管理员"})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

// GetPermissions 获取可用的权限列表
// @Summary 获取权限列表
// @Description 获取可用于授予员工的权限列表
// @Tags 商家管理
// @Accept json
// @Produce json
// @Success 200 {array} constants.PermissionDefinition
// @Router /api/v1/merchants/permissions [get]
func (ctrl *MerchantController) GetPermissions(c *gin.Context) {
	// 返回所有权限定义（前端会根据员工角色过滤）
	c.JSON(http.StatusOK, constants.GetAllPermissionDefinitions())
}
