package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guncv/Poll-Voting-Website/backend/util"
	"gorm.io/gorm"
)

// Server handles HTTP requests
type Server struct {
	config util.Config
	db     *gorm.DB
	app    *fiber.App
	logger util.LoggerInterface
}

// NewServer creates a new Fiber server with routes
func NewServer(config util.Config, db *gorm.DB) *Server {
	logger := util.Initialize(config.AppEnv)

	server := &Server{
		config: config,
		db:     db,
		app:    fiber.New(),
		logger: logger,
	}

	server.setupRoutes()

	return server
}

// setupRoutes defines all routes for the application
func (s *Server) setupRoutes() {
	api := s.app.Group("/api")

	api.Get("/health", s.HealthCheck)
}

// Start runs the Fiber app
func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}
