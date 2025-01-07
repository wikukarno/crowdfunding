package handler

import (
	"backend-crowdfunding/campaign"
	"backend-crowdfunding/helper"
	"backend-crowdfunding/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	// Log awal menerima request
	log.Println("Request received to create a campaign")

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// Log error validasi input
		log.Println("Validation Error:", err)
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Log input setelah validasi berhasil
	log.Println("Input after validation:", input)

	// Mendapatkan user dari context
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	// Log user yang sedang membuat campaign
	log.Println("Current user:", currentUser)

	// Membuat campaign
	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		// Log error saat membuat campaign
		log.Println("Service CreateCampaign Error:", err)
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Log campaign baru yang berhasil dibuat
	log.Println("New campaign created:", newCampaign)

	// Mengirimkan response
	response := helper.APIResponse("Campaign has been created", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	log.Println("Response sent:", response)
	c.JSON(http.StatusOK, response)
}
