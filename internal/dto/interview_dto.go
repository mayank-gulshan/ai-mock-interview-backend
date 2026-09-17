package dto

// The types below mirror Dto/InterviewDto.kt. NOTE: as in the original
// Kotlin project, none of these are currently referenced by any handler —
// see the NOTE in internal/models/session.go for context. Kept for 1:1
// source parity.

type StartSessionRequest struct {
	Role       string `json:"role"`
	Difficulty string `json:"difficulty"`
}

type StartSessionResponse struct {
	SessionID string `json:"sessionId"`
	Question  string `json:"question"`
}

type SubmitAnswerRequest struct {
	SessionID    string `json:"sessionId"`
	QuestionText string `json:"questionText"`
	UserAnswer   string `json:"userAnswer"`
}

type AIResponse struct {
	Question      string  `json:"question"`
	Rating        float64 `json:"rating"`
	Feedback      string  `json:"feedback"`
	IdealAnswer   string  `json:"idealAnswer"`
	AccuracyScore int     `json:"accuracyScore"`
	DepthScore    int     `json:"depthScore"`
	ClarityScore  int     `json:"clarityScore"`
}
