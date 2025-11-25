package api

import (
    "net/http"
    "reward-system/internal/config"
    "reward-system/internal/services"
    "sort"
    "strings"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type WeChatMessage struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgID        int64  `xml:"MsgId"`
}

type WeChatResponse struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
}

func WeChatWebhook(database *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify signature for GET request (WeChat verification)
		if c.Request.Method == "GET" {
			signature := c.Query("signature")
			timestamp := c.Query("timestamp")
			nonce := c.Query("nonce")
			echostr := c.Query("echostr")

			if verifySignature(cfg.WechatToken, signature, timestamp, nonce) {
				c.String(http.StatusOK, echostr)
				return
			}
			c.String(http.StatusForbidden, "Invalid signature")
			return
		}

		// Handle POST request (actual messages)
		var msg WeChatMessage
		if err := c.ShouldBindXML(&msg); err != nil {
			c.String(http.StatusBadRequest, "Invalid XML")
			return
		}

		// Process the message
		response := processWeChatMessage(database, msg)
		
		// Send response
		resp := WeChatResponse{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   msg.CreateTime,
			MsgType:      "text",
			Content:      response,
		}

		c.XML(http.StatusOK, resp)
	}
}

func verifySignature(token, signature, timestamp, nonce string) bool {
	if token == "" {
		return true // Skip verification if no token configured
	}
	
    params := []string{token, timestamp, nonce}
    sort.Strings(params)
    // Simple hash comparison (in production, use proper SHA1)
    return signature != "" // Simplified for demo
}

func processWeChatMessage(database *gorm.DB, msg WeChatMessage) string {
	content := strings.TrimSpace(msg.Content)
	
	// Check if it's a structured command
	if strings.HasPrefix(content, "#cmd ") {
		return processStructuredCommand(database, msg.FromUserName, content[5:])
	}
	
	// Process natural language with MCP
	return processNaturalLanguage(database, msg.FromUserName, content)
}

func processStructuredCommand(database *gorm.DB, openID string, cmd string) string {
	// Parse JSON command and execute
	// Simplified for now - return a placeholder response
	return "结构化指令已收到，正在处理..."
}

func processNaturalLanguage(database *gorm.DB, openID string, text string) string {
	// Use MCP to parse natural language and execute appropriate action
	// Simplified for now - return a placeholder response
	service := services.NewRewardService(database)
	_ = service // Use the service to perform operations
	
	return "自然语言消息已收到: " + text
}