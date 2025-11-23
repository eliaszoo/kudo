package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reward-system/internal/db"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestAPI(t *testing.T) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	if err := database.AutoMigrate(
		&db.Family{},
		&db.User{},
		&db.RewardType{},
		&db.Account{},
		&db.Transaction{},
		&db.AuditLog{},
	); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create test config
	cfg := &config.Config{
		APIToken: "test-token",
		Port:     "8080",
	}

	// Setup router
	router := SetupRouter(database, cfg)

	return router, database
}

func TestCreateRewardType(t *testing.T) {
	router, db := setupTestAPI(t)

	// Create test family
	family := &db.Family{Name: "Test Family"}
	if err := db.Create(family).Error; err != nil {
		t.Fatalf("Failed to create test family: %v", err)
	}

	// Test request body
	body := map[string]interface{}{
		"family_id":  family.ID,
		"name":       "Test Reward",
		"unit_kind":  "money",
		"unit_label": "元",
	}

	jsonBody, _ := json.Marshal(body)

	// Create request
	req, err := http.NewRequest("POST", "/api/v1/reward_types", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["code"] != float64(0) {
		t.Errorf("Expected code 0, got %v", response["code"])
	}
}

func TestGrantReward(t *testing.T) {
	router, db := setupTestAPI(t)

	// Create test data
	family := &db.Family{Name: "Test Family"}
	if err := db.Create(family).Error; err != nil {
		t.Fatalf("Failed to create test family: %v", err)
	}

	child := &db.User{
		FamilyID:    family.ID,
		Role:        "child",
		DisplayName: "Test Child",
	}
	if err := db.Create(child).Error; err != nil {
		t.Fatalf("Failed to create test child: %v", err)
	}

	rewardType := &db.RewardType{
		FamilyID: family.ID,
		Name:     "Test Reward",
		UnitKind: "money",
	}
	if err := db.Create(rewardType).Error; err != nil {
		t.Fatalf("Failed to create test reward type: %v", err)
	}

	// Test request body
	body := map[string]interface{}{
		"family_id":        family.ID,
		"child_id":         child.ID,
		"reward_type_id":   rewardType.ID,
		"value":            1000,
		"note":             "Test grant",
		"idempotency_key": "test-key-1",
	}

	jsonBody, _ := json.Marshal(body)

	// Create request
	req, err := http.NewRequest("POST", "/api/v1/rewards/grant", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["code"] != float64(0) {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	data := response["data"].(map[string]interface{})
	if data["transaction_id"] == nil {
		t.Error("Expected transaction_id in response")
	}

	if data["new_balance"] != float64(1000) {
		t.Errorf("Expected new_balance to be 1000, got %v", data["new_balance"])
	}
}

func TestGetBalance(t *testing.T) {
	router, db := setupTestAPI(t)

	// Create test data
	family := &db.Family{Name: "Test Family"}
	if err := db.Create(family).Error; err != nil {
		t.Fatalf("Failed to create test family: %v", err)
	}

	child := &db.User{
		FamilyID:    family.ID,
		Role:        "child",
		DisplayName: "Test Child",
	}
	if err := db.Create(child).Error; err != nil {
		t.Fatalf("Failed to create test child: %v", err)
	}

	rewardType := &db.RewardType{
		FamilyID: family.ID,
		Name:     "Test Reward",
		UnitKind: "money",
	}
	if err := db.Create(rewardType).Error; err != nil {
		t.Fatalf("Failed to create test reward type: %v", err)
	}

	// Create account with balance
	account := &db.Account{
		FamilyID:     family.ID,
		ChildID:      child.ID,
		RewardTypeID: rewardType.ID,
		Balance:      2500,
	}
	if err := db.Create(account).Error; err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}

	// Create request
	req, err := http.NewRequest("GET", "/api/v1/balances?family_id=1&child_id=1&reward_type_id=1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["code"] != float64(0) {
		t.Errorf("Expected code 0, got %v", response["code"])
	}

	data := response["data"].(map[string]interface{})
	if data["balance"] != float64(2500) {
		t.Errorf("Expected balance to be 2500, got %v", data["balance"])
	}
}

func TestAuthMiddleware(t *testing.T) {
	router, _ := setupTestAPI(t)

	// Test without auth header
	req, _ := http.NewRequest("GET", "/api/v1/balances?family_id=1&child_id=1&reward_type_id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	// Test with invalid auth header
	req, _ = http.NewRequest("GET", "/api/v1/balances?family_id=1&child_id=1&reward_type_id=1", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}