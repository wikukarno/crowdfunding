package transaction

import (
	"testing"

	"backend-crowdfunding/campaign"
	"backend-crowdfunding/payment"
	"backend-crowdfunding/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockPaymentService struct {
	enabled bool
	urlFn   func(payment.Transaction, user.User) (string, error)
}

func (m *mockPaymentService) Enabled() bool { return m.enabled }
func (m *mockPaymentService) GetPaymentURL(tx payment.Transaction, u user.User) (string, error) {
	if m.urlFn != nil {
		return m.urlFn(tx, u)
	}
	return "https://pay.example/redirect", nil
}

// directRunner runs the unit of work against the mocks without a real database
// transaction, which is all the service logic needs to be exercised.
type directRunner struct {
	repo     Repository
	campRepo campaign.Repository
}

func (r directRunner) Run(fn func(Repository, campaign.Repository) error) error {
	return fn(r.repo, r.campRepo)
}

type mockRepository struct {
	getByCampaignFn func(string) ([]Transaction, error)
	getByUserFn     func(string) ([]Transaction, error)
	getByIDFn       func(string) (Transaction, error)
	saveFn          func(Transaction) (Transaction, error)
	updateFn        func(Transaction) (Transaction, error)
}

func (m *mockRepository) GetByCampaignID(id string) ([]Transaction, error) {
	return m.getByCampaignFn(id)
}
func (m *mockRepository) GetByUserID(id string) ([]Transaction, error) { return m.getByUserFn(id) }
func (m *mockRepository) GetByID(id string) (Transaction, error)       { return m.getByIDFn(id) }
func (m *mockRepository) Save(tx Transaction) (Transaction, error)     { return m.saveFn(tx) }
func (m *mockRepository) Update(tx Transaction) (Transaction, error)   { return m.updateFn(tx) }
func (m *mockRepository) WithTx(*gorm.DB) Repository                   { return m }

type mockCampaignRepository struct {
	findByIDFn func(string) (campaign.Campaign, error)
	updateFn   func(campaign.Campaign) (campaign.Campaign, error)
}

func (m *mockCampaignRepository) FindAll() ([]campaign.Campaign, error)            { return nil, nil }
func (m *mockCampaignRepository) FindByUserID(string) ([]campaign.Campaign, error) { return nil, nil }
func (m *mockCampaignRepository) FindByID(id string) (campaign.Campaign, error) {
	return m.findByIDFn(id)
}
func (m *mockCampaignRepository) Save(c campaign.Campaign) (campaign.Campaign, error) {
	return c, nil
}
func (m *mockCampaignRepository) Update(c campaign.Campaign) (campaign.Campaign, error) {
	return m.updateFn(c)
}
func (m *mockCampaignRepository) CreateImage(i campaign.CampaignImage) (campaign.CampaignImage, error) {
	return i, nil
}
func (m *mockCampaignRepository) MarkAllImagesAsNonPrimary(string) (bool, error) {
	return true, nil
}
func (m *mockCampaignRepository) WithTx(*gorm.DB) campaign.Repository { return m }

func TestCreateTransactionDisabledDoesNotSave(t *testing.T) {
	saveCalled := false
	txRepo := &mockRepository{
		saveFn: func(tx Transaction) (Transaction, error) {
			saveCalled = true
			return tx, nil
		},
	}
	campaignRepo := &mockCampaignRepository{}
	pay := &mockPaymentService{enabled: false}
	service := NewService(txRepo, campaignRepo, pay, directRunner{txRepo, campaignRepo})

	_, err := service.CreateTransaction(
		CreateTransactionInput{Amount: 10000, CampaignID: "campaign-1"},
		user.User{ID: "user-1"},
	)

	require.ErrorIs(t, err, payment.ErrPaymentsDisabled)
	assert.False(t, saveCalled, "no transaction should be persisted while payments are disabled")
}

func TestCreateTransactionEnabledSavesPaymentURL(t *testing.T) {
	var savedStatus string
	txRepo := &mockRepository{
		saveFn: func(tx Transaction) (Transaction, error) {
			savedStatus = tx.Status
			return tx, nil
		},
		updateFn: func(tx Transaction) (Transaction, error) { return tx, nil },
	}
	campaignRepo := &mockCampaignRepository{}
	pay := &mockPaymentService{
		enabled: true,
		urlFn: func(payment.Transaction, user.User) (string, error) {
			return "https://pay.example/redirect", nil
		},
	}
	service := NewService(txRepo, campaignRepo, pay, directRunner{txRepo, campaignRepo})

	tx, err := service.CreateTransaction(
		CreateTransactionInput{Amount: 25000, CampaignID: "campaign-1"},
		user.User{ID: "user-1"},
	)

	require.NoError(t, err)
	assert.Equal(t, "pending", savedStatus)
	assert.Equal(t, "https://pay.example/redirect", tx.PaymentURL)
}

func TestProcessPaymentMarksSettlementAsPaid(t *testing.T) {
	var updatedCampaign campaign.Campaign
	txRepo := &mockRepository{
		getByIDFn: func(string) (Transaction, error) {
			return Transaction{ID: "tx-1", CampaignID: "campaign-5", Amount: 1000}, nil
		},
		updateFn: func(tx Transaction) (Transaction, error) { return tx, nil },
	}
	campaignRepo := &mockCampaignRepository{
		findByIDFn: func(string) (campaign.Campaign, error) {
			return campaign.Campaign{ID: "campaign-5", BackerCount: 2, CurrentAmount: 500}, nil
		},
		updateFn: func(c campaign.Campaign) (campaign.Campaign, error) {
			updatedCampaign = c
			return c, nil
		},
	}
	service := NewService(txRepo, campaignRepo, nil, directRunner{txRepo, campaignRepo})

	err := service.ProcessPayment(TransactionNotificationInput{
		OrderID:           "1",
		TransactionStatus: "settlement",
	})

	require.NoError(t, err)
	assert.Equal(t, 3, updatedCampaign.BackerCount)
	assert.Equal(t, 1500, updatedCampaign.CurrentAmount)
}

func TestProcessPaymentMarksCancel(t *testing.T) {
	var savedStatus string
	txRepo := &mockRepository{
		getByIDFn: func(string) (Transaction, error) {
			return Transaction{ID: "tx-1", CampaignID: "campaign-5"}, nil
		},
		updateFn: func(tx Transaction) (Transaction, error) {
			savedStatus = tx.Status
			return tx, nil
		},
	}
	campaignRepo := &mockCampaignRepository{
		findByIDFn: func(string) (campaign.Campaign, error) { return campaign.Campaign{ID: "campaign-5"}, nil },
		updateFn:   func(c campaign.Campaign) (campaign.Campaign, error) { return c, nil },
	}
	service := NewService(txRepo, campaignRepo, nil, directRunner{txRepo, campaignRepo})

	err := service.ProcessPayment(TransactionNotificationInput{
		OrderID:           "1",
		TransactionStatus: "cancel",
	})

	require.NoError(t, err)
	assert.Equal(t, "cancelled", savedStatus)
}
