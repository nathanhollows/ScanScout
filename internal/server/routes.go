package server

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nathanhollows/Rapua/filesystem"
	admin "github.com/nathanhollows/Rapua/internal/handlers/admin"
	players "github.com/nathanhollows/Rapua/internal/handlers/players"
	public "github.com/nathanhollows/Rapua/internal/handlers/public"
	"github.com/nathanhollows/Rapua/internal/middlewares"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/services"
)

func setupRouter(
	logger *slog.Logger,
	gameplayService services.GameplayService,
	gameManagerService services.GameManagerService,
	notificationService services.NotificationService,
) *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Compress(5))
	router.Use(middleware.CleanPath)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RedirectSlashes)

	// Create userServices for authentication
	userServices := services.NewUserServices(repositories.NewUserRepository())

	// Public routes
	publicHandler := public.NewPublicHandler(logger, *userServices)
	setupPublicRoutes(router, publicHandler)

	// Player routes
	playerHandler := players.NewPlayerHandler(logger, gameplayService, notificationService)
	setupPlayerRoutes(router, playerHandler)

	// Admin routes
	adminHandler := admin.NewAdminHandler(
		logger,
		gameManagerService,
		notificationService,
		*userServices)
	setupAdminRoutes(router, adminHandler)

	// Static files
	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "static"))}
	filesystem.FileServer(router, "/static", filesDir)

	return router
}

// Setup the player routes
func setupPlayerRoutes(router chi.Router, playerHandler *players.PlayerHandler) {
	teamRepo := repositories.NewTeamRepository()
	// Home route
	// Takes a GET request to show the home page
	// Takes a POST request to submit the home page form
	router.Get("/play", playerHandler.Play)
	router.Post("/play", playerHandler.PlayPost)

	// Show the next available locations
	router.Route("/next", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middlewares.TeamMiddleware(teamRepo,
				middlewares.LobbyMiddleware(teamRepo, next))
		})
		r.Get("/", playerHandler.Next)
		r.Post("/", playerHandler.Next)
	})

	router.Route("/blocks", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middlewares.TeamMiddleware(teamRepo,
				middlewares.LobbyMiddleware(teamRepo, next))
		})
		r.Post("/validate", playerHandler.ValidateBlock)
	})

	// Show the lobby page
	router.Route("/lobby", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middlewares.TeamMiddleware(teamRepo, next)
		})
		r.Get("/", playerHandler.Lobby)
	})

	// Check in to a location
	router.Route("/s", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middlewares.TeamMiddleware(teamRepo,
				middlewares.LobbyMiddleware(teamRepo, next))
		})
		r.Get("/{code:[A-z]{5}}", playerHandler.CheckIn)
		r.Post("/{code:[A-z]{5}}", playerHandler.CheckInPost)
	})

	// Check out of a location
	router.Route("/o", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middlewares.TeamMiddleware(teamRepo,
				middlewares.LobbyMiddleware(teamRepo, next))
		})
		r.Get("/", playerHandler.CheckOut)
		r.Get("/{code:[A-z]{5}}", playerHandler.CheckOut)
		r.Post("/{code:[A-z]{5}}", playerHandler.CheckOutPost)
	})

	router.Route("/checkins", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middlewares.TeamMiddleware(teamRepo,
				middlewares.LobbyMiddleware(teamRepo, next))
		})
		r.Get("/", playerHandler.MyCheckins)
		r.Get("/{id}", playerHandler.CheckInView)
	})

	router.Post("/dismiss/{ID}", playerHandler.DismissNotificationPost)

}

func setupPublicRoutes(router chi.Router, publicHandler *public.PublicHandler) {
	router.Get("/", publicHandler.Index)
	router.Get("/pricing", publicHandler.Pricing)
	router.Get("/about", publicHandler.About)
	router.Get("/contact", publicHandler.Contact)
	router.Get("/privacy", publicHandler.Privacy)

	router.Route("/login", func(r chi.Router) {
		r.Get("/", publicHandler.Login)
		r.Post("/", publicHandler.LoginPost)
	})
	router.Get("/logout", publicHandler.Logout)
	router.Route("/register", func(r chi.Router) {
		r.Get("/", publicHandler.Register)
		r.Post("/", publicHandler.RegisterPost)
	})
	router.Get("/forgot", publicHandler.ForgotPassword)
	router.Post("/forgot", publicHandler.ForgotPasswordPost)

	router.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", publicHandler.Auth)
		r.Get("/{provider}/callback", publicHandler.AuthCallback)
	})

	router.Route("/verify-email", func(r chi.Router) {
		r.Get("/", publicHandler.VerifyEmail)
		r.Get("/{token}", publicHandler.VerifyEmailWithToken)
		r.Get("/status", publicHandler.VerifyEmailStatus)
		r.Post("/resend", publicHandler.ResendEmailVerification)
	})

	router.NotFound(publicHandler.NotFound)

}

func setupAdminRoutes(router chi.Router, adminHandler *admin.AdminHandler) {
	router.Route("/admin", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middlewares.AdminAuthMiddleware(adminHandler.UserServices.AuthService, next)
		})
		r.Use(middlewares.AdminCheckInstanceMiddleware)

		r.Route("/quickstart", func(r chi.Router) {
			r.Get("/", adminHandler.Quickstart)
			r.Post("/dismiss", adminHandler.DismissQuickstart)
		})

		r.Get("/", adminHandler.Activity)
		r.Route("/activity", func(r chi.Router) {
			r.Get("/", adminHandler.Activity)
			r.Get("/teams", adminHandler.ActivityTeamsOverview)
			r.Get("/team/{teamCode}", adminHandler.TeamActivity)
		})

		r.Route("/locations", func(r chi.Router) {
			r.Get("/", adminHandler.Locations)
			r.Post("/reorder", adminHandler.ReorderLocations)
			r.Get("/new", adminHandler.LocationNew)
			r.Post("/new", adminHandler.LocationNewPost)
			r.Get("/{id}", adminHandler.LocationEdit)
			r.Post("/{id}", adminHandler.LocationEditPost)
			r.Delete("/{id}", adminHandler.LocationDelete)
			r.Get("/{id}/preview", adminHandler.LocationPreview)
			// Assets
			r.Get("/qr/{action}/{id}.{extension}", adminHandler.QRCode)
			r.Get("/qr-codes.zip", adminHandler.GenerateQRCodeArchive)
			r.Get("/poster/{id}.pdf", adminHandler.GeneratePoster)
			r.Get("/posters.pdf", adminHandler.GeneratePosters)
			// Blocks
			r.Route("/{location}/blocks", func(r chi.Router) {
				// r.Get("/", adminHandler.Blocks)
				// r.Post("/", adminHandler.BlocksPost)
				r.Post("/reorder", adminHandler.ReorderBlocks)
				r.Post("/new/{type}", adminHandler.BlockNewPost)
				r.Get("/{blockID}/edit", adminHandler.BlockEdit)
				r.Post("/{blockID}/update", adminHandler.BlockEditPost)
				r.Delete("/{blockID}/delete", adminHandler.BlockDelete)
			})
		})

		r.Route("/teams", func(r chi.Router) {
			r.Get("/", adminHandler.Teams)
			r.Post("/add", adminHandler.TeamsAdd)
		})

		r.Route("/experience", func(r chi.Router) {
			r.Get("/", adminHandler.Experience)
			r.Post("/", adminHandler.ExperiencePost)
			r.Post("/preview", adminHandler.ExperiencePreview)
		})

		r.Route("/instances", func(r chi.Router) {
			r.Get("/", adminHandler.Instances)
			r.Post("/new", adminHandler.InstancesCreate)
			r.Get("/{id}", adminHandler.Instances)
			r.Post("/{id}", adminHandler.Instances)
			r.Get("/{id}/switch", adminHandler.InstanceSwitch)
			r.Post("/delete", adminHandler.InstanceDelete)
			r.Post("/duplicate", adminHandler.InstanceDuplicate)
		})

		r.Route("/markdown", func(r chi.Router) {
			r.Get("/", adminHandler.MarkdownGuide)
			r.Post("/preview", adminHandler.PreviewMarkdown)
		})

		r.Route("/schedule", func(r chi.Router) {
			r.Get("/start", adminHandler.StartGame)
			r.Get("/stop", adminHandler.StopGame)
			r.Post("/", adminHandler.ScheduleGame)
		})

		r.Route("/notify", func(r chi.Router) {
			r.Post("/all", adminHandler.NotifyAllPost)
			r.Post("/team", adminHandler.NotifyTeamPost)
		})

		r.NotFound(adminHandler.NotFound)
	})
}
