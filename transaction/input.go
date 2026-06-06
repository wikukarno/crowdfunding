package transaction

type GetCampaignTransactionsInput struct {
	ID string `uri:"id" binding:"required"`
}

type CreateTransactionInput struct {
	Amount     int    `json:"amount" binding:"required"`
	CampaignID string `json:"campaign_id" binding:"required"`
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
