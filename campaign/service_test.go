package campaign

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockRepository struct {
	findAllFn      func() ([]Campaign, error)
	findByUserIDFn func(string) ([]Campaign, error)
	findByIDFn     func(string) (Campaign, error)
	saveFn         func(Campaign) (Campaign, error)
	updateFn       func(Campaign) (Campaign, error)
	createImageFn  func(CampaignImage) (CampaignImage, error)
	markNonPrimary func(string) (bool, error)
}

func (m *mockRepository) FindAll() ([]Campaign, error)               { return m.findAllFn() }
func (m *mockRepository) FindByUserID(id string) ([]Campaign, error) { return m.findByUserIDFn(id) }
func (m *mockRepository) FindByID(id string) (Campaign, error)       { return m.findByIDFn(id) }
func (m *mockRepository) Save(c Campaign) (Campaign, error)          { return m.saveFn(c) }
func (m *mockRepository) Update(c Campaign) (Campaign, error)        { return m.updateFn(c) }
func (m *mockRepository) CreateImage(i CampaignImage) (CampaignImage, error) {
	return m.createImageFn(i)
}
func (m *mockRepository) MarkAllImagesAsNonPrimary(id string) (bool, error) {
	return m.markNonPrimary(id)
}
func (m *mockRepository) WithTx(*gorm.DB) Repository { return m }

func TestGetCampaignsListsAllWhenNoUserFilter(t *testing.T) {
	allCalled := false
	repo := &mockRepository{
		findAllFn: func() ([]Campaign, error) {
			allCalled = true
			return []Campaign{{ID: "campaign-1"}}, nil
		},
	}
	service := NewService(repo)

	campaigns, err := service.GetCampaigns("")

	require.NoError(t, err)
	assert.True(t, allCalled)
	assert.Len(t, campaigns, 1)
}

func TestGetCampaignsFiltersByUser(t *testing.T) {
	var gotUserID string
	repo := &mockRepository{
		findByUserIDFn: func(id string) ([]Campaign, error) {
			gotUserID = id
			return nil, nil
		},
	}
	service := NewService(repo)

	_, err := service.GetCampaigns("user-9")

	require.NoError(t, err)
	assert.Equal(t, "user-9", gotUserID)
}

func TestCreateCampaignBuildsSlug(t *testing.T) {
	var saved Campaign
	repo := &mockRepository{
		saveFn: func(c Campaign) (Campaign, error) {
			saved = c
			return c, nil
		},
	}
	service := NewService(repo)

	_, err := service.CreateCampaign(CreateCampaignInput{
		Name: "Help Build School",
	}, "user-12")

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(saved.Slug, "help-build-school-"))
	assert.NotEmpty(t, saved.ID)
	assert.Equal(t, "user-12", saved.UserID)
}

func TestUpdateCampaignRejectsNonOwner(t *testing.T) {
	repo := &mockRepository{
		findByIDFn: func(string) (Campaign, error) {
			return Campaign{ID: "campaign-1", UserID: "owner"}, nil
		},
	}
	service := NewService(repo)

	_, err := service.UpdateCampaign(
		GetCampaignDetailInput{ID: "campaign-1"},
		CreateCampaignInput{},
		"other",
	)

	assert.EqualError(t, err, "not an owner of the campaign")
}
