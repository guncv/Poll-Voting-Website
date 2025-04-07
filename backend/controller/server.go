package controller

import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/sns"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors" // <--- import the cors middleware
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
func NewServer(cfg config.Config, db *gorm.DB, cacheService db.CacheService) *Server {
    logger := log.Initialize(cfg.AppEnv)
    healthService := service.NewHealthCheckService()

    // Notification
    notificationClient := NewNotificationClient(cfg.Notification, logger)
    notificationRepo := repository.NewNotificationRepository(notificationClient, cfg, logger)
    notificationService := service.NewNotificationService(notificationRepo, logger)

    // User
    userRepo := repository.NewUserRepository(db, logger)
    userService := service.NewUserService(userRepo, logger, notificationService)

    // Question
    questionRepo := repository.NewQuestionRepository(db, logger)
    // IMPORTANT: pass cacheService to the question service here
    questionService := service.NewQuestionService(questionRepo, cacheService, logger)

    // Create Fiber instance
    app := fiber.New()

    // Enable CORS
    app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:3000", // or your frontend URL
        AllowCredentials: true,
    }))

    // Build the Server
    server := &Server{
        config:             cfg,
        db:                 db,
        cache:              cacheService,
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

    // General question routes
    q.Post("/", s.CreateQuestion)
    q.Get("/", s.GetAllQuestions)

    // Specific routes
    // General question routes
    q.Post("/", s.CreateQuestion)
    q.Get("/", s.GetAllQuestions)
    q.Post("/vote", s.VoteForQuestion)

    // Specific routes
    q.Get("/last", s.GetLastArchivedQuestion)

    // Parameterized routes
    q.Get("/:id", s.GetQuestion)
    q.Delete("/:id", s.DeleteQuestion)

    // Cache routes
    c := q.Group("/cache")
    c.Post("/", s.CreateQuestionCache)
    c.Get("/today", s.GetAllTodayQuestionIDs)
    c.Get("/:id", s.GetQuestionCache)
    c.Delete("/:id", s.DeleteQuestionCache)

    // ========================================
    // Cache test routes
    // ========================================
    cache := api.Group("/cache")
    cache.Get("/:key", s.getCache)
    cache.Post("/:key", s.setCache)
}


// Start runs the Fiber app.
func (s *Server) Start(address string) error {
    return s.app.Listen(address)
}
