package transaction

import (
	"backend-crowdfunding/campaign"

	"gorm.io/gorm"
)

// TxRunner executes a unit of work where the transaction and campaign
// repositories share a single database transaction. The production
// implementation wraps a GORM transaction; tests can plug in a runner that
// simply calls the function against mocks.
type TxRunner interface {
	Run(fn func(txRepo Repository, campaignRepo campaign.Repository) error) error
}

type gormTxRunner struct {
	db                 *gorm.DB
	repository         Repository
	campaignRepository campaign.Repository
}

func NewTxRunner(db *gorm.DB, repository Repository, campaignRepository campaign.Repository) *gormTxRunner {
	return &gormTxRunner{db: db, repository: repository, campaignRepository: campaignRepository}
}

func (r *gormTxRunner) Run(fn func(Repository, campaign.Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(r.repository.WithTx(tx), r.campaignRepository.WithTx(tx))
	})
}
