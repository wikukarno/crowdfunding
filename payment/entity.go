package payment

type Transaction struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}
