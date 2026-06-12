package transaction

import (
	"errors"

	"backend-crowdfunding/campaign"
	"backend-crowdfunding/payment"
	"backend-crowdfunding/user"

	"github.com/google/uuid"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
	txRunner           TxRunner
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionsInput, userID string) ([]Transaction, error)
	GetTransactionsByUserID(userID string) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput, currentUser user.User) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service, txRunner TxRunner) *service {
	return &service{repository, campaignRepository, paymentService, txRunner}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionsInput, userID string) ([]Transaction, error) {

	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != userID {
		return []Transaction{}, errors.New("not an owner of the campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionsByUserID(userID string) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput, currentUser user.User) (Transaction, error) {
	// Bail out before writing a pending row when payments are in demo mode, so
	// we don't accumulate orphaned transactions that can never be paid.
	if !s.paymentService.Enabled() {
		return Transaction{}, payment.ErrPaymentsDisabled
	}

	transaction := Transaction{}
	transaction.ID = uuid.NewString()
	transaction.Amount = input.Amount
	transaction.CampaignID = input.CampaignID
	transaction.UserID = currentUser.ID
	transaction.Status = "pending"

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, currentUser)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	transactionID := input.OrderID

	// Updating the transaction status and bumping the campaign totals must
	// happen together, otherwise a crash between the two writes leaves a paid
	// donation that never counted towards the campaign.
	return s.txRunner.Run(func(txRepo Repository, campaignRepo campaign.Repository) error {
		trx, err := txRepo.GetByID(transactionID)
		if err != nil {
			return err
		}

		switch {
		case input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept":
			trx.Status = "paid"
		case input.TransactionStatus == "settlement":
			trx.Status = "paid"
		case input.TransactionStatus == "cancel", input.TransactionStatus == "deny", input.TransactionStatus == "expire":
			trx.Status = "cancelled"
		}

		updated, err := txRepo.Update(trx)
		if err != nil {
			return err
		}

		if updated.Status != "paid" {
			return nil
		}

		camp, err := campaignRepo.FindByID(updated.CampaignID)
		if err != nil {
			return err
		}

		camp.BackerCount++
		camp.CurrentAmount += updated.Amount

		_, err = campaignRepo.Update(camp)
		return err
	})
}
