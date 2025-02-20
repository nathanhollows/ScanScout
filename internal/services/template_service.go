package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/nathanhollows/Rapua/v3/db"
	"github.com/nathanhollows/Rapua/v3/models"
	"github.com/nathanhollows/Rapua/v3/repositories"
)

type TemplateService struct {
	transactor           db.Transactor
	locationService      LocationService
	userService          UserService
	teamService          TeamService
	instanceRepo         repositories.InstanceRepository
	instanceSettingsRepo repositories.InstanceSettingsRepository
}

func NewTemplateService(
	transactor db.Transactor,
	locationService LocationService,
	userService UserService,
	teamService TeamService,
	instanceRepo repositories.InstanceRepository,
	instanceSettingsRepo repositories.InstanceSettingsRepository,
) TemplateService {
	return TemplateService{
		transactor:           transactor,
		locationService:      locationService,
		userService:          userService,
		teamService:          teamService,
		instanceRepo:         instanceRepo,
		instanceSettingsRepo: instanceSettingsRepo,
	}
}

// CreateFromInstance creates a new template from an existing instance.
func (s *TemplateService) CreateFromInstance(ctx context.Context, userID, instanceID, name string) (*models.Instance, error) {
	if userID == "" {
		return nil, NewValidationError("userID")
	}

	oldInstance, err := s.instanceRepo.GetByID(ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("finding instance: %w", err)
	}

	if oldInstance.UserID != userID {
		return nil, ErrPermissionDenied
	}

	if oldInstance.IsTemplate {
		return nil, errors.New("cannot create a template from a template")
	}

	locations, err := s.locationService.FindByInstance(ctx, oldInstance.ID)
	if err != nil {
		return nil, fmt.Errorf("finding locations: %w", err)
	}

	if name == "" {
		return nil, NewValidationError("name")
	}
	if instanceID == "" {
		return nil, NewValidationError("id")
	}

	newInstance := &models.Instance{
		Name:       name,
		UserID:     userID,
		IsTemplate: true,
	}

	if err := s.instanceRepo.Create(ctx, newInstance); err != nil {
		return nil, fmt.Errorf("creating instance: %w", err)
	}

	// Copy locations
	for _, location := range locations {
		_, err := s.locationService.DuplicateLocation(ctx, location, newInstance.ID)
		if err != nil {
			return nil, fmt.Errorf("duplicating location: %w", err)
		}
	}

	// Copy settings
	settings := oldInstance.Settings
	settings.InstanceID = newInstance.ID
	if err := s.instanceSettingsRepo.Create(ctx, &settings); err != nil {
		return nil, fmt.Errorf("creating settings: %w", err)
	}

	return newInstance, nil
}

// Find retrieves all templates.
func (s *TemplateService) Find(ctx context.Context, userID string) ([]models.Instance, error) {
	instances, err := s.instanceRepo.FindTemplates(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("finding templates: %w", err)
	}
	return instances, nil
}

//
// // FindInstanceIDsForUser implements InstanceService.
// func (s *instanceService) FindInstanceIDsForUser(ctx context.Context, userID string) ([]string, error) {
// 	instances, err := s.instanceRepo.FindByUserID(ctx, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("finding instances for user: %w", err)
// 	}
//
// 	ids := make([]string, len(instances))
// 	for i, instance := range instances {
// 		ids[i] = instance.ID
// 	}
// 	return ids, nil
// }
//
// // DeleteInstance implements InstanceService.
// func (s *instanceService) DeleteInstance(ctx context.Context, user *models.User, instanceID string, confirmName string) (bool, error) {
// 	if user == nil {
// 		return false, ErrUserNotAuthenticated
// 	}
//
// 	// Check if the user has permission to delete the instance
// 	instance, err := s.instanceRepo.GetByID(ctx, instanceID)
// 	if err != nil {
// 		return false, fmt.Errorf("finding instance: %w", err)
// 	}
//
// 	if user.ID != instance.UserID {
// 		return false, ErrPermissionDenied
// 	}
//
// 	// If the name does not match the confirmation, return an error
// 	if confirmName != instance.Name {
// 		return false, errors.New("instance name does not match confirmation")
// 	}
//
// 	// Check if the user is currently using this instance
// 	if user.CurrentInstanceID == instance.ID {
// 		return false, errors.New("cannot delete an instance that is currently in use")
// 	}
//
// 	// Start transaction
// 	tx, err := s.transactor.BeginTx(ctx, &sql.TxOptions{})
// 	if err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("beginning transaction: %w", err)
// 	}
//
// 	err = s.instanceRepo.Delete(ctx, tx, instanceID)
// 	if err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("deleting instance: %w", err)
// 	}
//
// 	err = s.instanceSettingsRepo.Delete(ctx, tx, instanceID)
// 	if err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("deleting instance settings: %w", err)
// 	}
//
// 	err = s.teamService.DeleteByInstanceID(ctx, tx, instanceID)
// 	if err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("deleting teams: %w", err)
// 	}
//
// 	err = tx.Commit()
// 	if err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("committing transaction: %w", err)
// 	}
//
// 	return true, nil
// }
//
// // SwitchInstance implements InstanceService.
// func (s *instanceService) SwitchInstance(ctx context.Context, user *models.User, instanceID string) (*models.Instance, error) {
// 	if user == nil {
// 		return nil, ErrUserNotAuthenticated
// 	}
//
// 	instance, err := s.instanceRepo.GetByID(ctx, instanceID)
// 	if err != nil {
// 		return nil, errors.New("instance not found")
// 	}
//
// 	if instance.IsTemplate {
// 		return nil, errors.New("cannot switch to a template")
// 	}
//
// 	// Make sure the user has permission to switch to this instance
// 	if instance.UserID != user.ID {
// 		return nil, ErrPermissionDenied
// 	}
//
// 	user.CurrentInstanceID = instance.ID
// 	err = s.userService.UpdateUser(ctx, user)
// 	if err != nil {
// 		return nil, fmt.Errorf("updating user: %w", err)
// 	}
//
// 	return instance, nil
// }
