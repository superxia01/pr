package controllers

import (
	"fmt"
	"net/http"
	"pr-business/config"
	"pr-business/constants"
	"pr-business/models"
	"pr-business/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvitationController struct {
	db  *gorm.DB
	cfg *config.Config
}

// getUserRoles 从上下文或数据库获取用户角色
func (ctrl *InvitationController) getUserRoles(c *gin.Context) ([]string, error) {
	// 优先从上下文获取
	userID := c.GetString("userId")
	var user models.User
	if err := ctrl.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return []string(user.Roles), nil
}

func NewInvitationController(db *gorm.DB, cfg *config.Config) *InvitationController {
	return &InvitationController{
		db:  db,
		cfg: cfg,
	}
}

// CreateInvitationCodeRequest 创建邀请码请求
type CreateInvitationCodeRequest struct {
	CodeType  string `json:"codeType" binding:"required"`
	OwnerID   string `json:"ownerId" binding:"required"`
	OwnerType string `json:"ownerType" binding:"required"`
	MaxUses   int    `json:"maxUses"`
	ExpiresAt string `json:"expiresAt"` // ISO 8601格式
}

// UseInvitationCodeRequest 使用邀请码请求
type UseInvitationCodeRequest struct {
	Code   string `json:"code" binding:"required"`
	UserID string `json:"userId" binding:"required"`
}

// CreateInvitationCode 创建邀请码
func (ctrl *InvitationController) CreateInvitationCode(c *gin.Context) {
	var req CreateInvitationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 验证权限：只有超级管理员、服务商管理员、商家管理员可以创建邀请码
	roles, err := ctrl.getUserRoles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户角色失败"})
		return
	}
	if !ctrl.canCreateInvitationCode(roles, req.CodeType, req.OwnerType) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions to create this type of invitation code",
		})
		return
	}

	// 验证OwnerID是否存在
	if !ctrl.validateOwnerID(req.OwnerType, req.OwnerID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Owner not found",
		})
		return
	}

	// 生成固定邀请码
	code := utils.GenerateFixedInvitationCode(req.CodeType, req.OwnerID)

	// 检查邀请码是否已存在
	var existingCode models.InvitationCode
	result := ctrl.db.Where("code = ?", code).First(&existingCode)
	if result.Error == nil {
		// 邀请码已存在，返回已存在的邀请码
		c.JSON(http.StatusOK, existingCode)
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// 创建邀请码
	invitationCode := models.InvitationCode{
		Code:           code,
		Type:            req.CodeType,
		TargetRole:      req.CodeType, // 邀请码类型即为目标角色
		GeneratorID:      c.GetString("userId"),
		GeneratorType:    req.OwnerType, // 所有者类型即生成者类型
		OrganizationID:    nil, // 固定邀请码暂不绑定具体组织
		IsActive:        true,
		UseCount:        0,
	}

	if err := ctrl.db.Create(&invitationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create invitation code",
		})
		return
	}

	c.JSON(http.StatusCreated, invitationCode)
}

// GetInvitationCode 获取邀请码详情
func (ctrl *InvitationController) GetInvitationCode(c *gin.Context) {
	code := c.Param("code")

	var invitationCode models.InvitationCode
	result := ctrl.db.Where("code = ?", code).First(&invitationCode)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invitation code not found",
		})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, invitationCode)
}

// ListInvitationCodes 获取邀请码列表
func (ctrl *InvitationController) ListInvitationCodes(c *gin.Context) {
	// 查询参数
	codeType := c.Query("codeType")
	ownerID := c.Query("ownerId")
	status := c.Query("status")

	query := ctrl.db.Model(&models.InvitationCode{})

	// 权限过滤：用户只能看到自己创建的或自己类型的邀请码
	userID := c.GetString("userId")
	roles, err := ctrl.getUserRoles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户角色失败"})
		return
	}

	// 检查是否有超级管理员角色
	hasSuperAdmin := false
	for _, role := range roles {
		if role == "SUPER_ADMIN" {
			hasSuperAdmin = true
			break
		}
	}

	// 超级管理员可以看到所有邀请码
	if !hasSuperAdmin {
		// 其他角色只能看到自己的邀请码
		query = query.Where("owner_id = ?", userID)
	}

	// 应用过滤条件
	if codeType != "" {
		query = query.Where("code_type = ?", codeType)
	}
	if ownerID != "" {
		query = query.Where("owner_id = ?", ownerID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var invitationCodes []models.InvitationCode
	result := query.Order("created_at DESC").Find(&invitationCodes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"codes": invitationCodes,
		"total": len(invitationCodes),
	})
}

// UseInvitationCode 使用邀请码（从数据库验证）
func (ctrl *InvitationController) UseInvitationCode(c *gin.Context) {
	var req UseInvitationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 从数据库查找邀请码
	var invitationCode models.InvitationCode
	result := ctrl.db.Where("code = ?", req.Code).First(&invitationCode)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邀请码不存在",
		})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询邀请码失败",
		})
		return
	}

	// 检查邀请码是否激活
	if !invitationCode.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邀请码已被禁用",
		})
		return
	}

	targetRole := invitationCode.TargetRole

	// 验证组织ID是否有效（如果需要组织绑定）
	var orgID string
	if invitationCode.OrganizationID != nil {
		orgID = invitationCode.OrganizationID.String()
	}

	if utils.RequiresOrganizationBinding(targetRole) {
		// 验证组织ID是否存在
		if orgID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "此角色需要选择组织",
			})
			return
		}

		// 根据角色类型验证组织是否存在
		orgType := utils.GetOrganizationTypeByRole(targetRole)
		if orgType == "service_provider" {
			var sp models.ServiceProvider
			if err := ctrl.db.Where("id = ?", orgID).First(&sp).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "服务商不存在",
				})
				return
			}
		} else if orgType == "merchant" {
			var m models.Merchant
			if err := ctrl.db.Where("id = ?", orgID).First(&m).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "商家不存在",
				})
				return
			}
		}
	}

	// 更新使用次数
	invitationCode.UseCount += 1
	if err := ctrl.db.Save(&invitationCode).Error; err != nil {
		fmt.Printf("[UseInvitationCode] Failed to update use count: %v\n", err)
	}

	// 给予被邀请用户目标角色
	// 获取当前用户信息
	var user models.User
	if err := ctrl.db.Where("id = ?", c.GetString("userId")).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 检查用户是否已有该角色
	hasRole := false
	for _, role := range user.Roles {
		if role == targetRole {
			hasRole = true
			break
		}
	}

	// 如果没有该角色，添加角色
	if !hasRole {
		user.Roles = append(user.Roles, targetRole)
		if err := ctrl.db.Save(&user).Error; err != nil {
			fmt.Printf("[UseInvitationCode] Failed to add role to user: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加角色失败"})
			return
		}
		fmt.Printf("[UseInvitationCode] Added role %s to user %s\n", targetRole, user.ID)
	} else {
		fmt.Printf("[UseInvitationCode] User %s already has role %s\n", user.ID, targetRole)
	}

	// 如果需要组织绑定，设置用户为该组织的管理员
	if orgID != "" {
		orgType := utils.GetOrganizationTypeByRole(targetRole)
		if targetRole == "SERVICE_PROVIDER_ADMIN" && orgType == "service_provider" {
			// 设置用户为服务商的管理员
			if err := ctrl.db.Model(&models.ServiceProvider{}).
				Where("id = ?", orgID).
				Update("admin_id", user.ID).Error; err != nil {
				fmt.Printf("[UseInvitationCode] Failed to set service provider admin: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "设置服务商管理员失败"})
				return
			}
			fmt.Printf("[UseInvitationCode] Set user %s as admin of service provider %s\n", user.ID, orgID)
		} else if targetRole == "MERCHANT_ADMIN" && orgType == "merchant" {
			// 设置用户为商家的管理员
			if err := ctrl.db.Model(&models.Merchant{}).
				Where("id = ?", orgID).
				Update("admin_id", user.ID).Error; err != nil {
				fmt.Printf("[UseInvitationCode] Failed to set merchant admin: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "设置商家管理员失败"})
				return
			}
			fmt.Printf("[UseInvitationCode] Set user %s as admin of merchant %s\n", user.ID, orgID)
		} else if targetRole == "SERVICE_PROVIDER_STAFF" && orgType == "service_provider" {
			// 创建服务商员工记录
			var existingStaff models.ServiceProviderStaff
			err := ctrl.db.Where("user_id = ? AND provider_id = ?", user.ID, orgID).First(&existingStaff).Error
			if err == gorm.ErrRecordNotFound {
				// 不存在，创建新记录
				orgUUID, _ := uuid.Parse(orgID)
				staff := models.ServiceProviderStaff{
					UserID:     user.ID,
					ProviderID: orgUUID,
					Title:      "员工",
					Status:     "active",
				}
				if err := ctrl.db.Create(&staff).Error; err != nil {
					fmt.Printf("[UseInvitationCode] Failed to create service provider staff: %v\n", err)
					// 不返回错误，因为角色已经添加成功
				} else {
					fmt.Printf("[UseInvitationCode] Created service provider staff record for user %s, provider %s\n", user.ID, orgID)
				}
			} else if err != nil {
				fmt.Printf("[UseInvitationCode] Error checking service provider staff: %v\n", err)
			} else {
				fmt.Printf("[UseInvitationCode] Service provider staff already exists for user %s, provider %s\n", user.ID, orgID)
			}
		} else if targetRole == "MERCHANT_STAFF" && orgType == "merchant" {
			// 创建商家员工记录
			var existingStaff models.MerchantStaff
			err := ctrl.db.Where("user_id = ? AND merchant_id = ?", user.ID, orgID).First(&existingStaff).Error
			if err == gorm.ErrRecordNotFound {
				// 不存在，创建新记录
				orgUUID, _ := uuid.Parse(orgID)
				staff := models.MerchantStaff{
					UserID:     user.ID,
					MerchantID: orgUUID,
					Title:      "员工",
					Status:     "active",
				}
				if err := ctrl.db.Create(&staff).Error; err != nil {
					fmt.Printf("[UseInvitationCode] Failed to create merchant staff: %v\n", err)
					// 不返回错误，因为角色已经添加成功
				} else {
					fmt.Printf("[UseInvitationCode] Created merchant staff record for user %s, merchant %s\n", user.ID, orgID)
				}
			} else if err != nil {
				fmt.Printf("[UseInvitationCode] Error checking merchant staff: %v\n", err)
			} else {
				fmt.Printf("[UseInvitationCode] Merchant staff already exists for user %s, merchant %s\n", user.ID, orgID)
			}
		}
	}

	// 记录邀请关系（可选，用于追踪邀请历史）
	// TODO: 可以在这里创建 InvitationRelationship 记录

	c.JSON(http.StatusOK, gin.H{
		"message": "邀请码使用成功",
		"targetRole": targetRole,
		"organizationId": orgID,
	})
}

// DisableInvitationCode 禁用邀请码
func (ctrl *InvitationController) DisableInvitationCode(c *gin.Context) {
	code := c.Param("code")
	userID := c.GetString("userId")
	roles, err := ctrl.getUserRoles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户角色失败"})
		return
	}

	var invitationCode models.InvitationCode
	result := ctrl.db.Where("code = ?", code).First(&invitationCode)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invitation code not found",
		})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// 检查是否有超级管理员角色
	hasSuperAdmin := false
	for _, role := range roles {
		if role == "SUPER_ADMIN" {
			hasSuperAdmin = true
			break
		}
	}

	// 权限检查：只有创建者或超级管理员可以禁用
	if !hasSuperAdmin && invitationCode.GeneratorID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
		return
	}

	// 禁用邀请码
	invitationCode.IsActive = false
	if err := ctrl.db.Save(&invitationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to disable invitation code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Invitation code disabled successfully",
	})
}

// validateOwnerID 验证所有者是否存在
func (ctrl *InvitationController) validateOwnerID(ownerType, ownerID string) bool {
	// TODO: 实现实际的所有者验证
	// 这里简化处理，实际应该查询对应的表
	return true
}

// canCreateInvitationCode 检查是否可以创建指定类型的邀请码（基于用户拥有的所有角色）
func (ctrl *InvitationController) canCreateInvitationCode(roles []string, codeType, ownerType string) bool {
	// 检查用户的所有角色
	for _, role := range roles {
		switch role {
		case "SUPER_ADMIN":
			// 超级管理员可以创建所有类型的邀请码
			return true
		case "SERVICE_PROVIDER_ADMIN", "SP_ADMIN":
			// 服务商管理员可以创建商家、达人邀请码
			if codeType == "MERCHANT" || codeType == "CREATOR" || codeType == "STAFF" {
				return true
			}
		case "MERCHANT_ADMIN":
			// 商家管理员可以创建员工邀请码
			if codeType == "STAFF" {
				return true
			}
		}
	}
	return false
}

// GetMyFixedInvitationCodes 获取我的固定邀请码列表
// 返回基于用户所有角色可邀请的角色类型及其对应的邀请码
// GET /api/v1/invitations/fixed-codes
func (ctrl *InvitationController) GetMyFixedInvitationCodes(c *gin.Context) {
	userID := c.GetString("userId")

	// 获取用户信息
	var user models.User
	if err := ctrl.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 从数据库用户对象获取角色（models.Roles -> []string）
	roles := []string(user.Roles)

	// 获取所有可邀请的角色（已内置去重逻辑）
	invitableRoles := constants.GetInvitableRoles(roles)
	invitableRoleMap := make(map[string]constants.InvitableRole)
	for _, ir := range invitableRoles {
		invitableRoleMap[ir.Role] = ir
	}

	if len(invitableRoleMap) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"invitableRoles": []interface{}{},
			"userRoles":      roles,
		})
		return
	}

	// 获取可用的组织（服务商和商家）
	type OrganizationInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}

	organizations := make([]OrganizationInfo, 0)

	isSuperAdmin := utils.IsSuperAdmin(&user)
	isServiceProviderAdmin := utils.IsServiceProviderAdmin(&user)
	isMerchantAdmin := utils.IsMerchantAdmin(&user)

	// 超级管理员：获取所有服务商和商家（用于生成邀请码邀请别人成为管理员）
	// 服务商管理员：获取自己是管理员的服务商
	// 商家管理员：没有服务商
	var serviceProviders []models.ServiceProvider
	if isSuperAdmin {
		// 超级管理员：返回所有服务商
		if err := ctrl.db.Model(&models.ServiceProvider{}).Find(&serviceProviders).Error; err == nil {
			for _, sp := range serviceProviders {
				organizations = append(organizations, OrganizationInfo{
					ID:   sp.ID.String(),
					Name: sp.Name,
					Type: "service_provider",
				})
			}
		}
	} else if isServiceProviderAdmin {
		// 服务商管理员：返回自己是管理员的服务商
		if err := ctrl.db.Where("admin_id = ?", userID).Find(&serviceProviders).Error; err == nil {
			for _, sp := range serviceProviders {
				organizations = append(organizations, OrganizationInfo{
					ID:   sp.ID.String(),
					Name: sp.Name,
					Type: "service_provider",
				})
			}
		}
	}

	var merchants []models.Merchant

	if isSuperAdmin {
		// 超级管理员：返回所有商家
		if err := ctrl.db.Model(&models.Merchant{}).Find(&merchants).Error; err == nil {
			for _, m := range merchants {
				organizations = append(organizations, OrganizationInfo{
					ID:   m.ID.String(),
					Name: m.Name,
					Type: "merchant",
				})
			}
		}
	} else if isServiceProviderAdmin {
		// 服务商管理员：返回自己创建的商家（user_id = self）
		if err := ctrl.db.Where("user_id = ?", userID).Find(&merchants).Error; err == nil {
			for _, m := range merchants {
				organizations = append(organizations, OrganizationInfo{
					ID:   m.ID.String(),
					Name: m.Name,
					Type: "merchant",
				})
			}
		}
	} else if isMerchantAdmin {
		// 商家管理员：返回自己是管理员的商家
		if err := ctrl.db.Where("admin_id = ?", userID).Find(&merchants).Error; err == nil {
			for _, m := range merchants {
				organizations = append(organizations, OrganizationInfo{
					ID:   m.ID.String(),
					Name: m.Name,
					Type: "merchant",
				})
			}
		}
	}

	// 为每个可邀请的角色生成对应的邀请码（按角色和组织的组合）
	type InvitationCodeWithOrg struct {
		Code      string `json:"code"`
		Role      string `json:"role"`
		RoleLabel string `json:"roleLabel"`
		OrgID     string `json:"orgId,omitempty"`
		OrgName   string `json:"orgName,omitempty"`
		OrgType   string `json:"orgType,omitempty"`
	}

	codes := make([]InvitationCodeWithOrg, 0)

	for _, ir := range invitableRoleMap {
		needOrg := utils.RequiresOrganizationBinding(ir.Role)

		if needOrg {
			// 需要组织绑定：为该角色的每个组织生成一个邀请码
			orgType := utils.GetOrganizationTypeByRole(ir.Role)
			filteredOrgs := make([]OrganizationInfo, 0)

			// 筛选该角色对应的组织
			for _, org := range organizations {
				if org.Type == orgType {
					filteredOrgs = append(filteredOrgs, org)
				}
			}

			// 为每个组织生成邀请码
			for _, org := range filteredOrgs {
				// 查询是否已存在该邀请码（根据 generator_id + target_role + organization_id）
				var existingCode models.InvitationCode
				orgUUID, _ := uuid.Parse(org.ID)
				err := ctrl.db.Where("generator_id = ? AND target_role = ? AND organization_id = ?",
					userID, ir.Role, orgUUID).First(&existingCode).Error

				var code string
				if err == gorm.ErrRecordNotFound {
					// 不存在，生成新短码
					code = utils.GenerateUserFixedInvitationCode(userID, ir.Role, org.ID)

					newCode := models.InvitationCode{
						Code:             code,
						Type:             "FIXED",
						TargetRole:       ir.Role,
						GeneratorID:      userID,
						GeneratorType:    "user",
						OrganizationID:   &orgUUID,
						OrganizationType: &org.Type,
						IsActive:         true,
						UseCount:         0,
					}
					if err := ctrl.db.Create(&newCode).Error; err != nil {
						fmt.Printf("[GetMyFixedInvitationCodes] Failed to create invitation code %s: %v\n", code, err)
						continue
					}
					fmt.Printf("[GetMyFixedInvitationCodes] Created invitation code %s for role %s, org %s\n", code, ir.Role, org.Name)
				} else if err != nil {
					fmt.Printf("[GetMyFixedInvitationCodes] Error checking invitation code: %v\n", err)
					continue
				} else {
					// 已存在，使用现有邀请码
					code = existingCode.Code
					fmt.Printf("[GetMyFixedInvitationCodes] Reusing existing code %s for role %s, org %s\n", code, ir.Role, org.Name)
				}

				codes = append(codes, InvitationCodeWithOrg{
					Code:      code,
					Role:      ir.Role,
					RoleLabel: ir.Label,
					OrgID:     org.ID,
					OrgName:   org.Name,
					OrgType:   org.Type,
				})
			}
		} else {
			// 不需要组织绑定（如达人）：查询是否已存在该邀请码
			var existingCode models.InvitationCode
			err := ctrl.db.Where("generator_id = ? AND target_role = ? AND organization_id IS NULL",
				userID, ir.Role).First(&existingCode).Error

			var code string
			if err == gorm.ErrRecordNotFound {
				// 不存在，生成新短码
				code = utils.GenerateUserFixedInvitationCode(userID, ir.Role)

				newCode := models.InvitationCode{
					Code:          code,
					Type:          "FIXED",
					TargetRole:    ir.Role,
					GeneratorID:   userID,
					GeneratorType: "user",
					IsActive:      true,
					UseCount:      0,
				}
				if err := ctrl.db.Create(&newCode).Error; err != nil {
					fmt.Printf("[GetMyFixedInvitationCodes] Failed to create invitation code %s: %v\n", code, err)
					continue
				}
				fmt.Printf("[GetMyFixedInvitationCodes] Created invitation code %s for role %s\n", code, ir.Role)
			} else if err != nil {
				fmt.Printf("[GetMyFixedInvitationCodes] Error checking invitation code: %v\n", err)
				continue
			} else {
				// 已存在，使用现有邀请码
				code = existingCode.Code
				fmt.Printf("[GetMyFixedInvitationCodes] Reusing existing code %s for role %s\n", code, ir.Role)
			}

			codes = append(codes, InvitationCodeWithOrg{
				Code:      code,
				Role:      ir.Role,
				RoleLabel: ir.Label,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"invitationCodes": codes,
		"userRoles":       roles,
		"organizations":   organizations,
	})
}

// GetMyInvitations 获取我发出的邀请列表
// GET /api/v1/invitations/my
func (ctrl *InvitationController) GetMyInvitations(c *gin.Context) {
	userID := c.GetString("userId")

	// 查询我发出的邀请
	var invitations []models.InvitationRelationship
	if err := ctrl.db.Where("inviter_id = ?", userID).
		Preload("Invitee").
		Order("invited_at DESC").
		Find(&invitations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取邀请列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invitations": invitations,
		"total":      len(invitations),
	})
}
