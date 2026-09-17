// Command server is the entrypoint, mirroring
// AiMockInterviewApplication.kt's `fun main(args: Array<String>) { runApplication<...>(*args) }`.
package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/auth"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/config"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/db"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/gemini"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/interview"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/models"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/repository"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/security"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/server"
)

func main() {
	// Optional: load a local .env file for development convenience. In
	// production the same env vars (DB_PASS, JWT_SECRET_64, GEMINI_API_KEY,
	// PORT, ...) can be set directly, exactly as with the Spring Boot app.
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	gormDB, err := db.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Equivalent to spring.jpa.hibernate.ddl-auto=update.
	if err := gormDB.AutoMigrate(&models.User{}, &models.Session{}, &models.Answer{}); err != nil {
		log.Fatalf("failed to auto-migrate schema: %v", err)
	}

	// --- Repositories (Repository/*.kt) ---
	userRepo := repository.NewUserRepository(gormDB)
	_ = repository.NewSessionRepository(gormDB) // kept for parity; unused, see NOTE in internal/models/session.go
	_ = repository.NewAnswerRepository(gormDB)  // kept for parity; unused, see NOTE in internal/models/session.go

	// --- Security (Security/*.kt) ---
	hashEncoder := security.NewHashEncoder()
	jwtService, err := security.NewJWTService(cfg.JWTSecretBase64)
	if err != nil {
		log.Fatalf("invalid JWT_SECRET_64: %v", err)
	}

	// --- Auth feature (Security/AuthService.kt + AuthController.kt) ---
	authService := auth.NewService(userRepo, hashEncoder, jwtService)
	authHandler := auth.NewHandler(authService)

	// --- Gemini + interview feature (Rest/GeminiService.kt,
	//     Security/EvalutionService.kt, Rest/GeminiController.kt) ---
	geminiService := gemini.NewService(cfg.GeminiAPIKey, cfg.GeminiURL)
	evaluationService := interview.NewEvaluationService(geminiService)
	interviewHandler := interview.NewHandler(geminiService, evaluationService)

	router := server.NewRouter(server.Dependencies{
		AuthHandler:      authHandler,
		InterviewHandler: interviewHandler,
		JWTService:       jwtService,
		UserRepository:   userRepo,
	})

	// server.address=0.0.0.0 + server.port=${PORT:8080}: Gin's Run(":"+port)
	// listens on all interfaces by default, matching 0.0.0.0.
	addr := ":" + cfg.Port
	log.Printf("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
