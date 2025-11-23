package services

import (
	"fmt"
	"reward-system/internal/db"

	"gorm.io/gorm"
)

type RewardService struct {
	db *gorm.DB
}

func NewRewardService(db *gorm.DB) *RewardService {
	return &RewardService{db: db}
}

func (s *RewardService) GrantReward(familyID, childID, rewardTypeID uint64, value int64, note, idempotencyKey string) (map[string]interface{}, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check idempotency
	if idempotencyKey != "" {
		var existingTx db.Transaction
		if err := tx.Where("idempotency_key = ?", idempotencyKey).First(&existingTx).Error; err == nil {
			tx.Rollback()
			return map[string]interface{}{
				"transaction_id": existingTx.ID,
				"new_balance":      s.getAccountBalance(tx, childID, rewardTypeID),
			}, nil
		}
	}

	// Get or create account
	account, err := s.getOrCreateAccount(tx, familyID, childID, rewardTypeID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Lock account for update
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, account.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create transaction
	transaction := &db.Transaction{
		AccountID:      account.ID,
		Type:           "credit",
		Value:          value,
		Note:           note,
		CreatedBy:      childID, // In real implementation, this should be the guardian ID
		IdempotencyKey: idempotencyKey,
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update account balance
	account.Balance += value
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"transaction_id": transaction.ID,
		"new_balance":      account.Balance,
	}, nil
}

func (s *RewardService) SpendReward(familyID, childID, rewardTypeID uint64, value int64, note, idempotencyKey string) (map[string]interface{}, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check idempotency
	if idempotencyKey != "" {
		var existingTx db.Transaction
		if err := tx.Where("idempotency_key = ?", idempotencyKey).First(&existingTx).Error; err == nil {
			tx.Rollback()
			return map[string]interface{}{
				"transaction_id": existingTx.ID,
				"new_balance":      s.getAccountBalance(tx, childID, rewardTypeID),
			}, nil
		}
	}

	// Get account
	var account db.Account
	if err := tx.Where("child_id = ? AND reward_type_id = ?", childID, rewardTypeID).First(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("account not found")
	}

	// Lock account for update
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, account.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Check sufficient balance
	if account.Balance < value {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create transaction
	transaction := &db.Transaction{
		AccountID:      account.ID,
		Type:           "debit",
		Value:          value,
		Note:           note,
		CreatedBy:      childID, // In real implementation, this should be the guardian ID
		IdempotencyKey: idempotencyKey,
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update account balance
	account.Balance -= value
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"transaction_id": transaction.ID,
		"new_balance":      account.Balance,
	}, nil
}

func (s *RewardService) GetBalance(familyID, childID, rewardTypeID uint64) (int64, error) {
	var account db.Account
	if err := s.db.Where("child_id = ? AND reward_type_id = ?", childID, rewardTypeID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return account.Balance, nil
}

func (s *RewardService) ListTransactions(familyID, childID, rewardTypeID uint64, limit int, beforeID uint64) ([]db.Transaction, error) {
	query := s.db.Where("account_id IN (SELECT id FROM accounts WHERE child_id = ?)", childID)
	
	if rewardTypeID > 0 {
		query = query.Where("reward_type_id = ?", rewardTypeID)
	}
	
	if beforeID > 0 {
		query = query.Where("id < ?", beforeID)
	}
	
	var transactions []db.Transaction
	if err := query.Order("id DESC").Limit(limit).Find(&transactions).Error; err != nil {
		return nil, err
	}
	
	return transactions, nil
}

func (s *RewardService) AdjustTransaction(transactionID uint64, newValue *int64, newNote *string) error {
	return s.db.Model(&db.Transaction{}).Where("id = ?", transactionID).Updates(map[string]interface{}{
		"value": newValue,
		"note":  newNote,
	}).Error
}

func (s *RewardService) getOrCreateAccount(tx *gorm.DB, familyID, childID, rewardTypeID uint64) (*db.Account, error) {
	var account db.Account
	
	if err := tx.Where("child_id = ? AND reward_type_id = ?", childID, rewardTypeID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new account
			account = db.Account{
				FamilyID:     familyID,
				ChildID:      childID,
				RewardTypeID: rewardTypeID,
				Balance:      0,
			}
			if err := tx.Create(&account).Error; err != nil {
				return nil, err
			}
			return &account, nil
		}
		return nil, err
	}
	
	return &account, nil
}

func (s *RewardService) getAccountBalance(tx *gorm.DB, childID, rewardTypeID uint64) int64 {
	var account db.Account
	if err := tx.Where("child_id = ? AND reward_type_id = ?", childID, rewardTypeID).First(&account).Error; err != nil {
		return 0
	}
	return account.Balance
}