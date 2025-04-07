package controller

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gofiber/fiber/v2"
	"github.com/guncv/Poll-Voting-Website/backend/config"
	"github.com/guncv/Poll-Voting-Website/backend/db"
	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
	"github.com/guncv/Poll-Voting-Website/backend/service"
	"gorm.io/gorm"
)

// Server handles HTTP requests.
type Server struct {
	config             config.Config
	db                 *gorm.DB
	cache              db.CacheService
	app                *fiber.App
	logger             log.LoggerInterface
	healthCheckService service.HealthCheckService
	userService        service.UserService
	questionService    service.IQuestionService
}

func NewNotificationClient(cfg config.NotificationConfig, log log.LoggerInterface) *sns.Client {

	customCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		cfg.AccessKey,
		cfg.SecretKey,
		cfg.SessionToken,
	))

	return sns.New(sns.Options{
		Credentials: customCreds,
		Region:      cfg.Region,
	})
}

// NewServer creates a new Fiber server with injected dependencies.
func NewServer(cfg config.Config, db *gorm.DB, cacheService db.CacheService,) *Server {
	logger := log.Initialize(cfg.AppEnv)
	healthService := service.NewHealthCheckService()

	notificationClient := NewNotificationClient(cfg.Notification, logger)
	notificationRepo := repository.NewNotificationRepository(notificationClient, cfg, logger)
	notificationService := service.NewNotificationService(notificationRepo, logger)

	userRepo := repository.NewUserRepository(db, logger)
	userService := service.NewUserService(userRepo, logger, notificationService)

	questionRepo := repository.NewQuestionRepository(db, logger)
	questionService := service.NewQuestionService(questionRepo, cacheService, logger)

	server := &Server{
		config:             cfg,
		db:                 db,
		cache:              cacheService,
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

	// ========================================
	// Cache test routes
	// ========================================
	cache := api.Group("/cache")
	cache.Get("/:key", s.getCache)
	cache.Post("/:key", s.setCache)

	// ========================================
	// User routes
	// ========================================
	user := api.Group("/user")
	user.Post("/register", s.Register)
	user.Post("/login", s.Login)

	user.Use(JWTMiddleware)

	// Static
	user.Get("/profile", s.Profile)
	user.Get("/logout", s.Logout)

	// Dynamic
	user.Get("/:id", s.GetUser)
	user.Delete("/:id", s.DeleteUser)
	user.Put("/:id", s.UpdateUser)

	// ========================================
	// Question routes
	// ========================================
	q := api.Group("/question")
	q.Use(JWTMiddleware)

	// ✅ Specific routes first
	c := q.Group("/cache")
	c.Post("/", s.CreateQuestionCache)
	c.Get("/", s.GetAllTodayQuestionIDs)
	c.Get("/:id", s.GetQuestionCache)
	c.Delete("/:id", s.DeleteQuestionCache)

	// ✅ Then general question routes
	q.Post("/", s.CreateQuestion)
	q.Get("/", s.GetAllQuestions)

	// ⚠️ Parameterized routes last
	q.Get("/:id", s.GetQuestion)
	q.Delete("/:id", s.DeleteQuestion)

	q.Post("/vote", s.VoteForQuestion)
}


// Start runs the Fiber app.
func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}
