// Package gemini mirrors Rest/GeminiService.kt and Rest/AppConfig.kt (the
// latter only existed to provide a shared RestTemplate bean; in Go we just
// keep a single *http.Client on the Service, so AppConfig.kt has no
// standalone file of its own).
package gemini

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mayank-gulshan/ai-mock-interview-backend-go/internal/dto"
)

// Service mirrors Rest/GeminiService.kt.
type Service struct {
	httpClient *http.Client
	apiKey     string
	apiURL     string
}

func NewService(apiKey, apiURL string) *Service {
	return &Service{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		apiKey:     apiKey,
		apiURL:     apiURL,
	}
}

// GenerateInterviewQuestions mirrors:
//
//	fun generateInterviewQuestions(role: String, difficulty: String): List<QuestionAnswer>
func (s *Service) GenerateInterviewQuestions(role, difficulty string) ([]dto.QuestionAnswer, error) {
	prompt := fmt.Sprintf(
		"Generate exactly 10 technical interview questions with answers\nfor a %s level %s candidate.\nKeep each answer to 2-3 clear, concise sentences.",
		difficulty, role,
	)

	schema := dto.NewResponseSchema(dto.NewSchemaItem(
		map[string]dto.SchemaProperty{
			"question": {Type: "STRING"},
			"answer":   {Type: "STRING"},
		},
		[]string{"question", "answer"},
	))

	requestBody := dto.GeminiRequest{
		Contents: []dto.GeminiContent{
			{Parts: []dto.GeminiPart{{Text: prompt}}},
		},
		GenerationConfig: dto.NewGenerationConfig(schema),
	}

	jsonText, err := s.callGemini(requestBody)
	if err != nil {
		return nil, err
	}

	var questions []dto.QuestionAnswer
	if err := json.Unmarshal([]byte(jsonText), &questions); err != nil {
		return nil, err
	}
	return questions, nil
}

// GenerateEvaluation mirrors:
//
//	fun generateEvaluation(prompt: String): String
//
// Generic method for evaluation — no fixed schema, just returns whatever raw
// text Gemini generates so the caller can parse it into whatever shape it
// needs (e.g. score/feedback JSON).
func (s *Service) GenerateEvaluation(prompt string) (string, error) {
	requestBody := dto.GeminiRequest{
		Contents: []dto.GeminiContent{
			{Parts: []dto.GeminiPart{{Text: prompt}}},
		},
		// No responseSchema here — evaluation prompt itself instructs
		// Gemini to return plain JSON in the exact format we want.
	}

	return s.callGemini(requestBody)
}

// callGemini mirrors the private `callGemini(requestBody): String`.
func (s *Service) callGemini(requestBody dto.GeminiRequest) (string, error) {
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, s.apiURL, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp dto.GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		// mirrors `?: throw IllegalStateException("Empty response from Gemini API")`
		return "", errors.New("empty response from Gemini API")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
