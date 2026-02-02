package controllers

import (
	"net/http"
	"pr-business/config"
	"pr-business/models"
	"pr-business/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvitationController struct {
	db  *gorm.DB
	cfg *config.Config
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

	userID := c.GetString("userId")

	// 验证权限：只有超级管理员、服务商管理员、商家管理员可以创建邀请码
	currentRole := c.GetString("currentRole")
	if !ctrl.canCreateInvitationCode(currentRole, req.CodeType, req.OwnerType) {
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

	// 解析过期时间
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid expiresAt format, use ISO 8601",
			})
			return
		}
		expiresAt = &parsedTime
	}

	// 创建邀请码
	invitationCode := models.InvitationCode{
		Code:      code,
		CodeType:  req.CodeType,
		OwnerID:   req.OwnerID,
		OwnerType: req.OwnerType,
		Status:    "active",
		MaxUses:   req.MaxUses,
		UseCount:  0,
		ExpiresAt: expiresAt,
		UsedBy:    []string{},
		Metadata:  "{}",
	}

	if err := ctrl.db.Create(&invitationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create invitation code",
		})
		return
	}

	// 记录创建者
	ctrl.db.Model(&invitationCode).Update("metadata", gin.H{
		"createdBy": userID,
		"createdAt": time.Now(),
	})

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
	currentRole := c.GetString("currentRole")

	// 超级管理员可以看到所有邀请码
	if currentRole != "super_admin" {
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

// UseInvitationCode 使用邀请码
func (ctrl *InvitationController) UseInvitationCode(c *gin.Context) {
	var req UseInvitationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 查找邀请码
	var invitationCode models.InvitationCode
	result := ctrl.db.Where("code = ?", req.Code).First(&invitationCode)

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

	// 检查邀请码是否可用
	if !invitationCode.CanBeUsed() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invitation code is not available",
			"status": invitationCode.Status,
		})
		return
	}

	// 检查用户是否已使用过
	for _, usedBy := range invitationCode.UsedBy {
		if usedBy == req.UserID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "You have already used this invitation code",
			})
			return
		}
	}

	// 使用邀请码
	invitationCode.IncrementUse(req.UserID)
	if err := ctrl.db.Save(&invitationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to use invitation code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Invitation code used successfully",
		"code": invitationCode.Code,
		"codeType": invitationCode.CodeType,
		"ownerId": invitationCode.OwnerID,
	})
}

// DisableInvitationCode 禁用邀请码
func (ctrl *InvitationController) DisableInvitationCode(c *gin.Context) {
	code := c.Param("code")
	userID := c.GetString("userId")
	currentRole := c.GetString("currentRole")

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

	// 权限检查：只有创建者或超级管理员可以禁用
	if currentRole != "super_admin" && invitationCode.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
		return
	}

	// 禁用邀请码
	invitationCode.Status = "disabled"
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

// canCreateInvitationCode 检查是否可以创建指定类型的邀请码
func (ctrl *InvitationController) canCreateInvitationCode(currentRole, codeType, ownerType string) bool {
	switch currentRole {
	case "super_admin":
		// 超级管理员可以创建所有类型的邀请码
		return true
	case "service_provider_admin":
		// 服务商管理员可以创建商家、达人邀请码
		return codeType == "MERCHANT" || codeType == "CREATOR" || codeType == "STAFF"
	case "merchant_admin":
		// 商家管理员可以创建员工邀请码
		return codeType == "STAFF"
	default:
		return false
	}
}
