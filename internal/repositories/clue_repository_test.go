package repositories_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/migrate"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/models"
	"github.com/stretchr/testify/assert"
)

func setupClueRepo(t *testing.T) (repositories.ClueRepository, func()) {
	t.Helper()
	os.Setenv("DB_CONNECTION", "file::memory:?cache=shared")
	os.Setenv("DB_TYPE", "sqlite3")
	db := db.MustOpen()

	// Create tables
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	migrate.CreateTables(logger, db)

	clueRepo := repositories.NewClueRepository(db)
	return clueRepo, func() {
		db.Close()
	}
}

func TestClueRepository_Save(t *testing.T) {
	repo, cleanup := setupClueRepo(t)
	defer cleanup()
	ctx := context.Background()

	clue := &models.Clue{
		ID:         uuid.New().String(),
		InstanceID: "instance-1",
		LocationID: "location-1",
		Content:    "This is a test clue.",
	}

	err := repo.Save(ctx, clue)
	assert.NoError(t, err, "expected no error when saving clue")
}

func TestClueRepository_Delete(t *testing.T) {
	repo, cleanup := setupClueRepo(t)
	defer cleanup()
	ctx := context.Background()

	clue := &models.Clue{
		ID:         uuid.New().String(),
		InstanceID: "instance-1",
		LocationID: "location-1",
		Content:    "This is a test clue.",
	}

	// Save clue first
	err := repo.Save(ctx, clue)
	assert.NoError(t, err, "expected no error when saving clue")

	// Now delete it
	err = repo.Delete(ctx, clue.ID)
	assert.NoError(t, err, "expected no error when deleting clue")
}

func TestClueRepository_FindCluesByLocation(t *testing.T) {
	repo, cleanup := setupClueRepo(t)
	defer cleanup()
	ctx := context.Background()

	locationID := "location-1"
	clue1 := &models.Clue{
		InstanceID: "instance-1",
		LocationID: locationID,
		Content:    "Clue 1",
	}
	clue2 := &models.Clue{
		InstanceID: "instance-2",
		LocationID: locationID,
		Content:    "Clue 2",
	}

	// Save clues first
	err := repo.Save(ctx, clue1)
	assert.NoError(t, err, "expected no error when saving clue 1")
	err = repo.Save(ctx, clue2)
	assert.NoError(t, err, "expected no error when saving clue 2")

	clues, err := repo.FindCluesByLocation(ctx, locationID)
	assert.NoError(t, err, "expected no error when finding clues by location")
	assert.Len(t, clues, 2, "expected two clues to be found")
}
