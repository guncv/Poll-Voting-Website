package controller

import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/sns"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors" // <--- import the cors middleware
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
func NewServer(cfg config.Config, db *gorm.DB) *Server {
    logger := log.Initialize(cfg.AppEnv)
    healthService := service.NewHealthCheckService()

    notificationClient := NewNotificationClient(cfg.Notification, logger)
    notificationRepo := repository.NewNotificationRepository(notificationClient, cfg, logger)
    notificationService := service.NewNotificationService(notificationRepo, logger)

    userRepo := repository.NewUserRepository(db, logger)
    userService := service.NewUserService(userRepo, logger, notificationService)

    questionRepo := repository.NewQuestionRepository(db, logger)
    questionService := service.NewQuestionService(questionRepo, logger)

    // Create Fiber instance
    app := fiber.New()

    // Enable CORS
    app.Use(cors.New(cors.Config{
        // Change http://localhost:3000 to wherever your frontend is running
        AllowOrigins:     "http://localhost:3000",
        AllowCredentials: true,
    }))

    server := &Server{
        config:             cfg,
        db:                 db,
        app:                app,
        logger:             logger,
        healthCheckService: healthService,
        userService:        userService,
        questionService:    questionService,
    }

    // Set up routes on the fiber app
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
    user.Get("/logout", s.Logout)
    // Apply JWT middleware to protected routes.
    user.Use(JWTMiddleware)

    // Static
    user.Get("/profile", s.Profile)


    // Dynamic
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
