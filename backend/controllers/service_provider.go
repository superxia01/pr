package controllers

import (
	"net/http"
	"pr-business/models"
	"pr-business/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceProviderController struct {
	db *gorm.DB
}

func NewServiceProviderController(db *gorm.DB) *ServiceProviderController {
	return &ServiceProviderController{db: db}
}

// CreateServiceProviderRequest 创建服务商请求
type CreateServiceProviderRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
}

// UpdateServiceProviderRequest 更新服务商请求
type UpdateServiceProviderRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=active suspended inactive"`
}

// AddServiceProviderStaffRequest 添加服务商员工请求
type AddServiceProviderStaffRequest struct {
	UserID      string   `json:"userId" binding:"required"`
	Title       string      `json:"title" binding:"omitempty,max=50"`
	Permissions []string    `json:"permissions" binding:"required"`
}

// UpdateServiceProviderStaffPermissionRequest 更新员工权限请求
type UpdateServiceProviderStaffPermissionRequest struct {
	PermissionCode string `json:"permissionCode" binding:"required"`
	Action         string `json:"action" binding:"required,oneof=grant revoke"`
}

// CreateServiceProvider 创建服务商
// @Summary 创建服务商
// @Description 只有超级管理员可以创建服务商
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param request body CreateServiceProviderRequest true "创建服务商请求"
// @Success 200 {object} models.ServiceProvider
// @Router /api/v1/service-providers [post]
func (ctrl *ServiceProviderController) CreateServiceProvider(c *gin.Context) {
	var req CreateServiceProviderRequest
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

	// 权限检查：只有超级管理员可以创建服务商
	if !utils.HasRole(user, "super_admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限创建服务商"})
		return
	}

	// 创建服务商（admin_id 暂时为空，后续通过邀请码绑定管理员）
	provider := models.ServiceProvider{
		UserID:      user.ID, // 记录创建者
		Name:        req.Name,
		Description: req.Description,
		Status:      "active",
	}

	if err := ctrl.db.Create(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建服务商失败"})
		return
	}

	c.JSON(http.StatusOK, provider)
}
// GetServiceProviders 获取服务商列表
// @Summary 获取服务商列表
// @Description 获取服务商列表，支持过滤
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param status query string false "状态过滤"
// @Success 200 {array} models.ServiceProvider
// @Router /api/v1/service-providers [get]
func (ctrl *ServiceProviderController) GetServiceProviders(c *gin.Context) {
	var providers []models.ServiceProvider
	query := ctrl.db.Model(&models.ServiceProvider{})

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Preload("Admin").Preload("User").Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取服务商列表失败"})
		return
	}

	c.JSON(http.StatusOK, providers)
}

// GetServiceProvider 获取服务商详情
// @Summary 获取服务商详情
// @Description 根据ID获取服务商详情
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param id path string true "服务商ID"
// @Success 200 {object} models.ServiceProvider
// @Router /api/v1/service-providers/{id} [get]
func (ctrl *ServiceProviderController) GetServiceProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.ServiceProvider

	if err := ctrl.db.Where("id = ?", id).Preload("Admin").Preload("User").Preload("Staff").Preload("Merchants").First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务商不存在"})
		return
	}

	c.JSON(http.StatusOK, provider)
}

// UpdateServiceProvider 更新服务商信息
// @Summary 更新服务商信息
// @Description 只有超级管理员和该服务商的管理员可以更新
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param id path string true "服务商ID"
// @Param request body UpdateServiceProviderRequest true "更新服务商请求"
// @Success 200 {object} models.ServiceProvider
// @Router /api/v1/service-providers/{id} [put]
func (ctrl *ServiceProviderController) UpdateServiceProvider(c *gin.Context) {
	id := c.Param("id")
	var req UpdateServiceProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var provider models.ServiceProvider
	if err := ctrl.db.Where("id = ?", id).First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务商不存在"})
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
	hasPermission := utils.HasRole(user, "super_admin") ||
		(utils.HasRole(user, "provider_admin") && user.ID == provider.AdminID)

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新服务商信息"})
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
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := ctrl.db.Model(&provider).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新服务商失败"})
		return
	}

	c.JSON(http.StatusOK, provider)
}

// DeleteServiceProvider 删除服务商（软删除）
// @Summary 删除服务商
// @Description 只有超级管理员可以删除服务商
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param id path string true "服务商ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/service-providers/{id} [delete]
func (ctrl *ServiceProviderController) DeleteServiceProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.ServiceProvider

	if err := ctrl.db.Where("id = ?", id).First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务商不存在"})
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
	if !utils.HasRole(user, "super_admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除服务商"})
		return
	}

	if err := ctrl.db.Delete(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除服务商失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// AddServiceProviderStaff 添加服务商员工
// @Summary 添加服务商员工
// @Description 只有服务商管理员可以添加员工
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param id path string true "服务商ID"
// @Param request body AddServiceProviderStaffRequest true "添加员工请求"
// @Success 200 {object} models.ServiceProviderStaff
// @Router /api/v1/service-providers/{id}/staff [post]
func (ctrl *ServiceProviderController) AddServiceProviderStaff(c *gin.Context) {
	providerID := c.Param("id")
	var req AddServiceProviderStaffRequest
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

	// 获取服务商
	var provider models.ServiceProvider
	if err := ctrl.db.Where("id = ?", providerID).First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务商不存在"})
		return
	}

	// 权限检查：只有服务商管理员可以添加员工
	if !utils.HasRole(user, "super_admin") && provider.AdminID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限添加员工"})
		return
	}

	// 检查用户是否存在
	var targetUser models.User
	if err := ctrl.db.Where("id = ?", req.UserID).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查用户是否已经是员工（任何服务商的员工）
	var existingStaff models.ServiceProviderStaff
	if err := ctrl.db.Where("user_id = ?", req.UserID).First(&existingStaff).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该用户已经是服务商员工"})
		return
	}

	// 创建员工
	staff := models.ServiceProviderStaff{
		UserID:     req.UserID,
		ProviderID: provider.ID,
		Title:      req.Title,
		Status:     "active",
	}

	if err := ctrl.db.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加员工失败"})
		return
	}

	// 添加权限
	for _, permCode := range req.Permissions {
		permission := models.ServiceProviderStaffPermission{
			StaffID:        staff.ID,
			PermissionCode: permCode,
			GrantedBy:      staff.ID,
		}
		ctrl.db.Create(&permission)
	}

	// 为用户添加 SP_STAFF 角色
	if !utils.HasRole(&targetUser, "SP_STAFF") {
		targetUser.Roles = append(targetUser.Roles, "SP_STAFF")
		ctrl.db.Save(&targetUser)
	}

	c.JSON(http.StatusOK, staff)
}

// GetServiceProviderStaff 获取服务商员工列表
// @Summary 获取服务商员工列表
// @Description 获取指定服务商的员工列表
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param id path string true "服务商ID"
// @Success 200 {array} models.ServiceProviderStaff
// @Router /api/v1/service-providers/{id}/staff [get]
func (ctrl *ServiceProviderController) GetServiceProviderStaff(c *gin.Context) {
	providerID := c.Param("id")
	var staff []models.ServiceProviderStaff

	if err := ctrl.db.Where("provider_id = ?", providerID).Preload("User").Preload("Permissions").Find(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取员工列表失败"})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// UpdateServiceProviderStaffPermission 更新员工权限
// @Summary 更新员工权限
// @Description 授予或撤销员工权限
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param provider_id path string true "服务商ID"
// @Param staff_id path string true "员工ID"
// @Param request body UpdateServiceProviderStaffPermissionRequest true "更新权限请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/service-providers/{provider_id}/staff/{staff_id}/permissions [put]
func (ctrl *ServiceProviderController) UpdateServiceProviderStaffPermission(c *gin.Context) {
	providerID := c.Param("provider_id")
	staffID := c.Param("staff_id")
	var req UpdateServiceProviderStaffPermissionRequest
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

	// 获取服务商
	var provider models.ServiceProvider
	if err := ctrl.db.Where("id = ?", providerID).First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务商不存在"})
		return
	}

	// 权限检查：只有服务商管理员可以更新员工权限
	if !utils.HasRole(user, "super_admin") && provider.AdminID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新员工权限"})
		return
	}

	// 获取员工
	var staff models.ServiceProviderStaff
	if err := ctrl.db.Where("id = ? AND provider_id = ?", staffID, providerID).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "员工不存在"})
		return
	}

	if req.Action == "grant" {
		// 授予权限
		var existingPerm models.ServiceProviderStaffPermission
		if err := ctrl.db.Where("staff_id = ? AND permission_code = ?", staffID, req.PermissionCode).First(&existingPerm).Error; err != nil {
			permission := models.ServiceProviderStaffPermission{
				StaffID:        staff.ID,
				PermissionCode: req.PermissionCode,
				GrantedBy:      staff.ID,
			}
			ctrl.db.Create(&permission)
		}
	} else if req.Action == "revoke" {
		// 撤销权限
		ctrl.db.Where("staff_id = ? AND permission_code = ?", staffID, req.PermissionCode).Delete(&models.ServiceProviderStaffPermission{})
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限更新成功"})
}

// DeleteServiceProviderStaff 删除服务商员工
// @Summary 删除服务商员工
// @Description 从服务商移除员工
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Param provider_id path string true "服务商ID"
// @Param staff_id path string true "员工ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/service-providers/{provider_id}/staff/{staff_id} [delete]
func (ctrl *ServiceProviderController) DeleteServiceProviderStaff(c *gin.Context) {
	providerID := c.Param("provider_id")
	staffID := c.Param("staff_id")

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 获取服务商
	var provider models.ServiceProvider
	if err := ctrl.db.Where("id = ?", providerID).First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务商不存在"})
		return
	}

	// 权限检查：只有服务商管理员可以删除员工
	if !utils.HasRole(user, "super_admin") && provider.AdminID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除员工"})
		return
	}

	// 删除员工权限
	ctrl.db.Where("staff_id = ?", staffID).Delete(&models.ServiceProviderStaffPermission{})

	// 删除员工
	if err := ctrl.db.Where("id = ? AND provider_id = ?", staffID, providerID).Delete(&models.ServiceProviderStaff{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除员工失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetMyServiceProvider 获取当前用户管理的服务商
// @Summary 获取当前用户管理的服务商
// @Description 服务商管理员获取自己管理的服务商信息
// @Tags 服务商管理
// @Accept json
// @Produce json
// @Success 200 {object} models.ServiceProvider
// @Router /api/v1/service-provider/me [get]
func (ctrl *ServiceProviderController) GetMyServiceProvider(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	var provider models.ServiceProvider
	if err := ctrl.db.Where("admin_id = ?", user.ID).Preload("Admin").Preload("User").Preload("Staff").Preload("Merchants").First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "您不是服务商管理员"})
		return
	}

	c.JSON(http.StatusOK, provider)
}
