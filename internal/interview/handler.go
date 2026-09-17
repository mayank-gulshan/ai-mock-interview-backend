package interview

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/apperrors"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/dto"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/gemini"
)

// Handler mirrors Rest/GeminiController.kt's InterviewController
// (@RequestMapping("/api/interview")).
type Handler struct {
	geminiService     *gemini.Service
	evaluationService *EvaluationService
}

func NewHandler(geminiService *gemini.Service, evaluationService *EvaluationService) *Handler {
	return &Handler{geminiService: geminiService, evaluationService: evaluationService}
}

// RegisterRoutes mounts the same two endpoints as InterviewController.
// NOTE: SecurityConfig.kt requires authentication for anything outside
// /api/auth/**, so both routes here are protected by security.RequireAuth
// (see internal/server/router.go).
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/questions", h.getQuestions)
	rg.POST("/evaluate", h.evaluate)
}

// getQuestions mirrors:
//
//	@GetMapping("/questions")
//	fun getQuestions(
//	    @RequestParam role: String,
//	    @RequestParam(defaultValue = "medium") difficulty: String
//	): List<QuestionAnswer> = geminiService.generateInterviewQuestions(role, difficulty)
func (h *Handler) getQuestions(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		// NOTE: Spring's @RequestParam role: String (no default, not
		// nullable) makes `role` mandatory and returns 400 Bad Request if
		// missing. We replicate that instead of silently defaulting.
		c.JSON(http.StatusBadRequest, apperrors.ErrorResponse{Message: "Required request parameter 'role' is not present"})
		return
	}
	difficulty := c.DefaultQuery("difficulty", "medium")

	questions, err := h.geminiService.GenerateInterviewQuestions(role, difficulty)
	if err != nil {
		// NOTE: the Kotlin controller has no try/catch, so an exception here
		// (e.g. Gemini call failing, or the JSON not matching QuestionAnswer)
		// bubbles up as Spring Boot's default 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, apperrors.ErrorResponse{Message: "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// evaluate mirrors:
//
//	@PostMapping("/evaluate")
//	fun evaluate(@RequestBody request: EvaluationRequest): ResponseEntity<EvaluationResponse> {
//	    val response = evaluationService.evaluateAnswer(request)
//	    return ResponseEntity.ok(response)
//	}
func (h *Handler) evaluate(c *gin.Context) {
	var req dto.EvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperrors.ErrorResponse{Message: "Malformed request body"})
		return
	}

	resp := h.evaluationService.EvaluateAnswer(req)
	c.JSON(http.StatusOK, resp)
}
