package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"pr-business/models"
	"pr-business/services"
	"pr-business/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RechargeOrderController struct {
	db                *gorm.DB
	permissionService  *services.AccountPermissionService
	validatorService   *services.ValidatorService
}

func NewRechargeOrderController(db *gorm.DB) *RechargeOrderController {
	return &RechargeOrderController{
		db:                db,
		permissionService:  services.NewAccountPermissionService(db),
		validatorService:   services.NewValidatorService(db),
	}
}

// CreateRechargeOrderRequest 创建充值订单请求
type CreateRechargeOrderRequest struct {
	Amount        int    `json:"amount" binding:"required,min=1,max=1000000"`
	PaymentMethod string `json:"paymentMethod" binding:"required,oneof=alipay wechat bank"`
	PaymentProof string `json:"paymentProof" binding:"required"`
}

// CreateRechargeOrder 用户提交充值订单
func (ctrl *RechargeOrderController) CreateRechargeOrder(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	var req CreateRechargeOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确定用户的积分账户
	var accountType models.OwnerType
	var ownerID uuid.UUID

	// 只有商家管理员可以充值
	if !utils.IsMerchantAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有商家管理员可以充值"})
		return
	}

	// 获取商家信息
	var merchant models.Merchant
	if err := ctrl.db.Where("admin_id = ?", user.ID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商家信息不存在"})
		return
	}

	accountType = models.OwnerTypeOrgMerchant
	ownerID = merchant.ID

	// 获取或创建积分账户
	userID, err := uuid.Parse(user.AuthCenterUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID格式错误"})
		return
	}

	account, err := ctrl.findOrCreateAccount(ownerID, accountType, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取积分账户失败"})
		return
	}

	// 创建充值订单
	order := models.RechargeOrder{
		UserID:        user.ID,
		AccountID:     account.ID,
		Amount:         req.Amount,
		PaymentMethod:  req.PaymentMethod,
		PaymentProof:   req.PaymentProof,
		Status:         models.RechargeOrderStatusPending,
	}

	if err := ctrl.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建充值订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "充值订单已提交，等待管理员审核",
		"order": gin.H{
			"id":            order.ID,
			"amount":        order.Amount,
			"paymentMethod":  order.PaymentMethod,
			"status":        order.Status,
			"createdAt":     order.CreatedAt,
		},
	})
}

// GetRechargeOrders 获取充值订单列表
func (ctrl *RechargeOrderController) GetRechargeOrders(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 分页参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	status := c.Query("status")

	var orders []models.RechargeOrder
	var total int64

	query := ctrl.db.Model(&models.RechargeOrder{})

	// 超管可以查看所有订单，商家只能查看自己的
	if !utils.IsSuperAdmin(user) {
		query = query.Where("user_id = ?", user.ID)
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	offset := (pageInt - 1) * pageSizeInt
	if err := query.Preload("Account").Order("created_at DESC").Offset(offset).Limit(pageSizeInt).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"page":   pageInt,
		"page_size": pageSizeInt,
	})
}

// AuditRechargeOrderRequest 审核充值订单请求
type AuditRechargeOrderRequest struct {
	Approved      bool   `json:"approved" binding:"required"`
	RejectionNote string `json:"rejectionNote"`
}

// AuditRechargeOrder 超管审核充值订单
func (ctrl *RechargeOrderController) AuditRechargeOrder(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	user := currentUser.(*models.User)

	// 只有超管可以审核
	if !utils.IsSuperAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权审核充值订单"})
		return
	}

	id := c.Param("id")

	var req AuditRechargeOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.RechargeOrder
	if err := ctrl.db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "充值订单不存在"})
		return
	}

	// 检查状态
	if !order.CanAudit() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该订单已被审核"})
		return
	}

	// 获取积分账户
	var account models.CreditAccount
	if err := ctrl.db.Where("id = ?", order.AccountID).First(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "积分账户不存在"})
		return
	}

	err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		if req.Approved {
			// 通过：更新状态为approved
			order.Status = models.RechargeOrderStatusApproved
			order.AuditedBy = &user.ID
			order.AuditedAt = &now

			if err := tx.Save(&order).Error; err != nil {
				return err
			}

			// 充值入账
			balanceBefore := account.Balance
			account.Balance += order.Amount
			if err := tx.Save(&account).Error; err != nil {
				return err
			}

			// 记录充值流水
			transaction := models.CreditTransaction{
				AccountID:        account.ID,
				Type:             models.TransactionRecharge,
				Amount:           order.Amount,
				BalanceBefore:     balanceBefore,
				BalanceAfter:      account.Balance,
				Description:       fmt.Sprintf("充值订单：%s", order.ID),
			}
			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}

			// 更新订单状态为completed
			order.Status = models.RechargeOrderStatusCompleted
			order.ProcessedAt = &now
			if err := tx.Save(&order).Error; err != nil {
				return err
			}

		} else {
			// 拒绝
			order.Status = models.RechargeOrderStatusRejected
			order.RejectionNote = req.RejectionNote
			order.AuditedBy = &user.ID
			order.AuditedAt = &now

			if err := tx.Save(&order).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核失败"})
		return
	}

	statusText := "已通过"
	if !req.Approved {
		statusText = "已拒绝"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("充值订单%s", statusText),
		"order": gin.H{
			"id":       order.ID,
			"status":   order.Status,
			"auditedAt": order.AuditedAt,
		},
	})
}

// findOrCreateAccount 查找或创建账户（带权限）
func (ctrl *RechargeOrderController) findOrCreateAccount(ownerID uuid.UUID, ownerType models.OwnerType, userID uuid.UUID) (*models.CreditAccount, error) {
	var account models.CreditAccount
	err := ctrl.db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).First(&account).Error
	if err == nil {
		return &account, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		account := &models.CreditAccount{
			OwnerID:       ownerID,
			OwnerType:     ownerType,
			Balance:       0,
			FrozenBalance: 0,
		}

		if err := ctrl.db.Create(account).Error; err != nil {
			return nil, err
		}

		// 授予权限
		if ownerType == models.OwnerTypeOrgMerchant {
			if err := ctrl.permissionService.GrantMerchantAccountPermission(userID, account.ID, true); err != nil {
				return nil, err
			}
		}

		return account, nil
	}

	return nil, err
}
