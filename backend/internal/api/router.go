package api

import (
    "reward-system/internal/config"

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
		v1.GET("/reward_types", ListRewardTypes(database))
		v1.PATCH("/reward_types/:id", UpdateRewardType(database))

		// Families and users
		v1.GET("/families", ListFamilies(database))
		v1.POST("/families", CreateFamily(database))
		v1.GET("/users", ListUsers(database))
		v1.POST("/users", CreateUser(database))
		v1.DELETE("/users/:id", DeleteUser(database))
		
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