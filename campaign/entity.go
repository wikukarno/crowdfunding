package campaign

import (
	"backend-crowdfunding/user"
	"time"
)

type Campaign struct {
	ID               string `gorm:"type:char(36);primaryKey"`
	UserID           string `gorm:"type:char(36)"`
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User             user.User
}

type CampaignImage struct {
	ID         string `gorm:"type:char(36);primaryKey"`
	CampaignID string `gorm:"type:char(36)"`
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
