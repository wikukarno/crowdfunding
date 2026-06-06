package transaction

import (
	"backend-crowdfunding/campaign"
	"backend-crowdfunding/user"
	"time"
)

type Transaction struct {
	ID         string `json:"id" gorm:"type:char(36);primaryKey"`
	CampaignID string `json:"campaign_id" gorm:"type:char(36)"`
	UserID     string `json:"user_id" gorm:"type:char(36)"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	User       user.User
	PaymentURL string `json:"payment_url"`
	Campaign   campaign.Campaign
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
