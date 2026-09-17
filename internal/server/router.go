// Package server wires the Gin engine, mirroring the route-level access
// rules declared in Security/SecurityConfug.kt's securityFilterChain:
//
//	.authorizeHttpRequests {
//	    it.requestMatchers("/api/auth/**").permitAll()
//	        .anyRequest().authenticated()
//	}
package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/auth"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/interview"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/repository"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/security"
)

type Dependencies struct {
	AuthHandler      *auth.Handler
	InterviewHandler *interview.Handler
	JWTService       *security.JWTService
	UserRepository   *repository.UserRepository
}

// NewRouter builds the full Gin engine.
func NewRouter(deps Dependencies) *gin.Engine {
	r := gin.New()

	// gin.Logger() + gin.Recovery() give request logging and panic recovery
	// (-> 500) — the original relies on Spring Boot's embedded Tomcat access
	// logging (off by default, same as here unless enabled) and Spring's
	// default uncaught-exception -> 500 behavior respectively.
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// NOTE: the original SecurityConfig.kt never calls `.cors(...)` and
	// defines no CorsConfigurationSource bean, so Spring Security applies no
	// CORS headers at all — a browser-based frontend on a different origin
	// would be blocked by the browser's same-origin policy. That's very
	// likely an oversight rather than an intentional restriction for an
	// app meant to be called from a web frontend, so we add a permissive
	// CORS middleware here to make the Go service usable as-is. Tighten
	// AllowOrigins to your real frontend origin(s) before deploying.
	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")

	// "/api/auth/**".permitAll()
	authGroup := api.Group("/auth")
	deps.AuthHandler.RegisterRoutes(authGroup)

	// ".anyRequest().authenticated()" for everything else.
	interviewGroup := api.Group("/interview")
	interviewGroup.Use(security.RequireAuth(deps.JWTService, deps.UserRepository))
	deps.InterviewHandler.RegisterRoutes(interviewGroup)

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Max-Age", (12 * time.Hour).String())

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
