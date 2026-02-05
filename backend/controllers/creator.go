package controllers

import (
	"net/http"
	"pr-business/models"
	"pr-business/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatorController struct {
	db *gorm.DB
}

func NewCreatorController(db *gorm.DB) *CreatorController {
	return &CreatorController{db: db}
}

// UpdateCreatorRequest 更新达人信息请求
type UpdateCreatorRequest struct {
	Level            string `json:"level" binding:"omitempty,oneof=UGC KOC INF KOL"`
	FollowersCount   int    `json:"followersCount" binding:"omitempty,min=0"`
	WechatOpenID     string `json:"wechatOpenId" binding:"omitempty,max=100"`
	WechatNickname   string `json:"wechatNickname" binding:"omitempty,max=100"`
	WechatAvatar     string `json:"wechatAvatar" binding:"omitempty,max=500"`
	Status           string `json:"status" binding:"omitempty,oneof=active banned inactive"`
}

// GetCreators 获取达人列表
// @Summary 获取达人列表
// @Description 获取达人列表，支持过滤
// @Tags 达人管理
// @Accept json
// @Produce json
// @Param level query string false "等级过滤"
// @Param status query string false "状态过滤"
// @Param inviter_id query string false "邀请人ID过滤"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {array} models.Creator
// @Router /api/v1/creators [get]
func (ctrl *CreatorController) GetCreators(c *gin.Context) {
	var creators []models.Creator
	query := ctrl.db.Model(&models.Creator{})

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 等级过滤
	if level := c.Query("level"); level != "" {
		query = query.Where("level = ?", level)
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 邀请人过滤
	if inviterID := c.Query("inviter_id"); inviterID != "" {
		// 如果是服务商员工，只能看到自己邀请的达人
		if utils.HasRole(user, "provider_staff") {
			if inviterID != user.ID {
				c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看其他员工邀请的达人"})
				return
			}
		}
		query = query.Where("inviter_id = ?", inviterID)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	if err := query.Preload("User").Preload("Inviter").Offset(offset).Limit(pageSize).Find(&creators).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取达人列表失败"})
		return
	}

	// 获取总数
	var total int64
	ctrl.db.Model(&models.Creator{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"data":   creators,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	})
}

// GetCreator 获取达人详情
// @Summary 获取达人详情
// @Description 根据ID获取达人详情
// @Tags 达人管理
// @Accept json
// @Produce json
// @Param id path string true "达人ID"
// @Success 200 {object} models.Creator
// @Router /api/v1/creators/{id} [get]
func (ctrl *CreatorController) GetCreator(c *gin.Context) {
	id := c.Param("id")
	var creator models.Creator

	if err := ctrl.db.Where("id = ?", id).Preload("User").Preload("Inviter").First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "达人不存在"})
		return
	}

	c.JSON(http.StatusOK, creator)
}

// UpdateCreator 更新达人信息
// @Summary 更新达人信息
// @Description 达人可以更新自己的信息，管理员可以更新达人信息
// @Tags 达人管理
// @Accept json
// @Produce json
// @Param id path string true "达人ID"
// @Param request body UpdateCreatorRequest true "更新达人请求"
// @Success 200 {object} models.Creator
// @Router /api/v1/creators/{id} [put]
func (ctrl *CreatorController) UpdateCreator(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCreatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var creator models.Creator
	if err := ctrl.db.Where("id = ?", id).First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "达人不存在"})
		return
	}

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查：达人本人或管理员可以更新
	isCreator := user.ID == creator.UserID.String()
	isAdmin := utils.HasRole(user, "super_admin") || utils.HasRole(user, "provider_admin") || utils.HasRole(user, "provider_staff")

	if !isCreator && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新达人信息"})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})

	// 达人本人可以更新的字段
	if isCreator {
		if req.WechatOpenID != "" {
			updates["wechat_open_id"] = req.WechatOpenID
		}
		if req.WechatNickname != "" {
			updates["wechat_nickname"] = req.WechatNickname
		}
		if req.WechatAvatar != "" {
			updates["wechat_avatar"] = req.WechatAvatar
		}
	}

	// 管理员可以更新的字段
	if isAdmin {
		if req.Level != "" {
			updates["level"] = req.Level
		}
		if req.FollowersCount >= 0 {
			updates["followers_count"] = req.FollowersCount
		}
		if req.Status != "" {
			updates["status"] = req.Status
		}
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有可更新的字段"})
		return
	}

	if err := ctrl.db.Model(&creator).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新达人失败"})
		return
	}

	c.JSON(http.StatusOK, creator)
}

// GetCreatorLevelStats 获取达人等级统计
// @Summary 获取达人等级统计
// @Description 获取各等级达人数量统计
// @Tags 达人管理
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/creators/stats/level [get]
func (ctrl *CreatorController) GetCreatorLevelStats(c *gin.Context) {
	var stats []struct {
		Level string
		Count int64
	}

	if err := ctrl.db.Model(&models.Creator{}).
		Select("level, count(*) as count").
		Group("level").
		Scan(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GetCreatorInviterRelationship 获取达人邀请关系
// @Summary 获取达人邀请关系
// @Description 获取达人及其邀请人的关系信息
// @Tags 达人管理
// @Accept json
// @Produce json
// @Param id path string true "达人ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/creators/{id}/inviter [get]
func (ctrl *CreatorController) GetCreatorInviterRelationship(c *gin.Context) {
	id := c.Param("id")
	var creator models.Creator

	if err := ctrl.db.Where("id = ?", id).First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "达人不存在"})
		return
	}

	// 预加载邀请人信息
	if err := ctrl.db.Model(&creator).Preload("Inviter").First(&creator).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取邀请关系失败"})
		return
	}

	result := gin.H{
		"creator_id": creator.ID,
		"user_id":    creator.UserID,
		"inviter":    nil,
		"inviter_type": creator.InviterType,
		"relationship_broken": creator.InviterRelationshipBroken,
	}

	if creator.Inviter != nil {
		result["inviter"] = gin.H{
			"user_id":  creator.Inviter.ID,
			"nickname": creator.Inviter.Nickname,
			"roles":    creator.Inviter.Roles,
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetMyCreatorProfile 获取当前用户的达人资料
// @Summary 获取当前用户的达人资料
// @Description 达人获取自己的资料信息
// @Tags 达人管理
// @Accept json
// @Produce json
// @Success 200 {object} models.Creator
// @Router /api/v1/creator/me [get]
func (ctrl *CreatorController) GetMyCreatorProfile(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 查找主达人资料
	var creator models.Creator
	if err := ctrl.db.Where("user_id = ? AND is_primary = ?", user.ID, true).Preload("User").First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "您不是达人"})
		return
	}

	c.JSON(http.StatusOK, creator)
}

// UpdateMyCreatorProfile 更新当前用户的达人资料
// @Summary 更新当前用户的达人资料
// @Description 达人更新自己的资料信息
// @Tags 达人管理
// @Accept json
// @Produce json
// @Param request body UpdateCreatorRequest true "更新达人请求"
// @Success 200 {object} models.Creator
// @Router /api/v1/creator/me [put]
func (ctrl *CreatorController) UpdateMyCreatorProfile(c *gin.Context) {
	var req UpdateCreatorRequest
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

	// 查找主达人资料
	var creator models.Creator
	if err := ctrl.db.Where("user_id = ? AND is_primary = ?", user.ID, true).First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "您不是达人"})
		return
	}

	// 更新字段（达人本人可以更新的字段）
	updates := make(map[string]interface{})
	if req.WechatOpenID != "" {
		updates["wechat_open_id"] = req.WechatOpenID
	}
	if req.WechatNickname != "" {
		updates["wechat_nickname"] = req.WechatNickname
	}
	if req.WechatAvatar != "" {
		updates["wechat_avatar"] = req.WechatAvatar
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有可更新的字段"})
		return
	}

	if err := ctrl.db.Model(&creator).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新达人资料失败"})
		return
	}

	c.JSON(http.StatusOK, creator)
}

// BreakInviterRelationship 解除邀请关系
// @Summary 解除邀请关系
// @Description 达人解除与邀请人的关系
// @Tags 达人管理
// @Accept json
// @Produce json
// @Param id path string true "达人ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/creators/{id}/break-relationship [post]
func (ctrl *CreatorController) BreakInviterRelationship(c *gin.Context) {
	id := c.Param("id")

	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	var creator models.Creator
	if err := ctrl.db.Where("id = ?", id).First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "达人不存在"})
		return
	}

	// 权限检查：只有达人本人或管理员可以解除关系
	isCreator := user.ID == creator.UserID.String()
	isAdmin := utils.HasRole(user, "super_admin")

	if !isCreator && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限解除邀请关系"})
		return
	}

	// 更新关系状态
	if err := ctrl.db.Model(&creator).Update("inviter_relationship_broken", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解除关系失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已解除邀请关系"})
}
