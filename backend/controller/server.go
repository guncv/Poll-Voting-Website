package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guncv/Poll-Voting-Website/backend/config"
	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/service"
	"gorm.io/gorm"
)

// Server handles HTTP requests.
type Server struct {
	config             config.Config
	db                 *gorm.DB
	app                *fiber.App
	logger             log.LoggerInterface
	healthCheckService service.HealthCheckService
}

// NewServer creates a new Fiber server with injected dependencies.
func NewServer(cfg config.Config, db *gorm.DB) *Server {
	logger := log.Initialize(cfg.AppEnv)
	healthService := service.NewHealthCheckService()

	server := &Server{
		config:             cfg,
		db:                 db,
		app:                fiber.New(),
		logger:             logger,
		healthCheckService: healthService,
	}

	server.setupRoutes()

	return server
}

// setupRoutes defines all routes for the application.
func (s *Server) setupRoutes() {
	api := s.app.Group("/api")
	api.Get("/health", s.HealthCheck)
}

// Start runs the Fiber app.
func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}
