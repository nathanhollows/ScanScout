package services

import (
	"context"
	"errors"

	"github.com/nathanhollows/Rapua/internal/models"
)

type GameManagerService struct{}

func NewGameManagerService() *GameManagerService {
	return &GameManagerService{}
}

func (s *GameManagerService) CreateInstance(ctx context.Context, name string, user *models.User) (*models.Instance, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	instance := &models.Instance{
		Name:   name,
		UserID: user.ID,
	}

	if err := instance.Save(ctx); err != nil {
		return nil, err
	}

	user.CurrentInstanceID = instance.ID
	if err := user.Update(ctx); err != nil {
		return nil, err
	}

	return instance, nil
}

func (s *GameManagerService) SwitchInstance(ctx context.Context, user *models.User, instanceID string) (*models.Instance, error) {
	instance, err := models.FindInstanceByID(ctx, instanceID)
	if err != nil {
		return nil, errors.New("instance not found")
	}

	user.CurrentInstanceID = instance.ID
	if err := user.Update(ctx); err != nil {
		return nil, err
	}

	return instance, nil
}

func (s *GameManagerService) DeleteInstance(ctx context.Context, user *models.User, instanceID, confirmName string) error {
	instance, err := models.FindInstanceByID(ctx, instanceID)
	if err != nil {
		return errors.New("instance not found")
	}

	if user.ID != instance.UserID {
		return errors.New("you do not have permission to delete this instance")
	}

	if confirmName != instance.Name {
		return errors.New("instance name does not match confirmation")
	}

	if user.CurrentInstanceID == instance.ID {
		user.CurrentInstanceID = ""
		if err := user.Update(ctx); err != nil {
			return err
		}
	}

	if err := instance.Delete(ctx); err != nil {
		return err
	}

	return nil
}
