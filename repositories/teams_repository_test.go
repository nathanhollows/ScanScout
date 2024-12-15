package repositories_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/stretchr/testify/assert"
)

func setupTeamRepo(t *testing.T) (repositories.TeamRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	teamRepository := repositories.NewTeamRepository(db)
	return teamRepository, cleanup
}

func TestTeamRepository_InsertTeam(t *testing.T) {
	repo, cleanup := setupTeamRepo(t)
	defer cleanup()
	ctx := context.Background()

	// Check that teams without an ID are assigned a UUID
	sampleTeam := &models.Team{
		Code:       gofakeit.Password(false, true, false, false, false, 5),
		InstanceID: "instance-1",
	}

	err := repo.InsertBatch(ctx, []models.Team{*sampleTeam})
	assert.NoError(t, err, "expected no error when saving team")

	team, err := repo.FindTeamByCode(ctx, sampleTeam.Code)
	assert.NoError(t, err, "expected no error when finding team")
	assert.NotEmpty(t, team.ID, "expected team to have an ID")

	// Check that teams with duplicate codes are not allowed
	sampleTeam = &models.Team{
		ID:         gofakeit.UUID(),
		InstanceID: "instance-1",
	}

	err = repo.InsertBatch(ctx, []models.Team{*sampleTeam, *sampleTeam})
	assert.Error(t, err, "expected error when saving teams with duplicate codes")

	// Cleanup
	err = repo.DeleteByInstanceID(ctx, sampleTeam.InstanceID)
	assert.NoError(t, err, "expected no error when deleting team")
}

func TestTeamRepository_InsertAndUpdate(t *testing.T) {
	repo, cleanup := setupTeamRepo(t)
	defer cleanup()
	ctx := context.Background()

	sampleTeam := &models.Team{
		ID:         uuid.New().String(),
		Code:       gofakeit.Password(false, true, false, false, false, 5),
		InstanceID: "instance-1",
	}

	// Insert team first
	err := repo.InsertBatch(ctx, []models.Team{*sampleTeam})
	assert.NoError(t, err, "expected no error when saving team")

	// Check that the team was saved
	team, err := repo.FindTeamByCode(ctx, sampleTeam.Code)
	assert.NoError(t, err, "expected no error when finding team")

	// Update the team
	err = repo.Update(ctx, team)
	assert.NoError(t, err, "expected no error when updating team")

	// Cleanup
	err = repo.Delete(ctx, sampleTeam.InstanceID, sampleTeam.Code)
	assert.NoError(t, err, "expected no error when deleting team")
}

func TestTeamRepository_Delete(t *testing.T) {
	repo, cleanup := setupTeamRepo(t)
	defer cleanup()
	ctx := context.Background()

	sampleTeam := []models.Team{{
		ID:         uuid.New().String(),
		Code:       gofakeit.Password(false, true, false, false, false, 5),
		InstanceID: "instance-1",
	}}

	// Insert team first
	err := repo.Update(ctx, &sampleTeam[0])
	assert.NoError(t, err, "expected no error when saving team")

	// Now delete it
	err = repo.Delete(ctx, "instance-1", sampleTeam[0].Code)
	assert.NoError(t, err, "expected no error when deleting team")
}

func TestTeamRepository_FindAll(t *testing.T) {
	repo, cleanup := setupTeamRepo(t)
	defer cleanup()
	ctx := context.Background()

	instanceID := "instance-1"
	sampleTeams := []models.Team{
		{
			ID:         uuid.New().String(),
			Code:       gofakeit.Password(false, true, false, false, false, 5),
			InstanceID: instanceID,
		},
		{
			ID:         uuid.New().String(),
			Code:       gofakeit.Password(false, true, false, false, false, 5),
			InstanceID: instanceID,
		},
	}

	// Insert teams first
	err := repo.InsertBatch(ctx, sampleTeams)
	assert.NoError(t, err, "expected no error when saving team")

	teams, err := repo.FindAll(ctx, instanceID)
	assert.NoError(t, err, "expected no error when finding all teams")
	assert.Len(t, teams, len(sampleTeams), "expected correct number of teams to be found")

	// Cleanup
	for _, team := range teams {
		err = repo.Delete(ctx, instanceID, team.Code)
		assert.NoError(t, err, "expected no error when deleting team")
	}
}

func TestTeamRepository_FindAllWithScans(t *testing.T) {
	repo, cleanup := setupTeamRepo(t)
	defer cleanup()
	ctx := context.Background()

	instanceID := "instance-1"
	sampleTeams := []models.Team{
		{
			ID:         uuid.New().String(),
			Code:       gofakeit.Password(false, true, false, false, false, 5),
			InstanceID: instanceID,
		},
		{
			ID:         uuid.New().String(),
			Code:       gofakeit.Password(false, true, false, false, false, 5),
			InstanceID: instanceID,
		},
	}

	// Insert teams first

	err := repo.InsertBatch(ctx, sampleTeams)
	assert.NoError(t, err, "expected no error when saving team")

	teams, err := repo.FindAllWithScans(ctx, instanceID)
	assert.NoError(t, err, "expected no error when finding all teams with scans")
	assert.Len(t, teams, len(sampleTeams), "expected correct number of teams to be found")

	// Cleanup
	for _, team := range teams {
		err = repo.Delete(ctx, instanceID, team.Code)
		assert.NoError(t, err, "expected no error when deleting team")
	}
}

func TestTeamRepository_InsertBatch(t *testing.T) {
	repo, cleanup := setupTeamRepo(t)
	defer cleanup()
	ctx := context.Background()

	sampleTeams := []models.Team{
		{
			ID:         uuid.New().String(),
			Code:       gofakeit.Password(false, true, false, false, false, 5),
			InstanceID: "instance-1",
		},
		{
			ID:         uuid.New().String(),
			Code:       gofakeit.Password(false, true, false, false, false, 5),
			InstanceID: "instance-1",
		},
	}
	err := repo.InsertBatch(ctx, sampleTeams)
	assert.NoError(t, err, "expected no error when inserting batch of teams")

	// Check that the teams were saved
	for _, team := range sampleTeams {
		_, err = repo.FindTeamByCode(ctx, team.Code)
		assert.NoError(t, err, "expected no error when finding team")
	}

	// Cleanup
	for _, team := range sampleTeams {
		err = repo.Delete(ctx, team.InstanceID, team.Code)
		assert.NoError(t, err, "expected no error when deleting team")
	}
}

func TestTeamRepository_InsertBatch_UniqueConstraintError(t *testing.T) {
	repo, cleanup := setupTeamRepo(t)
	defer cleanup()
	ctx := context.Background()

	sampleTeams := []models.Team{{Code: "team1"}, {Code: "team2"}}
	err := repo.InsertBatch(ctx, sampleTeams)
	assert.NoError(t, err, "expected no error when inserting batch of teams")

	// Insert the same teams again to trigger unique constraint error
	err = repo.InsertBatch(ctx, sampleTeams)
	assert.Error(t, err, "expected unique constraint error when inserting duplicate batch of teams")
	assert.Contains(t, err.Error(), "UNIQUE constraint", "expected error message to indicate unique constraint violation")
}
