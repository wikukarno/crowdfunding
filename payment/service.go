package payment

import (
	"errors"

	"backend-crowdfunding/user"

	"github.com/veritrans/go-midtrans"
)

// ErrPaymentsDisabled is returned when a payment is requested while the Midtrans
// integration is turned off (demo mode). Callers should treat it as a
// "feature unavailable" signal rather than a hard failure.
var ErrPaymentsDisabled = errors.New("payments are disabled")

type service struct {
	serverKey string
	clientKey string
	enabled   bool
}

type Service interface {
	// Enabled reports whether the live payment flow is switched on.
	Enabled() bool
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService(serverKey, clientKey string, enabled bool) *service {
	return &service{serverKey: serverKey, clientKey: clientKey, enabled: enabled}
}

func (s *service) Enabled() bool { return s.enabled }

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	if !s.enabled {
		return "", ErrPaymentsDisabled
	}

	midclient := midtrans.NewClient()
	midclient.ServerKey = s.serverKey
	midclient.ClientKey = s.clientKey
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{Client: midclient}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.ID,
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}
