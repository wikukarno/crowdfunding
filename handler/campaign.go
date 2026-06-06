package handler

import (
	"fmt"
	"net/http"

	"backend-crowdfunding/campaign"
	"backend-crowdfunding/helper"
	"backend-crowdfunding/storage"
	"backend-crowdfunding/user"

	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	service  campaign.Service
	uploader storage.Uploader
}

func NewCampaignHandler(service campaign.Service, uploader storage.Uploader) *CampaignHandler {
	return &CampaignHandler{service, uploader}
}

// @Summary List campaigns, optionally filtered by user
// @Tags    campaigns
// @Param   user_id query string false "Filter by owner user id"
// @Success 200 {object} helper.Response
// @Router  /campaigns [get]
func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	userID := c.Query("user_id")

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		c.Error(err)
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

// @Summary Get campaign detail by id
// @Tags    campaigns
// @Param   id path string true "Campaign id"
// @Success 200 {object} helper.Response
// @Router  /campaigns/{id} [get]
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	if err := c.ShouldBindUri(&input); err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		c.Error(err)
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

// @Summary  Create a campaign
// @Tags     campaigns
// @Param    payload body campaign.CreateCampaignInput true "Campaign details"
// @Success  200 {object} helper.Response
// @Security BearerAuth
// @Router   /campaigns [post]
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	newCampaign, err := h.service.CreateCampaign(input, currentUser.ID)
	if err != nil {
		c.Error(err)
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign has been created", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

// @Summary  Update a campaign
// @Tags     campaigns
// @Param    id path string true "Campaign id"
// @Param    payload body campaign.CreateCampaignInput true "Updated details"
// @Success  200 {object} helper.Response
// @Security BearerAuth
// @Router   /campaigns/{id} [put]
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	if err := c.ShouldBindUri(&inputID); err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput
	if err := c.ShouldBindJSON(&inputData); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData, currentUser.ID)
	if err != nil {
		c.Error(err)
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign has been updated", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

// @Summary  Upload a campaign image
// @Tags     campaigns
// @Param    campaign_id formData string true "Campaign id"
// @Param    is_primary formData bool false "Mark as primary image"
// @Param    file formData file true "Image file"
// @Success  200 {object} helper.Response
// @Security BearerAuth
// @Router   /campaign-images [post]
func (h *CampaignHandler) UploadCampaignImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	if err := c.ShouldBind(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	opened, err := file.Open()
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	defer opened.Close()

	key := objectKey(fmt.Sprintf("campaigns/%s", userID), file.Filename)
	url, err := h.uploader.Upload(c.Request.Context(), key, opened, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		c.Error(err)
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := h.service.SaveCampaignImage(input, url, currentUser.ID); err != nil {
		c.Error(err)
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image has been uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
