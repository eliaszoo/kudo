package services

import (
	"testing"
	"reward-system/internal/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
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

	return database
}

func TestRewardService_GrantReward(t *testing.T) {
	db := setupTestDB(t)
	service := NewRewardService(db)

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

	// Test grant reward
	result, err := service.GrantReward(family.ID, child.ID, rewardType.ID, 1000, "Test grant", "test-key-1")
	if err != nil {
		t.Fatalf("Failed to grant reward: %v", err)
	}

	if result["transaction_id"] == nil {
		t.Error("Expected transaction_id to be set")
	}

	if result["new_balance"] != int64(1000) {
		t.Errorf("Expected balance to be 1000, got %v", result["new_balance"])
	}

	// Test idempotency
	result2, err := service.GrantReward(family.ID, child.ID, rewardType.ID, 1000, "Test grant", "test-key-1")
	if err != nil {
		t.Fatalf("Failed to grant reward with same idempotency key: %v", err)
	}

	if result["transaction_id"] != result2["transaction_id"] {
		t.Error("Expected same transaction_id for idempotent request")
	}
}

func TestRewardService_SpendReward(t *testing.T) {
	db := setupTestDB(t)
	service := NewRewardService(db)

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

	// First grant some reward
	_, err := service.GrantReward(family.ID, child.ID, rewardType.ID, 1000, "Initial grant", "test-grant")
	if err != nil {
		t.Fatalf("Failed to grant initial reward: %v", err)
	}

	// Test spend reward
	result, err := service.SpendReward(family.ID, child.ID, rewardType.ID, 500, "Test spend", "test-spend-1")
	if err != nil {
		t.Fatalf("Failed to spend reward: %v", err)
	}

	if result["new_balance"] != int64(500) {
		t.Errorf("Expected balance to be 500, got %v", result["new_balance"])
	}

	// Test insufficient balance
	_, err = service.SpendReward(family.ID, child.ID, rewardType.ID, 1000, "Test overspend", "test-overspend")
	if err == nil {
		t.Error("Expected error for insufficient balance")
	}

	if err.Error() != "insufficient balance" {
		t.Errorf("Expected 'insufficient balance' error, got: %v", err)
	}
}

func TestRewardService_GetBalance(t *testing.T) {
	db := setupTestDB(t)
	service := NewRewardService(db)

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

	// Test balance for non-existent account
	balance, err := service.GetBalance(family.ID, child.ID, rewardType.ID)
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}

	if balance != int64(0) {
		t.Errorf("Expected balance to be 0 for non-existent account, got %v", balance)
	}

	// Grant some reward and check balance
	_, err = service.GrantReward(family.ID, child.ID, rewardType.ID, 1500, "Test grant", "test-balance")
	if err != nil {
		t.Fatalf("Failed to grant reward: %v", err)
	}

	balance, err = service.GetBalance(family.ID, child.ID, rewardType.ID)
	if err != nil {
		t.Fatalf("Failed to get balance after grant: %v", err)
	}

	if balance != int64(1500) {
		t.Errorf("Expected balance to be 1500, got %v", balance)
	}
}

func TestRewardService_ListTransactions(t *testing.T) {
	db := setupTestDB(t)
	service := NewRewardService(db)

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

	// Create multiple transactions
	transactions := []struct {
		value int64
		type  string
		note  string
		key   string
	}{
		{1000, "credit", "First grant", "grant-1"},
		{500, "debit", "First spend", "spend-1"},
		{2000, "credit", "Second grant", "grant-2"},
		{300, "debit", "Second spend", "spend-2"},
	}

	for _, tx := range transactions {
		if tx.type == "credit" {
			_, err := service.GrantReward(family.ID, child.ID, rewardType.ID, tx.value, tx.note, tx.key)
			if err != nil {
				t.Fatalf("Failed to create grant transaction: %v", err)
			}
		} else {
			_, err := service.SpendReward(family.ID, child.ID, rewardType.ID, tx.value, tx.note, tx.key)
			if err != nil {
				t.Fatalf("Failed to create spend transaction: %v", err)
			}
		}
	}

	// Test listing transactions
	result, err := service.ListTransactions(family.ID, child.ID, rewardType.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list transactions: %v", err)
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 transactions, got %d", len(result))
	}

	// Verify transaction order (should be newest first)
	if result[0].Type != "debit" || result[0].Value != 300 {
		t.Error("Expected newest transaction to be debit of 300")
	}
}

func TestRewardService_CreateRewardType(t *testing.T) {
	db := setupTestDB(t)
	service := NewRewardService(db)

	// Create test family
	family := &db.Family{Name: "Test Family"}
	if err := db.Create(family).Error; err != nil {
		t.Fatalf("Failed to create test family: %v", err)
	}

	// Test create reward type
	rewardType := &db.RewardType{
		FamilyID:  family.ID,
		Name:      "Test Reward Type",
		UnitKind:  "money",
		UnitLabel: "元",
	}

	err := service.CreateRewardType(rewardType)
	if err != nil {
		t.Fatalf("Failed to create reward type: %v", err)
	}

	if rewardType.ID == 0 {
		t.Error("Expected reward type ID to be set after creation")
	}

	// Verify the reward type was created
	var found db.RewardType
	if err := db.First(&found, rewardType.ID).Error; err != nil {
		t.Fatalf("Failed to find created reward type: %v", err)
	}

	if found.Name != "Test Reward Type" {
		t.Errorf("Expected reward type name to be 'Test Reward Type', got %s", found.Name)
	}
}