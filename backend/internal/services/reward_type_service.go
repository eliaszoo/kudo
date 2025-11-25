package services

import (
    "reward-system/internal/db"
)

func (s *RewardService) CreateRewardType(rewardType *db.RewardType) error {
	return s.db.Create(rewardType).Error
}