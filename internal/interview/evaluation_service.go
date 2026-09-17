// Package interview mirrors Rest/GeminiController.kt (the InterviewController
// class) and Security/EvalutionService.kt (AnswerEvaluationService) — grouped
// together here since they form a single "interview" feature area.
package interview

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/dto"
	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/gemini"
)

// EvaluationService mirrors Security/EvalutionService.kt's AnswerEvaluationService.
type EvaluationService struct {
	geminiService *gemini.Service
}

func NewEvaluationService(geminiService *gemini.Service) *EvaluationService {
	return &EvaluationService{geminiService: geminiService}
}

// EvaluateAnswer mirrors `fun evaluateAnswer(request: EvaluationRequest): EvaluationResponse`.
func (s *EvaluationService) EvaluateAnswer(req dto.EvaluationRequest) dto.EvaluationResponse {
	prompt := s.buildEvaluationPrompt(req)
	raw, err := s.geminiService.GenerateEvaluation(prompt)
	if err != nil {
		// NOTE: the Kotlin version has no try/catch around geminiService call
		// itself (only around JSON parsing), so a Gemini/network failure
		// would propagate as an unhandled exception -> Spring's default 500.
		// We fall back to the same "could not evaluate" response the
		// original returns for a parse failure, which is a safe superset of
		// behavior for a failed upstream call and keeps this handler from
		// ever panicking.
		return fallbackEvaluationResponse()
	}
	return s.parseEvaluationJSON(raw)
}

// buildEvaluationPrompt mirrors the private `buildEvaluationPrompt`, with the
// exact same multi-line prompt template (Kotlin's trimIndent()).
func (s *EvaluationService) buildEvaluationPrompt(req dto.EvaluationRequest) string {
	return fmt.Sprintf(`You are an expert technical interviewer evaluating a candidate's answer.

Role: %s
Difficulty: %s
Question: %s
Candidate's Answer: %s

Evaluate the answer and respond ONLY with valid JSON in this exact format,
no markdown, no extra text:
{
  "score": <integer 0-100>,
  "feedback": "<2-3 sentence constructive feedback>",
  "idealAnswerSummary": "<1-2 sentence summary of what a strong answer would include>"
}`, req.Role, req.Difficulty, req.Question, req.Answer)
}

// parseEvaluationJson mirrors the private `parseEvaluationJson(raw): EvaluationResponse`.
func (s *EvaluationService) parseEvaluationJSON(raw string) dto.EvaluationResponse {
	cleaned := strings.TrimSpace(
		strings.ReplaceAll(strings.ReplaceAll(raw, "```json", ""), "```", ""),
	)

	var resp dto.EvaluationResponse
	if err := json.Unmarshal([]byte(cleaned), &resp); err != nil {
		return fallbackEvaluationResponse()
	}
	return resp
}

func fallbackEvaluationResponse() dto.EvaluationResponse {
	return dto.EvaluationResponse{
		Score:              0,
		Feedback:           "Could not evaluate answer, please try again.",
		IdealAnswerSummary: "",
	}
}
