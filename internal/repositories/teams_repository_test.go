package repositories_test

import (
	"context"
	"testing"

	"github.com/nathanhollows/Rapua/internal/models"
	db "github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestTeamRepository_Update(t *testing.T) {
	cleanup := db.SetupTestDB(t)
	defer cleanup()

	repo := repositories.NewTeamRepository()
	ctx := context.Background()

	sampleTeam := &models.Team{Code: "team1", InstanceID: "instance-1"}
	err := repo.Update(ctx, sampleTeam)
	assert.NoError(t, err, "expected no error when updating team")
}

func TestTeamRepository_Delete(t *testing.T) {
	cleanup := db.SetupTestDB(t)
	defer cleanup()

	repo := repositories.NewTeamRepository()
	ctx := context.Background()

	sampleTeam := &models.Team{Code: "team1", InstanceID: "instance-1"}

	// Insert team first
	err := repo.Update(ctx, sampleTeam)
	assert.NoError(t, err, "expected no error when saving team")

	// Now delete it
	err = repo.Delete(ctx, sampleTeam.Code)
	assert.NoError(t, err, "expected no error when deleting team")
}

func TestTeamRepository_FindAll(t *testing.T) {
	cleanup := db.SetupTestDB(t)
	defer cleanup()

	repo := repositories.NewTeamRepository()
	ctx := context.Background()

	instanceID := "instance-1"
	sampleTeams := []models.Team{
		{Code: "team1", InstanceID: instanceID},
		{Code: "team2", InstanceID: instanceID},
	}

	// Insert teams first
	err := repo.InsertBatch(ctx, sampleTeams)
	assert.NoError(t, err, "expected no error when saving team")

	teams, err := repo.FindAll(ctx, instanceID)
	assert.NoError(t, err, "expected no error when finding all teams")
	assert.Len(t, teams, len(sampleTeams), "expected correct number of teams to be found")
}

func TestTeamRepository_FindAllWithScans(t *testing.T) {
	cleanup := db.SetupTestDB(t)
	defer cleanup()

	repo := repositories.NewTeamRepository()
	ctx := context.Background()

	instanceID := "instance-1"
	sampleTeams := []models.Team{
		{Code: "team1", InstanceID: instanceID},
		{Code: "team2", InstanceID: instanceID},
	}

	// Insert teams first

	err := repo.InsertBatch(ctx, sampleTeams)
	assert.NoError(t, err, "expected no error when saving team")

	teams, err := repo.FindAllWithScans(ctx, instanceID)
	assert.NoError(t, err, "expected no error when finding all teams with scans")
	assert.Len(t, teams, len(sampleTeams), "expected correct number of teams to be found")
}

func TestTeamRepository_InsertBatch(t *testing.T) {
	cleanup := db.SetupTestDB(t)
	defer cleanup()

	repo := repositories.NewTeamRepository()
	ctx := context.Background()

	sampleTeams := []models.Team{{Code: "team1"}, {Code: "team2"}}
	err := repo.InsertBatch(ctx, sampleTeams)
	assert.NoError(t, err, "expected no error when inserting batch of teams")
}

func TestTeamRepository_InsertBatch_UniqueConstraintError(t *testing.T) {
	cleanup := db.SetupTestDB(t)
	defer cleanup()

	repo := repositories.NewTeamRepository()
	ctx := context.Background()

	sampleTeams := []models.Team{{Code: "team1"}, {Code: "team2"}}
	err := repo.InsertBatch(ctx, sampleTeams)
	assert.NoError(t, err, "expected no error when inserting batch of teams")

	// Insert the same teams again to trigger unique constraint error
	err = repo.InsertBatch(ctx, sampleTeams)
	assert.Error(t, err, "expected unique constraint error when inserting duplicate batch of teams")
	assert.Contains(t, err.Error(), "UNIQUE constraint", "expected error message to indicate unique constraint violation")
}
