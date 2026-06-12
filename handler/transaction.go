package handler

import (
	"errors"
	"net/http"

	"backend-crowdfunding/helper"
	"backend-crowdfunding/payment"
	"backend-crowdfunding/transaction"
	"backend-crowdfunding/user"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *TransactionHandler {
	return &TransactionHandler{service}
}

// @Summary  List transactions of a campaign
// @Tags     transactions
// @Param    id path string true "Campaign id"
// @Success  200 {object} helper.Response
// @Security BearerAuth
// @Router   /campaigns/{id}/transactions [get]
func (h *TransactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	if err := c.ShouldBindUri(&input); err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	transactions, err := h.service.GetTransactionByCampaignID(input, currentUser.ID)
	if err != nil {
		c.Error(err)
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaign's transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

// @Summary  List the current user's transactions
// @Tags     transactions
// @Success  200 {object} helper.Response
// @Security BearerAuth
// @Router   /transactions [get]
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	transactions, err := h.service.GetTransactionsByUserID(currentUser.ID)
	if err != nil {
		c.Error(err)
		response := helper.APIResponse("Failed to get user's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of user's transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

// @Summary  Create a donation transaction
// @Tags     transactions
// @Param    payload body transaction.CreateTransactionInput true "Transaction details"
// @Success  200 {object} helper.Response
// @Security BearerAuth
// @Router   /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIResponse("Failed to create transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	newTransaction, err := h.service.CreateTransaction(input, currentUser)
	if err != nil {
		// Demo mode: payments are switched off. Tell the client clearly so it
		// can surface its "contact support" prompt instead of a generic error.
		if errors.Is(err, payment.ErrPaymentsDisabled) {
			response := helper.APIResponse("Online payments are in demo mode. Contact support to enable a live demo.", http.StatusServiceUnavailable, "error", nil)
			c.JSON(http.StatusServiceUnavailable, response)
			return
		}

		c.Error(err)
		response := helper.APIResponse("Failed to create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Transaction has been created", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}

// @Summary Midtrans payment notification webhook
// @Tags    transactions
// @Param   payload body transaction.TransactionNotificationInput true "Notification payload"
// @Success 200 {object} transaction.TransactionNotificationInput
// @Router  /transactions/notification [post]
func (h *TransactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.service.ProcessPayment(input); err != nil {
		c.Error(err)
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input)
}
