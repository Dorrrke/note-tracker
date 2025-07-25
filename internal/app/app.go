package app

import (
	"github.com/Dorrrke/note-tracker/internal/config"
	"github.com/Dorrrke/note-tracker/internal/server"
	"github.com/Dorrrke/note-tracker/internal/service"
)

// App represents the main application structure that combines
// configuration, API server and repository components.
type App struct {
	cfg  config.Config
	API  *server.API
	repo service.Repository
}

// NewApp creates a new instance of App with the provided configuration,
// API server and repository.
//
// Parameters:
//   - cfg: Application configuration
//   - server: API server instance
//   - repo: Repository implementation for data access
//
// Returns:
//   - *App: Pointer to the newly created App instance
func NewApp(cfg config.Config, server *server.API, repo service.Repository) *App {
	return &App{
		cfg:  cfg,
		API:  server,
		repo: repo,
	}
}

// StartApp initializes and starts the application server.
//
// This method starts the HTTP server and begins listening for incoming requests.
// It returns an error if the server fails to start.
//
// Returns:
//   - error: Server startup error, if any
func (app *App) StartApp() error {
	if err := app.API.Start(); err != nil {
		return err
	}
	return nil
}
