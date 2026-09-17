package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/apperrors"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/dto"
)

// Handler mirrors Security/AuthController.kt (@RequestMapping("/api/auth")).
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes mounts the same three endpoints as AuthController, all
// permitAll per SecurityConfig.kt (`"/api/auth/**".permitAll()`).
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.register)
	rg.POST("/login", h.login)
	rg.POST("/refresh", h.refresh)
}

// register mirrors:
//
//	@PostMapping("/register")
//	fun register(@RequestBody request: RegisterRequest): ResponseEntity<AuthResponse> =
//	    ResponseEntity.ok(authService.register(request))
func (h *Handler) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperrors.ErrorResponse{Message: "Malformed request body"})
		return
	}

	resp, err := h.service.Register(req)
	if err != nil {
		status, body := apperrors.StatusFor(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// login mirrors the /login endpoint.
func (h *Handler) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperrors.ErrorResponse{Message: "Malformed request body"})
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		status, body := apperrors.StatusFor(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// refresh mirrors the /refresh endpoint.
func (h *Handler) refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperrors.ErrorResponse{Message: "Malformed request body"})
		return
	}

	resp, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		status, body := apperrors.StatusFor(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, resp)
}
