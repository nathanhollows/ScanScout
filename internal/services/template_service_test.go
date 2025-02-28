package services_test

import (
	"context"
	"testing"

	"github.com/nathanhollows/Rapua/v3/db"
	"github.com/nathanhollows/Rapua/v3/internal/services"
	"github.com/nathanhollows/Rapua/v3/models"
	"github.com/nathanhollows/Rapua/v3/repositories"
	"github.com/stretchr/testify/assert"
)

func setupTemplateService(t *testing.T) (services.TemplateService, services.InstanceService, func()) {
	t.Helper()
	dbc, cleanup := setupDB(t)

	transactor := db.NewTransactor(dbc)

	// Initialize repositories
	clueRepo := repositories.NewClueRepository(dbc)
	locationRepo := repositories.NewLocationRepository(dbc)
	markerRepo := repositories.NewMarkerRepository(dbc)
	blockStateRepo := repositories.NewBlockStateRepository(dbc)
	blockRepo := repositories.NewBlockRepository(dbc, blockStateRepo)
	instanceRepo := repositories.NewInstanceRepository(dbc)
	instanceSettingsRepo := repositories.NewInstanceSettingsRepository(dbc)
	shareLinkRepo := repositories.NewShareLinkRepository(dbc)
	checkInRepo := repositories.NewCheckInRepository(dbc)
	teamRepo := repositories.NewTeamRepository(dbc)

	// Initialize services
	locationService := services.NewLocationService(transactor, clueRepo, locationRepo, markerRepo, blockRepo)
	teamService := services.NewTeamService(transactor, teamRepo, checkInRepo, blockStateRepo, locationRepo)
	instanceService := services.NewInstanceService(
		transactor,
		locationService, teamService, instanceRepo, instanceSettingsRepo,
	)

	templateService := services.NewTemplateService(
		transactor,
		locationService,
		instanceRepo,
		instanceSettingsRepo,
		shareLinkRepo,
	)
	return templateService, instanceService, cleanup
}

func TestTemplateService_CreateTemplate(t *testing.T) {
	svc, instanceService, cleanup := setupTemplateService(t)
	defer cleanup()

	user := &models.User{ID: "user123", Password: "password", CurrentInstanceID: "instance123"}
	instance, err := instanceService.CreateInstance(context.Background(), "Game1", user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	t.Run("CreateTemplate", func(t *testing.T) {
		tests := []struct {
			name         string
			templateName string
			instanceID   string
			userID       string
			wantErr      bool
		}{
			{"Valid Template", "Template1", instance.ID, user.ID, false},
			{"Empty Template Name", "", instance.ID, user.ID, true},
			{"Invalid Instance ID", "Template1", "invalid", user.ID, true},
			{"Invalid User ID", "Template1", instance.ID, "invalid", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := svc.CreateFromInstance(context.Background(), tt.userID, tt.instanceID, tt.templateName)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestTemplateService_LaunchInstance(t *testing.T) {
	svc, instanceService, cleanup := setupTemplateService(t)
	defer cleanup()

	user := &models.User{ID: "user123", Password: "password", CurrentInstanceID: "instance123"}
	instance, err := instanceService.CreateInstance(context.Background(), "Game1", user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	template, err := svc.CreateFromInstance(context.Background(), user.ID, instance.ID, "Template1")
	assert.NoError(t, err)

	t.Run("LaunchInstance", func(t *testing.T) {
		tests := []struct {
			name         string
			templateID   string
			instanceName string
			userID       string
			wantErr      bool
		}{
			{"Valid Template", template.ID, "Game2", user.ID, false},
			{"Empty Template ID", "", "Game2", user.ID, true},
			{"Empty Instance Name", template.ID, "", user.ID, true},
			{"Invalid Template ID", "invalid", "Game2", user.ID, true},
			{"Invalid User ID", template.ID, "Game2", "invalid", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := svc.LaunchInstance(context.Background(), tt.userID, tt.templateID, tt.instanceName, false)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
