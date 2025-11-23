package api

import (
	"net/http"
	"reward-system/internal/db"
	"reward-system/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateRewardType(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FamilyID  uint64 `json:"family_id" binding:"required"`
			Name      string `json:"name" binding:"required"`
			UnitKind  string `json:"unit_kind" binding:"required,oneof=money time points custom"`
			UnitLabel string `json:"unit_label"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}

		rewardType := &db.RewardType{
			FamilyID:  req.FamilyID,
			Name:      req.Name,
			UnitKind:  req.UnitKind,
			UnitLabel: req.UnitLabel,
		}

		if err := database.Create(rewardType).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create reward type"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": rewardType})
	}
}

func GrantReward(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FamilyID         uint64 `json:"family_id" binding:"required"`
			ChildID          uint64 `json:"child_id" binding:"required"`
			RewardTypeID     uint64 `json:"reward_type_id" binding:"required"`
			Value            int64  `json:"value" binding:"required"`
			Note             string `json:"note"`
			IdempotencyKey   string `json:"idempotency_key"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}

		service := services.NewRewardService(database)
		result, err := service.GrantReward(req.FamilyID, req.ChildID, req.RewardTypeID, req.Value, req.Note, req.IdempotencyKey)
		
		if err != nil {
			if err.Error() == "insufficient balance" {
				c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "Insufficient balance"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
	}
}

func SpendReward(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FamilyID         uint64 `json:"family_id" binding:"required"`
			ChildID          uint64 `json:"child_id" binding:"required"`
			RewardTypeID     uint64 `json:"reward_type_id" binding:"required"`
			Value            int64  `json:"value" binding:"required"`
			Note             string `json:"note"`
			IdempotencyKey   string `json:"idempotency_key"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}

		service := services.NewRewardService(database)
		result, err := service.SpendReward(req.FamilyID, req.ChildID, req.RewardTypeID, req.Value, req.Note, req.IdempotencyKey)
		
		if err != nil {
			if err.Error() == "insufficient balance" {
				c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "Insufficient balance"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
	}
}

func GetBalance(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		familyID := c.Query("family_id")
		childID := c.Query("child_id")
		rewardTypeID := c.Query("reward_type_id")

		if familyID == "" || childID == "" || rewardTypeID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Missing required parameters"})
			return
		}

		service := services.NewRewardService(database)
		balance, err := service.GetBalance(parseUint(familyID), parseUint(childID), parseUint(rewardTypeID))
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"balance": balance}})
	}
}

func ListTransactions(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		familyID := c.Query("family_id")
		childID := c.Query("child_id")
		rewardTypeID := c.Query("reward_type_id")
		limit := c.DefaultQuery("limit", "20")
		beforeID := c.Query("before_id")

		if familyID == "" || childID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Missing required parameters"})
			return
		}

		service := services.NewRewardService(database)
		transactions, err := service.ListTransactions(parseUint(familyID), parseUint(childID), parseUint(rewardTypeID), parseInt(limit), parseUint(beforeID))
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": transactions})
	}
}

func AdjustTransaction(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionID := c.Param("id")
		
		var req struct {
			NewValue *int64  `json:"new_value"`
			NewNote  *string `json:"new_note"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}

		service := services.NewRewardService(database)
		err := service.AdjustTransaction(parseUint(transactionID), req.NewValue, req.NewNote)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"id": transactionID}})
	}
}

func parseUint(s string) uint64 {
	if s == "" {
		return 0
	}
	var result uint64
	fmt.Sscanf(s, "%d", &result)
	return result
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}