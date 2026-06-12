package campaign

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(UserID string) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput, userID string) (Campaign, error)
	UpdateCampaign(InputID GetCampaignDetailInput, inputData CreateCampaignInput, userID string) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string, userID string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(UserID string) ([]Campaign, error) {
	if UserID != "" {
		campaigns, err := s.repository.FindByUserID(UserID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput, userID string) (Campaign, error) {
	newCampaign := Campaign{}
	newCampaign.ID = uuid.NewString()
	newCampaign.Name = input.Name
	newCampaign.ShortDescription = input.ShortDescription
	newCampaign.Description = input.Description
	newCampaign.GoalAmount = input.GoalAmount
	newCampaign.Perks = input.Perks
	newCampaign.UserID = userID

	slugCandidate := fmt.Sprintf("%s %s", input.Name, newCampaign.ID)
	newCampaign.Slug = slug.Make(slugCandidate)

	campaign, err := s.repository.Save(newCampaign)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) UpdateCampaign(InputID GetCampaignDetailInput, inputData CreateCampaignInput, userID string) (Campaign, error) {
	campaign, err := s.repository.FindByID(InputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != userID {
		return campaign, errors.New("not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.GoalAmount = inputData.GoalAmount
	campaign.Perks = inputData.Perks

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string, userID string) (CampaignImage, error) {

	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	if campaign.UserID != userID {
		return CampaignImage{}, errors.New("not an owner of the campaign")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)

		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{}
	campaignImage.ID = uuid.NewString()
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
