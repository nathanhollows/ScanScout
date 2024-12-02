//go:generate npm run build

package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/migrate"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/server"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	if err := godotenv.Load(".env"); err != nil {
		logger.Warn("could not load .env file", "error", err)
	}

	initialiseFolders(logger)

	db := db.MustOpen()
	migrate.CreateTables(logger, db)

	// Initialise the repos
	blockStateRepo := repositories.NewBlockStateRepository(db)
	blockRepo := repositories.NewBlockRepository(db, blockStateRepo)
	checkInRepo := repositories.NewCheckInRepository(db)
	clueRepo := repositories.NewClueRepository(db)
	instanceRepo := repositories.NewInstanceRepository(db)
	instanceSettingsRepo := repositories.NewInstanceSettingsRepository(db)
	locationRepo := repositories.NewLocationRepository(db)
	markerRepo := repositories.NewMarkerRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Initialise the services
	assetGenerator := services.NewAssetGenerator()
	authService := services.NewAuthService(userRepo)
	blockService := services.NewBlockService(blockRepo, blockStateRepo)
	checkInService := services.NewCheckInService(checkInRepo, locationRepo, teamRepo)
	clueService := services.NewClueService(clueRepo, locationRepo)
	emailService := services.NewEmailService()
	locationService := services.NewLocationService(clueRepo, locationRepo, markerRepo, blockRepo)
	navigationService := services.NewNavigationService()
	notificationService := services.NewNotificationService(notificationRepo, teamRepo)
	teamService := services.NewTeamService(teamRepo)
	userService := services.NewUserService(userRepo)
	gameplayService := services.NewGameplayService(
		checkInService,
		locationService,
		teamService,
		blockService,
		navigationService,
		markerRepo,
	)
	gameManagerService := services.NewGameManagerService(
		locationService,
		userService,
		teamService,
		markerRepo,
		clueRepo,
		instanceRepo,
		instanceSettingsRepo,
	)

	sessions.Start()
	server.Start(
		logger,
		assetGenerator,
		authService,
		blockService,
		checkInService,
		clueService,
		emailService,
		gameManagerService,
		gameplayService,
		locationService,
		navigationService,
		notificationService,
		teamService,
		userService,
	)
}

func initialiseFolders(logger *slog.Logger) {
	// TODO: Make this configurable
	folders := []string{
		"assets/",
		"assets/codes/",
		"assets/codes/png/",
		"assets/codes/svg/",
		"assets/fonts/",
		"assets/posters/"}

	for _, folder := range folders {
		_, err := os.Stat(folder)
		if err != nil {
			// Attempt to create the directory
			err = os.MkdirAll(folder, 0755)
			if err != nil {
				logger.Error("could not create directory", "folder", folder, "error", err)
				os.Exit(1)
			}
		}
	}
}
