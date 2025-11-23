package api

import (
	"reward-system/internal/config"
	"reward-system/internal/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(database *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	
	// Middleware
	r.Use(CORSMiddleware())
	r.Use(AuthMiddleware(cfg))
	
	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Reward types
		v1.POST("/reward_types", CreateRewardType(database))
		
		// Rewards
		v1.POST("/rewards/grant", GrantReward(database))
		v1.POST("/rewards/spend", SpendReward(database))
		
		// Balances
		v1.GET("/balances", GetBalance(database))
		
		// Transactions
		v1.GET("/transactions", ListTransactions(database))
		v1.POST("/transactions/:id/adjust", AdjustTransaction(database))
		
		// WeChat webhook
		v1.POST("/wechat", WeChatWebhook(database, cfg))
		
		// MCP tools
		v1.POST("/mcp/tools", MCPToolsHandler(database))
	}
	
	return r
}