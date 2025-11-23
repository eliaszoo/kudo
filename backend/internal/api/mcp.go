package api

import (
	"net/http"
	"reward-system/internal/db"
	"reward-system/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MCPToolsHandler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Tool    string                 `json:"tool" binding:"required"`
			Params  map[string]interface{} `json:"params" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}

		service := services.NewRewardService(database)
		
		switch req.Tool {
		case "create_reward_type":
			handleCreateRewardType(c, service, req.Params)
		case "grant_reward":
			handleGrantReward(c, service, req.Params)
		case "spend_reward":
			handleSpendReward(c, service, req.Params)
		case "query_balance":
			handleQueryBalance(c, service, req.Params)
		case "list_transactions":
			handleListTransactions(c, service, req.Params)
		case "adjust_transaction":
			handleAdjustTransaction(c, service, req.Params)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Unknown tool"})
		}
	}
}

func handleCreateRewardType(c *gin.Context, service *services.RewardService, params map[string]interface{}) {
	familyID := uint64(params["family_id"].(float64))
	name := params["name"].(string)
	unitKind := params["unit_kind"].(string)
	unitLabel := ""
	if val, ok := params["unit_label"].(string); ok {
		unitLabel = val
	}

	rewardType := &db.RewardType{
		FamilyID:  familyID,
		Name:      name,
		UnitKind:  unitKind,
		UnitLabel: unitLabel,
	}

	if err := service.CreateRewardType(rewardType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"reward_type_id": rewardType.ID}})
}

func handleGrantReward(c *gin.Context, service *services.RewardService, params map[string]interface{}) {
	familyID := uint64(params["family_id"].(float64))
	childID := uint64(params["child_id"].(float64))
	rewardTypeID := uint64(params["reward_type_id"].(float64))
	value := int64(params["value"].(float64))
	
	note := ""
	if val, ok := params["note"].(string); ok {
		note = val
	}
	
	idempotencyKey := ""
	if val, ok := params["idempotency_key"].(string); ok {
		idempotencyKey = val
	}

	result, err := service.GrantReward(familyID, childID, rewardTypeID, value, note, idempotencyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

func handleSpendReward(c *gin.Context, service *services.RewardService, params map[string]interface{}) {
	familyID := uint64(params["family_id"].(float64))
	childID := uint64(params["child_id"].(float64))
	rewardTypeID := uint64(params["reward_type_id"].(float64))
	value := int64(params["value"].(float64))
	
	note := ""
	if val, ok := params["note"].(string); ok {
		note = val
	}
	
	idempotencyKey := ""
	if val, ok := params["idempotency_key"].(string); ok {
		idempotencyKey = val
	}

	result, err := service.SpendReward(familyID, childID, rewardTypeID, value, note, idempotencyKey)
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

func handleQueryBalance(c *gin.Context, service *services.RewardService, params map[string]interface{}) {
	familyID := uint64(params["family_id"].(float64))
	childID := uint64(params["child_id"].(float64))
	rewardTypeID := uint64(params["reward_type_id"].(float64))

	balance, err := service.GetBalance(familyID, childID, rewardTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"balance": balance}})
}

func handleListTransactions(c *gin.Context, service *services.RewardService, params map[string]interface{}) {
	familyID := uint64(params["family_id"].(float64))
	childID := uint64(params["child_id"].(float64))
	
	rewardTypeID := uint64(0)
	if val, ok := params["reward_type_id"].(float64); ok {
		rewardTypeID = uint64(val)
	}
	
	limit := 20
	if val, ok := params["limit"].(float64); ok {
		limit = int(val)
	}
	
	beforeID := uint64(0)
	if val, ok := params["before_id"].(float64); ok {
		beforeID = uint64(val)
	}

	transactions, err := service.ListTransactions(familyID, childID, rewardTypeID, limit, beforeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": transactions})
}

func handleAdjustTransaction(c *gin.Context, service *services.RewardService, params map[string]interface{}) {
	transactionID := uint64(params["transaction_id"].(float64))
	
	var newValue *int64
	if val, ok := params["new_value"].(float64); ok {
		v := int64(val)
		newValue = &v
	}
	
	var newNote *string
	if val, ok := params["new_note"].(string); ok {
		newNote = &val
	}

	err := service.AdjustTransaction(transactionID, newValue, newNote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"transaction_id": transactionID}})
}