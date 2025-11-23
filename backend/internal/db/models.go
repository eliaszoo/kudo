package db

import (
	"time"
)

type Family struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FamilyID     uint64    `gorm:"not null;index" json:"family_id"`
	Role         string    `gorm:"type:enum('guardian','child');not null" json:"role"`
	DisplayName  string    `gorm:"size:64;not null" json:"display_name"`
	WechatOpenID string    `gorm:"size:128;uniqueIndex" json:"wechat_openid,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	Family Family `gorm:"foreignKey:FamilyID" json:"family,omitempty"`
}

type RewardType struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FamilyID  uint64    `gorm:"not null;index" json:"family_id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	UnitKind  string    `gorm:"type:enum('money','time','points','custom');not null" json:"unit_kind"`
	UnitLabel string    `gorm:"size:32" json:"unit_label,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	Family Family `gorm:"foreignKey:FamilyID" json:"family,omitempty"`
	
	gorm:"uniqueIndex:uniq_family_name,composite:(family_id,name)"
}

type Account struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FamilyID     uint64    `gorm:"not null;index" json:"family_id"`
	ChildID      uint64    `gorm:"not null;index" json:"child_id"`
	RewardTypeID uint64    `gorm:"not null;index" json:"reward_type_id"`
	Balance      int64     `gorm:"default:0;not null" json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	Child      User       `gorm:"foreignKey:ChildID" json:"child,omitempty"`
	RewardType RewardType `gorm:"foreignKey:RewardTypeID" json:"reward_type,omitempty"`
	
	gorm:"uniqueIndex:uniq_acc,composite:(child_id,reward_type_id)"
}

type Transaction struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID       uint64    `gorm:"not null;index" json:"account_id"`
	Type            string    `gorm:"type:enum('credit','debit');not null" json:"type"`
	Value           int64     `gorm:"not null" json:"value"`
	Note            string    `gorm:"size:255" json:"note,omitempty"`
	CreatedBy       uint64    `gorm:"not null;index" json:"created_by"`
	IdempotencyKey  string    `gorm:"size:64;index" json:"idempotency_key,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	
	Account   Account `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Creator   User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

type AuditLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FamilyID  uint64    `gorm:"index" json:"family_id"`
	UserID    uint64    `gorm:"index" json:"user_id"`
	Action    string    `gorm:"size:32;not null" json:"action"`
	Payload   string    `gorm:"type:json" json:"payload,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}