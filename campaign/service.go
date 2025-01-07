package campaign

import (
	"fmt"
	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(UserID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(UserID int) ([]Campaign, error) {
	if UserID != 0 {
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

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	newCampaign := Campaign{}
	newCampaign.Name = input.Name
	newCampaign.ShortDescription = input.ShortDescription
	newCampaign.Description = input.Description
	newCampaign.GoalAmount = input.GoalAmount
	newCampaign.Perks = input.Perks
	newCampaign.UserID = input.User.ID

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	newCampaign.Slug = slug.Make(slugCandidate)

	campaign, err := s.repository.Save(newCampaign)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
