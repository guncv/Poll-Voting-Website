package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guncv/Poll-Voting-Website/backend/config"
	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
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
	userService        service.UserService
	questionService    service.QuestionService
}

// NewServer creates a new Fiber server with injected dependencies.
func NewServer(cfg config.Config, db *gorm.DB) *Server {
	logger := log.Initialize(cfg.AppEnv)
	healthService := service.NewHealthCheckService()

	userRepo := repository.NewUserRepository(db, logger)
	userService := service.NewUserService(userRepo, logger)

	questionRepo := repository.NewQuestionRepository(db, logger)
	questionService := service.NewQuestionService(questionRepo, logger)
	
	server := &Server{
		config:             cfg,
		db:                 db,
		app:                fiber.New(),
		logger:             logger,
		healthCheckService: healthService,
		userService:        userService,
		questionService:    questionService,
	}

	server.setupRoutes()

	return server
}

// setupRoutes defines all routes for the application.
func (s *Server) setupRoutes() {
	api := s.app.Group("/api")
	api.Get("/health", s.HealthCheck)

	user := api.Group("/user")
	user.Post("/register", s.Register)
	user.Post("/login", s.Login)
	user.Get("/profile", s.Profile)
	user.Get("/logout", s.Logout)
	user.Get("/:id", s.GetUser)
	user.Delete("/:id", s.DeleteUser)
	user.Put("/:id", s.UpdateUser)

	q := api.Group("/question")
	q.Post("/", s.CreateQuestion)
	q.Get("/", s.GetAllQuestions)
	q.Get("/:id", s.GetQuestion)
	q.Delete("/:id", s.DeleteQuestion)
}

// Start runs the Fiber app.
func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}
