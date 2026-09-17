package security

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/models"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/repository"
)

// ContextUserKey is the gin.Context key the authenticated *models.User is
// stored under, replacing Spring Security's SecurityContextHolder +
// UserPrincipal (see the NOTE in internal/models/user.go).
const ContextUserKey = "authenticatedUser"

// RequireAuth combines the behavior of three original Kotlin pieces into one
// Gin middleware, applied only to route groups that are NOT under
// /api/auth/** (which SecurityConfig.kt marks permitAll):
//
//  1. Security/JWTAuthFilter.kt: reads the "Authorization: Bearer <token>"
//     header, validates it's a well-formed *access* token, loads the user.
//  2. Security/SecurityConfug.kt: `.anyRequest().authenticated()` — any
//     request reaching here without a resolved user is rejected.
//  3. Security/CustonAuthEntryPoint.kt: on rejection, responds with
//     401 + JSON body {"message": "Unauthorized: <reason>"}, matching
//     `mapOf("message" to "Unauthorized: ${authException.message}")`.
func RequireAuth(jwtService *JWTService, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, reason := authenticate(c, jwtService, userRepo)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized: " + reason,
			})
			return
		}
		c.Set(ContextUserKey, user)
		c.Next()
	}
}

func authenticate(c *gin.Context, jwtService *JWTService, userRepo *repository.UserRepository) (*models.User, string) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		// NOTE: Spring's OncePerRequestFilter lets the request continue with
		// no authentication when the header is absent/malformed, and it's
		// the later `.anyRequest().authenticated()` check that ultimately
		// rejects it via the entry point. We collapse that two-step flow
		// into one middleware but produce the same end result: 401 with a
		// "Full authentication is required" style message.
		return nil, "Full authentication is required to access this resource"
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if !jwtService.ValidateToken(token) || !jwtService.IsAccessToken(token) {
		return nil, "Full authentication is required to access this resource"
	}

	userID, err := jwtService.ExtractUserID(token)
	if err != nil {
		return nil, "Full authentication is required to access this resource"
	}

	user, err := userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, "Full authentication is required to access this resource"
	}

	return user, ""
}

// CurrentUser retrieves the user set by RequireAuth, for handlers that need it.
func CurrentUser(c *gin.Context) *models.User {
	v, ok := c.Get(ContextUserKey)
	if !ok {
		return nil
	}
	u, ok := v.(*models.User)
	if !ok {
		return nil
	}
	return u
}
