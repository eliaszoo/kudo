package services

import (
	"reward-system/internal/db"
	"gorm.io/gorm"
)

func (s *RewardService) CreateRewardType(rewardType *db.RewardType) error {
	return s.db.Create(rewardType).Error
}