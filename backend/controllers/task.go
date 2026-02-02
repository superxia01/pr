package controllers

import (
	"net/http"
	"pr-business/models"
	"pr-business/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskController struct {
	db *gorm.DB
}

func NewTaskController(db *gorm.DB) *TaskController {
	return &TaskController{db: db}
}

// AcceptTaskRequest 接任务请求
type AcceptTaskRequest struct {
	Platform string `json:"platform" binding:"required"`
}

// SubmitTaskRequest 提交任务请求
type SubmitTaskRequest struct {
	PlatformURL  string   `json:"platformUrl" binding:"required,url"`
	Screenshots  string   `json:"screenshots"` // JSONB string
	Notes        string   `json:"notes"`
}

// AuditTaskRequest 审核任务请求
type AuditTaskRequest struct {
	Action     string `json:"action" binding:"required,oneof=approve reject"`
	AuditNote  string `json:"auditNote"`
}

// GetTasks 获取任务名额列表
// @Summary 获取任务名额列表
// @Description 根据营销活动ID获取任务名额列表
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param campaign_id query string false "营销活动ID"
// @Param status query string false "状态过滤"
// @Param creator_id query string false "达人ID过滤"
// @Success 200 {array} models.Task
// @Router /api/v1/tasks [get]
func (ctrl *TaskController) GetTasks(c *gin.Context) {
	var tasks []models.Task
	query := ctrl.db.Model(&models.Task{})

	// 营销活动过滤
	if campaignID := c.Query("campaign_id"); campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 达人过滤
	if creatorID := c.Query("creator_id"); creatorID != "" {
		query = query.Where("creator_id = ?", creatorID)
	}

	if err := query.Preload("Campaign").Preload("Creator").Preload("Auditor").Order("created_at DESC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务列表失败"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask 获取任务详情
// @Summary 获取任务详情
// @Description 根据ID获取任务详情
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} models.Task
// @Router /api/v1/tasks/{id} [get]
func (ctrl *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := ctrl.db.Where("id = ?", id).Preload("Campaign").Preload("Creator").Preload("Auditor").First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// AcceptTask 达人接任务
// @Summary 达人接任务
// @Description 达人接受一个开放的任务名额
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param request body AcceptTaskRequest true "接任务请求"
// @Success 200 {object} models.Task
// @Router /api/v1/tasks/{id}/accept [post]
func (ctrl *TaskController) AcceptTask(c *gin.Context) {
	id := c.Param("id")
	var req AcceptTaskRequest
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

	// 权限检查：只有达人可以接任务
	if !utils.HasRole(user, "CREATOR") {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有达人可以接任务"})
		return
	}

	// 获取任务
	var task models.Task
	if err := ctrl.db.Where("id = ?", id).Preload("Campaign").First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	// 检查任务状态
	if task.Status != models.TaskStatusOpen {
		c.JSON(http.StatusForbidden, gin.H{"error": "任务不可接"})
		return
	}

	// 检查营销活动状态
	if task.Campaign.Status != models.CampaignStatusOpen {
		c.JSON(http.StatusForbidden, gin.H{"error": "营销活动未开放"})
		return
	}

	// 检查截止时间
	if time.Now().After(task.Campaign.TaskDeadline) {
		c.JSON(http.StatusForbidden, gin.H{"error": "任务已过截止时间"})
		return
	}

	// 获取达人信息
	var creator models.Creator
	if err := ctrl.db.Where("user_id = ? AND is_primary = ?", user.ID, true).First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "达人信息不存在"})
		return
	}

	// 检查达人是否已经接了该营销活动的任务
	var existingTask models.Task
	if err := ctrl.db.Where("campaign_id = ? AND creator_id = ?", task.CampaignID, creator.ID).First(&existingTask).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "您已经接了该营销活动的任务"})
		return
	}

	// 开始事务
	tx := ctrl.db.Begin()

	// 更新任务状态
	now := time.Now()
	task.Status = models.TaskStatusAssigned
	task.CreatorID = &creator.ID
	task.AssignedAt = &now
	task.Platform = req.Platform
	task.InviterID = creator.InviterID
	task.InviterType = creator.InviterType

	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "接任务失败"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, task)
}

// SubmitTask 达人提交任务
// @Summary 达人提交任务
// @Description 达人提交完成任务的平台链接和截图
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param request body SubmitTaskRequest true "提交任务请求"
// @Success 200 {object} models.Task
// @Router /api/v1/tasks/{id}/submit [post]
func (ctrl *TaskController) SubmitTask(c *gin.Context) {
	id := c.Param("id")
	var req SubmitTaskRequest
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

	// 获取任务
	var task models.Task
	if err := ctrl.db.Where("id = ?", id).Preload("Campaign").First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	// 权限检查：只有任务所属的达人可以提交
	if task.CreatorID == nil || task.CreatorID.String() != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限提交此任务"})
		return
	}

	// 检查任务状态
	if task.Status != models.TaskStatusAssigned {
		c.JSON(http.StatusForbidden, gin.H{"error": "任务状态不允许提交"})
		return
	}

	// 检查截止时间
	if time.Now().After(task.Campaign.SubmissionDeadline) {
		c.JSON(http.StatusForbidden, gin.H{"error": "已过提交截止时间"})
		return
	}

	// 更新任务状态
	now := time.Now()
	task.Status = models.TaskStatusSubmitted
	task.PlatformURL = req.PlatformURL
	task.Screenshots = req.Screenshots
	task.Notes = req.Notes
	task.SubmittedAt = &now
	task.Version += 1

	if err := ctrl.db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交任务失败"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// AuditTask 服务商审核任务
// @Summary 服务商审核任务
// @Description 服务商审核达人提交的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param request body AuditTaskRequest true "审核任务请求"
// @Success 200 {object} models.Task
// @Router /api/v1/tasks/{id}/audit [post]
func (ctrl *TaskController) AuditTask(c *gin.Context) {
	id := c.Param("id")
	var req AuditTaskRequest
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

	// 权限检查：服务商员工可以审核
	if !utils.HasRole(user, "SP_ADMIN") && !utils.HasRole(user, "SP_STAFF") && !utils.HasRole(user, "SUPER_ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限审核任务"})
		return
	}

	// 获取任务
	var task models.Task
	if err := ctrl.db.Where("id = ?", id).Preload("Campaign").First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	// 检查任务状态
	if task.Status != models.TaskStatusSubmitted {
		c.JSON(http.StatusForbidden, gin.H{"error": "任务状态不允许审核"})
		return
	}

	// 更新任务状态
	now := time.Now()
	if req.Action == "approve" {
		task.Status = models.TaskStatusApproved
	} else if req.Action == "reject" {
		task.Status = models.TaskStatusRejected
		// 拒绝后，达人可以重新提交，状态转回 ASSIGNED
		task.Status = models.TaskStatusRejected
	}

	// 将 user.ID (string) 转换为 uuid.UUID
	auditorID, err := uuid.Parse(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
		return
	}
	task.AuditedBy = &auditorID
	task.AuditedAt = &now
	task.AuditNote = req.AuditNote
	task.Version += 1

	if err := ctrl.db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核失败"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetMyTasks 获取当前用户的任务列表
// @Summary 获取当前用户的任务列表
// @Description 达人获取自己的任务列表
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param status query string false "状态过滤"
// @Success 200 {array} models.Task
// @Router /api/v1/tasks/my [get]
func (ctrl *TaskController) GetMyTasks(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 获取达人信息
	var creator models.Creator
	if err := ctrl.db.Where("user_id = ? AND is_primary = ?", user.ID, true).First(&creator).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "达人信息不存在"})
		return
	}

	var tasks []models.Task
	query := ctrl.db.Where("creator_id = ?", creator.ID)

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Preload("Campaign").Order("created_at DESC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务列表失败"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTasksForReview 获取待审核任务列表（服务商）
// @Summary 获取待审核任务列表
// @Description 服务商获取待审核的任务列表
// @Tags 任务管理
// @Accept json
// @Produce json
// @Success 200 {array} models.Task
// @Router /api/v1/tasks/pending-review [get]
func (ctrl *TaskController) GetTasksForReview(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查
	if !utils.HasRole(user, "SP_ADMIN") && !utils.HasRole(user, "SP_STAFF") && !utils.HasRole(user, "SUPER_ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看待审核任务"})
		return
	}

	var tasks []models.Task
	query := ctrl.db.Where("status = ?", models.TaskStatusSubmitted)

	// 如果是服务商员工，只能看到自己负责的营销活动的任务
	if utils.HasRole(user, "SP_STAFF") && !utils.HasRole(user, "SUPER_ADMIN") {
		// 获取员工所属的服务商
		var providerStaff models.ServiceProviderStaff
		if err := ctrl.db.Where("user_id = ?", user.ID).First(&providerStaff).Error; err == nil {
			// 获取该服务商下的所有营销活动
			var campaigns []models.Campaign
			ctrl.db.Where("provider_id = ?", providerStaff.ProviderID).Find(&campaigns)

			campaignIDs := make([]string, len(campaigns))
			for i, campaign := range campaigns {
				campaignIDs[i] = campaign.ID.String()
			}

			if len(campaignIDs) > 0 {
				query = query.Where("campaign_id IN ?", campaignIDs)
			} else {
				// 如果没有营销活动，返回空列表
				c.JSON(http.StatusOK, []models.Task{})
				return
			}
		}
	}

	if err := query.Preload("Campaign").Preload("Creator").Order("submitted_at ASC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取待审核任务列表失败"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTaskHall 获取任务大厅（开放的任务）
// @Summary 获取任务大厅
// @Description 达人查看可接的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tasks/hall [get]
func (ctrl *TaskController) GetTaskHall(c *gin.Context) {
	// 获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 权限检查：只有达人可以查看任务大厅
	if !utils.HasRole(user, "CREATOR") {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有达人可以查看任务大厅"})
		return
	}

	// 分页
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if num, err := strconv.Atoi(p); err == nil {
			page = num
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if num, err := strconv.Atoi(ps); err == nil {
			pageSize = num
		}
	}

	offset := (page - 1) * pageSize

	var tasks []models.Task
	var total int64

	// 查询开放中的任务
	query := ctrl.db.Model(&models.Task{}).Where("status = ?", models.TaskStatusOpen)

	// 统计总数
	query.Count(&total)

	// 获取任务列表，预加载营销活动信息
	if err := query.Preload("Campaign.Merchant").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      tasks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
